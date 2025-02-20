package index

import "strconv"

func initIndexer(indexType IndexType, n int) Indexer {
	indexer := New(HASH)
	for i := 0; i < n; i++ {
		key := strconv.Itoa(i)
		indexer.Put([]byte(key), &Entry{})
	}
	return indexer
}
