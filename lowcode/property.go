package lowcode

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
)

const (
	GamePropWidth        = 1
	GamePropHeight       = 2
	GamePropCurPaytables = 3
	GamePropCurReels     = 4
	GamePropCurLineData  = 5
	GamePropCurLineNum   = 6
	GamePropCurBetIndex  = 7

	GamePropStepMulti = 100
	GamePropGameMulti = 101

	GamePropNextComponent   = 200
	GamePropRespinComponent = 201
)

var MapProperty map[string]int

func String2Property(str string) (int, error) {
	v, isok := MapProperty[str]
	if isok {
		return v, nil
	}

	goutils.Error("String2Property",
		zap.String("str", str),
		zap.Error(ErrInvalidGamePropertyString))

	return 0, ErrInvalidGamePropertyString
}

type GameProperty struct {
	Pool              *GamePropertyPool
	MapVals           map[int]int
	MapStrVals        map[int]string
	CurPaytables      *sgc7game.PayTables
	CurLineData       *sgc7game.LineData
	CurReels          *sgc7game.ReelsData
	MapIntValWeights  map[string]*sgc7game.ValWeights2
	MapStats          map[string]*sgc7stats.Feature
	mapInt            map[string]int
	mapStr            map[string]string
	mapGlobalStr      map[string]string
	MapComponentData  map[string]IComponentData
	HistoryComponents []IComponent
	RespinComponents  []string
	PoolScene         *sgc7game.GameScenePoolEx
}

func (gameProp *GameProperty) BuildGameParam(gp *GameParams) {
	if len(gameProp.RespinComponents) > 0 {
		gp.RespinComponents = make([]string, len(gameProp.RespinComponents))

		copy(gp.RespinComponents, gameProp.RespinComponents)
	}

	gp.SetGameProp(gameProp)
}

func (gameProp *GameProperty) OnNewGame(stake *sgc7game.Stake) error {
	gameProp.SetVal(GamePropGameMulti, 1)

	curBet := stake.CashBet / stake.CoinBet
	for i, v := range gameProp.Pool.Config.Bets {
		if v == int(curBet) {
			gameProp.SetVal(GamePropCurBetIndex, i)

			break
		}
	}

	gameProp.mapGlobalStr = make(map[string]string)

	return nil
}

func (gameProp *GameProperty) OnNewStep() error {
	gameProp.mapInt = make(map[string]int)
	gameProp.mapStr = make(map[string]string)

	gameProp.SetStrVal(GamePropNextComponent, "")
	gameProp.SetStrVal(GamePropRespinComponent, "")

	gameProp.SetVal(GamePropStepMulti, 1)

	gameProp.HistoryComponents = nil

	return nil
}

func (gameProp *GameProperty) TagScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return nil, -1
	}

	return pr.Scenes[si], si
}

func (gameProp *GameProperty) TagOtherScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetOtherScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return nil, -1
	}

	return pr.OtherScenes[si], si
}

func (gameProp *GameProperty) Respin(pr *sgc7game.PlayResult, gp *GameParams, respinComponent string, gs *sgc7game.GameScene, os *sgc7game.GameScene) {
	if gs != nil {
		// gp.LastScene = gs.Clone()
		gp.LastScene = gs.CloneEx(gameProp.PoolScene)
	}

	if os != nil {
		// gp.LastOtherScene = os.Clone()
		gp.LastOtherScene = os.CloneEx(gameProp.PoolScene)
	}

	gameProp.SetStrVal(GamePropRespinComponent, respinComponent)

	gp.NextStepFirstComponent = respinComponent
}

func (gameProp *GameProperty) onTriggerRespin(respinComponent string) error {
	if len(gameProp.RespinComponents) == 0 || gameProp.RespinComponents[len(gameProp.RespinComponents)-1] != respinComponent {
		gameProp.RespinComponents = append(gameProp.RespinComponents, respinComponent)
	}

	return nil
}

func (gameProp *GameProperty) onRespinEnding(respinComponent string) error {
	if len(gameProp.RespinComponents) > 0 && gameProp.RespinComponents[len(gameProp.RespinComponents)-1] == respinComponent {
		gameProp.RespinComponents = gameProp.RespinComponents[0 : len(gameProp.RespinComponents)-1]
	}

	return nil
}

