package mathtoolset

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func pos2int(x, y int) int {
	return y*10 + x
}

func int2pos(k int) (int, int) {
	return k % 10, k / 10
}

// calcNonWaysWinsInReels2 -
func calcNonWaysWinsInReels2(rd *sgc7game.ReelsData,
	symbol SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, x int, height int) int64 {

	num := int64(0)

	for y := 0; y < len(rd.Reels[x]); y++ {
		curmul := 0

		for ty := 0; ty < height; ty++ {
			off := y + ty
			if off >= len(rd.Reels[x]) {
				off -= len(rd.Reels[x])
			}

			if rd.Reels[x][off] == int(symbol) {
				curmul++
			} else if HasSymbol(wilds, SymbolType(rd.Reels[x][off])) {
				curmul++
			} else if symbolMapping != nil {
				ts, isok := symbolMapping.MapSymbols[SymbolType(rd.Reels[x][off])]
				if isok && ts == symbol {
					curmul++
				}
			}
		}

		if curmul <= 0 {
			num++
		}
	}

	return num
}

func getSymbolWithPos(overlaySyms *sgc7game.ValMapping2, x, y int) int {
	for k, v := range overlaySyms.MapVals {
		cx, cy := int2pos(k)

		if cx == x && cy == y {
			return v.Int()
		}
	}

	return -1
}

// calcWaysWinsInReels2 -
func calcWaysWinsInReels2(rd *sgc7game.ReelsData,
	symbol SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, overlaySyms *sgc7game.ValMapping2, x int, num int, height int) float64 {

	curwins := float64(0)

	if x < num-1 {
		for y := 0; y < len(rd.Reels[x]); y++ {
			curmul := float64(0)

			for ty := 0; ty < height; ty++ {
				off := y + ty
				if off >= len(rd.Reels[x]) {
					off -= len(rd.Reels[x])
				}

				cs := rd.Reels[x][off]

				if overlaySyms != nil {
					ocs := getSymbolWithPos(overlaySyms, x, off)
					if ocs >= 0 {
						cs = ocs
					}
				}

				csm := float64(1.0)
				if symMul != nil {
					cm, isok := symMul.MapVals[cs]
					if isok {
						csm = cm.Float64()
					}
				}

				if cs == int(symbol) {
					curmul += csm
				} else if HasSymbol(wilds, SymbolType(cs)) {
					curmul += csm
				} else if symbolMapping != nil {
					ts, isok := symbolMapping.MapSymbols[SymbolType(cs)]
					if isok && ts == symbol {
						curmul += csm
					}
				}
			}

			if curmul > 0 {
				curwin := calcWaysWinsInReels2(rd, symbol, wilds, symbolMapping, symMul, overlaySyms, x+1, num, height)

				curwins += curmul * curwin
			}
		}
	} else {
		lastnum := float64(1)
		if num < len(rd.Reels) {
			lastnum = float64(calcNonWaysWinsInReels2(rd, symbol, wilds, symbolMapping, num, height))

			for i := num + 1; i < len(rd.Reels); i++ {
				lastnum *= float64(len(rd.Reels[i]))
			}
		}

		for y := 0; y < len(rd.Reels[x]); y++ {
			curmul := float64(0)

			for ty := 0; ty < height; ty++ {
				off := y + ty
				if off >= len(rd.Reels[x]) {
					off -= len(rd.Reels[x])
				}

				cs := rd.Reels[x][off]

				if overlaySyms != nil {
					ocs := getSymbolWithPos(overlaySyms, x, ty)
					if ocs >= 0 {
						cs = ocs
					}
				}

				csm := float64(1.0)
				if symMul != nil {
					cm, isok := symMul.MapVals[cs]
					if isok {
						csm = cm.Float64()
					}
				}

				if cs == int(symbol) {
					curmul += csm
				} else if HasSymbol(wilds, SymbolType(cs)) {
					curmul += csm
				} else if symbolMapping != nil {
					ts, isok := symbolMapping.MapSymbols[SymbolType(cs)]
					if isok && ts == symbol {
						curmul += csm
					}
				}
			}

			if curmul > 0 {
				curwins += curmul * lastnum
			}
		}
	}

	return curwins
}

func AnalyzeReelsWaysEx2(paytables *sgc7game.PayTables, rd *sgc7game.ReelsData,
	symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, overlaySyms *sgc7game.ValMapping2,
	height int, bet int, mul int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range rd.Reels {
		ssws.TotalBet *= int64(len(rs))
	}

	ssws.TotalBet *= int64(mul)
	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		if symbolMapping.HasTarget(s) {
			arrPay, isok := paytables.MapPay[int(s)]
			if isok {
				for i := 0; i < len(arrPay); i++ {
					if arrPay[i] > 0 {
						cw := calcWaysWinsInReels2(rd, s, wilds, symbolMapping, symMul, overlaySyms, 0, i+1, height)

						sws.WinsNum[i] = int64(cw)
						sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(bet)

						ssws.TotalWins += sws.Wins[i]
					}
				}
			}
		} else {
			arrPay, isok := paytables.MapPay[int(s)]
			if isok {
				for i := 0; i < len(arrPay); i++ {
					if arrPay[i] > 0 {
						cw := calcWaysWinsInReels2(rd, s, wilds, nil, symMul, overlaySyms, 0, i+1, height)

						sws.WinsNum[i] = int64(cw)
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
