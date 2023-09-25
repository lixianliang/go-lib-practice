package sync

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAtomic(t *testing.T) {
	var n atomic.Int64
	var wg sync.WaitGroup

	// go.19 开始，新增了 Int64  Int32 Uint64 Uint32 Uintptr Bool Pointer 原子类型类型。
	// 以 atomic.Int64 为例就是对原来 atomic.AddInt64 atomic.LoadInt64 atomic.StoreInt64 atomic.SwapInt64 atomic.CompareAndSwapIn6 做了一层类型封装，更好用了。
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			n.Add(1)
		}()
	}
	wg.Wait()

	require.EqualValues(t, 1000, n.Load())
}

func TestPointer(t *testing.T) {
	type tempT struct {
		name string
		age  int
	}

	// 新的方式没有 unsafe 包的引入。
	var pt atomic.Pointer[tempT]
	temp := tempT{"name", 10}
	pt.Store(&temp)
	require.Equal(t, &temp, pt.Load())

	temp1 := tempT{"name1", 20}
	oldTemp := pt.Swap(&temp1)
	require.Equal(t, &temp, oldTemp)
	require.Equal(t, &temp1, pt.Load())

	ok := pt.CompareAndSwap(&temp, &temp1)
	require.False(t, ok)
	ok = pt.CompareAndSwap(&temp1, &temp)
	require.True(t, ok)
}

func TestValue(t *testing.T) {
	type tempT struct {
		name string
		age  int
	}

	// Value 类型都是以 interface{} 做参数。
	var value atomic.Value
	require.Panics(t, func() {
		_ = value.Load().(tempT)
	}, "load value is nil")

	// 直接 type 对应的类型会 panic，value 里面的值为 nil。
	val := value.Load()
	_, ok := val.(tempT)
	require.False(t, ok)

	// 存储一个 nil 值会 panic。
	require.Panics(t, func() {
		value.Store(nil)
	}, "store value is nil")

	// Store 后续的参数类型要与之前保持一致，否则会 panic。
	var temp tempT
	value.Store(temp)
	require.Panics(t, func() {
		value.Store("abc")
	}, "store value not match")

	temp1 := tempT{"name", 20}
	old := value.Swap(temp1)
	require.Equal(t, temp, old)
	require.Equal(t, temp1, value.Load().(tempT))

	ok = value.CompareAndSwap(temp, temp1)
	require.False(t, ok)
	ok = value.CompareAndSwap(temp1, temp)
	require.True(t, ok)
	require.Equal(t, temp, value.Load().(tempT))
}
