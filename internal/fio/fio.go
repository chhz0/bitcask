package fio

// FileIO 抽象文件操作接口
type FileIO interface {
	Open(path string) error
	ReadAt(offset int64) ([]byte, error)
	Write(buf []byte) (int, error)
	Sync() error
	Close() error
	FileID() uint32
	Path() string
}
