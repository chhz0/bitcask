package record

import (
	"encoding/binary"
	"github.com/chhz0/go-bitcask/internal/errs"
	"hash/crc32"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRecord_Encode_ValidRecord(t *testing.T) {
	lr := &LogRecord{
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
		Timestamp: 1234567890,
		Delete:    false,
	}

	encoded, err := lr.Encode()
	assert.NoError(t, err)

	// 检查编码后的数据长度是否正确
	assert.Equal(t, HeaderSize+len(lr.Key)+len(lr.Value), len(encoded))

	// 检查CRC校验码是否正确
	expectedCRC := crc32.ChecksumIEEE(encoded[4:])
	actualCRC := binary.LittleEndian.Uint32(encoded[0:4])
	assert.Equal(t, expectedCRC, actualCRC)

	// 检查timestamp是否正确
	expectedTimestamp := binary.LittleEndian.Uint64(encoded[4:12])
	assert.Equal(t, lr.Timestamp, uint64(expectedTimestamp))

	// 检查keySize和valueSize是否正确
	expectedKeySize := binary.LittleEndian.Uint16(encoded[12:14])
	assert.Equal(t, uint16(len(lr.Key)), expectedKeySize)
	expectedValueSize := binary.LittleEndian.Uint32(encoded[14:18])
	assert.Equal(t, uint32(len(lr.Value)), expectedValueSize)

	// 检查delete标志是否正确
	expectedDelete := encoded[18]
	assert.Equal(t, uint8(LogRecordNormal), expectedDelete)

	// 检查key和value是否正确
	expectedKey := encoded[HeaderSize : HeaderSize+len(lr.Key)]
	assert.Equal(t, lr.Key, expectedKey)
	expectedValue := encoded[HeaderSize+len(lr.Key):]
	assert.Equal(t, lr.Value, expectedValue)

}

func TestLogRecord_Encode_InvalidKeySize(t *testing.T) {
	lr := &LogRecord{
		Key:       make([]byte, MaxKeySize+1),
		Value:     []byte("test-value"),
		Timestamp: 1234567890,
		Delete:    false,
	}
	encoded, err := lr.Encode()
	assert.ErrorIs(t, err, errs.ErrInvalidRecord)
	assert.Nil(t, encoded)
}

func TestLogRecord_Encode_DeleteRecode(t *testing.T) {
	lr := &LogRecord{
		Key:       []byte("test-key"),
		Value:     []byte("test-value"),
		Timestamp: 1234567890,
		Delete:    true,
	}

	encoded, err := lr.Encode()
	assert.NoError(t, err)

	assert.Equal(t, uint8(LogRecordDeleted), encoded[18])
}

func TestDecode_InvalidHeaderLength_ShouldError(t *testing.T) {
	buf := make([]byte, HeaderSize-1) // 创建一个长度为HeaderSize-1的切片
	record, err := DecodeLogRecord(buf)
	assert.Nil(t, record)
	assert.ErrorIs(t, err, errs.ErrInvalidRecord)
}

func TestDecode_HeaderShortForKV_ShouldError(t *testing.T) {
	buf := make([]byte, HeaderSize+1)
	binary.LittleEndian.PutUint16(buf[12:14], 10) // 设置keySize为10
	binary.LittleEndian.PutUint32(buf[14:18], 10) // 设置valueSize为10
	record, err := DecodeLogRecord(buf)
	assert.Nil(t, record)
	assert.ErrorIs(t, err, errs.ErrInvalidRecord)
}

func TestDecode_ValidBuf_ShouldDecodeLogRecord(t *testing.T) {
	timestamp := uint64(1739779831104524900)
	key := []byte("test-key")
	value := []byte("test-value")

	lr := &LogRecord{
		Timestamp: timestamp,
		Key:       key,
		Value:     value,
		Delete:    false,
	}
	encoded, err := lr.Encode()
	assert.NoError(t, err)

	record, err := DecodeLogRecord(encoded)
	assert.NoError(t, err)
	assert.Equal(t, uint32(0x81612044), record.CRC)
	assert.Equal(t, timestamp, record.Timestamp)
	assert.Equal(t, key, record.Key)
	assert.Equal(t, value, record.Value)
	assert.False(t, record.Delete)
}
