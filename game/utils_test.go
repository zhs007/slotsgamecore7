package sgc7game

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
