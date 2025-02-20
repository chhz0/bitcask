package internal

import (
	"sync"

	"github.com/chhz0/go-bitcask/internal/fio"
)

// DataFile 管理磁盘数据文件的统一接口 .bitcask | .hint | filelock
type DataFile struct {
	io       fio.FileIO
	writeOff int64
	isActive bool
	mu       sync.RWMutex
}
