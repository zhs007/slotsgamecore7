package sgc7game

// GameScene - game scene
type GameScene struct {
	Arr [][]int `json:"arr"`
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

	return nil
}

// RandReels - random with reels
func (gs *GameScene) RandReels(game IGame, reelsName string) error {
	cfg := game.GetConfig()

	reels, isok := cfg.Reels[reelsName]
	if !isok {
		return ErrInvalidReels
	}

	for x, arr := range gs.Arr {
		cn, err := game.GetPlugin().Random(len(reels.Reels[x]))
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
