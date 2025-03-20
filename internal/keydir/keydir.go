package keydir

import "github.com/chhz0/go-bitcask/internal/keydir/idx"

// KeyDir 管理内存索引
type KeyDir struct {
	idx.Indexer
}

func NewKeyDir(indexType idx.IndexType) *KeyDir {
	return &KeyDir{idx.New(indexType)}
}

func (kr *KeyDir) Store(key []byte, value *idx.Index) {
	kr.Put(key, value)
}

func (kr *KeyDir) Load(key []byte) (*idx.Index, bool) {
	val, ok := kr.Get(key)

	return val, ok
}

func (kr *KeyDir) Remove(key []byte) (*idx.Index, bool) {
	return kr.Del(key)
}

func (kr *KeyDir) Keys() [][]byte {
	return kr.Indexer.Keys()
}
