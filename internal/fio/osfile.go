package fio

import (
	"fmt"
	"os"
	"path/filepath"
)

const formatFile = "%08d.bit"

type osFile struct {
	fd    *os.File
	fid   uint32 // 文件ID
	fpath string // 文件路径
}

func NewOsFile(fileID uint32, path string) File {
	return &osFile{
		fid:   fileID,
		fpath: path,
	}
}

// Open 打开/创建数据文件
func (f *osFile) Open() error {
	fullFile := filepath.Join(f.fpath, fmt.Sprintf(formatFile, f.fid))
	fd, err := os.OpenFile(
		fullFile,
		os.O_CREATE|os.O_RDWR|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	f.fd = fd
	return nil
}

func (f *osFile) ReadAt(b []byte, off int64) (n int, err error) {
	n, err = f.fd.ReadAt(b, off)
	return n, err
}

func (f *osFile) Write(b []byte) (n int, err error) {
	n, err = f.fd.Write(b)
	return n, err
}

func (f *osFile) Sync() error {
	return f.fd.Sync()
}

func (f *osFile) Close() error {
	return f.fd.Close()
}

func (f *osFile) Stat() (os.FileInfo, error) {
	return f.fd.Stat()
}

func (f *osFile) FileID() uint32 {
	return f.fid
}

func (f *osFile) Path() string {
	return f.fpath
}

func (f *osFile) Destroy() {
	_ = f.fd.Close()
	_ = os.Remove(filepath.Join(f.fpath, fmt.Sprintf(formatFile, f.fid)))
}
