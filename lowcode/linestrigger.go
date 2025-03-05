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

const LinesTriggerTypeName = "linesTrigger"

type LinesTriggerData struct {
	BasicComponentData
	NextComponent string
	SymbolNum     int
	WildNum       int
	RespinNum     int
	Wins          int
	WinMulti      int
}

// OnNewGame -
func (linesTriggerData *LinesTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	linesTriggerData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (linesTriggerData *LinesTriggerData) onNewStep() {
	linesTriggerData.CashWin = 0
	linesTriggerData.CoinWin = 0

	linesTriggerData.UsedResults = nil
	linesTriggerData.NextComponent = ""
	linesTriggerData.SymbolNum = 0
	linesTriggerData.WildNum = 0
	linesTriggerData.RespinNum = 0
	linesTriggerData.Wins = 0
	linesTriggerData.WinMulti = 1
}

// Clone
func (linesTriggerData *LinesTriggerData) Clone() IComponentData {
	target := &LinesTriggerData{
		BasicComponentData: linesTriggerData.CloneBasicComponentData(),
		NextComponent:      linesTriggerData.NextComponent,
		SymbolNum:          linesTriggerData.SymbolNum,
		WildNum:            linesTriggerData.WildNum,
		RespinNum:          linesTriggerData.RespinNum,
		Wins:               linesTriggerData.Wins,
		WinMulti:           linesTriggerData.WinMulti,
	}

	return target
}

// BuildPBComponentData
func (linesTriggerData *LinesTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.LinesTriggerData{
		BasicComponentData: linesTriggerData.BuildPBBasicComponentData(),
		NextComponent:      linesTriggerData.NextComponent,
		SymbolNum:          int32(linesTriggerData.SymbolNum),
		WildNum:            int32(linesTriggerData.WildNum),
		RespinNum:          int32(linesTriggerData.RespinNum),
		Wins:               int32(linesTriggerData.Wins),
		WinMulti:           int32(linesTriggerData.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (linesTriggerData *LinesTriggerData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolNum {
		return linesTriggerData.SymbolNum, true
	} else if key == CVWildNum {
		return linesTriggerData.WildNum, true
	} else if key == CVRespinNum {
		return linesTriggerData.RespinNum, true
	} else if key == CVWins {
		return linesTriggerData.Wins, true
	} else if key == CVResultNum || key == CVWinResultNum {
		return len(linesTriggerData.UsedResults), true
	}

	return 0, false
}

// LinesTriggerConfig - configuration for LinesTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type LinesTriggerConfig struct {
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
	LineData                        string                        `yaml:"linedata" json:"linedata"`                                           // linedata
	JumpToComponent                 string                        `yaml:"jumpToComponent" json:"jumpToComponent"`                             // jump to
	ForceToNext                     bool                          `yaml:"forceToNext" json:"forceToNext"`                                     // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards                          []*Award                      `yaml:"awards" json:"awards"`                                               // 新的奖励系统
	TargetMask                      string                        `yaml:"targetMask" json:"targetMask"`                                       // 如果是scatter这一组判断，可以把结果传递给一个mask
	IsReverse                       bool                          `yaml:"isReverse" json:"isReverse"`                                         // 如果isReverse，表示判定为否才触发
	PiggyBankComponent              string                        `yaml:"piggyBankComponent" json:"piggyBankComponent"`                       // piggyBank component
	OutputToComponent               string                        `yaml:"outputToComponent" json:"outputToComponent"`                         // 将结果给到一个 positionCollection
	IsAddRespinMode                 bool                          `yaml:"isAddRespinMode" json:"isAddRespinMode"`                             // 是否是增加respinNum模式，默认是增加triggerNum模式
	RespinNum                       int                           `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight                 string                        `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWeightVW               *sgc7game.ValWeights2         `yaml:"-" json:"-"`                                                         // respin number weight
	RespinNumWithScatterNum         map[int]int                   `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum   map[int]string                `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	RespinNumWeightWithScatterNumVW map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`                                                         // respin number weight with scatter number
}

// SetLinkComponent
func (cfg *LinesTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type LinesTrigger struct {
	*BasicComponent `json:"-"`
	Config          *LinesTriggerConfig `json:"config"`
}

// Init -
func (linesTrigger *LinesTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("LinesTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &LinesTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("LinesTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return linesTrigger.InitEx(cfg, pool)
}

// InitEx -
func (linesTrigger *LinesTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	linesTrigger.Config = cfg.(*LinesTriggerConfig)
	linesTrigger.Config.ComponentType = LinesTriggerTypeName

	linesTrigger.Config.OSMulType = ParseOtherSceneMultiType(linesTrigger.Config.OSMulTypeString)

	for _, s := range linesTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("LinesTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		linesTrigger.Config.SymbolCodes = append(linesTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range linesTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("LinesTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))

			return ErrIvalidSymbol
		}

		linesTrigger.Config.WildSymbolCodes = append(linesTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(linesTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("LinesTrigger.InitEx:WildSymbols",
			slog.String("SymbolTriggerType", linesTrigger.Config.Type),
			goutils.Err(ErrIvalidSymbolTriggerType))

		return ErrIvalidSymbolTriggerType
	}

	linesTrigger.Config.TriggerType = stt

	linesTrigger.Config.BetType = ParseBetType(linesTrigger.Config.BetTypeString)

	for _, award := range linesTrigger.Config.Awards {
		award.Init()
	}

	linesTrigger.Config.CheckWinType = ParseCheckWinType(linesTrigger.Config.StrCheckWinType)

	if linesTrigger.Config.RespinNumWeight != "" {
		vw2, err := pool.LoadIntWeights(linesTrigger.Config.RespinNumWeight, linesTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("LinesTrigger.InitEx:LoadIntWeights",
				slog.String("Weight", linesTrigger.Config.RespinNumWeight),
				goutils.Err(err))

			return err
		}

		linesTrigger.Config.RespinNumWeightVW = vw2
	}

	if len(linesTrigger.Config.RespinNumWeightWithScatterNum) > 0 {
		for k, v := range linesTrigger.Config.RespinNumWeightWithScatterNum {
			vw2, err := pool.LoadIntWeights(v, linesTrigger.Config.UseFileMapping)
			if err != nil {
				goutils.Error("LinesTrigger.InitEx:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))

				return err
			}

			linesTrigger.Config.RespinNumWeightWithScatterNumVW[k] = vw2
		}
	}

	if linesTrigger.Config.WinMulti <= 0 {
		linesTrigger.Config.WinMulti = 1
	}

	// if linesTrigger.Config.BetType == BTypeNoPay {
	// 	linesTrigger.Config.NeedDiscardResults = true
	// }

	linesTrigger.onInit(&linesTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (linesTrigger *LinesTrigger) procMask(gs *sgc7game.GameScene, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, ret *sgc7game.Result) error {

	if linesTrigger.Config.TargetMask != "" {
		mask := make([]bool, gs.Width)

		for i := 0; i < len(ret.Pos)/2; i++ {
			mask[ret.Pos[i*2]] = true
		}

		gameProp.UseComponent(linesTrigger.Config.TargetMask)

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, linesTrigger.Config.TargetMask, mask, false)
	}

	return nil
}

func (linesTrigger *LinesTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return linesTrigger.Config.SymbolCodes
}

func (linesTrigger *LinesTrigger) getLineData(gameProp *GameProperty, cd *LinesTriggerData) (*sgc7game.LineData, error) {
	ldname := cd.GetConfigVal(CCVLineData)
	if ldname != "" {
		v, isok := gameProp.Pool.Config.MapLinedate[ldname]
		if !isok {
			goutils.Error("LinesTrigger.getLineData:getLineData",
				slog.String("val", ldname),
				goutils.Err(ErrInvalidLineData))

			return nil, ErrInvalidLineData
		}

		return v, nil
	}

	if linesTrigger.Config.LineData == "" {
		return gameProp.CurLineData, nil
	}

	v, isok := gameProp.Pool.Config.MapLinedate[linesTrigger.Config.LineData]
	if !isok {
		goutils.Error("LinesTrigger.getLineData:getLineData",
			slog.String("val", linesTrigger.Config.LineData),
			goutils.Err(ErrInvalidLineData))

		return nil, ErrInvalidLineData
	}

	return v, nil
}

// canTrigger -
func (linesTrigger *LinesTrigger) canTrigger(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, _ *sgc7game.PlayResult, stake *sgc7game.Stake, cd *LinesTriggerData) (bool, []*sgc7game.Result) {
	// std := cd.(*LinesTriggerData)

	isTrigger := false
	lst := []*sgc7game.Result{}
	lstSym := linesTrigger.getSymbols(gameProp)

	if linesTrigger.Config.OSMulType != OSMTNone && os == nil {
		goutils.Error("LinesTrigger.canTrigger",
			goutils.Err(ErrInvalidOtherScene))

		return false, nil
	}

	funcCalcMulti := GetSymbolValMultiFunc(linesTrigger.Config.OSMulType)

	ld, err := linesTrigger.getLineData(gameProp, cd)
	if err != nil {
		goutils.Error("LinesTrigger.canTrigger",
			goutils.Err(err))

		return false, nil
	}

	if linesTrigger.Config.TriggerType == STTypeLines {
		if linesTrigger.Config.OSMulType == OSMTNone { // no otherscene multi
			if linesTrigger.Config.CheckWinType == CheckWinTypeCount {
				for _, cs := range linesTrigger.getSymbols(gameProp) {
					for i, v := range ld.Lines {
						ret := sgc7game.CountSymbolOnLine(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType), cs,
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return 1
							}, funcCalcMulti)
						if ret != nil {
							ret.LineIndex = i

							lst = append(lst, ret)
						}
					}
				}
			} else {
				for i, v := range ld.Lines {
					if linesTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
						ret := sgc7game.CalcLineEx(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
								return 1
							}, func(cursymbol int) int {
								return cursymbol
							})
						if ret != nil {
							ret.LineIndex = i

							lst = append(lst, ret)
						}
					}

					if linesTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
						ret := sgc7game.CalcLineRLEx(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
								return 1
							}, func(cursymbol int) int {
								return cursymbol
							})
						if ret != nil {
							ret.LineIndex = i

							lst = append(lst, ret)
						}
					}
				}
			}
		} else { // otherscene multi
			if linesTrigger.Config.CheckWinType == CheckWinTypeCount {
				for _, cs := range linesTrigger.getSymbols(gameProp) {
					for i, v := range ld.Lines {
						ret := sgc7game.CountSymbolOnLine(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType), cs,
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return os.Arr[x][y]
							}, funcCalcMulti)
						if ret != nil {
							ret.LineIndex = i

							// gameProp.ProcMulti(ret)

							lst = append(lst, ret)
						}
					}
				}
			} else {
				for i, v := range ld.Lines {
					isTriggerFull := false
					if linesTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
						ret := sgc7game.CalcLine2(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
								// return goutils.IndexOfIntSlice(linesTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return os.Arr[x][y]
							})
						if ret != nil {
							ret.LineIndex = i

							// gameProp.ProcMulti(ret)

							lst = append(lst, ret)

							if ret.SymbolNums == gs.Width {
								isTriggerFull = true
							}

							// if isSaveResult {
							// 	linesTrigger.AddResult(curpr, ret, &std.BasicComponentData)
							// }
						}
					}

					if !isTriggerFull && linesTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
						ret := sgc7game.CalcLineRL2(gs, gameProp.CurPaytables, v, gameProp.GetBet3(stake, linesTrigger.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
								// return goutils.IndexOfIntSlice(linesTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return os.Arr[x][y]
							})
						if ret != nil {
							ret.LineIndex = i

							// gameProp.ProcMulti(ret)

							lst = append(lst, ret)

							// if isSaveResult {
							// 	linesTrigger.AddResult(curpr, ret, &std.BasicComponentData)
							// }
						}
					}
				}
			}
		}

		if len(lst) > 0 {
			isTrigger = true
		}
	} else if linesTrigger.Config.TriggerType == STTypeCheckLines {
		if linesTrigger.Config.CheckWinType == CheckWinTypeCount {
			// for _, cs := range linesTrigger.Config.SymbolCodes {
			for _, cs := range linesTrigger.getSymbols(gameProp) {
				for i, v := range ld.Lines {
					ret := sgc7game.CountSymbolOnLine(gs, gameProp.CurPaytables, v, 0, cs,
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int) int {
							return cursymbol
						}, func(x, y int) int {
							return 1
						}, funcCalcMulti)
					if ret != nil {
						ret.LineIndex = i

						// gameProp.ProcMulti(ret)

						lst = append(lst, ret)
					}
				}
			}
		} else {
			for i, v := range ld.Lines {
				if linesTrigger.Config.CheckWinType != CheckWinTypeRightLeft {
					ret := sgc7game.CheckLine(gs, v, linesTrigger.Config.MinNum,
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
							// return goutils.IndexOfIntSlice(linesTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						// gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						// if isSaveResult {
						// 	linesTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						// }
					}
				}

				if linesTrigger.Config.CheckWinType != CheckWinTypeLeftRight {
					ret := sgc7game.CheckLineRL(gs, v, linesTrigger.Config.MinNum,
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(lstSym, cursymbol, 0) >= 0
							// return goutils.IndexOfIntSlice(linesTrigger.Config.ExcludeSymbolCodes, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(linesTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						// gameProp.ProcMulti(ret)

						lst = append(lst, ret)

						// if isSaveResult {
						// 	linesTrigger.AddResult(curpr, ret, &std.BasicComponentData)
						// }
					}
				}
			}
		}

		if len(lst) > 0 {
			isTrigger = true
		}
	}

	if linesTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

// procWins
func (linesTrigger *LinesTrigger) procWins(gameProp *GameProperty, curpr *sgc7game.PlayResult, std *LinesTriggerData, lst []*sgc7game.Result) (int, error) {
	if linesTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0

			linesTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds
		}

		return 0, nil
	}

	std.WinMulti = linesTrigger.GetWinMulti(&std.BasicComponentData)

	for _, v := range lst {
		v.OtherMul = std.WinMulti
		v.CoinWin *= std.WinMulti
		v.CashWin *= std.WinMulti

		std.Wins += v.CoinWin

		linesTrigger.AddResult(curpr, v, &std.BasicComponentData)

		std.SymbolNum += v.SymbolNums
		std.WildNum += v.Wilds
	}

	if std.Wins > 0 {
		if linesTrigger.Config.PiggyBankComponent != "" {
			cd := gameProp.GetCurComponentDataWithName(linesTrigger.Config.PiggyBankComponent)
			if cd == nil {
				goutils.Error("LinesTrigger.procWins:GetCurComponentDataWithName",
					slog.String("PiggyBankComponent", linesTrigger.Config.PiggyBankComponent),
					goutils.Err(ErrInvalidComponent))

				return 0, ErrInvalidComponent
			}

			cd.ChgConfigIntVal(CCVSavedMoney, std.Wins)

			for _, v := range lst {
				v.IsNoPayNow = true
			}

			gameProp.UseComponent(linesTrigger.Config.PiggyBankComponent)
		}
	}

	return std.Wins, nil
}

// calcRespinNum
func (linesTrigger *LinesTrigger) calcRespinNum(plugin sgc7plugin.IPlugin, ret *sgc7game.Result) (int, error) {

	if len(linesTrigger.Config.RespinNumWeightWithScatterNumVW) > 0 {
		vw2, isok := linesTrigger.Config.RespinNumWeightWithScatterNumVW[ret.SymbolNums]
		if isok {
			cr, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("LinesTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
					slog.Int("SymbolNum", ret.SymbolNums),
					goutils.Err(err))

				return 0, err
			}

			return cr.Int(), nil
		} else {
			goutils.Error("LinesTrigger.calcRespinNum:RespinNumWeightWithScatterNumVW",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}
	} else if len(linesTrigger.Config.RespinNumWithScatterNum) > 0 {
		v, isok := linesTrigger.Config.RespinNumWithScatterNum[ret.SymbolNums]
		if !isok {
			goutils.Error("LinesTrigger.calcRespinNum:RespinNumWithScatterNum",
				slog.Int("SymbolNum", ret.SymbolNums),
				goutils.Err(ErrInvalidSymbolNum))

			return 0, ErrInvalidSymbolNum
		}

		return v, nil
	} else if linesTrigger.Config.RespinNumWeightVW != nil {
		cr, err := linesTrigger.Config.RespinNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("LinesTrigger.calcRespinNum:RespinNumWeightVW",
				goutils.Err(err))

			return 0, err
		}

		return cr.Int(), nil
	} else if linesTrigger.Config.RespinNum > 0 {
		return linesTrigger.Config.RespinNum, nil
	}

	return 0, nil
}

// OnProcControllers -
func (linesTrigger *LinesTrigger) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(linesTrigger.Config.Awards) > 0 {
		gameProp.procAwards(plugin, linesTrigger.Config.Awards, curpr, gp)
	}
}

// procPositionCollection
func (linesTrigger *LinesTrigger) procPositionCollection(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	cd *LinesTriggerData) error {

	if linesTrigger.Config.OutputToComponent != "" {
		pcd := gameProp.GetComponentDataWithName(linesTrigger.Config.OutputToComponent)
		if pcd != nil {
			gameProp.UseComponent(linesTrigger.Config.OutputToComponent)
			pc := gameProp.Components.MapComponents[linesTrigger.Config.OutputToComponent]

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

// playgame
func (linesTrigger *LinesTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	std := cd.(*LinesTriggerData)
	std.onNewStep()

	gs := linesTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	os := linesTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := linesTrigger.canTrigger(gameProp, gs, os, curpr, stake, std)

	if isTrigger {
		linesTrigger.procWins(gameProp, curpr, std, lst)

		respinNum, err := linesTrigger.calcRespinNum(plugin, lst[0])
		if err != nil {
			goutils.Error("LinesTrigger.OnPlayGame:calcRespinNum",
				goutils.Err(err))

			return "", nil
		}

		std.RespinNum = respinNum

		err = linesTrigger.procMask(gs, gameProp, curpr, gp, plugin, lst[0])
		if err != nil {
			goutils.Error("LinesTrigger.OnPlayGame:procMask",
				goutils.Err(err))

			return "", err
		}

		err = linesTrigger.procPositionCollection(gameProp, curpr, std)
		if err != nil {
			goutils.Error("LinesTrigger.OnPlayGame:procPositionCollection",
				goutils.Err(err))

			return "", err
		}

		linesTrigger.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

		if linesTrigger.Config.JumpToComponent != "" {
			if gameProp.IsRespin(linesTrigger.Config.JumpToComponent) {
				// 如果jumpto是一个respin，那么就需要trigger respin
				if std.RespinNum == 0 {
					if linesTrigger.Config.ForceToNext {
						std.NextComponent = linesTrigger.Config.DefaultNextComponent
					} else {
						rn := gameProp.GetLastRespinNum(linesTrigger.Config.JumpToComponent)
						if rn > 0 {
							gameProp.TriggerRespin(plugin, curpr, gp, 0, linesTrigger.Config.JumpToComponent, !linesTrigger.Config.IsAddRespinMode)

							lst[0].Type = sgc7game.RTFreeGame
							lst[0].Value = rn
						}
					}
				} else {
					// 如果jumpto是respin，需要treigger这个respin
					gameProp.TriggerRespin(plugin, curpr, gp, std.RespinNum, linesTrigger.Config.JumpToComponent, !linesTrigger.Config.IsAddRespinMode)

					lst[0].Type = sgc7game.RTFreeGame
					lst[0].Value = std.RespinNum
				}
			}

			std.NextComponent = linesTrigger.Config.JumpToComponent

			nc := linesTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := linesTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := linesTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (linesTrigger *LinesTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	std := cd.(*LinesTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", linesTrigger.Name, std.NextComponent)
	}

	return nil
}

// // OnStatsWithPB -
// func (linesTrigger *LinesTrigger) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.LinesTriggerData)
// 	if !isok {
// 		goutils.Error("LinesTrigger.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return linesTrigger.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (linesTrigger *LinesTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	wins := int64(0)
// 	isTrigger := false

// 	for _, v := range lst {
// 		gp, isok := v.CurGameModParams.(*GameParams)
// 		if isok {
// 			curComponent, isok := gp.MapComponentMsgs[linesTrigger.Name]
// 			if isok {
// 				curwins, err := linesTrigger.OnStatsWithPB(feature, curComponent, v)
// 				if err != nil {
// 					goutils.Error("LinesTrigger.OnStats",
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
func (linesTrigger *LinesTrigger) NewComponentData() IComponentData {
	return &LinesTriggerData{}
}

func (linesTrigger *LinesTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return linesTrigger.Config.WinMulti
}

// NewStats2 -
func (linesTrigger *LinesTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (linesTrigger *LinesTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	linesTrigger.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*LinesTriggerData)

	s2.ProcStatsWins(linesTrigger.Name, int64(cd.Wins))
}

// GetAllLinkComponents - get all link components
func (linesTrigger *LinesTrigger) GetAllLinkComponents() []string {
	return []string{linesTrigger.Config.DefaultNextComponent, linesTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (linesTrigger *LinesTrigger) GetNextLinkComponents() []string {
	return []string{linesTrigger.Config.DefaultNextComponent, linesTrigger.Config.JumpToComponent}
}

// CanTriggerWithScene -
func (linesTrigger *LinesTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) {
	cd := icd.(*LinesTriggerData)

	return linesTrigger.canTrigger(gameProp, gs, nil, curpr, stake, cd)
}

func NewLinesTrigger(name string) IComponent {
	return &LinesTrigger{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "triggerType": "lines",
// "betType": "bet",
// "checkWinType": "left2right",
// "winMulti": 1,
// "symbols": [
//
//	"J",
//	"H",
//	"G",
//	"F",
//	"E",
//	"D",
//	"C",
//	"B",
//	"A",
//	"W"
//
// ],
// "wildSymbols": [
//
//	"W"
//
// ]
type jsonLinesTrigger struct {
	Symbols             []string `json:"symbols"`
	TriggerType         string   `json:"triggerType"`
	CheckWinType        string   `json:"checkWinType"`
	BetType             string   `json:"betType"`
	SymbolValsMulti     string   `json:"symbolValsMulti"`
	MinNum              int      `json:"minNum"`
	WildSymbols         []string `json:"wildSymbols"`
	WinMulti            int      `json:"winMulti"`
	PutMoneyInPiggyBank string   `json:"putMoneyInPiggyBank"`
	OutputToComponent   string   `json:"outputToComponent"`
}

func (jcfg *jsonLinesTrigger) build() *LinesTriggerConfig {
	cfg := &LinesTriggerConfig{
		Symbols:            jcfg.Symbols,
		Type:               jcfg.TriggerType,
		BetTypeString:      jcfg.BetType,
		StrCheckWinType:    jcfg.CheckWinType,
		OSMulTypeString:    jcfg.SymbolValsMulti,
		MinNum:             jcfg.MinNum,
		WildSymbols:        jcfg.WildSymbols,
		WinMulti:           jcfg.WinMulti,
		PiggyBankComponent: jcfg.PutMoneyInPiggyBank,
		OutputToComponent:  jcfg.OutputToComponent,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseLinesTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseLinesTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseLinesTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonLinesTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseLinesTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseLinesTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: LinesTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
