package sgc7game

import (
	"os"

	"github.com/bytedance/sonic"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type symbolWeightReels struct {
	SetType1 int `json:"settype1"`
	SetType2 int `json:"settype2"`
	Symbol   int `json:"symbolid"`
	R1       int `json:"r1"`
	R2       int `json:"r2"`
	R3       int `json:"r3"`
	R4       int `json:"r4"`
	R5       int `json:"r5"`
}

// SymbolWeightReelData - symbol weight reel data
type SymbolWeightReelData struct {
	Symbols    []int
	Weights    []int
	MaxWeights int
}

// SymbolWeightReelsData - symbol weight reels data
type SymbolWeightReelsData struct {
	Reels []*SymbolWeightReelData
}

// SymbolWeightReelsDataSet - symbol weight reels data
type SymbolWeightReelsDataSet struct {
	Arr []*SymbolWeightReelsData
}

// SymbolWeightReels - symbol weight reels
type SymbolWeightReels struct {
	Sets  []*SymbolWeightReelsDataSet
	Width int
}

// LoadSymbolWeightReels5JSON - load json file
func LoadSymbolWeightReels5JSON(fn string) (*SymbolWeightReels, error) {
	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var ri []symbolWeightReels
	err = sonic.Unmarshal(data, &ri)
	if err != nil {
		return nil, err
	}

	if len(ri) <= 0 {
		return nil, nil
	}

	p := &SymbolWeightReels{
		Width: 5,
	}

	for _, v := range ri {
		p.insData5(v)
	}

	return p, nil
}

// insData5
func (swr *SymbolWeightReels) insData5(data symbolWeightReels) error {
	si1 := data.SetType1 - 1
	if si1 < 0 {
		return ErrInvalidSymbolWeightReelsSetType1
	}

	if si1 >= len(swr.Sets) {
		for i := len(swr.Sets); i <= si1; i++ {
			swr.Sets = append(swr.Sets, &SymbolWeightReelsDataSet{})
		}
	}

	curset := swr.Sets[si1]

	si2 := data.SetType2 - 1
	if si2 < 0 {
		return ErrInvalidSymbolWeightReelsSetType2
	}

	if si2 >= len(curset.Arr) {
		for i := len(curset.Arr); i <= si2; i++ {
			cd := &SymbolWeightReelsData{}

			for j := 0; j < 5; j++ {
				cd.Reels = append(cd.Reels, &SymbolWeightReelData{})
			}

			curset.Arr = append(curset.Arr, cd)
		}
	}

	curdata := curset.Arr[si2]

	curdata.Reels[0].Symbols = append(curdata.Reels[0].Symbols, data.Symbol)
	curdata.Reels[0].Weights = append(curdata.Reels[0].Weights, data.R1)
	curdata.Reels[0].MaxWeights += data.R1

	curdata.Reels[1].Symbols = append(curdata.Reels[1].Symbols, data.Symbol)
	curdata.Reels[1].Weights = append(curdata.Reels[1].Weights, data.R2)
	curdata.Reels[1].MaxWeights += data.R2

	curdata.Reels[2].Symbols = append(curdata.Reels[2].Symbols, data.Symbol)
	curdata.Reels[2].Weights = append(curdata.Reels[2].Weights, data.R3)
	curdata.Reels[2].MaxWeights += data.R3

	curdata.Reels[3].Symbols = append(curdata.Reels[3].Symbols, data.Symbol)
	curdata.Reels[3].Weights = append(curdata.Reels[3].Weights, data.R4)
	curdata.Reels[3].MaxWeights += data.R4

	curdata.Reels[4].Symbols = append(curdata.Reels[4].Symbols, data.Symbol)
	curdata.Reels[4].Weights = append(curdata.Reels[4].Weights, data.R5)
	curdata.Reels[4].MaxWeights += data.R5

	return nil
}

// RandomScene -
func (swr *SymbolWeightReels) RandomScene(gs *GameScene, plugin sgc7plugin.IPlugin, si1 int, si2 int, nocheck bool) error {
	if gs.Width != swr.Width {
		return ErrInvalidSymbolWeightWidthReels
	}

	if si1 < 0 || si1 >= len(swr.Sets) {
		return ErrInvalidSymbolWeightReelsSetType1
	}

	curset := swr.Sets[si1]

	if si2 < 0 || si2 >= len(curset.Arr) {
		return ErrInvalidSymbolWeightReelsSetType2
	}

	curdata := curset.Arr[si2]

	if nocheck {
		for x, arr := range gs.Arr {
			for y := range arr {
				csi, err := RandWithWeights(plugin, curdata.Reels[x].MaxWeights, curdata.Reels[x].Weights)
				if err != nil {
					return err
				}

				gs.Arr[x][y] = curdata.Reels[x].Symbols[csi]
			}
		}
	} else {
		for x, arr := range gs.Arr {
			for y, v := range arr {
				if v == -1 {
					csi, err := RandWithWeights(plugin, curdata.Reels[x].MaxWeights, curdata.Reels[x].Weights)
					if err != nil {
						return err
					}

					gs.Arr[x][y] = curdata.Reels[x].Symbols[csi]
				}
			}
		}
	}

	return nil
}
