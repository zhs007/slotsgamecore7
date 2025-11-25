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

const MergeSPGridTypeName = "mergeSPGrid"

// MergeSPGridConfig is a placeholder configuration for MergeSPGrid
type MergeSPGridConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Source1              string              `yaml:"source1" json:"source1"`
	Source2              string              `yaml:"source2" json:"source2"`
	Formula              string              `yaml:"formula" json:"formula"`
	Output               string              `yaml:"output" json:"output"`
	MapControllers       map[string][]*Award `yaml:"controllers" json:"controllers"`
	core                 *MergeSPGridCore    `yaml:"-" json:"-"`
}

func (cfg *MergeSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type MergeSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *MergeSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *MergeSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MergeSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MergeSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MergeSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - init from config object
func (gen *MergeSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*MergeSPGridConfig)
	gen.Config.ComponentType = MergeSPGridTypeName

	mc, err := NewMergeSPGridCore(gen.Config.Formula)
	if err != nil {
		goutils.Error("MergeSPGrid.InitEx:NewMergeSPGridCore",
			goutils.Err(err))

		return err
	}

	gen.Config.core = mc

	for _, ctrls := range gen.Config.MapControllers {
		for _, ctrl := range ctrls {
			ctrl.Init()
		}
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

// ProcControllers -
func (gen *MergeSPGrid) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	arr, isok := gen.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, arr, curpr, gp)
	}
}

func (gen *MergeSPGrid) getSource1(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) (*sgc7game.GameScene, error) {
	if gen.Config.Source1 != "" {
		stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.Source1]
		if !isok {
			goutils.Error("MergeSPGrid.getSource1:MapSPGridStack",
				slog.String("SPGrid", gen.Config.Source1),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		spgrid := stackSPGrid.Stack.GetTopSPGridEx(gen.Config.Source1, curpr, prs)
		if spgrid == nil {
			goutils.Error("MergeSPGrid.getSource1:GetTopSPGridEx",
				slog.String("SPGrid", gen.Config.Source1),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		return spgrid, nil
	}

	os := gen.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	return os, nil
}

func (gen *MergeSPGrid) getSource2(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult) (*sgc7game.GameScene, error) {
	if gen.Config.Source2 != "" {
		stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.Source2]
		if !isok {
			goutils.Error("MergeSPGrid.getSource2:MapSPGridStack",
				slog.String("SPGrid", gen.Config.Source2),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		spgrid := stackSPGrid.Stack.GetTopSPGridEx(gen.Config.Source2, curpr, prs)
		if spgrid == nil {
			goutils.Error("MergeSPGrid.getSource2:GetTopSPGridEx",
				slog.String("SPGrid", gen.Config.Source2),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		return spgrid, nil
	}

	os := gen.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	return os, nil
}

func (gen *MergeSPGrid) getOutput(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, bcd *BasicComponentData) (*sgc7game.GameScene, error) {
	if gen.Config.Output != "" {
		stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.Output]
		if !isok {
			goutils.Error("MergeSPGrid.getOutput:MapSPGridStack",
				slog.String("SPGrid", gen.Config.Output),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		gs := gameProp.PoolScene.New2(stackSPGrid.Width, stackSPGrid.Height, 0)
		if gs == nil {
			goutils.Error("InitTropiCoolSPGrid.OnPlayGame:New2",
				slog.Int("Width", stackSPGrid.Width),
				slog.Int("Height", stackSPGrid.Height),
				slog.Int("EmptySymbolCode", 0),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, ErrInvalidComponentConfig
		}

		newspgrid := gs

		gen.AddSPGrid(gen.Config.Output, gameProp, curpr, newspgrid, bcd)

		return newspgrid, nil
	}

	os := gen.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	nos := os.CloneEx(gameProp.PoolScene)
	gen.AddOtherScene(gameProp, curpr, nos, bcd)

	return nos, nil
}

// OnPlayGame - placeholder: does nothing but advance
func (gen *MergeSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// placeholder implementation, does not change play result
	bcd := icd.(*BasicComponentData)
	source1, err := gen.getSource1(gameProp, curpr, prs)
	if err != nil {
		goutils.Error("MergeSPGrid.OnPlayGame:getSource1",
			goutils.Err(err))

		return "", err
	}

	source2, err := gen.getSource2(gameProp, curpr, prs)
	if err != nil {
		goutils.Error("MergeSPGrid.OnPlayGame:getSource2",
			goutils.Err(err))

		return "", err
	}

	output, err := gen.getOutput(gameProp, curpr, prs, bcd)
	if err != nil {
		goutils.Error("MergeSPGrid.OnPlayGame:getOutput",
			goutils.Err(err))

		return "", err
	}

	if source1.Width != source2.Width || source1.Height != source2.Height || source1.Width != output.Width || source1.Height != output.Height {
		goutils.Error("MergeSPGrid.OnPlayGame:invalid dimensions",
			slog.Int("source1Width", source1.Width),
			slog.Int("source1Height", source1.Height),
			slog.Int("source2Width", source2.Width),
			slog.Int("source2Height", source2.Height),
			slog.Int("outputWidth", output.Width),
			slog.Int("outputHeight", output.Height),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	for x := 0; x < source1.Width; x++ {
		for y := 0; y < source1.Height; y++ {
			s1 := source1.Arr[x][y]
			s2 := source2.Arr[x][y]

			val, err := gen.Config.core.CalcVal([]int{s1, s2})
			if err != nil {
				goutils.Error("MergeSPGrid.OnPlayGame:CalcVal",
					slog.Int("source1", s1),
					slog.Int("source2", s2),
					goutils.Err(err))

				return "", err
			}

			output.Arr[x][y] = val
		}
	}

	gen.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *MergeSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewMergeSPGrid(name string) IComponent {
	return &MergeSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "source2": "cg-spgrid-multi",
// "formula": "(source1 > 0 && source2 > 0) ? source1 * source2: 0"
type jsonMergeSPGrid struct {
	Source1 string `json:"source1"`
	Source2 string `json:"source2"`
	Formula string `json:"formula"`
	Output  string `json:"output"`
}

func (jcfg *jsonMergeSPGrid) build() *MergeSPGridConfig {
	return &MergeSPGridConfig{
		Source1: jcfg.Source1,
		Source2: jcfg.Source2,
		Formula: jcfg.Formula,
		Output:  jcfg.Output,
	}
}

func parseMergeSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseMergeSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseMergeSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonMergeSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseMergeSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: MergeSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
