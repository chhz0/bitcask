package index

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/chhz0/go-bitcask/internal/entry"
	"github.com/stretchr/testify/assert"
)

func TestShardMap_CURD(t *testing.T) {
	sm := NewShardMap(32, defaultHasher)
	key := []byte("exist-key")
	expectedEntry := &entry.Index{
		Fid:    1,
		Offset: 2,
		Size:   3,
	}

	sm.Put(key, expectedEntry)
	entry, ok := sm.Get(key)
	assert.Equal(t, expectedEntry, entry)
	assert.True(t, ok)

	entry, ok = sm.Get([]byte("not-exist-key"))
	assert.Nil(t, entry)
	assert.False(t, ok)

	entry, ok = sm.Del(key)
	assert.Equal(t, expectedEntry, entry)
	assert.True(t, ok)

	entry, ok = sm.Del(key)
	assert.Nil(t, entry)
	assert.False(t, ok)

	assert.Equal(t, 0, sm.Size())
}

const (
	numItems  = 1_000_000 // 测试数据量
	keySize   = 32        // key长度
	valueSize = 1024
)

// 初始化测试数据
func initShardMap(n int) *ShardMap {
	m := NewShardMap(32, defaultHasher).(*ShardMap)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		m.Put([]byte(key), &entry.Index{})
	}
	return m
}

// 基准测试：纯写入性能
func BenchmarkShardMap_Write(b *testing.B) {
	m := NewShardMap(32, defaultHasher).(*ShardMap)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(rand.Int63()))
		for pb.Next() {
			key := generateKey(r)
			m.Put(key, &entry.Index{
				// FileID:    1,
				// Offset:    1,
				// ValueSize: valueSize,
				// Timestamp: uint64(time.Now().UnixNano()),
			})
		}
	})
}

// 基准测试：纯读取性能
func BenchmarkShardMap_Read(b *testing.B) {
	m := initShardMap(1000000) // 预填充1百万数据
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i % 1000000)
			_, _ = m.Get([]byte(key))
			i++
		}
	})
}

// 基准测试：读写混合（50%读+50%写）
func BenchmarkShardMap_RW50(b *testing.B) {
	m := initShardMap(1000000)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				key := strconv.Itoa(i % 1000000)
				_, _ = m.Get([]byte(key))
			} else {
				key := strconv.Itoa(1000000 + i)
				m.Put([]byte(key), &entry.Index{})
			}
			i++
		}
	})
}

// 测试热点分片场景
func BenchmarkShardMap_HotShard(b *testing.B) {
	// 自定义哈希函数使所有Key命中同一分片
	hotHasher := func(string) uint32 { return 0 }
	m := NewShardMap(32, hotHasher).(*ShardMap)

	b.Run("Write-Hot", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := strconv.Itoa(i)
				m.Put([]byte(key), &entry.Index{})
				i++
			}
		})
	})

	b.Run("Read-Hot", func(b *testing.B) {
		initShardMapWithHasher(1000000, hotHasher) // 预填充
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := strconv.Itoa(i % 1000000)
				_, _ = m.Get([]byte(key))
				i++
			}
		})
	})
}

// 辅助函数：带自定义哈希的初始化
func initShardMapWithHasher(n int, hasher func(string) uint32) *ShardMap {
	m := NewShardMap(32, hasher).(*ShardMap)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		m.Put([]byte(key), &entry.Index{})
	}
	return m
}

// 测试不同分片数的影响
func BenchmarkShardMap_ShardCount(b *testing.B) {
	shardCounts := []int{8, 32, 64, 128, 256}

	for _, count := range shardCounts {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			// 临时修改分片数（需调整代码可见性）
			// origShards := shardCount
			// shardCount = count
			// defer func() { shardCount = origShards }()

			m := NewShardMap(count, defaultHasher).(*ShardMap)
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					if i%2 == 0 {
						key := strconv.Itoa(i)
						m.Put([]byte(key), &entry.Index{})
					} else {
						key := strconv.Itoa(i % 1000000)
						_, _ = m.Get([]byte(key))
					}
					i++
				}
			})
		})
	}
}

// 辅助函数：生成随机键
func generateKey(r *rand.Rand) []byte {
	b := make([]byte, keySize)
	r.Read(b)
	return b
}
