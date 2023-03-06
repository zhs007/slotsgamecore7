package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
)

type BasicComponentConfig struct {
	DefaultNextComponent     string   `yaml:"defaultNextComponent"`     // next component, if it is empty jump to ending
	DefaultFGRespinComponent string   `yaml:"defaultFGRespinComponent"` // respin component, if it is not empty and in FG
	TagScenes                []string `yaml:"tagScenes"`                // tag scenes
	TagOtherScenes           []string `yaml:"tagOtherScenes"`           // tag otherScenes
}

type BasicComponent struct {
	Name            string
	UsedScenes      []int
	UsedOtherScenes []int
	UsedResults     []int
	UsedPrizeScenes []int
	Config          *BasicComponentConfig
	CashWin         int64
	CoinWin         int
}

// onInit -
func (basicComponent *BasicComponent) onInit(cfg *BasicComponentConfig) {
	basicComponent.Config = cfg
}

// onStepEnd -
func (basicComponent *BasicComponent) onStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams) {
	if gameProp.GetVal(GamePropFGNum) > 0 && basicComponent.Config.DefaultFGRespinComponent != "" {
		gameProp.Respin(curpr, gp, basicComponent.Config.DefaultFGRespinComponent, nil, nil)
	} else {
		gameProp.SetStrVal(GamePropNextComponent, basicComponent.Config.DefaultNextComponent)
	}
}

// OnNewStep -
func (basicComponent *BasicComponent) OnNewStep() {
	basicComponent.UsedScenes = nil
	basicComponent.UsedOtherScenes = nil
	basicComponent.UsedResults = nil
	basicComponent.UsedPrizeScenes = nil
	basicComponent.CashWin = 0
	basicComponent.CoinWin = 0
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
	basicComponent.CoinWin += ret.CoinWin
	basicComponent.CashWin += int64(ret.CashWin)

	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	basicComponent.UsedResults = append(basicComponent.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

// BuildPBComponent -
func (basicComponent *BasicComponent) BuildPBComponent(gp *GameParams) {
	pb := &sgc7pb.ComponentData{}

	pb.Name = basicComponent.Name
	pb.CashWin = basicComponent.CashWin
	pb.CoinWin = int32(basicComponent.CoinWin)

	for _, v := range basicComponent.UsedOtherScenes {
		pb.UsedOtherScenes = append(pb.UsedOtherScenes, int32(v))
	}

	for _, v := range basicComponent.UsedScenes {
		pb.UsedScenes = append(pb.UsedScenes, int32(v))
	}

	for _, v := range basicComponent.UsedResults {
		pb.UsedResults = append(pb.UsedResults, int32(v))
	}

	for _, v := range basicComponent.UsedPrizeScenes {
		pb.UsedPrizeScenes = append(pb.UsedPrizeScenes, int32(v))
	}

	gp.MapComponents[pb.Name] = pb
}

func NewBasicComponent(name string) *BasicComponent {
	return &BasicComponent{
		Name: name,
	}
}
