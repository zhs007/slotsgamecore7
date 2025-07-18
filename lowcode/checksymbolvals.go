package lowcode

import (
	"log/slog"
	"os"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const CheckSymbolValsTypeName = "checkSymbolVals"

type CheckSymbolValsType int

const (
	CSVTypeEqu        CheckSymbolValsType = 0 // ==
	CSVTypeGreaterEqu CheckSymbolValsType = 1 // >=
	CSVTypeLessEqu    CheckSymbolValsType = 2 // <=
	CSVTypeGreater    CheckSymbolValsType = 3 // >
	CSVTypeLess       CheckSymbolValsType = 4 // <
	CSVTypeInAreaLR   CheckSymbolValsType = 5 // In [min, max]
	CSVTypeInAreaR    CheckSymbolValsType = 6 // In (min, max]
	CSVTypeInAreaL    CheckSymbolValsType = 7 // In [min, max)
	CSVTypeInArea     CheckSymbolValsType = 8 // In (min, max)
)

func parseCheckSymbolValsType(strType string) CheckSymbolValsType {
	switch strType {
	case ">=":
		return CSVTypeGreaterEqu
	case "<=":
		return CSVTypeLessEqu
	case ">":
		return CSVTypeGreater
	case "<":
		return CSVTypeLess
	case "in [min, max]":
		return CSVTypeInAreaLR
	case "in (min, max]":
		return CSVTypeInAreaR
	case "in [min, max)":
		return CSVTypeInAreaL
	case "in (min, max)":
		return CSVTypeInArea
	}

	return CSVTypeEqu
}

// CheckSymbolValsConfig - configuration for CheckSymbolVals
type CheckSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string              `yaml:"type" json:"type"`
	Type                 CheckSymbolValsType `yaml:"-" json:"-"`
	Value                int                 `yaml:"value" json:"value"`
	Min                  int                 `yaml:"min" json:"min"`
	Max                  int                 `yaml:"max" json:"max"`
	OutputToComponent    string              `yaml:"outputToComponent" json:"outputToComponent"`
}

// SetLinkComponent
func (cfg *CheckSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type CheckSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *CheckSymbolValsConfig `json:"config"`
}

// Init -
func (checkSymbolVals *CheckSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CheckSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CheckSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CheckSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return checkSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (checkSymbolVals *CheckSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	checkSymbolVals.Config = cfg.(*CheckSymbolValsConfig)
	checkSymbolVals.Config.ComponentType = CheckSymbolValsTypeName

	checkSymbolVals.Config.Type = parseCheckSymbolValsType(checkSymbolVals.Config.StrType)

	checkSymbolVals.onInit(&checkSymbolVals.Config.BasicComponentConfig)

	return nil
}

func (checkSymbolVals *CheckSymbolVals) checkVal(v int) bool {
	switch checkSymbolVals.Config.Type {
	case CSVTypeEqu:
		return v == checkSymbolVals.Config.Value
	case CSVTypeGreaterEqu:
		return v >= checkSymbolVals.Config.Value
	case CSVTypeLessEqu:
		return v <= checkSymbolVals.Config.Value
	case CSVTypeGreater:
		return v > checkSymbolVals.Config.Value
	case CSVTypeLess:
		return v < checkSymbolVals.Config.Value
	case CSVTypeInAreaLR:
		return v >= checkSymbolVals.Config.Min && v <= checkSymbolVals.Config.Max
	case CSVTypeInAreaR:
		return v > checkSymbolVals.Config.Min && v <= checkSymbolVals.Config.Max
	case CSVTypeInAreaL:
		return v >= checkSymbolVals.Config.Min && v < checkSymbolVals.Config.Max
	case CSVTypeInArea:
		return v > checkSymbolVals.Config.Min && v < checkSymbolVals.Config.Max
	}

	return false
}

// playgame
func (checkSymbolVals *CheckSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	os := checkSymbolVals.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	if os != nil {
		pc, isok := gameProp.Components.MapComponents[checkSymbolVals.Config.OutputToComponent]
		if isok {
			pccd := gameProp.GetComponentData(pc)

			for x, arr := range os.Arr {
				for y, v := range arr {
					if checkSymbolVals.checkVal(v) {
						pc.AddPos(pccd, x, y)
					}
				}
			}
		}
	}

	nc := checkSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (checkSymbolVals *CheckSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewCheckSymbolVals(name string) IComponent {
	return &CheckSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": ">",
// "outputToComponent": "bg-upmulpos",
// "value": 0
type jsonCheckSymbolVals struct {
	Type              string `json:"type"`
	Value             int    `json:"value"`
	Min               int    `json:"min"`
	Max               int    `json:"max"`
	OutputToComponent string `json:"outputToComponent"`
}

func (jcfg *jsonCheckSymbolVals) build() *CheckSymbolValsConfig {
	cfg := &CheckSymbolValsConfig{
		StrType:           strings.ToLower(jcfg.Type),
		OutputToComponent: jcfg.OutputToComponent,
		Value:             jcfg.Value,
		Min:               jcfg.Min,
		Max:               jcfg.Max,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseCheckSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCheckSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCheckSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCheckSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCheckSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CheckSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
