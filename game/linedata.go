package sgc7game

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type lineInfo struct {
	R1   int `json:"R1"`
	R2   int `json:"R2"`
	R3   int `json:"R3"`
	R4   int `json:"R4"`
	R5   int `json:"R5"`
	R6   int `json:"R6"`
	Line int `json:"line"`
}

// LineData - line data
type LineData struct {
	Lines [][]int `json:"lines"`
}

// isValidLI5 - is it valid lineInfo5
func isValidLI5(li5s []lineInfo) bool {
	if len(li5s) <= 0 {
		return false
	}

	// alllinezero := true
	for _, v := range li5s {
		if v.Line > 0 {
			// alllinezero = false

			return true
		}
	}

	return false
}

// LoadLine5JSON - load json file
func LoadLine5JSON(fn string) (*LineData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if !isValidLI5(li) {
		return nil, nil
	}

	d := &LineData{}
	for _, v := range li {
		cl := []int{v.R1, v.R2, v.R3, v.R4, v.R5}
		d.Lines = append(d.Lines, cl)
	}

	return d, nil
}

// LoadLine3JSON - load json file
func LoadLine3JSON(fn string) (*LineData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if !isValidLI5(li) {
		return nil, nil
	}

	d := &LineData{}
	for _, v := range li {
		cl := []int{v.R1, v.R2, v.R3}
		d.Lines = append(d.Lines, cl)
	}

	return d, nil
}

// LoadLine6JSON - load json file
func LoadLine6JSON(fn string) (*LineData, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if !isValidLI5(li) {
		return nil, nil
	}

	d := &LineData{}
	for _, v := range li {
		cl := []int{v.R1, v.R2, v.R3, v.R4, v.R5, v.R6}
		d.Lines = append(d.Lines, cl)
	}

	return d, nil
}

// LoadLineDataFromExcel - load xlsx file
func LoadLineDataFromExcel(fn string) (*LineData, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadLineDataFromExcel:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadLineDataFromExcel:GetSheetList",
			goutils.JSON("SheetList", lstname),
			zap.String("fn", fn),
			zap.Error(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadLineDataFromExcel:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	ld := &LineData{}

	// x -> ri
	mapli := make(map[int]int)
	maxli := 0

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				if colCell[0] == 'r' || colCell[0] == 'R' {
					iv, err := goutils.String2Int64(colCell[1:])
					if err != nil {
						goutils.Error("LoadLineDataFromExcel:String2Int64",
							zap.String("fn", fn),
							zap.String("header", colCell),
							zap.Error(err))

						return nil, err
					}

					if iv <= 0 {
						goutils.Error("LoadLineDataFromExcel",
							zap.String("info", "check iv"),
							zap.String("fn", fn),
							zap.String("header", colCell),
							zap.Error(ErrInvalidReelsExcelFile))

						return nil, ErrInvalidReelsExcelFile
					}

					mapli[x] = int(iv) - 1
					if int(iv) > maxli {
						maxli = int(iv)
					}
				}
			}

			if maxli != len(mapli) {
				goutils.Error("LoadLineDataFromExcel",
					zap.String("info", "check len"),
					zap.String("fn", fn),
					zap.Int("maxli", maxli),
					goutils.JSON("mapli", mapli),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if maxli <= 0 {
				goutils.Error("LoadLineDataFromExcel",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("maxli", maxli),
					goutils.JSON("mapli", mapli),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}
		} else {
			cld := make([]int, maxli)

			for x, colCell := range row {
				ri, isok := mapli[x]
				if isok {
					v, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadLineDataFromExcel:String2Int64",
							zap.String("val", colCell),
							zap.Error(err))

						return nil, err
					}

					cld[ri] = int(v)
				}
			}

			ld.Lines = append(ld.Lines, cld)
		}
	}

	return ld, nil
}
