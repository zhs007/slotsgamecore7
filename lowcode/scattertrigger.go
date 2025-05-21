package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const ScatterTriggerTypeName = "scatterTrigger"

type ScatterTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
	Wins          int
	WinMulti      int
}

// OnNewGame -
func (scatterTriggerData *ScatterTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	scatterTriggerData.BasicComponentData.OnNewGame(gameProp, component)

	scatterTrigger, isok := component.(*ScatterTrigger)
	if isok {
		if scatterTrigger.Config.Height > 0 {
			scatterTriggerData.SetConfigIntVal(CCVHeight, scatterTrigger.Config.Height)
		}
	}
}

// onNewStep -
func (scatterTriggerData *ScatterTriggerData) onNewStep() {
	scatterTriggerData.UsedResults = nil

	scatterTriggerData.NextComponent = ""
	scatterTriggerData.SymbolNum = 0
	scatterTriggerData.WildNum = 0
	scatterTriggerData.RespinNum = 0
	scatterTriggerData.Wins = 0
	scatterTriggerData.WinMulti = 1
}

// Clone
func (scatterTriggerData *ScatterTriggerData) Clone() IComponentData {
	target := &ScatterTriggerData{
		BasicComponentData: scatterTriggerData.CloneBasicComponentData(),
		NextComponent:      scatterTriggerData.NextComponent,
		SymbolNum:          scatterTriggerData.SymbolNum,
		WildNum:            scatterTriggerData.WildNum,
		RespinNum:          scatterTriggerData.RespinNum,
		Wins:               scatterTriggerData.Wins,
		WinMulti:           scatterTriggerData.WinMulti,
	}

	return target
}

