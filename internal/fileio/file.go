package fileio

import (
	"errors"
	"os"
)

var (
	ErrFileClosed   = errors.New("file is closed.")
	ErrInvalidSeek  = errors.New("invalid seed operation.")
	ErrReadOnly     = errors.New("file opened in read-only mode.")
	ErrWriteOnly    = errors.New("file opened in write-only mode.")
	ErrBufferTooBig = errors.New("requested buffer siez exceeds maximum.")
)

type FileIO interface {
	ReadAt(b []byte, off int64) (n int, err error)
	Write(b []byte) (n int, err error)
	Sync() error
	Close() error
	Size() (size int64, err error)
}

func Open(filename string) (FileIO, error) {
	return newFile(filename)
}

// func OpenWithBuf(filename string)   {}
// func OpenReadOnly(filename string)  {}
// func OpenWriteOnly(filename string) {}
// func OpenAppend(filename string)    {}

// type fileOptions struct {
// 	filename   string
// 	flag       int
// 	perm       os.FileMode
// 	bufferSize int
// }

// func open(fOpts *fileOptions) (*os.File, error) {

// 	return nil, nil
// }

type File struct {
	f *os.File
}

func newFile(fName string) (*File, error) {
	file, err := os.OpenFile(fName,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	return &File{f: file}, nil
}

func (f *File) ReadAt(b []byte, off int64) (n int, err error) {
	return f.f.ReadAt(b, off)
}

func (f *File) Write(b []byte) (n int, err error) {
	return f.f.Write(b)
}

func (f *File) Sync() error {
	return f.f.Sync()
}

func (f *File) Close() error {
	return f.f.Close()
}

func (f *File) Size() (size int64, err error) {
	stat, err := f.f.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil
}
