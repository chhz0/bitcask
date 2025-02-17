package index

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShardMap_CURD(t *testing.T) {
	sm := NewShardMap(32, defaultHasher)
	key := []byte("exist-key")
	expectedEntry := &Entry{
		FileID:    1,
		Offset:    2,
		ValueSize: 3,
		Timestamp: uint64(time.Now().UnixNano()),
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
