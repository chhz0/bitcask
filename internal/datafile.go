package internal

import (
	"sync"

	"github.com/chhz0/go-bitcask/internal/codec"
	"github.com/chhz0/go-bitcask/internal/fio"
)

// DataFile 管理磁盘数据文件 => .bitcask | .hint | filelock
type DataFile struct {
	Fio      fio.File
	Offset   int64
	Dec      *codec.Decoder
	Enc      *codec.Encoder
	RW       sync.RWMutex
	ReadOnly bool
}

func NewDataFile() *DataFile {

	return &DataFile{
		// Dec: codec.NewDecoder(),
		// Enc: codec.NewEncoder(),
	}
}
