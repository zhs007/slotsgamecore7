package mathtoolset

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type SymbolMapping struct {
	MapSymbols map[SymbolType]SymbolType
}

func (sm *SymbolMapping) IsSameKeys(vm *sgc7game.ValMapping2) bool {
	if len(sm.MapSymbols) == len(vm.MapVals) {
		for k := range sm.MapSymbols {
			_, isok := vm.MapVals[int(k)]
			if !isok {
				return false
			}
		}

		return true
	}

	return false
}

func (sm *SymbolMapping) Has(dst SymbolType, cur SymbolType) bool {
	s, isok := sm.MapSymbols[cur]
	if isok {
		return s == dst
	}

	return false
}

func (sm *SymbolMapping) HasTarget(cur SymbolType) bool {
	for _, v := range sm.MapSymbols {
		if v == cur {
			return true
		}
	}

	return false
}

func NewSymbolMapping() *SymbolMapping {
	return &SymbolMapping{
		MapSymbols: make(map[SymbolType]SymbolType),
	}
}
