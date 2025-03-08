package codec

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/chhz0/go-bitcask/internal/entry"
	"github.com/chhz0/go-bitcask/internal/errors"
)

const (
	checksumSize = 4
	keySize      = 4
	valueSize    = 4
)

// Encoder is a encoder for encode data entry.
// todo: add buffer pool for reduce memory allocations.
type Encoder struct{}

func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode 序列化日志记录为字节流，返回数据以及总长度
// todo: use variable-length storage method
func (e *Encoder) Encode(data *entry.DataEntry) ([]byte, error) {
	ks := len(data.Key)
	vs := len(data.Value)

	if ks > entry.MaxKeySize {
		return nil, errors.ErrInvalidRecord
	}

	ts := entry.HeaderSize + ks + vs
	buf := make([]byte, ts)

	bufWithoutCRC := buf[checksumSize:]
	binary.BigEndian.PutUint32(bufWithoutCRC[0:keySize], uint32(ks))
	binary.BigEndian.PutUint32(bufWithoutCRC[keySize:keySize+valueSize], uint32(vs))
	bufWithoutCRC[keySize+valueSize] = data.Type

	copy(bufWithoutCRC[keySize+valueSize+1:], data.Key)
	copy(bufWithoutCRC[keySize+valueSize+1+ks:], data.Value)

	crc := crc32.ChecksumIEEE(bufWithoutCRC)
	binary.BigEndian.PutUint32(buf[0:4], crc)

	return buf, nil
}
