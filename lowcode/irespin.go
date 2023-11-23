package lowcode

type IRespin interface {
	// AddRespinTimes -
	AddRespinTimes(gameProp *GameProperty, num int)
	// SaveRetriggerRespinNum -
	SaveRetriggerRespinNum(gameProp *GameProperty)
	// AddRetriggerRespinNum -
	AddRetriggerRespinNum(gameProp *GameProperty, num int)
	// Retrigger -
	Retrigger(gameProp *GameProperty)
}
