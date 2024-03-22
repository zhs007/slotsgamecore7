package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// calcScatterWinsInReels -
func calcScatterWinsInReels(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, num int, lst []int, ci int, height int) (int64, error) {
	if len(lst) == num {
		curwin := int64(1)

		for i := 0; i < len(rss.Reels); i++ {
			if goutils.IndexOfIntSlice(lst, i, 0) >= 0 {
				curwin *= int64(rss.GetScatterNum(i, symbol, IRSTypeSymbol, height))
			} else {
				curwin *= int64(rss.GetScatterNum(i, symbol, IRSTypeNoSymbol, height))
			}
		}

		return curwin, nil
	}

	if ci == len(rss.Reels) {
		return 0, nil
	}

	totalwin := int64(0)

	for t := 0; t <= 1; t++ {
		if t == 0 {
			lst = append(lst, ci)
		}

		cw, err := calcScatterWinsInReels(paytables, rss, symbol, num, lst, ci+1, height)
		if err != nil {
			goutils.Error("calcScatterWinsInReels:calcScatterWinsInReels",
				goutils.Err(err))

			return 0, err
		}

		if t == 0 {
			lst = lst[0 : len(lst)-1]
		}

		totalwin += cw
	}

	return totalwin, nil
}

// CalcScatterWinsInReels -
func CalcScatterWinsInReels(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, num int, height int) (int64, error) {
	curwins := int64(0)

	for t := 0; t <= 1; t++ {
		lst := []int{}

		if t == 0 {
			lst = append(lst, 0)
		}

		cw, err := calcScatterWinsInReels(paytables, rss, symbol, num, lst, 1, height)
		if err != nil {
			goutils.Error("CalcScatterWinsInReels:calcScatterWinsInReels",
				goutils.Err(err))

			return 0, err
		}

		curwins += cw
	}

	return curwins, nil
}

func AnalyzeReelsScatter(paytables *sgc7game.PayTables, reels *sgc7game.ReelsData,
	symbols []SymbolType, mapSymbols *SymbolsMapping, height int) (*SymbolsWinsStats, error) {

	rss, err := BuildReelsStats(reels, mapSymbols)
	if err != nil {
		goutils.Error("AnalyzeReelsWithScatter:BuildReelsStats",
			goutils.Err(err))

		return nil, err
	}

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, arr := range reels.Reels {
		ssws.TotalBet *= int64(len(arr))
	}

	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcScatterWinsInReels(paytables, rss, s, i+1, height)
					if err != nil {
						goutils.Error("AnalyzeReelsScatter:CalcScatterWinsInReels",
							goutils.Err(err))

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

func AnalyzeReelsScatterEx(paytables *sgc7game.PayTables, rss *ReelsStats,
	symbols []SymbolType, height int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range rss.Reels {
		ssws.TotalBet *= int64(rs.TotalSymbolNum)
	}

	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcScatterWinsInReels(paytables, rss, s, i+1, height)
					if err != nil {
						goutils.Error("AnalyzeReelsScatter:CalcScatterWinsInReels",
							goutils.Err(err))

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
