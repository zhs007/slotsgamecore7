package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"github.com/zhs007/slotsgamecore7/stats2"
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
	MapConfigVals         map[string]string
	MapConfigIntVals      map[string]int
	SrcScenes             []int
}

// OnNewGame -
func (basicComponentData *BasicComponentData) OnNewGame(gameProp *GameProperty, component IComponent) {
	basicComponentData.MapConfigVals = make(map[string]string)
	basicComponentData.MapConfigIntVals = make(map[string]int)
}

// OnNewStep -
func (basicComponentData *BasicComponentData) OnNewStep(gameProp *GameProperty, component IComponent) {
	basicComponentData.UsedScenes = nil
	basicComponentData.UsedOtherScenes = nil
	basicComponentData.UsedResults = nil
	basicComponentData.UsedPrizeScenes = nil
	basicComponentData.CashWin = 0
	basicComponentData.CoinWin = 0
	basicComponentData.TargetSceneIndex = -1
	basicComponentData.TargetOtherSceneIndex = -1
	basicComponentData.RNG = nil

	basicComponentData.initSrcScenes()
}

// GetVal -
func (basicComponentData *BasicComponentData) GetVal(key string) int {
	return 0
}

// SetVal -
func (basicComponentData *BasicComponentData) SetVal(key string, val int) {

}

// GetConfigVal -
func (basicComponentData *BasicComponentData) GetConfigVal(key string) string {
	return basicComponentData.MapConfigVals[key]
}

// SetConfigVal -
func (basicComponentData *BasicComponentData) SetConfigVal(key string, val string) {
	basicComponentData.MapConfigVals[key] = val
}

// GetConfigIntVal -
func (basicComponentData *BasicComponentData) GetConfigIntVal(key string) (int, bool) {
	ival, isok := basicComponentData.MapConfigIntVals[key]
	return ival, isok
}

// SetConfigIntVal -
func (basicComponentData *BasicComponentData) SetConfigIntVal(key string, val int) {
	basicComponentData.MapConfigIntVals[key] = val
}

// ChgConfigIntVal -
func (basicComponentData *BasicComponentData) ChgConfigIntVal(key string, off int) {
	basicComponentData.MapConfigIntVals[key] += off
}

// InitSrcScenes -
func (basicComponentData *BasicComponentData) initSrcScenes() {
	for i := range basicComponentData.SrcScenes {
		basicComponentData.SrcScenes[i] = -1
	}
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

	for _, v := range basicComponentData.SrcScenes {
		pbcd.SrcScenes = append(pbcd.SrcScenes, int32(v))
	}

	return pbcd
}

// GetResults -
func (basicComponentData *BasicComponentData) GetResults() []int {
	return basicComponentData.UsedResults
}

// 新思路：尽量弱化变量的概念，所有变量都放到component里面去，譬如循环、scene、分支等，这样逻辑会更清晰
type BasicComponentConfig struct {
	DefaultNextComponent   string            `yaml:"defaultNextComponent" json:"defaultNextComponent"`     // next component, if it is empty jump to ending
	TagScenes              []string          `yaml:"tagScenes" json:"tagScenes"`                           // tag scenes，v0.13开始弃用
	TagOtherScenes         []string          `yaml:"tagOtherScenes" json:"tagOtherScenes"`                 // tag otherScenes，v0.13开始弃用
	TargetScene            string            `yaml:"targetScene" json:"targetScene"`                       // target scenes，v0.13开始弃用
	TargetOtherScene       string            `yaml:"targetOtherScene" json:"targetOtherScene"`             // target otherscenes，v0.13开始弃用
	TagGlobalScenes        []string          `yaml:"tagGlobalScenes" json:"tagGlobalScenes"`               // tag global scenes，v0.13开始弃用
	TargetGlobalScene      string            `yaml:"targetGlobalScene" json:"targetGlobalScene"`           // target global scenes，v0.13开始弃用
	TagGlobalOtherScenes   []string          `yaml:"tagGlobalOtherScenes" json:"tagGlobalOtherScenes"`     // tag global other scenes，v0.13开始弃用
	TargetGlobalOtherScene string            `yaml:"targetGlobalOtherScene" json:"targetGlobalOtherScene"` // target global other scenes，v0.13开始弃用
	Scene2Components       []string          `yaml:"scene2Components" json:"scene2Components"`             // 新版本，关于scene换了个思路，用目标对象来驱动
	OtherScene2Components  []string          `yaml:"otherScene2Components" json:"otherScene2Components"`   // 新版本，关于other scene换了个思路，用目标对象来驱动
	TagRNG                 []string          `yaml:"tagRNG" json:"tagRNG"`                                 // tag RNG
	InitStrVals            map[string]string `yaml:"initStrVals" json:"initStrVals"`                       // 只要这个组件被执行，就会初始化这些strvals
	UseFileMapping         bool              `yaml:"useFileMapping" json:"useFileMapping"`                 // 兼容性配置，新配置应该一定用filemapping
	ComponentType          string            `yaml:"-" json:"componentType"`                               // 组件类型
	UseSceneV2             bool              `yaml:"useSceneV2" json:"useSceneV2"`                         // 新版本的scene
	TargetScenes3          [][]string        `yaml:"targetScenes3" json:"targetScenes3"`                   // target scenes V3
	UseSceneV3             bool              `yaml:"useSceneV3" json:"useSceneV3"`                         // 新版本的scene
	IsNeedCacheScene3      bool              `yaml:"isNeedCacheScene3" json:"isNeedCacheScene3"`           // 是否需要缓存scene
}

