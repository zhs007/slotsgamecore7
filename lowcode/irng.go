package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type FuncNewRNG func() IRNG

type IRNG interface {
	// Clone -
	Clone() IRNG
	// OnNewGame -
	OnNewGame(betMode int, plugin sgc7plugin.IPlugin) error
	// GetCurRNG -
	GetCurRNG(betMode int, gameProp *GameProperty, curComponent IComponent, cd IComponentData, fl IFeatureLevel) (bool, int, sgc7plugin.IPlugin, string)
	// OnChoiceBranch -
	OnChoiceBranch(betMode int, curComponent IComponent, branchName string) error
	// OnStepEnd -
	OnStepEnd(betMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) error
}
