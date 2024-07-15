package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// ValWeights2
type ValWeights2 struct {
	Vals      []IVal `json:"vals"`
	Weights   []int  `json:"weights"`
	MaxWeight int    `json:"maxWeight"`
}

func (vw *ValWeights2) SortBy(dst *ValWeights2) error {
	if len(vw.Vals) != len(dst.Vals) {
		goutils.Error("ValWeights.SortBy",
			goutils.Err(ErrInvalidValWeights))

		return ErrInvalidValWeights
	}

	vals := make([]IVal, len(vw.Vals))
	weights := make([]int, len(vw.Weights))

	for i, v := range dst.Vals {
		vals[i] = v
		weights[i] = vw.GetWeight(v)
	}

	maxweights := 0
	for _, v := range weights {
		maxweights += v
	}

	if maxweights != vw.MaxWeight {
		goutils.Error("ValWeights.SortBy",
			slog.Int("MaxWeight", vw.MaxWeight),
			slog.Int("NewMaxWeight", maxweights),
			goutils.Err(ErrInvalidValWeights))

		return ErrInvalidValWeights
	}

	vw.Vals = vals
	vw.Weights = weights

	return nil
}

func (vw *ValWeights2) GetWeight(val IVal) int {
	for i, v := range vw.Vals {
		if val.IsSame(v) {
			return vw.Weights[i]
		}
	}

	return 0
}

func (vw *ValWeights2) Add(val IVal, weight int) {
	vw.Vals = append(vw.Vals, val)
	vw.Weights = append(vw.Weights, weight)

	vw.MaxWeight += weight
}

func (vw *ValWeights2) ClearExcludeVal(val IVal) {
	vw.Vals = []IVal{val}
	vw.Weights = []int{1}
	vw.MaxWeight = 1
}

func (vw *ValWeights2) Reset(vals []IVal, weights []int) {
	vw.Vals = make([]IVal, len(vals))
	vw.Weights = make([]int, len(weights))

	copy(vw.Vals, vals)
	copy(vw.Weights, weights)

	vw.MaxWeight = 0

	for _, v := range vw.Weights {
		vw.MaxWeight += v
	}
}

func (vw *ValWeights2) CloneWithoutIntArray(arr []int) *ValWeights2 {
	if len(arr) == 0 {
		return vw.Clone()
	}

	nvw := &ValWeights2{
		Vals:      make([]IVal, len(vw.Vals)),
		Weights:   make([]int, len(vw.Weights)),
		MaxWeight: 0,
	}

	for i, v := range vw.Vals {
		if goutils.IndexOfIntSlice(arr, v.Int(), 0) < 0 {
			nvw.Vals = append(nvw.Vals, v)
			nvw.Weights = append(nvw.Weights, vw.Weights[i])
			nvw.MaxWeight += vw.Weights[i]
		}
	}

	if len(nvw.Vals) == 0 {
		return nil
	}

	return nvw
}

func (vw *ValWeights2) CloneWithIntArray(arr []int) *ValWeights2 {
	if len(arr) == 0 {
		return nil
	}

	nvw := &ValWeights2{
		Vals:      make([]IVal, len(vw.Vals)),
		Weights:   make([]int, len(vw.Weights)),
		MaxWeight: 0,
	}

	for i, v := range vw.Vals {
		if goutils.IndexOfIntSlice(arr, v.Int(), 0) >= 0 {
			nvw.Vals = append(nvw.Vals, v)
			nvw.Weights = append(nvw.Weights, vw.Weights[i])
			nvw.MaxWeight += vw.Weights[i]
		}
	}

	if len(nvw.Vals) == 0 {
		return nil
	}

	return nvw
}

func (vw *ValWeights2) Clone() *ValWeights2 {
	nvw := &ValWeights2{
		Vals:      make([]IVal, len(vw.Vals)),
		Weights:   make([]int, len(vw.Weights)),
		MaxWeight: vw.MaxWeight,
	}

	copy(nvw.Vals, vw.Vals)
	copy(nvw.Weights, vw.Weights)

	return nvw
}

func (vw *ValWeights2) getValidWeightNum() int {
	num := 0

	for _, v := range vw.Weights {
		if v > 0 {
			num++
		}
	}

	return num
}

func (vw *ValWeights2) Normalize() {
	vnum := vw.getValidWeightNum()
	if vnum != len(vw.Weights) {
		vals := make([]IVal, vnum)
		weights := make([]int, vnum)

		vw.MaxWeight = 0

		for i, v := range vw.Weights {
			if v > 0 {
				vals = append(vals, vw.Vals[i])
				weights = append(weights, v)

				vw.MaxWeight += v
			}
		}

		vw.Vals = vals
		vw.Weights = weights
	}
}

