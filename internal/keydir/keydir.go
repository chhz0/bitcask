package keydir

import (
	"github.com/chhz0/go-bitcask/internal/keydir/index"
)

// KeyDir 管理内存索引
type KeyDir struct {
	index.Indexer
	// iter.Iterator
}

func NewKeyDir(indexType string) *KeyDir {
	var indexer index.Indexer
	switch indexType {
	case "btree":
		indexer = index.New(index.BTREE)
	default:
		indexer = index.New(index.SHARDHASH)
	}

	return &KeyDir{indexer}
}
