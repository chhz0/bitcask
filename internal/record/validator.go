package record

import (
	"encoding/binary"
	"github.com/chhz0/go-bitcask/internal/errs"
	"hash/crc32"
)

func (lr *LogRecord) Validate() error {
	data, err := lr.Encode() // crc 会进行重新计算
	if err != nil {
		return err
	}
	if checkCRC(data) {
		return errs.ErrCRCValidation
	}
	return nil
}

func ValidateChecksum(data []byte) error {
	if len(data) < 4 {
		return errs.ErrInvalidRecord
	}

	if checkCRC(data) {
		return errs.ErrCRCValidation
	}

	return nil
}

func checkCRC(data []byte) bool {
	storedCRC := binary.LittleEndian.Uint32(data[0:4])
	expectedCRC := crc32.ChecksumIEEE(data[4:])

	return storedCRC != expectedCRC
}

func u32ToBytes(u32 uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(u32 >> 24)
	b[1] = byte(u32 >> 16)
	b[2] = byte(u32 >> 8)
	b[3] = byte(u32)
	return b
}

func u64ToBytes(u64 uint64) []byte {
	b := make([]byte, 8)
	for i := 0; i < 8; i++ {
		b[i] = byte(u64 >> (56 - i*8))
	}
	return b
}