func (vw *ValWeights2) RandVal(plugin sgc7plugin.IPlugin) (IVal, error) {
	if len(vw.Vals) == 1 {
		return vw.Vals[0], nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights2.RandVal:RandWithWeights",
			goutils.Err(err))

		return nil, err
	}

	return vw.Vals[ci], nil
}

func (vw *ValWeights2) RandValEx(plugin sgc7plugin.IPlugin) (IVal, int, error) {
	if len(vw.Vals) == 1 {
		return vw.Vals[0], 0, nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights2.RandVal:RandWithWeights",
			goutils.Err(err))

		return nil, 0, err
	}

	return vw.Vals[ci], ci, nil
}

// CloneExcludeVal - clone & exclude a val
func (vw *ValWeights2) CloneExcludeVal(val IVal) (*ValWeights2, error) {
	if len(vw.Vals) <= 1 {
		goutils.Error("ValWeights2.RandVal:CloneExcludeVal",
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	nvw := vw.Clone()

	for i, v := range vw.Vals {
		if val.IsSame(v) {
			nvw.Vals = append(nvw.Vals[0:i], nvw.Vals[i+1:]...)
			nvw.Weights = append(nvw.Weights[0:i], nvw.Weights[i+1:]...)
			nvw.MaxWeight -= vw.Weights[i]

			return nvw, nil
		}
	}

	goutils.Error("ValWeights.RandVal:CloneExcludeVal",
		goutils.Err(ErrInvalidValWeightsVal))

	return nil, ErrInvalidValWeightsVal
}

// RemoveVal - remove a val
func (vw *ValWeights2) RemoveVal(val IVal) error {
	for i, v := range vw.Vals {
		if val.IsSame(v) {
			weight := vw.Weights[i]

			vw.Vals = append(vw.Vals[0:i], vw.Vals[i+1:]...)
			vw.Weights = append(vw.Weights[0:i], vw.Weights[i+1:]...)

			vw.MaxWeight -= weight

			return nil
		}
	}

	goutils.Error("ValWeights.RandVal:RemoveVal",
		goutils.Err(ErrInvalidValWeightsVal))

	return ErrInvalidValWeightsVal
}

func NewValWeights2Ex() *ValWeights2 {
	return &ValWeights2{}
}

func NewValWeights2(vals []IVal, weights []int) (*ValWeights2, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewValWeights",
			slog.Int("vals", len(vals)),
			slog.Int("weights", len(weights)),
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	vw := &ValWeights2{
		Vals:      make([]IVal, len(vals)),
		Weights:   make([]int, len(vals)),
		MaxWeight: 0,
	}

	copy(vw.Vals, vals)
	copy(vw.Weights, weights)

	for _, v := range weights {
		vw.MaxWeight += v
	}

	return vw, nil
}

// LoadValWeights2FromExcel - load xlsx file
func LoadValWeights2FromExcel(fn string, headerVal string, headerWeight string, funcNew FuncNewIVal) (*ValWeights2, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadValWeightsFromExcel:OpenFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}
	defer f.Close()

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadValWeightsFromExcel:GetSheetList",
			slog.Any("SheetList", lstname),
			slog.String("fn", fn),
			goutils.Err(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadValWeightsFromExcel:GetRows",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	mapcolname := make(map[int]string)

	vals := []IVal{}
	weights := []int{}

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = strings.ToLower(colCell)
			}
		} else {
			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					if colname == headerVal {
						cv := funcNew()
						err := cv.ParseString(colCell)
						if err != nil {
							goutils.Error("LoadValWeightsFromExcel:ParseString",
								slog.String("val", colCell),
								goutils.Err(err))

							return nil, err
						}

						vals = append(vals, cv)
					} else if colname == headerWeight {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadValWeightsFromExcel:String2Int64",
								slog.String("weight", colCell),
								goutils.Err(err))

							return nil, err
						}

						weights = append(weights, int(v))
					}
				}
			}
		}
	}

	return NewValWeights2(vals, weights)
}

// LoadValWeights2FromExcelWithSymbols - load xlsx file
func LoadValWeights2FromExcelWithSymbols(fn string, headerVal string, headerWeight string, paytables *PayTables) (*ValWeights2, error) {
	vw, err := LoadValWeights2FromExcel(fn, headerVal, headerWeight, NewStrVal)
	if err != nil {
		goutils.Error("LoadValWeights2FromExcelWithSymbols:LoadValWeights2FromExcel",
			goutils.Err(err))

		return nil, err
	}

	vals := make([]IVal, len(vw.Vals))

	for i, v := range vw.Vals {
		vals[i] = NewIntValEx(paytables.MapSymbols[v.String()])
	}

	nvw, err := NewValWeights2(vals, vw.Weights)
	if err != nil {
		goutils.Error("LoadValWeights2FromExcelWithSymbols:NewValWeights2",
			goutils.Err(err))

		return nil, err
	}

	return nvw, nil
}
