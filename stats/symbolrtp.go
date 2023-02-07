package stats

import (
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

func NewSymbolRTP(s mathtoolset.SymbolType, maxSymbolNum int) *SymbolRTP {
	return &SymbolRTP{
		Symbol: s,
		Wins:   make([]int64, maxSymbolNum),
	}
}

type SymbolsRTP struct {
	MapSymbols map[mathtoolset.SymbolType]*SymbolRTP
	TotalBets  int64
}

func (ssrtp *SymbolsRTP) OnBet(bet int64) {
	ssrtp.TotalBets += bet
}

func (ssrtp *SymbolsRTP) OnWin(win *sgc7game.Result) {
	ssrtp.MapSymbols[mathtoolset.SymbolType(win.Symbol)].OnWin(win)
}

func NewSymbolsRTP(maxSymbolNum int, lst []mathtoolset.SymbolType) *SymbolsRTP {
	ssrtp := &SymbolsRTP{
		MapSymbols: make(map[mathtoolset.SymbolType]*SymbolRTP),
	}

	for _, v := range lst {
		ssrtp.MapSymbols[v] = NewSymbolRTP(v, maxSymbolNum)
	}

	return ssrtp
}
