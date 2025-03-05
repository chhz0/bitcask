package fio

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFile(t *testing.T) {
	file := NewOsFile(1, os.TempDir())
	defer file.Destroy()

	err := file.Open()
	assert.Nil(t, err)
}

func Test_File_Write(t *testing.T) {
	file := NewOsFile(1, os.TempDir())
	defer file.Destroy()

	err := file.Open()
	assert.Nil(t, err)

	n, err := file.Write([]byte(""))
	assert.Nil(t, err)
	assert.Equal(t, 0, n)

	n, err = file.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	n, err = file.Write([]byte("bitcask"))
	assert.Nil(t, err)
	assert.Equal(t, 7, n)
}

func Test_File_Read(t *testing.T) {
	file := NewOsFile(1, os.TempDir())
	defer file.Destroy()

	err := file.Open()
	assert.Nil(t, err)

	n, err := file.Write([]byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, 5, n)

	buf := make([]byte, 5)
	n, err = file.ReadAt(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte("hello"), buf)
}
