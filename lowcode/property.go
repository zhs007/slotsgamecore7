package lowcode

import (
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"github.com/zhs007/slotsgamecore7/stats2"
)

const (
	GamePropWidth        = 1
	GamePropHeight       = 2
	GamePropCurPaytables = 3
	GamePropCurReels     = 4
	GamePropCurLineData  = 5
	GamePropCurLineNum   = 6
	GamePropCurBetIndex  = 7
)

var MapProperty map[string]int

func String2Property(str string) (int, error) {
	v, isok := MapProperty[str]
	if isok {
		return v, nil
	}

	goutils.Error("String2Property",
		slog.String("str", str),
		goutils.Err(ErrInvalidGamePropertyString))

	return 0, ErrInvalidGamePropertyString
}

type HistoryComponentData struct {
	Component    IComponent
	ForeachIndex int
}

type GameProperty struct {
	CurBetMul              int
	Pool                   *GamePropertyPool
	MapVals                map[int]int
	MapStrVals             map[int]string
	CurPaytables           *sgc7game.PayTables
	CurLineData            *sgc7game.LineData
	CurReels               *sgc7game.ReelsData
	MapIntValWeights       map[string]*sgc7game.ValWeights2
	MapStats               map[string]*sgc7stats.Feature
	mapInt                 map[string]int
	mapStr                 map[string]string
	mapGlobalStr           map[string]string
	mapGlobalScene         map[string]*sgc7game.GameScene // v0.13开始弃用
	mapComponentScene      map[string]*sgc7game.GameScene
	mapComponentOtherScene map[string]*sgc7game.GameScene
	callStack              *CallStack
	RespinComponents       []string
	PoolScene              *sgc7game.GameScenePoolEx
	Components             *ComponentList
	SceneStack             *SceneStack
	OtherSceneStack        *SceneStack
	stats2Cache            *stats2.Cache
	usedComponent          []string
	rng                    IRNG
	featureLevel           IFeatureLevel
}

func (gameProp *GameProperty) GetBetMul() int {
	return gameProp.CurBetMul
}

func (gameProp *GameProperty) UseComponent(componentName string) {
	if goutils.IndexOfStringSlice(gameProp.usedComponent, componentName, 0) < 0 {
		gameProp.usedComponent = append(gameProp.usedComponent, componentName)
	}
}

func (gameProp *GameProperty) BuildGameParam(gp *GameParams) {
	if len(gameProp.RespinComponents) > 0 {
		gp.RespinComponents = make([]string, len(gameProp.RespinComponents))

		copy(gp.RespinComponents, gameProp.RespinComponents)
	}

	gp.SetGameProp(gameProp)
}

func (gameProp *GameProperty) OnNewGame(stake *sgc7game.Stake, curPlugin sgc7plugin.IPlugin) error {
	curBet := stake.CashBet / stake.CoinBet
	for i, v := range gameProp.Pool.Config.Bets {
		if v == int(curBet) {
			gameProp.SetVal(GamePropCurBetIndex, i)

			break
		}
	}

	gameProp.mapGlobalStr = make(map[string]string)
	gameProp.mapGlobalScene = make(map[string]*sgc7game.GameScene)
	gameProp.mapComponentScene = make(map[string]*sgc7game.GameScene)
	gameProp.mapComponentOtherScene = make(map[string]*sgc7game.GameScene)

	gameProp.callStack.OnNewGame()

	// gameProp.rng = gameProp.newRNG()

	gameProp.rng.OnNewGame(curPlugin)
	// gameProp = nil

	return nil
}

func (gameProp *GameProperty) OnNewStep() error {
	gameProp.mapInt = make(map[string]int)
	gameProp.mapStr = make(map[string]string)

	gameProp.callStack.OnNewStep()

	return nil
}

func (gameProp *GameProperty) TagGlobalScene(tag string, gs *sgc7game.GameScene) {
	gameProp.mapGlobalScene[tag] = gs
}

