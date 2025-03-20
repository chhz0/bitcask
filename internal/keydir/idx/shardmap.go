package idx

import (
	"hash/maphash"
	"sync"
)

// - 8核CPU：16-64分片
// - 16核CPU：64-256分片
const defaultShardCount = 64

// 分片结构
type shard struct {
	rw      sync.RWMutex
	entries map[string]*Index
	msize   int64
}

func (s *shard) load(key string) (*Index, bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()

	entry, ok := s.entries[key]
	return entry, ok
}

func (s *shard) store(key string, value *Index) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if _, ok := s.entries[key]; !ok {
		s.msize++
	}

	s.entries[key] = value
}

func (s *shard) delete(key string) (*Index, bool) {
	s.rw.Lock()
	defer s.rw.Unlock()

	if entry, ok := s.entries[key]; ok {
		delete(s.entries, key)
		s.msize--
		return entry, true
	}

	return nil, false
}

func (s *shard) size() int64 {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.msize
}

func (s *shard) forange(fn func(key string, value *Index) bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	for k, v := range s.entries {
		if !fn(k, v) {
			return
		}
	}
}

func (s *shard) clear() {
	s.rw.Lock()
	defer s.rw.Unlock()

	s.entries = make(map[string]*Index, 0)
	s.msize = 0
}

// ShardMap  并发安全的分片map
type ShardMap struct {
	shareds    []*shard            // 分片数组
	shardCount int                 // 分片数量
	hasher     func(string) uint64 // 哈希函数
}

func NewShardMap(shardCount int, hasher func(string) uint64) *ShardMap {
	shards := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shard{
			entries: make(map[string]*Index, 1024),
		}
	}

	return &ShardMap{
		shareds:    shards,
		shardCount: shardCount,
		hasher:     hasher,
	}
}

// Get implements Indexer.
func (sm *ShardMap) Get(key []byte) (*Index, bool) {
	k := string(key)
	shard := sm.getShard(k)

	return shard.load(k)
}

// Put implements Indexer.
func (sm *ShardMap) Put(key []byte, value *Index) {
	k := string(key)
	shard := sm.getShard(k)

	shard.store(k, value)
}

// Del implements Indexer.
func (sm *ShardMap) Del(key []byte) (*Index, bool) {
	k := string(key)
	shard := sm.getShard(k)
	return shard.delete(k)
}

// Size Len implements Indexer.
func (sm *ShardMap) Size() int64 {
	if sm.shardCount == 0 {
		return 0
	}

	var totalSize int64
	for _, shard := range sm.shareds {
		totalSize += shard.size()
	}

	return totalSize
}

// Close implements Indexer.
func (sm *ShardMap) Close() error {
	for _, shard := range sm.shareds {
		shard.clear()
	}
	return nil
}

func (sm *ShardMap) Keys() [][]byte {
	keys := make([][]byte, 0, sm.Size())
	for _, shard := range sm.shareds {
		shard.forange(func(key string, value *Index) bool {
			keys = append(keys, []byte(key))
			return true
		})
	}
	return keys
}

func (sm *ShardMap) getShard(key string) *shard {
	hashed := sm.hasher(key)
	return sm.shareds[hashed%uint64(sm.shardCount)]
}

var seed = maphash.MakeSeed()

func maphashFn(key string) uint64 {
	return maphash.String(seed, key)
}

// Iterator implements Indexer.
// func (sm *ShardMap) Iterator() Iterator {
// 	return &shardIterator{
// 		sm:       sm,
// 		shardIdx: 0,
// 		keys:     sm.Keys(),
// 		cursor:   -1,
// 	}
// }
