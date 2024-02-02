package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type IMask interface {
	// SetMask -
	SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error
	// SetMaskVal -
	SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error
	// SetMaskOnlyTrue -
	SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error
	// // GetMask -
	// GetMask(gameProp *GameProperty, cd IComponentData) []bool
}
