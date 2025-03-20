package idx

type IndexType int

const (
	HASH IndexType = iota
	BTREE
)

type Indexer interface {
	Get(key []byte) (idx *Index, b bool)
	Put(key []byte, value *Index)
	Del(key []byte) (old *Index, b bool)
	Keys() [][]byte
	// Iterator() Iterator
	// Size() int
	// Close() error
}

// Index is the key directory entry in memory.
type Index struct {
	FileID  uint32 // file id
	ValPos  int64  // value_pos in file
	ValSize uint64 // value_size of vlaue
}

// IndexWithTimeout is the key directory entry with expiration time in memory.
type IndexWithTimeout struct {
	Index
	Tstamp uint64 // timestamp
	Exp    uint64 // expiration time
}

// New 创建内存索引
func New(indexType IndexType) Indexer {
	switch indexType {
	// case BTREE:
	// 	return NewBTree(32)
	default:
		return NewShardMap(defaultShardCount, maphashFn)
	}
}
