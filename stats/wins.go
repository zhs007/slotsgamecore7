package stats

import (
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type winsdata struct {
	win   int
	times int64
}

type Wins struct {
	MapWins    map[int]int64
	TotalTimes int64
}

func (wins *Wins) genData() []*winsdata {
	lst := []*winsdata{}

	for k, v := range wins.MapWins {
		lst = append(lst, &winsdata{
			win:   k,
			times: v,
		})
	}

	return lst
}

func (wins *Wins) AddWin(win int) {
	wins.TotalTimes++

	_, isok := wins.MapWins[win]
	if isok {
		wins.MapWins[win]++
	} else {
		wins.MapWins[win] = 1
	}
}

func (wins *Wins) SaveSheet(f *excelize.File, sheet string) error {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), "times")

	lst := wins.genData()

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].win < lst[j].win
	})

	y := 1
	for _, v := range lst {
		f.SetCellValue(sheet, goutils.Pos2Cell(0, y), float64(v.win)/100.0)
		f.SetCellValue(sheet, goutils.Pos2Cell(1, y), v.times)

		y++
	}

	return nil
}

func NewWins() *Wins {
	return &Wins{
		MapWins: make(map[int]int64),
	}
}
