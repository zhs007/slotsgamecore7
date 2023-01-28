package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type SymbolStats struct {
	Symbol SymbolType
	Num    int
}

func NewSymbolStats(s SymbolType, num int) *SymbolStats {
	return &SymbolStats{
		Symbol: s,
		Num:    num,
	}
}

type ReelStats struct {
	MapSymbols map[SymbolType]*SymbolStats
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
		v, isok := rs.MapSymbols[SymbolType(s)]
		if !isok {
			rs.MapSymbols[SymbolType(s)] = NewSymbolStats(SymbolType(s), 1)
		} else {
			v.Num++
		}
	}

	return rs, nil
}

type ReelsStats struct {
	Reels []*ReelStats
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

	return rss, nil
}
