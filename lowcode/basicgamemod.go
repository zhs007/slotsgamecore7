package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
)

// BasicGameMod - basic gamemod
type BasicGameMod struct {
	*sgc7game.BasicGameMod
	GameProp   *GameProperty
	Components *ComponentList
}

// NewBasicGameMod - new BaseGame
func NewBasicGameMod(gameProp *GameProperty, cfgGameMod *GameModConfig, mgrComponent *ComponentMgr) *BasicGameMod {
	bgm := &BasicGameMod{
		sgc7game.NewBasicGameMod(cfgGameMod.Type, gameProp.Config.Width, gameProp.Config.Height),
		gameProp,
		NewComponentList(),
	}

	for _, v := range cfgGameMod.Components {
		c := mgrComponent.NewComponent(v)
		err := c.Init(v.Config, gameProp)
		if err != nil {
			goutils.Error("NewBasicGameMod:Init",
				zap.Error(err))

			return nil
		}

		bgm.Components.AddComponent(c)
	}

	return bgm
}

// OnPlay - on play
func (bgm *BasicGameMod) OnPlay(game sgc7game.IGame, plugin sgc7plugin.IPlugin, cmd string, param string,
	ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) (*sgc7game.PlayResult, error) {

	bgm.GameProp.OnNewStep()

	if cmd == "SPIN" {
		pr := &sgc7game.PlayResult{IsFinish: true, NextGameMod: "bg"}

		for i, v := range bgm.Components.Components {
			err := v.OnPlayGame(bgm.GameProp, pr, plugin, cmd, param, ps, stake, prs)
			if err != nil {
				goutils.Error("BasicGameMod.OnPlay:OnPlayGame",
					zap.Int("i", i),
					zap.Error(err))

				return nil, err
			}
		}

		return pr, nil
	}

	return nil, sgc7game.ErrInvalidCommand
}

// ResetConfig
func (bgm *BasicGameMod) ResetConfig(cfg *Config) {
	bgm.GameProp.Config = cfg
}

// OnAsciiGame - outpur to asciigame
func (bgm *BasicGameMod) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult) error {
	for _, v := range bgm.Components.Components {
		v.OnAsciiGame(bgm.GameProp, pr, lst, gameProp.MapSymbolColor)
	}

	return nil
}
