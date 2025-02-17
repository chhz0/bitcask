package index

type IndexType int

const (
	SHARDMAP IndexType = iota
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
	Snapshot() map[string]*Entry
	Iterator() Iterator
	Size() int
	Close() error
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
	case SHARDMAP:
		return NewShardMap(defaultShardCount, defaultHasher)
	case BTREE:
		return NewBtree()
	default:
		return NewShardMap(defaultShardCount, defaultHasher)
	}
}
