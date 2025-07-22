package datafile

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chhz0/go-bitcask/internal/datafile/entry"
)

const (
	dataFileSuffix = ".bit"
	dataFileFormat = "%08d"
)

type iDataFile interface {
	ReadAt(off int64) (*entry.DataEntry, error)
	Write(dataEntry *entry.DataEntry) (off int, err error)
	Sync() error
	Close() error
	Size() int64
	FileID() int
	ReadOnly() bool
	Destroy() error
}

type dataFile struct {
	file     *os.File
	fid      int
	fpath    string
	mu       sync.Mutex
	size     int64
	readOnly bool
}

var _ iDataFile = (*dataFile)(nil)

func newDataFile(dir string, fileID int, readOnly bool) (*dataFile, error) {
	filePath := filepath.Join(dir, fmt.Sprintf(dataFileFormat+dataFileSuffix, fileID))
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		return nil, err
	}

	return &dataFile{
		file:     file,
		fid:      fileID,
		fpath:    filePath,
		size:     stat.Size(),
		readOnly: readOnly,
	}, nil
}

func (f *dataFile) ReadAt(off int64) (*entry.DataEntry, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	header := make([]byte, entry.HeaderSize)
	_, err := f.file.ReadAt(header, off)
	if err != nil {
		return nil, err

	}

	flag := header[20]
	if flag == entry.DataDeleted {
		return nil, errors.New("data entry is deleted")
	}

	keySize := binary.BigEndian.Uint32(header[12:16])
	valueSize := binary.BigEndian.Uint32(header[16:20])

	dataEntryByte := make([]byte, entry.HeaderSize+int(keySize)+int(valueSize))
	_, err = f.file.ReadAt(dataEntryByte, off)
	if err != nil {
		return nil, err
	}

	return entry.Decode(dataEntryByte)
}

func (f *dataFile) Write(dataEntry *entry.DataEntry) (off int, err error) {
	if f.readOnly {
		return 0, errors.New("read-only file doesn't support write")
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if len(dataEntry.Key) > entry.MaxKeySize || len(dataEntry.Value) > entry.MaxValueSize {
		return 0, errors.New("key or value size exceeds the maximum limit")
	}

	dataEntry.Tstamp = uint64(time.Now().Unix())
	dataEntry.Flag = 0
	dataBytes, err := dataEntry.Encode()
	if err != nil {
		return 0, err
	}

	off = int(f.size)
	n, err := f.file.Write(dataBytes)
	if err != nil {
		return 0, err
	}

	atomic.AddInt64(&f.size, int64(n))
	return
}

func (f *dataFile) Sync() error {
	if f.readOnly {
		return errors.New("read-only file doesn't support sync")
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.file.Sync()
}

func (f *dataFile) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.file.Close()
}

func (f *dataFile) Size() int64 {
	return atomic.LoadInt64(&f.size)
}

func (f *dataFile) FileID() int {
	return f.fid
}

func (f *dataFile) ReadOnly() bool {
	return f.readOnly
}

func (f *dataFile) Destroy() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(f.fpath, fmt.Sprintf(dataFileFormat+dataFileSuffix, f.fid)))
}
