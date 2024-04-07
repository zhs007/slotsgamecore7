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

const (
	STCVWinMulti string = "winMulti" // 可以修改配置项里的winMulti
)

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
}

// onNewStep -
func (scatterTriggerData *ScatterTriggerData) onNewStep() {
	// scatterTriggerData.BasicComponentData.OnNewStep(gameProp, component)

	scatterTriggerData.UsedResults = nil

	scatterTriggerData.NextComponent = ""
	scatterTriggerData.SymbolNum = 0
	scatterTriggerData.WildNum = 0
	scatterTriggerData.RespinNum = 0
	scatterTriggerData.Wins = 0
	scatterTriggerData.WinMulti = 1
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

// GetVal -
func (scatterTriggerData *ScatterTriggerData) GetVal(key string) (int, bool) {
	if key == CVSymbolNum {
		return scatterTriggerData.SymbolNum, true
	} else if key == CVWildNum {
		return scatterTriggerData.WildNum, true
	} else if key == CVRespinNum {
		return scatterTriggerData.RespinNum, true
	} else if key == CVWins {
		return scatterTriggerData.Wins, true
	}

	return 0, false
}

// SetVal -
func (scatterTriggerData *ScatterTriggerData) SetVal(key string, val int) {
	if key == CVSymbolNum {
		scatterTriggerData.SymbolNum = val
	} else if key == CVWildNum {
		scatterTriggerData.WildNum = val
	} else if key == CVRespinNum {
		scatterTriggerData.RespinNum = val
	} else if key == CVWins {
		scatterTriggerData.Wins = val
	}
}

// ScatterTriggerConfig - configuration for ScatterTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type ScatterTriggerConfig struct {
	BasicComponentConfig        `yaml:",inline" json:",inline"`
	Symbols                     []string          `yaml:"symbols" json:"symbols"`                         // like scatter
	SymbolCodes                 []int             `yaml:"-" json:"-"`                                     // like scatter
	Type                        string            `yaml:"type" json:"type"`                               // like scatters
	TriggerType                 SymbolTriggerType `yaml:"-" json:"-"`                                     // SymbolTriggerType
	BetTypeString               string            `yaml:"betType" json:"betType"`                         // bet or totalBet or noPay
	BetType                     BetType           `yaml:"-" json:"-"`                                     // bet or totalBet or noPay
	MinNum                      int               `yaml:"minNum" json:"minNum"`                           // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols                 []string          `yaml:"wildSymbols" json:"wildSymbols"`                 // wild etc
	WildSymbolCodes             []int             `yaml:"-" json:"-"`                                     // wild symbolCode
	PosArea                     []int             `yaml:"posArea" json:"posArea"`                         // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
	CountScatterPayAs           string            `yaml:"countScatterPayAs" json:"countScatterPayAs"`     // countscatter时，按什么符号赔付
	SymbolCodeCountScatterPayAs int               `yaml:"-" json:"-"`                                     // countscatter时，按什么符号赔付
	WinMulti                    int               `yaml:"winMulti" json:"winMulti"`                       // winMulti，最后的中奖倍数，默认为1
	JumpToComponent             string            `yaml:"jumpToComponent" json:"jumpToComponent"`         // jump to
	PiggyBankComponent          string            `yaml:"piggyBankComponent" json:"piggyBankComponent"`   // piggyBank component
	ForceToNext                 bool              `yaml:"forceToNext" json:"forceToNext"`                 // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                      []*Award          `yaml:"awards" json:"awards"`                           // 新的奖励系统
	SymbolAwardsWeights         *AwardsWeights    `yaml:"symbolAwardsWeights" json:"symbolAwardsWeights"` // 每个中奖符号随机一组奖励
	TargetMask                  string            `yaml:"targetMask" json:"targetMask"`                   // 如果是scatter这一组判断，可以把结果传递给一个mask
	IsReverse                   bool              `yaml:"isReverse" json:"isReverse"`                     // 如果isReverse，表示判定为否才触发
	// NeedDiscardResults              bool                          `yaml:"needDiscardResults" json:"needDiscardResults"`                       // 如果needDiscardResults，表示抛弃results
	IsAddRespinMode                 bool                          `yaml:"isAddRespinMode" json:"isAddRespinMode"`                             // 是否是增加respinNum模式，默认是增加triggerNum模式
	RespinComponent                 string                        `yaml:"respinComponent" json:"respinComponent"`                             // respin component
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
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
		goutils.Error("SpSymbolTrigger.InitEx:WildSymbols",
			slog.String("SymbolTriggerType", scatterTrigger.Config.Type),
			goutils.Err(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	scatterTrigger.Config.TriggerType = stt

	scatterTrigger.Config.BetType = ParseBetType(scatterTrigger.Config.BetTypeString)

	for _, award := range scatterTrigger.Config.Awards {
		award.Init()
	}

	if scatterTrigger.Config.SymbolAwardsWeights != nil {
		scatterTrigger.Config.SymbolAwardsWeights.Init()
	}

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

	// if scatterTrigger.Config.BetType == BTypeNoPay {
	// 	scatterTrigger.Config.NeedDiscardResults = true
	// }

	scatterTrigger.onInit(&scatterTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
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

// // CanTrigger -
// func (scatterTrigger *ScatterTrigger) triggerScatter(gameProp *GameProperty, stake *sgc7game.Stake, gs *sgc7game.GameScene) *sgc7game.Result {
// 	return sgc7game.CalcScatter4(gs, gameProp.CurPaytables, scatterTrigger.Config.SymbolCodes[0], gameProp.GetBet2(stake, scatterTrigger.Config.BetType),
// 		func(scatter int, cursymbol int) bool {
// 			return goutils.IndexOfIntSlice(scatterTrigger.Config.SymbolCodes, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
// 		}, false)
// }

func (scatterTrigger *ScatterTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return scatterTrigger.Config.SymbolCodes
}

// CanTriggerWithScene -
func (scatterTrigger *ScatterTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	return scatterTrigger.canTrigger(gameProp, gs, nil, curpr, stake)
}

// CanTrigger -
func (scatterTrigger *ScatterTrigger) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, _ *sgc7game.GameScene, _ *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	isTrigger := false
	lst := []*sgc7game.Result{}

	symbols := scatterTrigger.getSymbols(gameProp)

	if scatterTrigger.Config.TriggerType == STTypeScatters {
		for _, s := range symbols {
			ret := sgc7game.CalcScatter4(gs, gameProp.CurPaytables, s, gameProp.GetBet2(stake, scatterTrigger.Config.BetType),
				func(scatter int, cursymbol int) bool {
					return cursymbol == scatter || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, false)

			if ret != nil {
				// if scatterTrigger.Config.BetType == BTypeNoPay {
				// 	ret.CoinWin = 0
				// 	ret.CashWin = 0
				// } else {
				// 	// gameProp.ProcMulti(ret)
				// }

				isTrigger = true

				lst = append(lst, ret)
			}
		}
	} else if scatterTrigger.Config.TriggerType == STTypeReelScatters {
		for _, s := range symbols {
			ret := sgc7game.CalcScatter4(gs, gameProp.CurPaytables, s, gameProp.GetBet2(stake, scatterTrigger.Config.BetType),
				func(scatter int, cursymbol int) bool {
					return cursymbol == scatter || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, true)

			if ret != nil {
				// if scatterTrigger.Config.BetType == BTypeNoPay {
				// 	ret.CoinWin = 0
				// 	ret.CashWin = 0
				// } else {
				// 	// gameProp.ProcMulti(ret)
				// }

				isTrigger = true

				lst = append(lst, ret)
			}
		}
	} else if scatterTrigger.Config.TriggerType == STTypeCountScatter {
		ret := sgc7game.CalcScatterEx(gs, symbols[0], scatterTrigger.Config.MinNum, func(scatter int, cursymbol int) bool {
			return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(scatterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		})

		if ret != nil {
			if scatterTrigger.Config.BetType != BTypeNoPay {
				// 	ret.CoinWin = 0
				// 	ret.CashWin = 0
				// } else {
				if scatterTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, scatterTrigger.Config.BetType)
				}

				// gameProp.ProcMulti(ret)
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
				// 	ret.CoinWin = 0
				// 	ret.CashWin = 0
				// } else {
				if scatterTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[scatterTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, scatterTrigger.Config.BetType)
				}

				// gameProp.ProcMulti(ret)
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

// // CanTrigger -
// func (scatterTrigger *ScatterTrigger) CanTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, isSaveResult bool) (bool, []*sgc7game.Result) {
// 	return scatterTrigger.canTrigger(gameProp, gs, curpr, stake, isSaveResult, scatterTrigger.getSymbols())
// }

// procWins
func (scatterTrigger *ScatterTrigger) procWins(gameProp *GameProperty, std *ScatterTriggerData, lst []*sgc7game.Result) (int, error) {
	if scatterTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0
		}

		return 0, nil
	}

	std.WinMulti = scatterTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin
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
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(scatterTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := scatterTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if scatterTrigger.Config.RespinNumWeightVW != nil {
		cr, err := scatterTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("ScatterTrigger.calcRespinNum:RespinNumWeightVW",
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if scatterTrigger.Config.RespinNum > 0 {
		return scatterTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// playgame
func (scatterTrigger *ScatterTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// scatterTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := icd.(*ScatterTriggerData)
	std.onNewStep()

	gs := scatterTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	os := scatterTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := scatterTrigger.canTrigger(gameProp, gs, os, curpr, stake)

	if isTrigger {
		scatterTrigger.procWins(gameProp, std, lst)

		// if !scatterTrigger.Config.NeedDiscardResults {
		for _, v := range lst {
			scatterTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds
		}
		// } else {
		// 	for _, v := range lst {
		// 		std.SymbolNum += v.SymbolNums
		// 		std.WildNum += v.Wilds
		// 	}
		// }

		respinNum, err := scatterTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", nil
		}

		std.RespinNum = respinNum

		err = scatterTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("ScatterTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		// if scatterTrigger.Config.TagSymbolNum != "" {
		// 	gameProp.TagInt(spSymbolTrigger.Config.TagSymbolNum, lst[0].SymbolNums)
		// }

		if len(scatterTrigger.Config.Awards) > 0 {
			gameProp.procAwards(plugin, scatterTrigger.Config.Awards, curpr, gp)
		}

		if scatterTrigger.Config.SymbolAwardsWeights != nil {
			for i := 0; i < lst[0].SymbolNums; i++ {
				node, err := scatterTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("ScatterTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
						goutils.Err(err))

					return "", err
				}

				gameProp.procAwards(plugin, node.Awards, curpr, gp)
			}
		}

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

			// if symbolTrigger.Config.RespinNumWeightWithScatterNum != nil {
			// 	v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, symbolTrigger.Config.RespinNumWeightWithScatterNum[lst[0].SymbolNums], symbolTrigger.Config.UseFileMapping, symbolTrigger.Config.JumpToComponent, true)
			// 	if err != nil {
			// 		goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
			// 			goutils.Err(err))

			// 		return nil
			// 	}

			// 	lst[0].Type = sgc7game.RTFreeGame
			// 	lst[0].Value = v
			// } else if len(symbolTrigger.Config.RespinNumWithScatterNum) > 0 {
			// 	gameProp.TriggerRespin(plugin, curpr, gp, symbolTrigger.Config.RespinNumWithScatterNum[lst[0].SymbolNums], symbolTrigger.Config.JumpToComponent, true)

			// 	lst[0].Type = sgc7game.RTFreeGame
			// 	lst[0].Value = symbolTrigger.Config.RespinNumWithScatterNum[lst[0].SymbolNums]
			// } else if symbolTrigger.Config.RespinNumWeight != "" {
			// 	v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, symbolTrigger.Config.RespinNumWeight, symbolTrigger.Config.UseFileMapping, symbolTrigger.Config.JumpToComponent, true)
			// 	if err != nil {
			// 		goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
			// 			goutils.Err(err))

			// 		return nil
			// 	}

			// 	lst[0].Type = sgc7game.RTFreeGame
			// 	lst[0].Value = v
			// } else if symbolTrigger.Config.RespinNum > 0 {
			// 	gameProp.TriggerRespin(plugin, curpr, gp, symbolTrigger.Config.RespinNum, symbolTrigger.Config.JumpToComponent, true)

			// 	lst[0].Type = sgc7game.RTFreeGame
			// 	lst[0].Value = symbolTrigger.Config.RespinNum
			// } else {
			// 	lst[0].Type = sgc7game.RTFreeGame
			// 	lst[0].Value = -1
			// }

			// if symbolTrigger.Config.ForceToNext {
			// 	std.NextComponent = symbolTrigger.Config.DefaultNextComponent
			// } else {
			// 	rn := gameProp.GetLastRespinNum(symbolTrigger.Config.JumpToComponent)
			// 	if rn > 0 {
			// 		gameProp.TriggerRespin(plugin, curpr, gp, 0, symbolTrigger.Config.JumpToComponent, true)

			// 		lst[0].Type = sgc7game.RTFreeGame
			// 		lst[0].Value = rn
			// 	}
			// }
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

// // OnStatsWithPB -
// func (scatterTrigger *ScatterTrigger) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.ScatterTriggerData)
// 	if !isok {
// 		goutils.Error("ScatterTrigger.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return scatterTrigger.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (scatterTrigger *ScatterTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	wins := int64(0)
// 	isTrigger := false

// 	for _, v := range lst {
// 		gp, isok := v.CurGameModParams.(*GameParams)
// 		if isok {
// 			curComponent, isok := gp.MapComponentMsgs[scatterTrigger.Name]
// 			if isok {
// 				curwins, err := scatterTrigger.OnStatsWithPB(feature, curComponent, v)
// 				if err != nil {
// 					goutils.Error("ScatterTrigger.OnStats",
// 						goutils.Err(err))

// 					continue
// 				}

// 				isTrigger = true
// 				wins += curwins
// 			}
// 		}
// 	}

// 	feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

// 	if feature.Parent != nil {
// 		totalwins := int64(0)

// 		for _, v := range lst {
// 			totalwins += v.CashWin
// 		}

// 		feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
// 	}

// 	return isTrigger, stake.CashBet, wins
// }

// NewComponentData -
func (scatterTrigger *ScatterTrigger) NewComponentData() IComponentData {
	return &ScatterTriggerData{}
}

func (scatterTrigger *ScatterTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(STCVWinMulti)
	if isok {
		return winMulti
	}

	return scatterTrigger.Config.WinMulti
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
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (scatterTrigger *ScatterTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache) {
	scatterTrigger.BasicComponent.OnStats2(icd, s2)

	cd := icd.(*ScatterTriggerData)

	s2.ProcStatsWins(scatterTrigger.Name, int64(cd.Wins))
}

// func (scatterTrigger *ScatterTrigger) getSymbols() []int {
// 	if scatterTrigger.dataForeachSymbol != nil {
// 		return []int{scatterTrigger.dataForeachSymbol.SymbolCode}
// 	}

// 	return scatterTrigger.Config.SymbolCodes
// }

// // CanTriggerWithScene -
// func (scatterTrigger *ScatterTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
// 	return scatterTrigger.canTrigger(gameProp, gs, nil, curpr, stake)
// }

// // OnEachSymbol - on foreach symbol
// func (scatterTrigger *ScatterTrigger) OnEachSymbol(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin, ps sgc7game.IPlayerState, stake *sgc7game.Stake,
// 	prs []*sgc7game.PlayResult, symbol int, cd IComponentData) (string, error) {

// 	std := cd.(*ScatterTriggerData)

// 	gs := scatterTrigger.GetTargetScene3(gameProp, curpr, prs, &std.BasicComponentData, scatterTrigger.Name, "", 0)

// 	isTrigger, lst := scatterTrigger.canTrigger(gameProp, gs, curpr, stake, !scatterTrigger.Config.NeedDiscardResults, []int{symbol})

// 	if isTrigger {
// 		scatterTrigger.procWins(std, lst)

// 		if !scatterTrigger.Config.NeedDiscardResults {
// 			for _, v := range lst {
// 				scatterTrigger.AddResult(curpr, v, &std.BasicComponentData)
// 			}
// 		}

// 		std.SymbolNum = lst[0].SymbolNums
// 		std.WildNum = lst[0].Wilds

// 		respinNum, err := scatterTrigger.calcRespinNum(plugin, lst[0])
// 		if err != nil {
// 			goutils.Error("ScatterTrigger.OnEachSymbol:calcRespinNum",
// 				goutils.Err(err))

// 			return "", nil
// 		}

// 		std.RespinNum = respinNum

// 		err = scatterTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
// 		if err != nil {
// 			goutils.Error("ScatterTrigger.OnPlayGame:procMask",
// 				goutils.Err(err))

// 			return "", err
// 		}

// 		// if scatterTrigger.Config.TagSymbolNum != "" {
// 		// 	gameProp.TagInt(spSymbolTrigger.Config.TagSymbolNum, lst[0].SymbolNums)
// 		// }

// 		if len(scatterTrigger.Config.Awards) > 0 {
// 			gameProp.procAwards(plugin, scatterTrigger.Config.Awards, curpr, gp)
// 		}

// 		if scatterTrigger.Config.SymbolAwardsWeights != nil {
// 			for i := 0; i < lst[0].SymbolNums; i++ {
// 				node, err := scatterTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
// 				if err != nil {
// 					goutils.Error("ScatterTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
// 						goutils.Err(err))

// 					return err
// 				}

// 				gameProp.procAwards(plugin, node.Awards, curpr, gp)
// 			}
// 		}

// 		if scatterTrigger.Config.JumpToComponent != "" {
// 			if gameProp.IsRespin(scatterTrigger.Config.JumpToComponent) {
// 				// 如果jumpto是一个respin，那么就需要trigger respin
// 				if std.RespinNum == 0 {
// 					if scatterTrigger.Config.ForceToNext {
// 						std.NextComponent = scatterTrigger.Config.DefaultNextComponent
// 					} else {
// 						rn := gameProp.GetLastRespinNum(scatterTrigger.Config.JumpToComponent)
// 						if rn > 0 {
// 							gameProp.TriggerRespin(plugin, curpr, gp, 0, scatterTrigger.Config.JumpToComponent, !scatterTrigger.Config.IsAddRespinMode)

// 							lst[0].Type = sgc7game.RTFreeGame
// 							lst[0].Value = rn
// 						}
// 					}
// 				} else {
// 					// 如果jumpto是respin，需要treigger这个respin
// 					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, scatterTrigger.Config.JumpToComponent, !scatterTrigger.Config.IsAddRespinMode)

// 					lst[0].Type = sgc7game.RTFreeGame
// 					lst[0].Value = std.RespinNum
// 				}
// 			}

// 			// if symbolTrigger.Config.RespinNumWeightWithScatterNum != nil {
// 			// 	v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, symbolTrigger.Config.RespinNumWeightWithScatterNum[lst[0].SymbolNums], symbolTrigger.Config.UseFileMapping, symbolTrigger.Config.JumpToComponent, true)
// 			// 	if err != nil {
// 			// 		goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
// 			// 			goutils.Err(err))

// 			// 		return nil
// 			// 	}

// 			// 	lst[0].Type = sgc7game.RTFreeGame
// 			// 	lst[0].Value = v
// 			// } else if len(symbolTrigger.Config.RespinNumWithScatterNum) > 0 {
// 			// 	gameProp.TriggerRespin(plugin, curpr, gp, symbolTrigger.Config.RespinNumWithScatterNum[lst[0].SymbolNums], symbolTrigger.Config.JumpToComponent, true)

// 			// 	lst[0].Type = sgc7game.RTFreeGame
// 			// 	lst[0].Value = symbolTrigger.Config.RespinNumWithScatterNum[lst[0].SymbolNums]
// 			// } else if symbolTrigger.Config.RespinNumWeight != "" {
// 			// 	v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, symbolTrigger.Config.RespinNumWeight, symbolTrigger.Config.UseFileMapping, symbolTrigger.Config.JumpToComponent, true)
// 			// 	if err != nil {
// 			// 		goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
// 			// 			goutils.Err(err))

// 			// 		return nil
// 			// 	}

// 			// 	lst[0].Type = sgc7game.RTFreeGame
// 			// 	lst[0].Value = v
// 			// } else if symbolTrigger.Config.RespinNum > 0 {
// 			// 	gameProp.TriggerRespin(plugin, curpr, gp, symbolTrigger.Config.RespinNum, symbolTrigger.Config.JumpToComponent, true)

// 			// 	lst[0].Type = sgc7game.RTFreeGame
// 			// 	lst[0].Value = symbolTrigger.Config.RespinNum
// 			// } else {
// 			// 	lst[0].Type = sgc7game.RTFreeGame
// 			// 	lst[0].Value = -1
// 			// }

// 			// if symbolTrigger.Config.ForceToNext {
// 			// 	std.NextComponent = symbolTrigger.Config.DefaultNextComponent
// 			// } else {
// 			// 	rn := gameProp.GetLastRespinNum(symbolTrigger.Config.JumpToComponent)
// 			// 	if rn > 0 {
// 			// 		gameProp.TriggerRespin(plugin, curpr, gp, 0, symbolTrigger.Config.JumpToComponent, true)

// 			// 		lst[0].Type = sgc7game.RTFreeGame
// 			// 		lst[0].Value = rn
// 			// 	}
// 			// }

// 			std.NextComponent = scatterTrigger.Config.JumpToComponent

// 			scatterTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

// 			return nil
// 		}
// 	}

// 	scatterTrigger.onStepEnd(gameProp, curpr, gp, "")
// }

func NewScatterTrigger(name string) IComponent {
	return &ScatterTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
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
//	},
type jsonScatterTrigger struct {
	Symbols                       []string   `json:"symbols"`
	TriggerType                   string     `json:"triggerType"`
	BetType                       string     `json:"betType"`
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
		// RespinNumWithScatterNum:       jcfg.RespinNumWithScatterNum,
		// RespinNumWeightWithScatterNum: jcfg.RespinNumWeightWithScatterNum,
	}

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

	for i := range cfg.PosArea {
		cfg.PosArea[i]--
	}

	// cfg.UseSceneV3 = true

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
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseScatterTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
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
