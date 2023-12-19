package lowcode

import (
	"fmt"
	"os"

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

const SymbolTriggerTypeName = "symbolTrigger"

type SymbolTriggerType int

const (
	STTypeUnknow             SymbolTriggerType = 0 // 非法
	STTypeLines              SymbolTriggerType = 1 // 线中奖判断，一定是判断全部线，且读paytable来判断是否可以中奖
	STTypeWays               SymbolTriggerType = 2 // ways中奖判断，且读paytable来判断是否可以中奖
	STTypeScatters           SymbolTriggerType = 3 // scatter中奖判断，且读paytable来判断是否可以中奖
	STTypeCountScatter       SymbolTriggerType = 4 // scatter判断，需要传入minnum，不读paytable
	STTypeCountScatterInArea SymbolTriggerType = 5 // 区域内的scatter判断，需要传入minnum，不读paytable
	STTypeCheckLines         SymbolTriggerType = 6 // 线判断，一定是判断全部线，需要传入minnum，不读paytable
	STTypeCheckWays          SymbolTriggerType = 7 // ways判断，需要传入minnum，不读paytable
)

func ParseSymbolTriggerType(str string) SymbolTriggerType {
	if str == "lines" {
		return STTypeLines
	} else if str == "ways" {
		return STTypeWays
	} else if str == "scatters" {
		return STTypeScatters
	} else if str == "countscatter" {
		return STTypeCountScatter
	} else if str == "countscatterInArea" {
		return STTypeCountScatterInArea
	} else if str == "checkLines" {
		return STTypeCheckLines
	} else if str == "checkWays" {
		return STTypeCheckWays
	}

	return STTypeUnknow
}

type BetType int

const (
	BTypeNoPay    BetType = 0
	BTypeBet      BetType = 1
	BTypeTotalBet BetType = 2
)

func ParseBetType(str string) BetType {
	if str == "bet" {
		return BTypeBet
	} else if str == "totalBet" {
		return BTypeTotalBet
	}

	return BTypeNoPay
}

const (
	STDVSymbolNum string = "symbolNum" // 触发后，中奖的符号数量
	STDVWildNum   string = "wildNum"   // 触发后，中奖符号里的wild数量
	STDVRespinNum string = "respinNum" // 触发后，如果有产生respin的逻辑，这就是最终respin的次数
)

type SymbolTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
}

