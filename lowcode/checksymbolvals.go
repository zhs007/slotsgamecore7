package lowcode

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
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
	if strType == ">=" {
		return CSVTypeGreaterEqu
	} else if strType == "<=" {
		return CSVTypeLessEqu
	} else if strType == ">" {
		return CSVTypeGreater
	} else if strType == "<" {
		return CSVTypeLess
	} else if strType == "In [min, max]" {
		return CSVTypeInAreaLR
	} else if strType == "In (min, max]" {
		return CSVTypeInAreaR
	} else if strType == "In [min, max)" {
		return CSVTypeInAreaL
	} else if strType == "In (min, max)" {
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &CheckSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CheckSymbolVals.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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
	if checkSymbolVals.Config.Type == CSVTypeEqu {
		return v == checkSymbolVals.Config.Value
	} else if checkSymbolVals.Config.Type == CSVTypeGreaterEqu {
		return v >= checkSymbolVals.Config.Value
	} else if checkSymbolVals.Config.Type == CSVTypeLessEqu {
		return v <= checkSymbolVals.Config.Value
	} else if checkSymbolVals.Config.Type == CSVTypeGreater {
		return v > checkSymbolVals.Config.Value
	} else if checkSymbolVals.Config.Type == CSVTypeLess {
		return v < checkSymbolVals.Config.Value
	} else if checkSymbolVals.Config.Type == CSVTypeInAreaLR {
		return v >= checkSymbolVals.Config.Min && v <= checkSymbolVals.Config.Max
	} else if checkSymbolVals.Config.Type == CSVTypeInAreaR {
		return v > checkSymbolVals.Config.Min && v <= checkSymbolVals.Config.Max
	} else if checkSymbolVals.Config.Type == CSVTypeInAreaL {
		return v >= checkSymbolVals.Config.Min && v < checkSymbolVals.Config.Max
	} else if checkSymbolVals.Config.Type == CSVTypeInArea {
		return v > checkSymbolVals.Config.Min && v < checkSymbolVals.Config.Max
	}

	return false
}

// playgame
func (checkSymbolVals *CheckSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	// cd := icd.(*BasicComponentData)

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

	// cd := icd.(*BasicComponentData)

	// if len(cd.UsedOtherScenes) > 0 {
	// 	asciigame.OutputOtherScene("The value of the symbols", pr.OtherScenes[cd.UsedOtherScenes[0]])
	// }

	return nil
}

// // OnStats
// func (checkSymbolVals *CheckSymbolVals) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

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
		StrType:           jcfg.Type,
		OutputToComponent: jcfg.OutputToComponent,
		Value:             jcfg.Value,
		Min:               jcfg.Min,
		Max:               jcfg.Max,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseCheckSymbolVals(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCheckSymbolVals:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCheckSymbolVals:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonCheckSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCheckSymbolVals:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CheckSymbolValsTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
