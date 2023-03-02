package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type BasicComponentConfig struct {
	DefaultNextComponent string `yaml:"defaultNextComponent"` // next component, if no jump to other
}

type BasicComponent struct {
	Name            string
	UsedScenes      []int
	UsedOtherScenes []int
	UsedResults     []int
}

// AddScene -
func (basicComponent *BasicComponent) OnNewStep() {
	basicComponent.UsedScenes = nil
	basicComponent.UsedOtherScenes = nil
	basicComponent.UsedResults = nil
}

// AddScene -
func (basicComponent *BasicComponent) AddScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, sc *sgc7game.GameScene, tag string) {
	si := len(curpr.Scenes)
	basicComponent.UsedScenes = append(basicComponent.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	if tag != "" {
		gameProp.TagScene(curpr, tag, si)
	}
}

// AddScene -
func (basicComponent *BasicComponent) AddOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, sc *sgc7game.GameScene, tag string) {
	si := len(curpr.OtherScenes)
	basicComponent.UsedOtherScenes = append(basicComponent.UsedOtherScenes, si)

	curpr.OtherScenes = append(curpr.OtherScenes, sc)

	if tag != "" {
		gameProp.TagOtherScene(curpr, tag, si)
	}
}

// AddResult -
func (basicComponent *BasicComponent) AddResult(curpr *sgc7game.PlayResult, ret *sgc7game.Result) {
	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	basicComponent.UsedResults = append(basicComponent.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

func NewBasicComponent(name string) *BasicComponent {
	return &BasicComponent{
		Name: name,
	}
}
