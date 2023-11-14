package mathtoolset

import (
	"math"
	"sort"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type WinWeightFitOptions struct {
	FuncGetDataNum func(any) int
	FuncSetWeight  func(any, int)
	WinScale       int // 在fit时，这个是比较精度，一般给100即可，就是精确到0.01
	RTPScale       int // 在fitending时，用这个替代winscale，一般来说，这个最好要精确到0.0001
	MaxFitTimes    int
	MinNodes       int // merge时，某一边节点数低于这个就需要merge
	MinSeeds       int // merge时，某一边seed数低于这个就需要merge
	TotalWeight    int // 总权重
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

func (wd *WinData) Clone() *WinData {
	return &WinData{
		Win:    wd.Win,
		Weight: wd.Weight,
		Data:   wd.Data,
	}
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

// checkTurnEx - 判断所有元素加权重后，是否会发生反转
func (wad *WinAreaData) checkTurnEx(avgWin float64, bet int, options *WinWeightFitOptions, isLess bool, lst []int, num int, isIgnoreEqu bool) bool {
	for _, i := range lst {
		wad.Wins[i].Weight += num
	}

	defer func() {
		for _, i := range lst {
			wad.Wins[i].Weight -= num
		}
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

	// wins经过排序，从小到大，这里会break，但希望lst从小到大，所以写法要注意
	for i := len(wad.Wins) - 1; i >= 0; i-- {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) <= avgWin {
			break
		}

		lst = append([]int{i}, lst...)
	}

	if len(lst) <= 0 {
		// 前面经过merge，不可能出现这种情况
		goutils.Error("WinAreaData.scaleUp:empty lst",
			zap.Error(ErrWinWeightScale))

		return false
	}

	for wad.checkTurn(avgWin, bet, options, true, lst[0], 1, false) {
		wad.scale(10)
	}

	retrynum := 0
retry:
	isneedscale := false
	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, true, i, 1, true) {
			// 因为排序，所以直接break
			isneedscale = true

			break
		}

		n := options.FuncGetDataNum(wad.Wins[i].Data)
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, true, i, n, true) {
				wad.Wins[i].Weight += n

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
			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if !isneedscale {
		goto retry
	}

	if retrynum < options.MaxFitTimes {
		wad.scale(10)
		retrynum++

		goto retry
	} else {
		goutils.Error("WinAreaData.scaleUp",
			zap.Error(ErrWinWeightScale))
	}

	return false
}

func (wad *WinAreaData) scaleDown(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []int{}

	// wins经过排序，从小到大，这里lst是从大到小，逻辑上是由近及远
	for i := 0; i < len(wad.Wins); i++ {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) >= avgWin {
			break
		}

		lst = append([]int{i}, lst...)
	}

	if len(lst) <= 0 {
		// 前面经过merge，不可能出现这种情况
		goutils.Error("WinAreaData.scaleDown",
			zap.Error(ErrWinWeightScale))

		return false
	}

	for wad.checkTurn(avgWin, bet, options, false, lst[0], 1, false) {
		wad.scale(10)
	}

	retrynum := 0
retry:
	isneedscale := false
	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, false, i, 1, true) {
			// 直接放弃，下一轮
			isneedscale = true

			break
		}

		n := options.FuncGetDataNum(wad.Wins[i].Data)
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, false, i, n, true) {
				wad.Wins[i].Weight += n

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
			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if !isneedscale {
		goto retry
	}

	if retrynum < options.MaxFitTimes {
		wad.scale(10)
		retrynum++

		goto retry
	} else {
		goutils.Error("WinAreaData.scaleDown",
			zap.Error(ErrWinWeightScale))
	}

	return false
}

func (wad *WinAreaData) scaleUp2(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []int{}

	// wins经过排序，从小到大，这里会break，但希望lst从小到大，所以写法要注意
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
			zap.Error(ErrWinWeightScale))

		return false
	}

	for wad.checkTurn(avgWin, bet, options, true, lst[len(lst)-1], 1, false) {
		wad.scale(10)
	}

	retrynum := 0
retry:
	isneedscale := false
	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, true, i, 1, true) {
			// 因为排序，所以直接break
			isneedscale = true

			continue
		}

		n := options.FuncGetDataNum(wad.Wins[i].Data)
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, true, i, n, true) {
				wad.Wins[i].Weight += n

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
			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if !isneedscale {
		goto retry
	}

	if retrynum < options.MaxFitTimes {
		wad.scale(10)
		retrynum++

		goto retry
	} else {
		goutils.Error("WinAreaData.scaleUp",
			zap.Error(ErrWinWeightScale))
	}

	return false
}

func (wad *WinAreaData) findWinWithMinWeight(win int) int {
	curweight := -1
	curi := -1

	for i, v := range wad.Wins {
		if v.Win == win {
			if curi == -1 {
				curi = i
				curweight = v.Weight
			} else {
				if v.Weight < curweight {
					curi = i
					curweight = v.Weight
				}
			}
		}
	}

	return curi
}