func (gameProp *GameProperty) ProcRespin(pr *sgc7game.PlayResult, gp *GameParams) {
	if len(gameProp.RespinComponents) > 0 {
		gp.NextStepFirstComponent = gameProp.RespinComponents[len(gameProp.RespinComponents)-1]

		pr.IsFinish = false

		if goutils.IndexOfStringSlice(gp.HistoryComponents, gp.NextStepFirstComponent, 0) < 0 {
			gp.AddComponentData(gp.NextStepFirstComponent, gameProp.MapComponentData[gp.NextStepFirstComponent])
		}
	} else if !pr.IsWait {
		pr.IsFinish = true
	}
}

func (gameProp *GameProperty) AddComponent2History(component IComponent, gp *GameParams) {
	for _, c := range gameProp.HistoryComponents {
		if c.GetName() == component.GetName() {
			return
		}
	}

	gameProp.HistoryComponents = append(gameProp.HistoryComponents, component)
	gp.HistoryComponents = append(gp.HistoryComponents, component.GetName())
}

func (gameProp *GameProperty) TriggerRespin(pr *sgc7game.PlayResult, gp *GameParams, respinNum int, respinComponent string) error {
	// if respinNum > 0 {
	component, isok := gameProp.Pool.MapComponents[respinComponent]
	if isok {
		respin, isok := component.(*Respin)
		if isok {
			respin.AddRespinTimes(gameProp, respinNum)

			gameProp.SetStrVal(GamePropRespinComponent, respinComponent)
			gameProp.onTriggerRespin(respinComponent)

			gp.NextStepFirstComponent = respinComponent
		}
	}
	// }

	return nil
}

func (gameProp *GameProperty) TriggerRespinWithWeights(pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, fn string, useFileMapping bool, respinComponent string) (int, error) {
	vw2, err := gameProp.GetIntValWeights(fn, useFileMapping)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:GetIntValWeights",
			zap.String("fn", fn),
			zap.Error(err))

		return 0, err
	}

	val, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:RandVal",
			zap.String("fn", fn),
			zap.Error(err))

		return 0, err
	}

	if val.Int() > 0 {
		gameProp.TriggerRespin(pr, gp, val.Int(), respinComponent)

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
				zap.String("Weight", fn),
				zap.Error(err))

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
				zap.String("val", val),
				zap.Error(ErrInvalidPaytables))

			return ErrInvalidPaytables
		}

		gameProp.CurPaytables = v
	} else if prop == GamePropCurLineData {
		v, isok := gameProp.Pool.Config.MapLinedate[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurLineData",
				zap.String("val", val),
				zap.Error(ErrInvalidPaytables))

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

func (gameProp *GameProperty) procAwards(plugin sgc7plugin.IPlugin, awards []*Award, curpr *sgc7game.PlayResult, gp *GameParams) {
	for _, v := range awards {
		gameProp.procAward(plugin, v, curpr, gp)
	}
}

func (gameProp *GameProperty) procAward(plugin sgc7plugin.IPlugin, award *Award, curpr *sgc7game.PlayResult, gp *GameParams) {
	if award.Type == AwardRespinTimes {
		component, isok := gameProp.Pool.MapComponents[award.StrParams[0]]
		if isok {
			respin, isok := component.(*Respin)
			if isok {
				respin.AddRespinTimes(gameProp, award.Vals[0])
			}
		}
	} else if award.Type == AwardGameMulti {
		gameProp.SetVal(GamePropGameMulti, award.Vals[0])
	} else if award.Type == AwardStepMulti {
		gameProp.SetVal(GamePropStepMulti, award.Vals[0])
	} else if award.Type == AwardInitMask {
		component, isok := gameProp.Pool.MapComponents[award.StrParams[0]]
		if isok {
			mask, isok := component.(*Mask)
			if isok {
				mask.ProcMask(plugin, gameProp, curpr, gp, award.StrParams[1])
			}
		}
	} else if award.Type == AwardTriggerRespin {
		gameProp.TriggerRespin(curpr, gp, award.Vals[0], award.StrParams[0])
	} else if award.Type == AwardCollector {
		component, isok := gameProp.Pool.MapComponents[award.StrParams[0]]
		if isok {
			collector, isok := component.(*Collector)
			if isok {
				err := collector.Add(plugin, award.Vals[0], nil, gameProp, curpr, gp, false)
				if err != nil {
					goutils.Error("GameProperty.procAward",
						zap.Error(err))

					return
				}
			}
		}
	} else if award.Type == AwardNoLevelUpCollector {
		component, isok := gameProp.Pool.MapComponents[award.StrParams[0]]
		if isok {
			collector, isok := component.(*Collector)
			if isok {
				err := collector.Add(plugin, award.Vals[0], nil, gameProp, curpr, gp, true)
				if err != nil {
					goutils.Error("GameProperty.procAward",
						zap.Error(err))

					return
				}
			}
		}
	} else if award.Type == AwardWeightGameRNG {
		vw, err := gameProp.GetIntValWeights(award.StrParams[0], true)
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardWeightGameRNG:GetIntValWeights",
				zap.Error(err))

			return
		}

		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("GameProperty.procAward:AwardWeightGameRNG:RandVal",
				zap.Error(err))

			return
		}

		gameProp.TagInt(award.StrParams[1], cr.Int())
	} else if award.Type == AwardPushSymbolCollection {
		component, isok := gameProp.Pool.MapComponents[award.StrParams[0]]
		if isok {
			symbolCollection, isok := component.(*SymbolCollection)
			if isok {
				for i := 0; i < award.Vals[0]; i++ {
					err := symbolCollection.Push(plugin, gameProp, gp)
					if err != nil {
						goutils.Error("GameProperty.procAward:AwardPushSymbolCollection:Push",
							zap.Error(err))

						return
					}
				}

				gameProp.AddComponent2History(component, gp)
			}
		}
	}
}

