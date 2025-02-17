package record

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

var (
	ErrInvalidRecord = errors.New("invalid record")
)

const (
	LogRecordDeleted byte = iota
	LogRecordNormal
)

const (
	// HeaderSize 文件日志记录头大小
	// crc[4] + timestamp[8] + keySize[2] + valueSize[4] + delete[1]
	HeaderSize = 4 + 8 + 2 + 4 + 1

	// MaxKeySize 文件日志记录最大key大小 65535
	MaxKeySize = 1<<16 - 1 // 2^16 - 1
)

// LogRecord 文件日志记录结构
type LogRecord struct {
	CRC       uint32 // crc校验码(Header + Key + Value)
	Timestamp uint64 // 时间戳 (unix时间戳)
	Key       []byte
	Value     []byte
	Delete    bool
}

// Encode 序列化日志记录为字节流，返回数据以及总长度
// todo: 解决header中由于keySize和valueSize是非固定长度，可能会造成空间浪费，但是因此会导致header变成无固定长度，需要重新实现encode和decode方法
func (lr *LogRecord) Encode() ([]byte, error) {
	keySize := len(lr.Key)
	valueSize := len(lr.Value)

	if keySize > MaxKeySize {
		return nil, ErrInvalidRecord
	}

	toalSize := HeaderSize + keySize + valueSize
	buf := make([]byte, toalSize)

	bufWithoutCRC := buf[4:]
	binary.LittleEndian.PutUint64(bufWithoutCRC[0:8], lr.Timestamp)
	binary.LittleEndian.PutUint16(bufWithoutCRC[8:10], uint16(keySize))
	binary.LittleEndian.PutUint32(bufWithoutCRC[10:14], uint32(valueSize))
	if lr.Delete {
		bufWithoutCRC[14] = LogRecordDeleted
	} else {
		bufWithoutCRC[14] = LogRecordNormal
	}
	copy(bufWithoutCRC[15:15+keySize], lr.Key)
	copy(bufWithoutCRC[15+keySize:], lr.Value)

	crc := crc32.ChecksumIEEE(bufWithoutCRC)
	binary.LittleEndian.PutUint32(buf[0:4], crc)

	return buf, nil
}

// DecodeLogRecord 解析字节流为日志记录
func DecodeLogRecord(buf []byte) (*LogRecord, error) {
	if len(buf) < HeaderSize {
		return nil, ErrInvalidRecord
	}

	crc := binary.LittleEndian.Uint32(buf[0:4])
	timestamp := binary.LittleEndian.Uint64(buf[4:12])
	keySize := binary.LittleEndian.Uint16(buf[12:14])
	valueSize := binary.LittleEndian.Uint32(buf[14:18])
	delete := buf[18] == LogRecordDeleted

	if len(buf) < HeaderSize+int(keySize)+int(valueSize) {
		return nil, ErrInvalidRecord
	}

	key := buf[HeaderSize : HeaderSize+keySize]
	value := buf[HeaderSize+keySize:]

	return &LogRecord{
		CRC:       crc,
		Timestamp: timestamp,
		Key:       key,
		Value:     value,
		Delete:    delete,
	}, nil
}
