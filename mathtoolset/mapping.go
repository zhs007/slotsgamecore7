package mathtoolset

type SymbolMapping struct {
	MapSymbols map[SymbolType]SymbolType
}

func (sm *SymbolMapping) Has(dst SymbolType, cur SymbolType) bool {
	s, isok := sm.MapSymbols[cur]
	if isok {
		return s == dst
	}

	return false
}

func NewSymbolMapping() *SymbolMapping {
	return &SymbolMapping{
		MapSymbols: make(map[SymbolType]SymbolType),
	}
}
