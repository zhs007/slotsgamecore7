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

func (wad *WinAreaData) checkWin(avgWin float64, bet int, options *WinWeightFitOptions) int {
	curaw := wad.calcAvgWin(bet)

	return options.cmpWin(curaw, avgWin)
}

// checkTurn - 判断最小修改后，是否会发生转向，就是小于变大于，如果发生转向，可能会需要整体权重放大
func (wad *WinAreaData) checkTurn(avgWin float64, bet int, options *WinWeightFitOptions, isLess bool, index int, num int, isIgnoreEqu bool) bool {
	wad.Wins[index].Weight += num
	defer func() {
		wad.Wins[index].Weight -= num
	}()

	curaw := wad.calcAvgWin(bet)

	if isLess {
		co := options.cmpWin(curaw, avgWin)
		if isIgnoreEqu {
			if co > 0 {
				return true
			}
		} else {
			if co >= 0 {
				return true
			}
		}

		return false
	}

	co := options.cmpWin(curaw, avgWin)
	if isIgnoreEqu {
		if co < 0 {
			return true
		}
	} else {
		if co <= 0 {
			return true
		}
	}

	return false
}

func (wad *WinAreaData) scaleUp(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []int{}

	// wins经过排序，从小到大，这里lst是从大到小，缩小也要注意维持逻辑一致
	for i := len(wad.Wins) - 1; i >= 0; i-- {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) <= avgWin {
			break
		}

		lst = append(lst, i)
	}

	if len(lst) <= 0 {
		// 前面经过merge，不可能出现这种情况
		goutils.Error("WinAreaData.scaleUp:empty lst",
			zap.Error(ErrWinWeightMerge))

		return false
	}

	for wad.checkTurn(avgWin, bet, options, true, lst[len(lst)-1], 1, false) {
		wad.scale(10)
	}
retry:
	isChg := false

	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, true, i, 1, true) {
			// 直接放弃，下一个
			continue
		}

		n := options.FuncGetDataNum(wad.Wins[i])
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, true, i, n, true) {
				wad.Wins[i].Weight += n
				isChg = true

				if wad.checkWin(avgWin, bet, options) == 0 {
					return true
				}
			} else {
				tn := -1
				for cn := 2; cn < n; cn++ {
					if wad.checkTurn(avgWin, bet, options, true, i, cn, true) {
						tn = cn - 1

						break
					}
				}

				isChg = true

				if tn < 0 {
					wad.Wins[i].Weight += n - 1
				} else {
					wad.Wins[i].Weight += tn
				}

				if wad.checkWin(avgWin, bet, options) == 0 {
					return true
				}
			}
		} else {
			isChg = true

			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if isChg {
		wad.scale(10)

		goto retry
	} else {
		goutils.Error("WinAreaData.scaleUp",
			zap.Error(ErrWinWeightMerge))
	}

	return false
}

func (wad *WinAreaData) scaleDown(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []int{}

	// wins经过排序，从小到大，这里lst是从小到大，逻辑上是由远及近
	for i := 0; i < len(wad.Wins); i++ {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) >= avgWin {
			break
		}

		lst = append(lst, i)
	}

	if len(lst) <= 0 {
		// 前面经过merge，不可能出现这种情况
		goutils.Error("WinAreaData.scaleDown",
			zap.Error(ErrWinWeightMerge))

		return false
	}

	for wad.checkTurn(avgWin, bet, options, false, lst[len(lst)-1], 1, false) {
		wad.scale(10)
	}
retry:
	isChg := false

	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, false, i, 1, true) {
			// 直接放弃，下一个
			continue
		}

		n := options.FuncGetDataNum(wad.Wins[i])
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, false, i, n, true) {
				wad.Wins[i].Weight += n
				isChg = true

				if wad.checkWin(avgWin, bet, options) == 0 {
					return true
				}
			} else {
				tn := -1
				for cn := 2; cn < n; cn++ {
					if wad.checkTurn(avgWin, bet, options, false, i, cn, true) {
						tn = cn - 1

						break
					}
				}

				isChg = true

				if tn < 0 {
					wad.Wins[i].Weight += n - 1
				} else {
					wad.Wins[i].Weight += tn
				}

				if wad.checkWin(avgWin, bet, options) == 0 {
					return true
				}
			}
		} else {
			isChg = true

			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if isChg {
		wad.scale(10)

		goto retry
	} else {
		goutils.Error("WinAreaData.scaleDown",
			zap.Error(ErrWinWeightMerge))
	}

	return false
}

func (wad *WinAreaData) scale(mul int) {
	for _, v := range wad.Wins {
		v.Weight *= mul
	}
}

func (wad *WinAreaData) initWeights(options *WinWeightFitOptions) {
	for _, v := range wad.Wins {
		n := options.FuncGetDataNum(v.Data)
		v.Weight = n
	}
}

func (wad *WinAreaData) Fit(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	wad.initWeights(options)

	curawin := wad.calcAvgWin(bet)

	wo := options.cmpWin(curawin, avgWin)
	if wo == 0 {
		return true
	}

	if wo < 0 {
		return wad.scaleUp(avgWin, bet, options)
	}

	return wad.scaleDown(avgWin, bet, options)
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

func (ww *WinWeight) getMaxIndex() int {
	maxi := 0

	for k := range ww.MapData {
		if k > maxi {
			maxi = k
		}
	}

	return maxi
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
				newi := wd.mergeAvgWins(si, i)
				ww.merge(si, i, newi)

				return newi, nil
			}

			if i == maxi {
				maxj := ww.getMaxIndex()
				for j := i + 1; j <= maxj; j++ {
					ret0, ret1 := ww.isValidData(si, j, aw, bet, options)
					if !ret0 {
						goutils.Error("WinWeight.mergeNext:less",
							zap.Error(ErrWinWeightMerge))

						return -1, ErrWinWeightMerge
					}

					if ret1 {
						newi := wd.mergeAvgWins(si, i)
						ww.merge(si, j, newi)

						return newi, nil
					}
				}
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

			if ret0 && ret1 {
				lasti = si
				si = i + 1

				continue
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

				lasti = si
				si = ni + 1
			} else if !ret0 {
				newi := wd.mergeAvgWins(si, i)
				ww.merge(si, i, newi)
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

func (ww *WinWeight) merge(mini, maxi int, newi int) error {
	if !(newi >= mini && newi <= maxi) {
		goutils.Error("WinWeight.merge",
			zap.Int("min index", mini),
			zap.Int("max index", maxi),
			zap.Int("new index", newi),
			zap.Error(ErrWinWeightMerge))

		return ErrWinWeightMerge
	}

	nwad := &WinAreaData{}

	for i := mini; i <= maxi; i++ {
		v, isok := ww.MapData[i]
		if isok {
			nwad.Wins = append(nwad.Wins, v.Wins...)

			delete(ww.MapData, i)
		}
	}

	ww.MapData[newi] = nwad

	return nil
}

func NewWinWeight() *WinWeight {
	return &WinWeight{
		MapData: make(map[int]*WinAreaData),
	}
}
