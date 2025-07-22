package datafile

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chhz0/go-bitcask/internal/datafile/entry"
	"github.com/chhz0/go-bitcask/internal/keydir"
)

const (
	hintFilePrefix     = "hint"
	hintFileTempFormat = "%d.%s.tmp"
	hintFileVersion    = 1
)

type hint struct {
	dir string
}

func NewHint(dir string) *hint {
	return &hint{dir}
}

func (h *hint) WriteHint(keyDir keydir.KeyDir) (string, error) {
	tempFile := filepath.Join(h.dir,
		fmt.Sprintf(hintFileTempFormat, time.Now().UnixNano(), hintFilePrefix),
	)
	f, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create hint file: %v", err)
	}
	defer f.Close()

	header := entry.HintHeader{
		Version: hintFileVersion,
	}
	if err := binary.Write(f, binary.BigEndian, &header); err != nil {
		return "", fmt.Errorf("failed to write hint header: %v", err)

	}

	// keys := keyDir.
	return "", nil
}

func (h *hint) FindLatestHint(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// var latestTime int64
	// var latesFile string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

	}

	return "", nil
}
