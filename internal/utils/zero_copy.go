package utils

import "unsafe"

type zero struct{}

func ZeroCopy() *zero {
	return &zero{}
}

// BytesToString converts a byte slice to a string without copying the underlying data.
func (*zero) BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (*zero) StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
