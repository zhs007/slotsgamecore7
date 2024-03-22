package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// CalcWaysWinsInReelsEx -
func CalcWaysWinsInReelsEx(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, num int, height int) (int64, error) {
	curwins := int64(1)

	for i := 0; i < num; i++ {
		curwins *= int64(rss.GetWaysNumEx(i, symbol, wilds, symbolMapping, IRSTypeSymbol, height))
	}

	if num < len(rss.Reels) {
		curwins *= int64(rss.GetWaysNumEx(num, symbol, wilds, symbolMapping, IRSTypeNoSymbol, height))

		for i := num + 1; i < len(rss.Reels); i++ {
			curwins *= int64(rss.Reels[i].TotalSymbolNum)
		}
	}

	return curwins, nil
}

// CalcWaysWinsInReels -
func CalcWaysWinsInReels(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, wilds []SymbolType, num int, height int) (int64, error) {
	curwins := int64(1)

	for i := 0; i < num; i++ {
		curwins *= int64(rss.GetWaysNum(i, symbol, wilds, IRSTypeSymbol, height))
	}

	if num < len(rss.Reels) {
		curwins *= int64(rss.GetWaysNum(num, symbol, wilds, IRSTypeNoSymbol, height))

		for i := num + 1; i < len(rss.Reels); i++ {
			curwins *= int64(rss.Reels[i].TotalSymbolNum)
		}
	}

	return curwins, nil
}

// AnalyzeReelsWaysEx - totalbet = reels length x mul, wins = symbol wins x mul
func AnalyzeReelsWaysEx(paytables *sgc7game.PayTables, rss *ReelsStats,
	symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, height int, bet int, mul int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range rss.Reels {
		ssws.TotalBet *= int64(rs.TotalSymbolNum)
	}

	ssws.TotalBet *= int64(mul)
	ssws.TotalWins = 0

	for _, s := range symbols {
		if symbolMapping.HasTarget(s) {
			sws := ssws.GetSymbolWinsStats(s)

			arrPay, isok := paytables.MapPay[int(s)]
			if isok {
				for i := 0; i < len(arrPay); i++ {
					if arrPay[i] > 0 {
						cw, err := CalcWaysWinsInReelsEx(paytables, rss, s, wilds, symbolMapping, i+1, height)
						if err != nil {
							goutils.Error("AnalyzeReelsWaysEx:CalcWaysWinsInReelsEx",
								goutils.Err(err))

							return nil, err
						}

						sws.WinsNum[i] = cw
						sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(bet)

						ssws.TotalWins += sws.Wins[i]
					}
				}
			}
		} else {
			sws := ssws.GetSymbolWinsStats(s)

			arrPay, isok := paytables.MapPay[int(s)]
			if isok {
				for i := 0; i < len(arrPay); i++ {
					if arrPay[i] > 0 {
						cw, err := CalcWaysWinsInReels(paytables, rss, s, wilds, i+1, height)
						if err != nil {
							goutils.Error("AnalyzeReelsWaysEx:CalcWaysWinsInReels",
								goutils.Err(err))

							return nil, err
						}

						sws.WinsNum[i] = cw
						sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(bet)

						ssws.TotalWins += sws.Wins[i]
					}
				}
			}
		}
	}

	ssws.onBuildEnd()

	return ssws, nil
}
