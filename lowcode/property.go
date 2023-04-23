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

	GamePropTriggerFG = 100
	GamePropFGNum     = 101

	GamePropNextComponent   = 200
	GamePropRespinComponent = 201

	GamePropCurMystery = 1000
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
	MapComponentData  map[string]IComponentData
	HistoryComponents []IComponent
	RespinComponents  []string
}

func (gameProp *GameProperty) OnNewStep() error {
	gameProp.mapInt = make(map[string]int)
	gameProp.mapStr = make(map[string]string)

	gameProp.SetStrVal(GamePropNextComponent, "")
	gameProp.SetStrVal(GamePropRespinComponent, "")

	gameProp.HistoryComponents = nil

	return nil
}

func (gameProp *GameProperty) TagScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return pr.Scenes[len(pr.Scenes)-1], len(pr.Scenes) - 1
	}

	return pr.Scenes[si], si
}

func (gameProp *GameProperty) TagOtherScene(pr *sgc7game.PlayResult, tag string, sceneIndex int) {
	gameProp.mapInt[tag] = sceneIndex
}

func (gameProp *GameProperty) GetOtherScene(pr *sgc7game.PlayResult, tag string) (*sgc7game.GameScene, int) {
	si, isok := gameProp.mapInt[tag]
	if !isok {
		return pr.OtherScenes[len(pr.OtherScenes)-1], len(pr.OtherScenes) - 1
	}

	return pr.OtherScenes[si], si
}

func (gameProp *GameProperty) Respin(pr *sgc7game.PlayResult, gp *GameParams, respinComponent string, gs *sgc7game.GameScene, os *sgc7game.GameScene) {
	if gs != nil {
		gp.LastScene = gs.Clone()
	}

	if os != nil {
		gp.LastOtherScene = os.Clone()
	}

	gameProp.SetStrVal(GamePropRespinComponent, respinComponent)

	gp.NextStepFirstComponent = respinComponent
}

func (gameProp *GameProperty) OnFGSpin() error {
	gameProp.SetVal(GamePropFGNum, gameProp.GetVal(GamePropFGNum)-1)

	return nil
}

func (gameProp *GameProperty) TriggerFG(pr *sgc7game.PlayResult, gp *GameParams, fgnum int, respinFirstComponent string) error {
	if fgnum > 0 {
		if gameProp.GetVal(GamePropTriggerFG) > 0 {
			gameProp.RetriggerFG(pr, gp, fgnum)
		} else {
			gameProp.SetVal(GamePropTriggerFG, 1)
			gameProp.SetVal(GamePropFGNum, fgnum)

			gameProp.SetStrVal(GamePropRespinComponent, respinFirstComponent)

			gp.NextStepFirstComponent = respinFirstComponent
		}
	}

	return nil
}

func (gameProp *GameProperty) onTriggerRespin(respinComponent string) error {
	gameProp.RespinComponents = append(gameProp.RespinComponents, respinComponent)

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
	} else if !pr.IsWait {
		pr.IsFinish = true
	}
}

func (gameProp *GameProperty) TriggerRespin(pr *sgc7game.PlayResult, gp *GameParams, respinNum int, respinComponent string) error {
	if respinNum > 0 {
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
	}

	return nil
}

func (gameProp *GameProperty) TriggerRespinWithWeights(pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, fn string, respinComponent string) error {
	vw2, err := gameProp.GetIntValWeights(fn)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:GetIntValWeights",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	val, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:RandVal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	if val.Int() > 0 {
		gameProp.TriggerRespin(pr, gp, val.Int(), respinComponent)
	}

	return nil
}

func (gameProp *GameProperty) RetriggerFG(pr *sgc7game.PlayResult, gp *GameParams, fgnum int) error {
	if fgnum > 0 {
		gameProp.AddVal(GamePropFGNum, fgnum)
	}

	return nil
}

func (gameProp *GameProperty) GetIntValWeights(fn string) (*sgc7game.ValWeights2, error) {
	vw2, isok := gameProp.MapIntValWeights[fn]
	if !isok {
		curvw2, err := sgc7game.LoadValWeights2FromExcel(gameProp.Pool.Config.GetPath(fn), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("GameProperty.GetIntValWeights:LoadValWeights2FromExcel",
				zap.String("fn", fn),
				zap.Error(err))

			return nil, err
		}

		gameProp.MapIntValWeights[fn] = curvw2

		vw2 = curvw2
	}

	return vw2, nil
}

func (gameProp *GameProperty) TriggerFGWithWeights(pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, fn string, respinFirstComponent string) error {
	vw2, err := gameProp.GetIntValWeights(fn)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:GetIntValWeights",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	val, err := vw2.RandVal(plugin)
	if err != nil {
		goutils.Error("GameProperty.TriggerFGWithWeights:RandVal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	if val.Int() > 0 {
		gameProp.SetVal(GamePropTriggerFG, 1)
		gameProp.SetVal(GamePropFGNum, val.Int())

		gameProp.SetStrVal(GamePropRespinComponent, respinFirstComponent)

		gp.NextStepFirstComponent = respinFirstComponent
	}

	return nil
}

func (gameProp *GameProperty) SetVal(prop int, val int) error {
	if prop == GamePropCurMystery {
		str := gameProp.CurPaytables.GetStringFromInt(val)

		gameProp.MapStrVals[prop] = str
	}

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
	if prop == GamePropCurMystery {
		v, isok := gameProp.CurPaytables.MapSymbols[val]
		if !isok {
			goutils.Error("GameProperty.SetStrVal:GamePropCurMystery",
				zap.String("val", val),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gameProp.MapVals[prop] = v
	} else if prop == GamePropCurPaytables {
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

func (gameProp *GameProperty) procAward(award *Award, curpr *sgc7game.PlayResult) {
	if award.AwardType == AwardRespinTimes {
		component, isok := gameProp.Pool.MapComponents[award.Config.StrParam]
		if isok {
			respin, isok := component.(*Respin)
			if isok {
				respin.AddRespinTimes(gameProp, award.Config.Val)
			}
		}
	}
}

func init() {
	MapProperty = make(map[string]int)

	MapProperty["width"] = GamePropWidth
	MapProperty["height"] = GamePropHeight
	MapProperty["paytables"] = GamePropCurPaytables
	MapProperty["reels"] = GamePropCurReels
	MapProperty["linedata"] = GamePropCurLineData

	MapProperty["triggerFG"] = GamePropTriggerFG
	MapProperty["FGNum"] = GamePropFGNum

	MapProperty["curMystery"] = GamePropCurMystery
}
