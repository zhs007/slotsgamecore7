package sgc7game

import (
	"testing"
	"sort"
	"github.com/xuri/excelize/v2"
	"github.com/stretchr/testify/assert"
)

func Test_ValMapping_Keys_Jules(t *testing.T) {
	vm := &ValMapping[int, int]{
		MapVals: map[int]int{1: 100, 2: 200, 3: 300},
	}
	keys := vm.Keys()
	sort.Ints(keys)
	assert.Equal(t, []int{1, 2, 3}, keys)
}

func Test_ValMapping_Clone_Jules(t *testing.T) {
	vm := &ValMapping[int, int]{
		MapVals: map[int]int{1: 100, 2: 200, 3: 300},
	}
	clone := vm.Clone()
	assert.Equal(t, vm.MapVals, clone.MapVals)

	clone.MapVals[1] = 101
	assert.NotEqual(t, vm.MapVals[1], clone.MapVals[1])
}

func Test_NewValMapping_Jules(t *testing.T) {
	// Test case 1: Valid input
	typevals1 := []int{1, 2}
	vals1 := []int{100, 200}
	vm1, err1 := NewValMapping(typevals1, vals1)
	assert.NoError(t, err1)
	assert.NotNil(t, vm1)
	assert.Equal(t, map[int]int{1: 100, 2: 200}, vm1.MapVals)

	// Test case 2: Mismatched lengths
	typevals2 := []int{1}
	vals2 := []int{100, 200}
	_, err2 := NewValMapping(typevals2, vals2)
	assert.Error(t, err2)
	assert.Equal(t, ErrInvalidValMapping, err2)
}

func Test_NewValMappingEx_Jules(t *testing.T) {
	vm := NewValMappingEx[int, int]()
	assert.NotNil(t, vm)
	assert.Empty(t, vm.MapVals)
}

func Test_LoadValMappingFromExcel_Jules(t *testing.T) {
	// Test case 1: Valid excel file
	vm, err := LoadValMappingFromExcel[int, int]("./testdata/valmapping.xlsx", "type", "val")
	assert.NoError(t, err)
	assert.NotNil(t, vm)
	assert.Equal(t, map[int]int{1: 100, 2: 200, 3: 300}, vm.MapVals)

	// Test case 2: File not found
	_, err = LoadValMappingFromExcel[int, int]("./testdata/nonexistent.xlsx", "type", "val")
	assert.Error(t, err)

	// Test case 3: Invalid header
	_, err = LoadValMappingFromExcel[int, int]("./testdata/valmapping.xlsx", "invalid", "val")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidValMapping, err)

	// Test case 4: Invalid data
	// create a temporary excel file with invalid data
	f := excelize.NewFile()
	index, _ := f.NewSheet("Sheet1")
	f.SetCellValue("Sheet1", "A1", "type")
	f.SetCellValue("Sheet1", "B1", "val")
	f.SetCellValue("Sheet1", "A2", "a")
	f.SetCellValue("Sheet1", "B2", "b")
	f.SetActiveSheet(index)
	err = f.SaveAs("./testdata/valmapping_invalid.xlsx")
	assert.NoError(t, err)

	_, err = LoadValMappingFromExcel[int, int]("./testdata/valmapping_invalid.xlsx", "type", "val")
	assert.Error(t, err)
}