func (gameProp *GameProperty) GetGlobalScene(tag string) *sgc7game.GameScene {
	gs, isok := gameProp.mapGlobalScene[tag]
	if !isok {
		return nil
	}

	return gs
}

func (gameProp *GameProperty) SetComponentScene(component string, gs *sgc7game.GameScene) {
	gameProp.mapComponentScene[component] = gs
}

func (gameProp *GameProperty) GetComponentScene(component string) *sgc7game.GameScene {
	gs, isok := gameProp.mapComponentScene[component]
	if !isok {
		return nil
	}

	return gs
}

func (gameProp *GameProperty) SetComponentOtherScene(component string, gs *sgc7game.GameScene) {
	gameProp.mapComponentOtherScene[component] = gs
}

func (gameProp *GameProperty) ClearComponentOtherScene(component string) {
	delete(gameProp.mapComponentOtherScene, component)
}

func (gameProp *GameProperty) GetComponentOtherScene(component string) *sgc7game.GameScene {
	gs, isok := gameProp.mapComponentOtherScene[component]
	if !isok {
		return nil
	}

	return gs
}

func (gameProp *GameProperty) TagScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return nil, -1
	}

	if si < len(pr.Scenes) {
		return pr.Scenes[si], si
	}

	return nil, -1
}

func (gameProp *GameProperty) TagOtherScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetOtherScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return nil, -1
	}

	if si < len(pr.OtherScenes) {
		return pr.OtherScenes[si], si
	}

	return nil, -1
}

func (gameProp *GameProperty) Respin(pr *sgc7game.PlayResult, gp *GameParams, respinComponent string, gs *sgc7game.GameScene, os *sgc7game.GameScene) {
	if gs != nil {
		gp.LastScene = gs.CloneEx(gameProp.PoolScene)
	}

	if os != nil {
		gp.LastOtherScene = os.CloneEx(gameProp.PoolScene)
	}
}

func (gameProp *GameProperty) hasRespin(respinComponent string) bool {
	if len(gameProp.RespinComponents) == 0 {
		return false
	}

	return goutils.IndexOfStringSlice(gameProp.RespinComponents, respinComponent, 0) >= 0
}

// onTriggerRespin -
func (gameProp *GameProperty) onTriggerRespin(respinComponent string) error {
	// 暂时不考虑respin的嵌套，respin的嵌套如果要处理，也需要callstack那个层面来处理
	if !gameProp.hasRespin(respinComponent) {
		gameProp.RespinComponents = append(gameProp.RespinComponents, respinComponent)
	}

	return nil
}

func (gameProp *GameProperty) removeRespin(respinComponent string) error {
	for i, respin := range gameProp.RespinComponents {
		if respin == respinComponent {
			gameProp.RespinComponents = append(gameProp.RespinComponents[:i], gameProp.RespinComponents[i+1:]...)
		}
	}

	return nil
}

func (gameProp *GameProperty) procEndingRespin() {
	if len(gameProp.RespinComponents) > 0 {
		ei := len(gameProp.RespinComponents) - 1
		for i := len(gameProp.RespinComponents) - 1; i >= 0; i-- {
			if gameProp.IsEndingRespin(gameProp.RespinComponents[i]) {
				ei--
			} else {
				break
			}
		}

		if ei >= 0 {
			gameProp.RespinComponents = gameProp.RespinComponents[:ei+1]
		}
	}
}

func (gameProp *GameProperty) ProcRespin(pr *sgc7game.PlayResult, gp *GameParams) {
	gameProp.procEndingRespin()

	if len(gameProp.RespinComponents) > 0 {
		gp.NextStepFirstComponent = gameProp.RespinComponents[len(gameProp.RespinComponents)-1]

		pr.IsFinish = false

		for _, v := range gameProp.usedComponent {
			if goutils.IndexOfStringSlice(gp.HistoryComponents, v, 0) < 0 {
				cd := gameProp.GetGlobalComponentDataWithName(v)
				gp.AddComponentData(v, cd)
			}
		}
	} else if !pr.IsWait {
		pr.IsFinish = true
	}
}

