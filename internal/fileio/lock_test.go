package fileio

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func Test_Lock(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "filelock")
	lock := NewLock(tmpDir)

	defer func() {
		_ = lock.unlock()
	}()

	if err := lock.TryLock(); err != nil {
		t.Logf("dir %s is lock: %v", tmpDir, err)
	} else {
		t.Logf("lock file is success: %s", tmpDir)
	}

	_ = lock.UnLock()
}

func Test_ConcurrentAccess_Lock(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "filelock")

	var wg sync.WaitGroup
	successCount := 0

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			lock := NewLock(tmpDir)
			if err := lock.TryLock(); err == nil {
				successCount++
				time.Sleep(100 * time.Millisecond)
				_ = lock.UnLock()
			}

		}(i)
	}

	wg.Wait()

	if successCount != 1 {
		t.Errorf("expect successCount is 1, but actually is %d", successCount)
	}
}
