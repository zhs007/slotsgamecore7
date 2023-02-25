package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

// BaseGame - base game
type BaseGame struct {
	*sgc7game.BasicGameMod
	cfg *Config
}

// NewBaseGame - new BaseGame
func NewBaseGame(cfg *Config, game sgc7game.IGame) *BaseGame {
	bg := &BaseGame{
		sgc7game.NewBasicGameMod("bg", cfg.Width, cfg.Height),
		cfg,
	}

	return bg
}

// OnPlay - on play
func (bg *BaseGame) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, error) {
	if cmd == "SPIN" {
		pr := &sgc7game.PlayResult{}

		return pr, nil
	}

	return nil, sgc7game.ErrInvalidCommand
}
