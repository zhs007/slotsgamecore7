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

const DropDownTropiCoolSPGridTypeName = "dropDownTropiCoolSPGrid"

// DropDownTropiCoolSPGridConfig - configuration for DropDownTropiCoolSPGrid
type DropDownTropiCoolSPGridConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SPGrid               string `yaml:"spGrid" json:"spGrid"`
	BlankSymbol          string `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int    `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *DropDownTropiCoolSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type DropDownTropiCoolSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *DropDownTropiCoolSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *DropDownTropiCoolSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("DropDownTropiCoolSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &DropDownTropiCoolSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("DropDownTropiCoolSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *DropDownTropiCoolSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*DropDownTropiCoolSPGridConfig)
	gen.Config.ComponentType = DropDownTropiCoolSPGridTypeName

	if gen.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.BlankSymbol]
		if !isok {
			goutils.Error("DropDownTropiCoolSPGrid.InitEx:BlankSymbol",
				slog.String("BlankSymbol", gen.Config.BlankSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.BlankSymbolCode = sc
	} else {
		gen.Config.BlankSymbolCode = -1
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *DropDownTropiCoolSPGrid) getSPGridSymbol(spgrid *sgc7game.GameScene, x int) int {
	if spgrid.Arr[x][spgrid.Height-1] == -1 {
		return -1
	}

	sym := spgrid.Arr[x][spgrid.Height-1]

	for y := spgrid.Height - 1; y > 0; y-- {
		spgrid.Arr[x][y] = spgrid.Arr[x][y-1]
	}

	spgrid.Arr[x][0] = -1

	return sym
}

// OnPlayGame - minimal implementation: does nothing but advance
func (gen *DropDownTropiCoolSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// This implementation intentionally does not modify the play result.
	// It simply ends this step and returns to the next component. It can be
	// extended later to implement drop-down / TropiCool-specific behaviour.
	bcd := icd.(*BasicComponentData)

	stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.SPGrid]
	if !isok {
		goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:MapSPGridStack",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	spgrid := stackSPGrid.Stack.GetTopSPGridEx(gen.Config.SPGrid, curpr, prs)
	if spgrid == nil {
		goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:GetTopSPGridEx",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	newspgrid := spgrid.CloneEx(gameProp.PoolScene)

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	ngs := gs.CloneEx(gameProp.PoolScene)

	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				sym := gen.getSPGridSymbol(newspgrid, x)
				if sym == -1 {
					break
				}

				ngs.Arr[x][y] = sym
			}
		}
	}

	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == gen.Config.BlankSymbolCode {
				ngs.Arr[x][y] = -1
			}
		}
	}

	gen.AddScene(gameProp, curpr, ngs, bcd)
	gen.AddSPGrid(gen.Config.SPGrid, gameProp, curpr, newspgrid, bcd)

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *DropDownTropiCoolSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewDropDownTropiCoolSPGrid(name string) IComponent {
	return &DropDownTropiCoolSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "spGrid": "bg-spgrid",
// "BlankSymbol": "BN"
type jsonDropDownTropiCoolSPGrid struct {
	SPGrid      string `json:"spGrid"`
	BlankSymbol string `json:"BlankSymbol"`
}

func (j *jsonDropDownTropiCoolSPGrid) build() *DropDownTropiCoolSPGridConfig {
	return &DropDownTropiCoolSPGridConfig{
		SPGrid:      j.SPGrid,
		BlankSymbol: j.BlankSymbol,
	}
}

func parseDropDownTropiCoolSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseDropDownTropiCoolSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseDropDownTropiCoolSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonDropDownTropiCoolSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseDropDownTropiCoolSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: DropDownTropiCoolSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
