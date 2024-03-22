package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
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
	Output                int
	StrOutput             string
}

// OnNewGame -
func (basicComponentData *BasicComponentData) OnNewGame(gameProp *GameProperty, component IComponent) {
	basicComponentData.MapConfigVals = make(map[string]string)
	basicComponentData.MapConfigIntVals = make(map[string]int)
}

// // OnNewStep -
// func (basicComponentData *BasicComponentData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	basicComponentData.UsedScenes = nil
// 	basicComponentData.UsedOtherScenes = nil
// 	basicComponentData.UsedResults = nil
// 	basicComponentData.UsedPrizeScenes = nil
// 	basicComponentData.CashWin = 0
// 	basicComponentData.CoinWin = 0
// 	basicComponentData.TargetSceneIndex = -1
// 	basicComponentData.TargetOtherSceneIndex = -1
// 	basicComponentData.RNG = nil

// 	basicComponentData.initSrcScenes()
// }

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

// ClearConfigIntVal -
func (basicComponentData *BasicComponentData) ClearConfigIntVal(key string) {
	delete(basicComponentData.MapConfigIntVals, key)
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
	pbcd.Output = int32(basicComponentData.Output)
	pbcd.StrOutput = basicComponentData.StrOutput

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

// GetOutput -
func (basicComponentData *BasicComponentData) GetOutput() int {
	return basicComponentData.Output
}

// GetStringOutput -
func (basicComponentData *BasicComponentData) GetStringOutput() string {
	return basicComponentData.StrOutput
}

// GetSymbols -
func (basicComponentData *BasicComponentData) GetSymbols() []int {
	return nil
}

// AddSymbol -
func (basicComponentData *BasicComponentData) AddSymbol(symbolCode int) {

}

// GetPos -
func (basicComponentData *BasicComponentData) GetPos() []int {
	return nil
}

// HasPos -
func (basicComponentData *BasicComponentData) HasPos(x int, y int) bool {
	return false
}

// AddPos -
func (basicComponentData *BasicComponentData) AddPos(x int, y int) {
}

// GetLastRespinNum -
func (basicComponentData *BasicComponentData) GetLastRespinNum() int {
	return 0
}

// IsRespinEnding -
func (basicComponentData *BasicComponentData) IsRespinEnding() bool {
	return false
}

// IsRespinStarted -
func (basicComponentData *BasicComponentData) IsRespinStarted() bool {
	return false
}

// // AddRetriggerRespinNum -
// func (basicComponentData *BasicComponentData) AddRetriggerRespinNum(num int) {

// }

// AddTriggerRespinAward -
func (basicComponentData *BasicComponentData) AddTriggerRespinAward(award *Award) {

}

// AddRespinTimes -
func (basicComponentData *BasicComponentData) AddRespinTimes(num int) {

}

// TriggerRespin
func (basicComponentData *BasicComponentData) TriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {

}

// PushTrigger -
func (basicComponentData *BasicComponentData) PushTriggerRespin(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, num int) {

}

// // SaveRetriggerRespinNum -
// func (basicComponentData *BasicComponentData) SaveRetriggerRespinNum()

// GetMask -
func (basicComponentData *BasicComponentData) GetMask() []bool {
	return nil
}

// ChgMask -
func (basicComponentData *BasicComponentData) ChgMask(curMask int, val bool) bool {
	return false
}

func (basicComponentData *BasicComponentData) PutInMoney(coins int) {

}

// 新思路：尽量弱化变量的概念，所有变量都放到component里面去，譬如循环、scene、分支等，这样逻辑会更清晰
type BasicComponentConfig struct {
	DefaultNextComponent string            `yaml:"defaultNextComponent" json:"defaultNextComponent"` // next component, if it is empty jump to ending
	TagRNG               []string          `yaml:"tagRNG" json:"tagRNG"`                             // tag RNG
	InitStrVals          map[string]string `yaml:"initStrVals" json:"initStrVals"`                   // 只要这个组件被执行，就会初始化这些strvals
	UseFileMapping       bool              `yaml:"useFileMapping" json:"useFileMapping"`             // 兼容性配置，新配置应该一定用filemapping
	ComponentType        string            `yaml:"-" json:"componentType"`                           // 组件类型
	TargetScenes3        [][]string        `yaml:"targetScenes3" json:"targetScenes3"`               // target scenes V3
	TargetOtherScenes3   [][]string        `yaml:"targetOtherScenes3" json:"targetOtherScenes3"`     // target scenes V3
}

type BasicComponent struct {
	Config      *BasicComponentConfig
	Name        string
	SrcSceneNum int
	// dataForeachSymbol *ForeachSymbolData
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

// // OnStats -
// func (basicComponent *BasicComponent) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // onInit -
// func (basicComponent *BasicComponent) onPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
// 	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

// 	// for k, v := range basicComponent.Config.InitStrVals {
// 	// 	gameProp.TagGlobalStr(k, v)
// 	// }

// 	return nil
// }

// onInit -
func (basicComponent *BasicComponent) onInit(cfg *BasicComponentConfig) {
	basicComponent.Config = cfg
}

// onStepEnd -
func (basicComponent *BasicComponent) onStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, nextComponent string) string {
	if nextComponent == "" {
		nextComponent = basicComponent.Config.DefaultNextComponent
	}

	// component, isok := gameProp.Components.MapComponents[nextComponent]
	// if isok && component.IsRespin() {
	// 	// gameProp.SetStrVal(GamePropRespinComponent, nextComponent)
	// 	// gameProp.onTriggerRespin(nextComponent)

	// 	// gp.NextStepFirstComponent = nextComponent

	// 	// gameProp.SetStrVal(GamePropNextComponent, "")

	// 	return nextComponent
	// }

	// gameProp.SetStrVal(GamePropNextComponent, nextComponent)

	return nextComponent
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
	// usi := len(basicCD.UsedScenes)
	basicCD.UsedScenes = append(basicCD.UsedScenes, si)

	curpr.Scenes = append(curpr.Scenes, sc)

	gameProp.SceneStack.Push(basicComponent.Name, si, sc)
}

// // ReTagScene -
// func (basicComponent *BasicComponent) ReTagScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
// 	si int, basicCD *BasicComponentData) {

// 	usi := len(basicCD.UsedScenes)
// 	basicCD.UsedScenes = append(basicCD.UsedScenes, si)

// 	if usi < len(basicComponent.Config.TagScenes) {
// 		gameProp.TagScene(curpr, basicComponent.Config.TagScenes[usi], si)
// 	}
// }

// AddOtherScene -
func (basicComponent *BasicComponent) AddOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	sc *sgc7game.GameScene, basicCD *BasicComponentData) {

	si := len(curpr.OtherScenes)
	// usi := len(basicCD.UsedOtherScenes)
	basicCD.UsedOtherScenes = append(basicCD.UsedOtherScenes, si)

	curpr.OtherScenes = append(curpr.OtherScenes, sc)

	gameProp.OtherSceneStack.Push(basicComponent.Name, si, sc)
}

