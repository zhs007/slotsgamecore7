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
