package ec

import (
	"crypto/md5"
	"fmt"
	"testing"

	"github.com/klauspost/reedsolomon"
	"github.com/stretchr/testify/require"
)

func TestReedSolomon(t *testing.T) {
	N, M := 6, 2
	enc, err := reedsolomon.New(N, M)
	require.NoError(t, err)
	enc2, err := reedsolomon.New(N, M)
	require.NoError(t, err)

	data := "0123456789abcdefghij=xyz"
	// 将数据切分为 6 块，同时生成 2 个校验块。
	shards, err := enc.Split([]byte(data))
	require.NoError(t, err)
	require.Len(t, shards, 8)
	shards2, err := enc2.Split([]byte(data))
	require.NoError(t, err)
	require.Len(t, shards2, 8)

	// enc.Join 支持 io.Writer，将 shard 数据块拼装写入 writer
	// newdatashards 设置需要更新的数据块，重新进行编码。
	// r.Update(shards, newdatashards)
	// ReconstructSome 可以设置需要恢复的块。

	// EC 编码。
	err = enc.Encode(shards)
	require.NoError(t, err)
	// 校验。
	ok, err := enc.Verify(shards)
	require.NoError(t, err)
	require.True(t, ok)

	var srcShards [8][]byte
	var crc32s [8][16]byte
	for i, shard := range shards {
		srcShards[i] = shard
		crc32s[i] = md5.Sum(shard)
	}
	for i, shard := range srcShards {
		if i < N {
			require.Equal(t, data[i*4:(i+1)*4], string(shard))
		} else {
			fmt.Println(string(shard))
		}
	}

	// 增量进行 EC 编码，校验块的数据每次都会变化。
	// 适合在条带 EC 直写，可减少内存使用，也可加快数据写入过程。
	for i, shard := range shards2[:6] {
		err = enc2.EncodeIdx(shard, i, shards2[6:])
		require.NoError(t, err)
	}
	for i, shard := range shards2 {
		require.Equal(t, srcShards[i], shard)
		require.Equal(t, crc32s[i], md5.Sum(shard))
	}

	cases := []struct {
		a, b int
	}{
		{0, 2},
		{0, 5},
		{3, 4},
		{0, 6},
		{1, 7},
	}

	for _, cs := range cases {
		// 设置块数据丢失。
		shards[cs.a] = nil
		shards[cs.b] = nil
		// 只恢复数据块。
		// err = enc.ReconstructData(shards)
		// 恢复数据和校验块。
		err = enc.Reconstruct(shards)
		require.NoError(t, err)
		// 如果数据块数据变动，虽然可以恢复但校验这步会失败。
		ok, err = enc.Verify(shards)
		require.NoError(t, err)
		require.True(t, ok)
		for i, shard := range shards {
			//fmt.Println(i, cs, string(srcShards[i]))
			require.Equal(t, srcShards[i], shard)
			require.Equal(t, crc32s[i], md5.Sum(shard))
		}
	}

	// 恢复失败。
	shards[0] = nil
	shards[6] = nil
	shards[7] = nil
	err = enc.Reconstruct(shards)
	// too few shards given
	require.Error(t, err)
}
