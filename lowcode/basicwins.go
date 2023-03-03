package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
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

// TriggerFeatureConfig - configuration for trigger feature
type TriggerFeatureConfig struct {
	TargetScene          string `yaml:"targetScene"`          // basicReels.mstery
	Symbol               string `yaml:"symbol"`               // like scatter
	Type                 string `yaml:"type"`                 // like scatters
	MinNum               int    `yaml:"minNum"`               // like 3
	Scripts              string `yaml:"scripts"`              // scripts
	FGNumWeight          string `yaml:"FGNumWeight"`          // FG number weight
	IsTriggerFG          bool   `yaml:"isTriggerFG"`          // is trigger FG
	BetType              string `yaml:"betType"`              // bet or totalBet
	RespinFirstComponent string `yaml:"respinFirstComponent"` // like fg-spin
}

// BasicWinsConfig - configuration for BasicWins
type BasicWinsConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MainType             string                  `yaml:"mainType"`       // lines or ways
	BetType              string                  `yaml:"betType"`        // bet or totalBet
	TargetScene          string                  `yaml:"targetScene"`    // basicReels.mstery
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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) {
	gs := gameProp.GetScene(curpr, tf.TargetScene)

	isTrigger := false
	if tf.Type == WinTypeScatters {
		ret := sgc7game.CalcScatter4(gs, gameProp.CurPaytables, gameProp.CurPaytables.MapSymbols[tf.Symbol], GetBet(stake, tf.BetType),
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter
			}, true)

		if ret != nil {
			basicWins.AddResult(curpr, ret)
			isTrigger = true
		}
	} else if tf.Type == WinTypeCountScatter {
		ret := sgc7game.CalcScatterEx(gs, gameProp.CurPaytables.MapSymbols[tf.Symbol], tf.MinNum, func(scatter int, cursymbol int) bool {
			return cursymbol == scatter
		})

		if ret != nil {
			basicWins.AddResult(curpr, ret)
			isTrigger = true
		}
	}

	if isTrigger {
		if tf.IsTriggerFG {
			gameProp.TriggerFGWithWeights(curpr, gp, plugin, tf.FGNumWeight, tf.RespinFirstComponent)
		}
	}
}

// Init -
func (basicWins *BasicWins) Init(fn string, gameProp *GameProperty) error {
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
		basicWins.ExcludeSymbols = append(basicWins.ExcludeSymbols, gameProp.CurPaytables.MapSymbols[v])
	}

	for _, v := range cfg.WildSymbols {
		basicWins.WildSymbols = append(basicWins.WildSymbols, gameProp.CurPaytables.MapSymbols[v])
	}

	return nil
}

// OnNewGame -
func (basicWins *BasicWins) OnNewGame(gameProp *GameProperty) error {
	return nil
}

// OnNewStep -
func (basicWins *BasicWins) OnNewStep(gameProp *GameProperty) error {

	basicWins.BasicComponent.OnNewStep()

	return nil
}

// playgame
func (basicWins *BasicWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	for _, v := range basicWins.Config.BeforMain {
		basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)
	}

	gs := gameProp.GetScene(curpr, basicWins.Config.TargetScene)

	if basicWins.Config.MainType == WinTypeWays {
		rets := sgc7game.CalcFullLineEx2(gs, gameProp.CurPaytables, GetBet(stake, basicWins.Config.BetType), func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
			return goutils.IndexOfIntSlice(basicWins.ExcludeSymbols, cursymbol, 0) < 0
		}, func(cursymbol int) bool {
			return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
		}, func(cursymbol int, startsymbol int) bool {
			if cursymbol == startsymbol {
				return true
			}

			return goutils.IndexOfIntSlice(basicWins.WildSymbols, cursymbol, 0) >= 0
		})

		for _, v := range rets {
			basicWins.AddResult(curpr, v)
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

				basicWins.AddResult(curpr, ret)
			}
		}
	}

	for _, v := range basicWins.Config.AfterMain {
		basicWins.ProcTriggerFeature(v, gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)
	}

	gameProp.SetStrVal(GamePropNextComponent, basicWins.Config.DefaultNextComponent)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicWins *BasicWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(basicWins.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

func NewBasicWins(name string) IComponent {
	basicWins := &BasicWins{
		BasicComponent: NewBasicComponent(name),
	}

	return basicWins
}
