package sgc7game

import sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"

// GameScene - game scene
type GameScene struct {
	Arr    [][]int `json:"arr"`
	Width  int     `json:"-"`
	Height int     `json:"-"`
}

// NewGameScene - random with reels
func NewGameScene(width int, height int) *GameScene {
	gs := &GameScene{}

	gs.Init(width, height)

	return gs
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

// RandReels - random with reels
func (gs *GameScene) RandReels(game IGame, plugin sgc7plugin.IPlugin, reelsName string) error {
	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReels
	}

	for x, arr := range gs.Arr {
		cn, err := plugin.Random(len(reels.Reels[x]))
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

// FuncForEachAround - function for ForEachAround
type FuncForEachAround func(x, y int, val int)

// ForEachAround - for each around positions
func (gs *GameScene) ForEachAround(x, y int, funcEachAround FuncForEachAround) {
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
