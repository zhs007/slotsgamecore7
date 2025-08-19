package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_NewArrValWeights_Jules(t *testing.T) {
	// Test case 1: Valid input
	arrvals1 := [][]int{{1, 2}, {3, 4}}
	weights1 := []int{10, 20}
	avw1, err1 := NewArrValWeights(arrvals1, weights1)
	assert.NoError(t, err1)
	assert.NotNil(t, avw1)
	assert.Equal(t, 30, avw1.MaxWeight)
	assert.Equal(t, 2, len(avw1.ArrVals))
	assert.Equal(t, 2, len(avw1.Weights))

	// Test case 2: Mismatched lengths
	arrvals2 := [][]int{{1, 2}}
	weights2 := []int{10, 20}
	_, err2 := NewArrValWeights(arrvals2, weights2)
	assert.Error(t, err2)
	assert.Equal(t, ErrInvalidValWeights, err2)
}

func Test_ArrValWeights_RandVal_Jules(t *testing.T) {
	// Test case: RandVal
	arrvals := [][]int{{1, 2}, {3, 4}}
	weights := []int{10, 20}
	avw, err := NewArrValWeights(arrvals, weights)
	assert.NoError(t, err)
	assert.NotNil(t, avw)

	plugin := sgc7plugin.NewMockPlugin()

	// Mock RandWithWeights to return the first index
	plugin.Cache = []int{0}
	val, err := avw.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2}, val)

	// Mock RandWithWeights to return the second index
	plugin.Cache = []int{10}
	val, err = avw.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 4}, val)
}
