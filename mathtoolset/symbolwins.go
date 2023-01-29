package mathtoolset

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type SymbolsWinsFileMode int

const (
	SWFModeRTP     SymbolsWinsFileMode = 1
	SWFModeWinsNum SymbolsWinsFileMode = 2
	SWFModeWins    SymbolsWinsFileMode = 3
)

type SymbolWinsStats struct {
	Symbol  SymbolType
	WinsNum []int64
	Wins    []int64
}

func newSymbolWinsStats(symbol SymbolType, num int) *SymbolWinsStats {
	return &SymbolWinsStats{
		Symbol:  symbol,
		WinsNum: make([]int64, num),
		Wins:    make([]int64, num),
	}
}

type SymbolsWinsStats struct {
	MapSymbols map[SymbolType]*SymbolWinsStats
	Symbols    []SymbolType
	Num        int
	Total      int64
}

func (ssws *SymbolsWinsStats) GetSymbolWinsStats(symbol SymbolType) *SymbolWinsStats {
	sws, isok := ssws.MapSymbols[symbol]
	if isok {
		return sws
	}

	ssws.MapSymbols[symbol] = newSymbolWinsStats(symbol, ssws.Num)
	ssws.Symbols = append(ssws.Symbols, symbol)

	return ssws.MapSymbols[symbol]
}

func (ssws *SymbolsWinsStats) buildSortedSymbols() {
	sort.Slice(ssws.Symbols, func(i, j int) bool {
		return ssws.Symbols[i] < ssws.Symbols[j]
	})
}

func (ssws *SymbolsWinsStats) SaveExcel(fn string, fm SymbolsWinsFileMode) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "symbol")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "total")

	si := 2

	for i := 0; i < ssws.Num; i++ {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+si, 0), fmt.Sprintf("X%v", i+1))
	}

	y := 1

	for _, s := range ssws.Symbols {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), int(s))
		f.SetCellValue(sheet, goutils.Pos2Cell(1, y), ssws.Total)

		for i := 0; i < ssws.Num; i++ {
			sws := ssws.GetSymbolWinsStats(s)

			if fm == SWFModeRTP {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), float64(sws.Wins[i])*100.0/float64(ssws.Total))
			} else if fm == SWFModeWinsNum {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), sws.WinsNum[i])
			} else {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), sws.Wins[i])
			}
		}

		y++
	}

	return f.SaveAs(fn)
}

func newSymbolsWinsStatsWithPaytables(paytables *sgc7game.PayTables, symbols []SymbolType) *SymbolsWinsStats {
	num := 0
	for _, arr := range paytables.MapPay {
		if len(arr) > num {
			num = len(arr)
		}
	}

	ssws := NewSymbolsWinsStats(num)

	for s := range paytables.MapPay {
		if HasSymbol(symbols, SymbolType(s)) {
			ssws.GetSymbolWinsStats(SymbolType(s))
		}
	}

	return ssws
}

// lst is like [S, W, S+W, No S+W, All]
func CalcSymbolWins(rss *ReelsStats, wilds []SymbolType, symbol SymbolType, symbol2 SymbolType, lst []InReelSymbolType) (int64, error) {
	curwins := int64(1)

	for i, t := range lst {
		cn := rss.GetNum(i, symbol, symbol2, wilds, t)
		if cn < 0 {
			goutils.Error("CalcSymbolWins:GetNum",
				zap.Int("InReelSymbolType", int(t)),
				zap.Error(ErrInvalidInReelSymbolType))

			return 0, ErrInvalidInReelSymbolType
		}

		curwins *= int64(cn)
	}

	return curwins, nil
}

// calcWildWins - 这个接口只能用于处理wild赢得，symbol必须是wild
func calcWildWins(paytables *sgc7game.PayTables, rss *ReelsStats, wilds []SymbolType, symbol SymbolType, num int) int64 {
	curwins := int64(1)

	// 如果数量最大，不需要处理排除的解
	if num == len(rss.Reels) {
		for i := 0; i < num; i++ {
			cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeWild)

			if cn <= 0 {
				return 0
			}

			curwins *= int64(cn)
		}
	}

	for i := 0; i < num; i++ {
		cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeWild)

		if cn <= 0 {
			return 0
		}

		curwins *= int64(cn)
	}

	// // 如果 A x 5 > W x 4，那么在计算 W x 4 时，就需要排除第5个图标是 A 的情况
	// wp := paytables.MapPay[int(symbol)][num-1]

	return curwins
}

