package mathtoolset2

import (
	"math/rand"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type stackSymbolData struct {
	symbol string
	num    int
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

	if len(newstack) == 0 {
		return nil, false
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
			newstack = rmstack(newstack, newstack[cr])

			return genStackSymbol(symbol, num, newstack)
		}

		return nil, false
	}

	ssds = append(ssds, nssds...)

	return ssds, true
}

func isUnable(ssds []*stackSymbolData, pres string) bool {
	mapSymbols := make(map[string]int)

	if pres != "" {
		mapSymbols[pres] = 1
	}

	for _, v := range ssds {
		_, isok := mapSymbols[v.symbol]
		if !isok {
			mapSymbols[v.symbol] = 1
		} else {
			mapSymbols[v.symbol]++
		}
	}

	if len(mapSymbols) == 1 {
		for _, v := range mapSymbols {
			if v > 1 {
				return true
			}
		}
	} else if len(mapSymbols) == 2 {
		v0 := -1
		for _, v := range mapSymbols {
			if v0 < 0 {
				v0 = v
			} else if v0 != v {
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

const maxretryssd = 9999

func genReelSSD(ssds []*stackSymbolData, pres string) ([]string, bool) {
	assds := []*stackSymbolData{}
	nssds := []*stackSymbolData{}

	for _, v := range ssds {
		if pres != v.symbol {
			assds = append(assds, v)
		} else {
			nssds = append(nssds, v)
		}
	}

	if len(assds) <= 0 {
		return nil, false
	}

	if isUnable(assds, pres) {
		return nil, false
	}

	cr := rand.Int() % len(assds)
	ssd := assds[cr]

	nassds := append(assds[:cr], assds[cr+1:]...)
	nextssds := append(nssds, nassds...)

	if len(nextssds) == 0 {
		nsyms := []string{}
		for i := 0; i < ssd.num; i++ {
			nsyms = append(nsyms, ssd.symbol)
		}

		return nsyms, true
	}

	syms, isok := genReelSSD(nextssds, ssd.symbol)
	if !isok {
		ci := chooseMore(assds)
		ssd = assds[ci]

		nassds = append(assds[:ci], assds[ci+1:]...)
		nextssds = append(nssds, nassds...)

		syms, isok = genReelSSD(nextssds, ssd.symbol)
		if !isok {
			return nil, false
		}
	}

	nsyms := []string{}
	for i := 0; i < ssd.num; i++ {
		nsyms = append(nsyms, ssd.symbol)
	}

	return append(nsyms, syms...), true
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
				goutils.Error("GenStackReels:genStackReel",
					zap.Error(ErrGenStackReel))

				return nil, ErrGenStackReel
			}

			ssds = append(ssds, nssds...)
		}
	}

	reel, isok := genReelSSD(ssds, "")
	if isok {
		return reel, nil
	}

	goutils.Error("GenStackReels:genReelSSD",
		zap.Error(ErrGenStackReel))

	return nil, ErrGenStackReel
}

func GenStackReels(fn string, stack []int, excludeSymbol []string) ([][]string, error) {
	rss2, err := LoadReelsStats2(fn)
	if err != nil {
		goutils.Error("GenStackReels:LoadReelsStats2",
			zap.Error(err))

		return nil, err
	}

	rd := [][]string{}

	for i, r := range rss2.Reels {
		reel, err := genStackReel(r, stack, excludeSymbol)
		if err != nil {
			goutils.Error("GenStackReels:genStackReel",
				zap.Int("reelIndex", i),
				zap.Error(err))

			return nil, err
		}

		rd = append(rd, reel)
	}

	return rd, nil
}
