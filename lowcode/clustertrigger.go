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

const ClusterTriggerTypeName = "clusterTrigger"

const (
	CTCVWinMulti string = "winMulti" // 可以修改配置项里的winMulti
)

type ClusterTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
	Wins          int
	WinMulti      int
}

// OnNewGame -
func (clusterTriggerData *ClusterTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	clusterTriggerData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (clusterTriggerData *ClusterTriggerData) onNewStep() {
	clusterTriggerData.UsedResults = nil
	clusterTriggerData.NextComponent = ""
	clusterTriggerData.SymbolNum = 0
	clusterTriggerData.WildNum = 0
	clusterTriggerData.RespinNum = 0
	clusterTriggerData.Wins = 0
	clusterTriggerData.WinMulti = 1
}

// BuildPBComponentData
func (clusterTriggerData *ClusterTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ClusterTriggerData{
		BasicComponentData: clusterTriggerData.BuildPBBasicComponentData(),
		NextComponent:      clusterTriggerData.NextComponent,
		SymbolNum:          int32(clusterTriggerData.SymbolNum),
		WildNum:            int32(clusterTriggerData.WildNum),
		RespinNum:          int32(clusterTriggerData.RespinNum),
		Wins:               int32(clusterTriggerData.Wins),
		WinMulti:           int32(clusterTriggerData.WinMulti),
	}

	return pbcd
}

// GetVal -
func (clusterTriggerData *ClusterTriggerData) GetVal(key string) int {
	if key == CVSymbolNum {
		return clusterTriggerData.SymbolNum
	} else if key == CVWildNum {
		return clusterTriggerData.WildNum
	} else if key == CVRespinNum {
		return clusterTriggerData.RespinNum
	} else if key == CVWins {
		return clusterTriggerData.Wins
	}

	return 0
}

// SetVal -
func (clusterTriggerData *ClusterTriggerData) SetVal(key string, val int) {
	if key == CVSymbolNum {
		clusterTriggerData.SymbolNum = val
	} else if key == CVWildNum {
		clusterTriggerData.WildNum = val
	} else if key == CVRespinNum {
		clusterTriggerData.RespinNum = val
	} else if key == CVWins {
		clusterTriggerData.Wins = val
	}
}

// ClusterTriggerConfig - configuration for ClusterTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type ClusterTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string `yaml:"symbols" json:"symbols"` // like scatter
	SymbolCodes          []int    `yaml:"-" json:"-"`             // like scatter
	// ExcludeSymbolCodes              []int                         `yaml:"-" json:"-"`                                                         // 在 lines 和 ways 里有用
	Type                            string                        `yaml:"type" json:"type"`                                                   // like scatters
	TriggerType                     SymbolTriggerType             `yaml:"-" json:"-"`                                                         // SymbolTriggerType
	BetTypeString                   string                        `yaml:"betType" json:"betType"`                                             // bet or totalBet or noPay
	BetType                         BetType                       `yaml:"-" json:"-"`                                                         // bet or totalBet or noPay
	MinNum                          int                           `yaml:"minNum" json:"minNum"`                                               // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols                     []string                      `yaml:"wildSymbols" json:"wildSymbols"`                                     // wild etc
	WildSymbolCodes                 []int                         `yaml:"-" json:"-"`                                                         // wild symbolCode
	WinMulti                        int                           `yaml:"winMulti" json:"winMulti"`                                           // winMulti，最后的中奖倍数，默认为1
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`                             // jump to
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`                                     // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                                               // 新的奖励系统
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
func (cfg *ClusterTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type ClusterTrigger struct {
	*BasicComponent `json:"-"`
	Config          *ClusterTriggerConfig `json:"config"`
}

// Init -
func (clusterTrigger *ClusterTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ClusterTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ClusterTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ClusterTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return clusterTrigger.InitEx(cfg, pool)
}

// InitEx -
func (clusterTrigger *ClusterTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	clusterTrigger.Config = cfg.(*ClusterTriggerConfig)
	clusterTrigger.Config.ComponentType = ClusterTriggerTypeName

	for _, s := range clusterTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ClusterTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		clusterTrigger.Config.SymbolCodes = append(clusterTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range clusterTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ClusterTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		clusterTrigger.Config.WildSymbolCodes = append(clusterTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(clusterTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("ClusterTrigger.InitEx:ParseSymbolTriggerType",
			slog.String("SymbolTriggerType", clusterTrigger.Config.Type),
			goutils.Err(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	clusterTrigger.Config.TriggerType = stt

	clusterTrigger.Config.BetType = ParseBetType(clusterTrigger.Config.BetTypeString)

	for _, award := range clusterTrigger.Config.Awards {
		award.Init()
	}

	// if clusterTrigger.Config.SymbolAwardsWeights != nil {
	// 	clusterTrigger.Config.SymbolAwardsWeights.Init()
	// }

	// clusterTrigger.Config.ExcludeSymbolCodes = GetExcludeSymbols(pool.DefaultPaytables, clusterTrigger.Config.SymbolCodes)

	// clusterTrigger.Config.CheckWinType = ParseCheckWinType(clusterTrigger.Config.StrCheckWinType)

	if clusterTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(clusterTrigger.Config.RespinNumWeight, clusterTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ClusterTrigger.InitEx:LoadIntWeights",
				slog.String("Weight", clusterTrigger.Config.RespinNumWeight),
				goutils.Err(err))

			return err
		}

		clusterTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(clusterTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range clusterTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, clusterTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("ClusterTrigger.InitEx:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))

				return err
			}

			clusterTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if clusterTrigger.Config.WinMulti <= 0 {
		clusterTrigger.Config.WinMulti = 1
	}

	// if clusterTrigger.Config.BetType == BTypeNoPay {
	// 	clusterTrigger.Config.NeedDiscardResults = true
	// }

	clusterTrigger.onInit(&clusterTrigger.Config.BasicComponentConfig)

	return nil
}

// // playgame
// func (clusterTrigger *ClusterTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
// 	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

// 	if clusterTrigger.Config.TargetMask != "" {
// 		mask := make([]bool, gs.Width)

// 		for i := 0; i < len(ret.Pos)/2; i++ {
// 			mask[ret.Pos[i*2]] = true
// 		}

// 		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, clusterTrigger.Config.TargetMask, mask, false)
// 	}

// 	return nil
// }

// // CanTrigger -
// func (clusterTrigger *ClusterTrigger) CanTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, isSaveResult bool) (bool, []*sgc7game.Result) {
// 	// std := gameProp.MapComponentData[clusterTrigger.Name].(*ClusterTriggerData)

// 	isTrigger := false
// 	lst := []*sgc7game.Result{}

// 	if clusterTrigger.Config.TriggerType == STTypeCluster {

// 		currets, err := sgc7game.CalcClusterResult(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, clusterTrigger.Config.BetType),
// 			func(cursymbol int) bool {
// 				return goutils.IndexOfIntSlice(clusterTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
// 			}, func(cursymbol int) bool {
// 				return goutils.IndexOfIntSlice(clusterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
// 			}, func(cursymbol int, startsymbol int) bool {
// 				if cursymbol == startsymbol {
// 					return true
// 				}

// 				return goutils.IndexOfIntSlice(clusterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
// 			}, func(cursymbol int) int {
// 				return cursymbol
// 			})
// 		if err != nil {
// 			goutils.Error("ClusterTrigger.CanTrigger:CalcClusterResult",
// 				goutils.Err(err))

// 			return false, nil
// 		}

// 		for _, v := range currets {
// 			gameProp.ProcMulti(v)

// 			// if isSaveResult {
// 			// 	waysTrigger.AddResult(curpr, v, &std.BasicComponentData)
// 			// }
// 		}

// 		lst = append(lst, currets...)

// 		if len(lst) > 0 {
// 			isTrigger = true
// 		}
// 	}

// 	if clusterTrigger.Config.IsReverse {
// 		isTrigger = !isTrigger
// 	}

// 	return isTrigger, lst
// }

// procWins
func (clusterTrigger *ClusterTrigger) procWins(std *ClusterTriggerData, lst []*sgc7game.Result) (int, error) {
	if clusterTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0
		}

		return 0, nil
	}

	std.WinMulti = clusterTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti

		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin
	}

	return std.Wins, nil
}

// calcRespinNum
func (clusterTrigger *ClusterTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(clusterTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := clusterTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ClusterTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					slog.Int("SymbolNum", ret.SymbolNums),
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("ClusterTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(clusterTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := clusterTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("ClusterTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if clusterTrigger.Config.RespinNumWeightVW != nil {
		cr, err := clusterTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("ClusterTrigger.calcRespinNum:RespinNumWeightVW",
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if clusterTrigger.Config.RespinNum > 0 {
		return clusterTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// playgame
func (clusterTrigger *ClusterTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// clusterTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := cd.(*ClusterTriggerData)
	std.onNewStep()

	gs := clusterTrigger.GetTargetScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := clusterTrigger.CanTriggerWithScene(gameProp, gs, curpr, stake)

	if isTrigger {
		clusterTrigger.procWins(std, lst)

		if !clusterTrigger.Config.NeedDiscardResults {
			for _, v := range lst {
				clusterTrigger.AddResult(curpr, v, &std.BasicComponentData)

				std.SymbolNum += v.SymbolNums
				std.WildNum += v.Wilds
			}
		} else {
			for _, v := range lst {
				std.SymbolNum += v.SymbolNums
				std.WildNum += v.Wilds
			}
		}

		respinNum, err := clusterTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("ClusterTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", err
		}

		std.RespinNum = respinNum

		// err = clusterTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		// if err != nil {
		// 	goutils.Error("ClusterTrigger.OnPlayGame:procMask",
		// 		goutils.Err(err))

		// 	return err
		// }

		// if symbolTrigger.Config.TagSymbolNum != "" {
		// 	gameProp.TagInt(symbolTrigger.Config.TagSymbolNum, lst[0].SymbolNums)
		// }

		if len(clusterTrigger.Config.Awards) > 0 {
			gameProp.procAwards(plugin, clusterTrigger.Config.Awards, curpr, gp)
		}

		// if clusterTrigger.Config.SymbolAwardsWeights != nil {
		// 	for i := 0; i < lst[0].SymbolNums; i++ {
		// 		node, err := clusterTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
		// 		if err != nil {
		// 			goutils.Error("ClusterTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
		// 				goutils.Err(err))

		// 			return err
		// 		}

		// 		gameProp.procAwards(plugin, node.Awards, curpr, gp)
		// 	}
		// }

		if clusterTrigger.Config.JumpToComponent != "" {
			if gameProp.IsRespin(clusterTrigger.Config.JumpToComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if clusterTrigger.Config.ForceToNext {
						std.NextComponent = clusterTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(clusterTrigger.Config.JumpToComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, clusterTrigger.Config.JumpToComponent, !clusterTrigger.Config.IsAddRespinMode)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, clusterTrigger.Config.JumpToComponent, !clusterTrigger.Config.IsAddRespinMode)

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

			std.NextComponent = clusterTrigger.Config.JumpToComponent

			nc := clusterTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := clusterTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := clusterTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (clusterTrigger *ClusterTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	std := cd.(*ClusterTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", clusterTrigger.Name, std.NextComponent)
	}

	return nil
}

// // OnStatsWithPB -
// func (clusterTrigger *ClusterTrigger) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.ClusterTriggerData)
// 	if !isok {
// 		goutils.Error("ClusterTrigger.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return clusterTrigger.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (clusterTrigger *ClusterTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	wins := int64(0)
// 	isTrigger := false

// 	for _, v := range lst {
// 		gp, isok := v.CurGameModParams.(*GameParams)
// 		if isok {
// 			curComponent, isok := gp.MapComponentMsgs[clusterTrigger.Name]
// 			if isok {
// 				curwins, err := clusterTrigger.OnStatsWithPB(feature, curComponent, v)
// 				if err != nil {
// 					goutils.Error("ClusterTrigger.OnStats",
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
func (clusterTrigger *ClusterTrigger) NewComponentData() IComponentData {
	return &ClusterTriggerData{}
}

func (clusterTrigger *ClusterTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(WTCVWinMulti)
	if isok {
		return winMulti
	}

	return clusterTrigger.Config.WinMulti
}

// NewStats2 -
func (clusterTrigger *ClusterTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (clusterTrigger *ClusterTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache) {
	clusterTrigger.BasicComponent.OnStats2(icd, s2)

	cd := icd.(*ClusterTriggerData)

	s2.ProcStatsWins(clusterTrigger.Name, int64(cd.Wins), true)
}

// GetAllLinkComponents - get all link components
func (clusterTrigger *ClusterTrigger) GetAllLinkComponents() []string {
	return []string{clusterTrigger.Config.DefaultNextComponent, clusterTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (clusterTrigger *ClusterTrigger) GetNextLinkComponents() []string {
	return []string{clusterTrigger.Config.DefaultNextComponent, clusterTrigger.Config.JumpToComponent}
}

func (clusterTrigger *ClusterTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return clusterTrigger.Config.SymbolCodes
}

// CanTriggerWithScene -
func (clusterTrigger *ClusterTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) (bool, []*sgc7game.Result) {
	isTrigger := false
	lst := []*sgc7game.Result{}

	if clusterTrigger.Config.TriggerType == STTypeCluster {

		symbols := clusterTrigger.getSymbols(gameProp)

		currets, err := sgc7game.CalcClusterResult(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, clusterTrigger.Config.BetType),
			func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(clusterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(clusterTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int) int {
				return cursymbol
			})
		if err != nil {
			goutils.Error("ClusterTrigger.CanTriggerWithScene:CalcClusterResult",
				goutils.Err(err))

			return false, nil
		}

		// for _, v := range currets {
		// 	gameProp.ProcMulti(v)
		// }

		lst = append(lst, currets...)

		if len(lst) > 0 {
			isTrigger = true
		}
	}

	if clusterTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

func NewClusterTrigger(name string) IComponent {
	return &ClusterTrigger{
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
type jsonClusterTrigger struct {
	Symbols     []string `json:"symbols"`
	TriggerType string   `json:"triggerType"`
	BetType     string   `json:"betType"`
	MinNum      int      `json:"minNum"`
	WildSymbols []string `json:"wildSymbols"`
	WinMulti    int      `json:"winMulti"`
}

func (jwt *jsonClusterTrigger) build() *ClusterTriggerConfig {
	cfg := &ClusterTriggerConfig{
		Symbols:       jwt.Symbols,
		Type:          jwt.TriggerType,
		BetTypeString: jwt.BetType,
		MinNum:        jwt.MinNum,
		WildSymbols:   jwt.WildSymbols,
		WinMulti:      jwt.WinMulti,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseClusterTrigger(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseClusterTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseClusterTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonClusterTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseClusterTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(gamecfg, ctrls)
		if err != nil {
			goutils.Error("parseClusterTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ClusterTriggerTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
