package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const BasicWinsTypeName = "basicWins"

const (
	WinTypeLines              = "lines"
	WinTypeWays               = "ways"
	WinTypeScatters           = "scatters"
	WinTypeCountScatter       = "countscatter"
	WinTypeCountScatterInArea = "countscatterInArea"

	BetTypeNormal   = "bet"
	BetTypeTotalBet = "totalBet"
	BetTypeNoPay    = "noPay"
)

func procSIWM(ret *sgc7game.Result, gs *sgc7game.GameScene, syms []int, mul int) {
	for i := 0; i < ret.SymbolNums; i++ {
		if goutils.IndexOfIntSlice(syms, gs.Arr[ret.Pos[i*2]][ret.Pos[i*2+1]], 0) >= 0 {
			ret.OtherMul = mul

			ret.CashWin *= mul
			ret.CoinWin *= mul

			return
		}
	}
}

type BasicWinsData struct {
	BasicComponentData
	NextComponent string
}

// OnNewGame -
func (basicWinsData *BasicWinsData) OnNewGame() {
	basicWinsData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (basicWinsData *BasicWinsData) OnNewStep() {
	basicWinsData.BasicComponentData.OnNewStep()

	basicWinsData.NextComponent = ""
}

// BuildPBComponentData
func (basicWinsData *BasicWinsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BasicWinsData{
		BasicComponentData: basicWinsData.BuildPBBasicComponentData(),
		NextComponent:      basicWinsData.NextComponent,
	}

	return pbcd
}

// TriggerFeatureConfig - configuration for trigger feature
type TriggerFeatureConfig struct {
	TargetScene                   string         `yaml:"targetScene" json:"targetScene"`                                     // like basicReels.mstery
	Symbol                        string         `yaml:"symbol" json:"symbol"`                                               // like scatter
	Type                          string         `yaml:"type" json:"type"`                                                   // like scatters
	MinNum                        int            `yaml:"minNum" json:"minNum"`                                               // like 3
	WildSymbols                   []string       `yaml:"wildSymbols" json:"wildSymbols"`                                     // wild etc
	WildSymbolCodes               []int          `yaml:"-" json:"-"`                                                         // wild symbolCode
	SIWMSymbols                   []string       `yaml:"SIWMSymbols" json:"SIWMSymbols"`                                     // SIWM就是如果有符号参与中奖考虑倍数，这里是SIWM的图标
	SIWMSymbolCodes               []int          `yaml:"-" json:"-"`                                                         //
	SIWMMul                       int            `yaml:"SIWMMul" json:"SIWMMul"`                                             // 这里是SIWM的倍数
	Scripts                       string         `yaml:"scripts" json:"scripts"`                                             // scripts
	RespinNum                     int            `yaml:"respinNum" json:"respinNum"`                                         // respin number
	RespinNumWeight               string         `yaml:"respinNumWeight" json:"respinNumWeight"`                             // respin number weight
	RespinNumWithScatterNum       map[int]int    `yaml:"respinNumWithScatterNum" json:"respinNumWithScatterNum"`             // respin number with scatter number
	RespinNumWeightWithScatterNum map[int]string `yaml:"respinNumWeightWithScatterNum" json:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	BetType                       string         `yaml:"betType" json:"betType"`                                             // bet or totalBet
	CountScatterPayAs             string         `yaml:"countScatterPayAs" json:"countScatterPayAs"`                         // countscatter时，按什么符号赔付
	SymbolCodeCountScatterPayAs   int            `yaml:"-" json:"-"`                                                         // countscatter时，按什么符号赔付
	RespinComponent               string         `yaml:"respinComponent" json:"respinComponent"`                             // like fg-spin
	NextComponent                 string         `yaml:"nextComponent" json:"nextComponent"`                                 // next component
	TagSymbolNum                  string         `yaml:"tagSymbolNum" json:"tagSymbolNum"`                                   // 这里可以将symbol数量记下来，别的地方能获取到
	Awards                        []*Award       `yaml:"awards" json:"awards"`                                               // 新的奖励系统
	SymbolAwardsWeights           *AwardsWeights `yaml:"symbolAwardsWeights" json:"symbolAwardsWeights"`                     // 每个中奖符号随机一组奖励
	IsNeedBreak                   bool           `yaml:"isNeedBreak" json:"isNeedBreak"`                                     // 如果触发，需要能break，不继续处理后续的trigger，仅限于当前队列
	PosArea                       []int          `yaml:"posArea" json:"posArea"`                                             // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
	IsSaveRetriggerRespinNum      bool           `yaml:"isSaveRetriggerRespinNum" json:"isSaveRetriggerRespinNum"`           // 如果配置了这个，触发respin以后，会将这次的respinnum缓存下来，后面可以直接用
	ForceToNext                   bool           `yaml:"forceToNext" json:"forceToNext"`                                     // 就算触发了respin，也要先执行next分支
	IsUseTriggerRespin2           bool           `yaml:"isUseTriggerRespin2" json:"isUseTriggerRespin2"`                     // 给true就用triggerRespin2
}

func (tfCfg *TriggerFeatureConfig) onInit(pool *GamePropertyPool) error {
	for _, award := range tfCfg.Awards {
		award.Init()
	}

	if tfCfg.SymbolAwardsWeights != nil {
		tfCfg.SymbolAwardsWeights.Init()
	}

	if tfCfg.CountScatterPayAs != "" {
		tfCfg.SymbolCodeCountScatterPayAs = pool.DefaultPaytables.MapSymbols[tfCfg.CountScatterPayAs]
	} else {
		tfCfg.SymbolCodeCountScatterPayAs = -1
	}

	for _, v := range tfCfg.WildSymbols {
		tfCfg.WildSymbolCodes = append(tfCfg.WildSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range tfCfg.SIWMSymbols {
		tfCfg.SIWMSymbolCodes = append(tfCfg.SIWMSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	return nil
}

// BasicWinsConfig - configuration for BasicWins
type BasicWinsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MainType             string                  `yaml:"mainType" json:"mainType"`             // lines or ways
	BetType              string                  `yaml:"betType" json:"betType"`               // bet or totalBet
	StrCheckWinType      string                  `yaml:"checkWinType" json:"checkWinType"`     // left2right or right2left or all
	CheckWinType         CheckWinType            `yaml:"-" json:"-"`                           //
	SIWMSymbols          []string                `yaml:"SIWMSymbols" json:"SIWMSymbols"`       // SIWM就是如果有符号参与中奖考虑倍数，这里是SIWM的图标
	SIWMSymbolCodes      []int                   `yaml:"-" json:"-"`                           //
	SIWMMul              int                     `yaml:"SIWMMul" json:"SIWMMul"`               // 这里是SIWM的倍数
	ExcludeSymbols       []string                `yaml:"excludeSymbols" json:"excludeSymbols"` // w/s etc
	WildSymbols          []string                `yaml:"wildSymbols" json:"wildSymbols"`       // wild etc
	BeforMain            []*TriggerFeatureConfig `yaml:"beforMain" json:"beforMain"`           // befor the maintype
	AfterMain            []*TriggerFeatureConfig `yaml:"afterMain" json:"afterMain"`           // after the maintype
	BeforMainTriggerName []string                `yaml:"-" json:"-"`                           // befor the maintype
	AfterMainTriggerName []string                `yaml:"-" json:"-"`                           // after the maintype
	IsRespinBreak        bool                    `yaml:"isRespinBreak" json:"isRespinBreak"`   // 如果触发了respin就不执行next，这个是兼容性配置，新游戏应该给true，维持逻辑的一致性
}

type BasicWins struct {
	*BasicComponent `json:"-"`
	Config          *BasicWinsConfig `json:"config"`
	ExcludeSymbols  []int            `json:"-"`
	WildSymbols     []int            `json:"-"`
}

// AddResult -
func (basicWins *BasicWins) ProcTriggerFeature(tf *TriggerFeatureConfig, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, bwd *BasicWinsData) *sgc7game.Result {

	if bwd.NextComponent != "" {
		return nil
	}

	gs, _ := gameProp.GetScene(curpr, tf.TargetScene)

	isTrigger := false
	var ret *sgc7game.Result

	if tf.Type == WinTypeScatters {
		ret = sgc7game.CalcScatter4(gs, gameProp.CurPaytables, gameProp.CurPaytables.MapSymbols[tf.Symbol], gameProp.GetBet(stake, tf.BetType),
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter || goutils.IndexOfIntSlice(tf.WildSymbolCodes, cursymbol, 0) >= 0
			}, true)

		if ret != nil {
			if tf.BetType == BetTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				gameProp.ProcMulti(ret)

				if len(tf.SIWMSymbolCodes) > 0 {
					procSIWM(ret, gs, tf.SIWMSymbolCodes, tf.SIWMMul)
				}
			}

			basicWins.AddResult(curpr, ret, &bwd.BasicComponentData)
			isTrigger = true
		}
	} else if tf.Type == WinTypeCountScatter {
		ret = sgc7game.CalcScatterEx(gs, gameProp.CurPaytables.MapSymbols[tf.Symbol], tf.MinNum, func(scatter int, cursymbol int) bool {
			return cursymbol == scatter || goutils.IndexOfIntSlice(tf.WildSymbolCodes, cursymbol, 0) >= 0
		})

		if ret != nil {
			if tf.BetType == BetTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if tf.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet(stake, tf.BetType)
				}

				gameProp.ProcMulti(ret)

				if len(tf.SIWMSymbolCodes) > 0 {
					procSIWM(ret, gs, tf.SIWMSymbolCodes, tf.SIWMMul)
				}
			}

			basicWins.AddResult(curpr, ret, &bwd.BasicComponentData)
			isTrigger = true
		}
	} else if tf.Type == WinTypeCountScatterInArea {
		ret = sgc7game.CountScatterInArea(gs, gameProp.CurPaytables.MapSymbols[tf.Symbol], tf.MinNum,
			func(x, y int) bool {
				return x >= tf.PosArea[0] && x <= tf.PosArea[1] && y >= tf.PosArea[2] && y <= tf.PosArea[3]
			},
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter || goutils.IndexOfIntSlice(tf.WildSymbolCodes, cursymbol, 0) >= 0
			})

		if ret != nil {
			if tf.BetType == BetTypeNoPay {
				ret.CoinWin = 0
				ret.CashWin = 0
			} else {
				if tf.SymbolCodeCountScatterPayAs > 0 {
					ret.Mul = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CoinWin = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1]
					ret.CashWin = gameProp.CurPaytables.MapPay[tf.SymbolCodeCountScatterPayAs][ret.SymbolNums-1] * gameProp.GetBet(stake, tf.BetType)
				}

				gameProp.ProcMulti(ret)

				if len(tf.SIWMSymbolCodes) > 0 {
					procSIWM(ret, gs, tf.SIWMSymbolCodes, tf.SIWMMul)
				}
			}

			basicWins.AddResult(curpr, ret, &bwd.BasicComponentData)
			isTrigger = true
		}
	}

	if isTrigger {
		if tf.TagSymbolNum != "" {
			gameProp.TagInt(tf.TagSymbolNum, ret.SymbolNums)
		}

		if len(tf.Awards) > 0 {
			gameProp.procAwards(plugin, tf.Awards, curpr, gp)
		}

		if tf.SymbolAwardsWeights != nil {
			for i := 0; i < ret.SymbolNums; i++ {
				node, err := tf.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("BasicWins.ProcTriggerFeature:SymbolAwardsWeights.RandVal",
						zap.Error(err))

					return nil
				}

				gameProp.procAwards(plugin, node.Awards, curpr, gp)
			}
		}

		if tf.RespinComponent != "" {
			if tf.RespinNumWeightWithScatterNum != nil {
				v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, tf.RespinNumWeightWithScatterNum[ret.SymbolNums], basicWins.Config.UseFileMapping, tf.RespinComponent, tf.IsUseTriggerRespin2)
				if err != nil {
					goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
						zap.Error(err))

					return nil
				}

				ret.Type = sgc7game.RTFreeGame
				ret.Value = v
			} else if len(tf.RespinNumWithScatterNum) > 0 {
				gameProp.TriggerRespin(plugin, curpr, gp, tf.RespinNumWithScatterNum[ret.SymbolNums], tf.RespinComponent, tf.IsUseTriggerRespin2)

				ret.Type = sgc7game.RTFreeGame
				ret.Value = tf.RespinNumWithScatterNum[ret.SymbolNums]
			} else if tf.RespinNumWeight != "" {
				v, err := gameProp.TriggerRespinWithWeights(curpr, gp, plugin, tf.RespinNumWeight, basicWins.Config.UseFileMapping, tf.RespinComponent, tf.IsUseTriggerRespin2)
				if err != nil {
					goutils.Error("BasicWins.ProcTriggerFeature:TriggerRespinWithWeights",
						zap.Error(err))

					return nil
				}

				ret.Type = sgc7game.RTFreeGame
				ret.Value = v
			} else if tf.RespinNum > 0 {
				gameProp.TriggerRespin(plugin, curpr, gp, tf.RespinNum, tf.RespinComponent, tf.IsUseTriggerRespin2)

				ret.Type = sgc7game.RTFreeGame
				ret.Value = tf.RespinNum
			} else {
				ret.Type = sgc7game.RTFreeGame
				ret.Value = -1
			}

			if tf.ForceToNext {
				bwd.NextComponent = tf.NextComponent
			}
		} else if tf.NextComponent != "" {
			bwd.NextComponent = tf.NextComponent
		}
	}

	return ret
}

// Init -
func (basicWins *BasicWins) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicWins.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BasicWinsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicWins.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return basicWins.InitEx(cfg, pool)
}

// InitEx -
func (basicWins *BasicWins) InitEx(cfg any, pool *GamePropertyPool) error {
	basicWins.Config = cfg.(*BasicWinsConfig)
	basicWins.Config.ComponentType = BasicWinsTypeName

	basicWins.Config.CheckWinType = ParseCheckWinType(basicWins.Config.StrCheckWinType)

	for _, v := range basicWins.Config.ExcludeSymbols {
		basicWins.ExcludeSymbols = append(basicWins.ExcludeSymbols, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range basicWins.Config.WildSymbols {
		basicWins.WildSymbols = append(basicWins.WildSymbols, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range basicWins.Config.SIWMSymbols {
		basicWins.Config.SIWMSymbolCodes = append(basicWins.Config.SIWMSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range basicWins.Config.BeforMain {
		v.onInit(pool)
	}

	for _, v := range basicWins.Config.AfterMain {
		v.onInit(pool)
	}

	basicWins.onInit(&basicWins.Config.BasicComponentConfig)

	return nil
}

// playgame
func (basicWins *BasicWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	basicWins.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bwd := gameProp.MapComponentData[basicWins.Name].(*BasicWinsData)

	rets := []*sgc7game.Result{}

	for _, v := range basicWins.Config.BeforMain {
		ret := basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, bwd)
		if ret != nil && v.IsNeedBreak {
			break
		}
	}

	gs := basicWins.GetTargetScene(gameProp, curpr, &bwd.BasicComponentData, "")

	if basicWins.Config.MainType == WinTypeWays {
		if basicWins.Config.BasicComponentConfig.TargetOtherScene != "" {
			os := basicWins.GetTargetOtherScene(gameProp, curpr, &bwd.BasicComponentData)

			if os != nil {
				currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet(stake, basicWins.Config.BetType),
					func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
						return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
					}, func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
					}, func(cursymbol int, startsymbol int) bool {
						if cursymbol == startsymbol {
							return true
						}

						return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
					}, func(x, y int) int {
						return os.Arr[x][y]
					})

				for _, v := range currets {
					gameProp.ProcMulti(v)

					if len(basicWins.Config.SIWMSymbolCodes) > 0 {
						procSIWM(v, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
					}
				}

				rets = append(rets, currets...)
			} else {
				currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, gameProp.GetBet(stake, basicWins.Config.BetType),
					func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
						return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
					}, func(cursymbol int) bool {
						return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
					}, func(cursymbol int, startsymbol int) bool {
						if cursymbol == startsymbol {
							return true
						}

						return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
					}, func(x, y int) int {
						return 1
					})

				for _, v := range currets {
					gameProp.ProcMulti(v)

					if len(basicWins.Config.SIWMSymbolCodes) > 0 {
						procSIWM(v, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
					}
				}

				rets = append(rets, currets...)
			}
		} else {
			currets := sgc7game.CalcFullLineEx(gs, gameProp.CurPaytables, gameProp.GetBet(stake, basicWins.Config.BetType),
				func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
					return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
				}, func(cursymbol int) bool {
					return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
				}, func(cursymbol int, startsymbol int) bool {
					if cursymbol == startsymbol {
						return true
					}

					return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
				})

			for _, v := range currets {
				gameProp.ProcMulti(v)

				if len(basicWins.Config.SIWMSymbolCodes) > 0 {
					procSIWM(v, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
				}
			}

			rets = append(rets, currets...)
		}
	} else if basicWins.Config.MainType == WinTypeLines {
		isDone := false
		if basicWins.Config.BasicComponentConfig.TargetOtherScene != "" {
			os := basicWins.GetTargetOtherScene(gameProp, curpr, &bwd.BasicComponentData)

			if os != nil {
				isDone = true

				for i, v := range gameProp.CurLineData.Lines {
					isTriggerFull := false
					if basicWins.Config.CheckWinType != CheckWinTypeRightLeft {
						ret := sgc7game.CalcLine2(gs, gameProp.CurPaytables, v, gameProp.GetBet(stake, basicWins.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return os.Arr[x][y]
							})
						if ret != nil {
							ret.LineIndex = i

							gameProp.ProcMulti(ret)

							if len(basicWins.Config.SIWMSymbolCodes) > 0 {
								procSIWM(ret, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
							}

							rets = append(rets, ret)

							if ret.SymbolNums == gs.Width {
								isTriggerFull = true
							}
						}
					}

					if !isTriggerFull && basicWins.Config.CheckWinType != CheckWinTypeLeftRight {
						ret := sgc7game.CalcLineRL2(gs, gameProp.CurPaytables, v, gameProp.GetBet(stake, basicWins.Config.BetType),
							func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
							}, func(cursymbol int) bool {
								return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
							}, func(cursymbol int, startsymbol int) bool {
								if cursymbol == startsymbol {
									return true
								}

								return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
							}, func(cursymbol int) int {
								return cursymbol
							}, func(x, y int) int {
								return os.Arr[x][y]
							})
						if ret != nil {
							ret.LineIndex = i

							gameProp.ProcMulti(ret)

							if len(basicWins.Config.SIWMSymbolCodes) > 0 {
								procSIWM(ret, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
							}

							rets = append(rets, ret)
						}
					}
				}
			}
		}

		if !isDone {
			for i, v := range gameProp.CurLineData.Lines {
				if basicWins.Config.CheckWinType != CheckWinTypeRightLeft {
					ret := sgc7game.CalcLineEx(gs, gameProp.CurPaytables, v, gameProp.GetBet(stake, basicWins.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
						}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
							return 1
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						if len(basicWins.Config.SIWMSymbolCodes) > 0 {
							procSIWM(ret, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
						}

						rets = append(rets, ret)
					}
				}

				if basicWins.Config.CheckWinType != CheckWinTypeLeftRight {
					ret := sgc7game.CalcLineRLEx(gs, gameProp.CurPaytables, v, gameProp.GetBet(stake, basicWins.Config.BetType),
						func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
						}, func(cursymbol int) bool {
							return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
						}, func(cursymbol int, startsymbol int) bool {
							if cursymbol == startsymbol {
								return true
							}

							return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
						}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
							return 1
						}, func(cursymbol int) int {
							return cursymbol
						})
					if ret != nil {
						ret.LineIndex = i

						gameProp.ProcMulti(ret)

						if len(basicWins.Config.SIWMSymbolCodes) > 0 {
							procSIWM(ret, gs, basicWins.Config.SIWMSymbolCodes, basicWins.Config.SIWMMul)
						}

						rets = append(rets, ret)
					}
				}
			}
		}
	}

	for _, v := range basicWins.Config.AfterMain {
		ret := basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, bwd)
		if ret != nil && v.IsNeedBreak {
			break
		}
	}

	for _, v := range rets {
		basicWins.AddResult(curpr, v, &bwd.BasicComponentData)
	}

	if basicWins.Config.IsRespinBreak {
		if gp.NextStepFirstComponent != "" {
			gameProp.SetStrVal(GamePropNextComponent, "")

			return nil
		}
	}

	basicWins.onStepEnd(gameProp, curpr, gp, bwd.NextComponent)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicWins *BasicWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[basicWins.Name].(*BasicWinsData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// OnStatsWithPB -
func (basicWins *BasicWins) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.BasicWinsData)
	if !isok {
		goutils.Error("BasicWins.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return basicWins.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// OnStats
func (basicWins *BasicWins) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	wins := int64(0)
	isTrigger := false

	for _, v := range lst {
		gp, isok := v.CurGameModParams.(*GameParams)
		if isok {
			curComponent, isok := gp.MapComponentMsgs[basicWins.Name]
			if isok {
				curwins, err := basicWins.OnStatsWithPB(feature, curComponent, v)
				if err != nil {
					goutils.Error("BasicWins.OnStats",
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
func (basicWins *BasicWins) NewComponentData() IComponentData {
	return &BasicWinsData{}
}

// EachUsedResults -
func (basicWins *BasicWins) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.BasicWinsData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicWins.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewBasicWins(name string) IComponent {
	basicWins := &BasicWins{
		BasicComponent: NewBasicComponent(name),
	}

	return basicWins
}
