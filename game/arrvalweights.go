package sgc7game

import (
	"log/slog"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// ArrValWeights
type ArrValWeights struct {
	ArrVals   [][]int
	Weights   []int
	MaxWeight int
}

func NewArrValWeights(arrvals [][]int, weights []int) (*ArrValWeights, error) {
	if len(arrvals) != len(weights) {
		goutils.Error("NewArrValWeights",
			slog.Int("vals", len(arrvals)),
			slog.Int("weights", len(weights)),
			goutils.Err(ErrInvalidValWeights))

		return nil, ErrInvalidValWeights
	}

	vw := &ArrValWeights{
		Weights:   make([]int, len(arrvals)),
		MaxWeight: 0,
	}

	for _, arr := range arrvals {
		carr := make([]int, len(arr))
		copy(carr, arr)

		vw.ArrVals = append(vw.ArrVals, carr)
	}

	copy(vw.Weights, weights)

	for _, v := range weights {
		vw.MaxWeight += v
	}

	return vw, nil
}

func (vw *ArrValWeights) RandVal(plugin sgc7plugin.IPlugin) ([]int, error) {
	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ArrValWeights.RandVal:RandWithWeights",
			goutils.Err(err))

		return nil, err
	}

	return vw.ArrVals[ci], nil
}

// LoadArrValWeightsFromExcel - load xlsx file
func LoadArrValWeightsFromExcel(fn string) (*ArrValWeights, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadArrValWeightsFromExcel:OpenFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadArrValWeightsFromExcel:GetSheetList",
			slog.Any("SheetList", lstname),
			slog.String("fn", fn),
			goutils.Err(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadArrValWeightsFromExcel:GetRows",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	mapcolname := make(map[int]string)

	arrvals := [][]int{}
	weights := []int{}

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = strings.ToLower(colCell)
			}
		} else {
			vals := []int{}

			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					if strings.Index(colname, "val") == 0 {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadArrValWeightsFromExcel:String2Int64",
								slog.String("val", colCell),
								goutils.Err(err))

							return nil, err
						}

						vals = append(vals, int(v))
					} else if colname == "weight" {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadArrValWeightsFromExcel:String2Int64",
								slog.String("weight", colCell),
								goutils.Err(err))

							return nil, err
						}

						weights = append(weights, int(v))
					}
				}
			}

			arrvals = append(arrvals, vals)
		}
	}

	return NewArrValWeights(arrvals, weights)
}