// calcNotWildWins - 这个接口只能用于处理非wild的赢得，symbol必须不是wild，且只处理 S 开头的情况
func calcNotWildWins(rss *ReelsStats, wilds []SymbolType, symbol SymbolType, num int) int64 {
	curwins := int64(rss.GetNum(0, symbol, -1, wilds, IRSTypeSymbol))

	for i := 1; i < num; i++ {
		cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeSymbolAndWild)

		if cn <= 0 {
			return 0
		}

		curwins *= int64(cn)
	}

	if num == len(rss.Reels) {
		return curwins
	}

	curwins *= int64(rss.GetNum(num, symbol, -1, wilds, IRSTypeNoSymbolAndNoWild))

	for i := num + 1; i < len(rss.Reels); i++ {
		cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeAll)

		if cn <= 0 {
			return 0
		}

		curwins *= int64(cn)
	}

	return curwins
}

// lst is like [S, W, S+W, All, All]
// 如果是 www 开头，且 w 作为 symbol符号，这里会需要减去比3w大的情况
func calcSymbolWinsFromList(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, ci int, num int, lst []InReelSymbolType) (int64, error) {

	if ci == num-1 {
		// 如果要计算 3a，如果是w开头，且前面至少有1个a，那么第4个只要不是a和w就好了
		if IsFirstWild(lst, num) {
			// 如果是 3w，用加法来算
			curwin := int64(0)
			ps := paytables.MapPay[int(symbol)][num-1]

			for _, s := range symbols {
				if s == symbol {
					continue
				}

				if HasSymbol(wilds, s) {
					continue
				}

				parr := paytables.MapPay[int(s)]
				cn := -1
				for j := num; j < len(parr); j++ {
					if parr[j] > ps {
						cn = j

						break
					}
				}

				// 如果4b大于3w，那么5b也一定大于3w，所以就彻底排除接下来如果是符号b的情况
				if cn == num {
					continue
				}

				lst[ci] = IRSTypeSymbol2
				for j := ci + 1; j < len(parr); j++ {
					lst[j] = IRSTypeAll
				}

				cw, err := CalcSymbolWins(rss, wilds, symbol, s, lst)
				if err != nil {
					goutils.Error("calcSymbolWinsFromList:CalcSymbolWins",
						zap.Error(err))

					return 0, err
				}

				curwin += cw

				// 如果 4b、5b 都小于 3w，那么需要加入接下来是b的情况
				if cn < 0 {
					continue
				} else {
					// 剩下就是4b小于3w，且5b大于3w，那么我们需要加入wwwbx，减去wwwbb和wwwbw
					// 如果是6个，则加入 wwwbxx，减去 wwwbbx 和 wwwbwx
					lst[ci+1] = IRSTypeSymbol2AndWild

					for j := ci + 2; j < len(parr); j++ {
						lst[j] = IRSTypeAll
					}

					cw0, err := CalcSymbolWins(rss, wilds, symbol, s, lst)
					if err != nil {
						goutils.Error("calcSymbolWinsFromList:CalcSymbolWins",
							zap.Error(err))

						return 0, err
					}

					curwin -= cw0
				}
			}

			return curwin, nil
		}

		lst[ci] = IRSTypeNoSymbolAndNoWild

		return CalcSymbolWins(rss, wilds, symbol, -1, lst)
	}

	curwins := int64(0)

	for t := IRSTypeSymbol; t < IRSTypeWild; t++ {
		lst[ci] = t

		cw, err := calcSymbolWinsFromList(paytables, rss, symbols, wilds, symbol, ci+1, num, lst)
		if err != nil {
			goutils.Error("calcSymbolWinsFromList:calcSymbolWinsFromList",
				zap.Int("ci", ci),
				zap.Int("InReelSymbolType", int(t)),
				zap.Error(err))

			return 0, err
		}

		curwins += cw
	}

	return curwins, nil
}