// procRespinBeforeStepEnding - 这里用来处理当前respin结束后，继续next的流程
func (gameProp *GameProperty) procRespinBeforeStepEnding(pr *sgc7game.PlayResult, gp *GameParams) (string, error) {
	if len(gameProp.RespinComponents) > 0 {
		nextComponent := ""
		for i := len(gameProp.RespinComponents) - 1; i >= 0; i-- {
			curRespin := gameProp.RespinComponents[i]

			cr, isok := gameProp.Components.MapComponents[curRespin]
			if isok {
				cd := gameProp.GetGlobalComponentData(cr)
				nc, err := cr.ProcRespinOnStepEnd(gameProp, pr, gp, cd, nextComponent == "")
				if err != nil {
					goutils.Error("GameProperty.procRespinBeforeStepEnding:ProcRespinOnStepEnd",
						slog.String("respin", curRespin),
						goutils.Err(err))

					return "", err
				}

				if nextComponent == "" {
					nextComponent = nc
				}
			}
		}

		return nextComponent, nil
	}

	return "", nil
}

func (gameProp *GameProperty) OnCallEnd(component IComponent, cd IComponentData, gp *GameParams) {
	if !component.IsRespin() && gAllowStats2 {
		if !gameProp.stats2Cache.HasFeature(component.GetName()) {
			gameProp.stats2Cache.AddFeature(component.GetName(), component.NewStats2(gameProp.Components.statsNodeData.GetParent(component.GetName())))
		}

		component.OnStats2(cd, gameProp.stats2Cache)
	}

	tag := gameProp.callStack.OnCallEnd(component, cd)
	gp.HistoryComponents = append(gp.HistoryComponents, tag)
}

func (gameProp *GameProperty) TriggerRespin(plugin sgc7plugin.IPlugin, pr *sgc7game.PlayResult, gp *GameParams, respinNum int, respinComponent string, usePushTrigger bool) error {
	component, isok := gameProp.Components.MapComponents[respinComponent]
	if isok {
		gameProp.UseComponent(respinComponent)

		cd := gameProp.GetGlobalComponentData(component)
		if usePushTrigger {
			cd.PushTriggerRespin(gameProp, plugin, pr, gp, respinNum)
		} else {
			cd.AddRespinTimes(respinNum)
		}
	}

	return nil
}

func (gameProp *GameProperty) TriggerRespinWithWeights(pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, fn string, useFileMapping bool, respinComponent string, usePushTrigger bool) (int, error) {
	vw2, err := gameProp.GetIntValWeights(fn, useFileMapping)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:GetIntValWeights",
			slog.String("fn", fn),
			goutils.Err(err))

		return 0, err
	}

	val, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:RandVal",
			slog.String("fn", fn),
			goutils.Err(err))

		return 0, err
	}

	if val.Int() > 0 {
		gameProp.TriggerRespin(plugin, pr, gp, val.Int(), respinComponent, usePushTrigger)

		return val.Int(), nil
	}

	return 0, nil
}

func (gameProp *GameProperty) GetIntValWeights(fn string, useFileMapping bool) (*sgc7game.ValWeights2, error) {
	vw2, isok := gameProp.MapIntValWeights[fn]
	if !isok {
		curvw2, err := gameProp.Pool.LoadIntWeights(fn, useFileMapping)
		if err != nil {
			goutils.Error("GameProperty.GetIntValWeights:LoadSymbolWeights",
				slog.String("Weight", fn),
				goutils.Err(err))

			return nil, err
		}

		gameProp.MapIntValWeights[fn] = curvw2

		vw2 = curvw2
	}

	return vw2, nil
}

func (gameProp *GameProperty) SetVal(prop int, val int) error {
	gameProp.MapVals[prop] = val

	return nil
}

func (gameProp *GameProperty) AddVal(prop int, val int) error {
	gameProp.MapVals[prop] += val

	return nil
}

func (gameProp *GameProperty) GetVal(prop int) int {
	return gameProp.MapVals[prop]
}

