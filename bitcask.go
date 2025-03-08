package bitcask

import (
	"sync"

	"github.com/chhz0/go-bitcask/internal"
)

type Bitcask struct {
	rw       sync.RWMutex         // rw lock
	active   *internal.DataFile   // active
	readOnly []*internal.DataFile // read only
	keydir   *internal.KeyDir     // keydir
	options  *Options
}

func Open(path string, opts ...Option) (*Bitcask, error) {

	return nil, nil
}
