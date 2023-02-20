package stats

import (
	"fmt"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
)

type Reels struct {
	Reels []*Reel
}

func (reels *Reels) Merge(src *Reels) {
	for i, v := range src.Reels {
		reels.Reels[i].Merge(v)
	}
}

func (reels *Reels) OnReelSymbols(mapSyms map[int][]mathtoolset.SymbolType) {
	for i, v := range reels.Reels {
		lst, isok := mapSyms[i]
		if isok {
			v.OnSymbols(lst)
		} else {
			v.OnSymbols(nil)
		}
	}
}

func (reels *Reels) OnScene(scene *sgc7game.GameScene) {
	for _, v := range reels.Reels {
		v.OnScene(scene)
	}
}

func (reels *Reels) GenSymbols() []mathtoolset.SymbolType {
	symbols := []mathtoolset.SymbolType{}

	for _, v := range reels.Reels {
		symbols = v.GenSymbols(symbols)
	}

	return symbols
}

func (reels *Reels) SaveSheet(f *excelize.File, sheet string) error {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "symbol")

	symbols := reels.GenSymbols()

	for i := range reels.Reels {
		f.SetCellValue(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("reel%v", i+1))
	}

	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})

	y := 1
	for _, s := range symbols {
		f.SetCellValue(sheet, goutils.Pos2Cell(0, y), s)

		for i := range reels.Reels {
			f.SetCellValue(sheet, goutils.Pos2Cell(i+1, y), reels.Reels[i].CalcHitRate(s))
		}

		y++
	}

	return nil
}

func NewReels(width int, lst []mathtoolset.SymbolType) *Reels {
	reels := &Reels{}

	for x := 0; x < width; x++ {
		r := NewReel(x, lst)

		reels.Reels = append(reels.Reels, r)
	}

	return reels
}