func (gameProp *GameProperty) procOtherSceneFeature(otherSceneFeature *OtherSceneFeature, curpr *sgc7game.PlayResult, os *sgc7game.GameScene) {
	if otherSceneFeature.Type == OtherSceneFeatureGameMulti {
		mul := 1

		for _, arr := range os.Arr {
			for _, v := range arr {
				mul *= v
			}
		}

		gameProp.SetVal(GamePropGameMulti, mul)
	} else if otherSceneFeature.Type == OtherSceneFeatureGameMultiSum {
		mul := 0

		for _, arr := range os.Arr {
			for _, v := range arr {
				if v > 1 {
					mul += v
				}
			}
		}

		gameProp.SetVal(GamePropGameMulti, mul)
	} else if otherSceneFeature.Type == OtherSceneFeatureStepMultiSum {
		mul := 0

		for _, arr := range os.Arr {
			for _, v := range arr {
				if v > 1 {
					mul += v
				}
			}
		}

		gameProp.SetVal(GamePropStepMulti, mul)
	} else if otherSceneFeature.Type == OtherSceneFeatureStepMulti {
		mul := 1

		for _, arr := range os.Arr {
			for _, v := range arr {
				mul *= v
			}
		}

		gameProp.SetVal(GamePropStepMulti, mul)
	}
}

func (gameProp *GameProperty) ProcMulti(ret *sgc7game.Result) {
	mul := gameProp.GetVal(GamePropStepMulti) * gameProp.GetVal(GamePropGameMulti)
	ret.CoinWin *= mul
	ret.CashWin *= mul
}

func (gameProp *GameProperty) GetBet(stake *sgc7game.Stake, bettype string) int {
	if bettype == BetTypeTotalBet {
		return int(stake.CoinBet) * gameProp.Pool.Config.TotalBetInWins[gameProp.GetVal(GamePropCurBetIndex)]
	}

	return int(stake.CoinBet)
}

func (gameProp *GameProperty) GetBet2(stake *sgc7game.Stake, bt BetType) int {
	if bt == BTypeTotalBet {
		return int(stake.CoinBet) * gameProp.Pool.Config.TotalBetInWins[gameProp.GetVal(GamePropCurBetIndex)]
	} else if bt == BTypeBet {
		return int(stake.CoinBet)
	}

	return 0
}

func init() {
	MapProperty = make(map[string]int)

	MapProperty["width"] = GamePropWidth
	MapProperty["height"] = GamePropHeight
	MapProperty["paytables"] = GamePropCurPaytables
	MapProperty["reels"] = GamePropCurReels
	MapProperty["linedata"] = GamePropCurLineData
}
