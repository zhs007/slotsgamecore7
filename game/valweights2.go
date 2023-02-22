package sgc7game

import (
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
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
			zap.Error(ErrInvalidValWeights))

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
			zap.Int("MaxWeight", vw.MaxWeight),
			zap.Int("NewMaxWeight", maxweights),
			zap.Error(ErrInvalidValWeights))

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

func (vw *ValWeights2) RandVal(plugin sgc7plugin.IPlugin) (IVal, error) {
	if len(vw.Vals) == 1 {
		return vw.Vals[0], nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights2.RandVal:RandWithWeights",
			zap.Error(err))

		return nil, err
	}

	return vw.Vals[ci], nil
}

// CloneExcludeVal - clone & exclude a val
func (vw *ValWeights2) CloneExcludeVal(val IVal) (*ValWeights2, error) {
	if len(vw.Vals) <= 1 {
		goutils.Error("ValWeights2.RandVal:CloneExcludeVal",
			zap.Error(ErrInvalidValWeights))

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
		zap.Error(ErrInvalidValWeightsVal))

	return nil, ErrInvalidValWeightsVal
}

func NewValWeights2Ex() *ValWeights2 {
	return &ValWeights2{}
}

func NewValWeights2(vals []IVal, weights []int) (*ValWeights2, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewValWeights",
			zap.Int("vals", len(vals)),
			zap.Int("weights", len(weights)),
			zap.Error(ErrInvalidValWeights))

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
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer f.Close()

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadValWeightsFromExcel:GetSheetList",
			goutils.JSON("SheetList", lstname),
			zap.String("fn", fn),
			zap.Error(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadValWeightsFromExcel:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

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
								zap.String("val", colCell),
								zap.Error(err))

							return nil, err
						}

						vals = append(vals, cv)
					} else if colname == headerWeight {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadValWeightsFromExcel:String2Int64",
								zap.String("weight", colCell),
								zap.Error(err))

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
