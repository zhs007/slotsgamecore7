package sgc7game

import (
	"os"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type payInfo struct {
	Code   int    `json:"Code"`
	Symbol string `json:"Symbol"`
	X1     int    `json:"X1"`
	X2     int    `json:"X2"`
	X3     int    `json:"X3"`
	X4     int    `json:"X4"`
	X5     int    `json:"X5"`
	X6     int    `json:"X6"`
	X7     int    `json:"X7"`
	X8     int    `json:"X8"`
	X9     int    `json:"X9"`
	X10    int    `json:"X10"`
	X11    int    `json:"X11"`
	X12    int    `json:"X12"`
	X13    int    `json:"X13"`
	X14    int    `json:"X14"`
	X15    int    `json:"X15"`
	X16    int    `json:"X16"`
	X17    int    `json:"X17"`
	X18    int    `json:"X18"`
	X19    int    `json:"X19"`
	X20    int    `json:"X20"`
	X21    int    `json:"X21"`
	X22    int    `json:"X22"`
	X23    int    `json:"X23"`
	X24    int    `json:"X24"`
	X25    int    `json:"X25"`
}

// PayTables - pay tables
type PayTables struct {
	MapPay     map[int][]int  `json:"paytables"`
	MapSymbols map[string]int `json:"symbols"`
}

func (pt *PayTables) GetStringFromInt(s int) string {
	for k, v := range pt.MapSymbols {
		if v == s {
			return k
		}
	}

	return ""
}

// LoadPayTables5JSON - load json file
func LoadPayTables5JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables3JSON - load json file
func LoadPayTables3JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables6JSON - load json file
func LoadPayTables6JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables15JSON - load json file
func LoadPayTables15JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6, v.X7, v.X8, v.X9, v.X10, v.X11, v.X12, v.X13, v.X14, v.X15}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPayTables25JSON - load json file
func LoadPayTables25JSON(fn string) (*PayTables, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var li []payInfo
	err = json.Unmarshal(data, &li)
	if err != nil {
		return nil, err
	}

	if len(li) <= 0 {
		return nil, nil
	}

	p := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	for _, v := range li {
		cl := []int{v.X1, v.X2, v.X3, v.X4, v.X5, v.X6, v.X7, v.X8, v.X9, v.X10, v.X11, v.X12, v.X13, v.X14, v.X15, v.X16, v.X17, v.X18, v.X19, v.X20, v.X21, v.X22, v.X23, v.X24, v.X25}
		p.MapPay[v.Code] = cl

		p.MapSymbols[v.Symbol] = v.Code
	}

	return p, nil
}

// LoadPaytablesFromExcel - load xlsx file
func LoadPaytablesFromExcel(fn string) (*PayTables, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadPaytablesFromExcel:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}
	defer f.Close()

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadPaytablesFromExcel:GetSheetList",
			goutils.JSON("SheetList", lstname),
			zap.String("fn", fn),
			zap.Error(ErrInvalidReelsExcelFile))

		return nil, ErrInvalidReelsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadPaytablesFromExcel:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	pt := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	// x -> ri
	mapli := make(map[int]int)
	maxli := 0
	codeIndex := -1
	symbolIndex := -1

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				if strings.ToLower(colCell) == "code" {
					codeIndex = x
				} else if strings.ToLower(colCell) == "symbol" {
					symbolIndex = x
				} else if colCell[0] == 'x' || colCell[0] == 'X' {
					iv, err := goutils.String2Int64(colCell[1:])
					if err != nil {
						goutils.Error("LoadPaytablesFromExcel:String2Int64",
							zap.String("fn", fn),
							zap.String("header", colCell),
							zap.Error(err))

						return nil, err
					}

					if iv <= 0 {
						goutils.Error("LoadPaytablesFromExcel",
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
				goutils.Error("LoadPaytablesFromExcel",
					zap.String("info", "check len"),
					zap.String("fn", fn),
					zap.Int("maxli", maxli),
					goutils.JSON("mapli", mapli),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if maxli <= 0 {
				goutils.Error("LoadPaytablesFromExcel",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("maxli", maxli),
					goutils.JSON("mapli", mapli),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if codeIndex < 0 {
				goutils.Error("LoadPaytablesFromExcel",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("codeIndex", codeIndex),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}

			if symbolIndex < 0 {
				goutils.Error("LoadPaytablesFromExcel",
					zap.String("info", "check empty"),
					zap.String("fn", fn),
					zap.Int("symbolIndex", symbolIndex),
					zap.Error(ErrInvalidReelsExcelFile))

				return nil, ErrInvalidReelsExcelFile
			}
		} else {
			cpd := make([]int, maxli)
			code := 0
			symbol := ""

			for x, colCell := range row {
				if x == codeIndex {
					v, err := goutils.String2Int64(colCell)
					if err != nil {
						goutils.Error("LoadPaytablesFromExcel:String2Int64",
							zap.String("val", colCell),
							zap.Error(err))

						return nil, err
					}

					code = int(v)
				} else if x == symbolIndex {
					symbol = colCell
				} else {
					ri, isok := mapli[x]
					if isok {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadPaytablesFromExcel:String2Int64",
								zap.String("val", colCell),
								zap.Error(err))

							return nil, err
						}

						cpd[ri] = int(v)
					}
				}
			}

			pt.MapPay[code] = cpd
			pt.MapSymbols[symbol] = code
		}
	}

	return pt, nil
}
