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

const ComponentValTriggerTypeName = "componentValTrigger"

type OperateType int

const (
	OTEqual        OperateType = 0 // ==
	OTGreaterEqual OperateType = 1 // >=
	OTLessEqual    OperateType = 2 // <=
	OTGreater      OperateType = 3 // >
	OTLess         OperateType = 4 // <
	OTNotEqual     OperateType = 5 // !=
)

func ParseOperateType(str string) OperateType {
	if str == ">=" {
		return OTGreaterEqual
	} else if str == "<=" {
		return OTLessEqual
	} else if str == ">" {
		return OTGreater
	} else if str == "<" {
		return OTLess
	} else if str == "!=" {
		return OTNotEqual
	}

	return OTEqual
}

// ComponentValTriggerConfig - configuration for ComponentValTrigger
type ComponentValTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ComponentVal         string      `yaml:"componentVal" json:"componentVal"`       // 用来检查的值，bg-wins.wins 这样的命名方式
	OperateString        string      `yaml:"operate" json:"operate"`                 // ==/>=/<=/>/</!=
	Operate              OperateType `yaml:"-" json:"-"`                             //
	TargetVal            int         `yaml:"targetVal" json:"targetVal"`             // 目标值
	JumpToComponent      string      `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
}

type ComponentValTrigger struct {
	*BasicComponent `json:"-"`
	Config          *ComponentValTriggerConfig `json:"config"`
}

// Init -
func (componentValTrigger *ComponentValTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ComponentValTrigger.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ComponentValTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ComponentValTrigger.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return componentValTrigger.InitEx(cfg, pool)
}

// InitEx -
func (componentValTrigger *ComponentValTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	componentValTrigger.Config = cfg.(*ComponentValTriggerConfig)
	componentValTrigger.Config.ComponentType = ComponentValTriggerTypeName

	componentValTrigger.Config.Operate = ParseOperateType(componentValTrigger.Config.OperateString)

	componentValTrigger.onInit(&componentValTrigger.Config.BasicComponentConfig)

	return nil
}

// playgame
func (componentValTrigger *ComponentValTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	componentValTrigger.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	isTrigger := false
	val, err := gameProp.GetComponentVal(componentValTrigger.Config.ComponentVal)
	if err != nil {
		goutils.Error("ComponentValTrigger.OnPlayGame:GetComponentVal",
			zap.String("ComponentVal", componentValTrigger.Config.ComponentVal),
			zap.Error(err))

		return err
	}

	switch componentValTrigger.Config.Operate {
	case OTEqual:
		isTrigger = (val == componentValTrigger.Config.TargetVal)
	case OTNotEqual:
		isTrigger = (val != componentValTrigger.Config.TargetVal)
	case OTGreaterEqual:
		isTrigger = (val >= componentValTrigger.Config.TargetVal)
	case OTLessEqual:
		isTrigger = (val <= componentValTrigger.Config.TargetVal)
	case OTGreater:
		isTrigger = (val > componentValTrigger.Config.TargetVal)
	case OTLess:
		isTrigger = (val < componentValTrigger.Config.TargetVal)
	}

	if isTrigger {
		componentValTrigger.onStepEnd(gameProp, curpr, gp, componentValTrigger.Config.JumpToComponent)
	} else {
		componentValTrigger.onStepEnd(gameProp, curpr, gp, "")
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (componentTrigger *ComponentValTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

// OnStats
func (componentTrigger *ComponentValTrigger) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewComponentValTrigger(name string) IComponent {
	return &ComponentValTrigger{
		BasicComponent: NewBasicComponent(name),
	}
}
