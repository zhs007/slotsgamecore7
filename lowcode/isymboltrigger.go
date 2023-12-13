package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type ISymbolTrigger interface {
	// CanTrigger -
	CanTrigger(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, stake *sgc7game.Stake, isSaveResult bool) (bool, []*sgc7game.Result)
}