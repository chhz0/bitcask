package index

import (
	"bytes"
	"sync"

	gbtree "github.com/google/btree"
)

// 新增对象池优化（放在文件底部）
var gbtreeItemPool = sync.Pool{
	New: func() interface{} {
		return new(gbtreeItem)
	},
}

func getGBtreeItem(key []byte, val *Index) *gbtreeItem {
	item := gbtreeItemPool.Get().(*gbtreeItem)
	item.key = bytes.Clone(key) // 深拷贝保证key安全
	item.value = val
	return item
}

func putGBtreeItem(item *gbtreeItem) {
	item.key = nil
	item.value = nil
	gbtreeItemPool.Put(item)
}

// 实现google/gbtree需要的Item接口
type gbtreeItem struct {
	key   []byte
	value *Index
}

func (bi gbtreeItem) Less(than gbtree.Item) bool {
	other := than.(*gbtreeItem)
	return bytes.Compare(bi.key, other.key) == -1
}

// Btree B树索引实现
type Btree struct {
	degree int // B树度数
	tree   *gbtree.BTree
	lock   sync.RWMutex
}

// // Close implements Indexer.
// func (bt *btree) Close() error {
// 	panic("unimplemented")
// }

func newBTree(degree int) *Btree {
	return &Btree{
		degree: degree,
		tree:   gbtree.New(degree),
	}
}

// Get 查找key
func (bt *Btree) Get(key []byte) (*Index, bool) {
	bt.lock.RLock()
	defer bt.lock.RUnlock()

	item := bt.tree.Get(&gbtreeItem{key: key})
	if item == nil {
		return nil, false
	}

	if it, ok := item.(*gbtreeItem); ok {
		return it.value, true
	}

	return nil, false
}

// Put 插入/更新
func (bt *Btree) Put(key []byte, value *Index) {
	bt.lock.Lock()
	defer bt.lock.Unlock()

	item := getGBtreeItem(key, value)
	defer putGBtreeItem(item)

	_ = bt.tree.ReplaceOrInsert(item)
}

// Del 删除key
func (bt *Btree) Del(key []byte) (*Index, bool) {
	bt.lock.Lock()
	defer bt.lock.Unlock()

	if item := bt.tree.Delete(&gbtreeItem{key: key}); item != nil {
		return item.(*gbtreeItem).value, true
	}
	return nil, false
}

// Size 数据量
func (bt *Btree) Size() int {
	return bt.tree.Len()
}

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
