package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewFloatValMapping(t *testing.T) {
	typevals := []int{1, 2}
	vals := []float32{1.1, 2.2}
	fvm, err := NewFloatValMapping(typevals, vals)
	assert.NoError(t, err)
	assert.NotNil(t, fvm)
	assert.Equal(t, float32(1.1), fvm.MapVals[1])
	assert.Equal(t, float32(2.2), fvm.MapVals[2])

	t.Logf("Test_NewFloatValMapping OK")
}

func Test_FloatValMapping_Clone(t *testing.T) {
	typevals := []int{1, 2}
	vals := []float32{1.1, 2.2}
	fvm, err := NewFloatValMapping(typevals, vals)
	assert.NoError(t, err)
	assert.NotNil(t, fvm)

	cloned := fvm.Clone()
	assert.NotNil(t, cloned)
	assert.Equal(t, fvm, cloned)

	t.Logf("Test_FloatValMapping_Clone OK")
}

func Test_LoadFloatValMappingFromExcel(t *testing.T) {
	fvm, err := LoadFloatValMappingFromExcel[int, float32]("../unittestdata/jules_test.xlsx", "type", "val")
	assert.NoError(t, err)
	assert.NotNil(t, fvm)
	assert.Equal(t, float32(1.1), fvm.MapVals[1])
	assert.Equal(t, float32(2.2), fvm.MapVals[2])

	t.Logf("Test_LoadFloatValMappingFromExcel OK")
}
