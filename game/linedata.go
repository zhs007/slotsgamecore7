package sgc7game

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/xuri/excelize/v2"
	goutils "github.com/zhs007/goutils"
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
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = sonic.Unmarshal(data, &li)
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
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = sonic.Unmarshal(data, &li)
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
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []lineInfo
	err = sonic.Unmarshal(data, &li)
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
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}
	defer f.Close()

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadLineDataFromExcel:GetSheetList",
			slog.Any("SheetList", lstname),
			slog.String("fn", fn),
			goutils.Err(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadLineDataFromExcel:GetRows",
			slog.String("fn", fn),
			goutils.Err(err))

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
							slog.String("fn", fn),
							slog.String("header", colCell),
							goutils.Err(err))

						return nil, err
					}

					if iv <= 0 {
						goutils.Error("LoadLineDataFromExcel",
							slog.String("info", "check iv"),
							slog.String("fn", fn),
							slog.String("header", colCell),
							goutils.Err(ErrInvalidReelsExcelFile))

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
					slog.String("info", "check len"),
					slog.String("fn", fn),
					slog.Int("maxli", maxli),
					slog.Any("mapli", mapli),
					goutils.Err(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if maxli <= 0 {
				goutils.Error("LoadLineDataFromExcel",
					slog.String("info", "check empty"),
					slog.String("fn", fn),
					slog.Int("maxli", maxli),
					slog.Any("mapli", mapli),
					goutils.Err(ErrInvalidReelsExcelFile))

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
							slog.String("val", colCell),
							goutils.Err(err))

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
