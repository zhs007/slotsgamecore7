package mathtoolset

type SymbolType int

func HasSymbol(symbols []SymbolType, symbol SymbolType) bool {
	for _, v := range symbols {
		if v == symbol {
			return true
		}
	}

	return false
}
