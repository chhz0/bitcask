package iter

import "github.com/chhz0/go-bitcask/internal/keydir/index"

type Iterator interface {
	Rewind()
	Next()
	Key() []byte
	Valid() bool
	Value() *index.Index
	Release()
}
