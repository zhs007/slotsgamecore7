package mathtoolset

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type ways3ReelData struct {
	totalNum    int
	noWinNum    int
	mapMul      map[float64]int
	waitGenData bool
}

func (ways3rd *ways3ReelData) clear() {
	ways3rd.noWinNum = -1
	ways3rd.mapMul = make(map[float64]int)
	ways3rd.waitGenData = true
}

func (ways3rd *ways3ReelData) eachMul(oneach func(mul float64, times int) float64) float64 {
	i64 := float64(0)

	for k, v := range ways3rd.mapMul {
		i64 += oneach(k, v)
	}

	return i64
}

func (ways3rd *ways3ReelData) addMul(mul float64) {
	ways3rd.mapMul[mul] = ways3rd.mapMul[mul] + 1
}

func (ways3rd *ways3ReelData) getNoWinData(rd *sgc7game.ReelsData,
	symbol SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, overlaySyms *sgc7game.ValMapping2, x int, height int) int {
	if ways3rd.noWinNum >= 0 {
		return ways3rd.noWinNum
	}

	num := 0

	for y := 0; y < len(rd.Reels[x]); y++ {
		curmul := 0

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

			if cs == int(symbol) {
				curmul++
			} else if HasSymbol(wilds, SymbolType(cs)) {
				curmul++
			} else if symbolMapping != nil {
				ts, isok := symbolMapping.MapSymbols[SymbolType(cs)]
				if isok && ts == symbol {
					curmul++
				}
			}
		}

		if curmul <= 0 {
			num++
		}
	}

	ways3rd.noWinNum = num

	return ways3rd.noWinNum
}

func (ways3rd *ways3ReelData) calcWins(rd *sgc7game.ReelsData,
	symbol SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, overlaySyms *sgc7game.ValMapping2, x int, height int) {

	if !ways3rd.waitGenData {
		return
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
			ways3rd.addMul(curmul)
		}
	}

	ways3rd.waitGenData = false
}

func newWays3ReelDataList(num int) []*ways3ReelData {
	lst := []*ways3ReelData{}

	for i := 0; i < num; i++ {
		lst = append(lst, &ways3ReelData{
			mapMul:      make(map[float64]int),
			noWinNum:    -1,
			waitGenData: true,
		})
	}

	return lst
}

func clearWay3ReelDataList(lst []*ways3ReelData) {
	for _, v := range lst {
		v.clear()
	}
}

