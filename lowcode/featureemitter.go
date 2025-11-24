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

const FeatureEmitterTypeName = "featureEmitter"

// FeatureEmitterConfig - configuration for FeatureEmitter (placeholder)
type FeatureEmitterConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Collect              string              `yaml:"collect" json:"collect"`
	ConsumedAmount       int                 `yaml:"consumedAmount" json:"consumedAmount"`
	JumpToComponent      string              `yaml:"jumpToComponent" json:"jumpToComponent"`
	MapControllers       map[string][]*Award `yaml:"mapControllers" json:"mapControllers"`
}

func (cfg *FeatureEmitterConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

// FeatureEmitter - placeholder component
type FeatureEmitter struct {
	*BasicComponent `json:"-"`
	Config          *FeatureEmitterConfig `json:"config"`
}

// Init - read yaml file
func (fe *FeatureEmitter) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FeatureEmitter.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &FeatureEmitterConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FeatureEmitter.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return fe.InitEx(cfg, pool)
}

// InitEx - initialize from config
func (fe *FeatureEmitter) InitEx(cfg any, pool *GamePropertyPool) error {
	fe.Config = cfg.(*FeatureEmitterConfig)
	fe.Config.ComponentType = FeatureEmitterTypeName

	for _, arr := range fe.Config.MapControllers {
		for _, aw := range arr {
			aw.Init()
		}
	}

	fe.onInit(&fe.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (fe *FeatureEmitter) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := fe.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

// OnPlayGame - placeholder behavior: does nothing
func (fe *FeatureEmitter) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	if fe.Config.Collect != "" {
		icollectdata := gameProp.GetComponentDataWithName(fe.Config.Collect)
		if icollectdata == nil {
			goutils.Error("FeatureEmitter.OnPlayGame:Collect:icollectdata==nil",
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		val := icollectdata.GetOutput()

		if val >= fe.Config.ConsumedAmount {

			icollectdata.ChgConfigIntVal(CCVValueNumNow, -fe.Config.ConsumedAmount)

			fe.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

			nc := fe.onStepEnd(gameProp, curpr, gp, fe.Config.JumpToComponent)

			return nc, nil
		}
	}

	nc := fe.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - output to asciigame (no-op placeholder)
func (fe *FeatureEmitter) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewFeatureEmitter(name string) IComponent {
	return &FeatureEmitter{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "collect": "co-collect",
// "consumedAmount": 1
type jsonFeatureEmitter struct {
	Collect        string `json:"collect"`
	ConsumedAmount int    `json:"consumedAmount"`
}

func (jcfg *jsonFeatureEmitter) build() *FeatureEmitterConfig {
	cfg := &FeatureEmitterConfig{
		Collect:        jcfg.Collect,
		ConsumedAmount: jcfg.ConsumedAmount,
	}

	return cfg
}

func parseFeatureEmitter(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseFeatureEmitter:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseFeatureEmitter:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonFeatureEmitter{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseFeatureEmitter:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseDropDownSymbols2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: FeatureEmitterTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
