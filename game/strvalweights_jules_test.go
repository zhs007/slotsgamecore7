package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_StrValWeights_jules(t *testing.T) {
	// Test NewStrValWeights
	svw1, err := NewStrValWeights([]string{"a", "b"}, []int{10, 20})
	assert.NoError(t, err)
	assert.NotNil(t, svw1)
	assert.Equal(t, 30, svw1.MaxWeight)

	_, err = NewStrValWeights([]string{"a"}, []int{10, 20})
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidValWeights, err)

	// Test Clone
	svw2 := svw1.Clone()
	assert.Equal(t, svw1.MaxWeight, svw2.MaxWeight)
	assert.Equal(t, len(svw1.Vals), len(svw2.Vals))
	assert.Equal(t, len(svw1.Weights), len(svw2.Weights))
	for i := range svw1.Vals {
		assert.Equal(t, svw1.Vals[i], svw2.Vals[i])
		assert.Equal(t, svw1.Weights[i], svw2.Weights[i])
	}

	// Test RandVal
	plugin := sgc7plugin.NewMockPlugin()
	plugin.Cache = []int{0}
	val, err := svw1.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, "a", val)

	plugin.Cache = []int{10}
	val, err = svw1.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, "b", val)

	svw3, _ := NewStrValWeights([]string{"c"}, []int{10})
	val, err = svw3.RandVal(plugin)
	assert.NoError(t, err)
	assert.Equal(t, "c", val)

	svwInvalid, _ := NewStrValWeights([]string{}, []int{})
	svwInvalid.MaxWeight = 0
	_, err = svwInvalid.RandVal(plugin)
	assert.Error(t, err)

	// Test RandIndex
	plugin.Cache = []int{0}
	idx, err := svw1.RandIndex(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 0, idx)

	plugin.Cache = []int{10}
	idx, err = svw1.RandIndex(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 1, idx)

	idx, err = svw3.RandIndex(plugin)
	assert.NoError(t, err)
	assert.Equal(t, 0, idx)

	_, err = svwInvalid.RandIndex(plugin)
	assert.Error(t, err)

	// Test CloneExcludeVal
	svw4, err := svw1.CloneExcludeVal("a")
	assert.NoError(t, err)
	assert.Equal(t, 20, svw4.MaxWeight)
	assert.Equal(t, 1, len(svw4.Vals))
	assert.Equal(t, "b", svw4.Vals[0])

	_, err = svw1.CloneExcludeVal("c")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidValWeightsVal, err)

	_, err = svw3.CloneExcludeVal("c")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidValWeights, err)

	// Test LoadStrValWeightsFromExcel
	svw5, err := LoadStrValWeightsFromExcel("../unittestdata/strvalweights.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, svw5)
	assert.Equal(t, 3, len(svw5.Vals))
	assert.Equal(t, 60, svw5.MaxWeight)


	_, err = LoadStrValWeightsFromExcel("filethatdoesnotexist.xlsx")
	assert.Error(t, err)

	_, err = LoadStrValWeightsFromExcel("../unittestdata/empty.json")
	assert.Error(t, err)
}
