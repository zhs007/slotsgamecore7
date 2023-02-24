package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadValWeights2FromExcel(t *testing.T) {
	vw, err := LoadValWeights2FromExcel("../unittestdata/mysteryweights.xlsx", "val", "weight", NewIntVal[int])
	assert.NoError(t, err)
	assert.NotNil(t, vw)

	assert.Equal(t, vw.MaxWeight, 7600)
	assert.Equal(t, len(vw.Vals), 9)
	assert.Equal(t, len(vw.Weights), 9)

	nvw, err := vw.CloneExcludeVal(&IntVal[int]{Val: 8})
	assert.NoError(t, err)

	assert.Equal(t, nvw.MaxWeight, 6600)
	assert.Equal(t, nvw.Vals[0].IsSame(&IntVal[int]{Val: 1}), true)
	assert.Equal(t, nvw.Vals[1].IsSame(&IntVal[int]{Val: 2}), true)
	assert.Equal(t, nvw.Vals[2].IsSame(&IntVal[int]{Val: 3}), true)
	assert.Equal(t, nvw.Vals[3].IsSame(&IntVal[int]{Val: 4}), true)
	assert.Equal(t, nvw.Vals[4].IsSame(&IntVal[int]{Val: 5}), true)
	assert.Equal(t, nvw.Vals[5].IsSame(&IntVal[int]{Val: 6}), true)
	assert.Equal(t, nvw.Vals[6].IsSame(&IntVal[int]{Val: 7}), true)
	assert.Equal(t, nvw.Vals[7].IsSame(&IntVal[int]{Val: 9}), true)
	assert.Equal(t, nvw.Weights, []int{500, 600, 700, 800, 1000, 1000, 1000, 1000})

	t.Logf("Test_LoadValWeights2FromExcel OK")
}

func Test_LoadValWeights2FromExcelWithSymbols(t *testing.T) {
	pt, err := LoadPaytablesFromExcel("../data/game001/paytables.xlsx")
	assert.NoError(t, err)

	vw, err := LoadValWeights2FromExcelWithSymbols("../data/game001/fgmysteryweights0.xlsx", "val", "weight", pt)
	assert.NoError(t, err)
	assert.NotNil(t, vw)

	assert.Equal(t, vw.MaxWeight, 10000)
	assert.Equal(t, len(vw.Vals), 10)
	assert.Equal(t, len(vw.Weights), 10)
	assert.Equal(t, vw.Vals[0].IsSame(&IntVal[int]{Val: 2}), true)
	assert.Equal(t, vw.Weights[0], 1065)

	t.Logf("Test_LoadValWeights2FromExcelWithSymbols OK")
}
