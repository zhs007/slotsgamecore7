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

const WaysTriggerTypeName = "waysTrigger"

const (
	WTCVWinMulti string = "winMulti" // 可以修改配置项里的winMulti
)

type WaysTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
	Wins          int
	WinMulti      int
}

// OnNewGame -
func (waysTriggerData *WaysTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	waysTriggerData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (waysTriggerData *WaysTriggerData) onNewStep() {
	// waysTriggerData.BasicComponentData.OnNewStep(gameProp, component)

	waysTriggerData.UsedResults = nil

	waysTriggerData.NextComponent = ""
	waysTriggerData.SymbolNum = 0
	waysTriggerData.WildNum = 0
	waysTriggerData.RespinNum = 0
	waysTriggerData.Wins = 0
	waysTriggerData.WinMulti = 1
}

// BuildPBComponentData
func (waysTriggerData *WaysTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WaysTriggerData{
		BasicComponentData: waysTriggerData.BuildPBBasicComponentData(),
		NextComponent:      waysTriggerData.NextComponent,
		SymbolNum:          int32(waysTriggerData.SymbolNum),
		WildNum:            int32(waysTriggerData.WildNum),
		RespinNum:          int32(waysTriggerData.RespinNum),
		Wins:               int32(waysTriggerData.Wins),
		WinMulti:           int32(waysTriggerData.WinMulti),
	}

	return pbcd
}

// GetVal -
func (waysTriggerData *WaysTriggerData) GetVal(key string) (int, bool) {
	if key == CVSymbolNum {
		return waysTriggerData.SymbolNum, true
	} else if key == CVWildNum {
		return waysTriggerData.WildNum, true
	} else if key == CVRespinNum {
		return waysTriggerData.RespinNum, true
	} else if key == CVWins {
		return waysTriggerData.Wins, true
	}

	return 0, false
}

// SetVal -
func (waysTriggerData *WaysTriggerData) SetVal(key string, val int) {
	if key == CVSymbolNum {
		waysTriggerData.SymbolNum = val
	} else if key == CVWildNum {
		waysTriggerData.WildNum = val
	} else if key == CVRespinNum {
		waysTriggerData.RespinNum = val
	} else if key == CVWins {
		waysTriggerData.Wins = val
	}
}

// WaysTriggerConfig - configuration for WaysTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WaysTriggerConfig struct {
	BasicComponentConfig            `yaml:",inline" json:",inline"`
	Symbols                         []string                      `yaml:"symbols" json:"symbols"`                                             // like scatter
	SymbolCodes                     []int                         `yaml:"-" json:"-"`                                                         // like scatter
	Type                            string                        `yaml:"type" json:"type"`                                                   // like scatters
	TriggerType                     SymbolTriggerType             `yaml:"-" json:"-"`                                                         // SymbolTriggerType
	BetTypeString                   string                        `yaml:"betType" json:"betType"`                                             // bet or totalBet or noPay
	BetType                         BetType                       `yaml:"-" json:"-"`                                                         // bet or totalBet or noPay
	OSMulTypeString                 string                        `yaml:"symbolValsMulti" json:"symbolValsMulti"`                             // OtherSceneMultiType
	OSMulType                       OtherSceneMultiType           `yaml:"-" json:"-"`                                                         // OtherSceneMultiType
	MinNum                          int                           `yaml:"minNum" json:"minNum"`                                               // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols                     []string                      `yaml:"wildSymbols" json:"wildSymbols"`                                     // wild etc
	WildSymbolCodes                 []int                         `yaml:"-" json:"-"`                                                         // wild symbolCode
	StrCheckWinType                 string                        `yaml:"checkWinType" json:"checkWinType"`                                   // left2right or right2left or all
	CheckWinType                    CheckWinType                  `yaml:"-" json:"-"`                                                         //
	WinMulti                        int                           `yaml:"winMulti" json:"winMulti"`                                           // winMulti，最后的中奖倍数，默认为1
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`                             // jump to
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`                                     // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                                               // 新的奖励系统
	SymbolAwardsWeights             *AwardsWeights                `yaml:"symbolAwardsWeights" json:"symbolAwardsWeights"`                     // 每个中奖符号随机一组奖励
	TargetMask                      string                        `yaml:"targetMask" json:"targetMask"`                                       // 如果是scatter这一组判断，可以把结果传递给一个mask
	IsReverse                       bool                          `yaml:"isReverse" json:"isReverse"`                                         // 如果isReverse，表示判定为否才触发
	PiggyBankComponent              string                        `yaml:"piggyBankComponent" json:"piggyBankComponent"`                       // piggyBank component
	IsAddRespinMode                 bool                          `yaml:"isAddRespinMode" json:"isAddRespinMode"`                             // 是否是增加respinNum模式，默认是增加triggerNum模式
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
}

