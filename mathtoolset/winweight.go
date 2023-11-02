package mathtoolset

import (
	"math"
	"sort"
)

type WinWeightFitOptions struct {
	MapDataNum     map[int]int
	DataNum        int
	FuncGetDataNum func(any) int
	FuncSetWeight  func(any, int)
	WinScale       int
	MaxFitTimes    int
}

func (wwfo *WinWeightFitOptions) cmpWin(win0, win1 float64) int {
	w0 := int64(win0 * float64(wwfo.WinScale))
	w1 := int64(win1 * float64(wwfo.WinScale))

	if w0 == w1 {
		return 0
	}

	if w0 > w1 {
		return 1
	}

	return -1
}

type WinData struct {
	Win    int `yaml:"win" json:"win"`
	Weight int `yaml:"weight" json:"weight"`
	Data   any `yaml:"data" json:"data"`
}

type WinAreaData struct {
	Wins         []*WinData `yaml:"wins" json:"wins"`
	TotalWeights int        `yaml:"maxWeight" json:"maxWeight"`
	AvgWin       float64    `yaml:"avgWin" json:"avgWin"`
	Percent      float64    `yaml:"percent" json:"percent"`
}

func (wad *WinAreaData) calcAvgWin(bet int) float64 {
	totalwin := float64(0)
	totalweights := 0

	for _, v := range wad.Wins {
		totalwin += float64(v.Win) / float64(bet) * float64(v.Weight)
		totalweights += v.Weight
	}

	return totalwin / float64(totalweights)
}

func (wad *WinAreaData) checkUp(wd *WinData, avgWin float64, bet int, maxweight int) bool {
	return false
}

func (wad *WinAreaData) up(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []*WinData{}

	for i := len(wad.Wins) - 1; i >= 0; i-- {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) <= avgWin {
			break
		}

		lst = append(lst, v)
	}

	if len(lst) <= 0 {
		return false
	}

	for _, v := range lst {
		n := options.FuncGetDataNum(v.Data)
		if n > 1 {
			wad.checkUp(v, avgWin, bet, n)
		}
	}

	return false

	// v := wad.Wins[lasti]

	// if float64(v.Win)/float64(bet) < avgWin {
	// 	return -1
	// }

	// n := options.FuncGetDataNum(v.Data)

	// if n > 1 {

	// }

	// return lasti - 1
}

func (wad *WinAreaData) Fit(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	for _, v := range wad.Wins {
		n := options.FuncGetDataNum(v.Data)
		v.Weight = n
	}

	curawin := wad.calcAvgWin(bet)
	times := 0
	// lasti := -1

	for {
		wo := options.cmpWin(curawin, avgWin)
		if wo == 0 {
			break
		}

		if wo < 0 {

		}

		times++

		if times >= options.MaxFitTimes {
			return false
		}
	}

	return true
}

type WinWeight struct {
	MapData map[int]*WinAreaData
}

func (ww *WinWeight) Add(win int, bet int, data any) {
	k := -1
	if win > 0 {
		fwin := float64(win) / float64(bet)

		k = int(math.Floor(fwin))
	}

	v, isok := ww.MapData[k]
	if isok {
		v.Wins = append(v.Wins, &WinData{
			Win:  win,
			Data: data,
		})
	} else {
		v = &WinAreaData{}

		v.Wins = append(v.Wins, &WinData{
			Win:  win,
			Data: data,
		})

		ww.MapData[k] = v
	}
}

func (ww *WinWeight) sort() {
	for _, v := range ww.MapData {
		sort.Slice(v.Wins, func(i, j int) bool {
			return v.Wins[i].Win < v.Wins[j].Win
		})
	}
}

func (ww *WinWeight) Fit(wd *WinningDistribution, bet int, options *WinWeightFitOptions) (*WinWeight, error) {
	ww.sort()

	target := NewWinWeight()

	for k, v := range wd.AvgWins {
		wwv, isok := ww.MapData[k]
		if isok {
			if wwv.Fit(v.AvgWin, bet, options) {
				delete(wd.AvgWins, k)
			}
		}
	}

	return target, nil
}

func (ww *WinWeight) merge(mini, maxi int) {
	nwad := &WinAreaData{}

	for i := mini; i <= maxi; i++ {
		v, isok := ww.MapData[i]
		if isok {
			nwad.Wins = append(nwad.Wins, v.Wins...)
		}
	}
}

func NewWinWeight() *WinWeight {
	return &WinWeight{
		MapData: make(map[int]*WinAreaData),
	}
}
