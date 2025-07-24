package codec

import (
	"encoding/binary"
	"errors"

	"github.com/chhz0/bitcask/internal"
)

var (
	ErrInvalidHeader  = errors.New("invalid entry header size.")
	ErrIncompleteRead = errors.New("incomplate data read.")
	ErrCRCValidation  = errors.New("crc validation failed.")
)

func Decode(b []byte) (*internal.Entry, error) {
	if len(b) < headerSize {
		return nil, ErrInvalidHeader
	}

	crc := binary.BigEndian.Uint32(b[0:crcSize])
	tstamp := int64(binary.BigEndian.Uint64(b[crcSize:bufTstampEndIdx]))
	ksz := binary.BigEndian.Uint32(b[bufTstampEndIdx:bufKszEndIdx])
	vsz := binary.BigEndian.Uint32(b[bufKszEndIdx:bufVszEndIdx])

	totalSize := headerSize + ksz + vsz
	if len(b) < int(totalSize) {
		return nil, ErrIncompleteRead
	}

	if !verifyCRC(b[crcSize:totalSize], crc) {
		return nil, ErrCRCValidation
	}

	keyStart := headerSize
	keyEnd := keyStart + int(ksz)
	key := make([]byte, ksz)
	copy(key, b[keyStart:keyEnd])

	value := []byte{}
	if vsz > 0 {
		value = make([]byte, vsz)
		copy(value, b[keyEnd:keyEnd+int(vsz)])
	}

	return &internal.Entry{
		CRC:    crc,
		Tstamp: tstamp,
		Key:    key,
		Val:    value,
	}, nil
}
