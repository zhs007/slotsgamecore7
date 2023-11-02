package mathtoolset

import (
	"math"
	"sort"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type WinWeightFitOptions struct {
	MapDataNum     map[int]int
	DataNum        int
	FuncGetDataNum func(any) int
	FuncSetWeight  func(any, int)
	WinScale       int
	MaxFitTimes    int
	MinNodes       int // merge时，某一边节点数低于这个就需要merge
	MinSeeds       int // merge时，某一边seed数低于这个就需要merge
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

func (ww *WinWeight) isValidData(si int, ci int, avgwin float64, bet int, options *WinWeightFitOptions) (bool, bool) {
	lessnum := 0
	lessseednum := 0
	bignum := 0
	bigseednum := 0
	equnum := 0

	for i := si; i <= ci; i++ {
		wd, isok := ww.MapData[i]
		if isok {
			for _, v := range wd.Wins {
				n := options.FuncGetDataNum(v.Data)

				cw := float64(v.Win) / float64(bet)
				if cw > avgwin {
					bignum++
					bigseednum += n
				} else if cw < avgwin {
					lessnum++
					lessseednum += n
				} else {
					equnum++
				}
			}
		}
	}

	if bignum == 0 && lessnum == 0 && equnum > 0 {
		return true, true
	}

	ret0 := true
	ret1 := true

	if lessnum < options.MinNodes || lessseednum < options.MinSeeds {
		ret0 = false
	}

	if bignum < options.MinNodes || bigseednum < options.MinSeeds {
		ret1 = false
	}

	return ret0, ret1
}

func (ww *WinWeight) mergeNext(wd *WinningDistribution, bet int, options *WinWeightFitOptions, si int, ci int, maxi int) (int, error) {
	for i := ci; i <= maxi; i++ {
		_, isok := wd.AvgWins[i]
		if isok {
			aw := wd.getAvgWin(si, i)

			ret0, ret1 := ww.isValidData(si, i, aw, bet, options)
			if !ret0 {
				goutils.Error("WinWeight.mergeNext:less",
					zap.Error(ErrWinWeightMerge))

				return -1, ErrWinWeightMerge
			}

			if ret1 {
				ww.merge(si, i)
				wd.mergeAvgWins(si, i)

				return i, nil
			}
		}
	}

	goutils.Error("WinWeight.mergeNext:unless",
		zap.Error(ErrWinWeightMerge))

	return -1, ErrWinWeightMerge
}

func (ww *WinWeight) mergeWith(wd *WinningDistribution, bet int, options *WinWeightFitOptions) error {
	si := 0
	lasti := 0
	maxi := wd.getMax()

	for i := 0; i <= maxi; i++ {
		_, isok := wd.AvgWins[i]
		if isok {
			aw := wd.getAvgWin(si, i)

			ret0, ret1 := ww.isValidData(si, i, aw, bet, options)
			if si == 0 && !ret0 {
				goutils.Error("WinWeight.mergeWith:0",
					zap.Error(ErrWinWeightMerge))

				return ErrWinWeightMerge
			}

			if !ret0 {
				si = lasti
			}

			if !ret1 {
				ni, err := ww.mergeNext(wd, bet, options, si, i, maxi)
				if err != nil {
					goutils.Error("WinWeight.mergeWith:mergeNext",
						zap.Int("si", si),
						zap.Int("i", i),
						zap.Error(err))

					return err
				}

				i = ni
				si = ni + 1
			} else if !ret0 {
				ww.merge(si, i)
				wd.mergeAvgWins(si, i)
			}
		}
	}

	return nil
}

func (ww *WinWeight) Fit(wd *WinningDistribution, bet int, options *WinWeightFitOptions) (*WinWeight, error) {
	err := ww.mergeWith(wd, bet, options)
	if err != nil {
		goutils.Error("WinWeight.Fit:mergeWith",
			zap.Error(err))

		return nil, err
	}

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
