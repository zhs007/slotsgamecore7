package lowcode

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
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

// OnNewStep -
func (waysTriggerData *WaysTriggerData) OnNewStep(gameProp *GameProperty, component IComponent) {
	waysTriggerData.BasicComponentData.OnNewStep(gameProp, component)

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
func (waysTriggerData *WaysTriggerData) GetVal(key string) int {
	if key == STDVSymbolNum {
		return waysTriggerData.SymbolNum
	} else if key == STDVWildNum {
		return waysTriggerData.WildNum
	} else if key == STDVRespinNum {
		return waysTriggerData.RespinNum
	} else if key == STDVWins {
		return waysTriggerData.Wins
	}

	return 0
}

// SetVal -
func (waysTriggerData *WaysTriggerData) SetVal(key string, val int) {
	if key == STDVSymbolNum {
		waysTriggerData.SymbolNum = val
	} else if key == STDVWildNum {
		waysTriggerData.WildNum = val
	} else if key == STDVRespinNum {
		waysTriggerData.RespinNum = val
	} else if key == STDVWins {
		waysTriggerData.Wins = val
	}
}

// WaysTriggerConfig - configuration for WaysTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WaysTriggerConfig struct {
	BasicComponentConfig            `yaml:",inline" json:",inline"`
	Symbols                         []string                      `yaml:"symbols" json:"symbols"`                                             // like scatter
	SymbolCodes                     []int                         `yaml:"-" json:"-"`                                                         // like scatter
	ExcludeSymbolCodes              []int                         `yaml:"-" json:"-"`                                                         // 在 lines 和 ways 里有用
	Type                            string                        `yaml:"type" json:"type"`                                                   // like scatters
	TriggerType                     SymbolTriggerType             `yaml:"-" json:"-"`                                                         // SymbolTriggerType
	BetTypeString                   string                        `yaml:"betType" json:"betType"`                                             // bet or totalBet or noPay
	BetType                         BetType                       `yaml:"-" json:"-"`                                                         // bet or totalBet or noPay
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
	NeedDiscardResults              bool                          `yaml:"needDiscardResults" json:"needDiscardResults"`                       // 如果needDiscardResults，表示抛弃results
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WaysTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WaysTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return waysTrigger.InitEx(cfg, pool)
}

