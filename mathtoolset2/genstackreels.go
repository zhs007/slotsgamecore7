package mathtoolset2

import (
	"io"
	"log/slog"
	"math/rand"

	"github.com/zhs007/goutils"
)

type stackSymbolData struct {
	symbol string
	num    int
}

func (ssd *stackSymbolData) clone() *stackSymbolData {
	return &stackSymbolData{
		symbol: ssd.symbol,
		num:    ssd.num,
	}
}

func rmstack(stack []int, num int) []int {
	ns := []int{}

	for _, v := range stack {
		if v != num {
			ns = append(ns, v)
		}
	}

	return ns
}

func genStackSymbol(symbol string, num int, stack []int) ([]*stackSymbolData, bool) {
	newstack := []int{}
	for _, v := range stack {
		if v < num {
			newstack = append(newstack, v)
		} else if v == num {
			return []*stackSymbolData{
				{
					symbol: symbol,
					num:    v,
				},
			}, true
		}
	}

	// symbol数量太少
	if len(newstack) == 0 {
		return nil, false
	}

	if len(newstack) == 1 {
		// 已经失败了，不能切分完全
		if num%newstack[0] != 0 {
			return nil, false
		}
	}

	ssds := []*stackSymbolData{}

	cr := rand.Int() % len(newstack)
	ssds = append(ssds, &stackSymbolData{
		symbol: symbol,
		num:    newstack[cr],
	})

	nssds, isok := genStackSymbol(symbol, num-newstack[cr], newstack)
	if !isok {
		if len(newstack) > 1 {
			// 如果失败了，那么表示这一次的选择彻底错误，那么换一个选择
			cr++
			if cr >= len(newstack) {
				cr = 0
			}
			// newstack = rmstack(newstack, newstack[cr])

			return genStackSymbol(symbol, num, newstack)
		}

		return nil, false
	}

	ssds = append(ssds, nssds...)

	return ssds, true
}

func isUnable(ssds []*stackSymbolData) bool {
	// 这个用来简单判断一下是否可以分，只判断了剩下2个symbol和剩下1个symbol的情况
	mapSymbols := make(map[string]int)

	// for _, v := range excludeSymbols {
	// 	mapSymbols[v] = 1
	// }

	for _, v := range ssds {
		_, isok := mapSymbols[v.symbol]
		if !isok {
			mapSymbols[v.symbol] = 1
		} else {
			mapSymbols[v.symbol]++
		}
	}

	if len(mapSymbols) == 1 {
		// 如果只剩下1个symbol，那么有2组就只能失败
		for _, v := range mapSymbols {
			if v > 1 {
				return true
			}
		}
	} else if len(mapSymbols) == 2 {
		// 如果剩下2个symbol，那么如果2组数量不相同，也只能失败
		v0 := -1
		for _, v := range mapSymbols {
			if v0 < 0 {
				v0 = v
			} else if AbsInt(v0, v) > 1 {
				return true
			}
		}
	}

	return false
}

func chooseMore(ssds []*stackSymbolData) int {
	mapSymbols := make(map[string]int)

	for _, v := range ssds {
		_, isok := mapSymbols[v.symbol]
		if !isok {
			mapSymbols[v.symbol] = 1
		} else {
			mapSymbols[v.symbol]++
		}
	}

	maxv := -1
	maxs := ""

	for k, v := range mapSymbols {
		if v > maxv {
			maxs = k
			maxv = v
		}
	}

	for i, v := range ssds {
		if v.symbol == maxs {
			return i
		}
	}

	return -1
}

