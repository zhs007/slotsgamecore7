package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	BasicComponentConfig `yaml:",inline"`
	ReelSetsWeight       string `yaml:"reelSetWeight"`
	IsFGMainSpin         bool   `yaml:"isFGMainSpin"`
}

type BasicReels struct {
	*BasicComponent
	Config         *BasicReelsConfig
	ReelSetWeights *sgc7game.ValWeights2
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

	basicReels.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (basicReels *BasicReels) OnNewGame(gameProp *GameProperty) error {
	return nil
}

// OnNewStep -
func (basicReels *BasicReels) OnNewStep(gameProp *GameProperty) error {

	basicReels.BasicComponent.OnNewStep()

	return nil
}

// playgame
func (basicReels *BasicReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	if basicReels.ReelSetWeights != nil {
		val, si, err := basicReels.ReelSetWeights.RandValEx(plugin)
		if err != nil {
			goutils.Error("BasicReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Error(err))

			return err
		}

		basicReels.AddRNG(gameProp, si)

		rd, isok := gameProp.Config.MapReels[val.String()]
		if !isok {
			goutils.Error("BasicReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
	}

	sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	if err != nil {
		goutils.Error("BasicReels.OnPlayGame:NewGameScene",
			zap.Error(err))

		return err
	}

	sc.RandReelsWithReelData(gameProp.CurReels, plugin)

	basicReels.AddScene(gameProp, curpr, sc)

	if basicReels.Config.IsFGMainSpin {
		gameProp.OnFGSpin()
	}

	basicReels.onStepEnd(gameProp, curpr, gp)

	basicReels.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels *BasicReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(basicReels.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[basicReels.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (basicReels *BasicReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewBasicReels(name string) IComponent {
	basicReels := &BasicReels{
		BasicComponent: NewBasicComponent(name),
	}

	return basicReels
}
