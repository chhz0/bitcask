package index

import "github.com/chhz0/go-bitcask/internal/entry"

type IndexType int

const (
	HASH IndexType = iota
	BTREE
	// BPTREE
	// SKIPLIST
	// LSM
)

type Indexer interface {
	Get(key []byte) (idx *entry.Index, b bool)
	Put(key []byte, value *entry.Index)
	Del(key []byte) (old *entry.Index, b bool)
	Iterator() Iterator
	Size() int
	Close() error

	// Scan(start int, end int) <-chan Entry
	// Snapshot() map[string]*Entry
	// Has(key []byte) bool
	// Keys() [][]byte
}

type Iterator interface {
	Rewind()
	Next()
	Key() []byte
	Valid() bool
	Value() *entry.Index
	Release()
}

// New 创建内存索引
func New(indexType IndexType) Indexer {
	switch indexType {
	case HASH:
		return NewShardMap(defaultShardCount, fnv32a)
	// case BTREE:
	// 	return NewBTree(32)
	default:
		return NewShardMap(defaultShardCount, defaultHasher)
	}
}
