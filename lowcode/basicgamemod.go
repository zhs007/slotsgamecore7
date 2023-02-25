package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
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
		c := mgrComponent.NewComponent(v.Type)
		bgm.Components.AddComponent(c)
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