func (gameProp *GameProperty) SetStrVal(prop int, val string) error {
	if prop == GamePropCurPaytables {
		v, isok := gameProp.Pool.Config.MapPaytables[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurPaytables",
				slog.String("val", val),
				goutils.Err(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurPaytables = v
	} else if prop == GamePropCurLineData {
		v, isok := gameProp.Pool.Config.MapLinedate[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurLineData",
				slog.String("val", val),
				goutils.Err(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurLineData = v

		gameProp.SetVal(GamePropCurLineNum, len(gameProp.CurLineData.Lines))
	}

	gameProp.MapStrVals[prop] = val

	return nil
}

func (gameProp *GameProperty) GetStrVal(prop int) string {
	return gameProp.MapStrVals[prop]
}

func (gameProp *GameProperty) TagInt(tag string, val int) {
	gameProp.mapInt[tag] = val
}

func (gameProp *GameProperty) GetTagInt(tag string) int {
	return gameProp.mapInt[tag]
}

func (gameProp *GameProperty) TagStr(tag string, val string) {
	gameProp.mapStr[tag] = val
}

func (gameProp *GameProperty) GetTagStr(tag string) string {
	return gameProp.mapStr[tag]
}

func (gameProp *GameProperty) TagGlobalStr(tag string, val string) {
	gameProp.mapGlobalStr[tag] = val
}

func (gameProp *GameProperty) GetTagGlobalStr(tag string) string {
	return gameProp.mapGlobalStr[tag]
}

func (gameProp *GameProperty) IsRespin(componentName string) bool {
	component, isok := gameProp.Components.MapComponents[componentName]
	if isok {
		return component.IsRespin()
	}

	return false
}

func (gameProp *GameProperty) GetComponentVal(componentVal string) (int, error) {
	arr := strings.Split(componentVal, ".")
	if len(arr) != 2 {
		goutils.Error("GameProperty.GetComponentVal",
			slog.String("componentVal", componentVal),
			goutils.Err(ErrInvalidComponentVal))

		return 0, ErrInvalidComponentVal
	}

	component, isok := gameProp.Components.MapComponents[arr[0]]
	if !isok {
		goutils.Error("GameProperty.GetComponentVal",
			slog.String("component", arr[0]),
			goutils.Err(ErrInvalidComponent))

		return 0, ErrInvalidComponent
	}

	cd := gameProp.callStack.GetComponentData(gameProp, component)

	v, isok := cd.GetVal(arr[1])
	if !isok {
		goutils.Error("GameProperty.GetComponentVal:GetVal",
			slog.String("componentVal", componentVal),
			goutils.Err(ErrInvalidComponentVal))

		return 0, ErrInvalidComponentVal
	}

	return v, nil
}

func (gameProp *GameProperty) GetComponentVal2(component string, val string) (int, error) {
	ic, isok := gameProp.Components.MapComponents[component]
	if !isok {
		goutils.Error("GameProperty.GetComponentVal",
			slog.String("component", component),
			goutils.Err(ErrInvalidComponent))

		return 0, ErrInvalidComponent
	}

	cd := gameProp.callStack.GetComponentData(gameProp, ic)

	v, isok := cd.GetVal(val)
	if !isok {
		goutils.Error("GameProperty.GetComponentVal:GetVal",
			slog.String("component", component),
			slog.String("val", val),
			goutils.Err(ErrInvalidComponentVal))

		return 0, ErrInvalidComponentVal
	}

	return v, nil
}

func (gameProp *GameProperty) procAwards(plugin sgc7plugin.IPlugin, awards []*Award, curpr *sgc7game.PlayResult, gp *GameParams) {
	for _, v := range awards {
		gameProp.procAward(plugin, v, curpr, gp, false)
	}
}

func (gameProp *GameProperty) procAward(plugin sgc7plugin.IPlugin, award *Award, curpr *sgc7game.PlayResult, gp *GameParams, skipTriggerRespin bool) {
	if !skipTriggerRespin && award.OnTriggerRespin != "" {
		component, isok := gameProp.Components.MapComponents[award.OnTriggerRespin]
		if !isok {
			goutils.Error("GameProperty.procAward:OnTriggerRespin",
				goutils.Err(ErrInvalidComponent))

			return
		}

		if !component.IsRespin() {
			goutils.Error("GameProperty.procAward:OnTriggerRespin:IsRespin",
				goutils.Err(ErrNotRespin))

			return
		}

		cd := gameProp.GetGlobalComponentData(component)
		cd.AddTriggerRespinAward(award)

		return
	}

	if award.Type == AwardRespinTimes {
		component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		if isok {
			cd := gameProp.GetGlobalComponentData(component)
			gameProp.UseComponent(award.StrParams[0])

			cd.AddRespinTimes(award.Vals[0])
		}
	} else if award.Type == AwardTriggerRespin {
		gameProp.TriggerRespin(plugin, curpr, gp, award.Vals[0], award.StrParams[0], false)
		// component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		// if isok {
		// 	cd := gameProp.GetGlobalComponentData(component)
		// 	cd.TriggerRespin(gameProp, plugin, curpr, gp)
		// }
	} else if award.Type == AwardAddRetriggerRespinNum {
		component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		if isok {
			cd := gameProp.GetGlobalComponentData(component)
			cd.ChgConfigIntVal(CCVRetriggerRespinNum, award.Vals[0])
		}
	} else if award.Type == AwardSetMaskVal {
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.Pool.SetMaskVal(plugin, gameProp, curpr, gp, award.StrParams[0], award.Vals[0], award.Vals[1] != 0)
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetMaskVal:SetMaskVal",
				goutils.Err(err))

			return
		}
	} else if award.Type == AwardTriggerRespin2 {
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, award.StrParams[0], award.GetVal(gameProp, 0))
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardTriggerRespin2:PushTrigger",
				goutils.Err(err))

			return
		}
	} else if award.Type == AwardSetComponentConfigVal {
		err := gameProp.SetComponentConfigVal(award.StrParams[0], award.StrParams[1])
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:SetComponentConfigVal",
				goutils.Err(err))

			return
		}
	} else if award.Type == AwardSetComponentConfigIntVal {
		err := gameProp.SetComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0))
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:AwardSetComponentConfigIntVal",
				goutils.Err(err))

			return
		}
	} else if award.Type == AwardChgComponentConfigIntVal {
		err := gameProp.ChgComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0))
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:AwardChgComponentConfigIntVal",
				goutils.Err(err))

			return
		}
	}
}

