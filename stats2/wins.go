package stats2

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsWins struct {
	TotalWin int64
	TotalBet int64
}

func (wins *StatsWins) Clone() *StatsWins {
	return &StatsWins{
		TotalWin: wins.TotalWin,
		TotalBet: wins.TotalBet,
	}
}

func (wins *StatsWins) Merge(src *StatsWins) {
	wins.TotalBet += src.TotalBet
	wins.TotalWin += src.TotalWin
}

func (wins *StatsWins) SaveSheet(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "bet")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "rtp")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), wins.TotalWin)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), wins.TotalBet)
	if wins.TotalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(wins.TotalWin)/float64(wins.TotalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}
}