// BuildPBComponentData
func (scatterTriggerData *ScatterTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ScatterTriggerData{
		BasicComponentData: scatterTriggerData.BuildPBBasicComponentData(),
		NextComponent:      scatterTriggerData.NextComponent,
		SymbolNum:          int32(scatterTriggerData.SymbolNum),
		WildNum:            int32(scatterTriggerData.WildNum),
		RespinNum:          int32(scatterTriggerData.RespinNum),
		Wins:               int32(scatterTriggerData.Wins),
		WinMulti:           int32(scatterTriggerData.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (scatterTriggerData *ScatterTriggerData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolNum {
		return scatterTriggerData.SymbolNum, true
	} else if key == CVWildNum {
		return scatterTriggerData.WildNum, true
	} else if key == CVRespinNum {
		return scatterTriggerData.RespinNum, true
	} else if key == CVWins {
		return scatterTriggerData.Wins, true
	} else if key == CVResultNum || key == CVWinResultNum {
		return len(scatterTriggerData.UsedResults), true
	}

	return 0, false
}

// ScatterTriggerConfig - configuration for ScatterTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type ScatterTriggerConfig struct {
	BasicComponentConfig            `yaml:",inline" json:",inline"`
	Symbols                         []string                      `yaml:"symbols" json:"symbols"`                       // like scatter
	SymbolCodes                     []int                         `yaml:"-" json:"-"`                                   // like scatter
	Type                            string                        `yaml:"type" json:"type"`                             // like scatters
	TriggerType                     SymbolTriggerType             `yaml:"-" json:"-"`                                   // SymbolTriggerType
	BetTypeString                   string                        `yaml:"betType" json:"betType"`                       // bet or totalBet or noPay
	BetType                         BetType                       `yaml:"-" json:"-"`                                   // bet or totalBet or noPay
	OSMulTypeString                 string                        `yaml:"symbolValsMulti" json:"symbolValsMulti"`       // OtherSceneMultiType
	OSMulType                       OtherSceneMultiType           `yaml:"-" json:"-"`                                   // OtherSceneMultiType
	MinNum                          int                           `yaml:"minNum" json:"minNum"`                         // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols                     []string                      `yaml:"wildSymbols" json:"wildSymbols"`               // wild etc
	WildSymbolCodes                 []int                         `yaml:"-" json:"-"`                                   // wild symbolCode
	PosArea                         []int                         `yaml:"posArea" json:"posArea"`                       // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
	CountScatterPayAs               string                        `yaml:"countScatterPayAs" json:"countScatterPayAs"`   // countscatter时，按什么符号赔付
	SymbolCodeCountScatterPayAs     int                           `yaml:"-" json:"-"`                                   // countscatter时，按什么符号赔付
	WinMulti                        int                           `yaml:"winMulti" json:"winMulti"`                     // winMulti，最后的中奖倍数，默认为1
	Height                          int                           `yaml:"Height" json:"Height"`                         // Height
	MaxHeight                       int                           `yaml:"MaxHeight" json:"MaxHeight"`                   // MaxHeight
	IsReversalHeight                bool                          `yaml:"isReversalHeight" json:"isReversalHeight"`     // isReversalHeight
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`       // jump to
	PiggyBankComponent              string                        `yaml:"piggyBankComponent" json:"piggyBankComponent"` // piggyBank component
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`               // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                         // 新的奖励系统
	TargetMask                      string                        `yaml:"targetMask" json:"targetMask"`                 // 如果是scatter这一组判断，可以把结果传递给一个mask
	OutputToComponent               string                        `yaml:"outputToComponent" json:"outputToComponent"`   // 将结果给到一个 positionCollection
	ReelsCollector                  string                        `yaml:"reelsCollector" json:"reelsCollector"`
	IsReverse                       bool                          `yaml:"isReverse" json:"isReverse"`                                         // 如果isReverse，表示判定为否才触发
	IsAddRespinMode                 bool                          `yaml:"isAddRespinMode" json:"isAddRespinMode"`                             // 是否是增加respinNum模式，默认是增加triggerNum模式
	RespinComponent                 string                        `yaml:"respinComponent" json:"respinComponent"`                             // respin component
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
	MapAwards                       map[int][]*Award              `yaml:"mapAwards" json:"mapAwards"`                                         // 新的奖励系统
}

// SetLinkComponent
func (cfg *ScatterTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type ScatterTrigger struct {
	*BasicComponent `json:"-"`
	Config          *ScatterTriggerConfig `json:"config"`
}

// Init -
func (scatterTrigger *ScatterTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ScatterTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ScatterTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ScatterTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return scatterTrigger.InitEx(cfg, pool)
}

// InitEx -
func (scatterTrigger *ScatterTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	scatterTrigger.Config = cfg.(*ScatterTriggerConfig)
	scatterTrigger.Config.ComponentType = ScatterTriggerTypeName

	scatterTrigger.Config.OSMulType = ParseOtherSceneMultiType(scatterTrigger.Config.OSMulTypeString)

	for _, s := range scatterTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ScatterTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		scatterTrigger.Config.SymbolCodes = append(scatterTrigger.Config.SymbolCodes, sc)
	}

	sc, isok := pool.DefaultPaytables.MapSymbols[scatterTrigger.Config.CountScatterPayAs]
	if isok {
		scatterTrigger.Config.SymbolCodeCountScatterPayAs = sc
	}

	for _, s := range scatterTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ScatterTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		scatterTrigger.Config.WildSymbolCodes = append(scatterTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(scatterTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("ScatterTrigger.InitEx:WildSymbols",
			slog.String("SymbolTriggerType", scatterTrigger.Config.Type),
			goutils.Err(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	scatterTrigger.Config.TriggerType = stt

	scatterTrigger.Config.BetType = ParseBetType(scatterTrigger.Config.BetTypeString)

	for _, award := range scatterTrigger.Config.Awards {
		award.Init()
	}

	for _, lst := range scatterTrigger.Config.MapAwards {
		for _, award := range lst {
			award.Init()
		}
	}

	// if scatterTrigger.Config.SymbolAwardsWeights != nil {
	// 	scatterTrigger.Config.SymbolAwardsWeights.Init()
	// }

	if scatterTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(scatterTrigger.Config.RespinNumWeight, scatterTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ScatterTrigger.InitEx:LoadIntWeights",
				slog.String("Weight", scatterTrigger.Config.RespinNumWeight),
				goutils.Err(err))

			return err
		}

		scatterTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(scatterTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range scatterTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, scatterTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("ScatterTrigger.InitEx:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))

				return err
			}

			scatterTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if scatterTrigger.Config.WinMulti <= 0 {
		scatterTrigger.Config.WinMulti = 1
	}

	scatterTrigger.onInit(&scatterTrigger.Config.BasicComponentConfig)

	return nil
}

// procMask
func (scatterTrigger *ScatterTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if scatterTrigger.Config.TargetMask != "" {
		gameProp.UseComponent(scatterTrigger.Config.TargetMask)

		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, scatterTrigger.Config.TargetMask, mask, false)
	}

	return nil
}

// procReelsCollector
func (scatterTrigger *ScatterTrigger) procReelsCollector(gs *sgc7game.GameScene, gameProp *GameProperty, ips sgc7game.IPlayerState,
	ret *sgc7game.Result, stake *sgc7game.Stake) error {

	if scatterTrigger.Config.ReelsCollector != "" {
		gameProp.UseComponent(scatterTrigger.Config.ReelsCollector)

		reelsData := make([]int, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			reelsData[ret.Pos[i*2]]++
		}

		betMethod := stake.CashBet / stake.CoinBet

		ps, isok := ips.(*PlayerState)
		if !isok {
			goutils.Error("ScatterTrigger.procReelsCollector:PlayerState",
				goutils.Err(ErrIvalidPlayerState))

			return ErrIvalidPlayerState
		}

		return gameProp.Pool.ChgReelsCollector(gameProp, scatterTrigger.Config.ReelsCollector, ps, int(betMethod), int(stake.CoinBet), reelsData)
	}

	return nil
}

// procPositionCollection
func (scatterTrigger *ScatterTrigger) procPositionCollection(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	cd *ScatterTriggerData) error {

	if scatterTrigger.Config.OutputToComponent != "" {
		pcd := gameProp.GetComponentDataWithName(scatterTrigger.Config.OutputToComponent)
		if pcd != nil {
			gameProp.UseComponent(scatterTrigger.Config.OutputToComponent)
			pc := gameProp.Components.MapComponents[scatterTrigger.Config.OutputToComponent]

			for _, ri := range cd.UsedResults {
				ret := curpr.Results[ri]

				for i := 0; i < len(ret.Pos)/2; i++ {
					pc.AddPos(pcd, ret.Pos[i*2], ret.Pos[i*2+1])
				}
			}
		}
	}

	return nil
}

func (scatterTrigger *ScatterTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return scatterTrigger.Config.SymbolCodes
}

// CanTriggerWithScene -
func (scatterTrigger *ScatterTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) {
	return scatterTrigger.canTrigger(gameProp, gs, nil, curpr, stake)
}

// CanTrigger -
func (scatterTrigger *ScatterTrigger) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, _ *sgc7game.GameScene, _ *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	icd := gameProp.GetComponentData(scatterTrigger)
	std := icd.(*ScatterTriggerData)

	isTrigger := false
	lst := []*sgc7game.Result{}

	symbols := scatterTrigger.getSymbols(gameProp)

	if scatterTrigger.Config.TriggerType == STTypeScatters {
		for _, s := range symbols {
			ret := sgc7game.CalcScatter5(gs, gameProp.CurPaytables, s, gameProp.GetBet3(stake, scatterTrigger.Config.BetType),
				func(scatter int, cursymbol int) bool {
					return cursymbol == scatter || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, false, scatterTrigger.GetHeight(&std.BasicComponentData), scatterTrigger.Config.IsReversalHeight)

			if ret != nil {
				isTrigger = true

				lst = append(lst, ret)
			}
		}
	} else if scatterTrigger.Config.TriggerType == STTypeReelScatters {
		for _, s := range symbols {
			ret := sgc7game.CalcScatter5(gs, gameProp.CurPaytables, s, gameProp.GetBet3(stake, scatterTrigger.Config.BetType),
				func(scatter int, cursymbol int) bool {
					return cursymbol == scatter || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, true, scatterTrigger.GetHeight(&std.BasicComponentData), scatterTrigger.Config.IsReversalHeight)

			if ret != nil {
				isTrigger = true

				lst = append(lst, ret)
			}
		}
	} else if scatterTrigger.Config.TriggerType == STTypeCountScatter {
		ret := sgc7game.CalcScatterEx2(gs, symbols[0], scatterTrigger.Config.MinNum, func(scatter int, cursymbol int) bool {
			return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		}, scatterTrigger.GetHeight(&std.BasicComponentData), scatterTrigger.Config.IsReversalHeight)

		if ret != nil {
			if scatterTrigger.Config.BetType != BTypeNoPay {
				if scatterTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet3(stake, scatterTrigger.Config.BetType)
				}
			}

			isTrigger = true

			lst = append(lst, ret)
		}
	} else if scatterTrigger.Config.TriggerType == STTypeCountScatterReels {
		ret := sgc7game.CalcReelScatterEx2(gs, symbols[0], scatterTrigger.Config.MinNum, func(scatter int, cursymbol int) bool {
			return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		}, scatterTrigger.GetHeight(&std.BasicComponentData), scatterTrigger.Config.IsReversalHeight)

		if ret != nil {
			if scatterTrigger.Config.BetType != BTypeNoPay {
				if scatterTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet3(stake, scatterTrigger.Config.BetType)
				}
			}

			isTrigger = true

			lst = append(lst, ret)
		}
	} else if scatterTrigger.Config.TriggerType == STTypeCountScatterInArea {
		ret := sgc7game.CountScatterInArea(gs, symbols[0], scatterTrigger.Config.MinNum,
			func(x, y int) bool {
				return IsInPosArea(x, y, scatterTrigger.Config.PosArea)
			},
			func(scatter int, cursymbol int) bool {
				return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		if ret != nil {
			if scatterTrigger.Config.BetType != BTypeNoPay {
				if scatterTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet3(stake, scatterTrigger.Config.BetType)
				}
			}

			isTrigger = true

			lst = append(lst, ret)
		}
	}

	if scatterTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

// procWins
func (scatterTrigger *ScatterTrigger) procWins(gameProp *GameProperty, curpr *sgc7game.PlayResult, std *ScatterTriggerData, lst []*sgc7game.Result) (int, error) {
	if scatterTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0

			scatterTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds
		}

		return 0, nil
	}

	std.WinMulti = scatterTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin

		scatterTrigger.AddResult(curpr, v, &std.BasicComponentData)

		std.SymbolNum += v.SymbolNums
		std.WildNum += v.Wilds
	}

	if std.Wins > 0 {
		if scatterTrigger.Config.PiggyBankComponent != "" {
			cd := gameProp.GetCurComponentDataWithName(scatterTrigger.Config.PiggyBankComponent)
			if cd == nil {
				goutils.Error("ScatterTrigger.procWins:GetCurComponentDataWithName",
					slog.String("PiggyBankComponent", scatterTrigger.Config.PiggyBankComponent),
					goutils.Err(ErrInvalidComponent))

				return 0, ErrInvalidComponent
			}

			cd.ChgConfigIntVal(CCVSavedMoney, std.Wins)

			for _, v := range lst {
				v.IsNoPayNow = true
			}

			gameProp.UseComponent(scatterTrigger.Config.PiggyBankComponent)
		}
	}

	return std.Wins, nil
}

// calcRespinNum
func (scatterTrigger *ScatterTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(scatterTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := scatterTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					slog.Int("SymbolNum", ret.SymbolNums),
					slog.String("componentName", scatterTrigger.GetName()),
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				slog.String("componentName", scatterTrigger.GetName()),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(scatterTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := scatterTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				slog.String("componentName", scatterTrigger.GetName()),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if scatterTrigger.Config.RespinNumWeightVW != nil {
		cr, err := scatterTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWeightVW",
				slog.String("componentName", scatterTrigger.GetName()),
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if scatterTrigger.Config.RespinNum > 0 {
		return scatterTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// OnProcControllers -
func (scatterTrigger *ScatterTrigger) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	lst, isok := scatterTrigger.Config.MapAwards[val]
	if isok {
		gameProp.procAwards(plugin, lst, curpr, gp)
	}

	if len(scatterTrigger.Config.Awards) > 0 {
		gameProp.procAwards(plugin, scatterTrigger.Config.Awards, curpr, gp)
	}
}

// playgame
func (scatterTrigger *ScatterTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	std := icd.(*ScatterTriggerData)
	std.onNewStep()

	gs := scatterTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	os := scatterTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := scatterTrigger.canTrigger(gameProp, gs, os, curpr, stake)

	if isTrigger {
		scatterTrigger.procWins(gameProp, curpr, std, lst)

		respinNum, err := scatterTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", err
		}

		std.RespinNum = respinNum

		err = scatterTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		err = scatterTrigger.procReelsCollector(gs, gameProp, ps, lst[0], stake)
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:procReelsCollector",
				goutils.Err(err))

			return "", err
		}

		err = scatterTrigger.procPositionCollection(gameProp, curpr, std)
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:procPositionCollection",
				goutils.Err(err))

			return "", err
		}

		scatterTrigger.ProcControllers(gameProp, plugin, curpr, gp, lst[0].SymbolNums, "")

		if scatterTrigger.Config.RespinComponent != "" {
			if gameProp.IsRespin(scatterTrigger.Config.RespinComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if scatterTrigger.Config.ForceToNext {
						std.NextComponent = scatterTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(scatterTrigger.Config.RespinComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, scatterTrigger.Config.RespinComponent, !scatterTrigger.Config.IsAddRespinMode)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, scatterTrigger.Config.RespinComponent, !scatterTrigger.Config.IsAddRespinMode)

					lst[0].Type = sgc7game.RTFreeGame
					lst[0].Value = std.RespinNum
				}
			}
		}

		if scatterTrigger.Config.JumpToComponent != "" {
			std.NextComponent = scatterTrigger.Config.JumpToComponent

			nc := scatterTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := scatterTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := scatterTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (scatterTrigger *ScatterTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*ScatterTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", scatterTrigger.Name, std.NextComponent)
	}

	return nil
}

// NewComponentData -
func (scatterTrigger *ScatterTrigger) NewComponentData() IComponentData {
	return &ScatterTriggerData{}
}

func (scatterTrigger *ScatterTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return scatterTrigger.Config.WinMulti
}

func (scatterTrigger *ScatterTrigger) GetHeight(basicCD *BasicComponentData) int {
	height, isok := basicCD.GetConfigIntVal(CCVHeight)
	if isok {
		if height > scatterTrigger.Config.MaxHeight {
			return scatterTrigger.Config.MaxHeight
		}

		return height
	}

	if scatterTrigger.Config.Height > scatterTrigger.Config.MaxHeight {
		return scatterTrigger.Config.MaxHeight
	}

	return scatterTrigger.Config.Height
}

// GetAllLinkComponents - get all link components
func (scatterTrigger *ScatterTrigger) GetAllLinkComponents() []string {
	return []string{scatterTrigger.Config.DefaultNextComponent, scatterTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (scatterTrigger *ScatterTrigger) GetNextLinkComponents() []string {
	return []string{scatterTrigger.Config.DefaultNextComponent, scatterTrigger.Config.JumpToComponent}
}

// NewStats2 -
func (scatterTrigger *ScatterTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins, stats2.OptIntVal})
}

// OnStats2
func (scatterTrigger *ScatterTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	scatterTrigger.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*ScatterTriggerData)

	s2.ProcStatsWins(scatterTrigger.Name, int64(cd.Wins))

	s2.ProcStatsIntVal(scatterTrigger.Name, cd.SymbolNum)
}

func NewScatterTrigger(name string) IComponent {
	return &ScatterTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//		"triggerType": "countscatter",
//		"betType": "bet",
//		"triggerRespinType": "respinNum",
//		"winMulti": 1,
//		"minNum": 3,
//		"countScatterPayAs": "S",
//		"symbols": [
//			"W"
//		],
//		"respinComponent": "fg-start",
//		"genRespinType": "number",
//		"respinNum": 10,
//		"putMoneyInPiggyBank": "bg-piggybank"
// 		"reelsCollector": "bg-collect"

type jsonScatterTrigger struct {
	Symbols                       []string   `json:"symbols"`
	TriggerType                   string     `json:"triggerType"`
	BetType                       string     `json:"betType"`
	SymbolValsMulti               string     `json:"symbolValsMulti"`
	MinNum                        int        `json:"minNum"`
	WildSymbols                   []string   `json:"wildSymbols"`
	PosArea                       []int      `json:"posArea"`
	CountScatterPayAs             string     `json:"countScatterPayAs"`
	WinMulti                      int        `json:"winMulti"`
	TargetMask                    string     `json:"targetMask"`
	TriggerRespinType             string     `json:"triggerRespinType"`
	RespinComponent               string     `json:"respinComponent"`
	PutMoneyInPiggyBank           string     `json:"putMoneyInPiggyBank"`
	GenRespinType                 string     `json:"genRespinType"`
	RespinNum                     int        `json:"respinNum"`
	RespinNumWeight               string     `json:"respinNumWeight"`
	RespinNumWithScatterNum       [][]int    `json:"respinNumWithScatterNum"`
	RespinNumWeightWithScatterNum [][]string `json:"respinNumWeightWithScatterNum"`
	Height                        int        `json:"Height"`
	MaxHeight                     int        `json:"MaxHeight"`
	IsReversalHeight              bool       `json:"isReversalHeight"`
	OutputToComponent             string     `json:"outputToComponent"`
	ReelsCollector                string     `json:"reelsCollector"`
}

func (jcfg *jsonScatterTrigger) build() *ScatterTriggerConfig {
	cfg := &ScatterTriggerConfig{
		Symbols:            jcfg.Symbols,
		Type:               jcfg.TriggerType,
		BetTypeString:      jcfg.BetType,
		MinNum:             jcfg.MinNum,
		WildSymbols:        jcfg.WildSymbols,
		PosArea:            jcfg.PosArea,
		CountScatterPayAs:  jcfg.CountScatterPayAs,
		WinMulti:           jcfg.WinMulti,
		TargetMask:         jcfg.TargetMask,
		PiggyBankComponent: jcfg.PutMoneyInPiggyBank,
		IsAddRespinMode:    jcfg.TriggerRespinType == "respinNum",
		RespinComponent:    jcfg.RespinComponent,
		RespinNum:          jcfg.RespinNum,
		RespinNumWeight:    jcfg.RespinNumWeight,
		OSMulTypeString:    jcfg.SymbolValsMulti,
		Height:             jcfg.Height,
		MaxHeight:          jcfg.MaxHeight,
		IsReversalHeight:   jcfg.IsReversalHeight,
		OutputToComponent:  jcfg.OutputToComponent,
		ReelsCollector:     jcfg.ReelsCollector,
	}

	if jcfg.TriggerRespinType != "none" {
		if jcfg.RespinNumWithScatterNum != nil {
			cfg.RespinNumWithScatterNum = make(map[int]int)
			for _, arr := range jcfg.RespinNumWithScatterNum {
				cfg.RespinNumWithScatterNum[arr[0]] = arr[1]
			}
		}

		if jcfg.RespinNumWeightWithScatterNum != nil {
			cfg.RespinNumWeightWithScatterNum = make(map[int]string)
			for _, arr := range jcfg.RespinNumWeightWithScatterNum {
				i64, err := goutils.String2Int64(arr[0])
				if err != nil {
					goutils.Error("jsonScatterTrigger:RespinNumWeightWithScatterNum:String2Int64",
						goutils.Err(err))

					return nil
				}

				cfg.RespinNumWeightWithScatterNum[int(i64)] = arr[1]
			}
		}
	}

	for i := range cfg.PosArea {
		cfg.PosArea[i]--
	}

	return cfg
}

func parseScatterTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseScatterTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseScatterTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonScatterTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseScatterTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, mapAwards, err := parseScatterTriggerControllers(ctrls)
		if err != nil {
			goutils.Error("parseScatterTrigger:parseScatterTriggerControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ScatterTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
