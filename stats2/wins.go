package stats2

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsWins struct {
	TotalWin      int64          `json:"totalWin"`
	MapWinTimes   map[int]int    `json:"mapWinTimes"`
	MapWinTimesEx map[string]int `json:"MapWinTimesEx"`
	MapWinEx      map[string]int `json:"MapWinEx"`
}

func (wins *StatsWins) AddWin(win int64) {
	wins.TotalWin += win
	wins.MapWinTimes[int(win)]++
}

func (wins *StatsWins) genRange(totalBet int) {
	wins.MapWinTimesEx = make(map[string]int)
	wins.MapWinEx = make(map[string]int)

	for win, times := range wins.MapWinTimes {
		if win == 0 {
			wins.MapWinTimesEx["noWins"] += times
			wins.MapWinEx["noWins"] = 0
		} else {
			curWins := float64(win) / float64(totalBet)

			for i := 0; i < len(gWinRange); i++ {
				if curWins > float64(gWinRange[i]) {
					if i < len(gWinRange)-1 && curWins <= float64(gWinRange[i+1]) {
						k := fmt.Sprintf("(%v,%v]", gWinRange[i], gWinRange[i+1])
						wins.MapWinTimesEx[k] += times
						wins.MapWinEx[k] += win * times

						break
					} else if i == len(gWinRange)-1 {
						k := fmt.Sprintf(">%v", gWinRange[i])
						wins.MapWinTimesEx[k] += times
						wins.MapWinEx[k] += win * times

						break
					}
				}
			}
		}
	}
}

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

	wins.genRange(int(totalBet))

	tx := sx + 8

	f.SetCellValue(sheet, goutils.Pos2Cell(tx+1, sy+5), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(tx+2, sy+5), "times")
	f.SetCellValue(sheet, goutils.Pos2Cell(tx+3, sy+5), "trigger chance")
	f.SetCellValue(sheet, goutils.Pos2Cell(tx+4, sy+5), "total wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(tx+5, sy+5), "rtp")

	{
		y = 6
		k := "noWins"

		v := wins.MapWinTimesEx[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(tx+1, sy+y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(tx+2, sy+y), v)

		if totalTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+3, sy+y), float64(v)/float64(totalTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+3, sy+y), 0)
		}

		f.SetCellValue(sheet, goutils.Pos2Cell(tx+4, sy+y), 0)

		if totalBet > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+5, sy+y), 0)
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+5, sy+y), 0)
		}
	}

	y++
	for i := 0; i < len(gWinRange); i++ {
		var k string
		if i < len(gWinRange)-1 {
			k = fmt.Sprintf("(%v,%v]", gWinRange[i], gWinRange[i+1])
		} else {
			k = fmt.Sprintf(">%v", gWinRange[i])
		}

		v := wins.MapWinTimesEx[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(tx+1, sy+y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(tx+2, sy+y), v)

		if totalTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+3, sy+y), float64(v)/float64(totalTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+3, sy+y), 0)
		}

		f.SetCellValue(sheet, goutils.Pos2Cell(tx+4, sy+y), wins.MapWinEx[k])

		if totalBet > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+5, sy+y), float64(wins.MapWinEx[k])/float64(totalBet))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(tx+5, sy+y), 0)
		}

		y++
	}
}

func NewStatsWins() *StatsWins {
	return &StatsWins{
		MapWinTimes: make(map[int]int),
	}
}
