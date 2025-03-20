package iter

import "github.com/chhz0/go-bitcask/internal/keydir/idx"

type Iterator interface {
	Rewind()
	Next()
	Key() []byte
	Valid() bool
	Value() *idx.Index
	Release()
}