// ClearOtherScene -
func (basicComponent *BasicComponent) ClearOtherScene(gameProp *GameProperty) {
	// if basicComponent.Config.UseSceneV2 {
	// 	if len(basicComponent.Config.OtherScene2Components) > 0 {
	// 		for _, v := range basicComponent.Config.OtherScene2Components {
	// 			gameProp.ClearComponentOtherScene(v)
	// 		}
	// 	}
	// }
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

// // AddRNG -
// func (basicComponent *BasicComponent) AddRNG(gameProp *GameProperty, rng int, basicCD *BasicComponentData) {
// 	i := len(basicCD.RNG)

// 	basicCD.RNG = append(basicCD.RNG, rng)

// 	if len(basicComponent.Config.TagRNG) > i {
// 		gameProp.TagInt(basicComponent.Config.TagRNG[i], rng)
// 	}
// }

// // OnStatsWithPB -
// func (basicComponent *BasicComponent) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.BasicComponentData)
// 	if !isok {
// 		goutils.Error("BasicComponent.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return basicComponent.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStatsWithComponent -
// func (basicComponent *BasicComponent) OnStatsWithPBBasicComponentData(feature *sgc7stats.Feature, pbComponent *sgc7pb.ComponentData, pr *sgc7game.PlayResult) int64 {
// 	wins := int64(0)

// 	for _, v := range pbComponent.UsedResults {
// 		ret := pr.Results[v]

// 		feature.Symbols.OnWin(ret)

// 		wins += int64(ret.CashWin)
// 	}

// 	if pbComponent.TargetScene >= 0 {
// 		feature.Reels.OnScene(pr.Scenes[pbComponent.TargetScene])
// 	}

// 	return wins
// }

// // GetTargetScene -
// func (basicComponent *BasicComponent) GetTargetScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, targetScene string) *sgc7game.GameScene {
// 	if targetScene == "" {
// 		if basicComponent.Config.TargetGlobalScene != "" {
// 			return gameProp.GetGlobalScene(basicComponent.Config.TargetGlobalScene)
// 		} else {
// 			targetScene = basicComponent.Config.TargetScene
// 		}
// 	}

// 	gs, si := gameProp.GetScene(curpr, targetScene)

// 	if si >= 0 {
// 		basicCD.TargetSceneIndex = si
// 	}

// 	return gs
// }

// // GetTargetOtherScene -
// func (basicComponent *BasicComponent) GetTargetOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData) *sgc7game.GameScene {
// 	gs, si := gameProp.GetOtherScene(curpr, basicComponent.Config.TargetOtherScene)

// 	if si >= 0 {
// 		basicCD.TargetOtherSceneIndex = si
// 	}

// 	return gs
// }

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
			goutils.Err(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

// ProcRespinOnStepEnd - 现在只有respin需要特殊处理结束，如果多层respin嵌套时，只要新的有next，就不会继续结束respin
func (basicComponent *BasicComponent) ProcRespinOnStepEnd(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, canRemove bool) (string, error) {
	return "", nil
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

// IsForeach -
func (basicComponent *BasicComponent) IsForeach() bool {
	return false
}

// // IsTriggerRespin -
// func (basicComponent *BasicComponent) IsTriggerRespin() bool {
// 	return false
// }

// IsMask -
func (basicComponent *BasicComponent) IsMask() bool {
	return false
}

// func (basicComponent *BasicComponent) GetTargetScene2(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, component string, tag string) *sgc7game.GameScene {
// 	if basicComponent.Config.UseSceneV2 {
// 		return gameProp.GetComponentScene(component)
// 	}

// 	if tag == "" {
// 		if basicComponent.Config.TargetGlobalScene != "" {
// 			return gameProp.GetGlobalScene(basicComponent.Config.TargetGlobalScene)
// 		} else {
// 			tag = basicComponent.Config.TargetScene
// 		}
// 	}

// 	gs, si := gameProp.GetScene(curpr, tag)

// 	if si >= 0 {
// 		basicCD.TargetSceneIndex = si
// 	}

// 	return gs
// }

func (basicComponent *BasicComponent) GetTargetScene3(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, si int) *sgc7game.GameScene {
	return gameProp.SceneStack.GetTargetScene3(gameProp, basicComponent.Config, si, curpr, prs)
}

// func (basicComponent *BasicComponent) GetTargetOtherScene2(gameProp *GameProperty, curpr *sgc7game.PlayResult, basicCD *BasicComponentData, component string, tag string) *sgc7game.GameScene {
// 	if basicComponent.Config.UseSceneV2 {
// 		return gameProp.GetComponentOtherScene(component)
// 	}

// 	if tag == "" {
// 		tag = basicComponent.Config.TargetOtherScene
// 	}

// 	gs, _ := gameProp.GetOtherScene(curpr, tag)

// 	return gs
// }

func (basicComponent *BasicComponent) GetTargetOtherScene3(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, si int) *sgc7game.GameScene {
	return gameProp.OtherSceneStack.GetTargetScene3(gameProp, basicComponent.Config, si, curpr, prs)
}

// NewStats2 -
func (basicComponent *BasicComponent) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, nil)
}

// OnStats2
func (basicComponent *BasicComponent) OnStats2(icd IComponentData, s2 *stats2.Cache) {
	s2.ProcStatsTrigger(basicComponent.Name)
}

// // OnStats2Trigger
// func (basicComponent *BasicComponent) OnStats2Trigger(s2 *Stats2) {

// }

// // GetSymbols -
// func (basicComponent *BasicComponent) GetSymbols(gameProp *GameProperty) []int {
// 	return nil
// }

// // AddSymbol -
// func (basicComponent *BasicComponent) AddSymbol(gameProp *GameProperty, symbol int) {

// }

// // OnEachSymbol - on foreach symbol
// func (basicComponent *BasicComponent) OnEachSymbol(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState,
// 	stake *sgc7game.Stake, prs []*sgc7game.PlayResult, symbol int, cd IComponentData) (string, error) {
// 	return "", nil
// }

// EachSymbols - foreach symbols
func (basicComponent *BasicComponent) EachSymbols(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
	prs []*sgc7game.PlayResult, cd IComponentData) error {
	return nil
}

// AddPos -
func (basicComponent *BasicComponent) AddPos(cd IComponentData, x int, y int) {

}

// // GetComponentData -
// func (basicComponent *BasicComponent) GetComponentData(gameProp *GameProperty) IComponentData {
// 	return gameProp.GetCurComponentData(basicComponent)

// 	if basicComponent.dataForeachSymbol != nil {
// 		return gameProp.MapComponentData[fmt.Sprintf("%v:%v", basicComponent.Name, basicComponent.dataForeachSymbol.Index)]
// 	}

// 	return gameProp.MapComponentData[basicComponent.Name]
// }

// // SetForeachSymbolData -
// func (basicComponent *BasicComponent) SetForeachSymbolData(data *ForeachSymbolData) {
// 	basicComponent.dataForeachSymbol = data
// }

// OnGameInited - on game inited
func (basicComponent *BasicComponent) OnGameInited(components *ComponentList) error {
	return nil
}

// GetAllLinkComponents - get all link components
func (basicComponent *BasicComponent) GetAllLinkComponents() []string {
	return []string{basicComponent.Config.DefaultNextComponent}
}

// GetNextLinkComponents - get next link components
func (basicComponent *BasicComponent) GetNextLinkComponents() []string {
	return []string{basicComponent.Config.DefaultNextComponent}
}

// GetChildLinkComponents - get child link components
func (basicComponent *BasicComponent) GetChildLinkComponents() []string {
	return nil
}

// CanTriggerWithScene -
func (basicComponent *BasicComponent) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	goutils.Error("BasicComponent.CanTriggerWithScene",
		goutils.Err(ErrInvalidComponent))

	return false, nil
}

// SetMask -
func (basicComponent *BasicComponent) SetMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error {
	goutils.Error("BasicComponent.SetMask",
		goutils.Err(ErrInvalidComponent))

	return ErrInvalidComponent
}

// SetMaskVal -
func (basicComponent *BasicComponent) SetMaskVal(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, index int, mask bool) error {
	goutils.Error("BasicComponent.SetMaskVal",
		goutils.Err(ErrInvalidComponent))

	return ErrInvalidComponent
}

// SetMaskOnlyTrue -
func (basicComponent *BasicComponent) SetMaskOnlyTrue(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd IComponentData, mask []bool) error {
	goutils.Error("BasicComponent.SetMaskOnlyTrue",
		goutils.Err(ErrInvalidComponent))

	return ErrInvalidComponent
}

func NewBasicComponent(name string, srcSceneNum int) *BasicComponent {
	return &BasicComponent{
		Name:        name,
		SrcSceneNum: srcSceneNum,
	}
}
