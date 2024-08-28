package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func NewBasicRNG() IRNG {
	return &BasicRNG{}
}

type BasicRNG struct {
	curPlugin sgc7plugin.IPlugin
}

func (rng *BasicRNG) Clone() IRNG {
	return &BasicRNG{
		curPlugin: rng.curPlugin,
	}
}

// OnNewGame -
func (rng *BasicRNG) OnNewGame(betMode int, plugin sgc7plugin.IPlugin) error {
	rng.curPlugin = plugin

	return nil
}

// GetCurRNG -
func (rng *BasicRNG) GetCurRNG(betMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {
	return false, -1, rng.curPlugin, ""
}

// OnChoiceBranch -
func (rng *BasicRNG) OnChoiceBranch(betMode int, curComponent IComponent, branchName string) error {
	return nil
}

// OnStepEnd -
func (rng *BasicRNG) OnStepEnd(betMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error {
	return nil
}
