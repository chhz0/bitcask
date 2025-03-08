package internal

import "github.com/chhz0/go-bitcask/internal/index"

// KeyDir 管理内存索引
type KeyDir struct {
	idx index.Indexer
}

func NewKeyDir(indexType index.IndexType) *KeyDir {

	return &KeyDir{idx: index.New(indexType)}
}