// OnNewGame -
func (symbolTriggerData *SymbolTriggerData) OnNewGame() {
	symbolTriggerData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (symbolTriggerData *SymbolTriggerData) OnNewStep() {
	symbolTriggerData.BasicComponentData.OnNewStep()

	symbolTriggerData.NextComponent = ""
	symbolTriggerData.SymbolNum = 0
	symbolTriggerData.WildNum = 0
	symbolTriggerData.RespinNum = 0
}

// BuildPBComponentData
func (symbolTriggerData *SymbolTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolTriggerData{
		BasicComponentData: symbolTriggerData.BuildPBBasicComponentData(),
		NextComponent:      symbolTriggerData.NextComponent,
		SymbolNum:          int32(symbolTriggerData.SymbolNum),
		WildNum:            int32(symbolTriggerData.WildNum),
		RespinNum:          int32(symbolTriggerData.RespinNum),
	}

	return pbcd
}

// GetVal -
func (symbolTriggerData *SymbolTriggerData) GetVal(key string) int {
	if key == STDVSymbolNum {
		return symbolTriggerData.SymbolNum
	} else if key == STDVWildNum {
		return symbolTriggerData.WildNum
	} else if key == STDVRespinNum {
		return symbolTriggerData.RespinNum
	}

	return 0
}

// SetVal -
func (symbolTriggerData *SymbolTriggerData) SetVal(key string, val int) {
	if key == STDVSymbolNum {
		symbolTriggerData.SymbolNum = val
	} else if key == STDVWildNum {
		symbolTriggerData.WildNum = val
	} else if key == STDVRespinNum {
		symbolTriggerData.RespinNum = val
	}
}

// SymbolTriggerConfig - configuration for SymbolTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type SymbolTriggerConfig struct {
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
	PosArea                         []int                         `yaml:"posArea" json:"posArea"`                                             // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
	CountScatterPayAs               string                        `yaml:"countScatterPayAs" json:"countScatterPayAs"`                         // countscatter时，按什么符号赔付
	SymbolCodeCountScatterPayAs     int                           `yaml:"-" json:"-"`                                                         // countscatter时，按什么符号赔付
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`                             // jump to
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`                                     // 如果触发，默认跳转jump to，这里可以强制走next分支
	TagSymbolNum                    string                        `yaml:"tagSymbolNum" json:"tagSymbolNum"`                                   // 这里可以将symbol数量记下来，别的地方能获取到
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                                               // 新的奖励系统
	SymbolAwardsWeights             *AwardsWeights                `yaml:"symbolAwardsWeights" json:"symbolAwardsWeights"`                     // 每个中奖符号随机一组奖励
	TargetMask                      string                        `yaml:"targetMask" json:"targetMask"`                                       // 如果是scatter这一组判断，可以把结果传递给一个mask
	IsReverse                       bool                          `yaml:"isReverse" json:"isReverse"`                                         // 如果isReverse，表示判定为否才触发
	NeedDiscardResults              bool                          `yaml:"needDiscardResults" json:"needDiscardResults"`                       // 如果needDiscardResults，表示抛弃results
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
}

type SymbolTrigger struct {
	*BasicComponent `json:"-"`
	Config          *SymbolTriggerConfig `json:"config"`
}

// Init -
func (symbolTrigger *SymbolTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolTrigger.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolTrigger.InitEx(cfg, pool)
}

