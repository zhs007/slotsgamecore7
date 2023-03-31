package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

// CalcWaysWinsInReels -
func CalcWaysWinsInReels(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, wilds []SymbolType, num int, height int) (int64, error) {
	curwins := int64(1)

	for i := 0; i < num; i++ {
		curwins *= int64(rss.GetWaysNum(i, symbol, wilds, IRSTypeSymbol, height))
	}

	for i := num; i < len(rss.Reels); i++ {
		curwins *= int64(rss.GetWaysNum(i, symbol, wilds, IRSTypeNoSymbol, height))
	}

	return curwins, nil
}

func AnalyzeReelsWaysEx(paytables *sgc7game.PayTables, rss *ReelsStats,
	symbols []SymbolType, wilds []SymbolType, height int, bet int, mul int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range rss.Reels {
		ssws.TotalBet *= int64(rs.TotalSymbolNum)
	}

	ssws.TotalBet *= int64(mul)
	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcWaysWinsInReels(paytables, rss, s, wilds, i+1, height)
					if err != nil {
						goutils.Error("AnalyzeReelsWaysEx:CalcWaysWinsInReels",
							zap.Error(err))

						return nil, err
					}

					sws.WinsNum[i] = cw
					sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i]

					ssws.TotalWins += int64(arrPay[i]) * sws.WinsNum[i]
				}
			}
		}
	}

	ssws.onBuildEnd()

	return ssws, nil
}
