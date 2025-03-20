package idx

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShardMap_CURD(t *testing.T) {
	sm := NewShardMap(defaultShardCount, maphashFn)
	key := []byte("exist-key")
	expectedEntry := &Index{
		FileID:  1,
		ValPos:  2,
		ValSize: 3,
	}

	sm.Put(key, expectedEntry)
	entry, ok := sm.Get(key)
	assert.Equal(t, expectedEntry, entry)
	assert.True(t, ok)

	entry, ok = sm.Get([]byte("not-exist-key"))
	assert.Nil(t, entry)
	assert.False(t, ok)

	assert.Equal(t, int64(1), sm.Size())

	entry, ok = sm.Del(key)
	assert.Equal(t, expectedEntry, entry)
	assert.True(t, ok)

	entry, ok = sm.Del(key)
	assert.Nil(t, entry)
	assert.False(t, ok)

	assert.Equal(t, int64(0), sm.Size())
}

const (
	numItems  = 1_000_000 // 测试数据量
	keySize   = 32        // key长度
	valueSize = 1024
)

// 初始化测试数据
func initShardMap(n int) *ShardMap {
	m := NewShardMap(defaultShardCount, maphashFn)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		m.Put([]byte(key), &Index{})
	}
	return m
}

// 辅助函数：生成随机键
func generateKey(r *rand.Rand) []byte {
	b := make([]byte, keySize)
	r.Read(b)
	return b
}

// BenchmarkShardMap_Write 基准测试：纯写入性能
func BenchmarkShardMap_Write(b *testing.B) {
	m := NewShardMap(defaultShardCount, maphashFn)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(rand.Int63()))
		for pb.Next() {
			key := generateKey(r)
			m.Put(key, &Index{
				FileID:  1,
				ValSize: valueSize,
				ValPos:  1,
			})
		}
	})
}

// BenchmarkShardMap_Read 基准测试：纯读取性能
func BenchmarkShardMap_Read(b *testing.B) {
	m := initShardMap(numItems)
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

// BenchmarkShardMap_RW50 基准测试：读写混合（50%读+50%写）
func BenchmarkShardMap_RW50(b *testing.B) {
	m := initShardMap(numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%2 == 0 {
				key := strconv.Itoa(i % numItems)
				_, _ = m.Get([]byte(key))
			} else {
				key := strconv.Itoa(numItems + i)
				m.Put([]byte(key), &Index{})
			}
			i++
		}
	})
}

// BenchmarkShardMap_HotShard 测试热点分片场景
func BenchmarkShardMap_HotShard(b *testing.B) {
	// 自定义哈希函数使所有Key命中同一分片
	hotHasher := func(string) uint64 { return 0 }
	m := NewShardMap(defaultShardCount, maphashFn)

	b.Run("Write-Hot", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := strconv.Itoa(i)
				m.Put([]byte(key), &Index{})
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
func initShardMapWithHasher(n int, hasher func(string) uint64) *ShardMap {
	m := NewShardMap(defaultShardCount, hasher)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		m.Put([]byte(key), &Index{})
	}
	return m
}

// BenchmarkShardMap_ShardCount 测试不同分片数的影响
func BenchmarkShardMap_ShardCount(b *testing.B) {
	shardCounts := []int{32, 64, 128, 256, 512, 1024}

	for _, count := range shardCounts {
		b.Run(strconv.Itoa(count), func(b *testing.B) {
			// 临时修改分片数（需调整代码可见性）
			// origShards := shardCount
			// shardCount = count
			// defer func() { shardCount = origShards }()

			m := NewShardMap(count, maphashFn)
			b.RunParallel(func(pb *testing.PB) {
				i := 0
				for pb.Next() {
					if i%2 == 0 {
						key := strconv.Itoa(i)
						m.Put([]byte(key), &Index{})
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
