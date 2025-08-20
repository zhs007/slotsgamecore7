package sgc7game

import (
	"reflect"
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_ValWeights2_SortBy_jules(t *testing.T) {
	vw1, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw2, err := NewValWeights2([]IVal{NewIntValEx(2), NewIntValEx(1)}, []int{20, 10})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}

	err = vw1.SortBy(vw2)
	if err != nil {
		t.Fatalf("SortBy failed: %v", err)
	}

	if !reflect.DeepEqual(vw1.Vals, vw2.Vals) {
		t.Errorf("Test_ValWeights2_SortBy_jules() vw1.Vals = %v, want %v", vw1.Vals, vw2.Vals)
	}
	if !reflect.DeepEqual(vw1.Weights, vw2.Weights) {
		t.Errorf("Test_ValWeights2_SortBy_jules() vw1.Weights = %v, want %v", vw1.Weights, vw2.Weights)
	}

	vw3, err := NewValWeights2([]IVal{NewIntValEx(3)}, []int{30})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	err = vw1.SortBy(vw3)
	if err == nil {
		t.Error("Test_ValWeights2_SortBy_jules() expected an error for mismatched lengths")
	}

	vw4, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw5, err := NewValWeights2([]IVal{NewIntValEx(2), NewIntValEx(3)}, []int{10, 10})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	err = vw4.SortBy(vw5)
	if err == nil {
		t.Error("Test_ValWeights2_SortBy_jules() expected an error for invalid max weight")
	}
}

func Test_ValWeights2_GetWeight_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	if w := vw.GetWeight(NewIntValEx(1)); w != 10 {
		t.Errorf("Test_ValWeights2_GetWeight_jules() GetWeight(1) = %v, want %v", w, 10)
	}
	if w := vw.GetWeight(NewIntValEx(3)); w != 0 {
		t.Errorf("Test_ValWeights2_GetWeight_jules() GetWeight(3) = %v, want %v", w, 0)
	}
}

func Test_ValWeights2_Add_jules(t *testing.T) {
	vw := NewValWeights2Ex()
	vw.Add(NewIntValEx(1), 10)
	if len(vw.Vals) != 1 || vw.Vals[0].Int() != 1 || len(vw.Weights) != 1 || vw.Weights[0] != 10 || vw.MaxWeight != 10 {
		t.Errorf("Test_ValWeights2_Add_jules() failed after first add")
	}
	vw.Add(NewIntValEx(2), 20)
	if len(vw.Vals) != 2 || vw.Vals[1].Int() != 2 || len(vw.Weights) != 2 || vw.Weights[1] != 20 || vw.MaxWeight != 30 {
		t.Errorf("Test_ValWeights2_Add_jules() failed after second add")
	}
}

func Test_ValWeights2_ClearExcludeVal_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw.ClearExcludeVal(NewIntValEx(1))
	if len(vw.Vals) != 1 || vw.Vals[0].Int() != 1 || len(vw.Weights) != 1 || vw.Weights[0] != 1 || vw.MaxWeight != 1 {
		t.Errorf("Test_ValWeights2_ClearExcludeVal_jules() failed")
	}
}

func Test_ValWeights2_ResetMaxWeight_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw.MaxWeight = 100
	vw.ResetMaxWeight()
	if vw.MaxWeight != 30 {
		t.Errorf("Test_ValWeights2_ResetMaxWeight_jules() MaxWeight = %v, want %v", vw.MaxWeight, 30)
	}
}

func Test_ValWeights2_Reset_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw.Reset([]IVal{NewIntValEx(3)}, []int{30})
	if len(vw.Vals) != 1 || vw.Vals[0].Int() != 3 || len(vw.Weights) != 1 || vw.Weights[0] != 30 || vw.MaxWeight != 30 {
		t.Errorf("Test_ValWeights2_Reset_jules() failed")
	}
}