type BasicComponent struct {
	Config            *BasicComponentConfig
	Name              string
	SrcSceneNum       int
	dataForeachSymbol *ForeachSymbolData
}

// Init -
func (basicComponent *BasicComponent) Init(fn string, pool *GamePropertyPool) error {
	return nil
}

// InitEx -
func (basicComponent *BasicComponent) InitEx(cfg any, pool *GamePropertyPool) error {
	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicComponent *BasicComponent) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnPlayGame - on playgame
func (basicComponent *BasicComponent) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {
	return nil
}

// OnStats -
func (basicComponent *BasicComponent) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
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
	if nextComponent == "" {
		nextComponent = basicComponent.Config.DefaultNextComponent
	}

	component, isok := gameProp.Components.MapComponents[nextComponent]
	if isok && component.IsRespin() {
		gameProp.SetStrVal(GamePropRespinComponent, nextComponent)
		gameProp.onTriggerRespin(nextComponent)

		gp.NextStepFirstComponent = nextComponent

		gameProp.SetStrVal(GamePropNextComponent, "")

		return
	}

	gameProp.SetStrVal(GamePropNextComponent, nextComponent)
}

// // OnNewGame -
// func (basicComponent *BasicComponent) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.GetCurComponentData(basicComponent)

// 	cd.OnNewGame()

// 	return nil
// }

// // OnNewStep -
// func (basicComponent *BasicComponent) OnNewStep(gameProp *GameProperty) error {
// 	cd := gameProp.GetCurComponentData(basicComponent)

// 	cd.OnNewStep()

// 	return nil
// }

