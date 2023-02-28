package lowcode

import (
	"fmt"
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
	IsBeforMystery bool   `yaml:"isBeforMystery"` // is befor mystery
	Symbol         string `yaml:"symbol"`         // like scatter
	Type           string `yaml:"type"`           // like scatters
	MinNum         int    `yaml:"minNum"`         // like 3
	Scripts        string `yaml:"scripts"`        // scripts
	FGNumWeight    string `yaml:"FGNumWeight"`    // FG number weight
	IsTriggerFG    bool   `yaml:"isTriggerFG"`    // is trigger FG
	BetType        string `yaml:"betType"`        // bet or totalBet
}

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	MainType       string                  `yaml:"mainType"`       // lines or ways
	BetType        string                  `yaml:"betType"`        // bet or totalBet
	ExcludeSymbols []string                `yaml:"excludeSymbols"` // w/s etc
	WildSymbols    []string                `yaml:"wildSymbols"`    // wild etc
	ReelSetsWeight string                  `yaml:"reelSetWeight"`
	MysteryWeight  string                  `yaml:"mysteryWeight"`
	Mystery        string                  `yaml:"mystery"`
	BeforMain      []*TriggerFeatureConfig `yaml:"beforMain"` // befor the maintype
	AfterMain      []*TriggerFeatureConfig `yaml:"afterMain"` // after the maintype
}

type BasicReels struct {
	Config         *BasicReelsConfig
	ReelSetWeights *sgc7game.ValWeights2
	MysteryWeights *sgc7game.ValWeights2
	MysterySymbol  int
	UsedScenes     []int
	UsedResults    []int
	ExcludeSymbols []int
	WildSymbols    []int
}

// AddScene -
func (basicReels *BasicReels) AddScene(curpr *sgc7game.PlayResult, sc *sgc7game.GameScene) {
	basicReels.UsedScenes = append(basicReels.UsedScenes, len(curpr.Scenes))

	curpr.Scenes = append(curpr.Scenes, sc)
}

// AddResult -
func (basicReels *BasicReels) AddResult(curpr *sgc7game.PlayResult, ret *sgc7game.Result) {
	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	basicReels.UsedResults = append(basicReels.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

// AddResult -
func (basicReels *BasicReels) ProcTriggerFeature(tf *TriggerFeatureConfig, gameProp *GameProperty, curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) {
	lastsi := basicReels.UsedScenes[len(basicReels.UsedScenes)-1]
	if tf.IsBeforMystery {
		lastsi = basicReels.UsedScenes[0]
	}

	isTrigger := false
	if tf.Type == WinTypeScatters {
		ret := sgc7game.CalcScatter4(curpr.Scenes[lastsi], gameProp.CurPaytables, gameProp.CurPaytables.MapSymbols[tf.Symbol], GetBet(stake, tf.BetType),
			func(scatter int, cursymbol int) bool {
				return cursymbol == scatter
			}, true)

		if ret != nil {
			basicReels.AddResult(curpr, ret)
			isTrigger = true
		}
	} else if tf.Type == WinTypeCountScatter {
		ret := sgc7game.CalcScatterEx(curpr.Scenes[lastsi], gameProp.CurPaytables.MapSymbols[tf.Symbol], tf.MinNum, func(scatter int, cursymbol int) bool {
			return cursymbol == scatter
		})

		if ret != nil {
			basicReels.AddResult(curpr, ret)
			isTrigger = true
		}
	}

	if isTrigger {
		if tf.IsTriggerFG {
			gameProp.TriggerFGWithWeights(tf.FGNumWeight)
		}
	}
}

// Init -
func (basicReels *BasicReels) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BasicReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	basicReels.Config = cfg

	if basicReels.Config.ReelSetsWeight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(basicReels.Config.ReelSetsWeight, "val", "weight", sgc7game.NewStrVal)
		if err != nil {
			goutils.Error("BasicReels.Init:LoadValWeights2FromExcel",
				zap.String("ReelSetsWeight", basicReels.Config.ReelSetsWeight),
				zap.Error(err))

			return err
		}

		basicReels.ReelSetWeights = vw2
	}

	if basicReels.Config.MysteryWeight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(basicReels.Config.MysteryWeight, "val", "weight", gameProp.CurPaytables)
		if err != nil {
			goutils.Error("BasicReels.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", basicReels.Config.MysteryWeight),
				zap.Error(err))

			return err
		}

		basicReels.MysteryWeights = vw2
	}

	basicReels.MysterySymbol = gameProp.CurPaytables.MapSymbols[basicReels.Config.Mystery]

	for _, v := range cfg.ExcludeSymbols {
		basicReels.ExcludeSymbols = append(basicReels.ExcludeSymbols, gameProp.CurPaytables.MapSymbols[v])
	}

	for _, v := range cfg.WildSymbols {
		basicReels.WildSymbols = append(basicReels.WildSymbols, gameProp.CurPaytables.MapSymbols[v])
	}

	return nil
}

