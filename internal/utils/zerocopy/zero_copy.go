// zerocopy provides zero-copy operations for byte slices.
// use unsafe may be dangerous, if you don't know what you are doing, don't use it.
package zerocopy

import (
	"unsafe"
)

// BytesToString converts a byte slice to a string without copying the underlying data.
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

func U32ToBytes(u32 uint32) []byte {
	return (*[4]byte)(unsafe.Pointer(&u32))[:]
}

func IntToU32(i int) uint32 {
	return *(*uint32)(unsafe.Pointer(&i))
}

func I64ToInt(i64 int64) int {
	return *(*int)(unsafe.Pointer(&i64))
}

func U64ToU32(u64 uint64) uint32 {
	return *(*uint32)(unsafe.Pointer(&u64))
}