func Test_ValWeights2_CloneWithoutIntArray_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2), NewIntValEx(3)}, []int{10, 20, 30})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	nvw := vw.CloneWithoutIntArray([]int{2})
	if len(nvw.Vals) != 2 || nvw.MaxWeight != 40 || nvw.GetWeight(NewIntValEx(2)) != 0 {
		t.Errorf("Test_ValWeights2_CloneWithoutIntArray_jules() failed")
	}
	nvw2 := vw.CloneWithoutIntArray(nil)
	if !reflect.DeepEqual(vw, nvw2) {
		t.Errorf("Test_ValWeights2_CloneWithoutIntArray_jules() failed for nil array")
	}
	nvw3 := vw.CloneWithoutIntArray([]int{1, 2, 3})
	if nvw3 != nil {
		t.Errorf("Test_ValWeights2_CloneWithoutIntArray_jules() expected nil for all elements excluded")
	}
}

func Test_ValWeights2_CloneWithIntArray_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2), NewIntValEx(3)}, []int{10, 20, 30})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	nvw := vw.CloneWithIntArray([]int{2})
	if len(nvw.Vals) != 1 || nvw.MaxWeight != 20 || nvw.GetWeight(NewIntValEx(2)) != 20 {
		t.Errorf("Test_ValWeights2_CloneWithIntArray_jules() failed")
	}
	nvw2 := vw.CloneWithIntArray(nil)
	if nvw2 != nil {
		t.Errorf("Test_ValWeights2_CloneWithIntArray_jules() expected nil for nil array")
	}
	nvw3 := vw.CloneWithIntArray([]int{4, 5, 6})
	if nvw3 != nil {
		t.Errorf("Test_ValWeights2_CloneWithIntArray_jules() expected nil for no elements included")
	}
}

func Test_ValWeights2_Normalize_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2), NewIntValEx(3)}, []int{10, 0, 30})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vw.Normalize()
	if len(vw.Vals) != 2 || vw.MaxWeight != 40 || vw.GetWeight(NewIntValEx(2)) != 0 {
		t.Errorf("Test_ValWeights2_Normalize_jules() failed")
	}
	vw.Normalize() // should not change anything
	if len(vw.Vals) != 2 || vw.MaxWeight != 40 || vw.GetWeight(NewIntValEx(2)) != 0 {
		t.Errorf("Test_ValWeights2_Normalize_jules() failed on second normalize")
	}
}

func Test_ValWeights2_RandVal_jules(t *testing.T) {
	plugin := sgc7plugin.NewMockPlugin()
	plugin.Cache = []int{0, 15, 5, 25, 1, 18}

	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}

	v, err := vw.RandVal(plugin)
	if err != nil {
		t.Fatalf("RandVal failed: %v", err)
	}
	if v.Int() != 1 {
		t.Errorf("Test_ValWeights2_RandVal_jules() RandVal() = %v, want %v", v.Int(), 1)
	}

	v, i, err := vw.RandValEx(plugin)
	if err != nil {
		t.Fatalf("RandValEx failed: %v", err)
	}
	if v.Int() != 2 || i != 1 {
		t.Errorf("Test_ValWeights2_RandVal_jules() RandValEx() = %v, %v, want %v, %v", v.Int(), i, 2, 1)
	}

	vw1, _ := NewValWeights2([]IVal{NewIntValEx(1)}, []int{10})
	v, err = vw1.RandVal(plugin)
	if err != nil || v.Int() != 1 {
		t.Errorf("Test_ValWeights2_RandVal_jules() RandVal() with single value failed")
	}
	v, i, err = vw1.RandValEx(plugin)
	if err != nil || v.Int() != 1 || i != 0 {
		t.Errorf("Test_ValWeights2_RandVal_jules() RandValEx() with single value failed")
	}
}

