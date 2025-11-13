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

type SPGridStack struct {
	Width  int
	Height int
	Stack  *SceneStack
}

type GameProperty struct {
	CurBetMul                        int
	Pool                             *GamePropertyPool
	MapVals                          map[int]int
	MapStrVals                       map[int]string
	CurPaytables                     *sgc7game.PayTables
	CurLineData                      *sgc7game.LineData
	CurReels                         *sgc7game.ReelsData
	MapIntValWeights                 map[string]*sgc7game.ValWeights2
	MapStats                         map[string]*sgc7stats.Feature
	mapInt                           map[string]int
	mapStr                           map[string]string
	mapGlobalStr                     map[string]string
	mapGlobalScene                   map[string]*sgc7game.GameScene // v0.13开始弃用
	mapComponentScene                map[string]*sgc7game.GameScene
	mapComponentOtherScene           map[string]*sgc7game.GameScene
	callStack                        *CallStack
	RespinComponents                 []string
	PoolScene                        *sgc7game.GameScenePoolEx
	Components                       *ComponentList
	SceneStack                       *SceneStack
	OtherSceneStack                  *SceneStack
	stats2Cache                      *stats2.Cache
	usedComponent                    []string
	rng                              IRNG
	featureLevel                     IFeatureLevel
	lstNeedOnStepEndStats2Components []string
	MapSPGridStack                   map[string]*SPGridStack
	posPool                          *PosPool
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

	gameProp.usedComponent = nil
	gameProp.lstNeedOnStepEndStats2Components = nil

	// gameProp.rng = gameProp.newRNG()

	gameProp.rng.OnNewGame(int(curBet), curPlugin)
	// gameProp = nil

	return nil
}

func (gameProp *GameProperty) OnNewStep() error {
	gameProp.mapInt = make(map[string]int)
	gameProp.mapStr = make(map[string]string)

	// gameProp.callStack = gameProp.callStack.OnNewStep()
	gameProp.callStack.OnNewStep()

	gameProp.usedComponent = nil

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

		if ei < 0 {
			gameProp.RespinComponents = nil
		} else if ei >= 0 && ei < len(gameProp.RespinComponents)-1 {
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
		//! 这里不能全部遍历完，因为有可能在 next 分支里，会产生外层 respin的增加，这里只能处理当前这一级
		ci := len(gameProp.RespinComponents) - 1
		curRespin := gameProp.RespinComponents[ci]

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

		return nextComponent, nil
	}

	return "", nil
}

