package sgc7game

import (
	"context"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// FuncCountSymbolExIsSymbol -
type FuncCountSymbolExIsSymbol func(cursymbol int, x, y int) bool

// GameScene - game scene
type GameScene struct {
	Arr    [][]int `json:"arr"`
	Width  int     `json:"-"`
	Height int     `json:"-"`
}

// NewGameScene - new a GameScene
func NewGameScene(width int, height int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.Init(width, height)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

// NewGameSceneWithArr2 - new a GameScene
func NewGameSceneWithArr2(arr [][]int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.InitWithArr2(arr)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

// NewGameSceneWithArr - new a GameScene
func NewGameSceneWithArr(w, h int, arr []int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.InitWithArr(w, h, arr)
	if err != nil {
		return nil, err
	}

	return gs, nil
}

// NewGameSceneWithReels - new a GameScene
func NewGameSceneWithReels(reels *ReelsData, w, h int, arr []int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.Init(w, h)
	if err != nil {
		return nil, err
	}

	gs.Fill(reels, arr)

	return gs, nil
}

// Init - init scene
func (gs *GameScene) Init(w int, h int) error {
	gs.Arr = nil
	for x := 0; x < w; x++ {
		gs.Arr = append(gs.Arr, []int{})

		for y := 0; y < h; y++ {
			gs.Arr[x] = append(gs.Arr[x], -1)
		}
	}

	gs.Width = w
	gs.Height = h

	return nil
}

// InitWithArr2 - init scene
func (gs *GameScene) InitWithArr2(arr [][]int) error {
	gs.Arr = nil
	gs.Width = len(arr)
	gs.Height = len(arr[0])

	for _, l := range arr {
		if len(l) != gs.Height {
			return ErrInvalidArray
		}

		gs.Arr = append(gs.Arr, l)
	}

	return nil
}

// InitWithArr - init scene
func (gs *GameScene) InitWithArr(w int, h int, arr []int) error {
	if len(arr) != w*h {
		return ErrInvalidArray
	}

	gs.Width = w
	gs.Height = h
	gs.Arr = nil

	for x := 0; x < w; x++ {
		gs.Arr = append(gs.Arr, arr[x*h:(x+1)*h])
	}

	return nil
}

// RandReels - random with reels
func (gs *GameScene) RandReels(game IGame, plugin sgc7plugin.IPlugin, reelsName string) error {
	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReels
	}

	for x, arr := range gs.Arr {
		cn, err := plugin.Random(context.Background(), len(reels.Reels[x]))
		if err != nil {
			return err
		}

		for y := range arr {
			gs.Arr[x][y] = reels.Reels[x][cn]

			cn++
			if cn >= len(reels.Reels[x]) {
				cn -= len(reels.Reels[x])
			}
		}
	}

	return nil
}

// FuncForEach - function for ForEach
type FuncForEach func(x, y int, val int)

// ForEachAround - for each around positions
func (gs *GameScene) ForEachAround(x, y int, funcEachAround FuncForEach) {
	if x >= 0 && x < gs.Width && y >= 0 && y < gs.Height {
		for ox := -1; ox <= 1; ox++ {
			for oy := -1; oy <= 1; oy++ {
				if ox == 0 && oy == 0 {
					continue
				}

				if x+ox >= 0 && x+ox < gs.Width && y+oy >= 0 && y+oy < gs.Height {
					funcEachAround(x+ox, y+oy, gs.Arr[x+ox][y+oy])
				}
			}
		}
	}
}

// ForEach - for each all positions
func (gs *GameScene) ForEach(funcEach FuncForEach) {
	for x, l := range gs.Arr {
		for y, v := range l {
			funcEach(x, y, v)
		}
	}
}

// CountSymbol - count a symbol
func (gs *GameScene) CountSymbol(s int) int {
	nums := 0
	for _, l := range gs.Arr {
		for _, v := range l {
			if v == s {
				nums++
			}
		}
	}

	return nums
}

// CountSymbols - count some symbols
func (gs *GameScene) CountSymbols(arr []int) []int {
	narr := make([]int, len(arr))
	for _, l := range gs.Arr {
		for _, v := range l {
			i := sgc7utils.IndexOfIntSlice(arr, v, 0)
			if i >= 0 {
				narr[i]++
			}
		}
	}

	return narr
}

// Clone - clone
func (gs *GameScene) Clone() *GameScene {
	ngs := &GameScene{
		Arr:    make([][]int, gs.Width),
		Width:  gs.Width,
		Height: gs.Height,
	}

	for i := 0; i < gs.Width; i++ {
		ngs.Arr[i] = make([]int, gs.Height)
		copy(ngs.Arr[i], gs.Arr[i])
	}

	return ngs
}

// Fill - fill with reels and indexs
func (gs *GameScene) Fill(reels *ReelsData, arr []int) {
	for x, v := range arr {
		for y := 0; y < gs.Height; y++ {
			gs.Arr[x][y] = reels.Reels[x][v]

			v++
			if v >= len(reels.Reels[x]) {
				v -= len(reels.Reels[x])
			}
		}
	}
}

// CountSymbolEx - count a symbol
func (gs *GameScene) CountSymbolEx(issymbol FuncCountSymbolExIsSymbol) int {
	nums := 0
	for x, l := range gs.Arr {
		for y, v := range l {
			if issymbol(v, x, y) {
				nums++
			}
		}
	}

	return nums
}
