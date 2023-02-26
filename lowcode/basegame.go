package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// BaseGame - base game
type BaseGame struct {
	*BasicGameMod
}

// NewBaseGame - new BaseGame
func NewBaseGame(gameProp *GameProperty, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) sgc7game.IGameMod {
	bg := &BaseGame{
		NewBasicGameMod(gameProp, cfgGameMod, mgrComponent),
	}

	return bg
}

// OnPlay - on play
func (bg *BaseGame) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, error) {
	if cmd == "SPIN" {
		pr := &sgc7game.PlayResult{IsFinish: true, NextGameMod: "bg"}

		return pr, nil
	}

	return nil, sgc7game.ErrInvalidCommand
}
