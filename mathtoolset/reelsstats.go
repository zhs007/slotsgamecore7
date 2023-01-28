package mathtoolset

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type SymbolStats struct {
	Symbol SymbolType
	Num    int
}

func newSymbolStats(s SymbolType, num int) *SymbolStats {
	return &SymbolStats{
		Symbol: s,
		Num:    num,
	}
}

type ReelStats struct {
	MapSymbols     map[SymbolType]*SymbolStats
	TotalSymbolNum int
}

func (rs *ReelStats) GetSymbolStats(symbol SymbolType) *SymbolStats {
	v, isok := rs.MapSymbols[symbol]
	if isok {
		return v
	}

	rs.MapSymbols[symbol] = newSymbolStats(symbol, 0)

	return rs.MapSymbols[symbol]
}

func NewReelStats() *ReelStats {
	return &ReelStats{
		MapSymbols: make(map[SymbolType]*SymbolStats),
	}
}

func BuildReelStats(reel []int) (*ReelStats, error) {
	if len(reel) == 0 {
		goutils.Error("ReelStats.AnalyzeReel",
			zap.Error(ErrInvalidReel))

		return nil, ErrInvalidReel
	}

	rs := NewReelStats()

	for _, s := range reel {
		ss := rs.GetSymbolStats(SymbolType(s))
		ss.Num++
	}

	rs.TotalSymbolNum = len(reel)

	return rs, nil
}

type ReelsStats struct {
	Reels   []*ReelStats
	Symbols []SymbolType
}

func (rss *ReelsStats) HasSymbols(symbol SymbolType) bool {
	for _, v := range rss.Symbols {
		if v == symbol {
			return true
		}
	}

	return false
}

func (rss *ReelsStats) buildSortedSymbols() {
	rss.Symbols = nil

	for _, r := range rss.Reels {
		for s := range r.MapSymbols {
			if !rss.HasSymbols(s) {
				rss.Symbols = append(rss.Symbols, s)
			}
		}
	}

	sort.Slice(rss.Symbols, func(i, j int) bool {
		return rss.Symbols[i] < rss.Symbols[j]
	})
}

func (rss *ReelsStats) SaveExcel(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "symbol")
	for i := range rss.Reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("reel%v", i+1))
	}

	y := 1

	for _, s := range rss.Symbols {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), int(s))

		for i, reel := range rss.Reels {
			statsSymbol, isok := reel.MapSymbols[s]
			if isok {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), statsSymbol.Num)
			} else {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), 0)
			}
		}

		y++
	}

	for i, rs := range rss.Reels {
		f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), rs.TotalSymbolNum)
	}

	return f.SaveAs(fn)
}

func BuildReelsStats(reels *sgc7game.ReelsData) (*ReelsStats, error) {
	rss := &ReelsStats{
		Reels: make([]*ReelStats, len(reels.Reels)),
	}

	for i, r := range reels.Reels {
		rs, err := BuildReelStats(r)
		if err != nil {
			goutils.Error("BuildReelsStats:BuildReelStats",
				zap.Error(err))

			return nil, err
		}

		rss.Reels[i] = rs
	}

	rss.buildSortedSymbols()

	return rss, nil
}
