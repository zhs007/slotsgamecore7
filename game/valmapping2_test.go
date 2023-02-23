package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadValMapping2FromExcel(t *testing.T) {
	vam, err := LoadValMapping2FromExcel("../unittestdata/s5b5.xlsx", "index", "values", NewIntArrVal[int])
	assert.NoError(t, err)
	assert.NotNil(t, vam)

	assert.Equal(t, len(vam.MapVals), 15)

	arr1 := vam.MapVals[1].IntArr()
	assert.Equal(t, len(arr1), 5)
	assert.Equal(t, arr1[0], 1)
	assert.Equal(t, arr1[1], 1)
	assert.Equal(t, arr1[2], 3)
	assert.Equal(t, arr1[3], 3)
	assert.Equal(t, arr1[4], 2)

	arr10 := vam.MapVals[10].IntArr()
	assert.Equal(t, len(arr10), 5)
	assert.Equal(t, arr10[0], 2)
	assert.Equal(t, arr10[1], 1)
	assert.Equal(t, arr10[2], 3)
	assert.Equal(t, arr10[3], 1)
	assert.Equal(t, arr10[4], 2)

	arr15 := vam.MapVals[15].IntArr()
	assert.Equal(t, len(arr15), 5)
	assert.Equal(t, arr15[0], 2)
	assert.Equal(t, arr15[1], 1)
	assert.Equal(t, arr15[2], 4)
	assert.Equal(t, arr15[3], 1)
	assert.Equal(t, arr15[4], 3)

	t.Logf("Test_LoadValMapping2FromExcel OK")
}