// playgame
func (basicReels *BasicReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	if basicReels.ReelSetWeights != nil {
		val, err := basicReels.ReelSetWeights.RandVal(plugin)
		if err != nil {
			goutils.Error("BasicReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Error(err))

			return err
		}

		rd, isok := gameProp.Config.MapReels[val.String()]
		if !isok {
			goutils.Error("BasicReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
	}

	sc, err := sgc7game.NewGameScene(gameProp.MapVal[GamePropWidth], gameProp.MapVal[GamePropHeight])
	if err != nil {
		goutils.Error("BasicReels.OnPlayGame:NewGameScene",
			zap.Error(err))

		return err
	}

	sc.RandReelsWithReelData(gameProp.CurReels, plugin)

	basicReels.AddScene(curpr, sc)

	if basicReels.MysteryWeights != nil {
		curm, err := basicReels.MysteryWeights.RandVal(plugin)
		if err != nil {
			goutils.Error("BasicReels.OnPlayGame:RandVal",
				zap.Error(err))

			return err
		}

		gameProp.MapVal[GamePropCurMystery] = curm.Int()

		sc2 := sc.Clone()
		sc2.ReplaceSymbol(basicReels.MysterySymbol, gameProp.MapVal[GamePropCurMystery])

		basicReels.AddScene(curpr, sc2)
	}

	return nil
}

// pay
func (basicReels *BasicReels) OnPay(gameProp *GameProperty, curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	for _, v := range basicReels.Config.BeforMain {
		basicReels.ProcTriggerFeature(v, gameProp, curpr, plugin, cmd, param, ps, stake, prs)
	}

	if basicReels.Config.MainType == WinTypeWays {
		lastsi := basicReels.UsedScenes[len(basicReels.UsedScenes)-1]

		rets := sgc7game.CalcFullLineEx2(curpr.Scenes[lastsi], gameProp.CurPaytables, GetBet(stake, basicReels.Config.BetType), func(cursymbol int, scene *sgc7game.GameScene, x, y int) bool {
			return goutils.IndexOfIntSlice(basicReels.ExcludeSymbols, cursymbol, 0) < 0
		}, func(cursymbol int) bool {
			return goutils.IndexOfIntSlice(basicReels.WildSymbols, cursymbol, 0) >= 0
		}, func(cursymbol int, startsymbol int) bool {
			if cursymbol == startsymbol {
				return true
			}

			return goutils.IndexOfIntSlice(basicReels.WildSymbols, cursymbol, 0) >= 0
		})

		for _, v := range rets {
			basicReels.AddResult(curpr, v)
		}
	} else if basicReels.Config.MainType == WinTypeLines {
		lastsi := basicReels.UsedScenes[len(basicReels.UsedScenes)-1]

		for i, v := range gameProp.CurLineData.Lines {
			ret := sgc7game.CalcLineEx(curpr.Scenes[lastsi], gameProp.CurPaytables, v, GetBet(stake, basicReels.Config.BetType), func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(basicReels.ExcludeSymbols, cursymbol, 0) < 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(basicReels.WildSymbols, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(basicReels.WildSymbols, cursymbol, 0) >= 0
			}, func(scene *sgc7game.GameScene, result *sgc7game.Result) int {
				return 1
			}, func(cursymbol int) int {
				return cursymbol
			})
			if ret != nil {
				ret.LineIndex = i

				basicReels.AddResult(curpr, ret)
			}
		}
	}

	for _, v := range basicReels.Config.AfterMain {
		basicReels.ProcTriggerFeature(v, gameProp, curpr, plugin, cmd, param, ps, stake, prs)
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels *BasicReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(basicReels.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[basicReels.UsedScenes[0]], mapSymbolColor)

		if basicReels.MysteryWeights != nil {
			fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(gameProp.MapVal[GamePropCurMystery]))
			asciigame.OutputScene("after symbols", pr.Scenes[basicReels.UsedScenes[1]], mapSymbolColor)
		}
	}

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(basicReels.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

func NewBasicReels() IComponent {
	return &BasicReels{}
}