func (gameProp *GameProperty) GetBet2(stake *sgc7game.Stake, bt BetType) int {
	if bt == BTypeTotalBet {
		return int(stake.CoinBet) * gameProp.Pool.Config.TotalBetInWins[gameProp.GetVal(GamePropCurBetIndex)]
	} else if bt == BTypeBet {
		return int(stake.CoinBet)
	}

	return 0
}

func (gameProp *GameProperty) GetLastRespinNum(respinComponent string) int {
	component, isok := gameProp.Components.MapComponents[respinComponent]
	if isok {
		cd := gameProp.GetGlobalComponentData(component)

		return cd.GetLastRespinNum()
	}

	return 0
}

// CanTrigger -
func (gameProp *GameProperty) CanTrigger(componentName string, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	component, isok := gameProp.Components.MapComponents[componentName]
	if isok {
		isTrigger, _ := component.CanTriggerWithScene(gameProp, gs, curpr, stake)

		return isTrigger
	}

	return false
}

func (gameProp *GameProperty) IsInCurCallStack(componentName string) bool {
	return gameProp.callStack.IsInCurCallStack(componentName)
}

func (gameProp *GameProperty) IsEndingRespin(componentName string) bool {
	component, isok := gameProp.Components.MapComponents[componentName]
	if isok {
		cd := gameProp.GetGlobalComponentData(component)
		return cd.IsRespinEnding()
	}

	return false
}

