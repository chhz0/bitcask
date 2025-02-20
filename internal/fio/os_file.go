package fio

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type osFile struct {
	file     *os.File
	fileID   uint32
	filePath string
	fileSize int64
	mu       sync.Mutex
}

func NewOsFile(fileID uint32, path string) *osFile {
	return &osFile{
		fileID:   fileID,
		filePath: filepath.Join(path, fmt.Sprintf("%08d.cask", fileID)),
	}
}