// AddScene -
func (basicComponent *BasicComponent) AddScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene, basicCD *BasicComponentData) {

	si := len(curpr.Scenes)
	usi := len(basicCD.UsedScenes)
	basicCD.UsedScenes = append(basicCD.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	if basicComponent.Config.UseSceneV3 {
		gameProp.SceneStack.Push(basicComponent.Name, si, sc, basicComponent.Config.IsNeedCacheScene3)
	} else if basicComponent.Config.UseSceneV2 {
		if len(basicComponent.Config.Scene2Components) > 0 {
			for _, v := range basicComponent.Config.Scene2Components {
				gameProp.SetComponentScene(v, sc)
			}
		}
	} else {
		if usi < len(basicComponent.Config.TagScenes) {
			gameProp.TagScene(curpr, basicComponent.Config.TagScenes[usi], si)
		}

		if usi < len(basicComponent.Config.TagGlobalScenes) {
			gameProp.TagGlobalScene(basicComponent.Config.TagGlobalScenes[usi], sc)
		}
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

	if basicComponent.Config.UseSceneV2 {
		if len(basicComponent.Config.OtherScene2Components) > 0 {
			for _, v := range basicComponent.Config.OtherScene2Components {
				gameProp.SetComponentOtherScene(v, sc)
			}
		}
	} else {
		if usi < len(basicComponent.Config.TagOtherScenes) {
			gameProp.TagOtherScene(curpr, basicComponent.Config.TagOtherScenes[usi], si)
		}
	}
}

// ClearOtherScene -
func (basicComponent *BasicComponent) ClearOtherScene(gameProp *GameProperty) {
	if basicComponent.Config.UseSceneV2 {
		if len(basicComponent.Config.OtherScene2Components) > 0 {
			for _, v := range basicComponent.Config.OtherScene2Components {
				gameProp.ClearComponentOtherScene(v)
			}
		}
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
		if basicComponent.Config.TargetGlobalScene != "" {
			return gameProp.GetGlobalScene(basicComponent.Config.TargetGlobalScene)
		} else {
			targetScene = basicComponent.Config.TargetScene
		}
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
	bcd := &BasicComponentData{}

	if basicComponent.SrcSceneNum > 0 {
		bcd.SrcScenes = make([]int, basicComponent.SrcSceneNum)
	}

	return bcd
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

// // SetMask -
// func (basicComponent *BasicComponent) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, mask []bool) error {
// 	return ErrNotMask
// }

// // GetMask -
// func (basicComponent *BasicComponent) GetMask(gameProp *GameProperty) []bool {
// 	return nil
// }

// IsRespin -
func (basicComponent *BasicComponent) IsRespin() bool {
	return false
}

// IsMask -
func (basicComponent *BasicComponent) IsMask() bool {
	return false
}

func (basicComponent *BasicComponent) GetTargetScene2(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, component string, tag string) *sgc7game.GameScene {
	if basicComponent.Config.UseSceneV2 {
		return gameProp.GetComponentScene(component)
	}

	if tag == "" {
		if basicComponent.Config.TargetGlobalScene != "" {
			return gameProp.GetGlobalScene(basicComponent.Config.TargetGlobalScene)
		} else {
			tag = basicComponent.Config.TargetScene
		}
	}

	gs, si := gameProp.GetScene(curpr, tag)

	if si >= 0 {
		basicCD.TargetSceneIndex = si
	}

	return gs
}

func (basicComponent *BasicComponent) GetTargetScene3(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, basicCD *BasicComponentData, component string, tag string, si int) *sgc7game.GameScene {
	if basicComponent.Config.UseSceneV3 {
		return gameProp.SceneStack.GetTargetScene3(gameProp, basicComponent.Config, si, curpr, prs)
	}

	if basicComponent.Config.UseSceneV2 {
		return gameProp.GetComponentScene(component)
	}

	if tag == "" {
		if basicComponent.Config.TargetGlobalScene != "" {
			return gameProp.GetGlobalScene(basicComponent.Config.TargetGlobalScene)
		} else {
			tag = basicComponent.Config.TargetScene
		}
	}

	gs, si := gameProp.GetScene(curpr, tag)

	if si >= 0 {
		basicCD.TargetSceneIndex = si
	}

	return gs
}

func (basicComponent *BasicComponent) GetTargetOtherScene2(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, component string, tag string) *sgc7game.GameScene {
	if basicComponent.Config.UseSceneV2 {
		return gameProp.GetComponentOtherScene(component)
	}

	if tag == "" {
		tag = basicComponent.Config.TargetOtherScene
	}

	gs, _ := gameProp.GetOtherScene(curpr, tag)

	return gs
}

// NewStats2 -
func (basicComponent *BasicComponent) NewStats2() *stats2.Feature {
	return nil
}

// OnStats2
func (basicComponent *BasicComponent) OnStats2(icd IComponentData, s2 *stats2.Stats) {
}

// // OnStats2Trigger
// func (basicComponent *BasicComponent) OnStats2Trigger(s2 *Stats2) {

// }

// GetSymbols -
func (basicComponent *BasicComponent) GetSymbols(gameProp *GameProperty) []int {
	return nil
}

// AddSymbol -
func (basicComponent *BasicComponent) AddSymbol(gameProp *GameProperty, symbol int) {

}

// OnEachSymbol - on foreach symbol
func (basicComponent *BasicComponent) OnEachSymbol(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState,
	stake *sgc7game.Stake, prs []*sgc7game.PlayResult, symbol int, cd IComponentData) (string, error) {
	return "", nil
}

// ForEachSymbols - foreach symbols
func (basicComponent *BasicComponent) ForeachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
	prs []*sgc7game.PlayResult) error {
	return nil
}

// // GetComponentData -
// func (basicComponent *BasicComponent) GetComponentData(gameProp *GameProperty) IComponentData {
// 	return gameProp.GetCurComponentData(basicComponent)

// 	if basicComponent.dataForeachSymbol != nil {
// 		return gameProp.MapComponentData[fmt.Sprintf("%v:%v", basicComponent.Name, basicComponent.dataForeachSymbol.Index)]
// 	}

// 	return gameProp.MapComponentData[basicComponent.Name]
// }

// SetForeachSymbolData -
func (basicComponent *BasicComponent) SetForeachSymbolData(data *ForeachSymbolData) {
	basicComponent.dataForeachSymbol = data
}

// OnGameInited - on game inited
func (basicComponent *BasicComponent) OnGameInited(components *ComponentList) error {
	return nil
}

// GetAllLinkComponents - get all link components
func (basicComponent *BasicComponent) GetAllLinkComponents() []string {
	return []string{basicComponent.Config.DefaultNextComponent}
}

func NewBasicComponent(name string, srcSceneNum int) *BasicComponent {
	return &BasicComponent{
		Name:        name,
		SrcSceneNum: srcSceneNum,
	}
}
