package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// BasicGameMod - basic gamemod
type BasicGameMod struct {
	*sgc7game.BasicGameMod
	GameProp *GameProperty
}

// NewBasicGameMod - new BaseGame
func NewBasicGameMod(name string, gameProp *GameProperty) *BasicGameMod {
	bgm := &BasicGameMod{
		sgc7game.NewBasicGameMod(name, gameProp.Config.Width, gameProp.Config.Height),
		gameProp,
	}

	return bgm
}

// OnPlay - on play
func (bgm *BasicGameMod) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, error) {

	if cmd == "SPIN" {
		pr := &sgc7game.PlayResult{}

		return pr, nil
	}

	return nil, sgc7game.ErrInvalidCommand
}

// ResetConfig
func (bgm *BasicGameMod) ResetConfig(cfg *Config) {
	bgm.GameProp.Config = cfg
}
