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

const WeightTriggerTypeName = "weightTrigger"

// WeightTriggerConfig - configuration for WeightTrigger
type WeightTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	NextComponents       []string `yaml:"nextComponents" json:"nextComponents"`
	RespinNums           []int    `yaml:"respinNums" json:"respinNums"`
	WeightSet            string   `yaml:"weightSet" json:"weightSet"`
	IsUseTriggerRespin2  bool     `yaml:"isUseTriggerRespin2" json:"isUseTriggerRespin2"` // 给true就用triggerRespin2
}

type WeightTrigger struct {
	*BasicComponent `json:"-"`
	Config          *WeightTriggerConfig  `json:"config"`
	WeightSet       *sgc7game.ValWeights2 `json:"-"`
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

	return weightTrigger.InitEx(cfg, pool)
}

// InitEx -
func (weightTrigger *WeightTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	weightTrigger.Config = cfg.(*WeightTriggerConfig)
	weightTrigger.Config.ComponentType = WeightTriggerTypeName

	if weightTrigger.Config.WeightSet != "" {
		vw2, err := pool.LoadIntWeights(weightTrigger.Config.WeightSet, weightTrigger.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightTrigger.Init:LoadValWeights",
				zap.String("Weight", weightTrigger.Config.WeightSet),
				zap.Error(err))

			return err
		}

		weightTrigger.WeightSet = vw2
	}

	weightTrigger.onInit(&weightTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (weightTrigger *WeightTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	// cd := gameProp.MapComponentData[weightTrigger.Name].(*BasicComponentData)

	rv, err := weightTrigger.WeightSet.RandVal(plugin)
	if err != nil {
		goutils.Error("WeightTrigger.OnPlayGame:RandVal",
			zap.Error(err))

		return "", err
	}

	setIndex := rv.Int()

	if len(weightTrigger.Config.RespinNums) == len(weightTrigger.Config.NextComponents) {
		if weightTrigger.Config.RespinNums[setIndex] > 0 {
			gameProp.TriggerRespin(plugin, curpr, gp, weightTrigger.Config.RespinNums[setIndex], weightTrigger.Config.NextComponents[setIndex], weightTrigger.Config.IsUseTriggerRespin2)
		}
		// 	weightTrigger.onStepEnd(gameProp, curpr, gp, "")
		// } else {
		// 	weightTrigger.onStepEnd(gameProp, curpr, gp, weightTrigger.Config.NextComponents[setIndex])
		// }
	}

	nc := weightTrigger.onStepEnd(gameProp, curpr, gp, weightTrigger.Config.NextComponents[setIndex])

	// gp.AddComponentData(weightTrigger.Name, cd)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightTrigger *WeightTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// OnStats
func (weightTrigger *WeightTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewWeightTrigger(name string) IComponent {
	return &WeightTrigger{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
