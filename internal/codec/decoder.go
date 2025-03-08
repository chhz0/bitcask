package codec

import (
	"encoding/binary"

	"github.com/chhz0/go-bitcask/internal/entry"
	"github.com/chhz0/go-bitcask/internal/errors"
)

const (
	headerSize            = 4 + 4 + 4 + 1
	headerSizeWithTimeout = headerSize + 8 + 8
)

type Decoder struct{}

func NewDecoder() *Decoder { return &Decoder{} }

func (d *Decoder) Decode(data []byte) (*entry.DataEntry, error) {
	if len(data) < headerSize {
		return nil, errors.ErrInvalidRecord
	}

	crc := binary.BigEndian.Uint32(data[0:checksumSize])
	ks := binary.BigEndian.Uint32(data[checksumSize : checksumSize+keySize])
	vs := binary.BigEndian.Uint32(data[checksumSize+keySize : checksumSize+keySize+valueSize])
	t := data[checksumSize+keySize+valueSize]

	if len(data) < headerSize+int(ks)+int(vs) {
		return nil, errors.ErrInvalidRecord
	}

	key := data[headerSize : headerSize+ks]
	value := data[headerSize+ks:]

	return &entry.DataEntry{
		Header: &entry.Header{
			CRC:  crc,
			Ks:   ks,
			Vs:   vs,
			Type: t,
		},
		Key:   key,
		Value: value,
	}, nil
}

func (d *Decoder) DecodeHeader(data []byte) (*entry.Header, error) {
	return &entry.Header{}, nil
}