func (wad *WinAreaData) scaleUpEnding(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	// 最后的缩放逻辑，为了拟合rtp，这里不能再整体放大倍数了
	// 因为是最后的缩放了，所以从近端开始，而且一个win，一次只放大权重最小的1个

	lst := []int{}

	// wins经过排序，从小到大，这里会break，但希望lst从小到大，所以写法要注意
	for i := len(wad.Wins) - 1; i >= 0; i-- {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) <= avgWin {
			break
		}

		lst = append([]int{i}, lst...)
	}

	if len(lst) <= 0 {
		// 最后的缩放，如果lst为空就没办法了
		goutils.Error("WinAreaData.scaleUpEnding:empty lst",
			zap.Error(ErrWinWeightScale))

		return false
	}

	if wad.checkTurn(avgWin, bet, options, true, lst[0], 1, false) {
		// 最后的缩放，如果最近端都不能放大，则算失败
		goutils.Error("WinAreaData.scaleUpEnding:check 0 cannot scaleup",
			zap.Error(ErrWinWeightScale))

		return false
	}

	for tn := options.TotalWeight / 100; tn > 100; tn /= 10 {
		for !wad.checkTurnEx(avgWin, bet, options, true, lst, tn, false) {
			for _, i := range lst {
				wad.Wins[i].Weight += tn
			}
		}
	}

retry:
	chgnum := 0
	prewin := -1

	for _, i := range lst {
		// 放大时则一档先只变一个
		if wad.Wins[i].Win <= prewin {
			continue
		}

		prewin = wad.Wins[i].Win

		// 先找到这个win里权重最小的
		curi := wad.findWinWithMinWeight(prewin)

		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, true, curi, 1, true) {
			// 因为排序，所以直接break
			break
		}

		chgnum++
		wad.Wins[curi].Weight++

		if wad.checkWin(avgWin, bet, options) == 0 {
			return true
		}
	}

	if chgnum > 0 {
		goto retry
	}

	goutils.Error("WinAreaData.scaleUpEnding:cannot scaleup",
		zap.Error(ErrWinWeightScale))

	return false
}

func (wad *WinAreaData) scaleDownEnding(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	// 最后的缩放逻辑，为了拟合rtp，这里不能再整体放大倍数了
	// 因为是最后的缩放了，所以从近端开始，而且一个win，一次只放大权重最小的1个
	lst := []int{}

	// wins经过排序，从小到大，这里lst是从大到小，逻辑上是由近及远
	for i := 0; i < len(wad.Wins); i++ {
		v := wad.Wins[i]

		if float64(v.Win)/float64(bet) >= avgWin {
			break
		}

		lst = append([]int{i}, lst...)
	}

	if len(lst) <= 0 {
		// 最后的缩放，如果lst为空就没办法了
		goutils.Error("WinAreaData.scaleDownEnding:empty lst",
			zap.Error(ErrWinWeightScale))

		return false
	}

	if wad.checkTurn(avgWin, bet, options, false, lst[0], 1, false) {
		// 最后的缩放，如果最近端都不能放大，则算失败
		goutils.Error("WinAreaData.scaleDownEnding:check 0 cannot scaleup",
			zap.Error(ErrWinWeightScale))

		return false
	}

	for tn := options.TotalWeight / 100; tn > 100; tn /= 10 {
		for !wad.checkTurnEx(avgWin, bet, options, false, lst, tn, false) {
			for _, i := range lst {
				wad.Wins[i].Weight += tn
			}
		}
	}

retry:
	chgnum := 0
	prewin := -1

	for _, i := range lst {
		// 一般情况下，缩小时改的数据量会比较小，所以每个都会变，而放大时则一档先只变一个
		if prewin >= 0 && wad.Wins[i].Win > prewin {
			continue
		}

		prewin = wad.Wins[i].Win

		// 先找到这个win里权重最小的
		curi := wad.findWinWithMinWeight(prewin)

		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, false, curi, 1, true) {
			// 因为排序，所以直接break
			break
		}

		chgnum++
		wad.Wins[curi].Weight++

		if wad.checkWin(avgWin, bet, options) == 0 {
			return true
		}
	}

	if chgnum > 0 {
		goto retry
	}

	goutils.Error("WinAreaData.scaleDownEnding:cannot scaleup",
		zap.Error(ErrWinWeightScale))

	return false
}