// InitEx -
func (symbolTrigger *SymbolTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolTrigger.Config = cfg.(*SymbolTriggerConfig)
	symbolTrigger.Config.ComponentType = SymbolTriggerTypeName

	for _, s := range symbolTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolTrigger.InitEx:Symbol",
				zap.String("symbol", s),
				zap.Error(ErrIvalidSymbol))
		}

		symbolTrigger.Config.SymbolCodes = append(symbolTrigger.Config.SymbolCodes, sc)
	}

	sc, isok := pool.DefaultPaytables.MapSymbols[symbolTrigger.Config.CountScatterPayAs]
	if isok {
		symbolTrigger.Config.SymbolCodeCountScatterPayAs = sc
	}

	for _, s := range symbolTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolTrigger.InitEx:WildSymbols",
				zap.String("symbol", s),
				zap.Error(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		symbolTrigger.Config.WildSymbolCodes = append(symbolTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(symbolTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("SymbolTrigger.InitEx:WildSymbols",
			zap.String("SymbolTriggerType", symbolTrigger.Config.Type),
			zap.Error(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	symbolTrigger.Config.TriggerType = stt

	symbolTrigger.Config.BetType = ParseBetType(symbolTrigger.Config.BetTypeString)

	for _, award := range symbolTrigger.Config.Awards {
		award.Init()
	}

	if symbolTrigger.Config.SymbolAwardsWeights != nil {
		symbolTrigger.Config.SymbolAwardsWeights.Init()
	}

	if symbolTrigger.Config.TriggerType == STTypeLines || symbolTrigger.Config.TriggerType == STTypeWays {
		symbolTrigger.Config.ExcludeSymbolCodes = GetExcludeSymbols(pool.DefaultPaytables, symbolTrigger.Config.SymbolCodes)
	}

	symbolTrigger.Config.CheckWinType = ParseCheckWinType(symbolTrigger.Config.StrCheckWinType)

	if symbolTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(symbolTrigger.Config.RespinNumWeight, symbolTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("SymbolTrigger.InitEx:LoadIntWeights",
				zap.String("Weight", symbolTrigger.Config.RespinNumWeight),
				zap.Error(err))

			return err
		}

		symbolTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(symbolTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range symbolTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, symbolTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("SymbolTrigger.InitEx:LoadIntWeights",
					zap.String("Weight", v),
					zap.Error(err))

				return err
			}

			symbolTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	symbolTrigger.onInit(&symbolTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolTrigger *SymbolTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if symbolTrigger.Config.TargetMask != "" {
		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, symbolTrigger.Config.TargetMask, mask)
	}

	return nil
}

// CanTrigger -
func (symbolTrigger *SymbolTrigger) triggerScatter(gameProp *GameProperty, stake *sgc7game.Stake, gs *sgc7game.GameScene) *sgc7game.Result {
	return sgc7game.CalcScatter4(gs, gameProp.CurPaytables, symbolTrigger.Config.SymbolCodes[0], gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
		func(scatter int, cursymbol int) bool {
			return goutils.IndexOfIntSlice(symbolTrigger.Config.SymbolCodes, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		}, true)
}

// CanTrigger -
func (symbolTrigger *SymbolTrigger) CanTrigger(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, stake *sgc7game.Stake, isSaveResult bool) (bool, []*sgc7game.Result) {
	std := gameProp.MapComponentData[symbolTrigger.Name].(*SymbolTriggerData)

	gs := symbolTrigger.GetTargetScene2(gameProp, curpr, &std.BasicComponentData, symbolTrigger.Name, "")

	isTrigger := false
	lst := []*sgc7game.Result{}

	if symbolTrigger.Config.TriggerType == STTypeLines {
		os := symbolTrigger.GetTargetOtherScene2(gameProp, curpr, &std.BasicComponentData, symbolTrigger.Name, "")

		if os != nil {
			for i, v := range gameProp.CurLineData.Lines {
				isTriggerFull := false
				if symbolTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
					ret := sgc7game.CalcLine2(gs, gameProp.CurPaytables, v, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int) int {
							return cursymbol
						}, func(x, y int) int {
							return os.Arr[x][y]
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						if ret.SymbolNums == gs.Width {
							isTriggerFull = true
						}

						if isSaveResult {
							symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						}
					}
				}

				if !isTriggerFull && symbolTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
					ret := sgc7game.CalcLineRL2(gs, gameProp.CurPaytables, v, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int) int {
							return cursymbol
						}, func(x, y int) int {
							return os.Arr[x][y]
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						if isSaveResult {
							symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						}
					}
				}
			}
		} else {
			for i, v := range gameProp.CurLineData.Lines {
				if symbolTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
					ret := sgc7game.CalcLineEx(gs, gameProp.CurPaytables, v, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
							return 1
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						if isSaveResult {
							symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						}
					}
				}

				if symbolTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
					ret := sgc7game.CalcLineRLEx(gs, gameProp.CurPaytables, v, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
							return 1
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						if isSaveResult {
							symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						}
					}
				}
			}
		}
	} else if symbolTrigger.Config.TriggerType == STTypeCheckLines {

		for i, v := range gameProp.CurLineData.Lines {
			if symbolTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
				ret := sgc7game.CheckLine(gs, v, symbolTrigger.Config.MinNum,
					func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
					}, func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
					}, func(cursymbol int, startsymbol int) bool {
						if cursymbol == startsymbol {
							return true
						}

						return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
					}, func(cursymbol int) int {
						return cursymbol
					})
				if ret != nil {
					ret.LineIndex = i

					gameProp.ProcMulti(ret)

					lst = append(lst, ret)

					if isSaveResult {
						symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
					}
				}
			}

			if symbolTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
				ret := sgc7game.CheckLineRL(gs, v, symbolTrigger.Config.MinNum,
					func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
					}, func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
					}, func(cursymbol int, startsymbol int) bool {
						if cursymbol == startsymbol {
							return true
						}

						return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
					}, func(cursymbol int) int {
						return cursymbol
					})
				if ret != nil {
					ret.LineIndex = i

					gameProp.ProcMulti(ret)

					lst = append(lst, ret)

					if isSaveResult {
						symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
					}
				}
			}
		}

	} else if symbolTrigger.Config.TriggerType == STTypeWays {
		os := symbolTrigger.GetTargetOtherScene2(gameProp, curpr, &std.BasicComponentData, symbolTrigger.Name, "")

		if os != nil {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(x, y int) int {
					return os.Arr[x][y]
				})

			for _, v := range currets {
				gameProp.ProcMulti(v)

				if isSaveResult {
					symbolTrigger.AddResult(curpr, v, &std.BasicComponentData)
				}
			}

			lst = append(lst, currets...)
		} else {
			currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
				}, func(x, y int) int {
					return 1
				})

			for _, v := range currets {
				gameProp.ProcMulti(v)

				if isSaveResult {
					symbolTrigger.AddResult(curpr, v, &std.BasicComponentData)
				}
			}

			lst = append(lst, currets...)
		}

	} else if symbolTrigger.Config.TriggerType == STTypeCheckWays {
		currets := sgc7game.CheckWays(gs, symbolTrigger.Config.MinNum,
			func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
				return goutils.IndexOfIntSlice(symbolTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		for _, v := range currets {
			gameProp.ProcMulti(v)

			if isSaveResult {
				symbolTrigger.AddResult(curpr, v, &std.BasicComponentData)
			}
		}

		lst = append(lst, currets...)

	} else if symbolTrigger.Config.TriggerType == STTypeScatters {
		ret := symbolTrigger.triggerScatter(gameProp, stake, gs)
		// for _, s := range symbolTrigger.Config.SymbolCodes {
		// ret := sgc7game.CalcScatter4(gs, gameProp.CurPaytables, symbolTrigger.Config.SymbolCodes[0], gameProp.GetBet2(stake, symbolTrigger.Config.BetType),
		// 	func(scatter int, cursymbol int) bool {
		// 		return goutils.IndexOfIntSlice(symbolTrigger.Config.SymbolCodes, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		// 	}, true)

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				gameProp.ProcMulti(ret)
			}

			if isSaveResult {
				symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			}

			isTrigger = true

			lst = append(lst, ret)
		}
		// }
	} else if symbolTrigger.Config.TriggerType == STTypeCountScatter {
		// for _, s := range symbolTrigger.Config.SymbolCodes {
		ret := sgc7game.CalcScatterEx(gs, symbolTrigger.Config.SymbolCodes[0], symbolTrigger.Config.MinNum, func(scatter int, cursymbol int) bool {
			return goutils.IndexOfIntSlice(symbolTrigger.Config.SymbolCodes, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
		})

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if symbolTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, symbolTrigger.Config.BetType)
				}

				gameProp.ProcMulti(ret)
			}

			if isSaveResult {
				symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			}

			isTrigger = true

			lst = append(lst, ret)
		}
		// }
	} else if symbolTrigger.Config.TriggerType == STTypeCountScatterInArea {
		// for _, s := range symbolTrigger.Config.SymbolCodes {
		ret := sgc7game.CountScatterInArea(gs, symbolTrigger.Config.SymbolCodes[0], symbolTrigger.Config.MinNum,
			func(x, y int) bool {
				return x >= symbolTrigger.Config.PosArea[0] && x <= symbolTrigger.Config.PosArea[1] && y >= symbolTrigger.Config.PosArea[2] && y <= symbolTrigger.Config.PosArea[3]
			},
			func(scatter int, cursymbol int) bool {
				return goutils.IndexOfIntSlice(symbolTrigger.Config.SymbolCodes, cursymbol, 0) >= 0 || goutils.IndexOfIntSlice(symbolTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			})

		if ret != nil {
			if symbolTrigger.Config.BetType == BTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if symbolTrigger.Config.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[symbolTrigger.Config.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet2(stake, symbolTrigger.Config.BetType)
				}

				gameProp.ProcMulti(ret)
			}

			if isSaveResult {
				symbolTrigger.AddResult(curpr, ret, &std.BasicComponentData)
			}

			isTrigger = true

			lst = append(lst, ret)
		}
		// }
	}

	if symbolTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

