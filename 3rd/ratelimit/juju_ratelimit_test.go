package ratelimit

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/juju/ratelimit"
	"github.com/stretchr/testify/require"
)

func TestJujuRateLimit(t *testing.T) {

	bucket := ratelimit.NewBucketWithRate(10, 10)
	start := time.Now()
	bucket.Wait(10)
	// wait time 为 0ms
	log.Println("wait: ", time.Now().Sub(start))

	// 等待时间为 1000ms
	log.Println(bucket.Take(10))

	// 超过 capacity 不会报错，等待时间为 2200ms
	log.Println(bucket.Take(12))

	require.False(t, false)

	// 极端情况下，write buf 的长度远大于 cap，里面的直接 sleep 是不合理的，这个出现在 cap 比较小的情况，刚好 write buf 数据比较大。
	// 里面使用了 wait => Take => sleep。
	w := ratelimit.Writer(io.Discard, bucket)
	// return w.w.Write(buf)
	log.Println(w)

}
