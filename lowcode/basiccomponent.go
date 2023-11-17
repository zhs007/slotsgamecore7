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

// OnNewStep -
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
	return &sgc7pb.BasicComponentData{
		BasicComponentData: basicComponentData.BuildPBBasicComponentData(),
	}
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
	DefaultNextComponent string            `yaml:"defaultNextComponent" json:"defaultNextComponent"` // next component, if it is empty jump to ending
	TagScenes            []string          `yaml:"tagScenes" json:"tagScenes"`                       // tag scenes
	TagOtherScenes       []string          `yaml:"tagOtherScenes" json:"tagOtherScenes"`             // tag otherScenes
	TargetScene          string            `yaml:"targetScene" json:"targetScene"`                   // target scenes
	TargetOtherScene     string            `yaml:"targetOtherScene" json:"targetOtherScene"`         // target otherscenes
	TagRNG               []string          `yaml:"tagRNG" json:"tagRNG"`                             // tag RNG
	InitStrVals          map[string]string `yaml:"initStrVals" json:"initStrVals"`                   // 只要这个组件被执行，就会初始化这些strvals
	UseFileMapping       bool              `yaml:"useFileMapping" json:"useFileMapping"`             // 兼容性配置，新配置应该一定用filemapping
	ComponentType        string            `yaml:"-" json:"componentType"`                           // 组件类型
}

type BasicComponent struct {
	Config *BasicComponentConfig
	Name   string
}

// onInit -
func (basicComponent *BasicComponent) onPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	for k, v := range basicComponent.Config.InitStrVals {
		gameProp.TagGlobalStr(k, v)
	}

	return nil
}

// onInit -
func (basicComponent *BasicComponent) onInit(cfg *BasicComponentConfig) {
	basicComponent.Config = cfg
}

// onStepEnd -
func (basicComponent *BasicComponent) onStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, nextComponent string) {
	if nextComponent != "" {
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
func (basicComponent *BasicComponent) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.BasicComponentData)
	if !isok {
		goutils.Error("BasicComponent.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return basicComponent.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
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
func (basicComponent *BasicComponent) GetTargetScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, targetScene string) *sgc7game.GameScene {
	if targetScene == "" {
		targetScene = basicComponent.Config.TargetScene
	}

	gs, si := gameProp.GetScene(curpr, targetScene)

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
	pbcd := &sgc7pb.BasicComponentData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicComponent.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

// OnPlayGame - on playgame
func (basicComponent *BasicComponent) OnPlayGameEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {
	return nil
}

// GetName -
func (basicComponent *BasicComponent) GetName() string {
	return basicComponent.Name
}

func NewBasicComponent(name string) *BasicComponent {
	return &BasicComponent{
		Name: name,
	}
}
