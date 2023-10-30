package slice

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDemo(t *testing.T) {
	type foo struct {
		a, b int
	}

	arr := []foo{
		{
			1, 3,
		},
		{
			5, 6,
		},
	}

	for i := range arr {
		arr[i].a += 3
		arr[i].b = 10
	}

	require.EqualValues(t, []foo{
		{4, 10},
		{8, 10},
	}, arr)
}
