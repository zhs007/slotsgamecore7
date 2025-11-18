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

const CalculatorTypeName = "calculator"

// CalculatorConfig - configuration for Calculator
type CalculatorConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Input1               int                 `yaml:"input1" json:"input1"`
	Input2               int                 `yaml:"input2" json:"input2"`
	Formula              string              `yaml:"formula" json:"formula"`
	calculator           *CalculatorCore     `yaml:"-" json:"-"`
	MapControllers       map[string][]*Award `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent - set link, only "next" supported
func (cfg *CalculatorConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Calculator struct {
	*BasicComponent `json:"-"`
	Config          *CalculatorConfig `json:"config"`
}

// Init - load from yaml file
func (c *Calculator) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Calculator.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CalculatorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Calculator.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return c.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (c *Calculator) InitEx(cfg any, pool *GamePropertyPool) error {
	c.Config = cfg.(*CalculatorConfig)
	c.Config.ComponentType = CalculatorTypeName

	cc, err := NewCalculatorCore(c.Config.Formula)
	if err != nil {
		goutils.Error("Calculator.InitEx:NewCalculatorCore",
			slog.String("formula", c.Config.Formula),
			goutils.Err(err))

		return err
	}

	c.Config.calculator = cc

	for _, ctrls := range c.Config.MapControllers {
		for _, ctrl := range ctrls {
			ctrl.Init()
		}
	}

	c.onInit(&c.Config.BasicComponentConfig)

	return nil
}

// ProcControllers -
func (c *Calculator) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	arr, isok := c.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, arr, curpr, gp)
	}
}

func (c *Calculator) getInput1(_ *GameProperty, basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVInput1)
	if isok {
		return v
	}

	return c.Config.Input1
}

func (c *Calculator) getInput2(_ *GameProperty, basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVInput2)
	if isok {
		return v
	}

	return c.Config.Input2
}

// OnPlayGame - placeholder: does nothing and move on to next
func (c *Calculator) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	bcd := icd.(*BasicComponentData)

	inputs := []int{
		c.getInput1(gameProp, bcd),
		c.getInput2(gameProp, bcd),
	}

	ret, err := c.Config.calculator.CalcVal(inputs)
	if err != nil {
		goutils.Error("Calculator.OnPlayGame:CalcVal",
			goutils.Err(err))

		return "", err
	}

	bcd.Output = ret

	c.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

	// move to next component
	nc := c.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - placeholder ascii representation
func (c *Calculator) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewCalculator(name string) IComponent {
	return &Calculator{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

// "input1": 0,
// "input2": 0,
// "formula": "input1 - 1"
type jsonCalculator struct {
	Input1  int    `json:"input1"`
	Input2  int    `json:"input2"`
	Formula string `json:"formula"`
}

func (j *jsonCalculator) build() *CalculatorConfig {
	return &CalculatorConfig{
		Input1:  j.Input1,
		Input2:  j.Input2,
		Formula: strings.ToLower(j.Formula),
	}
}

func parseCalculator(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCalculator:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCalculator:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCalculator{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCalculator:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapControllers, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseCalculator:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapControllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CalculatorTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
