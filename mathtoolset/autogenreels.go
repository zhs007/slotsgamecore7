package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

func AutoGenReels(w, h int, paytables *sgc7game.PayTables, syms []SymbolType, wilds []SymbolType,
	totalBet int, lineNum int) (*sgc7game.ReelsData, error) {

	rss, err := BuildBasicReelsStatsEx(w, syms)
	if err != nil {
		goutils.Error("AutoGenReels:BuildBasicReelsStatsEx",
			zap.Error(err))

		return nil, err
	}

	ssws, err := AnalyzeReelsWithLineEx(paytables, rss, syms, wilds, totalBet, lineNum)
	if err != nil {
		goutils.Error("AutoGenReels:AnalyzeReelsWithLineEx",
			zap.Error(err))

		return nil, err
	}

	goutils.Info("AutoGenReels",
		zap.Float64("rtp", ssws.CountRTP()))

	return nil, nil
}
