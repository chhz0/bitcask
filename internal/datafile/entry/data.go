package entry

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/chhz0/go-bitcask/internal/errors"
)

const (
	DataDeleted byte = iota
	DataNormal
	DataTransction
)

const (
	// HeaderSize 文件日志记录头大小
	HeaderSize = 4 + 8 + 4 + 4 + 1

	// MaxKeySize 文件日志记录最大key大小 2^32 - 1
	MaxKeySize = 1<<32 - 1

	// MaxValueSize 文件日志记录最大value大小 2^32 - 1
	MaxValueSize = 1<<32 - 1
)

// DataHeader 数据文件日志记录头结构, 大小为13字节
// crc(4) + timestamp(8) + keySize(4) + valueSize(4) + flag(1)
type DataHeader struct {
	CRC    uint32 // crc校验码
	Tstamp uint64 // 时间戳
	Ksz    uint32 // keySize
	Vsz    uint32 // valueSize -1 代表删除
	Flag   byte   // 标志位()
}

// DataEntry is datafile 's data entry
type DataEntry struct {
	*DataHeader
	Key   []byte
	Value []byte
}

func (e *DataEntry) Encode() ([]byte, error) {
	ks := len(e.Key)
	vs := len(e.Value)

	buf := make([]byte, HeaderSize+ks+vs)
	withoutCRC := buf[4:]

	binary.BigEndian.PutUint64(withoutCRC[:8], e.Tstamp)
	binary.BigEndian.PutUint32(withoutCRC[8:12], uint32(ks))
	binary.BigEndian.PutUint32(withoutCRC[12:16], uint32(vs))

	withoutCRC[16] = e.Flag

	copy(withoutCRC[17:], e.Key)
	copy(withoutCRC[17+ks:], e.Value)

	crc := crc32.ChecksumIEEE(withoutCRC)
	binary.LittleEndian.PutUint32(buf[:4], crc)

	return buf, nil
}

func Decode(data []byte) (*DataEntry, error) {
	if len(data) < HeaderSize {
		return nil, errors.ErrInvalidRecord
	}
	timestamp := binary.BigEndian.Uint64(data[4:12])
	keySize := binary.BigEndian.Uint32(data[12:16])
	valueSize := binary.BigEndian.Uint32(data[16:20])
	flag := data[20]

	if len(data) < HeaderSize+int(keySize)+int(valueSize) {
		return nil, errors.ErrInvalidRecord
	}

	dataEntry := &DataEntry{
		DataHeader: &DataHeader{
			CRC:    binary.LittleEndian.Uint32(data[:4]),
			Tstamp: timestamp,
			Ksz:    keySize,
			Vsz:    valueSize,
			Flag:   flag,
		},
		Key:   make([]byte, keySize),
		Value: make([]byte, valueSize),
	}

	copy(dataEntry.Key, data[HeaderSize:HeaderSize+keySize])
	copy(dataEntry.Value, data[HeaderSize+keySize:HeaderSize+keySize+valueSize])

	if dataEntry.CRC != calculateCRC(dataEntry) {
		return nil, errors.ErrCRCValidation
	}

	return dataEntry, nil
}

func calculateCRC(e *DataEntry) uint32 {
	buf := make([]byte, 8+4+4+1)
	binary.BigEndian.PutUint64(buf[:8], e.Tstamp)
	binary.BigEndian.PutUint32(buf[8:12], e.Ksz)
	binary.BigEndian.PutUint32(buf[12:16], e.Vsz)
	buf[16] = e.Flag

	crc := crc32.NewIEEE()
	crc.Write(buf)
	crc.Write(e.Key)
	crc.Write(e.Value)
	return crc.Sum32()
}

// HeaderWithTimeout 数据文件日志记录头结构, 大小为29字节
// crc[4] + keySize[4] + valueSize[4] + tstamp[8] + exp[8]
type HeaderWithTimeout struct {
	*DataHeader
	Tstamp uint64 // 时间戳
	Exp    uint64 // 过期时间
}

type DataEntryWithTimeout struct {
	*HeaderWithTimeout
	Key   []byte
	Value []byte
}
