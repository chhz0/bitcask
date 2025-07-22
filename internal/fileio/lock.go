package fileio

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

var (
	ErrAlreadyLocked = errors.New("directory is already locked")
	ErrNotLocked     = errors.New("directory is not locked")
	ErrLockTimeout   = errors.New("lock acquisition timed out")
	ErrLockFailed    = errors.New("failed to acquisition lock")
	ErrUnSupported   = errors.New("operation not supported on this platform")
)

// Lock is a dir lock :)
type Lock struct {
	path    string
	file    *os.File
	locked  bool
	timeout time.Duration
}

func NewLock(dir string) *Lock {
	return &Lock{
		path:   filepath.Clean(dir),
		locked: false,
	}
}

func (l *Lock) Lock() error {
	return l.lock(false, 0)
}

func (l *Lock) TryLock() error {
	return l.lock(true, 0)
}

func (l *Lock) LockWithTimeout(timeout time.Duration) error {
	return l.lock(false, timeout)
}

func (l *Lock) UnLock() error {
	if !l.locked {
		return ErrNotLocked
	}

	if err := l.unlock(); err != nil {
		return err
	}

	if err := l.file.Close(); err != nil {
		return err
	}

	l.locked = false
	return nil
}

func (l *Lock) IsLocked() bool {
	return l.locked
}

func (l *Lock) Path() string {
	return l.path
}

func (l *Lock) lock(nonblocking bool, timeout time.Duration) error {
	if l.locked {
		return ErrAlreadyLocked
	}

	if err := os.MkdirAll(l.path, 0755); err != nil {
		return err
	}

	lockFile := filepath.Join(l.path, ".flock")
	file, err := os.OpenFile(lockFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	fd := int(file.Fd())

	if timeout > 0 {
		deadline := time.Now().Add(timeout)

		for {
			err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
			if err == nil {
				break
			}

			if time.Now().After(deadline) {
				return ErrLockTimeout
			}

			time.Sleep(50 * time.Millisecond)
		}
	} else if nonblocking {
		err = syscall.Flock(fd, syscall.LOCK_EX|syscall.LOCK_NB)
	} else {
		err = syscall.Flock(fd, syscall.LOCK_EX)
	}

	if err != nil {
		file.Close()
		return err
	}

	l.file = file
	l.locked = true
	l.timeout = timeout
	return nil
}

func (l *Lock) unlock() error {
	return syscall.Flock(int(l.file.Fd()), syscall.LOCK_UN)
}
