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

const GenSymbolValsTypeName = "genSymbolVals"

type GenSymbolValsType int

const (
	GSVTypeBasic  GenSymbolValsType = 0
	GSVTypeWeight GenSymbolValsType = 1
)

func parseGenSymbolValsType(strType string) GenSymbolValsType {
	if strType == "weight" {
		return GSVTypeWeight
	}

	return GSVTypeBasic
}

// GenSymbolValsConfig - configuration for GenSymbolVals
type GenSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`
	Type                 GenSymbolValsType     `yaml:"-" json:"-"`
	DefaultVal           int                   `yaml:"defaultVal" json:"defaultVal"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2 `json:"-"`
}

// SetLinkComponent
func (cfg *GenSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolValsConfig `json:"config"`
}

// Init -
func (genSymbolVals *GenSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolVals *GenSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolVals.Config = cfg.(*GenSymbolValsConfig)
	genSymbolVals.Config.ComponentType = GenSymbolValsTypeName

	genSymbolVals.Config.Type = parseGenSymbolValsType(genSymbolVals.Config.StrType)

	if genSymbolVals.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(genSymbolVals.Config.Weight, genSymbolVals.Config.UseFileMapping)
		if err != nil {
			goutils.Error("GenSymbolVals.Init:LoadStrWeights",
				slog.String("Weight", genSymbolVals.Config.Weight),
				goutils.Err(err))

			return err
		}

		genSymbolVals.Config.WeightVW = vw2
	} else {
		goutils.Error("GenSymbolVals.InitEx:Weight",
			goutils.Err(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	genSymbolVals.onInit(&genSymbolVals.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genSymbolVals *GenSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*BasicComponentData)

	if genSymbolVals.Config.Type == GSVTypeBasic {
		os := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolVals.Config.DefaultVal)

		genSymbolVals.AddOtherScene(gameProp, curpr, os, cd)
	} else if genSymbolVals.Config.Type == GSVTypeWeight {
		os := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), 0)

		for x, arr := range os.Arr {
			for y := range arr {
				ival, err := genSymbolVals.Config.WeightVW.RandVal(plugin)
				if err != nil {
					goutils.Error("GenSymbolVals.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				os.Arr[x][y] = ival.Int()

			}
		}

		genSymbolVals.AddOtherScene(gameProp, curpr, os, cd)
	}

	nc := genSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (genSymbolVals *GenSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("GenSymbolVals", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

func NewGenSymbolVals(name string) IComponent {
	return &GenSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isAlwaysGen": true,
// "type": "weight",
// "weight": "hotzone"
type jsonGenSymbolVals struct {
	Type       string `json:"type"`
	Weight     string `json:"weight"`
	DefaultVal int    `json:"defaultVal"`
}

func (jcfg *jsonGenSymbolVals) build() *GenSymbolValsConfig {
	cfg := &GenSymbolValsConfig{
		StrType:    jcfg.Type,
		Weight:     jcfg.Weight,
		DefaultVal: jcfg.DefaultVal,
	}

	return cfg
}

func parseGenSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
