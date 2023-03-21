package lowcode

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
)

type BasicComponentConfig struct {
	DefaultNextComponent     string   `yaml:"defaultNextComponent"`     // next component, if it is empty jump to ending
	DefaultFGRespinComponent string   `yaml:"defaultFGRespinComponent"` // respin component, if it is not empty and in FG
	TagScenes                []string `yaml:"tagScenes"`                // tag scenes
	TagOtherScenes           []string `yaml:"tagOtherScenes"`           // tag otherScenes
	TargetScene              string   `yaml:"targetScene"`              // target scenes
	TargetOtherScene         string   `yaml:"targetOtherScene"`         // target otherscenes
	TagRNG                   []string `yaml:"tagRNG"`                   // tag RNG
}

type BasicComponent struct {
	Config                *BasicComponentConfig
	Name                  string
	UsedScenes            []int
	UsedOtherScenes       []int
	UsedResults           []int
	UsedPrizeScenes       []int
	CashWin               int64
	CoinWin               int
	TargetSceneIndex      int
	TargetOtherSceneIndex int
	RNG                   []int
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
	basicComponent.TargetSceneIndex = -1
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

// AddRNG -
func (basicComponent *BasicComponent) AddRNG(gameProp *GameProperty, rng int) {
	i := len(basicComponent.RNG)

	basicComponent.RNG = append(basicComponent.RNG, rng)

	if len(basicComponent.Config.TagRNG) > i {
		gameProp.MapInt[basicComponent.Config.TagRNG[i]] = rng
	}
}

// BuildPBComponent -
func (basicComponent *BasicComponent) BuildPBComponent(gp *GameParams) {
	pb := &sgc7pb.ComponentData{}

	pb.Name = basicComponent.Name
	pb.CashWin = basicComponent.CashWin
	pb.CoinWin = int32(basicComponent.CoinWin)
	pb.TargetScene = int32(basicComponent.TargetSceneIndex)

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

// BuildPBComponent -
func (basicComponent *BasicComponent) OnStatsWithComponent(feature *sgc7stats.Feature, pbComponent *sgc7pb.ComponentData, pr *sgc7game.PlayResult) int64 {
	wins := int64(0)

	for _, v := range pbComponent.UsedResults {
		ret := pr.Results[v]

		feature.Symbols.OnWin(ret)

		wins += int64(ret.CashWin)
	}

	if pbComponent.TargetScene >= 0 {
		feature.Reels.OnScene(pr.Scenes[pbComponent.TargetScene])
	}

	return wins
}

// GetTargetScene -
func (basicComponent *BasicComponent) GetTargetScene(gameProp *GameProperty, curpr *sgc7game.PlayResult) *sgc7game.GameScene {
	gs, si := gameProp.GetScene(curpr, basicComponent.Config.TargetScene)

	basicComponent.TargetSceneIndex = si

	return gs
}

// GetTargetOtherScene -
func (basicComponent *BasicComponent) GetTargetOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult) *sgc7game.GameScene {
	gs, si := gameProp.GetOtherScene(curpr, basicComponent.Config.TargetOtherScene)

	basicComponent.TargetOtherSceneIndex = si

	return gs
}

func NewBasicComponent(name string) *BasicComponent {
	return &BasicComponent{
		Name: name,
	}
}
