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

const WeightTrigger2TypeName = "weightTrigger2"

// WeightTrigger2Config - configuration for WeightTrigger2
type WeightTrigger2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	TriggerWeight        string                `yaml:"triggerWeight" json:"triggerWeight"`
	TriggerWeightVW      *sgc7game.ValWeights2 `json:"-"`
	JumpToComponent      string                `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
}

type WeightTrigger2 struct {
	*BasicComponent `json:"-"`
	Config          *WeightTrigger2Config `json:"config"`
}

// Init -
func (weightTrigger2 *WeightTrigger2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightTrigger2.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WeightTrigger2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightTrigger2.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return weightTrigger2.InitEx(cfg, pool)
}

// InitEx -
func (weightTrigger2 *WeightTrigger2) InitEx(cfg any, pool *GamePropertyPool) error {
	weightTrigger2.Config = cfg.(*WeightTrigger2Config)
	weightTrigger2.Config.ComponentType = WeightTrigger2TypeName

	if weightTrigger2.Config.TriggerWeight != "" {
		vw2, err := pool.LoadIntWeights(weightTrigger2.Config.TriggerWeight, weightTrigger2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightTrigger2.Init:LoadValWeights",
				zap.String("Weight", weightTrigger2.Config.TriggerWeight),
				zap.Error(err))

			return err
		}

		weightTrigger2.Config.TriggerWeightVW = vw2
	}

	weightTrigger2.onInit(&weightTrigger2.Config.BasicComponentConfig)

	return nil
}

// playgame
func (weightTrigger2 *WeightTrigger2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	weightTrigger2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	rv, err := weightTrigger2.Config.TriggerWeightVW.RandVal(plugin)
	if err != nil {
		goutils.Error("WeightTrigger2.OnPlayGame:RandVal",
			zap.Error(err))

		return err
	}

	if rv.Int() != 0 {
		weightTrigger2.onStepEnd(gameProp, curpr, gp, weightTrigger2.Config.JumpToComponent)
	} else {
		weightTrigger2.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (weightTrigger2 *WeightTrigger2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (weightTrigger2 *WeightTrigger2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewWeightTrigger2(name string) IComponent {
	return &WeightTrigger2{
		BasicComponent: NewBasicComponent(name),
	}
}
