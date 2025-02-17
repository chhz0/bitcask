package utils

import "unsafe"

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

func IntToUint32(i int) uint32 {
	return *(*uint32)(unsafe.Pointer(&i))
}

func Int64ToInt(i64 int64) int {
	return *(*int)(unsafe.Pointer(&i64))
}

func UInt64ToUInt32(u64 uint64) uint32 {
	return *(*uint32)(unsafe.Pointer(&u64))
}
