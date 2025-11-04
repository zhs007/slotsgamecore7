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

const GenSPGridTypeName = "genSPGrid"

// GenSPGridConfig - configuration for GenSPGrid
type GenSPGridConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Width                int `yaml:"width" json:"width"`
	Height               int `yaml:"height" json:"height"`
}

// SetLinkComponent
func (cfg *GenSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *GenSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *GenSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *GenSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*GenSPGridConfig)
	gen.Config.ComponentType = GenSPGridTypeName

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

// OnPlayGame - minimal implementation: does nothing but advance
func (gen *GenSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// This implementation intentionally does not modify the play result.
	// It simply ends this step and returns ErrComponentDoNothing so the caller
	// can proceed. Later this can be extended to actually generate SPGrid scenes.
	curpr.SPGrid[gen.GetName()] = []*sgc7game.GameScene{}

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *GenSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewGenSPGrid(name string) IComponent {
	return &GenSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// json representation used by editor
type jsonGenSPGrid struct {
	Width  int `json:"Width"`
	Height int `json:"Height"`
}

func (j *jsonGenSPGrid) build() *GenSPGridConfig {
	return &GenSPGridConfig{
		Width:  j.Width,
		Height: j.Height,
	}
}

func parseGenSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
