package asciigame

import (
	"github.com/fatih/color"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type FuncOnSelectColor func() bool
type FuncGetSymbolString func(int) string

func SelectColor(onselect FuncOnSelectColor, c1 *color.Color, c2 *color.Color) *color.Color {
	if onselect() {
		return c1
	}

	return c2
}

type SymbolColorMap struct {
	MapSymbols        map[int]*color.Color
	PayTables         *sgc7game.PayTables
	OnGetSymbolString FuncGetSymbolString
}

func (mapSymbolColor *SymbolColorMap) AddSymbolColor(s int, c *color.Color) {
	mapSymbolColor.MapSymbols[s] = c
}

func (mapSymbolColor *SymbolColorMap) GetSymbolString(s int) string {
	c, isok := mapSymbolColor.MapSymbols[s]
	if isok {
		return FormatColorString(mapSymbolColor.OnGetSymbolString(s), c)
	}

	return mapSymbolColor.OnGetSymbolString(s)
}

func NewSymbolColorMap(paytables *sgc7game.PayTables) *SymbolColorMap {
	scm := &SymbolColorMap{
		MapSymbols: make(map[int]*color.Color),
		PayTables:  paytables,
	}

	scm.OnGetSymbolString = func(s int) string {
		return scm.PayTables.GetStringFromInt(s)
	}

	return scm
}
