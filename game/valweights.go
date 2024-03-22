package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// ValWeights
type ValWeights struct {
	Vals      []int `json:"vals"`
	Weights   []int `json:"weights"`
	MaxWeight int   `json:"maxWeight"`
}

func (vw *ValWeights) SortBy(dst *ValWeights) error {
	if len(vw.Vals) != len(dst.Vals) {
		goutils.Error("ValWeights.SortBy",
			goutils.Err(ErrInvalidValWeights))

		return ErrInvalidValWeights
	}

	vals := make([]int, len(vw.Vals))
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

func (vw *ValWeights) GetWeight(val int) int {
	for i, v := range vw.Vals {
		if v == val {
			return vw.Weights[i]
		}
	}

	return 0
}

func (vw *ValWeights) Add(val int, weight int) {
	vw.Vals = append(vw.Vals, val)
	vw.Weights = append(vw.Weights, weight)

	vw.MaxWeight += weight
}

func (vw *ValWeights) ClearExcludeVal(val int) {
	vw.Vals = []int{val}
	vw.Weights = []int{1}
	vw.MaxWeight = 1
}

func (vw *ValWeights) Reset(vals []int, weights []int) {
	vw.Vals = make([]int, len(vals))
	vw.Weights = make([]int, len(weights))

	copy(vw.Vals, vals)
	copy(vw.Weights, weights)

	vw.MaxWeight = 0

	for _, v := range vw.Weights {
		vw.MaxWeight += v
	}
}

func (vw *ValWeights) Clone() *ValWeights {
	nvw := &ValWeights{
		Vals:      make([]int, len(vw.Vals)),
		Weights:   make([]int, len(vw.Weights)),
		MaxWeight: vw.MaxWeight,
	}

	copy(nvw.Vals, vw.Vals)
	copy(nvw.Weights, vw.Weights)

	return nvw
}

func (vw *ValWeights) RandVal(plugin sgc7plugin.IPlugin) (int, error) {
	if len(vw.Vals) == 1 {
		return vw.Vals[0], nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights.RandVal:RandWithWeights",
			goutils.Err(err))

		return 0, err
	}

	return vw.Vals[ci], nil
}

// CloneExcludeVal - clone & exclude a val
func (vw *ValWeights) CloneExcludeVal(val int) (*ValWeights, error) {
	if len(vw.Vals) <= 1 {
		goutils.Error("ValWeights.RandVal:CloneExcludeVal",
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	nvw := vw.Clone()

	for i, v := range vw.Vals {
		if v == val {
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

func NewValWeightsEx() *ValWeights {
	return &ValWeights{}
}

func NewValWeights(vals []int, weights []int) (*ValWeights, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewValWeights",
			slog.Int("vals", len(vals)),
			slog.Int("weights", len(weights)),
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	vw := &ValWeights{
		Vals:      make([]int, len(vals)),
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

// LoadValWeightsFromExcel - load xlsx file
func LoadValWeightsFromExcel(fn string) (*ValWeights, error) {
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

	vals := []int{}
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
					if colname == "val" {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadValWeightsFromExcel:String2Int64",
								slog.String("val", colCell),
								goutils.Err(err))

							return nil, err
						}

						vals = append(vals, int(v))
					} else if colname == "weight" {
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

	return NewValWeights(vals, weights)
}