// SetLinkComponent
func (cfg *WaysTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type WaysTrigger struct {
	*BasicComponent `json:"-"`
	Config          *WaysTriggerConfig `json:"config"`
}

// Init -
func (waysTrigger *WaysTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WaysTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WaysTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WaysTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return waysTrigger.InitEx(cfg, pool)
}

// InitEx -
func (waysTrigger *WaysTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	waysTrigger.Config = cfg.(*WaysTriggerConfig)
	waysTrigger.Config.ComponentType = WaysTriggerTypeName

	waysTrigger.Config.OSMulType = ParseOtherSceneMultiType(waysTrigger.Config.OSMulTypeString)

	for _, s := range waysTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		waysTrigger.Config.SymbolCodes = append(waysTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range waysTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		waysTrigger.Config.WildSymbolCodes = append(waysTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(waysTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("WaysTrigger.InitEx:ParseSymbolTriggerType",
			slog.String("SymbolTriggerType", waysTrigger.Config.Type),
			goutils.Err(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	waysTrigger.Config.TriggerType = stt

	waysTrigger.Config.BetType = ParseBetType(waysTrigger.Config.BetTypeString)

	for _, award := range waysTrigger.Config.Awards {
		award.Init()
	}

	if waysTrigger.Config.SymbolAwardsWeights != nil {
		waysTrigger.Config.SymbolAwardsWeights.Init()
	}

	// waysTrigger.Config.ExcludeSymbolCodes = GetExcludeSymbols(pool.DefaultPaytables, waysTrigger.Config.SymbolCodes)

	waysTrigger.Config.CheckWinType = ParseCheckWinType(waysTrigger.Config.StrCheckWinType)

	if waysTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(waysTrigger.Config.RespinNumWeight, waysTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
				slog.String("Weight", waysTrigger.Config.RespinNumWeight),
				goutils.Err(err))

			return err
		}

		waysTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(waysTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range waysTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, waysTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))

				return err
			}

			waysTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if waysTrigger.Config.WinMulti <= 0 {
		waysTrigger.Config.WinMulti = 1
	}

	// if waysTrigger.Config.BetType == BTypeNoPay {
	// 	waysTrigger.Config.NeedDiscardResults = true
	// }

	waysTrigger.onInit(&waysTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (waysTrigger *WaysTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if waysTrigger.Config.TargetMask != "" {
		gameProp.UseComponent(waysTrigger.Config.TargetMask)

		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, waysTrigger.Config.TargetMask, mask, false)
	}

	return nil
}

func (waysTrigger *WaysTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return waysTrigger.Config.SymbolCodes
}

// CanTriggerWithScene -
func (waysTrigger *WaysTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	return waysTrigger.canTrigger(gameProp, gs, nil, curpr, stake)
}

// CanTrigger -
func (waysTrigger *WaysTrigger) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, _ *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	// std := gameProp.MapComponentData[waysTrigger.Name].(*WaysTriggerData)

	isTrigger := false
	lst := []*sgc7game.Result{}
	symbols := waysTrigger.getSymbols(gameProp)

	if waysTrigger.Config.TriggerType == STTypeWays {
		// os := waysTrigger.GetTargetOtherScene2(gameProp, curpr, &std.BasicComponentData, waysTrigger.Name, "")

		if os != nil {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, waysTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					// return true
					return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(x, y int) int {
					return os.Arr[x][y]
				})

			// for _, v := range currets {
			// 	gameProp.ProcMulti(v)

			// 	// if isSaveResult {
			// 	// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
			// 	// }
			// }

			lst = append(lst, currets...)
		} else {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, waysTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(x, y int) int {
					return 1
				})

			// for _, v := range currets {
			// 	gameProp.ProcMulti(v)

			// 	// if isSaveResult {
			// 	// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
			// 	// }
			// }

			lst = append(lst, currets...)
		}

		if len(lst) > 0 {
			isTrigger = true
		}
	} else if waysTrigger.Config.TriggerType == STTypeCheckWays {
		currets := sgc7game.CheckWays(gs, waysTrigger.Config.MinNum,
			func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
				return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		// for _, v := range currets {
		// 	// gameProp.ProcMulti(v)

		// 	// if isSaveResult {
		// 	// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
		// 	// }
		// }

		lst = append(lst, currets...)

		if len(lst) > 0 {
			isTrigger = true
		}
	}

	if waysTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

// procWins
func (waysTrigger *WaysTrigger) procWins(gameProp *GameProperty, std *WaysTriggerData, lst []*sgc7game.Result) (int, error) {
	if waysTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0
		}

		return 0, nil
	}

	std.WinMulti = waysTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin
	}

	if std.Wins > 0 {
		if waysTrigger.Config.PiggyBankComponent != "" {
			cd := gameProp.GetCurComponentDataWithName(waysTrigger.Config.PiggyBankComponent)
			if cd == nil {
				goutils.Error("ScatterTrigger.procWins:GetCurComponentDataWithName",
					slog.String("PiggyBankComponent", waysTrigger.Config.PiggyBankComponent),
					goutils.Err(ErrInvalidComponent))

				return 0, ErrInvalidComponent
			}

			cd.ChgConfigIntVal(CCVSavedMoney, std.Wins)

			gameProp.UseComponent(waysTrigger.Config.PiggyBankComponent)
		}
	}

	return std.Wins, nil
}

