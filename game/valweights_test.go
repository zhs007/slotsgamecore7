package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadValWeightsFromExcel(t *testing.T) {
	vw, err := LoadValWeightsFromExcel("../unittestdata/mysteryweights.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, vw)

	assert.Equal(t, vw.MaxWeight, 7600)
	assert.Equal(t, len(vw.Vals), 9)
	assert.Equal(t, len(vw.Weights), 9)

	nvw, err := vw.CloneExcludeVal(8)
	assert.NoError(t, err)

	assert.Equal(t, nvw.MaxWeight, 6600)
	assert.Equal(t, nvw.Vals, []int{1, 2, 3, 4, 5, 6, 7, 9})
	assert.Equal(t, nvw.Weights, []int{500, 600, 700, 800, 1000, 1000, 1000, 1000})

	t.Logf("Test_LoadValWeightsFromExcel OK")
}
