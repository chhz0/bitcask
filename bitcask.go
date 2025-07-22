package bitcask

import (
	"os"
	"sync"

	"github.com/chhz0/go-bitcask/internal/keydir"
)

type Bitcask struct {
	rw        sync.RWMutex // rw lock
	options   *Options
	isMerging bool
}

// Open 打开或者创建 Bitcask, 支持读写 写入同步等
// 仅支持单线程
func Open(dir string, opts ...Option) (*Bitcask, error) {
	o := &Options{
		Dir:         dir,
		MaxFileSize: 1 << 30, // 1GB
		SyncOnWrite: false,
		ReadOnly:    false,
	}

	for _, opt := range opts {
		opt(o)
	}

	// check config file && options
	// if config file no exists, create a new  default config file
	if err := checkOrMKdir(dir); err != nil {
		return nil, ErrCheckOrMkdir
	}

	// create a new bitcask instance
	bitcask := &Bitcask{
		rw:      sync.RWMutex{},
		options: o,
	}

	// try to get file lock

	// loadavtivefile

	// loadKeydir

	return bitcask, nil
}

// OpenReadOnly 以只读模式打开 Bitcask
func OpenReadOnly(dir string) *Bitcask {
	return nil
}

// Put Stores a key and a value in the bitcask datastore
func (b *Bitcask) Put(key []byte, value []byte) bool {
	return true
}

// Get Reads a value by key from a datastore
func (b *Bitcask) Get(key []byte) (bool, []byte) {
	return true, nil
}

// Delete Removes a key from the datastore
func (b *Bitcask) Delete(key []byte) bool {
	return true
}

// Close a bitcask data store and flushes all pending writes to disk
func (b *Bitcask) Close() error {
	return nil
}

// ListKey Returns list of all keys
func (b *Bitcask) ListKeys() ([][]byte, error) {
	return nil, nil
}

// Sync Force any writes to sync to disk
func (b *Bitcask) Sync() error {
	return nil
}

// Merge Merge several data files within a Bitcask datastore into a more compact form.
// Also, produce hintfiles for faster startup.
func (b *Bitcask) Merge() error {
	return nil
}

// Fold over all K/V pairs in a Bitcask datastore.
// → Acc Fun is expected to be of the form: F(K,V,Acc0) → Acc
func (b *Bitcask) Fold(fn func([]byte, []byte, any) any, acc any) any {
	return nil
}

func checkOrMKdir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func loadDataFiles(dir string, keyDir *keydir.KeyDir) error {

	return nil
}