// calcSymbolFirstWildWins - 这个接口只能用于处理非wild赢得，symbol必须不是wild，且以wild开头
func calcSymbolFirstWildWins(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, num int) (int64, error) {

	lst := NewInReelSymbolTypeArr(num)
	lst[0] = IRSTypeWild

	return calcSymbolWinsFromList(paytables, rss, symbols, wilds, symbol, 1, num, lst)
}

// CalcSymbolWinsInReelsWithLine -
// symbol 不可能是 wild
func CalcSymbolWinsInReelsWithLine(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, num int) (int64, error) {

	curwins := calcNotWildWins(rss, wilds, symbol, num)

	cw, err := calcSymbolFirstWildWins(paytables, rss, symbols, wilds, symbol, num)
	if err != nil {
		goutils.Error("CalcSymbolWinsInReelsWithLine:calcSymbolFirstWildWins",
			zap.Error(err))

		return 0, err
	}

	curwins += cw

	return curwins, nil

	// // if symbol is wild
	// if HasSymbol(wilds, symbol) {
	// 	for i := 0; i < num; i++ {
	// 		cn := rss.GetSymbolNum(i, symbol, wilds)

	// 		if cn <= 0 {
	// 			return 0
	// 		}

	// 		curwins *= int64(cn)
	// 	}

	// 	if num == len(rss.Reels) {
	// 		return curwins
	// 	}

	// 	curwins *= int64(rss.GetReelLengthNoSymbol(num, symbol, wilds))

	// 	for i := num + 1; i < len(rss.Reels); i++ {
	// 		curwins *= int64(rss.GetReelLength(i))
	// 	}

	// 	return curwins
	// }

	// // if wildnum <= 0 || wildnum >= len(rss.Reels) {

	// // }

	// for i := 0; i < num; i++ {
	// 	ss := rss.Reels[i].GetSymbolStats(symbol)

	// 	wildnum := 0
	// 	for _, w := range wilds {
	// 		if w == symbol {
	// 			continue
	// 		}

	// 		ws := rss.Reels[i].GetSymbolStats(w)
	// 		if ws.Num > 0 {
	// 			wildnum += ws.Num
	// 		}
	// 	}

	// 	if ss.Num <= 0 && wildnum <= 0 {
	// 		return 0
	// 	}

	// 	curwins *= int64(ss.Num + wildnum)
	// }

	// if num == len(rss.Reels) {
	// 	return curwins
	// }

	// for i := num; i < len(rss.Reels); i++ {
	// 	curwins *= int64(rss.Reels[i].TotalSymbolNum)
	// }

	// return curwins
}

func analyzeWildNum(paytables *sgc7game.PayTables, symbol SymbolType, num int, wild SymbolType) int {
	sp := paytables.MapPay[int(symbol)][num-1]
	warr := paytables.MapPay[int(wild)]

	for i := 0; i < len(warr); i++ {
		if warr[i] >= sp {
			return i + 1
		}
	}

	return 0
}

func AnalyzeReelsWithLine(paytables *sgc7game.PayTables, reels *sgc7game.ReelsData,
	symbols []SymbolType, wilds []SymbolType, betMul int, lineNum int) (*SymbolsWinsStats, error) {

	rss, err := BuildReelsStats(reels)
	if err != nil {
		goutils.Error("AnalyzeReelsWithLine:BuildReelsStats",
			zap.Error(err))

		return nil, err
	}

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.Total = 1
	for _, arr := range reels.Reels {
		ssws.Total *= int64(len(arr))
	}

	ssws.Total *= int64(betMul)

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcSymbolWinsInReelsWithLine(paytables, rss, symbols, wilds, s, i+1)
					if err != nil {
						goutils.Error("AnalyzeReelsWithLine:CalcSymbolWinsInReelsWithLine",
							zap.Error(err))

						return nil, err
					}

					sws.WinsNum[i] = cw
					sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(lineNum)
				}
			}
		}
	}

	ssws.buildSortedSymbols()

	return ssws, nil
}

func NewSymbolsWinsStats(num int) *SymbolsWinsStats {
	return &SymbolsWinsStats{
		MapSymbols: make(map[SymbolType]*SymbolWinsStats),
		Num:        num,
	}
}
