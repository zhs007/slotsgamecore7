package sgc7game

import (
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BasicGameMod - basic gameMod
type BasicGameMod struct {
	Name   string
	Width  int
	Height int
}

// NewBasicGameMod - new a BasicGameMod
func NewBasicGameMod(name string, w int, h int) *BasicGameMod {
	return &BasicGameMod{
		Name:   name,
		Width:  w,
		Height: h,
	}
}

// GetName - get mode name
func (mod *BasicGameMod) GetName() string {
	return mod.Name
}

// OnPlay - on play
func (mod *BasicGameMod) OnPlay(game IGame, plugin sgc7plugin.IPlugin, cmd string, param string, stake *Stake, prs []*PlayResult) (*PlayResult, error) {
	return nil, ErrInvalidCommand
}

// RandomScene - on random scene
func (mod *BasicGameMod) RandomScene(game IGame, plugin sgc7plugin.IPlugin, reelsName string, gs *GameScene) (*GameScene, error) {
	if mod.Width > 0 && mod.Height > 0 {
		if gs == nil {
			cs, err := NewGameScene(mod.Width, mod.Height)
			if err != nil {
				sgc7utils.Error("sgc7game.BasicGameMod.RandomScene:NewGameScene",
					zap.Int("width", mod.Width),
					zap.Int("height", mod.Height),
					zap.Error(err))

				return nil, err
			}

			gs = cs
		}

		err := gs.RandReels(game, plugin, reelsName)
		if err != nil {
			return nil, err
		}

		return gs, nil
	}

	return nil, ErrInvalidWHGameMod
}

// NewPlayResult - new a PlayResult
func (mod *BasicGameMod) NewPlayResult(gamemodparams interface{}) *PlayResult {
	return &PlayResult{
		CurGameMod:       mod.Name,
		CurGameModParams: gamemodparams,
	}
}
