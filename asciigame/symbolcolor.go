package asciigame

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type SymbolColorMap struct {
	MapSymbols map[int]Color
	PayTables  *sgc7game.PayTables
}

func (mapSymbolColor *SymbolColorMap) AddSymbolColor(s int, c Color) {
	mapSymbolColor.MapSymbols[s] = c
}

func (mapSymbolColor *SymbolColorMap) GetSymbolString(s int) string {
	c, isok := mapSymbolColor.MapSymbols[s]
	if isok {
		return FormatColorString(mapSymbolColor.PayTables.GetStringFromInt(s), c)
	}

	return mapSymbolColor.PayTables.GetStringFromInt(s)
}

func NewSymbolColorMap(paytables *sgc7game.PayTables) *SymbolColorMap {
	return &SymbolColorMap{
		MapSymbols: make(map[int]Color),
		PayTables:  paytables,
	}
}
