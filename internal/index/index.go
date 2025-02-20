package index

type IndexType int

const (
	HASH IndexType = iota
	BTREE
	BPTREE
	SKIPLIST
	LSM
)

type Indexer interface {
	Get(key []byte) (*Entry, bool)
	Put(key []byte, value *Entry)
	Del(key []byte) (*Entry, bool)

	// Scan(start int, end int) <-chan Entry
	// Snapshot() map[string]*Entry

	// Iterator return Iterator
	Iterator() Iterator
	// Has(key []byte) bool
	// Keys() [][]byte

	Size() int
	Close() error
}
type Iterator interface {
	Rewind()
	Next()
	Key() []byte
	Valid() bool
	Value() *Entry
	Release()
}

// Entry 索引记录结构
type Entry struct {
	FileID    uint32 // 文件ID
	Offset    int64  // 偏移量
	ValueSize uint32 // 值大小
	Timestamp uint64 // 时间戳
}

// New 创建内存索引
func New(indexType IndexType) Indexer {
	switch indexType {
	case HASH:
		return NewShardMap(defaultShardCount, defaultHasher)
	case BTREE:
		return NewBTree(32)
	default:
		return NewShardMap(defaultShardCount, defaultHasher)
	}
}