// calcWaysWinsInReels2 -
func calcWaysWinsInReels3(ways3data []*ways3ReelData, rd *sgc7game.ReelsData, symbol SymbolType, wilds []SymbolType,
	symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, overlaySyms *sgc7game.ValMapping2, x int, num int, height int) float64 {

	// curwins := float64(0)

	ways3data[x].calcWins(rd, symbol, wilds, symbolMapping, symMul, overlaySyms, x, height)

	if x < num-1 {
		return ways3data[x].eachMul(func(mul float64, times int) float64 {
			curwin := calcWaysWinsInReels3(ways3data, rd, symbol, wilds, symbolMapping, symMul, overlaySyms, x+1, num, height)

			return mul * float64(times) * curwin
		})

		// for y := 0; y < len(rd.Reels[x]); y++ {
		// 	curmul := float64(0)

		// 	for ty := 0; ty < height; ty++ {
		// 		off := y + ty
		// 		if off >= len(rd.Reels[x]) {
		// 			off -= len(rd.Reels[x])
		// 		}

		// 		cs := rd.Reels[x][off]

		// 		if overlaySyms != nil {
		// 			ocs := getSymbolWithPos(overlaySyms, x, ty)
		// 			if ocs >= 0 {
		// 				cs = ocs
		// 			}
		// 		}

		// 		csm := float64(1.0)
		// 		if symMul != nil {
		// 			cm, isok := symMul.MapVals[cs]
		// 			if isok {
		// 				csm = cm.Float64()
		// 			}
		// 		}

		// 		if cs == int(symbol) {
		// 			curmul += csm
		// 		} else if HasSymbol(wilds, SymbolType(cs)) {
		// 			curmul += csm
		// 		} else if symbolMapping != nil {
		// 			ts, isok := symbolMapping.MapSymbols[SymbolType(cs)]
		// 			if isok && ts == symbol {
		// 				curmul += csm
		// 			}
		// 		}
		// 	}

		// 	if curmul > 0 {
		// 		curwin := calcWaysWinsInReels3(ways3data, rd, symbol, wilds, symbolMapping, symMul, overlaySyms, x+1, num, height)

		// 		curwins += curmul * curwin
		// 	}
		// }
	}

	lastnum := float64(1)
	if num < len(rd.Reels) {
		lastnum = float64(ways3data[num].getNoWinData(rd, symbol, wilds, symbolMapping, overlaySyms, num, height))

		for i := num + 1; i < len(rd.Reels); i++ {
			lastnum *= float64(len(rd.Reels[i]))
		}
	}

	return ways3data[x].eachMul(func(mul float64, times int) float64 {
		// curwin := calcWaysWinsInReels3(ways3data, rd, symbol, wilds, symbolMapping, symMul, overlaySyms, x+1, num, height)

		return mul * float64(times) * lastnum
	})

	// for y := 0; y < len(rd.Reels[x]); y++ {
	// 	curmul := float64(0)

	// 	for ty := 0; ty < height; ty++ {
	// 		off := y + ty
	// 		if off >= len(rd.Reels[x]) {
	// 			off -= len(rd.Reels[x])
	// 		}

	// 		cs := rd.Reels[x][off]

	// 		if overlaySyms != nil {
	// 			ocs := getSymbolWithPos(overlaySyms, x, ty)
	// 			if ocs >= 0 {
	// 				cs = ocs
	// 			}
	// 		}

	// 		csm := float64(1.0)
	// 		if symMul != nil {
	// 			cm, isok := symMul.MapVals[cs]
	// 			if isok {
	// 				csm = cm.Float64()
	// 			}
	// 		}

	// 		if cs == int(symbol) {
	// 			curmul += csm
	// 		} else if HasSymbol(wilds, SymbolType(cs)) {
	// 			curmul += csm
	// 		} else if symbolMapping != nil {
	// 			ts, isok := symbolMapping.MapSymbols[SymbolType(cs)]
	// 			if isok && ts == symbol {
	// 				curmul += csm
	// 			}
	// 		}
	// 	}

	// 	if curmul > 0 {
	// 		curwins += curmul * lastnum
	// 	}
	// }
	// }

	// return curwins
}

func AnalyzeReelsWaysEx3(paytables *sgc7game.PayTables, rd *sgc7game.ReelsData,
	symbols []SymbolType, wilds []SymbolType, symbolMapping *SymbolMapping, symMul *sgc7game.ValMapping2, overlaySyms *sgc7game.ValMapping2,
	height int, bet int, mul int) (*SymbolsWinsStats, error) {

	ways3data := newWays3ReelDataList(len(rd.Reels))

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for i, rs := range rd.Reels {
		ssws.TotalBet *= int64(len(rs))

		ways3data[i].totalNum = len(rs)
	}

	ssws.TotalBet *= int64(mul)
	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)
		clearWay3ReelDataList(ways3data)

		if symbolMapping.HasTarget(s) {
			arrPay, isok := paytables.MapPay[int(s)]
			if isok {
				for i := 0; i < len(arrPay); i++ {
					if arrPay[i] > 0 {
						cw := calcWaysWinsInReels3(ways3data, rd, s, wilds, symbolMapping, symMul, overlaySyms, 0, i+1, height)

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
						cw := calcWaysWinsInReels3(ways3data, rd, s, wilds, nil, symMul, overlaySyms, 0, i+1, height)

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
