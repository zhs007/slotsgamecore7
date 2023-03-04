package lowcode

import sgc7game "github.com/zhs007/slotsgamecore7/game"

type BasicComponentConfig struct {
	DefaultNextComponent string   `yaml:"defaultNextComponent"` // next component, if it is empty jump to other
	TagScenes            []string `yaml:"tagScenes"`            // tag scenes
	TagOtherScenes       []string `yaml:"tagOtherScenes"`       // tag otherScenes
}

type BasicComponent struct {
	Name            string
	UsedScenes      []int
	UsedOtherScenes []int
	UsedResults     []int
	Config          *BasicComponentConfig
}

// onInit -
func (basicComponent *BasicComponent) onInit(cfg *BasicComponentConfig) {
	basicComponent.Config = cfg
}

// OnNewStep -
func (basicComponent *BasicComponent) OnNewStep() {
	basicComponent.UsedScenes = nil
	basicComponent.UsedOtherScenes = nil
	basicComponent.UsedResults = nil
}

// AddScene -
func (basicComponent *BasicComponent) AddScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene) {

	si := len(curpr.Scenes)
	usi := len(basicComponent.UsedScenes)
	basicComponent.UsedScenes = append(basicComponent.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	if usi < len(basicComponent.Config.TagScenes) {
		gameProp.TagScene(curpr, basicComponent.Config.TagScenes[usi], si)
	}
}

// AddScene -
func (basicComponent *BasicComponent) AddOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene) {

	si := len(curpr.OtherScenes)
	usi := len(basicComponent.UsedOtherScenes)
	basicComponent.UsedOtherScenes = append(basicComponent.UsedOtherScenes, si)

	curpr.OtherScenes = append(curpr.OtherScenes, sc)

	if usi < len(basicComponent.Config.TagOtherScenes) {
		gameProp.TagOtherScene(curpr, basicComponent.Config.TagOtherScenes[usi], si)
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
