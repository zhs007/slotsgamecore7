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

const ComponentTriggerTypeName = "componentTrigger"

// ComponentTriggerConfig - configuration for ComponentTrigger
type ComponentTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	CheckRunComponents   []string `yaml:"checkRunComponents" json:"checkRunComponents"` // 这一组components只要有1个已经运行过就算触发
	JumpToComponent      string   `yaml:"jumpToComponent" json:"jumpToComponent"`       // jump to
	IsReverse            bool     `yaml:"isReverse" json:"isReverse"`                   // 如果isReverse，表示判定为否才触发
}

type ComponentTrigger struct {
	*BasicComponent `json:"-"`
	Config          *ComponentTriggerConfig `json:"config"`
}

// Init -
func (componentTrigger *ComponentTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ComponentTrigger.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ComponentTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ComponentTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return componentTrigger.InitEx(cfg, pool)
}

// InitEx -
func (componentTrigger *ComponentTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	componentTrigger.Config = cfg.(*ComponentTriggerConfig)
	componentTrigger.Config.ComponentType = ComponentTriggerTypeName

	componentTrigger.onInit(&componentTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (componentTrigger *ComponentTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	componentTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	isTrigger := false
	for _, cn := range componentTrigger.Config.CheckRunComponents {
		if gameProp.InHistoryComponents(cn) {

			isTrigger = true
			break
		}
	}

	if componentTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	if isTrigger {
		componentTrigger.onStepEnd(gameProp, curpr, gp, componentTrigger.Config.JumpToComponent)
	} else {
		componentTrigger.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (componentTrigger *ComponentTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (componentTrigger *ComponentTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewComponentTrigger(name string) IComponent {
	return &ComponentTrigger{
		BasicComponent: NewBasicComponent(name),
	}
}
