package stats

import (
	"sort"
	"sync"

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
	//sync.map is map[interface{}]interface{}
	//we not sure the front end send to us 
	sync *sync.RWMutex
}

func (wins *Wins) Clone() *Wins {
	nw := &Wins{
		MapWins:    make(map[int]int64),
		TotalTimes: wins.TotalTimes,
	}
	wins.sync.Lock()
    defer wins.sync.Unlock()
	for k, v := range wins.MapWins {
		wins.MapWins[k] = v
	}

	return nw
}

func (wins *Wins) Merge(src *Wins) {
	wins.TotalTimes += src.TotalTimes

	wins.sync.Lock()
    defer wins.sync.Unlock()
	for k, v := range src.MapWins {
		_, isok := wins.MapWins[k]
		if isok {
			wins.MapWins[k] += v
		} else {
			wins.MapWins[k] = v
		}
	}
}

func (wins *Wins) genData() []*winsdata {
	lst := []*winsdata{}

	wins.sync.Lock()
    defer wins.sync.Unlock()
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

	wins.sync.Lock()
    defer wins.sync.Unlock()
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

	wins.sync.Lock()
    defer wins.sync.Unlock()
	
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
