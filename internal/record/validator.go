package record

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

var ErrCRCValidation = errors.New("crc validation failed")

func (lr *LogRecord) Validate() error {
	data, err := lr.Encode() // crc 会进行重新计算
	if err != nil {
		return err
	}
	if checkCRC(data) {
		return ErrCRCValidation
	}
	return nil
}

func ValidateChecksum(data []byte) error {
	if len(data) < 4 {
		return ErrInvalidRecord
	}

	if checkCRC(data) {
		return ErrCRCValidation
	}

	return nil
}

func checkCRC(data []byte) bool {
	storedCRC := binary.LittleEndian.Uint32(data[0:4])
	expectedCRC := crc32.ChecksumIEEE(data[4:])

	return storedCRC != expectedCRC
}
