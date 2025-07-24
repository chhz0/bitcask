package codec

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test Case gen by Doubao AI
func TestEncode_NormalCase(t *testing.T) {
	// 准备测试数据
	key := []byte("testKey")
	val := []byte("testValue")
	expectedTstamp := uint64(time.Now().Unix())

	// 执行编码
	result := Encode(key, val, false)

	// 验证总长度
	expectedLen := headerSize + len(key) + len(val)
	assert.Equal(t, expectedLen, len(result),
		"The total length after encoding does not match")

	// 验证时间戳（允许1秒误差，避免时间戳恰好跨秒）
	tstamp := binary.BigEndian.Uint64(result[crcSize:bufTstampEndIdx])
	assert.True(t, tstamp >= expectedTstamp && tstamp <= expectedTstamp+1,
		"Timestamp does not fall within the expected range")

	// 验证键长度
	ksz := binary.BigEndian.Uint32(result[bufTstampEndIdx:bufKszEndIdx])
	assert.Equal(t, uint32(len(key)), ksz,
		"Key length encoding error")

	// 验证值长度
	vsz := binary.BigEndian.Uint32(result[bufKszEndIdx:bufVszEndIdx])
	assert.Equal(t, uint32(len(val)), vsz,
		"Value length encoding error")

	// 验证键内容
	keyStart := bufVszEndIdx
	keyEnd := keyStart + len(key)
	assert.True(t, bytes.Equal(key, result[keyStart:keyEnd]),
		"Key content does not match")

	// 验证值内容
	valStart := keyEnd
	valEnd := valStart + len(val)
	assert.True(t, bytes.Equal(val, result[valStart:valEnd]),
		"Value content does not match")

	// 验证CRC校验
	expectedCRC := calculateCRC(result[crcSize:])
	actualCRC := binary.BigEndian.Uint32(result[:crcSize])
	assert.Equal(t, expectedCRC, actualCRC,
		"CRC checksum value does not match")
}

func TestEncode_DeletedKey(t *testing.T) {
	// 准备测试数据
	key := []byte("deletedKey")
	val := []byte("ignoredValue") // deleted=true时val应被忽略

	// 执行编码（标记为删除）
	result := Encode(key, val, true)

	// 验证总长度（删除时不含值内容）
	expectedLen := headerSize + len(key)
	assert.Equal(t, expectedLen, len(result),
		"Wrong total length when deleting tags")

	// 验证值长度被标记为-1（uint32存储的是补码，实际读取时需处理）
	vsz := binary.BigEndian.Uint32(result[bufKszEndIdx:bufVszEndIdx])
	assert.Equal(t, uint32(0x00), vsz,
		"The value length should be 0 when deleting a mark")

	// 验证值内容不存在
	keyStart := bufVszEndIdx
	keyEnd := keyStart + len(key)
	assert.Equal(t, len(result), keyEnd,
		"The value content should not be included after the removal tag")

	// 验证键内容正确
	assert.True(t, bytes.Equal(key, result[keyStart:keyEnd]),
		"Wrong key content when deleting a tag")

	// 验证CRC校验
	expectedCRC := calculateCRC(result[crcSize:])
	actualCRC := binary.BigEndian.Uint32(result[:crcSize])
	assert.Equal(t, expectedCRC, actualCRC,
		"CRC check error when deleting mark")
}

func TestEncode_EmptyKeyAndValue(t *testing.T) {
	// 准备测试数据（空键空值）
	key := []byte("")
	val := []byte("")

	// 执行编码
	result := Encode(key, val, false)

	// 验证总长度（仅包含头部）
	assert.Equal(t, headerSize, len(result),
		"Total length error when key value is empty")

	// 验证键长度为0
	ksz := binary.BigEndian.Uint32(result[bufTstampEndIdx:bufKszEndIdx])
	assert.Zero(t, ksz,
		"Empty key length should be 0")

	// 验证值长度为0
	vsz := binary.BigEndian.Uint32(result[bufKszEndIdx:bufVszEndIdx])
	assert.Zero(t, vsz,
		"Empty value length should be 0")

	// 验证CRC校验
	expectedCRC := calculateCRC(result[crcSize:])
	actualCRC := binary.BigEndian.Uint32(result[:crcSize])
	assert.Equal(t, expectedCRC, actualCRC,
		"CRC check error when key value is empty")
}

func TestEncode_LargeData(t *testing.T) {
	// 准备大尺寸测试数据
	key := bytes.Repeat([]byte("k"), 1024*10)  // 10KB键
	val := bytes.Repeat([]byte("v"), 1024*100) // 100KB值

	// 执行编码
	result := Encode(key, val, false)

	// 验证总长度
	expectedLen := headerSize + len(key) + len(val)
	assert.Equal(t, expectedLen, len(result),
		"Total length error when using large data")

	// 验证键长度
	ksz := binary.BigEndian.Uint32(result[bufTstampEndIdx:bufKszEndIdx])
	assert.Equal(t, uint32(len(key)), ksz,
		"Large key length encoding error")

	// 验证值长度
	vsz := binary.BigEndian.Uint32(result[bufKszEndIdx:bufVszEndIdx])
	assert.Equal(t, uint32(len(val)), vsz,
		"Large value length encoding error")

	// 验证键内容
	keyStart := bufVszEndIdx
	keyEnd := keyStart + len(key)
	assert.True(t, bytes.Equal(key, result[keyStart:keyEnd]),
		"Key content does not match")

	// 验证值内容
	valStart := keyEnd
	valEnd := valStart + len(val)
	assert.True(t, bytes.Equal(val, result[valStart:valEnd]),
		"Large value content does not match")
}
