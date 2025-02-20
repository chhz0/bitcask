package index

import (
	"sync"

	"github.com/cespare/xxhash/v2"
	"github.com/chhz0/go-bitcask/internal/utils/zerocopy"
)

// - 8核CPU：16-64分片
// - 16核CPU：64-256分片
const defaultShardCount = 32

// 分片结构
type shardEntry struct {
	rw      sync.RWMutex
	entries map[string]*Entry
	size    int64
}

func (se *shardEntry) load(key string) (*Entry, bool) {
	se.rw.RLock()
	defer se.rw.RUnlock()
	entry, ok := se.entries[key]
	return entry, ok
}

func (se *shardEntry) store(key string, value *Entry) {
	se.rw.Lock()
	defer se.rw.Unlock()

	se.entries[key] = value
	se.size++
}

func (se *shardEntry) delete(key string) (*Entry, bool) {
	se.rw.Lock()
	defer se.rw.Unlock()

	entry, ok := se.entries[key]
	if !ok {
		return nil, ok
	}

	delete(se.entries, key)
	se.size--
	return entry, ok
}

func (se *shardEntry) ssize() int64 {
	se.rw.RLock()
	defer se.rw.RUnlock()
	return se.size
}

func (se *shardEntry) forange(fn func(key string, value *Entry) bool) {
	se.rw.RLock()
	defer se.rw.RUnlock()
	for k, v := range se.entries {
		if !fn(k, v) {
			return
		}
	}
}

func (se *shardEntry) clear() {
	se.rw.Lock()
	defer se.rw.Unlock()

	se.entries = make(map[string]*Entry, 1024)
	se.size = 0
}

// ShardMap  并发安全的分片map
type ShardMap struct {
	shareds    []*shardEntry       // 分片数组
	shardCount int                 // 分片数量
	hasher     func(string) uint32 // 哈希函数
}

func NewShardMap(shardCount int, hasher func(string) uint32) Indexer {
	shards := make([]*shardEntry, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shardEntry{entries: make(map[string]*Entry, 1024)}
	}

	return &ShardMap{
		shareds:    shards,
		shardCount: shardCount,
		hasher:     hasher,
	}
}

// Get implements Indexer.
func (sm *ShardMap) Get(key []byte) (*Entry, bool) {
	k := zerocopy.BytesToString(key)
	shard := sm.getShard(k)

	return shard.load(k)
}

// Put implements Indexer.
func (sm *ShardMap) Put(key []byte, value *Entry) {
	k := zerocopy.BytesToString(key)
	shard := sm.getShard(k)

	shard.store(k, value)
}

// Del  implements Indexer.
func (sm *ShardMap) Del(key []byte) (*Entry, bool) {
	k := zerocopy.BytesToString(key)
	shard := sm.getShard(k)
	return shard.delete(k)
}

// Size Len implements Indexer.
func (sm *ShardMap) Size() int {
	if sm.shardCount == 0 {
		return 0
	}

	var totalSize int64
	for _, shard := range sm.shareds {
		totalSize += shard.ssize()
	}

	return zerocopy.I64ToInt(totalSize)
}

// Scan implements Indexer.
func (sm *ShardMap) Scan(start int, end int) <-chan Entry {
	panic("unimplemented")
}

// Snapshot implements Indexer.
func (sm *ShardMap) Snapshot() map[string]*Entry {
	snap := make(map[string]*Entry, sm.Size())
	sm.iter(func(key string, value *Entry) bool {
		snap[key] = value
		return true
	})
	return snap
}

func (sm *ShardMap) iter(fn func(key string, value *Entry) bool) {
	for _, shard := range sm.shareds {
		shard.forange(fn)
	}
}

// Close implements Indexer.
func (sm *ShardMap) Close() error {
	for _, shard := range sm.shareds {
		shard.clear()
	}
	return nil
}

func (sm *ShardMap) Keys() []string {
	keys := make([]string, 0, sm.Size())
	for _, shard := range sm.shareds {
		shard.forange(func(key string, value *Entry) bool {
			keys = append(keys, key)
			return true
		})
	}
	return keys
}

// Reshard 动态扩容 todo
func (sm *ShardMap) Reshard(newCount int) {
	// 1. 全局锁保护重建过程
	// 2. 创建新的分片集合
	// 3. 数据迁移
	// 4. 原子替换分片引用
}

func (sm *ShardMap) getShard(key string) *shardEntry {
	hashed := sm.hasher(key)
	return sm.shareds[hashed%uint32(sm.shardCount)]
}

// Iterator implements Indexer.
func (sm *ShardMap) Iterator() Iterator {
	return &shardIterator{
		sm:       sm,
		shardIdx: 0,
		keys:     sm.Keys(),
		cursor:   -1,
	}
}

type shardIterator struct {
	sm       *ShardMap
	shardIdx int
	keys     []string
	cursor   int
	sync.RWMutex
}

// Rewind implements Iterator.
func (sit *shardIterator) Rewind() {
	sit.Lock()
	defer sit.Unlock()
	sit.shardIdx = 0
	sit.cursor = -1
	sit.loadKeys()
}

// Valid implements Iterator.
func (sit *shardIterator) Valid() bool {
	return sit.cursor < len(sit.keys)-1 && sit.shardIdx < len(sit.sm.shareds)-1
}

// Next implements Iterator.
func (sit *shardIterator) Next() {
	sit.Lock()
	defer sit.Unlock()

	if sit.cursor < len(sit.keys)-1 {
		sit.cursor++
	} else {
		sit.shardIdx++
		sit.cursor = 0
		sit.loadKeys()
	}
}

// Value implements Iterator.
func (sit *shardIterator) Value() *Entry {
	shard := sit.sm.shareds[sit.shardIdx]

	entry, ok := shard.load(sit.keys[sit.cursor])
	if !ok {
		return nil
	}
	return entry
}

// Release implements Iterator.
func (sit *shardIterator) Release() {
	sit.sm = nil
	sit.shardIdx = 0
	sit.keys = nil
	sit.cursor = -1
}

// Key implements Iterator.
func (sit *shardIterator) Key() []byte {
	if sit.cursor >= 0 && sit.cursor < len(sit.keys) {
		return []byte(sit.keys[sit.cursor])
	}
	return nil
}

func (sit *shardIterator) loadKeys() {
	shard := sit.sm.shareds[sit.shardIdx]

	sit.keys = make([]string, 0, shard.ssize())
	shard.forange(func(key string, value *Entry) bool {
		sit.keys = append(sit.keys, key)
		return true
	})
}

func defaultHasher(key string) uint32 {
	return zerocopy.U64ToU32(xxhash.Sum64String(key))
}

func xxxHasher(key string) uint32 {
	return zerocopy.U64ToU32(xxhash.Sum64String(key))
}
