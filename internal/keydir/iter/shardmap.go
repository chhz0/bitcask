package iter

// type shardIterator struct {
// 	sm       *ShardMap
// 	shardIdx int
// 	keys     []string
// 	cursor   int
// 	sync.RWMutex
// }

// // Rewind implements Iterator.
// func (sit *shardIterator) Rewind() {
// 	sit.Lock()
// 	defer sit.Unlock()
// 	sit.shardIdx = 0
// 	sit.cursor = -1
// 	sit.loadKeys()
// }

// // Valid implements Iterator.
// func (sit *shardIterator) Valid() bool {
// 	return sit.cursor < len(sit.keys)-1 && sit.shardIdx < len(sit.sm.shareds)-1
// }

// // Next implements Iterator.
// func (sit *shardIterator) Next() {
// 	sit.Lock()
// 	defer sit.Unlock()

// 	if sit.cursor < len(sit.keys)-1 {
// 		sit.cursor++
// 	} else {
// 		sit.shardIdx++
// 		sit.cursor = 0
// 		sit.loadKeys()
// 	}
// }

// // Value implements Iterator.
// func (sit *shardIterator) Value() *entry.Index {
// 	shard := sit.sm.shareds[sit.shardIdx]

// 	entry, ok := shard.load(sit.keys[sit.cursor])
// 	if !ok {
// 		return nil
// 	}
// 	return entry
// }

// // Release implements Iterator.
// func (sit *shardIterator) Release() {
// 	sit.sm = nil
// 	sit.shardIdx = 0
// 	sit.keys = nil
// 	sit.cursor = -1
// }

// // Key implements Iterator.
// func (sit *shardIterator) Key() []byte {
// 	if sit.cursor >= 0 && sit.cursor < len(sit.keys) {
// 		return []byte(sit.keys[sit.cursor])
// 	}
// 	return nil
// }

// func (sit *shardIterator) loadKeys() {
// 	shard := sit.sm.shareds[sit.shardIdx]

// 	sit.keys = make([]string, 0, shard.ssize())
// 	shard.forange(func(key string, value *entry.Index) bool {
// 		sit.keys = append(sit.keys, key)
// 		return true
// 	})
// }
