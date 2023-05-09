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

const (
	WinTypeLines        = "lines"
	WinTypeWays         = "ways"
	WinTypeScatters     = "scatters"
	WinTypeCountScatter = "countscatter"

	BetTypeNormal   = "bet"
	BetTypeTotalBet = "totalBet"
)

func GetBet(stake *sgc7game.Stake, bettype string) int {
	if bettype == BetTypeTotalBet {
		return int(stake.CashBet)
	}

	return int(stake.CoinBet)
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
	TargetScene                   string         `yaml:"targetScene"`                   // like basicReels.mstery
	Symbol                        string         `yaml:"symbol"`                        // like scatter
	Type                          string         `yaml:"type"`                          // like scatters
	MinNum                        int            `yaml:"minNum"`                        // like 3
	Scripts                       string         `yaml:"scripts"`                       // scripts
	RespinNum                     int            `yaml:"respinNum"`                     // respin number
	RespinNumWeight               string         `yaml:"respinNumWeight"`               // respin number weight
	RespinNumWithScatterNum       map[int]int    `yaml:"respinNumWithScatterNum"`       // respin number with scatter number
	RespinNumWeightWithScatterNum map[int]string `yaml:"respinNumWeightWithScatterNum"` // respin number weight with scatter number
	BetType                       string         `yaml:"betType"`                       // bet or totalBet
	RespinComponent               string         `yaml:"respinComponent"`               // like fg-spin
	NextComponent                 string         `yaml:"nextComponent"`                 // next component
	TagSymbolNum                  string         `yaml:"tagSymbolNum"`                  // 这里可以将symbol数量记下来，别的地方能获取到
	Awards                        []*Award       `yaml:"awards"`                        // 新的奖励系统
	SymbolAwardsWeights           *AwardsWeights `yaml:"symbolAwardsWeights"`           // 每个中奖符号随机一组奖励
}

// BasicWinsConfig - configuration for BasicWins
type BasicWinsConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MainType             string                  `yaml:"mainType"`       // lines or ways
	BetType              string                  `yaml:"betType"`        // bet or totalBet
	ExcludeSymbols       []string                `yaml:"excludeSymbols"` // w/s etc
	WildSymbols          []string                `yaml:"wildSymbols"`    // wild etc
	BeforMain            []*TriggerFeatureConfig `yaml:"beforMain"`      // befor the maintype
	AfterMain            []*TriggerFeatureConfig `yaml:"afterMain"`      // after the maintype
}

