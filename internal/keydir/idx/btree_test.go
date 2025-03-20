package idx

import (
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkBtree_Write(b *testing.B) {
	bt := NewBTree(32)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		r := rand.New(rand.NewSource(rand.Int63()))
		for pb.Next() {
			key := generateKey(r)
			bt.Put(key, &Index{
				FileID:  1,
				ValSize: valueSize,
				ValPos:  1,
			})
		}
	})
}

func initBtree(n int) *Btree {
	bt := NewBTree(32)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		bt.Put([]byte(key), &Index{})
	}
	return bt
}

func BenchmarkBtree_Read(b *testing.B) {
	bt := initBtree(numItems)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := strconv.Itoa(i % numItems)
			_, _ = bt.Get([]byte(key))
			i++
		}
	})
}
