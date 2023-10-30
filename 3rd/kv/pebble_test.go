package kv

import (
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/stretchr/testify/require"
)

func TestPebble(t *testing.T) {
	// 不支持事务，事务在 cockroachdb 层做。
	// 支持 iter 读取，支持 range 删除。
	// 范围查找更友好一点。
	db, err := pebble.Open("./pb.db", &pebble.Options{})
	require.NoError(t, err)

	key, val := []byte("foo"), []byte("bar")
	err = db.Set(key, val, &pebble.WriteOptions{Sync: true})
	require.NoError(t, err)
	gval, closer, err := db.Get(key)
	require.NoError(t, err)
	require.Equal(t, val, gval)
	err = closer.Close()
	require.NoError(t, err)

	err = db.Close()
	require.NoError(t, err)
}
