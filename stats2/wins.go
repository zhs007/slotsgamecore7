package stats2

import (
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsWins struct {
	TotalWin    int64       `json:"totalWin"`
	MapWinTimes map[int]int `json:"mapWinTimes"`
}

func (wins *StatsWins) AddWin(win int64) {
	wins.TotalWin += win
	wins.MapWinTimes[int(win)]++
}

// func (wins *StatsWins) Clone() *StatsWins {
// 	return &StatsWins{
// 		TotalWin: wins.TotalWin,
// 		TotalBet: wins.TotalBet,
// 	}
// }

func (wins *StatsWins) Merge(src *StatsWins) {
	wins.TotalWin += src.TotalWin

	for k, v := range src.MapWinTimes {
		wins.MapWinTimes[k] += v
	}
}

func (wins *StatsWins) SaveSheet(f *excelize.File, sheet string, totalBet int64) {
	wins.saveSheet(f, sheet, 0, 0, totalBet)
}

func (wins *StatsWins) saveSheet(f *excelize.File, sheet string, sx, sy int, totalBet int64) {
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+0, sy+0), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+0, sy+1), "bet")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+0, sy+2), "rtp")

	f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+0), wins.TotalWin)
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+1), totalBet)
	if totalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+2), float64(wins.TotalWin)/float64(totalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+2), 0)
	}

	totalTimes := int64(0)
	lstwins := []int{}
	for k, v := range wins.MapWinTimes {
		totalTimes += int64(v)
		lstwins = append(lstwins, k)
	}

	sort.Slice(lstwins, func(i, j int) bool {
		return lstwins[i] < lstwins[j]
	})

	f.SetCellValue(sheet, goutils.Pos2Cell(sx+3, sy+5), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+4, sy+5), "times")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+5), "trigger chance")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+6, sy+5), "total wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+7, sy+5), "rtp")

	y := 6
	for _, k := range lstwins {
		v := wins.MapWinTimes[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+3, sy+y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+4, sy+y), v)

		if totalTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), float64(v)/float64(totalTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), 0)
		}

		f.SetCellValue(sheet, goutils.Pos2Cell(sx+6, sy+y), k*v)

		if totalBet > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+7, sy+y), float64(k*v)/float64(totalBet))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+7, sy+y), 0)
		}

		y++
	}
}

func NewStatsWins() *StatsWins {
	return &StatsWins{
		MapWinTimes: make(map[int]int),
	}
}
