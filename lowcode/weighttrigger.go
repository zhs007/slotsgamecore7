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

// WeightTriggerConfig - configuration for WeightTrigger
type WeightTriggerConfig struct {
	BasicComponentConfig `yaml:",inline"`
	NextComponents       []string `yaml:"nextComponents"`
	WeightSet            string   `yaml:"weightSet"`
}

type WeightTrigger struct {
	*BasicComponent
	Config    *WeightTriggerConfig
	WeightSet *sgc7game.ValWeights2
}

// Init -
func (weightTrigger *WeightTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightTrigger.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WeightTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	weightTrigger.Config = cfg

	if cfg.WeightSet != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(cfg.WeightSet), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("WeightTrigger.Init:LoadValWeights2FromExcel",
				zap.String("Weight", cfg.WeightSet),
				zap.Error(err))

			return err
		}

		weightTrigger.WeightSet = vw2
	}

	weightTrigger.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (weightTrigger *WeightTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[weightTrigger.Name].(*BasicComponentData)

	rv, err := weightTrigger.WeightSet.RandVal(plugin)
	if err != nil {
		goutils.Error("WeightTrigger.OnPlayGame:RandVal",
			zap.Error(err))

		return err
	}

	setIndex := rv.Int()

	weightTrigger.onStepEnd(gameProp, curpr, gp, weightTrigger.Config.NextComponents[setIndex])

	gp.AddComponentData(weightTrigger.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (weightTrigger *WeightTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (weightTrigger *WeightTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewWeightTrigger(name string) IComponent {
	return &WeightTrigger{
		BasicComponent: NewBasicComponent(name),
	}
}
