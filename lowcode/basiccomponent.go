package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/types/known/anypb"
)

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
func (basicComponent *BasicComponent) onStepEnd(_ *GameProperty, _ *sgc7game.PlayResult, _ *GameParams, nextComponent string) string {
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

	gameProp.SceneStack.Push(basicComponent.Name, sc)
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

	gameProp.OtherSceneStack.Push(basicComponent.Name, sc)
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

	// curpr.CashWin += int64(ret.CashWin)
	// curpr.CoinWin += ret.CoinWin

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
func (basicComponent *BasicComponent) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	s2.ProcStatsTrigger(basicComponent.Name)
}

// IsNeedOnStepEndStats2 - 除respin外，如果也有component也需要在stepEnd调用的话，这里需要返回true
func (basicComponent *BasicComponent) IsNeedOnStepEndStats2() bool {
	return false
}

// OnProcControllers -
func (basicComponent *BasicComponent) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {

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

// OnPlayGameWithSet - on playgame with a set
func (basicComponent *BasicComponent) OnPlayGameWithSet(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData, set int) (string, error) {
	goutils.Error("BasicComponent.OnPlayGameWithSet",
		goutils.Err(ErrInvalidSetComponent))

	return "", ErrInvalidSetComponent
}

// GetBranchNum -
func (basicComponent *BasicComponent) GetBranchNum() int {
	return 0
}

// GetBranchWeights -
func (basicComponent *BasicComponent) GetBranchWeights() []int {
	return nil
}

func NewBasicComponent(name string, srcSceneNum int) *BasicComponent {
	return &BasicComponent{
		Name:        name,
		SrcSceneNum: srcSceneNum,
	}
}
