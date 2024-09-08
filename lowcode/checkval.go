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

const CheckValTypeName = "checkVal"

type CheckValData struct {
	BasicComponentData
	IsTrigger bool
}

// OnNewGame -
func (checkValData *CheckValData) OnNewGame(gameProp *GameProperty, component IComponent) {
	checkValData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (checkValData *CheckValData) Clone() IComponentData {
	target := &CheckValData{
		BasicComponentData: checkValData.CloneBasicComponentData(),
		IsTrigger:          checkValData.IsTrigger,
	}

	return target
}

// BuildPBComponentData
func (checkValData *CheckValData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.CheckValData{
		BasicComponentData: checkValData.BuildPBBasicComponentData(),
		IsTrigger:          checkValData.IsTrigger,
	}

	return pbcd
}

// GetValEx -
func (checkValData *CheckValData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	return 0, false
}

// CheckValConfig - configuration for CheckVal
type CheckValConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SourceVal            []string `yaml:"sourceVal" json:"sourceVal"`
	TargetVal            []string `yaml:"targetVal" json:"targetVal"`
	JumpToComponent      string   `yaml:"jumpToComponent" json:"jumpToComponent"`
	Controllers          []*Award `yaml:"controllers" json:"controllers"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *CheckValConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type CheckVal struct {
	*BasicComponent `json:"-"`
	Config          *CheckValConfig `json:"config"`
}

// Init -
func (checkVal *CheckVal) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("CheckVal.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CheckValConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("CheckVal.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return checkVal.InitEx(cfg, pool)
}

// InitEx -
func (checkVal *CheckVal) InitEx(cfg any, pool *GamePropertyPool) error {
	checkVal.Config = cfg.(*CheckValConfig)
	checkVal.Config.ComponentType = CheckValTypeName

	for _, ctrl := range checkVal.Config.Controllers {
		ctrl.Init()
	}

	checkVal.onInit(&checkVal.Config.BasicComponentConfig)

	return nil
}

// playgame
func (checkVal *CheckVal) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*CheckValData)

	sv, err := gameProp.GetComponentStrVal2(checkVal.Config.SourceVal[0], checkVal.Config.SourceVal[1])
	if err != nil {
		goutils.Error("CheckVal.OnPlayGame:GetComponentStrVal2",
			slog.Any("sourceVal", checkVal.Config.SourceVal),
			goutils.Err(err))

		return "", err
	}

	tv, err := gameProp.GetComponentStrVal2(checkVal.Config.TargetVal[0], checkVal.Config.TargetVal[1])
	if err != nil {
		goutils.Error("CheckVal.OnPlayGame:GetComponentStrVal2",
			slog.Any("targetVal", checkVal.Config.TargetVal),
			goutils.Err(err))

		return "", err
	}

	nextComponent := ""

	if sv == tv {
		cd.IsTrigger = true

		nextComponent = checkVal.Config.JumpToComponent

		if len(checkVal.Config.Controllers) > 0 {
			gameProp.procAwards(plugin, checkVal.Config.Controllers, curpr, gp)
		}
	}

	nc := checkVal.onStepEnd(gameProp, curpr, gp, nextComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (checkVal *CheckVal) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*CheckValData)

	fmt.Printf("CheckVal %v, got %v\n", checkVal.GetName(), cd.IsTrigger)

	return nil
}

// NewComponentData -
func (checkVal *CheckVal) NewComponentData() IComponentData {
	return &CheckValData{}
}

// GetAllLinkComponents - get all link components
func (checkVal *CheckVal) GetAllLinkComponents() []string {
	return []string{checkVal.Config.DefaultNextComponent, checkVal.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (checkVal *CheckVal) GetNextLinkComponents() []string {
	return []string{checkVal.Config.DefaultNextComponent, checkVal.Config.JumpToComponent}
}

func NewCheckVal(name string) IComponent {
	return &CheckVal{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "sourceVal": [
//
//	"fg-next",
//	"value"
//
// ],
// "targetVal": [
//
//	"bg-selectfg",
//	"value"
//
// ]
type jsonCheckVal struct {
	SourceVal []string `json:"sourceVal"`
	TargetVal []string `json:"targetVal"`
}

func (jcfg *jsonCheckVal) build() *CheckValConfig {
	cfg := &CheckValConfig{
		SourceVal: jcfg.SourceVal,
		TargetVal: jcfg.TargetVal,
	}

	return cfg
}

func parseCheckVal(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCheckVal:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCheckVal:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCheckVal{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCheckVal:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseCheckVal:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CheckValTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
