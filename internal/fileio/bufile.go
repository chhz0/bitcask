package fileio

import (
	"bufio"
	"os"
)

const (
	BufferSiez    = 64 * 1024
	MaxBufferSize = 1 << 20
)

type Bufile struct {
	f *os.File
	// r      *bufio.Reader
	w *bufio.Writer
}

func newBufile(f string, bufSize int) (*Bufile, error) {
	file, err := os.OpenFile(f,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}

	return &Bufile{
		f: file,
		// r: bufio.NewReaderSize(file, bufSize),
		w: bufio.NewWriterSize(file, bufSize),
	}, nil
}

// ReadAt implements FileIO.
func (bf *Bufile) ReadAt(b []byte, off int64) (n int, err error) {
	return bf.f.ReadAt(b, off)
}

// Write implements FileIO.
func (bf *Bufile) Write(b []byte) (n int, err error) {
	return bf.w.Write(b)
}

// Sync implements FileIO.
func (bf *Bufile) Sync() error {
	if err := bf.w.Flush(); err != nil {
		return err
	}
	return bf.f.Sync()
}

// Close implements FileIO.
func (bf *Bufile) Close() error {
	if err := bf.Sync(); err != nil {
		_ = bf.f.Close()
		return err
	}

	return bf.f.Close()
}

// Size implements FileIO.
func (bf *Bufile) Size() (size int64, err error) {
	stat, err := bf.f.Stat()
	if err != nil {
		return 0, err
	}

	return stat.Size() + int64(bf.w.Buffered()), nil
}

var _ FileIO = (*Bufile)(nil)
