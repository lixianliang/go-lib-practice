package ratelimit

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"
)

func TestXRate(t *testing.T) {
	limiter := rate.NewLimiter(0, 0)
	require.EqualValues(t, 0, int(limiter.Tokens()))

	// limit = 10, burst = 10
	limiter = rate.NewLimiter(10, 10)

	ok := limiter.AllowN(time.Now(), 4)
	require.True(t, ok)
	require.EqualValues(t, 6, int64(limiter.Tokens()))

	ok = limiter.AllowN(time.Now(), 5)
	require.True(t, ok)
	// token 只剩 1。
	require.EqualValues(t, 1, int(limiter.Tokens()))

	// 900 ms token = 10
	require.EqualValues(t, 10, int64(limiter.TokensAt(time.Now().Add(900*time.Millisecond))))
	// 1000 ms token = 10, burst = 10
	require.EqualValues(t, 10, int64(limiter.TokensAt(time.Now().Add(time.Second))))

	// n 不要大于 burst
	ok = limiter.AllowN(time.Now().Add(time.Second), 11)
	require.False(t, ok)

	ok = limiter.AllowN(time.Now().Add(time.Second), 10)
	require.True(t, ok)
	// 当前 token = 0
	require.EqualValues(t, 0, int64(limiter.Tokens()))

	// 1.1s 后 max(token) = 11
	WrapSetBurst(limiter, 11)
	require.EqualValues(t, 11, int64(limiter.TokensAt(time.Now().Add(time.Second+100*time.Millisecond))))
	require.EqualValues(t, 11, int64(limiter.TokensAt(time.Now().Add(time.Second+200*time.Millisecond))))

	// 当前不允许通过
	ok = limiter.AllowN(time.Now(), 1)
	require.False(t, ok)

	// 请求的 n 也要小于容量
	reserve := limiter.ReserveN(time.Now(), 15)
	require.False(t, reserve.OK())
}

func TestXRateWaitTime(t *testing.T) {
	ctx := context.Background()
	limiter := rate.NewLimiter(10, 10)
	err := limiter.WaitN(ctx, 10)
	require.NoError(t, err)

	// 此时 token = 0
	require.EqualValues(t, 0, int(limiter.Tokens()))

	// wait 100ms
	err = WrapWait(limiter, 1, "1")
	require.NoError(t, err)

	// wait 500ms
	err = WrapWait(limiter, 5, "2")
	require.NoError(t, err)

	// error, waitN <= burst
	err = WrapWait(limiter, 15, "3")
	require.Error(t, err)

	// wait 1000ms
	err = WrapWait(limiter, 10, "4")
	require.NoError(t, err)
}

func TestXRateWaitTime4Concurrency(t *testing.T) {
	ctx := context.Background()
	limiter := rate.NewLimiter(10, 10)
	err := limiter.WaitN(ctx, 10)
	require.NoError(t, err)
	// 此时 token = 0
	require.EqualValues(t, 0, int(limiter.Tokens()))

	// 如果有多个并发 wait，
	// 从日志看，如果 waitN <= burst/2, 则优先选择 waitN, 其它 waitN 继续等待，等待时间会为 sum(wait1,...,waitN)
	// 如果所有等待都 waitN >= burst/2，则由小到大选择
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		// wait 500ms + 600ms
		err = WrapWait(limiter, 5, "1")
		require.NoError(t, err)
	}()

	// 如果这里 sleep，会严格按顺序的，这里tokens 是个负数，如果不做sleep  token是个小数，说明上面的 1 也还没到 wait阶段。
	time.Sleep(time.Millisecond)
	log.Printf("tokens: %f\n", limiter.Tokens())

	wg.Add(1)
	go func() {
		defer wg.Done()

		// wait 600ms
		err = WrapWait(limiter, 6, "2")
		require.NoError(t, err)
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = WrapWait(limiter, 9, "3")
		require.NoError(t, err)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err = WrapWait(limiter, 8, "4")
		require.NoError(t, err)
	}()

	// 等待一段时间扩大 limiter，已经等待的 waitN 则按之前 limiter 继续运行。
	time.Sleep(100 * time.Millisecond)
	wg.Add(1)
	go func() {
		defer wg.Done()

		ResetLimiter(limiter, 20)
	}()

	// 这里总等待时间为 800ms + 900ms - 100ms + 500ms = 2600ms，相差 500ms？
	// 这里看上去还是继续在原来的处理速度处理 wait，貌似自己测试这里也有比较小的概率按新的limiter处理，需要看 ResetLimiter 处理时间。
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = WrapWait(limiter, 10, "5")
		require.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)
	// 总时间为 1950ms =  800ms + 900ms - 200ms  + 750ms = 2250ms，相差 300ms，这里生效的时间应该是在 id1 id2 完成之前。
	wg.Add(1)
	go func() {
		defer wg.Done()

		err = WrapWait(limiter, 15, "6")
		require.NoError(t, err)
	}()

	wg.Wait()

	// 总结
	// write 直接先查看 burst ，需要通过 for 循环write
	// 若 waitN 时代，需要做一个 退化 write 处理
}

func TestXRateLimitZero(t *testing.T) {
	// limit = 10, burst = 10
	limiter := rate.NewLimiter(10, 10)

	ok := limiter.AllowN(time.Now(), 5)
	require.True(t, ok)
	require.EqualValues(t, 5, int64(limiter.Tokens()))

	// 可用 token 重置后为 0.
	limiter.SetBurst(0)
	require.EqualValues(t, 0, int64(limiter.Tokens()))

	// wait 会失败。
	err := limiter.WaitN(context.Background(), 1)
	require.Error(t, err)
}

func WrapWait(limiter *rate.Limiter, n int, id string) error {
	start := time.Now()

	defer func() {
		log.Printf("id: %s, wait time: %v\n", id, time.Now().Sub(start))
	}()

	return limiter.WaitN(context.Background(), n)
}

func WrapSetBurst(limiter *rate.Limiter, n int) {
	start := time.Now()

	defer func() {
		log.Printf("SetBurst wait time: %v\n", time.Now().Sub(start))
	}()

	limiter.SetBurst(n)
}

func ResetLimiter(limiter *rate.Limiter, n int) {
	start := time.Now()
	limiter.SetBurst(n)
	startLimit := time.Now()
	limiter.SetLimit(rate.Limit(n))

	end := time.Now()
	log.Printf("ResetLimiter wait time: %v, %v:%v\n", end.Sub(start), startLimit.Sub(start), end.Sub(startLimit))
}
