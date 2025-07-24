package codec

import "hash/crc32"

// var crcTable = crc32.MakeTable(crc32.IEEE)

func calculateCRC(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func verifyCRC(data []byte, crc uint32) bool {
	return crc32.ChecksumIEEE(data) == crc
}