func (gameProp *GameProperty) IsStartedRespin(componentName string) bool {
	component, isok := gameProp.Components.MapComponents[componentName]
	if isok {
		cd := gameProp.GetGlobalComponentData(component)
		return cd.IsRespinStarted()
	}

	return false
}

func (gameProp *GameProperty) SetComponentConfigVal(componentConfigValName string, val string) error {
	arr := strings.Split(componentConfigValName, ".")
	if len(arr) != 2 {
		goutils.Error("GameProperty.SetComponentConfigVal",
			slog.String("componentConfigValName", componentConfigValName),
			goutils.Err(ErrInvalidComponentVal))

		return ErrInvalidComponentVal
	}

	cd := gameProp.GetCurComponentDataWithName(arr[0])
	if cd == nil {
		goutils.Error("GameProperty.SetComponentConfigVal",
			slog.String("componentConfigVal", arr[0]),
			goutils.Err(ErrInvalidComponent))

		return ErrInvalidComponent
	}

	gameProp.UseComponent(arr[0])

	cd.SetConfigVal(arr[1], val)

	return nil
}

func (gameProp *GameProperty) SetComponentConfigIntVal(componentConfigValName string, val int) error {
	arr := strings.Split(componentConfigValName, ".")
	if len(arr) != 2 {
		goutils.Error("GameProperty.SetComponentConfigIntVal",
			slog.String("componentConfigValName", componentConfigValName),
			goutils.Err(ErrInvalidComponentVal))

		return ErrInvalidComponentVal
	}

	cd := gameProp.GetCurComponentDataWithName(arr[0])
	if cd == nil {
		goutils.Error("GameProperty.SetComponentConfigIntVal",
			slog.String("componentConfigVal", arr[0]),
			goutils.Err(ErrInvalidComponent))

		return ErrInvalidComponent
	}

	gameProp.UseComponent(arr[0])

	cd.SetConfigIntVal(arr[1], val)

	return nil
}

func (gameProp *GameProperty) ChgComponentConfigIntVal(componentConfigValName string, off int) error {
	arr := strings.Split(componentConfigValName, ".")
	if len(arr) != 2 {
		goutils.Error("GameProperty.SetComponentConfigIntVal",
			slog.String("componentConfigValName", componentConfigValName),
			goutils.Err(ErrInvalidComponentVal))

		return ErrInvalidComponentVal
	}

	cd := gameProp.GetCurComponentDataWithName(arr[0])
	if cd == nil {
		goutils.Error("GameProperty.ChgComponentConfigIntVal",
			slog.String("componentConfigVal", arr[0]),
			goutils.Err(ErrInvalidComponent))

		return ErrInvalidComponent
	}

	gameProp.UseComponent(arr[0])

	cd.ChgConfigIntVal(arr[1], off)

	return nil
}

func (gameProp *GameProperty) GetComponentSymbols(componentName string) []int {
	cd := gameProp.GetCurComponentDataWithName(componentName)
	if cd == nil {
		goutils.Error("GameProperty.GetComponentSymbols",
			slog.String("componentConfigVal", componentName),
			goutils.Err(ErrInvalidComponent))

		return nil
	}

	return cd.GetSymbols()
}

func (gameProp *GameProperty) GetComponentPos(componentName string) []int {
	cd := gameProp.GetCurComponentDataWithName(componentName)
	if cd == nil {
		goutils.Error("GameProperty.GetComponentPos",
			slog.String("componentConfigVal", componentName),
			goutils.Err(ErrInvalidComponent))

		return nil
	}

	return cd.GetPos()
}

func (gameProp *GameProperty) AddComponentSymbol(componentName string, symbolCode int) {
	cd := gameProp.GetCurComponentDataWithName(componentName)
	if cd == nil {
		goutils.Error("GameProperty.AddComponentSymbol",
			slog.String("componentConfigVal", componentName),
			goutils.Err(ErrInvalidComponent))

		return
	}

	gameProp.UseComponent(componentName)

	cd.AddSymbol(symbolCode)
}

