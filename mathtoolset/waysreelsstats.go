package mathtoolset

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func buildWaysSymbolStatsKey(symbol SymbolType, numInWindow int) int {
	return int(symbol)*100 + numInWindow
}

func unpackWaysSymbolStatsKey(key int) (SymbolType, int) {
	return SymbolType(key / 100), key % 100
}

type FuncOnWaysSymbolStats func(*WaysSymbolStats)

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

func (wrs *WaysReelStats) GetSymbolKeys(symbol SymbolType) []int {
	arr := []int{}

	for k := range wrs.MapSymbols {
		if k >= int(symbol)*100 && k < (int(symbol)+1)*100 {
			arr = append(arr, k)
		}
	}

	return arr
}

func (wrs *WaysReelStats) EachSymbol(symbol SymbolType, onEach FuncOnWaysSymbolStats) {
	for k, v := range wrs.MapSymbols {
		if k >= int(symbol)*100 && k < (int(symbol)+1)*100 {
			onEach(v)
		}
	}
}

func (wrs *WaysReelStats) GetNumWithSymbolNumInWindow(symbol SymbolType, numInWindow int) int {
	wss, isok := wrs.MapSymbols[buildWaysSymbolStatsKey(symbol, numInWindow)]
	if isok {
		return wss.Num
	}

	return 0
}

func (wrs *WaysReelStats) GetSymbolStats(symbol SymbolType, numInWindow int) *WaysSymbolStats {
	k := buildWaysSymbolStatsKey(symbol, numInWindow)
	v, isok := wrs.MapSymbols[k]
	if isok {
		return v
	}

	wrs.MapSymbols[k] = newWaysSymbolStats(symbol, numInWindow, 0)

	return wrs.MapSymbols[k]
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

	wrs.TotalSymbolNum = len(reel)

	return wrs
}

func newWaysReelStatsWithReelEx(reel []int, height int, symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping) *WaysReelStats {
	wrs := newWaysReelStats()

	for _, s := range symbols {
		for y := range reel {
			num := CountSymbolInReelEx(s, reel, y, height, wilds, symbolMapping)

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

	wrs.TotalSymbolNum = len(reel)

	return wrs
}

type WaysReelsStats struct {
	Reels   []*WaysReelStats
	Symbols []SymbolType
	Keys    []int
	Height  int
}

// func (wrss *WaysReelsStats) GetWaysNum(reelindex int, symbol SymbolType, numInWindow int, wilds []SymbolType, wildNumInWindow int, irstype InReelSymbolType, height int) int {
// 	ss := wrss.Reels[reelindex].GetSymbolStats(symbol, numInWindow)

// 	wildnum := 0
// 	for _, w := range wilds {
// 		if w == symbol {
// 			continue
// 		}

// 		ws := wrss.Reels[reelindex].GetSymbolStats(w, wildNumInWindow)
// 		if ws.Num > 0 {
// 			wildnum += ws.Num
// 		}
// 	}

// 	if irstype == IRSTypeSymbol {
// 		return (wildnum + ss.Num) * height
// 	}

// 	if irstype == IRSTypeNoSymbol {
// 		return wrss.Reels[reelindex].TotalSymbolNum - (wildnum+ss.Num)*height
// 	}

// 	return -1
// }

func (wrss *WaysReelsStats) GetNonWaysNum(reelindex int, symbol SymbolType) int {
	num := 0

	wrss.Reels[reelindex].EachSymbol(symbol, func(wss *WaysSymbolStats) {
		num += wss.Num
	})

	return wrss.Reels[reelindex].TotalSymbolNum - num
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

	wrss.Keys = nil

	for _, wrs := range wrss.Reels {
		for k := range wrs.MapSymbols {
			if goutils.IndexOfIntSlice(wrss.Keys, k, 0) < 0 {
				wrss.Keys = append(wrss.Keys, k)
			}
		}
	}
}

func (wrss *WaysReelsStats) SaveExcel(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "symbol")
	for i := range wrss.Reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("reel%v", i+1))
	}

	y := 1

	for _, k := range wrss.Keys {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), k)

		for i, reel := range wrss.Reels {
			statsSymbol, isok := reel.MapSymbols[k]
			if isok {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), statsSymbol.Num)
			} else {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), 0)
			}
		}

		y++
	}

	for i, rs := range wrss.Reels {
		f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), rs.TotalSymbolNum)
	}

	return f.SaveAs(fn)
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

// BuildWaysReelsStatsEx - 只计算symbols里的symbol，且把wild直接计算进去
func BuildWaysReelsStatsEx(rd *sgc7game.ReelsData, height int, symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping) *WaysReelsStats {
	wrss := NewWaysReelsStats(height)

	for _, r := range rd.Reels {
		wrs := newWaysReelStatsWithReelEx(r, height, symbols, wilds, symbolMapping)

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
