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
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "bet")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "rtp")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), wins.TotalWin)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), totalBet)
	if totalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(wins.TotalWin)/float64(totalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
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

	f.SetCellValue(sheet, goutils.Pos2Cell(3, 5), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(4, 5), "times")
	f.SetCellValue(sheet, goutils.Pos2Cell(5, 5), "trigger chance")
	f.SetCellValue(sheet, goutils.Pos2Cell(6, 5), "total wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(7, 5), "rtp")

	y := 6
	for _, k := range lstwins {
		v := wins.MapWinTimes[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(3, y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(4, y), v)

		if totalTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(5, y), float64(v)/float64(totalTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(5, y), 0)
		}

		f.SetCellValue(sheet, goutils.Pos2Cell(6, y), k*v)

		if totalBet > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(7, y), float64(k*v)/float64(totalBet))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(7, y), 0)
		}

		y++
	}
}

func NewStatsWins() *StatsWins {
	return &StatsWins{
		MapWinTimes: make(map[int]int),
	}
}
