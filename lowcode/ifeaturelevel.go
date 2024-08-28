package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type FuncNewFeatureLevel func(bet int) IFeatureLevel

type IFeatureLevel interface {
	// Init -
	Init()
	// OnStepEnd -
	OnStepEnd(gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult)
	// CountLevel -
	CountLevel() int
}
