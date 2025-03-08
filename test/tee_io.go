package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
)

type TeeR struct {
	s string
	i int64
}

func (tr *TeeR) Read(p []byte) (n int, err error) {
	if tr.i >= int64(len(tr.s)) {
		return 0, io.EOF
	}
	n = copy(p, tr.s[tr.i:])
	tr.i += int64(n)
	return
}

type TeeW struct{}

func (tw *TeeW) Write(p []byte) (n int, err error) {
	fmt.Printf("%s", p)
	return len(p), nil
}

func main() {
	var r io.Reader
	var w io.Writer

	r = &TeeR{"Hello, Reader", 0}
	w = &TeeW{}
	// 分流1
	pr, pw := io.Pipe()
	r = io.TeeReader(r, pw)
	repErrCh := make(chan error, 1)
	go func() {
		var repErr error
		defer func() {
			pr.CloseWithError(repErr)
			repErrCh <- repErr
		}()

		_, repErr = io.CopyBuffer(os.Stderr, pr, make([]byte, 1))
	}()

	crcW := crc32.NewIEEE()
	r = io.TeeReader(r, crcW)

	io.CopyBuffer(w, r, make([]byte, 1))
	fmt.Printf("crc32: %x\n", crcW.Sum32())
	pw.Close()
	<-repErrCh
}
