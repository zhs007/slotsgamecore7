package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type BasicComponentData struct {
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

// OnNewGame -
func (basicComponentData *BasicComponentData) OnNewGame() {
}

// OnNewGame -
func (basicComponentData *BasicComponentData) OnNewStep() {
	basicComponentData.UsedScenes = nil
	basicComponentData.UsedOtherScenes = nil
	basicComponentData.UsedResults = nil
	basicComponentData.UsedPrizeScenes = nil
	basicComponentData.CashWin = 0
	basicComponentData.CoinWin = 0
	basicComponentData.TargetSceneIndex = -1
	basicComponentData.TargetOtherSceneIndex = -1
	basicComponentData.RNG = nil
}

// BuildPBComponentData
func (basicComponentData *BasicComponentData) BuildPBComponentData() proto.Message {
	return basicComponentData.BuildPBBasicComponentData()
}

// BuildPBBasicComponentData
func (basicComponentData *BasicComponentData) BuildPBBasicComponentData() *sgc7pb.ComponentData {
	pbcd := &sgc7pb.ComponentData{}

	pbcd.CashWin = basicComponentData.CashWin
	pbcd.CoinWin = int32(basicComponentData.CoinWin)
	pbcd.TargetScene = int32(basicComponentData.TargetSceneIndex)

	for _, v := range basicComponentData.UsedOtherScenes {
		pbcd.UsedOtherScenes = append(pbcd.UsedOtherScenes, int32(v))
	}

	for _, v := range basicComponentData.UsedScenes {
		pbcd.UsedScenes = append(pbcd.UsedScenes, int32(v))
	}

	for _, v := range basicComponentData.UsedResults {
		pbcd.UsedResults = append(pbcd.UsedResults, int32(v))
	}

	for _, v := range basicComponentData.UsedPrizeScenes {
		pbcd.UsedPrizeScenes = append(pbcd.UsedPrizeScenes, int32(v))
	}

	return pbcd
}

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
	Config *BasicComponentConfig
	Name   string
}

// onInit -
func (basicComponent *BasicComponent) onInit(cfg *BasicComponentConfig) {
	basicComponent.Config = cfg
}

// onStepEnd -
func (basicComponent *BasicComponent) onStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, nextComponent string) {
	if gameProp.GetVal(GamePropFGNum) > 0 && basicComponent.Config.DefaultFGRespinComponent != "" {
		gameProp.Respin(curpr, gp, basicComponent.Config.DefaultFGRespinComponent, nil, nil)
	} else if nextComponent != "" {
		gameProp.SetStrVal(GamePropNextComponent, nextComponent)
	} else {
		gameProp.SetStrVal(GamePropNextComponent, basicComponent.Config.DefaultNextComponent)
	}
}

// OnNewGame -
func (basicComponent *BasicComponent) OnNewGame(gameProp *GameProperty) error {
	return nil
}

// OnNewStep -
func (basicComponent *BasicComponent) OnNewStep(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[basicComponent.Name]

	cd.OnNewStep()

	return nil
}

// AddScene -
func (basicComponent *BasicComponent) AddScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene, basicCD *BasicComponentData) {

	si := len(curpr.Scenes)
	usi := len(basicCD.UsedScenes)
	basicCD.UsedScenes = append(basicCD.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	if usi < len(basicComponent.Config.TagScenes) {
		gameProp.TagScene(curpr, basicComponent.Config.TagScenes[usi], si)
	}
}

// ReTagScene -
func (basicComponent *BasicComponent) ReTagScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	si int, basicCD *BasicComponentData) {

	usi := len(basicCD.UsedScenes)
	basicCD.UsedScenes = append(basicCD.UsedScenes, si)

	if usi < len(basicComponent.Config.TagScenes) {
		gameProp.TagScene(curpr, basicComponent.Config.TagScenes[usi], si)
	}
}

// AddOtherScene -
func (basicComponent *BasicComponent) AddOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene, basicCD *BasicComponentData) {

	si := len(curpr.OtherScenes)
	usi := len(basicCD.UsedOtherScenes)
	basicCD.UsedOtherScenes = append(basicCD.UsedOtherScenes, si)

	curpr.OtherScenes = append(curpr.OtherScenes, sc)

	if usi < len(basicComponent.Config.TagOtherScenes) {
		gameProp.TagOtherScene(curpr, basicComponent.Config.TagOtherScenes[usi], si)
	}
}

// AddResult -
func (basicComponent *BasicComponent) AddResult(curpr *sgc7game.PlayResult, ret *sgc7game.Result, basicCD *BasicComponentData) {
	basicCD.CoinWin += ret.CoinWin
	basicCD.CashWin += int64(ret.CashWin)

	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	basicCD.UsedResults = append(basicCD.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

// AddRNG -
func (basicComponent *BasicComponent) AddRNG(gameProp *GameProperty, rng int, basicCD *BasicComponentData) {
	i := len(basicCD.RNG)

	basicCD.RNG = append(basicCD.RNG, rng)

	if len(basicComponent.Config.TagRNG) > i {
		gameProp.TagInt(basicComponent.Config.TagRNG[i], rng)
	}
}

// OnStatsWithPB -
func (basicComponent *BasicComponent) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.ComponentData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicComponent.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return basicComponent.OnStatsWithPBBasicComponentData(feature, pbcd, pr), nil
}

// OnStatsWithComponent -
func (basicComponent *BasicComponent) OnStatsWithPBBasicComponentData(feature *sgc7stats.Feature, pbComponent *sgc7pb.ComponentData, pr *sgc7game.PlayResult) int64 {
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
func (basicComponent *BasicComponent) GetTargetScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData) *sgc7game.GameScene {
	gs, si := gameProp.GetScene(curpr, basicComponent.Config.TargetScene)

	if si >= 0 {
		basicCD.TargetSceneIndex = si
	}

	return gs
}

// GetTargetOtherScene -
func (basicComponent *BasicComponent) GetTargetOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData) *sgc7game.GameScene {
	gs, si := gameProp.GetOtherScene(curpr, basicComponent.Config.TargetOtherScene)

	if si >= 0 {
		basicCD.TargetOtherSceneIndex = si
	}

	return gs
}

// NewComponentData -
func (basicComponent *BasicComponent) NewComponentData() IComponentData {
	return &BasicComponentData{}
}

// EachUsedResults -
func (basicComponent *BasicComponent) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.ComponentData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicComponent.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.UsedResults {
		oneach(pr.Results[v])
	}
}

// OnPlayGame - on playgame
func (basicComponent *BasicComponent) OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {
	return nil
}

func NewBasicComponent(name string) *BasicComponent {
	return &BasicComponent{
		Name: name,
	}
}
