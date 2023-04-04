package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type cwwirNode struct {
	Symbol      SymbolType
	NumInWindow int
	Num         int
}

type cwwirReel struct {
	SymbolNodes [][]*cwwirNode
}

func buildCalcWaysWinsInReels2Data(wrss *WaysReelsStats, symbol SymbolType, num int) *cwwirReel {
	cwwirr := &cwwirReel{}

	for i := 0; i < num; i++ {
		nodes := []*cwwirNode{}

		wrss.Reels[i].EachSymbol(symbol, func(wss *WaysSymbolStats) {
			if wss.Num > 0 {
				nodes = append(nodes, &cwwirNode{
					Symbol:      symbol,
					NumInWindow: wss.NumInWindow,
					Num:         wss.Num,
				})
			}
		})

		cwwirr.SymbolNodes = append(cwwirr.SymbolNodes, nodes)
	}

	return cwwirr
}

type funcOnEachCWWIRReel func(int)

func deepEachCWWIRReel(cwwirr *cwwirReel, i int, curmul int, oneach funcOnEachCWWIRReel) {
	if i < len(cwwirr.SymbolNodes)-1 {
		for _, v := range cwwirr.SymbolNodes[i] {
			deepEachCWWIRReel(cwwirr, i+1, curmul*v.NumInWindow*v.Num, oneach)
		}
	} else {
		for _, v := range cwwirr.SymbolNodes[i] {
			oneach(curmul * v.NumInWindow * v.Num)
		}
	}
}

func eachCWWIRReel(cwwirr *cwwirReel, oneach funcOnEachCWWIRReel) {
	deepEachCWWIRReel(cwwirr, 0, 1, oneach)
}

// CalcWaysWinsInReels2 -
func CalcWaysWinsInReels2(wrss *WaysReelsStats, symbol SymbolType, num int, height int) (int64, error) {
	curwins := int64(0)

	cwwirr := buildCalcWaysWinsInReels2Data(wrss, symbol, num)

	lastnum := int64(1)
	if num < len(wrss.Reels) {
		lastnum = int64(wrss.GetNonWaysNum(num, symbol))

		for i := num + 1; i < len(wrss.Reels); i++ {
			lastnum *= int64(wrss.Reels[i].TotalSymbolNum)
		}
	}

	eachCWWIRReel(cwwirr, func(mul int) {
		curwins += (int64(mul) * lastnum)
	})

	return curwins, nil
}

func AnalyzeReelsWaysEx2(paytables *sgc7game.PayTables, wrss *WaysReelsStats,
	symbols []SymbolType, height int, bet int, mul int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range wrss.Reels {
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
					cw, err := CalcWaysWinsInReels2(wrss, s, i+1, height)
					if err != nil {
						goutils.Error("AnalyzeReelsWaysEx2:CalcWaysWinsInReels",
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
