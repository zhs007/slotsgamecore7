package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type FuncNewFeatureLevel func() IFeatureLevel

type IFeatureLevel interface {
	// Init -
	Init()
	// OnResult -
	OnResult(pr *sgc7game.PlayResult)
}