// InitEx -
func (waysTrigger *WaysTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	waysTrigger.Config = cfg.(*WaysTriggerConfig)
	waysTrigger.Config.ComponentType = WaysTriggerTypeName

	for _, s := range waysTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:Symbol",
				zap.String("symbol", s),
				zap.Error(ErrIvalidSymbol))
		}

		waysTrigger.Config.SymbolCodes = append(waysTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range waysTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WaysTrigger.InitEx:WildSymbols",
				zap.String("symbol", s),
				zap.Error(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		waysTrigger.Config.WildSymbolCodes = append(waysTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(waysTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("WaysTrigger.InitEx:ParseSymbolTriggerType",
			zap.String("SymbolTriggerType", waysTrigger.Config.Type),
			zap.Error(ErrIvalidSymbolTriggerType))

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

	waysTrigger.Config.ExcludeSymbolCodes = GetExcludeSymbols(pool.DefaultPaytables, waysTrigger.Config.SymbolCodes)

	waysTrigger.Config.CheckWinType = ParseCheckWinType(waysTrigger.Config.StrCheckWinType)

	if waysTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(waysTrigger.Config.RespinNumWeight, waysTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
				zap.String("Weight", waysTrigger.Config.RespinNumWeight),
				zap.Error(err))

			return err
		}

		waysTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(waysTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range waysTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, waysTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("WaysTrigger.InitEx:LoadIntWeights",
					zap.String("Weight", v),
					zap.Error(err))

				return err
			}

			waysTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if waysTrigger.Config.WinMulti <= 0 {
		waysTrigger.Config.WinMulti = 1
	}

	if waysTrigger.Config.BetType == BTypeNoPay {
		waysTrigger.Config.NeedDiscardResults = true
	}

	waysTrigger.onInit(&waysTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (waysTrigger *WaysTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if waysTrigger.Config.TargetMask != "" {
		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, waysTrigger.Config.TargetMask, mask, false)
	}

	return nil
}

// CanTrigger -
func (waysTrigger *WaysTrigger) CanTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, isSaveResult bool) (bool, []*sgc7game.Result) {
	std := gameProp.MapComponentData[waysTrigger.Name].(*WaysTriggerData)

	isTrigger := false
	lst := []*sgc7game.Result{}

	if waysTrigger.Config.TriggerType == STTypeWays {
		os := waysTrigger.GetTargetOtherScene2(gameProp, curpr, &std.BasicComponentData, waysTrigger.Name, "")

		if os != nil {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, waysTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
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

			for _, v := range currets {
				gameProp.ProcMulti(v)

				// if isSaveResult {
				// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
				// }
			}

			lst = append(lst, currets...)
		} else {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, waysTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(waysTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
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

			for _, v := range currets {
				gameProp.ProcMulti(v)

				// if isSaveResult {
				// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
				// }
			}

			lst = append(lst, currets...)
		}

		if len(lst) > 0 {
			isTrigger = true
		}
	} else if waysTrigger.Config.TriggerType == STTypeCheckWays {
		currets := sgc7game.CheckWays(gs, waysTrigger.Config.MinNum,
			func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
				return goutils.IndexOfIntSlice(waysTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(waysTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		for _, v := range currets {
			gameProp.ProcMulti(v)

			// if isSaveResult {
			// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
			// }
		}

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
func (waysTrigger *WaysTrigger) procWins(std *WaysTriggerData, lst []*sgc7game.Result) (int, error) {
	std.WinMulti = waysTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin
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
					zap.Int("SymbolNum", ret.SymbolNums),
					zap.Error(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				zap.Int("SymbolNum", ret.SymbolNums),
				zap.Error(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(waysTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := waysTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWithScatterNum",
				zap.Int("SymbolNum", ret.SymbolNums),
				zap.Error(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if waysTrigger.Config.RespinNumWeightVW != nil {
		cr, err := waysTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("WaysTrigger.calcRespinNum:RespinNumWeightVW",
				zap.Error(err))

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	waysTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := gameProp.MapComponentData[waysTrigger.Name].(*WaysTriggerData)

	gs := waysTrigger.GetTargetScene3(gameProp, curpr, prs, &std.BasicComponentData, waysTrigger.Name, "", 0)

	isTrigger, lst := waysTrigger.CanTrigger(gameProp, gs, curpr, stake, !waysTrigger.Config.NeedDiscardResults)

	if isTrigger {
		waysTrigger.procWins(std, lst)

		if !waysTrigger.Config.NeedDiscardResults {
			for _, v := range lst {
				waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
			}
		}

		std.SymbolNum = lst[0].SymbolNums
		std.WildNum = lst[0].Wilds

		respinNum, err := waysTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:calcRespinNum",
				zap.Error(err))

			return nil
		}

		std.RespinNum = respinNum

		err = waysTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("WaysTrigger.OnPlayGame:procMask",
				zap.Error(err))

			return err
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
						zap.Error(err))

					return err
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
			// 			zap.Error(err))

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
			// 			zap.Error(err))

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

			waysTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nil
		}
	}

	waysTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (waysTrigger *WaysTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	std := gameProp.MapComponentData[waysTrigger.Name].(*WaysTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", waysTrigger.Name, std.NextComponent)
	}

	return nil
}

// OnStatsWithPB -
func (waysTrigger *WaysTrigger) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.WaysTriggerData)
	if !isok {
		goutils.Error("WaysTrigger.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return waysTrigger.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// OnStats
func (waysTrigger *WaysTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	wins := int64(0)
	isTrigger := false

	for _, v := range lst {
		gp, isok := v.CurGameModParams.(*GameParams)
		if isok {
			curComponent, isok := gp.MapComponentMsgs[waysTrigger.Name]
			if isok {
				curwins, err := waysTrigger.OnStatsWithPB(feature, curComponent, v)
				if err != nil {
					goutils.Error("WaysTrigger.OnStats",
						zap.Error(err))

					continue
				}

				isTrigger = true
				wins += curwins
			}
		}
	}

	feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

	if feature.Parent != nil {
		totalwins := int64(0)

		for _, v := range lst {
			totalwins += v.CashWin
		}

		feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
	}

	return isTrigger, stake.CashBet, wins
}

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

// GetAllLinkComponents - get all link components
func (waysTrigger *WaysTrigger) GetAllLinkComponents() []string {
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
	Symbols     []string `json:"symbols"`
	TriggerType string   `json:"triggerType"`
	BetType     string   `json:"betType"`
	MinNum      int      `json:"minNum"`
	WildSymbols []string `json:"wildSymbols"`
	WinMulti    int      `json:"winMulti"`
}

func (jwt *jsonWaysTrigger) build() *WaysTriggerConfig {
	cfg := &WaysTriggerConfig{
		Symbols:       jwt.Symbols,
		Type:          jwt.TriggerType,
		BetTypeString: jwt.BetType,
		MinNum:        jwt.MinNum,
		WildSymbols:   jwt.WildSymbols,
		WinMulti:      jwt.WinMulti,
	}

	cfg.UseSceneV3 = true

	return cfg
}

func parseWaysTrigger(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWaysTrigger:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWaysTrigger:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonWaysTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWaysTrigger:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(gamecfg, ctrls)
		if err != nil {
			goutils.Error("parseWaysTrigger:parseControllers",
				zap.Error(err))

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

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
