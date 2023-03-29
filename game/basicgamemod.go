package sgc7game

import (
	goutils "github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
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
func (mod *BasicGameMod) OnPlay(game IGame, plugin sgc7plugin.IPlugin, cmd string, param string, ps IPlayerState, stake *Stake, prs []*PlayResult, gameData any) (*PlayResult, error) {
	return nil, ErrInvalidCommand
}

// RandomScene - on random scene
func (mod *BasicGameMod) RandomScene(game IGame, plugin sgc7plugin.IPlugin, reelsName string, gs *GameScene) (*GameScene, error) {
	if mod.Width > 0 && mod.Height > 0 {
		if gs == nil {
			cs, err := NewGameScene(mod.Width, mod.Height)
			if err != nil {
				goutils.Error("sgc7game.BasicGameMod.RandomScene:NewGameScene",
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
func (mod *BasicGameMod) NewPlayResult(gamemodparams any) *PlayResult {
	return &PlayResult{
		CurGameMod:       mod.Name,
		CurGameModParams: gamemodparams,
	}
}

// NewPlayResult2 - new a PlayResult
func (mod *BasicGameMod) NewPlayResult2(gamemodparams any, prs []*PlayResult, parentIndex int, modType string) *PlayResult {
	ci := GetPlayResultCurIndex(prs)
	pr := NewPlayResult(mod.Name, ci, parentIndex, modType)
	pr.CurGameModParams = gamemodparams

	return pr
}
