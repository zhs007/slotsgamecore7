package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type funcOnEachCWWIRSM func(float64)
type funcCWWIRSMCalcCurMulti func(int) float64
type funcCWWIRSMCalcCurMultiEx func(int, SymbolType) float64

func deepEachCWWIRSM(i int, max int, curmul float64, symbols []int, func0 funcCWWIRSMCalcCurMulti, func1 funcCWWIRSMCalcCurMultiEx, oneach funcOnEachCWWIRSM) {
	if i < max-1 {
		deepEachCWWIRSM(i+1, max, curmul*func0(i), symbols, func0, func1, oneach)

		for _, v := range symbols {
			deepEachCWWIRSM(i+1, max, curmul*func1(i, SymbolType(v)), symbols, func0, func1, oneach)
		}
	} else {
		oneach(curmul * func0(i))

		for _, v := range symbols {
			oneach(curmul * func1(i, SymbolType(v)))
		}
	}
}

// calcWaysWinsInReelsSymbolMulti - symMul is map[int]float
func calcWaysWinsInReelsSymbolMulti(paytables *sgc7game.PayTables, rss *ReelsStats, symbol SymbolType, wilds []SymbolType,
	symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, num int, height int) (int64, error) {

	curwins := int64(0)
	syms := symMul.Keys()

	lastnum := int64(1)
	if num < len(rss.Reels) {
		lastnum = int64(rss.GetWaysNumEx(num, symbol, wilds, symbolMapping, IRSTypeSymbol, height))

		for i := num + 1; i < len(rss.Reels); i++ {
			lastnum *= int64(rss.Reels[i].TotalSymbolNum)
		}
	}

	deepEachCWWIRSM(0, num, 1, syms, func(i int) float64 {
		return float64(rss.GetWaysNum(i, symbol, wilds, IRSTypeSymbol, height))
	}, func(i int, s SymbolType) float64 {
		return float64(rss.GetWaysNum(i, s, nil, IRSTypeSymbol, height)) * symMul.MapVals[int(s)].Float64()
	}, func(mul float64) {
		curwins += int64(mul * float64(lastnum))
	})

	// for i := 0; i < num; i++ {
	// 	curwins *= int64(rss.GetWaysNumEx(i, symbol, wilds, symbolMapping, IRSTypeSymbol, height))
	// }

	// if num < len(rss.Reels) {
	// 	curwins *= int64(rss.GetWaysNumEx(num, symbol, wilds, symbolMapping, IRSTypeNoSymbol, height))

	// 	for i := num + 1; i < len(rss.Reels); i++ {
	// 		curwins *= int64(rss.Reels[i].TotalSymbolNum)
	// 	}
	// }

	return curwins, nil
}

// AnalyzeReelsWaysSymbolMulti - totalbet = reels length x mul, wins = symbol wins x mul
func AnalyzeReelsWaysSymbolMulti(paytables *sgc7game.PayTables, rss *ReelsStats,
	symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, height int, bet int, mul int) (*SymbolsWinsStats, error) {

	if symMul == nil || symMul.IsEmpty() {
		return AnalyzeReelsWaysEx(paytables, rss, symbols, wilds, symbolMapping, height, bet, mul)
	}

	//!! 现在只处理图标映射和图标倍数是一种图标的情况，别的分支没处理正确
	//!! 这个接口只是一个近似计算，不考虑同时出现的情况，所以轮子最好能回避这点，最好w、s、m都不会同时出现，否则误差将非常大
	if symbolMapping != nil && symMul != nil && symbolMapping.IsSameKeys(symMul) {
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
							cw, err := calcWaysWinsInReelsSymbolMulti(paytables, rss, s, wilds, symbolMapping, symMul, i+1, height)
							if err != nil {
								goutils.Error("AnalyzeReelsWaysSymbolMulti:calcWaysWinsInReelsSymbolMulti",
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
								goutils.Error("AnalyzeReelsWaysSymbolMulti:CalcWaysWinsInReels",
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

	goutils.Error("AnalyzeReelsWaysSymbolMulti",
		goutils.Err(ErrUnimplementedCode))

	return nil, ErrUnimplementedCode
}
