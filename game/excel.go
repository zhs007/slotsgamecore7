package sgc7game

import (
	"io"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type FuncProcHeader func(x int, str string) string
type FuncProcData func(x int, y int, header string, data string) error

func LoadExcel(fn string, sheet string, onheader FuncProcHeader, ondata FuncProcData) error {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadExcel:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}
	defer f.Close()

	if sheet == "" {
		sheet = f.GetSheetName(0)
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadExcel:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	mapcolname := make(map[int]string)

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = onheader(x, colCell)
			}
		} else {
			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					err := ondata(x, y, colname, colCell)
					if err != nil {
						goutils.Error("LoadExcel:ondata",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("header", colname),
							zap.String("val", colCell),
							zap.Error(err))

						return err
					}
				}
			}
		}
	}

	return nil
}

func LoadExcelWithReader(reader io.Reader, sheet string, onheader FuncProcHeader, ondata FuncProcData) error {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		goutils.Error("LoadExcelWithReader:OpenReader",
			zap.Error(err))

		return err
	}
	defer f.Close()

	if sheet == "" {
		sheet = f.GetSheetName(0)
	}

	rows, err := f.GetRows(sheet)
	if err != nil {
		goutils.Error("LoadExcelWithReader:GetRows",
			zap.Error(err))

		return err
	}

	mapcolname := make(map[int]string)

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = onheader(x, colCell)
			}
		} else {
			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					err := ondata(x, y, colname, colCell)
					if err != nil {
						goutils.Error("LoadExcelWithReader:ondata",
							zap.Int("x", x),
							zap.Int("y", y),
							zap.String("header", colname),
							zap.String("val", colCell),
							zap.Error(err))

						return err
					}
				}
			}
		}
	}

	return nil
}
