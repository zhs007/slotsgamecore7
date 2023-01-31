package mathtoolset

type SymbolsMapping struct {
	MapSymbols map[SymbolType]SymbolType
}

func (mapSymbols *SymbolsMapping) Add(src SymbolType, dest SymbolType) {
	mapSymbols.MapSymbols[src] = dest
}

func (mapSymbols *SymbolsMapping) Has(s SymbolType) bool {
	_, isok := mapSymbols.MapSymbols[s]

	return isok
}

func NewSymbolsMapping() *SymbolsMapping {
	return &SymbolsMapping{
		MapSymbols: make(map[SymbolType]SymbolType),
	}
}
