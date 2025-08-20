package sgc7game

import (
	"reflect"
	"sort"
	"testing"
)

func Test_ValArrMapping_Keys_jules(t *testing.T) {
	vm := &ValArrMapping[int, int]{
		MapVals: map[int][]int{
			1: {10, 20},
			2: {30, 40},
			3: {50, 60},
		},
	}
	keys := vm.Keys()
	sort.Ints(keys)
	if !reflect.DeepEqual(keys, []int{1, 2, 3}) {
		t.Errorf("Test_ValArrMapping_Keys_jules() Keys = %v, want %v", keys, []int{1, 2, 3})
	}
}

func Test_ValArrMapping_Clone_jules(t *testing.T) {
	vm := &ValArrMapping[int, int]{
		MapVals: map[int][]int{
			1: {10, 20},
			2: {30, 40},
		},
	}
	clone := vm.Clone()
	if !reflect.DeepEqual(vm, clone) {
		t.Errorf("Test_ValArrMapping_Clone_jules() Clone = %v, want %v", clone, vm)
	}
	clone.MapVals[1][0] = 100
	if vm.MapVals[1][0] == 100 {
		t.Error("Test_ValArrMapping_Clone_jules() failed, it's a shallow copy")
	}
}

func Test_NewValArrMapping_jules(t *testing.T) {
	typevals := []int{1, 2}
	vals := [][]int{{10, 20}, {30, 40}}
	vm, err := NewValArrMapping(typevals, vals)
	if err != nil {
		t.Fatalf("Test_NewValArrMapping_jules() NewValArrMapping failed: %v", err)
	}
	if len(vm.MapVals) != 2 {
		t.Errorf("Test_NewValArrMapping_jules() len(MapVals) = %v, want %v", len(vm.MapVals), 2)
	}

	_, err = NewValArrMapping([]int{1}, [][]int{{1}, {2}})
	if err == nil {
		t.Error("Test_NewValArrMapping_jules() expected an error for mismatched lengths")
	}
}

func Test_NewValArrMappingEx_jules(t *testing.T) {
	vm := NewValArrMappingEx[int, int]()
	if vm.MapVals == nil {
		t.Error("Test_NewValArrMappingEx_jules() failed, MapVals is nil")
	}
	if len(vm.MapVals) != 0 {
		t.Errorf("Test_NewValArrMappingEx_jules() len(MapVals) = %v, want %v", len(vm.MapVals), 0)
	}
}

func Test_LoadValArrMappingFromExcel_jules(t *testing.T) {
	vm, err := LoadValArrMappingFromExcel[int, int]("./testdata/valarrmapping.xlsx", "type", "val")
	if err != nil {
		t.Fatalf("Test_LoadValArrMappingFromExcel_jules() failed: %v", err)
	}
	if len(vm.MapVals) != 3 {
		t.Errorf("Test_LoadValArrMappingFromExcel_jules() len(MapVals) = %v, want %v", len(vm.MapVals), 3)
	}
	if !reflect.DeepEqual(vm.MapVals[1], []int{100, 101, 102}) {
		t.Errorf("Test_LoadValArrMappingFromExcel_jules() MapVals[1] = %v, want %v", vm.MapVals[1], []int{100, 101, 102})
	}
	if !reflect.DeepEqual(vm.MapVals[2], []int{200, 201}) {
		t.Errorf("Test_LoadValArrMappingFromExcel_jules() MapVals[2] = %v, want %v", vm.MapVals[2], []int{200, 201})
	}
	if !reflect.DeepEqual(vm.MapVals[3], []int{300}) {
		t.Errorf("Test_LoadValArrMappingFromExcel_jules() MapVals[3] = %v, want %v", vm.MapVals[3], []int{300})
	}

	// Test with a non-existent file
	_, err = LoadValArrMappingFromExcel[int, int]("./testdata/nonexistent.xlsx", "type", "val")
	if err == nil {
		t.Error("Test_LoadValArrMappingFromExcel_jules() expected an error for non-existent file")
	}
}
