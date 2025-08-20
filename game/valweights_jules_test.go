package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_ValWeights_SortBy_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	dst, _ := NewValWeights([]int{3, 1, 2}, []int{1, 1, 1})

	err := vw.SortBy(dst)
	assert.NoError(t, err)
	assert.Equal(t, []int{3, 1, 2}, vw.Vals)
	assert.Equal(t, []int{30, 10, 20}, vw.Weights)
	assert.Equal(t, 60, vw.MaxWeight)

	// test error case
	dst2, _ := NewValWeights([]int{3, 1}, []int{1, 1})
	err = vw.SortBy(dst2)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidValWeights, err)

	vw2, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	dst3, _ := NewValWeights([]int{3, 1, 4}, []int{1, 1, 1})
	err = vw2.SortBy(dst3)
	assert.Error(t, err)
}

func Test_ValWeights_GetWeight_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	assert.Equal(t, 10, vw.GetWeight(1))
	assert.Equal(t, 20, vw.GetWeight(2))
	assert.Equal(t, 30, vw.GetWeight(3))
	assert.Equal(t, 0, vw.GetWeight(4))
}

func Test_ValWeights_Add_Jules(t *testing.T) {
	vw := NewValWeightsEx()
	vw.Add(1, 10)
	assert.Equal(t, []int{1}, vw.Vals)
	assert.Equal(t, []int{10}, vw.Weights)
	assert.Equal(t, 10, vw.MaxWeight)
}

func Test_ValWeights_ClearExcludeVal_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	vw.ClearExcludeVal(2)
	assert.Equal(t, []int{2}, vw.Vals)
	assert.Equal(t, []int{1}, vw.Weights)
	assert.Equal(t, 1, vw.MaxWeight)
}

func Test_ValWeights_Reset_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	vw.Reset([]int{4, 5}, []int{40, 50})
	assert.Equal(t, []int{4, 5}, vw.Vals)
	assert.Equal(t, []int{40, 50}, vw.Weights)
	assert.Equal(t, 90, vw.MaxWeight)
}

func Test_ValWeights_Clone_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})
	clone := vw.Clone()
	assert.Equal(t, vw, clone)
	clone.Add(4, 40)
	assert.NotEqual(t, vw, clone)
}

func Test_ValWeights_RandVal_Jules(t *testing.T) {
	plugin := sgc7plugin.NewMockPlugin()
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})

	plugin.Cache = []int{0}
	val, err := vw.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)

	plugin.Cache = []int{10}
	val, err = vw.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 2, val)

	plugin.Cache = []int{30}
	val, err = vw.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 3, val)

	vw2, _ := NewValWeights([]int{1}, []int{10})
	val, err = vw2.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 1, val)
}

func Test_ValWeights_CloneExcludeVal_Jules(t *testing.T) {
	vw, _ := NewValWeights([]int{1, 2, 3}, []int{10, 20, 30})

	nvw, err := vw.CloneExcludeVal(2)
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 3}, nvw.Vals)
	assert.Equal(t, []int{10, 30}, nvw.Weights)
	assert.Equal(t, 40, nvw.MaxWeight)

	_, err = vw.CloneExcludeVal(4)
	assert.Error(t, err)

	vw2, _ := NewValWeights([]int{1}, []int{10})
	_, err = vw2.CloneExcludeVal(1)
	assert.Error(t, err)
}

func Test_NewValWeights_Jules(t *testing.T) {
	vw, err := NewValWeights([]int{1, 2}, []int{10, 20})
	assert.NoError(t, err)
	assert.NotNil(t, vw)

	_, err = NewValWeights([]int{1}, []int{10, 20})
	assert.Error(t, err)
}

func Test_LoadValWeightsFromExcel_Jules(t *testing.T) {
	vw, err := LoadValWeightsFromExcel("./testdata/valweights.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, vw)
	assert.Equal(t, []int{1, 2, 3}, vw.Vals)
	assert.Equal(t, []int{10, 20, 30}, vw.Weights)
	assert.Equal(t, 60, vw.MaxWeight)

	_, err = LoadValWeightsFromExcel("./testdata/nonexistent.xlsx")
	assert.Error(t, err)
}
