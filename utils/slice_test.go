package sgc7utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IndexOfIntSlice(t *testing.T) {
	ci := IndexOfIntSlice([]int{1, 2, 3}, 3, 0)
	assert.Equal(t, ci, 2)

	ci = IndexOfIntSlice(nil, 3, -1)
	assert.Equal(t, ci, -1)

	ci = IndexOfIntSlice([]int{1, 2, 3}, 3, 5)
	assert.Equal(t, ci, -1)

	ci = IndexOfIntSlice([]int{1, 2, 3}, 3, -100)
	assert.Equal(t, ci, 2)

	t.Logf("Test_IndexOfIntSlice OK")
}

func Test_IndexOfInt2Slice(t *testing.T) {
	ci := IndexOfInt2Slice([]int{1, 2, 3, 4, 5, 6}, 3, 4, 0)
	assert.Equal(t, ci, 1)

	ci = IndexOfInt2Slice(nil, 3, 4, -1)
	assert.Equal(t, ci, -1)

	ci = IndexOfInt2Slice([]int{1, 2, 3, 4, 5, 6}, 3, 4, 5)
	assert.Equal(t, ci, -1)

	ci = IndexOfInt2Slice([]int{1, 2, 3, 4, 5, 6}, 3, 4, -100)
	assert.Equal(t, ci, 1)

	ci = IndexOfInt2Slice([]int{1, 2, 3, 4, 5, 6}, 2, 3, -100)
	assert.Equal(t, ci, -1)

	t.Logf("Test_IndexOfInt2Slice OK")
}

func Test_IndexOfStringSlice(t *testing.T) {
	ci := IndexOfStringSlice([]string{"1", "2", "3"}, "3", 0)
	assert.Equal(t, ci, 2)

	ci = IndexOfStringSlice(nil, "3", -1)
	assert.Equal(t, ci, -1)

	ci = IndexOfStringSlice([]string{"1", "2", "3"}, "3", 5)
	assert.Equal(t, ci, -1)

	ci = IndexOfStringSlice([]string{"1", "2", "3"}, "3", -100)
	assert.Equal(t, ci, 2)

	t.Logf("Test_IndexOfStringSlice OK")
}

func Test_InsUniqueIntSlice(t *testing.T) {
	arr := InsUniqueIntSlice([]int{1, 2, 3}, 3)
	assert.Equal(t, len(arr), 3)

	arr = InsUniqueIntSlice([]int{}, 3)
	assert.Equal(t, len(arr), 1)

	arr = InsUniqueIntSlice(nil, 3)
	assert.Equal(t, len(arr), 1)

	arr = InsUniqueIntSlice([]int{1, 2, 3}, 4)
	assert.Equal(t, len(arr), 4)

	t.Logf("Test_InsUniqueIntSlice OK")
}

func Test_IntArr2ToInt32Arr(t *testing.T) {
	arr := IntArr2ToInt32Arr([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	assert.Equal(t, len(arr), 9)
	assert.Equal(t, arr[0], int32(1))
	assert.Equal(t, arr[1], int32(2))
	assert.Equal(t, arr[2], int32(3))
	assert.Equal(t, arr[3], int32(4))
	assert.Equal(t, arr[4], int32(5))
	assert.Equal(t, arr[5], int32(6))
	assert.Equal(t, arr[6], int32(7))
	assert.Equal(t, arr[7], int32(8))
	assert.Equal(t, arr[8], int32(9))

	t.Logf("Test_IntArr2ToInt32Arr OK")
}
