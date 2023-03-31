package mathtoolset

import sgc7game "github.com/zhs007/slotsgamecore7/game"

func GetSymbols(arr []string, paytables *sgc7game.PayTables) []SymbolType {
	symbols := []SymbolType{}

	for _, v := range arr {
		s, isok := paytables.MapSymbols[v]
		if isok {
			symbols = append(symbols, SymbolType(s))
		}
	}

	return symbols
}

// CountSymbolInReel - count symbol number in reel，[stop, stop + height)
func CountSymbolInReel(symbol SymbolType, reel []int, stop int, height int) int {
	if stop < 0 {
		for {
			stop += len(reel)

			if stop >= 0 {
				break
			}
		}
	}

	if stop >= len(reel) {
		for {
			stop -= len(reel)

			if stop < len(reel) {
				break
			}
		}
	}

	num := 0

	for i := 0; i < height; i++ {
		if reel[stop] == int(symbol) {
			num++
		}

		stop++
		if stop >= len(reel) {
			stop -= len(reel)
		}
	}

	return num
}
