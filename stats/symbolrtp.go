package stats

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
)

type SymbolRTP struct {
	Symbol mathtoolset.SymbolType
	Wins   []int64
}

func (srtp *SymbolRTP) OnWin(win *sgc7game.Result) {
	if win.Symbol == int(srtp.Symbol) {
		srtp.Wins[win.SymbolNums-1] += int64(win.CashWin)
	}
}

func (srtp *SymbolRTP) CalcRTP(totalBets int64, num int) float64 {
	return float64(srtp.Wins[num-1]) / float64(totalBets)
}

func NewSymbolRTP(s mathtoolset.SymbolType, maxSymbolWinNum int) *SymbolRTP {
	return &SymbolRTP{
		Symbol: s,
		Wins:   make([]int64, maxSymbolWinNum),
	}
}

type SymbolsRTP struct {
	MapSymbols      map[mathtoolset.SymbolType]*SymbolRTP
	TotalBets       int64
	MaxSymbolWinNum int
}

func (ssrtp *SymbolsRTP) GenSymbols() []mathtoolset.SymbolType {
	symbols := []mathtoolset.SymbolType{}

	for k := range ssrtp.MapSymbols {
		symbols = append(symbols, k)
	}

	return symbols
}

func (ssrtp *SymbolsRTP) OnBet(bet int64) {
	ssrtp.TotalBets += bet
}

func (ssrtp *SymbolsRTP) OnWin(win *sgc7game.Result) {
	ssrtp.MapSymbols[mathtoolset.SymbolType(win.Symbol)].OnWin(win)
}

func (ssrtp *SymbolsRTP) SaveSheet(f *excelize.File, sheet string) error {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "symbol")

	symbols := ssrtp.GenSymbols()

	for i := 0; i < ssrtp.MaxSymbolWinNum; i++ {
		f.SetCellValue(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("X%v", i+1))
	}

	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})

	y := 1
	for _, s := range symbols {
		f.SetCellValue(sheet, goutils.Pos2Cell(0, y), s)

		for i := 0; i < ssrtp.MaxSymbolWinNum; i++ {
			f.SetCellValue(sheet, goutils.Pos2Cell(i+1, y), ssrtp.MapSymbols[s].CalcRTP(ssrtp.TotalBets, i+1))
		}

		y++
	}

	return nil
}

func NewSymbolsRTP(maxSymbolWinNum int, lst []mathtoolset.SymbolType) *SymbolsRTP {
	ssrtp := &SymbolsRTP{
		MapSymbols:      make(map[mathtoolset.SymbolType]*SymbolRTP),
		MaxSymbolWinNum: maxSymbolWinNum,
	}

	for _, v := range lst {
		ssrtp.MapSymbols[v] = NewSymbolRTP(v, maxSymbolWinNum)
	}

	return ssrtp
}