func (gameProp *GameProperty) AddComponentPos(componentName string, pos []int) {
	cd := gameProp.GetCurComponentDataWithName(componentName)
	if cd == nil {
		goutils.Error("GameProperty.AddComponentPos",
			slog.String("componentConfigVal", componentName),
			goutils.Err(ErrInvalidComponent))

		return
	}

	gameProp.UseComponent(componentName)

	for i := 0; i < len(pos)/2; i++ {
		cd.AddPos(pos[i*2], pos[i*2+1])
	}
}

func (gameProp *GameProperty) ForceComponentBranch(componentName string, branchIndex int) {
	if gIsReleaseMode {
		return
	}

	cd := gameProp.GetCurComponentDataWithName(componentName)
	if cd == nil {
		goutils.Error("GameProperty.ForceComponentBranch",
			slog.String("componentConfigVal", componentName),
			goutils.Err(ErrInvalidComponent))

		return
	}

	cd.ForceBranch(branchIndex)
}

func (gameProp *GameProperty) onStepEnd(gp *GameParams, pr *sgc7game.PlayResult, _ []*sgc7game.PlayResult) {
	gameProp.featureLevel.OnStepEnd(gameProp, gp, pr)

	if gAllowStats2 {
		for _, v := range gp.HistoryComponents {
			ic, isok := gameProp.Components.MapComponents[v]
			if isok && ic.IsRespin() {
				if !gameProp.stats2Cache.HasFeature(v) {
					gameProp.stats2Cache.AddFeature(v, ic.NewStats2(gameProp.Components.statsNodeData.GetParent(v)))
				}

				ic.OnStats2(gameProp.GetComponentData(ic), gameProp.stats2Cache)
			}
		}
	}

	if pr.IsFinish {
		gameProp.PoolScene.Reset()
	}
}

func (gameProp *GameProperty) GetMask(name string) ([]bool, error) {
	ic, isok := gameProp.Components.MapComponents[name]
	if !isok || !ic.IsMask() {
		goutils.Error("GameProperty.GetMask",
			slog.String("name", name),
			goutils.Err(ErrInvalidComponentName))

		return nil, ErrInvalidComponentName
	}

	cd := gameProp.GetComponentData(ic)
	mask := cd.GetMask()

	return mask, nil
}

func (gameProp *GameProperty) GetCurComponentData(ic IComponent) IComponentData {
	return gameProp.callStack.GetCurComponentData(gameProp, ic)
}

func (gameProp *GameProperty) GetCurComponentDataWithName(componentName string) IComponentData {
	ic := gameProp.Components.MapComponents[componentName]

	if ic != nil {
		return gameProp.callStack.GetCurComponentData(gameProp, ic)
	}

	return nil
}

func (gameProp *GameProperty) GetGlobalComponentDataWithName(componentName string) IComponentData {
	ic := gameProp.Components.MapComponents[componentName]

	if ic != nil {
		return gameProp.callStack.GetGlobalComponentData(gameProp, ic)
	}

	return nil
}

func (gameProp *GameProperty) GetGlobalComponentData(icomponent IComponent) IComponentData {
	return gameProp.callStack.GetGlobalComponentData(gameProp, icomponent)
}

func (gameProp *GameProperty) GetComponentDataWithName(componentName string) IComponentData {
	ic := gameProp.Components.MapComponents[componentName]

	if ic != nil {
		return gameProp.callStack.GetComponentData(gameProp, ic)
	}

	return nil
}

func (gameProp *GameProperty) GetComponentData(icomponent IComponent) IComponentData {
	return gameProp.callStack.GetComponentData(gameProp, icomponent)
}

func (gameProp *GameProperty) GetCurCallStackSymbol() int {
	return gameProp.callStack.GetCurCallStackSymbol()
}

func init() {
	MapProperty = make(map[string]int)

	MapProperty["width"] = GamePropWidth
	MapProperty["height"] = GamePropHeight
	MapProperty["paytables"] = GamePropCurPaytables
	MapProperty["reels"] = GamePropCurReels
	MapProperty["linedata"] = GamePropCurLineData
}
