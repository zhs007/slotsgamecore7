package mathtoolset

import sgc7game "github.com/zhs007/slotsgamecore7/game"

func buildWaysSymbolStatsKey(symbol SymbolType, numInWindow int) int {
	return int(symbol)*100 + numInWindow
}

func unpackWaysSymbolStatsKey(key int) (SymbolType, int) {
	return SymbolType(key / 100), key % 100
}

type WaysSymbolStats struct {
	Symbol      SymbolType
	NumInWindow int
	Num         int
}

func newWaysSymbolStats(s SymbolType, numInWindow int, num int) *WaysSymbolStats {
	return &WaysSymbolStats{
		Symbol:      s,
		NumInWindow: numInWindow,
		Num:         num,
	}
}

type WaysReelStats struct {
	MapSymbols     map[int]*WaysSymbolStats
	TotalSymbolNum int
}

func (wrs *WaysReelStats) GetNumWithSymbolNumInWindow(symbol SymbolType, numInWindow int) int {
	wss, isok := wrs.MapSymbols[buildWaysSymbolStatsKey(symbol, numInWindow)]
	if isok {
		return wss.Num
	}

	return 0
}

func newWaysReelStats() *WaysReelStats {
	return &WaysReelStats{
		MapSymbols: make(map[int]*WaysSymbolStats),
	}
}

func newWaysReelStatsWithReel(reel []int, height int) *WaysReelStats {
	wrs := newWaysReelStats()

	symbols := []SymbolType{}

	for _, v := range reel {
		if !HasSymbol(symbols, SymbolType(v)) {
			symbols = append(symbols, SymbolType(v))
		}
	}

	for _, s := range symbols {
		for y := range reel {
			num := CountSymbolInReel(s, reel, y, height)

			if num > 0 {
				key := buildWaysSymbolStatsKey(s, num)
				wss, isok := wrs.MapSymbols[key]
				if !isok {
					wss = newWaysSymbolStats(s, num, 1)

					wrs.MapSymbols[key] = wss
				} else {
					wss.Num++
				}
			}
		}
	}

	return wrs
}

type WaysReelsStats struct {
	Reels   []*WaysReelStats
	Symbols []SymbolType
	Height  int
}

func (wrss *WaysReelsStats) rebuildSymbols() {
	wrss.Symbols = nil

	for _, wrs := range wrss.Reels {
		for k := range wrs.MapSymbols {
			s, _ := unpackWaysSymbolStatsKey(k)

			if !HasSymbol(wrss.Symbols, s) {
				wrss.Symbols = append(wrss.Symbols, s)
			}
		}
	}
}

func BuildWaysReelsStats(rd *sgc7game.ReelsData, height int) *WaysReelsStats {
	wrss := NewWaysReelsStats(height)

	for _, r := range rd.Reels {
		wrs := newWaysReelStatsWithReel(r, height)

		wrss.Reels = append(wrss.Reels, wrs)
	}

	wrss.rebuildSymbols()

	return wrss
}

func NewWaysReelsStats(height int) *WaysReelsStats {
	return &WaysReelsStats{
		Height: height,
	}
}
