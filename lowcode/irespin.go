package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

type IRespin interface {
	// AddRespinTimes -
	AddRespinTimes(gameProp *GameProperty, num int)
	// SaveRetriggerRespinNum -
	SaveRetriggerRespinNum(gameProp *GameProperty)
	// // AddRetriggerRespinNum -
	// AddRetriggerRespinNum(gameProp *GameProperty, num int)
	// // Retrigger -
	// Retrigger(gameProp *GameProperty)
	// Trigger -
	Trigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams)
	// PushTrigger -
	PushTrigger(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int)
	// // AddTriggerAward -
	// AddTriggerAward(gameProp *GameProperty, award *Award)

	// // GetLastRespinNum -
	// GetLastRespinNum(gameProp *GameProperty) int

	// // IsEnding -
	// IsEnding(gameProp *GameProperty) bool
	// // IsStarted -
	// IsStarted(gameProp *GameProperty) bool
}
