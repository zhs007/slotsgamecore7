package sgc7game

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValMapping2_IsEmpty_jules(t *testing.T) {
	vm1, err := NewValMapping2([]int{1}, []IVal{NewIntValEx[int](10)})
	assert.NoError(t, err)
	assert.NotNil(t, vm1)
	assert.False(t, vm1.IsEmpty())

	vm2 := NewValMappingEx2()
	assert.NotNil(t, vm2)
	assert.True(t, vm2.IsEmpty())

	t.Logf("Test_ValMapping2_IsEmpty_jules OK")
}

func Test_ValMapping2_Keys_jules(t *testing.T) {
	vm1, err := NewValMapping2([]int{1, 2, 3}, []IVal{NewIntValEx[int](10), NewIntValEx[int](20), NewIntValEx[int](30)})
	assert.NoError(t, err)
	assert.NotNil(t, vm1)

	keys := vm1.Keys()
	sort.Ints(keys)
	assert.Equal(t, []int{1, 2, 3}, keys)

	vm2 := NewValMappingEx2()
	assert.NotNil(t, vm2)
	assert.Empty(t, vm2.Keys())

	t.Logf("Test_ValMapping2_Keys_jules OK")
}

func Test_ValMapping2_Clone_jules(t *testing.T) {
	vm1, err := NewValMapping2([]int{1}, []IVal{NewIntValEx[int](10)})
	assert.NoError(t, err)
	assert.NotNil(t, vm1)

	vm2 := vm1.Clone()
	assert.NotNil(t, vm2)
	assert.NotSame(t, vm1, vm2)
	assert.Equal(t, vm1.MapVals[1], vm2.MapVals[1])

	// As IVal is an interface, the values are pointers, so they are the same.
	// This is the expected behavior of Clone.
	// If a deep copy of IVal is needed, it should be done outside.
	assert.Same(t, vm1.MapVals[1], vm2.MapVals[1])

	vm1.MapVals[1] = NewIntValEx[int](20)
	assert.NotEqual(t, vm1.MapVals[1], vm2.MapVals[1])

	t.Logf("Test_ValMapping2_Clone_jules OK")
}

func Test_NewValMapping2_jules(t *testing.T) {
	vm1, err := NewValMapping2([]int{1, 2}, []IVal{NewIntValEx[int](10), NewIntValEx[int](20)})
	assert.NoError(t, err)
	assert.NotNil(t, vm1)
	assert.Equal(t, 2, len(vm1.MapVals))
	assert.Equal(t, int64(10), vm1.MapVals[1].Int64())
	assert.Equal(t, int64(20), vm1.MapVals[2].Int64())

	vm2, err := NewValMapping2([]int{1}, []IVal{NewIntValEx[int](10), NewIntValEx[int](20)})
	assert.Error(t, err)
	assert.Nil(t, vm2)
	assert.Equal(t, ErrInvalidValMapping, err)

	t.Logf("Test_NewValMapping2_jules OK")
}

func Test_NewValMappingEx2_jules(t *testing.T) {
	vm := NewValMappingEx2()
	assert.NotNil(t, vm)
	assert.NotNil(t, vm.MapVals)
	assert.Empty(t, vm.MapVals)

	t.Logf("Test_NewValMappingEx2_jules OK")
}

func Test_LoadValMapping2FromExcel_jules_nonexistent(t *testing.T) {
	vam, err := LoadValMapping2FromExcel("./testdata/nonexistent.xlsx", "index", "values", NewIntArrVal[int])
	assert.Error(t, err)
	assert.Nil(t, vam)

	t.Logf("Test_LoadValMapping2FromExcel_jules_nonexistent OK")
}

func Test_LoadValMapping2FromExcel_jules_invalid(t *testing.T) {
	vam, err := LoadValMapping2FromExcel("./testdata/valmapping_invalid.xlsx", "type", "val", NewIntVal[int])
	assert.Error(t, err)
	assert.Nil(t, vam)

	t.Logf("Test_LoadValMapping2FromExcel_jules_invalid OK")
}
