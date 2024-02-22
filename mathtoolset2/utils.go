package mathtoolset2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

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

func NewExcelFile(reels [][]string) *excelize.File {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "line")
	for i := range reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("R%v", i+1))
	}

	maxj := 0

	for i, reel := range reels {
		if maxj < len(reel) {
			maxj = len(reel)
		}

		for j, v := range reel {
			f.SetCellStr(sheet, goutils.Pos2Cell(i+1, j+1), v)
		}
	}

	for i := 0; i < maxj; i++ {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, i+1), i)
	}

	return f
}

func SaveReels(fn string, reels [][]string) error {
	f := NewExcelFile(reels)

	return f.SaveAs(fn)
}
