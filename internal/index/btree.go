package index

import (
	"bytes"
	"sync"

	"github.com/chhz0/go-bitcask/internal/entry"
	gbtree "github.com/google/btree"
)

// 实现google/gbtree需要的Item接口
type gbtreeItem struct {
	key   []byte
	value *entry.Index
}

func (bi gbtreeItem) Less(than gbtree.Item) bool {
	other := than.(*gbtreeItem)
	return bytes.Compare(bi.key, other.key) < 0
}

// btree B树索引实现
type btree struct {
	degree int // B树度数
	tree   *gbtree.BTree
	lock   sync.RWMutex
	size   int
}

// // Close implements Indexer.
// func (bt *btree) Close() error {
// 	panic("unimplemented")
// }

// // Snapshot implements Indexer.
// func (bt *btree) Snapshot() map[string]*Entry {
// 	panic("unimplemented")
// }

// func NewBTree(degree int) Indexer {
// 	return &btree{
// 		degree: degree,
// 		tree:   gbtree.New(degree),
// 	}
// }

// // Get 查找key
// func (bt *btree) Get(key []byte) (*Entry, bool) {
// 	bt.lock.RLock()
// 	defer bt.lock.RUnlock()

// 	item := bt.tree.Get(&gbtreeItem{key: key})
// 	if item == nil {
// 		return nil, false
// 	}
// 	return item.(*gbtreeItem).value, true
// }

// // Put 插入/更新
// func (bt *btree) Put(key []byte, value *Entry) {
// 	bt.lock.Lock()
// 	defer bt.lock.Unlock()

// 	item := &gbtreeItem{key: key, value: value}
// 	if old := bt.tree.ReplaceOrInsert(item); old == nil {
// 		bt.size++
// 	}
// }

// // Del 删除key
// func (bt *btree) Del(key []byte) (*Entry, bool) {
// 	bt.lock.Lock()
// 	defer bt.lock.Unlock()

// 	if item := bt.tree.Delete(&gbtreeItem{key: key}); item != nil {
// 		bt.size--
// 		return item.(*gbtreeItem).value, true
// 	}
// 	return nil, false
// }

// // Size 数据量
// func (bt *btree) Size() int {
// 	bt.lock.RLock()
// 	defer bt.lock.RUnlock()
// 	return bt.size
// }

// // Iterator 实现范围遍历的迭代器
// func (bt *btree) Iterator() Iterator {
// 	return &gbtreeIterator{
// 		bt:     bt,
// 		cursor: -1,
// 		keys:   bt.getAllKeys(),
// 	}
// }

// // 获取所有键（按序）
// func (bt *btree) getAllKeys() [][]byte {
// 	bt.lock.RLock()
// 	defer bt.lock.RUnlock()

// 	keys := make([][]byte, 0, bt.size)
// 	bt.tree.Ascend(func(i gbtree.Item) bool {
// 		keys = append(keys, i.(*gbtreeItem).key)
// 		return true
// 	})
// 	return keys
// }

// // gbtree迭代器实现
// type gbtreeIterator struct {
// 	bt     *btree
// 	cursor int      // 当前索引
// 	keys   [][]byte // 快照键集合
// }

// func (bi *gbtreeIterator) Rewind() {
// 	bi.cursor = -1
// }

// func (bi *gbtreeIterator) Valid() bool {
// 	return bi.cursor >= 0 && bi.cursor < len(bi.keys)
// }

// func (bi *gbtreeIterator) Next() {
// 	bi.cursor++
// }

// func (bi *gbtreeIterator) Key() []byte {
// 	if !bi.Valid() {
// 		return nil
// 	}
// 	return bi.keys[bi.cursor]
// }

// func (bi *gbtreeIterator) Value() *Entry {
// 	if !bi.Valid() {
// 		return nil
// 	}
// 	v, _ := bi.bt.Get(bi.keys[bi.cursor])
// 	return v
// }

// func (bi *gbtreeIterator) Release() {
// 	bi.keys = nil
// 	bi.bt = nil
// }
