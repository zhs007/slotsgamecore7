package stats

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type FuncOnSaveExcel func(*excelize.File) error

func SaveExcel(fn string, lst []*Feature, onSave FuncOnSaveExcel) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "gamemod")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "playtimes")
	f.SetCellStr(sheet, goutils.Pos2Cell(2, 0), "bet")
	f.SetCellStr(sheet, goutils.Pos2Cell(3, 0), "wins")
	f.SetCellStr(sheet, goutils.Pos2Cell(4, 0), "rtp")
	f.SetCellStr(sheet, goutils.Pos2Cell(5, 0), "triggertimes")
	f.SetCellStr(sheet, goutils.Pos2Cell(6, 0), "hit rate")

	y := 1

	for _, v := range lst {
		f.SetCellValue(sheet, goutils.Pos2Cell(0, y), v.Name)
		f.SetCellValue(sheet, goutils.Pos2Cell(1, y), v.PlayTimes)
		f.SetCellValue(sheet, goutils.Pos2Cell(2, y), v.TotalBets)
		f.SetCellValue(sheet, goutils.Pos2Cell(3, y), v.TotalWins)
		f.SetCellValue(sheet, goutils.Pos2Cell(4, y), float64(v.TotalWins)/float64(v.TotalBets))
		f.SetCellValue(sheet, goutils.Pos2Cell(5, y), v.TriggerTimes)
		f.SetCellValue(sheet, goutils.Pos2Cell(6, y), float64(v.TriggerTimes)/float64(v.PlayTimes))

		y++
	}

	if onSave != nil {
		onSave(f)
	}

	return f.SaveAs(fn)
}
