package fio

import "os"

type Mmap struct {
	file *os.File
	mmap []byte
}
