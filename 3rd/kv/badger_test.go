package kv

import (
	"testing"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
)

func TestBadgerTxn(t *testing.T) {
	// https://note.iawen.com/note/graph/badger_base
	db, err := badger.Open(badger.DefaultOptions("./bg.db"))
	require.NoError(t, err)

	// 支持读写事务、支持 iter 读取，支持 batch。
	// 支持 stream
	// 单独读取更友好一点，原因 key val 分离。
	// 读写事务。
	key, val := []byte("bar"), []byte("bar")
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(key, val)
		return err
	})
	require.NoError(t, err)
	gval, err := get(db, key)
	require.NoError(t, err)
	require.Equal(t, val, gval)

	// 设置 ttl 。
	key = []byte("abc")
	val = []byte("123")
	err = db.Update(func(txn *badger.Txn) error {
		entry := badger.NewEntry(key, val).WithMeta(byte(1)).WithTTL(time.Second)
		err := txn.SetEntry(entry)
		return err
	})
	require.NoError(t, err)
	gval, err = get(db, key)
	require.NoError(t, err)
	require.Equal(t, val, gval)

	time.Sleep(time.Second)
	gval, err = get(db, key)
	require.EqualError(t, badger.ErrKeyNotFound, err.Error())
	return

	err = db.Close()
	require.NoError(t, err)
}

func get(db *badger.DB, key []byte) ([]byte, error) {
	// 只读事务。
	var gval []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		gval, err = item.ValueCopy(nil)
		return err
	})
	return gval, err
}