// calcRespinNum
func (waysTrigger *WaysTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(waysTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := waysTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					slog.Int("SymbolNum", ret.SymbolNums),
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(waysTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := waysTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if waysTrigger.Config.RespinNumWeightVW != nil {
		cr, err := waysTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightVW",
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if waysTrigger.Config.RespinNum > 0 {
		return waysTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// playgame
func (waysTrigger *WaysTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// waysTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := icd.(*WaysTriggerData)
	std.onNewStep()

	gs := waysTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	os := waysTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := waysTrigger.canTrigger(gameProp, gs, os, curpr, stake)

	if isTrigger {
		waysTrigger.procWins(gameProp, std, lst)

		// if !waysTrigger.Config.NeedDiscardResults {
		for _, v := range lst {
			waysTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds
		}
		// } else {
		// 	for _, v := range lst {
		// 		std.SymbolNum += v.SymbolNums
		// 		std.WildNum += v.Wilds
		// 	}
		// }

		respinNum, err := waysTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", nil
		}

		std.RespinNum = respinNum

		err = waysTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		// if symbolTrigger.Config.TagSymbolNum != "" {
		// 	gameProp.TagInt(symbolTrigger.Config.TagSymbolNum, lst[0].SymbolNums)
		// }

		if len(waysTrigger.Config.Awards) > 0 {
			gameProp.procAwards(plugin, waysTrigger.Config.Awards, curpr, gp)
		}

		if waysTrigger.Config.SymbolAwardsWeights != nil {
			for i := 0; i < lst[0].SymbolNums; i++ {
				node, err := waysTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("WaysTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
						goutils.Err(err))

					return "", err
				}

				gameProp.procAwards(plugin, node.Awards, curpr, gp)
			}
		}

		if waysTrigger.Config.JumpToComponent != "" {
			if gameProp.IsRespin(waysTrigger.Config.JumpToComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if waysTrigger.Config.ForceToNext {
						std.NextComponent = waysTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(waysTrigger.Config.JumpToComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, waysTrigger.Config.JumpToComponent, !waysTrigger.Config.IsAddRespinMode)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, waysTrigger.Config.JumpToComponent, !waysTrigger.Config.IsAddRespinMode)

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

			std.NextComponent = waysTrigger.Config.JumpToComponent

			nc := waysTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := waysTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := waysTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (waysTrigger *WaysTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WaysTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", waysTrigger.Name, std.NextComponent)
	}

	return nil
}

// // OnStatsWithPB -
// func (waysTrigger *WaysTrigger) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.WaysTriggerData)
// 	if !isok {
// 		goutils.Error("WaysTrigger.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return waysTrigger.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (waysTrigger *WaysTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	wins := int64(0)
// 	isTrigger := false

// 	for _, v := range lst {
// 		gp, isok := v.CurGameModParams.(*GameParams)
// 		if isok {
// 			curComponent, isok := gp.MapComponentMsgs[waysTrigger.Name]
// 			if isok {
// 				curwins, err := waysTrigger.OnStatsWithPB(feature, curComponent, v)
// 				if err != nil {
// 					goutils.Error("WaysTrigger.OnStats",
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
func (waysTrigger *WaysTrigger) NewComponentData() IComponentData {
	return &WaysTriggerData{}
}

func (waysTrigger *WaysTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(WTCVWinMulti)
	if isok {
		return winMulti
	}

	return waysTrigger.Config.WinMulti
}

// NewStats2 -
func (waysTrigger *WaysTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (waysTrigger *WaysTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache) {
	waysTrigger.BasicComponent.OnStats2(icd, s2)

	cd := icd.(*WaysTriggerData)

	s2.ProcStatsWins(waysTrigger.Name, int64(cd.Wins))
}

// GetAllLinkComponents - get all link components
func (waysTrigger *WaysTrigger) GetAllLinkComponents() []string {
	return []string{waysTrigger.Config.DefaultNextComponent, waysTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (waysTrigger *WaysTrigger) GetNextLinkComponents() []string {
	return []string{waysTrigger.Config.DefaultNextComponent, waysTrigger.Config.JumpToComponent}
}

func NewWaysTrigger(name string) IComponent {
	return &WaysTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"triggerType": "lines",
//		"betType": "bet",
//		"checkWinType": "left2right",
//		"symbols": [
//			"WL",
//			"A",
//			"B",
//			"C",
//			"D",
//			"E",
//			"F",
//			"G",
//			"H",
//			"J",
//			"K",
//			"L"
//		],
//		"wildSymbols": [
//			"WL"
//		]
//	},
type jsonWaysTrigger struct {
	Symbols             []string `json:"symbols"`
	TriggerType         string   `json:"triggerType"`
	BetType             string   `json:"betType"`
	SymbolValsMulti     string   `json:"symbolValsMulti"`
	MinNum              int      `json:"minNum"`
	WildSymbols         []string `json:"wildSymbols"`
	WinMulti            int      `json:"winMulti"`
	PutMoneyInPiggyBank string   `json:"putMoneyInPiggyBank"`
}

func (jcfg *jsonWaysTrigger) build() *WaysTriggerConfig {
	cfg := &WaysTriggerConfig{
		Symbols:            jcfg.Symbols,
		Type:               jcfg.TriggerType,
		BetTypeString:      jcfg.BetType,
		MinNum:             jcfg.MinNum,
		WildSymbols:        jcfg.WildSymbols,
		WinMulti:           jcfg.WinMulti,
		PiggyBankComponent: jcfg.PutMoneyInPiggyBank,
		OSMulTypeString:    jcfg.SymbolValsMulti,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWaysTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWaysTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWaysTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWaysTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWaysTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseWaysTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WaysTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
