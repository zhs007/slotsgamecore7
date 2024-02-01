package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type ISymbolTrigger interface {
	// CanTrigger -
	CanTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, isSaveResult bool, cd IComponentData) (bool, []*sgc7game.Result)
}
