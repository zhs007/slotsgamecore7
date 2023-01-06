package sgc7game

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

// ValWeights
type ValWeights struct {
	Vals      []int
	Weights   []int
	MaxWeight int
}

func NewValWeights(vals []int, weights []int) (*ValWeights, error) {
	if len(vals) != len(weights) {
		goutils.Error("NewValWeights",
			zap.Int("vals", len(vals)),
			zap.Int("weights", len(weights)),
			zap.Error(ErrInvalidValWeights))

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

func (vw *ValWeights) RandVal(plugin sgc7plugin.IPlugin) (int, error) {
	ci, err := RandWithWeights(plugin, vw.MaxWeight, vw.Weights)
	if err != nil {
		goutils.Error("ValWeights.RandVal:RandWithWeights",
			zap.Error(err))

		return 0, err
	}

	return vw.Vals[ci], nil
}

// LoadValWeightsFromExcel - load xlsx file
func LoadValWeightsFromExcel(fn string) (*ValWeights, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadValWeightsFromExcel:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

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

	vals := []int{}
	weights := []int{}

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = colCell
			}
		} else {
			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					if colname == "val" {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadValWeightsFromExcel:String2Int64",
								zap.String("val", colCell),
								zap.Error(err))

							return nil, err
						}

						vals = append(vals, int(v))
					} else if colname == "weight" {
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

	return NewValWeights(vals, weights)
}
