package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadArrValWeightsFromExcel(t *testing.T) {
	vw, err := LoadArrValWeightsFromExcel("../unittestdata/arrvalweights.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, vw)

	assert.Equal(t, vw.MaxWeight, 9750)
	assert.Equal(t, len(vw.ArrVals), 15)
	assert.Equal(t, len(vw.Weights), 15)

	assert.Equal(t, vw.ArrVals[0][0], 2)
	assert.Equal(t, vw.ArrVals[0][1], 2)
	assert.Equal(t, vw.ArrVals[0][2], 2)
	assert.Equal(t, vw.ArrVals[0][3], 1)
	assert.Equal(t, vw.ArrVals[0][4], 1)

	t.Logf("Test_LoadArrValWeightsFromExcel OK")
}