// OnCallEnd - call after the component onPlay
func (gameProp *GameProperty) OnCallEnd(component IComponent, cd IComponentData, gp *GameParams, pr *sgc7game.PlayResult) {
	if !component.IsRespin() && gAllowStats2 {
		if !gameProp.stats2Cache.HasFeature(component.GetName()) {
			gameProp.stats2Cache.AddFeature(component.GetName(),
				component.NewStats2(gameProp.Components.statsNodeData.GetParent(component.GetName())),
				false)
		}

		componentName := component.GetName()
		if component.IsNeedOnStepEndStats2() &&
			goutils.IndexOfStringSlice(gameProp.lstNeedOnStepEndStats2Components, componentName, 0) < 0 {

			gameProp.lstNeedOnStepEndStats2Components = append(gameProp.lstNeedOnStepEndStats2Components, componentName)
		} else {
			component.OnStats2(cd, gameProp.stats2Cache, gameProp, gp, pr, false)
		}
	}

	tag := gameProp.callStack.OnCallEnd(component, cd)

	gp.HistoryComponents = append(gp.HistoryComponents, tag)

	if gAllowFullComponentHistory {
		gp.HistoryComponentsEx = append(gp.HistoryComponentsEx, tag)
	}
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
	switch prop {
	case GamePropCurPaytables:
		v, isok := gameProp.Pool.Config.MapPaytables[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurPaytables",
				slog.String("val", val),
				goutils.Err(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurPaytables = v
	case GamePropCurLineData:
		v, isok := gameProp.Pool.Config.MapLinedate[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurLineData",
				slog.String("val", val),
				goutils.Err(ErrInvalidLineData))

			return ErrInvalidLineData
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

	v, isok := cd.GetValEx(strings.ToLower(arr[1]), GCVTypeNormal)
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

	v, isok := cd.GetValEx(strings.ToLower(val), GCVTypeNormal)
	if !isok {
		goutils.Error("GameProperty.GetComponentVal:GetVal",
			slog.String("component", component),
			slog.String("val", val),
			goutils.Err(ErrInvalidComponentVal))

		return 0, ErrInvalidComponentVal
	}

	return v, nil
}

func (gameProp *GameProperty) GetComponentStrVal(componentVal string) (string, error) {
	arr := strings.Split(componentVal, ".")
	if len(arr) != 2 {
		goutils.Error("GameProperty.GetComponentVal",
			slog.String("componentVal", componentVal),
			goutils.Err(ErrInvalidComponentVal))

		return "", ErrInvalidComponentVal
	}

	component, isok := gameProp.Components.MapComponents[arr[0]]
	if !isok {
		goutils.Error("GameProperty.GetComponentVal",
			slog.String("component", arr[0]),
			goutils.Err(ErrInvalidComponent))

		return "", ErrInvalidComponent
	}

	cd := gameProp.callStack.GetComponentData(gameProp, component)

	v, isok := cd.GetStrVal(strings.ToLower(arr[1]))
	if !isok {
		goutils.Error("GameProperty.GetComponentVal:GetStrVal",
			slog.String("componentVal", componentVal),
			goutils.Err(ErrInvalidComponentVal))

		return "", ErrInvalidComponentVal
	}

	return v, nil
}

func (gameProp *GameProperty) GetComponentStrVal2(component string, val string) (string, error) {
	ic, isok := gameProp.Components.MapComponents[component]
	if !isok {
		goutils.Error("GameProperty.GetComponentStrVal2",
			slog.String("component", component),
			goutils.Err(ErrInvalidComponent))

		return "", ErrInvalidComponent
	}

	cd := gameProp.callStack.GetComponentData(gameProp, ic)

	v, isok := cd.GetStrVal(val)
	if !isok {
		goutils.Error("GameProperty.GetComponentStrVal2:GetVal",
			slog.String("component", component),
			slog.String("val", val),
			goutils.Err(ErrInvalidComponentVal))

		return "", ErrInvalidComponentVal
	}

	return v, nil
}

func (gameProp *GameProperty) procAwards(plugin sgc7plugin.IPlugin, awards []*Award, curpr *sgc7game.PlayResult, gp *GameParams) {
	for _, v := range awards {
		gameProp.procAward(plugin, v, curpr, gp, false)
	}
}

func (gameProp *GameProperty) RunController(award *Award, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams) {
	switch award.Type {
	case AwardSetComponentConfigVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.SetComponentConfigVal(award.StrParams[0], award.StrParams[1])
		if err != nil {
			goutils.Error("GameProperty.RunController:AwardSetComponentConfigVal:SetComponentConfigVal",
				goutils.Err(err))

			return
		}
	case AwardSetComponentConfigIntVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.SetComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0), func(componentName string, valName string, val int) bool {
			if valName == CCVLastTriggerNum {
				err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, componentName, award.GetVal(gameProp, 0))
				if err != nil {
					goutils.Error("GameProperty.RunController:AwardSetComponentConfigIntVal:PushTrigger",
						goutils.Err(err))
				}

				return true
			}

			return false
		}, func(componentName string, valName string, val int) {
			switch valName {
			case CCVForceValNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ProcControllers(gameProp, plugin, curpr, gp, val, "")
				}
			case CCVClearNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ClearData(gameProp.GetComponentDataWithName(componentName), true)
				}
			case CCVValueNumNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.OnUpdateDataWithPlayerState(gameProp.Pool, gameProp, plugin, curpr, gp, gp.ps, int(gp.stake.CashBet/gp.stake.CoinBet), int(gp.stake.CoinBet), gameProp.GetComponentDataWithName(componentName))
				}
			}
		})
		if err != nil {
			goutils.Error("GameProperty.RunController:AwardSetComponentConfigVal:AwardSetComponentConfigIntVal",
				goutils.Err(err))

			return
		}
	case AwardChgComponentConfigIntVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.ChgComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0), func(componentName string, valName string, off int) bool {
			if valName == CCVLastTriggerNum {
				err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, componentName, award.GetVal(gameProp, 0))
				if err != nil {
					goutils.Error("GameProperty.RunController:ChgComponentConfigIntVal:PushTrigger",
						goutils.Err(err))

				}

				return true
			}

			return false
		}, func(componentName string, valName string, off int, val int) {
			switch valName {
			case CCVForceValNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ProcControllers(gameProp, plugin, curpr, gp, val, "")
				}
			case CCVClearNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ClearData(gameProp.GetComponentDataWithName(componentName), true)
				}
			case CCVValueNumNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.OnUpdateDataWithPlayerState(gameProp.Pool, gameProp, plugin, curpr, gp, gp.ps, int(gp.stake.CashBet/gp.stake.CoinBet), int(gp.stake.CoinBet), gameProp.GetComponentDataWithName(componentName))
				}
			}
		})
		if err != nil {
			goutils.Error("GameProperty.RunController:AwardSetComponentConfigVal:AwardChgComponentConfigIntVal",
				goutils.Err(err))

			return
		}
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

	switch award.Type {
	case AwardRespinTimes:
		component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		if isok {
			cd := gameProp.GetGlobalComponentData(component)
			gameProp.UseComponent(award.StrParams[0])

			cd.AddRespinTimes(award.Vals[0])
		}
	case AwardTriggerRespin:
		gameProp.TriggerRespin(plugin, curpr, gp, award.Vals[0], award.StrParams[0], false)
		// component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		// if isok {
		// 	cd := gameProp.GetGlobalComponentData(component)
		// 	cd.TriggerRespin(gameProp, plugin, curpr, gp)
		// }
	case AwardAddRetriggerRespinNum:
		component, isok := gameProp.Components.MapComponents[award.StrParams[0]]
		if isok {
			cd := gameProp.GetGlobalComponentData(component)
			cd.ChgConfigIntVal(CCVRetriggerRespinNum, award.Vals[0])
		}
	case AwardSetMaskVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.Pool.SetMaskVal(plugin, gameProp, curpr, gp, award.StrParams[0], award.Vals[0], award.Vals[1] != 0)
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetMaskVal:SetMaskVal",
				goutils.Err(err))

			return
		}
	case AwardTriggerRespin2:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, award.StrParams[0], award.GetVal(gameProp, 0))
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardTriggerRespin2:PushTrigger",
				goutils.Err(err))

			return
		}
	case AwardSetComponentConfigVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.SetComponentConfigVal(award.StrParams[0], award.GetStringVal(gameProp, 0))
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:SetComponentConfigVal",
				goutils.Err(err))

			return
		}
	case AwardSetComponentConfigIntVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.SetComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0), func(componentName string, valName string, val int) bool {
			if valName == CCVLastTriggerNum {
				err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, componentName, award.GetVal(gameProp, 0))
				if err != nil {
					goutils.Error("GameProperty.procAward:AwardSetComponentConfigIntVal:PushTrigger",
						goutils.Err(err))

				}

				return true
			}

			return false
		}, func(componentName string, valName string, val int) {
			switch valName {
			case CCVForceValNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ProcControllers(gameProp, plugin, curpr, gp, val, "")
				}
			case CCVClearNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ClearData(gameProp.GetComponentDataWithName(componentName), true)
				}
			case CCVValueNumNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.OnUpdateDataWithPlayerState(gameProp.Pool, gameProp, plugin, curpr, gp, gp.ps, int(gp.stake.CashBet/gp.stake.CoinBet), int(gp.stake.CoinBet), gameProp.GetComponentDataWithName(componentName))
				}
			}
		})
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:AwardSetComponentConfigIntVal",
				goutils.Err(err))

			return
		}
	case AwardChgComponentConfigIntVal:
		gameProp.UseComponent(award.StrParams[0])
		err := gameProp.ChgComponentConfigIntVal(award.StrParams[0], award.GetVal(gameProp, 0), func(componentName string, valName string, off int) bool {
			if valName == CCVLastTriggerNum {
				err := gameProp.Pool.PushTrigger(gameProp, plugin, curpr, gp, componentName, 0)
				if err != nil {
					goutils.Error("GameProperty.procAward:AwardChgComponentConfigIntVal:PushTrigger",
						goutils.Err(err))

				}

				return true
			}

			return false
		}, func(componentName string, valName string, off int, val int) {
			switch valName {
			case CCVForceValNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ProcControllers(gameProp, plugin, curpr, gp, val, "")
				}
			case CCVClearNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.ClearData(gameProp.GetComponentDataWithName(componentName), true)
				}
			case CCVValueNumNow:
				component := gameProp.Components.MapComponents[componentName]
				if component != nil {
					component.OnUpdateDataWithPlayerState(gameProp.Pool, gameProp, plugin, curpr, gp, gp.ps, int(gp.stake.CashBet/gp.stake.CoinBet), int(gp.stake.CoinBet), gameProp.GetComponentDataWithName(componentName))
				}
			}
		})
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardSetComponentConfigVal:AwardChgComponentConfigIntVal",
				goutils.Err(err))

			return
		}
	}
}

