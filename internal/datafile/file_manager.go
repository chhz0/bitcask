package datafile

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/chhz0/go-bitcask/internal/datafile/entry"
)

type FileManager struct {
	dir         string
	activeFile  iDataFile
	readOnly    map[int]iDataFile
	readLock    sync.RWMutex
	maxFileSize int64
	fileLock    sync.Mutex
}

func NewDataFileManager(dir string, maxFileSize int64) (*FileManager, error) {
	if err := checkOrMKdir(dir); err != nil {
		return nil, err
	}

	maxFileID := 0
	entris, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("read dir failed")
	}

	for _, entry := range entris {
		if strings.HasSuffix(entry.Name(), DataFileSuffix) {
			fileIDStr := strings.TrimSuffix(entry.Name(), DataFileSuffix)
			fileID, err := strconv.Atoi(fileIDStr)
			if err == nil && fileID > maxFileID {
				maxFileID = fileID
			}
		}
	}

	activeFile, err := newDataFile(dir, maxFileID, false)
	if err != nil {
		return nil, err
	}

	readOnlys := make(map[int]iDataFile)
	for i := 0; i < maxFileID; i++ {
		if file, err := newDataFile(dir, i, true); err == nil {
			readOnlys[i] = file
		}
	}

	return &FileManager{
		dir:         dir,
		activeFile:  activeFile,
		readOnly:    readOnlys,
		readLock:    sync.RWMutex{},
		maxFileSize: maxFileSize,
		fileLock:    sync.Mutex{},
	}, nil
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

func (m *FileManager) ReadAt(fileID int, off int64) (*entry.DataEntry, error) {
	if fileID == m.activeFile.FileID() {
		return m.activeFile.ReadAt(off)
	}

	m.readLock.RLock()
	file, ok := m.readOnly[fileID]
	m.readLock.RUnlock()
	if !ok {
		return nil, errors.New("file not found")
	}

	return file.ReadAt(off)
}

func (m *FileManager) Write(dataEntry *entry.DataEntry) (off int, err error) {
	m.fileLock.Lock()
	defer m.fileLock.Unlock()

	if m.activeFile.Size() >= m.maxFileSize {
		if err = m.rotateActiveFile(); err != nil {
			return 0, err
		}
	}
	off, err = m.activeFile.Write(dataEntry)
	if err != nil {
		return 0, err
	}

	return off, nil
}

func (m *FileManager) rotateActiveFile() error {
	// close active file
	oldActiveID := m.activeFile.FileID()
	if err := m.activeFile.Close(); err != nil {
		return err
	}

	// transfer active file to read only
	readOnlyFile, err := newDataFile(m.dir, oldActiveID, true)
	if err != nil {
		return err
	}

	m.readLock.Lock()
	m.readOnly[oldActiveID] = readOnlyFile
	m.readLock.Unlock()

	// create new active file
	newActiveID := oldActiveID + 1
	newActiveFile, err := newDataFile(m.dir, newActiveID, false)
	if err != nil {
		return err
	}

	m.activeFile = newActiveFile
	return nil
}

func (m *FileManager) Sync() error {
	return m.activeFile.Sync()
}

func (m *FileManager) Close() error {
	m.fileLock.Lock()
	defer m.fileLock.Unlock()

	if err := m.activeFile.Close(); err != nil {
		return err
	}

	var errs []error
	for _, file := range m.readOnly {
		if err := file.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
