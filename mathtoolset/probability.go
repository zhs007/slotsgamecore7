package mathtoolset

import sgc7game "github.com/zhs007/slotsgamecore7/game"

// countScatterNumTimes -
func countScatterNumTimes(rss *ReelsStats, symbol SymbolType, num int, ci int, height int) int64 {
	if ci == len(rss.Reels)-1 {
		if num > 1 {
			return 0
		}

		if num > 0 {
			cn := rss.GetScatterNum(ci, symbol, IRSTypeSymbol, height)

			return int64(cn)
		}

		cnn := rss.GetScatterNum(ci, symbol, IRSTypeNoSymbol, height)

		return int64(cnn)
	}

	tn := int64(0)

	if num > 0 {
		cn := rss.GetScatterNum(ci, symbol, IRSTypeSymbol, height)
		if cn > 0 {
			nn := countScatterNumTimes(rss, symbol, num-1, ci+1, height)

			tn += int64(cn) * nn
		}
	}

	cnn := rss.GetScatterNum(ci, symbol, IRSTypeNoSymbol, height)
	if cnn > 0 {
		nn := countScatterNumTimes(rss, symbol, num, ci+1, height)

		tn += int64(cnn) * nn
	}

	return tn
}

func CalcScatterProbability(rss *ReelsStats, symbol SymbolType, num int, height int) float64 {

	totalnum := int64(1)

	for _, rs := range rss.Reels {
		totalnum *= int64(rs.TotalSymbolNum)
	}

	symbolnum := countScatterNumTimes(rss, symbol, num, 0, height)

	return float64(symbolnum) / float64(totalnum)
}

func CalcProbWithWeights(vw *sgc7game.ValWeights2, probs []float64) float64 {
	ret := float64(0)

	for i, v := range probs {
		ret += v * float64(vw.Weights[i]) / float64(vw.MaxWeight)
	}

	return ret
}

type scReelData struct {
	mapNum map[int]int
}

func (scrd *scReelData) add(num int) {
	scrd.mapNum[num] = scrd.mapNum[num] + 1
}

func (scrd *scReelData) build(rd *sgc7game.ReelsData, symbol SymbolType, symbolMapping *SymbolMapping, overlaySyms *sgc7game.ValMapping2, x int, height int) {
	if len(scrd.mapNum) > 0 {
		return
	}

	for y := range rd.Reels[x] {

		cn := 0
		for ty := 0; ty < height; ty++ {
			off := y + ty

			if off >= len(rd.Reels[x]) {
				off -= len(rd.Reels[x])
			}

			s := rd.Reels[x][off]

			if overlaySyms != nil {
				ocs := getSymbolWithPos(overlaySyms, x, ty)
				if ocs >= 0 {
					s = ocs
				}
			}

			if symbolMapping != nil {
				ts, isok := symbolMapping.MapSymbols[SymbolType(s)]
				if isok {
					s = int(ts)
				}
			}

			if s == int(symbol) {
				cn++
			}
		}

		scrd.add(cn)
	}
}

func newScReelData(num int) []*scReelData {
	lst := []*scReelData{}

	for i := 0; i < num; i++ {
		lst = append(lst, &scReelData{
			mapNum: make(map[int]int),
		})
	}

	return lst
}

// countScatterNumTimes -
func countScatterNumTimesWithReels(scrds []*scReelData, rd *sgc7game.ReelsData, symbol SymbolType, symbolMapping *SymbolMapping, overlaySyms *sgc7game.ValMapping2, num int, x int, height int) int64 {
	scrds[x].build(rd, symbol, symbolMapping, overlaySyms, x, height)

	if x == len(rd.Reels)-1 {
		if num > 1 {
			return 0
		}

		if num > 0 {
			return int64(scrds[x].mapNum[1])
		}

		return int64(scrds[x].mapNum[0])
	}

	tn := int64(0)

	if num > 0 {
		cn := int64(scrds[x].mapNum[1])
		if cn > 0 {
			nn := countScatterNumTimesWithReels(scrds, rd, symbol, symbolMapping, overlaySyms, num-1, x+1, height)

			tn += int64(cn) * nn
		}
	}

	cnn := int64(scrds[x].mapNum[0])
	if cnn > 0 {
		nn := countScatterNumTimesWithReels(scrds, rd, symbol, symbolMapping, overlaySyms, num, x+1, height)

		tn += int64(cnn) * nn
	}

	return tn
}

func CalcScatterProbabilitWithReels(rd *sgc7game.ReelsData, symbol SymbolType, symbolMapping *SymbolMapping, overlaySyms *sgc7game.ValMapping2, num int, height int) float64 {
	scrds := newScReelData(len(rd.Reels))

	totalnum := int64(1)

	for _, r := range rd.Reels {
		totalnum *= int64(len(r))
	}

	symbolnum := countScatterNumTimesWithReels(scrds, rd, symbol, symbolMapping, overlaySyms, num, 0, height)

	return float64(symbolnum) / float64(totalnum)
}
