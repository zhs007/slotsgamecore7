package lowcode

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const IntValMappingTypeName = "intValMapping"

// IntValMappingConfig - configuration for IntValMapping
type IntValMappingConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ValMapping           string                `yaml:"valMapping" json:"valMapping"`
	ValMappingVM         *sgc7game.ValMapping2 `yaml:"-" json:"-"`
	InputVal             int                   `yaml:"inputVal" json:"inputVal"`
	ComponentOutput      string                `yaml:"componentOutput" json:"componentOutput"`
}

// SetLinkComponent
func (cfg *IntValMappingConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type IntValMapping struct {
	*BasicComponent `json:"-"`
	Config          *IntValMappingConfig `json:"config"`
}

// Init -
func (intValMapping *IntValMapping) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("IntValMapping.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &IntValMappingConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("IntValMapping.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return intValMapping.InitEx(cfg, pool)
}

// InitEx -
func (intValMapping *IntValMapping) InitEx(cfg any, pool *GamePropertyPool) error {
	intValMapping.Config = cfg.(*IntValMappingConfig)
	intValMapping.Config.ComponentType = IntValMappingTypeName

	if intValMapping.Config.ValMapping != "" {
		vm2 := pool.LoadIntMapping(intValMapping.Config.ValMapping)
		if vm2 == nil {
			goutils.Error("IntValMapping.Init:LoadIntMapping",
				zap.String("ValMapping", intValMapping.Config.ValMapping),
				zap.Error(ErrInvalidIntValMappingFile))

			return ErrInvalidIntValMappingFile
		}

		intValMapping.Config.ValMappingVM = vm2
	} else {
		goutils.Error("IntValMapping.InitEx:ValMapping",
			zap.Error(ErrInvalidIntValMappingFile))

		return ErrInvalidIntValMappingFile
	}

	intValMapping.onInit(&intValMapping.Config.BasicComponentConfig)

	return nil
}

func (intValMapping *IntValMapping) getInput(basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVInputVal)
	if isok {
		return v
	}

	return intValMapping.Config.InputVal
}

// playgame
func (intValMapping *IntValMapping) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*BasicComponentData)

	// cd.Output = 0

	in := intValMapping.getInput(cd)

	mv, isok := intValMapping.Config.ValMappingVM.MapVals[in]
	if !isok {
		goutils.Error("IntValMapping.OnPlayGame:ValMappingVM",
			zap.Error(ErrInvalidIntValMappingValue))

		return "", ErrInvalidIntValMappingValue
	}

	cd.Output = mv.Int()

	nc := intValMapping.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (intValMapping *IntValMapping) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	fmt.Printf("rollSymbol %v, %v => %v \n", intValMapping.GetName(), intValMapping.getInput(cd), cd.Output)

	return nil
}

// OnStats
func (intValMapping *IntValMapping) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewIntValMapping(name string) IComponent {
	return &IntValMapping{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"valMapping": "valmapping",
//		"inputVal": 3
//	},
type jsonIntValMapping struct {
	ValMapping string `json:"valMapping"`
	InputVal   int    `json:"inputVal"`
}

func (jcfg *jsonIntValMapping) build() *IntValMappingConfig {
	cfg := &IntValMappingConfig{
		ValMapping: jcfg.ValMapping,
		InputVal:   jcfg.InputVal,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseIntValMapping(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseIntValMapping:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValMapping:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonIntValMapping{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseIntValMapping:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: IntValMappingTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