func Test_LoadValWeights2FromExcel_jules(t *testing.T) {
	vw, err := LoadValWeights2FromExcel("./testdata/valweights2.xlsx", "val", "weight", NewStrVal)
	if err != nil {
		t.Fatalf("LoadValWeights2FromExcel failed: %v", err)
	}
	if len(vw.Vals) != 3 {
		t.Errorf("Test_LoadValWeights2FromExcel_jules() len(Vals) = %v, want %v", len(vw.Vals), 3)
	}
	if vw.MaxWeight != 60 {
		t.Errorf("Test_LoadValWeights2FromExcel_jules() MaxWeight = %v, want %v", vw.MaxWeight, 60)
	}

	// Test with a non-existent file
	_, err = LoadValWeights2FromExcel("./testdata/nonexistent.xlsx", "val", "weight", NewStrVal)
	if err == nil {
		t.Error("Test_LoadValWeights2FromExcel_jules() expected an error for non-existent file")
	}
}

func Test_LoadValWeights2FromExcelWithSymbols_jules(t *testing.T) {
	paytables := &PayTables{
		MapSymbols: map[string]int{
			"a": 1,
			"b": 2,
			"c": 3,
		},
	}
	vw, err := LoadValWeights2FromExcelWithSymbols("./testdata/valweights2.xlsx", "val", "weight", paytables)
	if err != nil {
		t.Fatalf("LoadValWeights2FromExcelWithSymbols failed: %v", err)
	}
	if len(vw.Vals) != 3 {
		t.Errorf("Test_LoadValWeights2FromExcelWithSymbols_jules() len(Vals) = %v, want %v", len(vw.Vals), 3)
	}
	if vw.MaxWeight != 60 {
		t.Errorf("Test_LoadValWeights2FromExcelWithSymbols_jules() MaxWeight = %v, want %v", vw.MaxWeight, 60)
	}
	if vw.Vals[0].Int() != 1 {
		t.Errorf("Test_LoadValWeights2FromExcelWithSymbols_jules() Vals[0] = %v, want %v", vw.Vals[0].Int(), 1)
	}
}

func Test_ValWeights2_CloneExcludeVal_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2), NewIntValEx(3)}, []int{10, 20, 30})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}

	nvw, err := vw.CloneExcludeVal(NewIntValEx(2))
	if err != nil {
		t.Fatalf("CloneExcludeVal failed: %v", err)
	}
	if len(nvw.Vals) != 2 || nvw.MaxWeight != 40 {
		t.Errorf("Test_ValWeights2_CloneExcludeVal_jules() wrong length or maxweight")
	}
	if nvw.GetWeight(NewIntValEx(2)) != 0 {
		t.Errorf("Test_ValWeights2_CloneExcludeVal_jules() val 2 should be excluded")
	}

	// exclude non-existent value
	_, err = vw.CloneExcludeVal(NewIntValEx(4))
	if err == nil {
		t.Error("Test_ValWeights2_CloneExcludeVal_jules() expected an error for excluding non-existent val")
	}

	// exclude from single-element vw
	vw1, _ := NewValWeights2([]IVal{NewIntValEx(1)}, []int{10})
	_, err = vw1.CloneExcludeVal(NewIntValEx(1))
	if err == nil {
		t.Error("Test_ValWeights2_CloneExcludeVal_jules() expected an error for excluding from single element vw")
	}
}

func Test_ValWeights2_RemoveVal_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}

	err = vw.RemoveVal(NewIntValEx(1))
	if err != nil {
		t.Fatalf("RemoveVal failed: %v", err)
	}
	if len(vw.Vals) != 1 || vw.MaxWeight != 20 || vw.GetWeight(NewIntValEx(1)) != 0 {
		t.Errorf("Test_ValWeights2_RemoveVal_jules() failed")
	}

	err = vw.RemoveVal(NewIntValEx(3))
	if err == nil {
		t.Error("Test_ValWeights2_RemoveVal_jules() expected an error for removing non-existent val")
	}
}

func Test_ValWeights2_GetIntVals_jules(t *testing.T) {
	vw, err := NewValWeights2([]IVal{NewIntValEx(1), NewIntValEx(2)}, []int{10, 20})
	if err != nil {
		t.Fatalf("NewValWeights2 failed: %v", err)
	}
	vals := vw.GetIntVals()
	if !reflect.DeepEqual(vals, []int{1, 2}) {
		t.Errorf("Test_ValWeights2_GetIntVals_jules() GetIntVals() = %v, want %v", vals, []int{1, 2})
	}
}
