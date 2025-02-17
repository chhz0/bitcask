package index

import (
	"hash/fnv"
	"sync"
	"sync/atomic"

	"github.com/chhz0/caskv/internal/utils"
)

// - 8核CPU：16-64分片
// - 16核CPU：64-256分片
const defaultShardCount = 32

// 分片结构
type shardEntry struct {
	sync.RWMutex
	entries map[string]*Entry
	size    int64
	
}

// ShardMap  并发安全的分片map
type ShardMap struct {
	shareds []*shardEntry       // 分片数组
	hasher  func(string) uint32 // 哈希函数
	size    atomic.Int64
}

// Get implements Indexer.
func (sm *ShardMap) Get(key []byte) (*Entry, bool) {
	k := utils.BytesToString(key)
	shard := sm.getShard(k)
	shard.RLock()
	defer shard.RUnlock()

	entry, ok := shard.entries[k]
	if !ok {
		return nil, false
	}

	return entry, true
}

// Put implements Indexer.
func (sm *ShardMap) Put(key []byte, value *Entry) {
	k := utils.BytesToString(key)
	shard := sm.getShard(k)
	shard.Lock()
	defer shard.Unlock()

	shard.entries[k] = value
	shard.size++
}

// Delete implements Indexer.
func (sm *ShardMap) Del(key []byte) (*Entry, bool) {
	k := utils.BytesToString(key)
	shard := sm.getShard(k)
	shard.Lock()
	defer shard.Unlock()

	entry, ok := shard.entries[k]
	if !ok {
		return nil, false
	}

	delete(shard.entries, k)
	shard.size--
	return entry, true
}

// Len implements Indexer.
func (sm *ShardMap) Size() int {
	if len(sm.shareds) == 0 {
		return 0
	}

	var totalSize int64
	shardSize := make(chan int64, len(sm.shareds))

	var wg sync.WaitGroup
	wg.Add(len(sm.shareds))

	for _, shard := range sm.shareds {
		go func(s *shardEntry) {
			defer wg.Done()

			s.RLock()
			defer s.RUnlock()
			shardSize <- s.size
		}(shard)
	}

	go func() {
		wg.Wait()
		close(shardSize)
	}()

	for size := range shardSize {
		totalSize += size
	}

	return int(totalSize)
}

// Scan implements Indexer.
func (sm *ShardMap) Scan(start int, end int) <-chan Entry {
	panic("unimplemented")
}

// Snapshot implements Indexer.
func (sm *ShardMap) Snapshot() map[string]*Entry {
	snap := make(map[string]*Entry, sm.size.Load())
	sm.iter(func(key string, value *Entry) bool {
		snap[key] = value
		return true
	})
	return snap
}

func (sm *ShardMap) iter(fn func(key string, value *Entry) bool) {
	for _, shard := range sm.shareds {
		shard.RLock()
		items := make(map[string]*Entry, len(shard.entries))
		for k, v := range shard.entries {
			items[k] = v
		}
		shard.RUnlock()

		// 处理快照数据
		for k, v := range items {
			if !fn(k, v) {
				return
			}
		}
	}
}

// Close implements Indexer.
func (sm *ShardMap) Close() error {
	for _, shard := range sm.shareds {
		shard.Lock()
		shard.entries = make(map[string]*Entry)
		shard.Unlock()
	}
	sm.size.Store(0)
	return nil
}

func (sm *ShardMap) Keys() []string {
	keys := make([]string, 0, sm.size.Load())
	mutex := &sync.Mutex{}

	var wg sync.WaitGroup
	wg.Add(len(sm.shareds))

	for _, shard := range sm.shareds {
		go func(s *shardEntry) {
			defer wg.Done()

			s.RLock()
			defer s.RUnlock()
			for k := range s.entries {
				mutex.Lock()
				keys = append(keys, k)
				mutex.Unlock()
			}
		}(shard)
	}

	wg.Wait()
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
	return sm.shareds[hashed%uint32(len(sm.shareds))]
}

func NewShardMap(shardCount int, hasher func(string) uint32) Indexer {
	shards := make([]*shardEntry, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shardEntry{entries: make(map[string]*Entry)}
	}

	return &ShardMap{
		shareds: shards,
		hasher:  hasher,
	}
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
	shard.RLock()
	defer shard.RUnlock()
	return shard.entries[sit.keys[sit.cursor]]
}

// Close implements Iterator.
func (sit *shardIterator) Close() {
	sit.keys = nil
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
	shard.RLock()
	defer shard.RUnlock()

	sit.keys = make([]string, 0, len(shard.entries))
	for k := range shard.entries {
		sit.keys = append(sit.keys, k)
	}
}

func defaultHasher(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
