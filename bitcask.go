package bitcask

import (
	"os"
	"sync"

	"github.com/chhz0/go-bitcask/internal/datafile"
	"github.com/chhz0/go-bitcask/internal/keydir"
	"github.com/gofrs/flock"
)

const (
	SHARDHASH = "shardhash"
	BTREE     = "btree"
)

type Bitcask struct {
	rw        sync.RWMutex // rw lock
	flock     *flock.Flock
	fileMgr   *datafile.FileManager
	keydir    *keydir.KeyDir // keydir
	options   *Options
	isMerging bool
}

// Open a new or an existing bitcask datastore
func Open(dirPath string, opts ...Option) (*Bitcask, error) {
	o := &Options{
		Dir:         dirPath,
		MaxFileSize: 1 << 30, // 1GB
		SyncOnWrite: false,
		ReadOnly:    false,
	}

	for _, opt := range opts {
		opt(o)
	}

	// check config file && options
	// if config file no exists, create a new  default config file
	if err := checkOrMKdir(dirPath); err != nil {
		return nil, ErrCheckOrMkdir
	}

	// create a new bitcask instance
	bitcask := &Bitcask{
		rw:      sync.RWMutex{},
		keydir:  keydir.NewKeyDir(SHARDHASH),
		options: o,
	}

	// try to get file lock

	// loadavtivefile

	// loadKeydir

	return bitcask, nil
}

func loadDataFiles(dir string, keyDir *keydir.KeyDir) error {

	return nil
}

// Put Stores a key and a value in the bitcask datastore
func (b *Bitcask) Put(key []byte, value []byte) error {
	return nil
}

// Get Reads a value by key from a datastore
func (b *Bitcask) Get(key []byte) ([]byte, error) {
	return nil, nil
}

// Delete Removes a key from the datastore
func (b *Bitcask) Delete(key []byte) error {
	return nil
}

// Close a bitcask data store and flushes all pending writes to disk
func (b *Bitcask) Close() error {
	return nil
}

// ListKey Returns list of all keys
func (b *Bitcask) ListKeys() error {
	return nil
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