// playgame
func (symbolTrigger *SymbolTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(symbolTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := symbolTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("SymbolTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					zap.Int("SymbolNum", ret.SymbolNums),
					zap.Error(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("SymbolTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				zap.Int("SymbolNum", ret.SymbolNums),
				zap.Error(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(symbolTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := symbolTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("SymbolTrigger.calcRespinNum:RespinNumWithScatterNum",
				zap.Int("SymbolNum", ret.SymbolNums),
				zap.Error(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if symbolTrigger.Config.RespinNumWeightVW != nil {
		cr, err := symbolTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("SymbolTrigger.calcRespinNum:RespinNumWeightVW",
				zap.Error(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if symbolTrigger.Config.RespinNum > 0 {
		return symbolTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// playgame
func (symbolTrigger *SymbolTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := gameProp.MapComponentData[symbolTrigger.Name].(*SymbolTriggerData)

	gs := symbolTrigger.GetTargetScene2(gameProp, curpr, &std.BasicComponentData, symbolTrigger.Name, "")

	isTrigger, lst := symbolTrigger.CanTrigger(gameProp, curpr, gp, stake, !symbolTrigger.Config.NeedDiscardResults)

	if isTrigger {
		std.SymbolNum = lst[0].SymbolNums
		std.WildNum = lst[0].Wilds

		respinNum, err := symbolTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("SymbolTrigger.OnPlayGame:calcRespinNum",
				zap.Error(err))

			return nil
		}

		std.RespinNum = respinNum

		err = symbolTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("SymbolTrigger.OnPlayGame:procMask",
				zap.Error(err))

			return err
		}

		if symbolTrigger.Config.TagSymbolNum != "" {
			gameProp.TagInt(symbolTrigger.Config.TagSymbolNum, lst[0].SymbolNums)
		}

		if len(symbolTrigger.Config.Awards) > 0 {
			gameProp.procAwards(plugin, symbolTrigger.Config.Awards, curpr, gp)
		}

		if symbolTrigger.Config.SymbolAwardsWeights != nil {
			for i := 0; i < lst[0].SymbolNums; i++ {
				node, err := symbolTrigger.Config.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("SymbolTrigger.OnPlayGame:SymbolAwardsWeights.RandVal",
						zap.Error(err))

					return err
				}

				gameProp.procAwards(plugin, node.Awards, curpr, gp)
			}
		}

		if symbolTrigger.Config.JumpToComponent != "" {
			if gameProp.IsRespin(symbolTrigger.Config.JumpToComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if symbolTrigger.Config.ForceToNext {
						std.NextComponent = symbolTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(symbolTrigger.Config.JumpToComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, symbolTrigger.Config.JumpToComponent, true)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, symbolTrigger.Config.JumpToComponent, true)

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

			std.NextComponent = symbolTrigger.Config.JumpToComponent

			symbolTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nil
		}
	}

	symbolTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolTrigger *SymbolTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	std := gameProp.MapComponentData[symbolTrigger.Name].(*SymbolTriggerData)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v", symbolTrigger.Name, std.NextComponent)
	}

	return nil
}

// OnStats
func (symbolTrigger *SymbolTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (symbolTrigger *SymbolTrigger) NewComponentData() IComponentData {
	return &SymbolTriggerData{}
}

func NewSymbolTrigger(name string) IComponent {
	return &SymbolTrigger{
		BasicComponent: NewBasicComponent(name),
	}
}
