package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const IntValMappingTypeName = "intValMapping"

type IntValMappingData struct {
	BasicComponentData
	cfg *IntValMappingConfig
}

// OnNewGame -
func (svmd *IntValMappingData) OnNewGame(gameProp *GameProperty, component IComponent) {
	svmd.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (svmd *IntValMappingData) Clone() IComponentData {
	target := &IntValMappingData{
		BasicComponentData: svmd.CloneBasicComponentData(),
		cfg:                svmd.cfg,
	}

	return target
}

// BuildPBComponentData
func (svmd *IntValMappingData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BasicComponentData{
		BasicComponentData: svmd.BuildPBBasicComponentData(),
	}

	return pbcd
}

// ChgConfigIntVal -
func (svmd *IntValMappingData) ChgConfigIntVal(key string, off int) int {
	if key == CCVInputVal {
		if svmd.cfg.InputVal > 0 {
			svmd.MapConfigIntVals[key] = svmd.cfg.InputVal
		}
	}

	return svmd.BasicComponentData.ChgConfigIntVal(key, off)
}

func (svmd *IntValMappingData) getInputVal() int {
	input, isok := svmd.GetConfigIntVal(CCVInputVal)
	if isok {
		return input
	}

	return svmd.cfg.InputVal
}

// IntValMappingConfig - configuration for IntValMapping
type IntValMappingConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ValMapping           string                `yaml:"valMapping" json:"valMapping"`
	ValMappingVM         *sgc7game.ValMapping2 `yaml:"-" json:"-"`
	InputVal             int                   `yaml:"inputVal" json:"inputVal"`
	ComponentOutput      string                `yaml:"componentOutput" json:"componentOutput"`
	Controllers          []*Award              `yaml:"controllers" json:"controllers"` // 新的奖励系统
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &IntValMappingConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("IntValMapping.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
				slog.String("ValMapping", intValMapping.Config.ValMapping),
				goutils.Err(ErrInvalidIntValMappingFile))

			return ErrInvalidIntValMappingFile
		}

		intValMapping.Config.ValMappingVM = vm2
	} else {
		goutils.Error("IntValMapping.InitEx:ValMapping",
			goutils.Err(ErrInvalidIntValMappingFile))

		return ErrInvalidIntValMappingFile
	}

	for _, ctrl := range intValMapping.Config.Controllers {
		ctrl.Init()
	}

	intValMapping.onInit(&intValMapping.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (intValMapping *IntValMapping) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(intValMapping.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, intValMapping.Config.Controllers, curpr, gp)
	}
}

// playgame
func (intValMapping *IntValMapping) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*IntValMappingData)

	// cd.Output = 0

	in := cd.getInputVal()

	mv, isok := intValMapping.Config.ValMappingVM.MapVals[in]
	if !isok {
		goutils.Error("IntValMapping.OnPlayGame:ValMappingVM",
			goutils.Err(ErrInvalidIntValMappingValue))

		return "", ErrInvalidIntValMappingValue
	}

	cd.Output = mv.Int()

	intValMapping.ProcControllers(gameProp, plugin, curpr, gp, cd.Output, "")

	nc := intValMapping.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (intValMapping *IntValMapping) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*IntValMappingData)

	fmt.Printf("rollSymbol %v, %v => %v \n", intValMapping.GetName(), cd.getInputVal(), cd.Output)

	return nil
}

// NewComponentData -
func (intValMapping *IntValMapping) NewComponentData() IComponentData {
	return &IntValMappingData{
		cfg: intValMapping.Config,
	}
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

	return cfg
}

func parseIntValMapping(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseIntValMapping:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValMapping:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonIntValMapping{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseIntValMapping:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		controllers, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseIntValMapping:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = controllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: IntValMappingTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
