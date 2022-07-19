package sgc7game

import (
	"context"

	goutils "github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

// FuncCountSymbolExIsSymbol -
type FuncCountSymbolExIsSymbol func(cursymbol int, x, y int) bool

// GameScene - game scene
type GameScene struct {
	Arr      [][]int `json:"arr"`
	Width    int     `json:"-"`
	Height   int     `json:"-"`
	Indexes  []int   `json:"indexes"`
	ValidRow []int   `json:"validrow"`
	HeightEx []int   `json:"-"`
}

// isArrEx - is a GameSceneEx
func isArrEx(arr [][]int) bool {
	h := len(arr[0])

	for _, v := range arr {
		if h != len(v) {
			return true
		}
	}

	return false
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

// NewGameSceneEx - new a GameScene
func NewGameSceneEx(heights []int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.InitEx(heights)
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

// NewGameSceneWithArr2Ex - new a GameScene
func NewGameSceneWithArr2Ex(arr [][]int) (*GameScene, error) {
	gs := &GameScene{}

	err := gs.InitWithArr2Ex(arr)
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

// InitEx - init scene
func (gs *GameScene) InitEx(h []int) error {
	gs.Arr = nil
	gs.Height = 0
	for x := 0; x < len(h); x++ {
		gs.Arr = append(gs.Arr, []int{})

		for y := 0; y < h[x]; y++ {
			gs.Arr[x] = append(gs.Arr[x], -1)
		}

		if gs.Height < h[x] {
			gs.Height = h[x]
		}
	}

	gs.Width = len(h)
	gs.HeightEx = h

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

// InitWithArr2Ex - init scene
func (gs *GameScene) InitWithArr2Ex(arr [][]int) error {
	if isArrEx(arr) {
		gs.Arr = nil
		gs.Width = len(arr)
		gs.Height = len(arr[0])

		for _, l := range arr {
			if len(l) > gs.Height {
				gs.Height = len(l)
			}

			gs.Arr = append(gs.Arr, l)

			gs.HeightEx = append(gs.HeightEx, len(l))
		}

		return nil
	}

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
		return ErrInvalidReelsName
	}

	if gs.Indexes == nil {
		gs.Indexes = make([]int, 0, gs.Width)
	} else {
		gs.Indexes = gs.Indexes[0:0:cap(gs.Indexes)]
	}

	for x, arr := range gs.Arr {
		cn, err := plugin.Random(context.Background(), len(reels.Reels[x]))
		if err != nil {
			return err
		}

		gs.Indexes = append(gs.Indexes, cn)

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

// RandReelsEx - random with reels
func (gs *GameScene) RandReelsEx(game IGame, plugin sgc7plugin.IPlugin, reelsName string, rpd *ReelsPosData, nums int) error {
	if nums < 0 {
		nums = 0
	}

	if nums > gs.Width {
		nums = gs.Width
	}

	rarr := []int{}

	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReelsName
	}

	if gs.Indexes == nil {
		gs.Indexes = make([]int, 0, gs.Width)
	} else {
		gs.Indexes = gs.Indexes[0:0:cap(gs.Indexes)]
	}

	for x := 0; x < gs.Width; x++ {
		rarr = append(rarr, x)
	}

	for i := 0; i < gs.Width-nums; i++ {
		cr, err := plugin.Random(context.Background(), len(rarr))
		if err != nil {
			goutils.Error("GameScene.RandReelsEx:Random",
				zap.Error(err))

			return err
		}

		rarr = append(rarr[0:cr], rarr[cr+1:]...)
	}

	for x, arr := range gs.Arr {
		var cn int
		var err error

		if goutils.IndexOfIntSlice(rarr, x, 0) < 0 {
			cn, err = plugin.Random(context.Background(), len(reels.Reels[x]))
		} else {
			cn, err = rpd.RandReel(context.Background(), plugin, x)
		}

		// cn, err := rpd.RandReel(context.Background(), plugin, x)
		// cn, err := plugin.Random(context.Background(), len(reels.Reels[x]))
		if err != nil {
			return err
		}

		gs.Indexes = append(gs.Indexes, cn)

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

// RandReelsEx - random with reels
func (gs *GameScene) RandReelsEx2(game IGame, plugin sgc7plugin.IPlugin, reelsName string, rpd0 *ReelsPosData, rpd1 *ReelsPosData, nums int) error {
	if nums < 0 {
		nums = 0
	}

	if nums > gs.Width {
		nums = gs.Width
	}

	rarr := []int{}

	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReelsName
	}

	if gs.Indexes == nil {
		gs.Indexes = make([]int, 0, gs.Width)
	} else {
		gs.Indexes = gs.Indexes[0:0:cap(gs.Indexes)]
	}

	for x := 0; x < gs.Width; x++ {
		rarr = append(rarr, x)
	}

	for i := 0; i < gs.Width-nums; i++ {
		cr, err := plugin.Random(context.Background(), len(rarr))
		if err != nil {
			goutils.Error("GameScene.RandReelsEx:Random",
				zap.Error(err))

			return err
		}

		rarr = append(rarr[0:cr], rarr[cr+1:]...)
	}

	for x, arr := range gs.Arr {
		var cn int
		var err error

		if goutils.IndexOfIntSlice(rarr, x, 0) < 0 {
			cn, err = rpd1.RandReel(context.Background(), plugin, x)
			// cn, err = plugin.Random(context.Background(), len(reels.Reels[x]))
		} else {
			cn, err = rpd0.RandReel(context.Background(), plugin, x)
		}

		// cn, err := rpd.RandReel(context.Background(), plugin, x)
		// cn, err := plugin.Random(context.Background(), len(reels.Reels[x]))
		if err != nil {
			return err
		}

		gs.Indexes = append(gs.Indexes, cn)

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

// ResetReelIndex - reset reel with index
// 	某些游戏里，可能会出现重新移动某一轴，这个就是移动某一轴的接口
func (gs *GameScene) ResetReelIndex(game IGame, reelsName string, x int, index int) error {
	if x < 0 || x >= gs.Width {
		return ErrInvalidSceneX
	}

	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReelsName
	}

	if gs.Indexes != nil {
		gs.Indexes[x] = index
	}

	for ; index < 0; index += len(reels.Reels[x]) {
	}

	for ; index >= len(reels.Reels[x]); index -= len(reels.Reels[x]) {
	}

	for y := range gs.Arr[x] {
		gs.Arr[x][y] = reels.Reels[x][index]

		index++
		if index >= len(reels.Reels[x]) {
			index -= len(reels.Reels[x])
		}
	}

	return nil
}

// FuncForEach - function for ForEach
type FuncForEach func(x, y int, val int)

// ForEachAround - for each around positions
func (gs *GameScene) ForEachAround(x, y int, funcEachAround FuncForEach) {
	if len(gs.HeightEx) > 0 {
		if x >= 0 && x < gs.Width && y >= 0 && y < gs.HeightEx[x] {
			for ox := -1; ox <= 1; ox++ {
				for oy := -1; oy <= 1; oy++ {
					if ox == 0 && oy == 0 {
						continue
					}

					if x+ox >= 0 && x+ox < gs.Width && y+oy >= 0 && y+oy < gs.HeightEx[x+ox] {
						funcEachAround(x+ox, y+oy, gs.Arr[x+ox][y+oy])
					}
				}
			}
		}

		return
	}

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
			i := goutils.IndexOfIntSlice(arr, v, 0)
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
		ngs.Arr[i] = make([]int, len(gs.Arr[i]))
		copy(ngs.Arr[i], gs.Arr[i])
	}

	if len(gs.HeightEx) > 0 {
		ngs.HeightEx = make([]int, len(gs.HeightEx))
		copy(ngs.HeightEx, gs.HeightEx)
	}

	return ngs
}

// Fill - fill with reels and indexs
func (gs *GameScene) Fill(reels *ReelsData, arr []int) {
	for x, v := range arr {
		for y := 0; y < len(gs.Arr[x]); y++ {
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

// HasSymbol - has a symbol
func (gs *GameScene) HasSymbol(s int) bool {
	for _, l := range gs.Arr {
		for _, v := range l {
			if v == s {
				return true
			}
		}
	}

	return false
}

// SetReels - set reels
func (gs *GameScene) SetReels(game IGame, reelsName string, pos []int) error {
	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReelsName
	}

	if gs.Indexes == nil {
		gs.Indexes = make([]int, 0, gs.Width)
	} else {
		gs.Indexes = gs.Indexes[0:0:cap(gs.Indexes)]
	}

	for x, arr := range gs.Arr {
		cn := pos[x]
		// cn, err := plugin.Random(context.Background(), len(reels.Reels[x]))
		// if err != nil {
		// 	return err
		// }

		gs.Indexes = append(gs.Indexes, cn)

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
