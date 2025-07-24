package codec

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Case gen by Doubao AI
func TestDecode_NormalCase(t *testing.T) {
	// 准备测试数据
	key := []byte("normalKey")
	val := []byte("normalValue")
	encoded := Encode(key, val, false)

	// 执行解码
	entry, err := Decode(encoded)

	// 基础验证
	require.NoError(t, err,
		"Normal data decoding should not return errors")
	assert.NotNil(t, entry,
		"The decoded result should not be nil")

	// 验证元数据
	assert.Equal(t, uint32(calculateCRC(encoded[crcSize:])), entry.CRC,
		"CRC value does not match")
	assert.InDelta(t, time.Now().Unix(), entry.Tstamp, 1,
		"Timestamp deviation is too large") // 允许1秒误差

	// 验证键值
	assert.Equal(t, key, entry.Key,
		"The decoded key does not match")
	assert.Equal(t, val, entry.Val,
		"The decoded value does not match")
}

func TestDecode_DeletedEntry(t *testing.T) {
	// 准备删除标记的数据（vsz=-1，实际存储为0xffffffff）
	key := []byte("deletedKey")
	val := []byte("shouldBeIgnored")
	encoded := Encode(key, val, true)

	// 执行解码
	entry, err := Decode(encoded)

	// 验证
	require.NoError(t, err, "Decoding of data marked for deletion should not return an error")
	assert.Equal(t, key, entry.Key, "删除条Delete entry key mismatch目键不匹配")
	assert.Empty(t, entry.Val, "To delete an entry the value should be empty") // 因为vsz=0时Decode中vsz>0不成立
}

func TestDecode_EmptyKeyAndValue(t *testing.T) {
	// 准备空键空值数据
	encoded := Encode([]byte(""), []byte(""), false)

	// 执行解码
	entry, err := Decode(encoded)

	// 验证
	require.NoError(t, err, "Decoding a nil key value should not return an error")
	assert.Empty(t, entry.Key, "Nil key decoding error")
	assert.Empty(t, entry.Val, "Nil value decoding error")
}

func TestDecode_InvalidHeader(t *testing.T) {
	// 测试场景1：数据长度小于头部大小
	invalidData := make([]byte, headerSize-1)
	entry, err := Decode(invalidData)
	assert.ErrorIs(t, err, ErrInvalidHeader, "Should return an invalid header error")
	assert.Nil(t, entry, "Should return nil for invalid header")

	// 测试场景2：头部完整但数据不完整（总长度不足）
	validHeader := Encode([]byte("short"), []byte("longvalue"), false)
	truncatedData := validHeader[:headerSize+3] // 截断键值部分
	entry, err = Decode(truncatedData)
	assert.ErrorIs(t, err, ErrIncompleteRead, "Should return an incomplete data error")
	assert.Nil(t, entry, "If the data is incomplete, nil should be returned.")
}

func TestDecode_CRCValidationFailed(t *testing.T) {
	// 生成正常数据后篡改内容，导致CRC校验失败
	validData := Encode([]byte("crcKey"), []byte("crcVal"), false)
	tamperedData := make([]byte, len(validData))
	copy(tamperedData, validData)
	tamperedData[headerSize] ^= 0x01 // 篡改一个字节

	// 执行解码
	entry, err := Decode(tamperedData)

	// 验证
	assert.ErrorIs(t, err, ErrCRCValidation, "If the CRC check fails, the corresponding error should be returned")
	assert.Nil(t, entry, "If the CRC check fails, nil should be returned.")
}

func TestDecode_ZeroValueSize(t *testing.T) {
	// 测试值长度为0的情况（非删除，显式存储vsz=0）
	key := []byte("zeroValKey")
	val := []byte("") // 空值但非删除
	encoded := Encode(key, val, false)

	// 执行解码
	entry, err := Decode(encoded)

	// 验证
	require.NoError(t, err, "Decoding should not fail when value length is 0")
	assert.Equal(t, key, entry.Key, "Key decoding error")
	assert.Empty(t, entry.Val, "Value should be empty")
	assert.Equal(t, uint32(0), binary.BigEndian.Uint32(encoded[bufKszEndIdx:bufVszEndIdx]),
		"vsz should be stored as 0")
}

func TestDecode_LargeData(t *testing.T) {
	// 测试大键大值的解码
	largeKey := make([]byte, 1024*5) // 5KB键
	for i := range largeKey {
		largeKey[i] = byte(i % 256)
	}
	largeVal := make([]byte, 1024*100) // 100KB值
	for i := range largeVal {
		largeVal[i] = byte((i + 128) % 256)
	}

	encoded := Encode(largeKey, largeVal, false)
	entry, err := Decode(encoded)

	require.NoError(t, err, "Decoding big data should not go wrong")
	assert.Equal(t, largeKey, entry.Key, "Big key decoding mismatch")
	assert.Equal(t, largeVal, entry.Val, "Big value decoding mismatch")
}
