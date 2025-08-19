package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_RandWithWeights_Jules(t *testing.T) {
	plugin := sgc7plugin.NewMockPlugin()

	// Test case 1: Basic test
	plugin.SetCache([]int{10})
	idx, err := RandWithWeights(plugin, 100, []int{20, 30, 50})
	assert.NoError(t, err)
	assert.Equal(t, 0, idx)

	plugin.SetCache([]int{25})
	idx, err = RandWithWeights(plugin, 100, []int{20, 30, 50})
	assert.NoError(t, err)
	assert.Equal(t, 1, idx)

	plugin.SetCache([]int{75})
	idx, err = RandWithWeights(plugin, 100, []int{20, 30, 50})
	assert.NoError(t, err)
	assert.Equal(t, 2, idx)

	// Test case 2: Edge case - empty weights
	idx, err = RandWithWeights(plugin, 100, []int{})
	assert.Error(t, err)
	assert.Equal(t, -1, idx)
	assert.Equal(t, ErrInvalidWeights, err)

	// Test case 3: Edge case - max is 0
	idx, err = RandWithWeights(plugin, 0, []int{20, 30, 50})
	assert.Error(t, err)
	assert.Equal(t, -1, idx)
	assert.Equal(t, ErrInvalidWeights, err)
}

func Test_RandList_Jules(t *testing.T) {
	plugin := sgc7plugin.NewMockPlugin()

	// Test case 1: Basic test
	// The mock random function returns the value from the cache modulo the range.
	// So to get index 2 from a list of 5, we need a value that gives 2 when modulo 5. Let's just use 2.
	// Then the list has 4 elements, to get index 0, we need a value that gives 0 mod 4. Let's use 0.
	plugin.SetCache([]int{2, 0})
	lst, err := RandList(plugin, []int{1, 2, 3, 4, 5}, 2)
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 1}, lst)

	// Test case 2: num >= len(arr)
	lst, err = RandList(plugin, []int{1, 2, 3}, 3)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, lst)

	lst, err = RandList(plugin, []int{1, 2, 3}, 4)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, lst)

	// Test case 3: Edge case - empty list
	lst, err = RandList(plugin, []int{}, 2)
	assert.Error(t, err)
	assert.Nil(t, lst)
	assert.Equal(t, ErrInvalidParam, err)

	// Test case 4: Edge case - num is 0
	lst, err = RandList(plugin, []int{1, 2, 3}, 0)
	assert.Error(t, err)
	assert.Nil(t, lst)
	assert.Equal(t, ErrInvalidParam, err)
}
