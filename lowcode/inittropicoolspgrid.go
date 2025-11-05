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

const InitTropiCoolSPGridTypeName = "initTropiCoolSPGrid"

// InitTropiCoolSPGridConfig - configuration for InitTropiCoolSPGrid
type InitTropiCoolSPGridConfig struct {
    BasicComponentConfig `yaml:",inline" json:",inline"`
    Width  int `yaml:"width" json:"width"`
    Height int `yaml:"height" json:"height"`
}

// SetLinkComponent
func (cfg *InitTropiCoolSPGridConfig) SetLinkComponent(link string, componentName string) {
    if link == "next" {
        cfg.DefaultNextComponent = componentName
    }
}

type InitTropiCoolSPGrid struct {
    *BasicComponent `json:"-"`
    Config          *InitTropiCoolSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *InitTropiCoolSPGrid) Init(fn string, pool *GamePropertyPool) error {
    data, err := os.ReadFile(fn)
    if err != nil {
        goutils.Error("InitTropiCoolSPGrid.Init:ReadFile",
            slog.String("fn", fn),
            goutils.Err(err))

        return err
    }

    cfg := &InitTropiCoolSPGridConfig{}

    err = yaml.Unmarshal(data, cfg)
    if err != nil {
        goutils.Error("InitTropiCoolSPGrid.Init:Unmarshal",
            slog.String("fn", fn),
            goutils.Err(err))

        return err
    }

    return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *InitTropiCoolSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
    gen.Config = cfg.(*InitTropiCoolSPGridConfig)
    gen.Config.ComponentType = InitTropiCoolSPGridTypeName

    gen.onInit(&gen.Config.BasicComponentConfig)

    return nil
}

// OnPlayGame - minimal implementation: does nothing but advance
func (gen *InitTropiCoolSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
    cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

    // Initialize an SPGrid for this component on the current play result
    curpr.SPGrid[gen.GetName()] = &sgc7game.SPGrid{
        Width:  gen.Config.Width,
        Height: gen.Config.Height,
    }

    nc := gen.onStepEnd(gameProp, curpr, gp, "")

    return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *InitTropiCoolSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
    return nil
}

func NewInitTropiCoolSPGrid(name string) IComponent {
    return &InitTropiCoolSPGrid{
        BasicComponent: NewBasicComponent(name, 1),
    }
}

// json representation used by editor
type jsonInitTropiCoolSPGrid struct {
    Width  int `json:"Width"`
    Height int `json:"Height"`
}

func (j *jsonInitTropiCoolSPGrid) build() *InitTropiCoolSPGridConfig {
    return &InitTropiCoolSPGridConfig{
        Width:  j.Width,
        Height: j.Height,
    }
}

func parseInitTropiCoolSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
    cfg, label, _, err := getConfigInCell(cell)
    if err != nil {
        goutils.Error("parseInitTropiCoolSPGrid:getConfigInCell",
            goutils.Err(err))

        return "", err
    }

    buf, err := cfg.MarshalJSON()
    if err != nil {
        goutils.Error("parseInitTropiCoolSPGrid:MarshalJSON",
            goutils.Err(err))

        return "", err
    }

    data := &jsonInitTropiCoolSPGrid{}

    err = sonic.Unmarshal(buf, data)
    if err != nil {
        goutils.Error("parseInitTropiCoolSPGrid:Unmarshal",
            goutils.Err(err))

        return "", err
    }

    cfgd := data.build()

    gamecfg.mapConfig[label] = cfgd
    gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

    ccfg := &ComponentConfig{
        Name: label,
        Type: InitTropiCoolSPGridTypeName,
    }

    gamecfg.Components = append(gamecfg.Components, ccfg)

    return label, nil
}