func genReelSSD(ssds []*stackSymbolData, excludeSymbols []string, firstSymbol string) ([]*stackSymbolData, bool) {
	if isUnable(ssds) {
		return nil, false
	}

	assds := []*stackSymbolData{}
	nssds := []*stackSymbolData{}

	for _, v := range ssds {
		if goutils.IndexOfStringSlice(excludeSymbols, v.symbol, 0) < 0 {
			// if pres != v.symbol {
			assds = append(assds, v)
		} else {
			nssds = append(nssds, v)
		}
	}

	if len(assds) <= 0 {
		return nil, false
	}

	cr := rand.Int() % len(assds)
	ssd := assds[cr].clone()

	nassds := []*stackSymbolData{}
	nassds = append(nassds, assds[:cr]...)
	nassds = append(nassds, assds[cr+1:]...)
	nassds = append(nassds, nssds...)

	nextExcludeSym := []string{ssd.symbol}

	if firstSymbol == "" {
		firstSymbol = ssd.symbol
	}

	if len(nassds) == 0 {
		// nsyms := []string{}
		// for i := 0; i < ssd.num; i++ {
		// 	nsyms = append(nsyms, ssd.symbol)
		// }

		return []*stackSymbolData{ssd}, true
	} else if len(nassds) == 1 {
		nextExcludeSym = append(nextExcludeSym, firstSymbol)
	}

	target, isok := genReelSSD(nassds, nextExcludeSym, firstSymbol)
	if !isok {
		// 如果失败了，应该选择最多的那个
		ci := chooseMore(assds)
		if ci < 0 {
			return nil, false
		}

		if ssd.symbol == assds[ci].symbol {
			return nil, false
		}

		ssd = assds[ci].clone()

		nassds1 := []*stackSymbolData{}
		nassds1 = append(nassds1, assds[:ci]...)
		nassds1 = append(nassds1, assds[ci+1:]...)
		nassds1 = append(nassds1, nssds...)

		// nassds = append(assds[:ci], assds[ci+1:]...)
		// nextssds = append(nssds, nassds...)

		nextExcludeSym = []string{ssd.symbol}
		if len(nassds1) == 1 {
			nextExcludeSym = append(nextExcludeSym, firstSymbol)
		}
		// nextssds = append(nextssds, ssd)

		// ssd = ssd1

		target, isok = genReelSSD(nassds1, nextExcludeSym, firstSymbol)
		if !isok {
			return nil, false
		}
	}

	return append([]*stackSymbolData{ssd}, target...), true
	// nsyms := []string{}
	// for i := 0; i < ssd.num; i++ {
	// 	nsyms = append(nsyms, ssd.symbol)
	// }

	// return append(nsyms, syms...), true
}

func checkSR(rs2 *ReelStats2, ssds []*stackSymbolData) bool {
	tn0 := 0
	for _, v := range rs2.MapSymbols {
		tn0 += v
	}

	tn1 := 0
	for _, v := range ssds {
		tn1 += v.num
	}

	return tn0 == tn1
}

func checkSR2(rs2 *ReelStats2, arr []string) bool {
	tn0 := 0
	for _, v := range rs2.MapSymbols {
		tn0 += v
	}

	return tn0 == len(arr)
}

func genStackReel(rs2 *ReelStats2, stack []int, excludeSymbol []string) ([]string, error) {
	ssds := []*stackSymbolData{}

	for _, s := range excludeSymbol {
		for i := 0; i < rs2.MapSymbols[s]; i++ {
			ssds = append(ssds, &stackSymbolData{
				symbol: s,
				num:    1,
			})
		}
	}

	for s, n := range rs2.MapSymbols {
		if n <= 0 {
			continue
		}

		if goutils.IndexOfStringSlice(excludeSymbol, s, 0) < 0 {
			nssds, isok := genStackSymbol(s, n, stack)
			if !isok {
				goutils.Error("genStackReel:genStackSymbol",
					goutils.Err(ErrGenStackReel))

				return nil, ErrGenStackReel
			}

			ssds = append(ssds, nssds...)
		}
	}

	isok := checkSR(rs2, ssds)
	if !isok {
		goutils.Error("genStackReel:checkSR",
			goutils.Err(ErrGenStackReel))

		return nil, ErrGenStackReel
	}

	target, isok := genReelSSD(ssds, nil, "")
	if isok {
		isok := checkSR(rs2, target)
		if !isok {
			goutils.Error("genStackReel:checkSR:2",
				goutils.Err(ErrGenStackReel))

			return nil, ErrGenStackReel
		}

		reel := []string{}
		for _, v := range target {
			for i := 0; i < v.num; i++ {
				reel = append(reel, v.symbol)
			}
		}

		if !checkSR2(rs2, reel) {
			goutils.Error("genStackReel:checkSR2",
				goutils.Err(ErrGenStackReel))

			return nil, ErrGenStackReel
		}

		return reel, nil
	}

	goutils.Error("genStackReel:genReelSSD",
		goutils.Err(ErrGenStackReel))

	return nil, ErrGenStackReel
}

func GenStackReels(reader io.Reader, stack []int, excludeSymbol []string) ([][]string, error) {
	rss2, err := LoadReelsStats2(reader)
	if err != nil {
		goutils.Error("GenStackReels:LoadReelsStats2",
			goutils.Err(err))

		return nil, err
	}

	rd := [][]string{}

	for i, r := range rss2.Reels {
		reel, err := genStackReel(r, stack, excludeSymbol)
		if err != nil {
			goutils.Error("GenStackReels:genStackReel",
				slog.Int("reelIndex", i),
				goutils.Err(err))

			return nil, err
		}

		rd = append(rd, reel)
	}

	return rd, nil
}
