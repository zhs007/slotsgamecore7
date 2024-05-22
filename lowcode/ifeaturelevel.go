package lowcode

type FuncNewFeatureLevel func() IFeatureLevel

type IFeatureLevel interface {
	// Init -
	Init()
	// OnStepEnd -
	OnStepEnd(gameProp *GameProperty, gp *GameParams)
	// CountLevel -
	CountLevel() int
}
