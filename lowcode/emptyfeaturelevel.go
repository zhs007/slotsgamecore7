package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type EmptyFeatureLevel struct {
}

// Init -
func (fl *EmptyFeatureLevel) Init() {

}

// OnResult -
func (fl *EmptyFeatureLevel) OnResult(pr *sgc7game.PlayResult) {

}

// CountLevel -
func (fl *EmptyFeatureLevel) CountLevel() int {
	return 0
}

func NewEmptyFeatureLevel() IFeatureLevel {
	return &EmptyFeatureLevel{}
}
