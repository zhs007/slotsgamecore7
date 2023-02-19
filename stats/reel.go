package stats

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
)

type SymbolStats struct {
	Symbol       mathtoolset.SymbolType
	TriggerTimes int64
}

func (ss *SymbolStats) CalcHitRate(totalTimes int64) float64 {
	return float64(ss.TriggerTimes) / float64(totalTimes)
}

func NewSymbolStats(s mathtoolset.SymbolType) *SymbolStats {
	return &SymbolStats{
		Symbol: s,
	}
}

type Reel struct {
	Index      int
	MapSymbols map[mathtoolset.SymbolType]*SymbolStats
	TotalTimes int64
}

func (reel *Reel) CalcHitRate(s mathtoolset.SymbolType) float64 {
	return reel.MapSymbols[s].CalcHitRate(reel.TotalTimes)
}

func (reel *Reel) OnScene(scene *sgc7game.GameScene) {
	reel.TotalTimes++

	for y := 0; y < scene.Height; y++ {
		s := scene.Arr[reel.Index][y]

		reel.MapSymbols[mathtoolset.SymbolType(s)].TriggerTimes++
	}
}

func (reel *Reel) OnSymbols(lst []mathtoolset.SymbolType) {
	reel.TotalTimes++

	for _, s := range lst {
		reel.MapSymbols[mathtoolset.SymbolType(s)].TriggerTimes++
	}
}

func (reel *Reel) GenSymbols(symbols []mathtoolset.SymbolType) []mathtoolset.SymbolType {
	for k := range reel.MapSymbols {
		if !mathtoolset.HasSymbol(symbols, k) {
			symbols = append(symbols, k)
		}
	}

	return symbols
}

func NewReel(i int, lst []mathtoolset.SymbolType) *Reel {
	r := &Reel{
		Index:      i,
		MapSymbols: make(map[mathtoolset.SymbolType]*SymbolStats),
	}

	for _, v := range lst {
		r.MapSymbols[v] = NewSymbolStats(v)
	}

	return r
}
