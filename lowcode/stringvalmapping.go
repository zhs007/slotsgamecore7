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

const StringValMappingTypeName = "stringValMapping"

type StringValMappingData struct {
	BasicComponentData
	cfg *StringValMappingConfig
}

// OnNewGame -
func (svmd *StringValMappingData) OnNewGame(gameProp *GameProperty, component IComponent) {
	svmd.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (svmd *StringValMappingData) Clone() IComponentData {
	target := &StringValMappingData{
		BasicComponentData: svmd.CloneBasicComponentData(),
		cfg:                svmd.cfg,
	}

	return target
}

// BuildPBComponentData
func (svmd *StringValMappingData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.BasicComponentData{
		BasicComponentData: svmd.BuildPBBasicComponentData(),
	}

	return pbcd
}

// ChgConfigIntVal -
func (svmd *StringValMappingData) ChgConfigIntVal(key string, off int) int {
	if key == CCVInputVal {
		if svmd.cfg.InputVal > 0 {
			svmd.MapConfigIntVals[key] = svmd.cfg.InputVal
		}
	}

	return svmd.BasicComponentData.ChgConfigIntVal(key, off)
}

func (svmd *StringValMappingData) getInputVal() int {
	input, isok := svmd.GetConfigIntVal(CCVInputVal)
	if isok {
		return input
	}

	return svmd.cfg.InputVal
}

// StringValMappingConfig placeholder for mapping int -> string
type StringValMappingConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ValMapping           string                `yaml:"valMapping" json:"valMapping"`
	ValMappingVM         *sgc7game.ValMapping2 `yaml:"-" json:"-"`
	InputVal             int                   `yaml:"inputVal" json:"inputVal"`
	ComponentOutput      string                `yaml:"componentOutput" json:"componentOutput"`
	Controllers          []*Award              `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent - set only "next" support
func (cfg *StringValMappingConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type StringValMapping struct {
	*BasicComponent `json:"-"`
	Config          *StringValMappingConfig `json:"config"`
}

// Init - read yaml configuration file
func (svm *StringValMapping) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("StringValMapping.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &StringValMappingConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("StringValMapping.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return svm.InitEx(cfg, pool)
}

// InitEx - init with parsed config
func (svm *StringValMapping) InitEx(cfg any, pool *GamePropertyPool) error {
	svm.Config = cfg.(*StringValMappingConfig)
	svm.Config.ComponentType = StringValMappingTypeName

	if svm.Config.ValMapping != "" {
		vm2 := pool.LoadIntMapping(svm.Config.ValMapping)
		if vm2 == nil {
			goutils.Error("StringValMapping.Init:LoadIntMapping",
				slog.String("ValMapping", svm.Config.ValMapping),
				goutils.Err(ErrInvalidStringValMappingFile))

			return ErrInvalidStringValMappingFile
		}

		svm.Config.ValMappingVM = vm2
	} else {
		goutils.Error("StringValMapping.InitEx:ValMapping",
			goutils.Err(ErrInvalidStringValMappingFile))

		return ErrInvalidStringValMappingFile
	}

	for _, ctrl := range svm.Config.Controllers {
		ctrl.Init()
	}

	svm.onInit(&svm.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (svm *StringValMapping) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(svm.Config.Controllers) > 0 {
		gameProp.procAwards(plugin, svm.Config.Controllers, curpr, gp)
	}
}

// OnPlayGame - placeholder: maps input int -> string and set StrOutput
func (svm *StringValMapping) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*StringValMappingData)

	in := cd.getInputVal()

	mv, isok := svm.Config.ValMappingVM.MapVals[in]
	if !isok {
		goutils.Error("StringValMapping.OnPlayGame:ValMappingVM",
			goutils.Err(ErrInvalidStringValMappingValue))

		return "", ErrInvalidStringValMappingValue
	}

	cd.StrOutput = mv.String()

	svm.ProcControllers(gameProp, plugin, curpr, gp, -1, cd.StrOutput)

	nc := svm.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - show info
func (svm *StringValMapping) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*StringValMappingData)

	fmt.Printf("StringValMapping %v, %v => %v\n", svm.GetName(), cd.getInputVal(), cd.StrOutput)

	return nil
}

// NewComponentData -
func (svm *StringValMapping) NewComponentData() IComponentData {
	return &StringValMappingData{
		cfg: svm.Config,
	}
}

func NewStringValMapping(name string) IComponent {
	return &StringValMapping{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "inputValType": "number",
// "inputVal": 4,
// "valMapping": "heightmapping"
type jsonStringValMapping struct {
	ValMapping string `json:"valMapping"`
	InputVal   int    `json:"inputVal"`
}

func (jcfg *jsonStringValMapping) build() *StringValMappingConfig {
	cfg := &StringValMappingConfig{
		ValMapping: jcfg.ValMapping,
		InputVal:   jcfg.InputVal,
	}

	return cfg
}

func parseStringValMapping(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseStringValMapping:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseStringValMapping:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonStringValMapping{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseStringValMapping:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		controllers, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseStringValMapping:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = controllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: StringValMappingTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
