package utils

import "unsafe"

// BytesToString converts a byte slice to a string without copying the underlying data.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
