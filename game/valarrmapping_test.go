package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadValArrMappingFromExcel(t *testing.T) {
	vam, err := LoadValArrMappingFromExcel[int, int]("../unittestdata/s5b5.xlsx", "index", "values")
	assert.NoError(t, err)
	assert.NotNil(t, vam)

	assert.Equal(t, len(vam.MapVals), 15)

	assert.Equal(t, len(vam.MapVals[1]), 5)
	assert.Equal(t, vam.MapVals[1][0], 1)
	assert.Equal(t, vam.MapVals[1][1], 1)
	assert.Equal(t, vam.MapVals[1][2], 3)
	assert.Equal(t, vam.MapVals[1][3], 3)
	assert.Equal(t, vam.MapVals[1][4], 2)

	assert.Equal(t, len(vam.MapVals[10]), 5)
	assert.Equal(t, vam.MapVals[10][0], 2)
	assert.Equal(t, vam.MapVals[10][1], 1)
	assert.Equal(t, vam.MapVals[10][2], 3)
	assert.Equal(t, vam.MapVals[10][3], 1)
	assert.Equal(t, vam.MapVals[10][4], 2)

	assert.Equal(t, len(vam.MapVals[15]), 5)
	assert.Equal(t, vam.MapVals[15][0], 2)
	assert.Equal(t, vam.MapVals[15][1], 1)
	assert.Equal(t, vam.MapVals[15][2], 4)
	assert.Equal(t, vam.MapVals[15][3], 1)
	assert.Equal(t, vam.MapVals[15][4], 3)

	t.Logf("Test_LoadValArrMappingFromExcel OK")
}
