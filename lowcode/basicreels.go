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

// TriggerFeatureConfig - configuration for trigger feature
type TriggerFeatureConfig struct {
	Symbol  string   `yaml:"symbol"`  // like scatter
	Type    string   `yaml:"type"`    // like scatters
	Scripts []string `yaml:"scripts"` // scripts
}

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	MainType       string                  `yaml:"mainType"`       // lines or ways
	ExcludeSymbols []string                `yaml:"excludeSymbols"` // w/s etc
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

	curpr.Scenes = append(curpr.Scenes, sc)

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

		curpr.Scenes = append(curpr.Scenes, sc2)
	}

	return nil
}

// pay
func (basicReels *BasicReels) OnPay(gameProp *GameProperty, curpr *sgc7game.PlayResult, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels *BasicReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	asciigame.OutputScene("initial symbols", pr.Scenes[0], mapSymbolColor)

	if basicReels.MysteryWeights != nil {
		fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(gameProp.MapVal[GamePropCurMystery]))
		asciigame.OutputScene("after symbols", pr.Scenes[1], mapSymbolColor)
	}

	return nil
}

func NewBasicReels() IComponent {
	return &BasicReels{}
}
