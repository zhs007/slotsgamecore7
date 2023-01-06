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

	t.Logf("Test_LoadValWeightsFromExcel OK")
}
