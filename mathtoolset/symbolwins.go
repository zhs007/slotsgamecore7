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
	RTPMode     SymbolsWinsFileMode = 1
	WinsNumMode SymbolsWinsFileMode = 2
	WinsMode    SymbolsWinsFileMode = 3
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

			if fm == RTPMode {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), float64(sws.Wins[i])*100.0/float64(ssws.Total))
			} else if fm == WinsNumMode {
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

func CalcSymbolWinsInReelsWithLine(rss *ReelsStats, symbol SymbolType, num int) int64 {
	curwins := int64(1)

	for i := 0; i < num; i++ {
		ss := rss.Reels[i].GetSymbolStats(symbol)
		if ss.Num <= 0 {
			return 0
		}

		curwins *= int64(ss.Num)
	}

	for i := num; i < len(rss.Reels); i++ {
		curwins *= int64(rss.Reels[i].TotalSymbolNum)
	}

	return curwins
}

func AnalyzeReelsWithLine(paytables *sgc7game.PayTables, reels *sgc7game.ReelsData,
	symbols []SymbolType, betMul int, lineNum int) (*SymbolsWinsStats, error) {

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
					sws.WinsNum[i] = CalcSymbolWinsInReelsWithLine(rss, s, i+1)
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
