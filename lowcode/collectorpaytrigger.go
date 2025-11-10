package lowcode

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const CollectorPayTriggerTypeName = "collectorPayTrigger"

// CollectorPayTriggerConfig - configuration for CollectorPayTrigger
type CollectorPayTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Collector            string   `yaml:"collector" json:"collector"`             // collector component name
	Value                int      `yaml:"value" json:"value"`                     // threshold value to trigger
	Compare              string   `yaml:"compare" json:"compare"`                 // compare operator: ge(>=) or eq(==), default ge
	Awards               []*Award `yaml:"awards" json:"awards"`                   // awards to proc when trigger
	JumpToComponent      string   `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
	ForceToNext          bool     `yaml:"forceToNext" json:"forceToNext"`
	IsReverse            bool     `yaml:"isReverse" json:"isReverse"`
}

func (cfg *CollectorPayTriggerConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type CollectorPayTrigger struct {
	*BasicComponent `json:"-"`
	Config          *CollectorPayTriggerConfig `json:"config"`
}

// Init - load from file
func (cpt *CollectorPayTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CollectorPayTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CollectorPayTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CollectorPayTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return cpt.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (cpt *CollectorPayTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	cpt.Config = cfg.(*CollectorPayTriggerConfig)
	cpt.Config.ComponentType = CollectorPayTriggerTypeName

	for _, a := range cpt.Config.Awards {
		if a != nil {
			a.Init()
		}
	}

	if cpt.Config.Compare == "" {
		cpt.Config.Compare = "ge"
	}

	cpt.onInit(&cpt.Config.BasicComponentConfig)

	return nil
}

// OnPlayGame - check collector value and proc awards when reach threshold
func (cpt *CollectorPayTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// find collector component data
	cd := gameProp.GetComponentDataWithName(cpt.Config.Collector)
	if cd == nil {
		goutils.Error("CollectorPayTrigger.OnPlayGame:GetComponentDataWithName",
			slog.String("collector", cpt.Config.Collector),
			goutils.Err(ErrInvalidComponent))
		return "", ErrInvalidComponent
	}

	val, _ := cd.GetValEx(CVValue, GCVTypeNormal)

	isTrigger := false
	switch cpt.Config.Compare {
	case "eq":
		isTrigger = (val == cpt.Config.Value)
	default: // ge
		isTrigger = (val >= cpt.Config.Value)
	}

	if cpt.Config.IsReverse {
		isTrigger = !isTrigger
	}

	nc := ""
	if isTrigger {
		if cpt.Config.Awards != nil {
			gameProp.procAwards(plugin, cpt.Config.Awards, curpr, gp)
		}

		if cpt.Config.ForceToNext {
			nc = cpt.onStepEnd(gameProp, curpr, gp, "")
		} else {
			nc = cpt.onStepEnd(gameProp, curpr, gp, cpt.Config.JumpToComponent)
		}
	} else {
		nc = cpt.onStepEnd(gameProp, curpr, gp, "")
	}

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (cpt *CollectorPayTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewCollectorPayTrigger(name string) IComponent {
	return &CollectorPayTrigger{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// json representation used by editor
type jsonCollectorPayTrigger struct {
	Collector string `json:"collector"`
	Value     int    `json:"value"`
	Compare   string `json:"compare"`
}

func (j *jsonCollectorPayTrigger) build() *CollectorPayTriggerConfig {
	return &CollectorPayTriggerConfig{
		Collector: j.Collector,
		Value:     j.Value,
		Compare:   j.Compare,
	}
}

func parseCollectorPayTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCollectorPayTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCollectorPayTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CollectorPayTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
