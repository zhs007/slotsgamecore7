package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type SimpleRNG struct {
	curPlugin sgc7plugin.IPlugin
	plugin    *sgc7plugin.FastPlugin
	weights   []int
	curIndex  int
}

func (rng *SimpleRNG) Clone() IRNG {
	return &SimpleRNG{
		plugin: rng.plugin,
	}
}

// OnNewGame -
func (rng *SimpleRNG) OnNewGame(curBetMode int, plugin sgc7plugin.IPlugin) error {
	rng.curPlugin = plugin

	return nil
}

// GetCurRNG -
func (rng *SimpleRNG) GetCurRNG(curBetMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string) {

	return false, -1, rng.curPlugin, ""
}

func (rng *SimpleRNG) IsNeedIterate() bool {
	return false 
}

func (rng *SimpleRNG) IsIterateEnding() bool {
	return true
}

// OnChoiceBranch -
func (rng *SimpleRNG) OnChoiceBranch(curBetMode int, curComponent IComponent, branchName string) error {
	return nil
}

// OnStepEnd -
func (rng *SimpleRNG) OnStepEnd(curBetMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error {
	return nil
}

func NewSimpleRNG(iterateComponent string) IRNG {
	return &SimpleRNG{
		plugin: sgc7plugin.NewFastPlugin(),
	}
}
