package bitcask

import "sync"

type bitcask struct {
	rw sync.RWMutex
}