func (wad *WinAreaData) scaleDown2(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	lst := []int{}

	// wins经过排序，从小到大，这里lst是从大到小，逻辑上是由近及远
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

	retrynum := 0
retry:
	isneedscale := false
	for _, i := range lst {
		// 首先看加1是否就会跳
		if wad.checkTurn(avgWin, bet, options, false, i, 1, true) {
			// 直接放弃，下一轮
			isneedscale = true

			continue
		}

		n := options.FuncGetDataNum(wad.Wins[i].Data)
		if n > 1 {
			// 再看加满是否会跳，如果加满不会跳，就直接加满
			if !wad.checkTurn(avgWin, bet, options, false, i, n, true) {
				wad.Wins[i].Weight += n

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
			wad.Wins[i].Weight++

			if wad.checkWin(avgWin, bet, options) == 0 {
				return true
			}
		}
	}

	if !isneedscale {
		goto retry
	}

	if retrynum < options.MaxFitTimes {
		wad.scale(10)
		retrynum++

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

func (wad *WinAreaData) countTotalWeight() {
	wad.TotalWeights = 0

	for _, v := range wad.Wins {
		wad.TotalWeights += v.Weight
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

func (wad *WinAreaData) Fit2(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	wad.initWeights(options)

	curawin := wad.calcAvgWin(bet)

	wo := options.cmpWin(curawin, avgWin)
	if wo == 0 {
		return true
	}

	if wo < 0 {
		return wad.scaleUp2(avgWin, bet, options)
	}

	return wad.scaleDown2(avgWin, bet, options)
}

func (wad *WinAreaData) Format(options *WinWeightFitOptions) {
	totalweight := 0.0
	for _, v := range wad.Wins {
		totalweight += float64(v.Weight)
	}

	for _, v := range wad.Wins {
		v.Weight = int(math.Floor(float64(v.Weight) * float64(options.TotalWeight) / totalweight))
	}
}

func (wad *WinAreaData) FitEnding(avgWin float64, bet int, options *WinWeightFitOptions) bool {
	// 先排序
	sort.Slice(wad.Wins, func(i, j int) bool {
		return wad.Wins[i].Win < wad.Wins[j].Win
	})

	curawin := wad.calcAvgWin(bet)

	// 因为前面拟合用于win，这里本质上是rtp的拟合
	srcWinScale := options.WinScale
	options.WinScale = options.RTPScale
	defer func() {
		options.WinScale = srcWinScale
	}()

	defer func() {
		wad.Format(options)

		for _, v := range wad.Wins {
			options.FuncSetWeight(v.Data, v.Weight)
		}
	}()

	wo := options.cmpWin(curawin, avgWin)
	if wo == 0 {
		return true
	}

	if wo < 0 {
		return wad.scaleUpEnding(avgWin, bet, options)
	}

	return wad.scaleDownEnding(avgWin, bet, options)
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

				// 如果到最后都合不上，也没办法了，就是最后一个解可能无解
				newi := wd.mergeAvgWins(si, i)
				ww.merge(si, i, newi)

				return newi, nil
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
				if si != i {
					newi := wd.mergeAvgWins(si, i)
					ww.merge(si, i, newi)

					lasti = newi
					si = i + 1

					continue
				}

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

				lasti = newi
				si = i + 1
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

	lstwd := []*WinData{}

	maxi := ww.getMaxIndex()

	for k, v := range wd.AvgWins {
		wwv, isok := ww.MapData[k]
		if isok {
			if wwv.Fit(v.AvgWin, bet, options) || k == maxi {
				target.MapData[k] = wwv

				wwv.countTotalWeight()

				for _, win := range wwv.Wins {
					nw := win.Clone()

					cw := v.Percent * float64(options.TotalWeight) * float64(win.Weight) / float64(wwv.TotalWeights)
					if cw < 0.5 {
						nw.Weight = 1

						options.FuncSetWeight(win.Data, 0)
					} else if cw < 1 {
						nw.Weight = 1

						options.FuncSetWeight(win.Data, 1)
					} else {
						nw.Weight = int(cw)

						options.FuncSetWeight(win.Data, int(cw))
					}

					lstwd = append(lstwd, nw)
				}
			}
		}
	}

	nwad := &WinAreaData{
		Wins: lstwd,
	}

	nwad.FitEnding(wd.getAllAvgWin(), bet, options)

	return target, nil
}

func (ww *WinWeight) Fit2(wd *WinningDistribution, bet int, options *WinWeightFitOptions) (*WinWeight, error) {
	err := ww.mergeWith(wd, bet, options)
	if err != nil {
		goutils.Error("WinWeight.Fit2:mergeWith",
			zap.Error(err))

		return nil, err
	}

	ww.sort()

	target := NewWinWeight()

	lstwd := []*WinData{}

	maxi := ww.getMaxIndex()

	for k, v := range wd.AvgWins {
		wwv, isok := ww.MapData[k]
		if isok {
			if wwv.Fit2(v.AvgWin, bet, options) || k == maxi {
				target.MapData[k] = wwv

				wwv.countTotalWeight()

				for _, win := range wwv.Wins {
					nw := win.Clone()

					cw := v.Percent * float64(options.TotalWeight) * float64(win.Weight) / float64(wwv.TotalWeights)
					if cw < 0.5 {
						nw.Weight = 1

						options.FuncSetWeight(win.Data, 0)
					} else if cw < 1 {
						nw.Weight = 1

						options.FuncSetWeight(win.Data, 1)
					} else {
						nw.Weight = int(cw)

						options.FuncSetWeight(win.Data, int(cw))
					}

					lstwd = append(lstwd, nw)
				}
			}
		}
	}

	nwad := &WinAreaData{
		Wins: lstwd,
	}

	nwad.FitEnding(wd.getAllAvgWin(), bet, options)

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
