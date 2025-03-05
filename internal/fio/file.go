package fio

import "os"

// File 抽象文件操作接口
type File interface {
	Open() error
	ReadAt(b []byte, off int64) (n int, err error)
	Write(b []byte) (n int, err error)
	Sync() error
	Close() error
	Stat() (os.FileInfo, error)
	FileID() uint32
	Path() string

	Destroy()
}