type BasicWins struct {
	*BasicComponent
	Config         *BasicWinsConfig
	ExcludeSymbols []int
	WildSymbols    []int
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
		ret = sgc7game.CalcScatter4(gs, gameProp.CurPaytables, gameProp.CurPaytables.MapSymbols[tf.Symbol], GetBet(stake, tf.BetType),
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter
			}, true)

		if ret != nil {
			gameProp.ProcMulti(ret)

			basicWins.AddResult(curpr, ret, &bwd.BasicComponentData)
			isTrigger = true
		}
	} else if tf.Type == WinTypeCountScatter {
		ret = sgc7game.CalcScatterEx(gs, gameProp.CurPaytables.MapSymbols[tf.Symbol], tf.MinNum, func(scatter int, cursymbol int) bool {
			return cursymbol == scatter
		})

		if ret != nil {
			gameProp.ProcMulti(ret)

			basicWins.AddResult(curpr, ret, &bwd.BasicComponentData)
			isTrigger = true
		}
	}

	if isTrigger {
		if tf.TagSymbolNum != "" {
			gameProp.TagInt(tf.TagSymbolNum, ret.SymbolNums)
		}

		if len(tf.Awards) > 0 {
			gameProp.procAwards(tf.Awards, curpr, gp)
		}

		if tf.SymbolAwardsWeights != nil {
			for i := 0; i < ret.SymbolNums; i++ {
				node, err := tf.SymbolAwardsWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("BasicWins.ProcTriggerFeature:SymbolAwardsWeights.RandVal",
						zap.Error(err))

					return nil
				}

				gameProp.procAwards(node.Awards, curpr, gp)
			}
		}

		if tf.RespinComponent != "" {
			if tf.RespinNumWeightWithScatterNum != nil {
				gameProp.TriggerRespinWithWeights(curpr, gp, plugin, tf.RespinNumWeightWithScatterNum[ret.SymbolNums], tf.RespinComponent)
			} else if len(tf.RespinNumWithScatterNum) > 0 {
				gameProp.TriggerRespin(curpr, gp, tf.RespinNumWithScatterNum[ret.SymbolNums], tf.RespinComponent)
			} else if tf.RespinNumWeight != "" {
				gameProp.TriggerRespinWithWeights(curpr, gp, plugin, tf.RespinNumWeight, tf.RespinComponent)
			} else {
				gameProp.TriggerRespin(curpr, gp, tf.RespinNum, tf.RespinComponent)
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

	basicWins.Config = cfg

	for _, v := range cfg.ExcludeSymbols {
		basicWins.ExcludeSymbols = append(basicWins.ExcludeSymbols, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range cfg.WildSymbols {
		basicWins.WildSymbols = append(basicWins.WildSymbols, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range cfg.BeforMain {
		for _, award := range v.Awards {
			award.Init()
		}

		if v.SymbolAwardsWeights != nil {
			v.SymbolAwardsWeights.Init()
		}
	}

	for _, v := range cfg.AfterMain {
		for _, award := range v.Awards {
			award.Init()
		}

		if v.SymbolAwardsWeights != nil {
			v.SymbolAwardsWeights.Init()
		}
	}

	basicWins.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (basicWins *BasicWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	basicWins.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bwd := gameProp.MapComponentData[basicWins.Name].(*BasicWinsData)

	rets := []*sgc7game.Result{}

	for _, v := range basicWins.Config.BeforMain {
		basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, bwd)
	}

	gs := basicWins.GetTargetScene(gameProp, curpr, &bwd.BasicComponentData, "")

	if basicWins.Config.MainType == WinTypeWays {
		if basicWins.Config.BasicComponentConfig.TargetOtherScene != "" {
			os := basicWins.GetTargetOtherScene(gameProp, curpr, &bwd.BasicComponentData)

			if os != nil {
				currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, GetBet(stake, basicWins.Config.BetType), func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
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
				}

				rets = append(rets, currets...)
			} else {
				currets := sgc7game.CalcFullLineExWithMulti(gs, gameProp.CurPaytables, GetBet(stake, basicWins.Config.BetType), func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
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
				}

				rets = append(rets, currets...)
			}
		} else {
			currets := sgc7game.CalcFullLineEx(gs, gameProp.CurPaytables, GetBet(stake, basicWins.Config.BetType), func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
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
			}

			rets = append(rets, currets...)
		}
	} else if basicWins.Config.MainType == WinTypeLines {
		for i, v := range gameProp.CurLineData.Lines {
			ret := sgc7game.CalcLineEx(gs, gameProp.CurPaytables, v, GetBet(stake, basicWins.Config.BetType), func(cursymbol int) bool {
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

				rets = append(rets, ret)

				// basicWins.AddResult(curpr, ret, cd)
			}
		}
	}

	if basicWins.Config.BasicComponentConfig.TargetOtherScene != "" && basicWins.Config.MainType == WinTypeLines {
		os := basicWins.GetTargetOtherScene(gameProp, curpr, &bwd.BasicComponentData)

		if os != nil {
			for _, v := range rets {
				mul := 1

				for i := 0; i < len(v.Pos)/2; i++ {
					mul *= os.Arr[v.Pos[i*2]][v.Pos[i*2+1]]
				}

				v.OtherMul = mul

				v.CashWin *= mul
				v.CoinWin *= mul
			}
		}
	}

	for _, v := range basicWins.Config.AfterMain {
		basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs, bwd)
	}

	for _, v := range rets {
		basicWins.AddResult(curpr, v, &bwd.BasicComponentData)
	}

	basicWins.onStepEnd(gameProp, curpr, gp, bwd.NextComponent)

	// gp.AddComponentData(basicWins.Name, bwd)

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
func (basicWins *BasicWins) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.BasicWinsData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicWins.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
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
			curComponent, isok := gp.MapComponents[basicWins.Name]
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