func (gameProp *GameProperty) GetBet2(stake *sgc7game.Stake, bt BetType) int {
	switch bt {
	case BTypeTotalBet:
		return int(stake.CoinBet) * gameProp.Pool.Config.TotalBetInWins[gameProp.GetVal(GamePropCurBetIndex)]
	case BTypeBet:
		return int(stake.CoinBet)
	}

	return 0
}

func (gameProp *GameProperty) GetBet3(stake *sgc7game.Stake, bt BetType) int {
	if bt == BTypeTotalBet || bt == BTypeBet {
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
		icd := gameProp.GetComponentData(component)
		isTrigger, _ := component.CanTriggerWithScene(gameProp, gs, curpr, stake, icd)

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
	if len(arr) < 2 {
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

	str := arr[1]
	if len(arr) > 2 {
		for i := 2; i < len(arr); i++ {
			str += "." + arr[i]
		}
	}

	str = strings.ToLower(str)

	cd.SetConfigVal(str, val)

	return nil
}

func (gameProp *GameProperty) SetComponentConfigIntVal(componentConfigValName string, val int, onProc FuncOnChgComponentIntVal, onProced FuncOnSettedComponentIntVal) error {
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

	arr[1] = strings.ToLower(arr[1])

	if onProc != nil {
		if onProc(arr[0], arr[1], val) {
			return nil
		}
	}

	cd.SetConfigIntVal(arr[1], val)

	if onProced != nil {
		onProced(arr[0], arr[1], val)
	}

	return nil
}

func (gameProp *GameProperty) ChgComponentConfigIntVal(componentConfigValName string, off int, onProc FuncOnChgComponentIntVal, onProced FuncOnChgedComponentIntVal) error {
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

	arr[1] = strings.ToLower(arr[1])

	if onProc != nil {
		if onProc(arr[0], arr[1], off) {
			return nil
		}
	}

	nval := cd.ChgConfigIntVal(arr[1], off)

	if onProced != nil {
		onProced(arr[0], arr[1], off, nval)
	}

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

// func (gameProp *GameProperty) ForceComponentBranch(componentName string, branchIndex int) {
// 	if gIsReleaseMode {
// 		return
// 	}

// 	cd := gameProp.GetCurComponentDataWithName(componentName)
// 	if cd == nil {
// 		goutils.Error("GameProperty.ForceComponentBranch",
// 			slog.String("componentConfigVal", componentName),
// 			goutils.Err(ErrInvalidComponent))

// 		return
// 	}

// 	cd.ForceBranch(branchIndex)
// }

func (gameProp *GameProperty) onStepEnd(curBetMode int, gp *GameParams, pr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) {
	pr.CashWin = 0
	pr.CashWin = 0

	for _, v := range pr.Results {
		if !v.IsNoPayNow {
			pr.CashWin += int64(v.CashWin)
			pr.CoinWin += v.CoinWin
		}
	}

	gameProp.featureLevel.OnStepEnd(gameProp, gp, pr)
	gameProp.rng.OnStepEnd(curBetMode, gp, pr, prs)

	if gAllowStats2 {
		for _, v := range gp.RespinComponents {
			ic, isok := gameProp.Components.MapComponents[v]
			if isok && ic.IsRespin() {
				if !gameProp.stats2Cache.HasFeature(v) {
					gameProp.stats2Cache.AddFeature(v,
						ic.NewStats2(gameProp.Components.statsNodeData.GetParent(v)),
						true)
				}

				// ic.OnStats2(gameProp.GetComponentData(ic), gameProp.stats2Cache, gameProp, gp)
			}
		}

		for _, v := range gameProp.stats2Cache.RespinArr {
			ic, isok := gameProp.Components.MapComponents[v]
			if isok {
				ic.OnStats2(gameProp.GetComponentData(ic), gameProp.stats2Cache, gameProp, gp, pr, true)
			}
		}

		for _, v := range gameProp.lstNeedOnStepEndStats2Components {
			ic, isok := gameProp.Components.MapComponents[v]
			if isok {
				ic.OnStats2(gameProp.GetComponentData(ic), gameProp.stats2Cache, gameProp, gp, pr, true)
			}
		}

		gameProp.stats2Cache.OnStepEnd(gp.RespinComponents)
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
