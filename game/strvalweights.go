package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// StrValWeights
type StrValWeights struct {
	Vals      []string
	Weights   []int
	MaxWeight int
}

func NewStrValWeights(vals []string, weights []int) (*StrValWeights, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewStrValWeights",
			slog.Int("vals", len(vals)),
			slog.Int("weights", len(weights)),
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	vw := &StrValWeights{
		Vals:      make([]string, len(vals)),
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

func (vw *StrValWeights) Clone() *StrValWeights {
	nvw := &StrValWeights{
		Vals:      make([]string, len(vw.Vals)),
		Weights:   make([]int, len(vw.Weights)),
		MaxWeight: vw.MaxWeight,
	}

	copy(nvw.Vals, vw.Vals)
	copy(nvw.Weights, vw.Weights)

	return nvw
}

func (vw *StrValWeights) RandVal(plugin sgc7plugin.IPlugin) (string, error) {
	if len(vw.Vals) == 1 {
		return vw.Vals[0], nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("StrValWeights.RandVal:RandWithWeights",
			goutils.Err(err))

		return "", err
	}

	return vw.Vals[ci], nil
}

func (vw *StrValWeights) RandIndex(plugin sgc7plugin.IPlugin) (int, error) {
	if len(vw.Vals) == 1 {
		return 0, nil
	}

	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("StrValWeights.RandVal:RandWithWeights",
			goutils.Err(err))

		return -1, err
	}

	return ci, nil
}

// CloneExcludeVal - clone & exclude a val
func (vw *StrValWeights) CloneExcludeVal(val string) (*StrValWeights, error) {
	if len(vw.Vals) <= 1 {
		goutils.Error("StrValWeights.RandVal:CloneExcludeVal",
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

	goutils.Error("StrValWeights.RandVal:CloneExcludeVal",
		goutils.Err(ErrInvalidValWeightsVal))

	return nil, ErrInvalidValWeightsVal
}

// LoadStrValWeightsFromExcel - load xlsx file
func LoadStrValWeightsFromExcel(fn string) (*StrValWeights, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadStrValWeightsFromExcel:OpenFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}
	defer f.Close()

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadStrValWeightsFromExcel:GetSheetList",
			slog.Any("SheetList", lstname),
			slog.String("fn", fn),
			goutils.Err(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadStrValWeightsFromExcel:GetRows",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	mapcolname := make(map[int]string)

	vals := []string{}
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
						vals = append(vals, colCell)
					} else if colname == "weight" {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadStrValWeightsFromExcel:String2Int64",
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

	return NewStrValWeights(vals, weights)
}
