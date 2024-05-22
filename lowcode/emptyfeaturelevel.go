package lowcode

type EmptyFeatureLevel struct {
}

// Init -
func (fl *EmptyFeatureLevel) Init() {

}

// OnStepEnd -
func (fl *EmptyFeatureLevel) OnStepEnd(gameProp *GameProperty, gp *GameParams) {

}

// CountLevel -
func (fl *EmptyFeatureLevel) CountLevel() int {
	return 0
}

func NewEmptyFeatureLevel() IFeatureLevel {
	return &EmptyFeatureLevel{}
}
