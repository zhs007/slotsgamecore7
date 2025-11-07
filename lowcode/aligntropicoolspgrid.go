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

const AlignTropiCoolSPGridTypeName = "alignTropiCoolSPGrid"

// AlignTropiCoolSPGridConfig - configuration for AlignTropiCoolSPGrid
type AlignTropiCoolSPGridConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SPGrid               string `yaml:"spGrid" json:"spGrid"`
	BlankSymbol          string `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int    `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *AlignTropiCoolSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type AlignTropiCoolSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *AlignTropiCoolSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *AlignTropiCoolSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("AlignTropiCoolSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &AlignTropiCoolSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("AlignTropiCoolSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *AlignTropiCoolSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*AlignTropiCoolSPGridConfig)
	gen.Config.ComponentType = AlignTropiCoolSPGridTypeName

	if gen.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.BlankSymbol]
		if !isok {
			goutils.Error("AlignTropiCoolSPGrid.InitEx:BlankSymbol",
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

// OnPlayGame - minimal implementation: does nothing but advance
func (gen *AlignTropiCoolSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	bcd := icd.(*BasicComponentData)

	// Reuse NewSPGrid generation, but if Align is true, attempt a simple post-processing
	// step to align special symbols to the top-left (a lightweight approximation of
	// TropiCool alignment). This keeps the change low-risk while providing the
	// requested distinct component behavior.
	stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.SPGrid]
	if !isok {
		goutils.Error("AlignTropiCoolSPGrid.OnPlayGame:MapSPGridStack",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	spgrid := stackSPGrid.Stack.GetTopSPGridEx(gen.Config.SPGrid, curpr, prs)
	if spgrid == nil {
		goutils.Error("AlignTropiCoolSPGrid.OnPlayGame:GetTopSPGridEx",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	newspgrid := spgrid

	ismoved := false
	for x := 0; x < newspgrid.Width-1; {
		isnone := true
		for y := 0; y < newspgrid.Height; y++ {
			if newspgrid.Arr[x][y] != -1 {
				isnone = false

				break
			}
		}

		if isnone {
			if newspgrid == spgrid {
				newspgrid = spgrid.CloneEx(gameProp.PoolScene)
			}

			for tx := x; tx < newspgrid.Width-1; tx++ {
				for y := 0; y < newspgrid.Height; y++ {
					newspgrid.Arr[tx][y] = newspgrid.Arr[tx+1][y]
				}
			}

			for y := 0; y < newspgrid.Height; y++ {
				newspgrid.Arr[newspgrid.Width-1][y] = -1
			}

			ismoved = true
		} else {
			x++
		}
	}

	if ismoved {
		gen.AddSPGrid(gen.Config.SPGrid, gameProp, curpr, newspgrid, bcd)
	}

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *AlignTropiCoolSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewAlignTropiCoolSPGrid(name string) IComponent {
	return &AlignTropiCoolSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "spGrid": "bg-spgrid",
// "BlankSymbol": "BN"
type jsonAlignTropiCoolSPGrid struct {
	SPGrid      string `json:"spGrid"`
	BlankSymbol string `json:"BlankSymbol"`
}

func (j *jsonAlignTropiCoolSPGrid) build() *AlignTropiCoolSPGridConfig {
	return &AlignTropiCoolSPGridConfig{
		SPGrid:      j.SPGrid,
		BlankSymbol: j.BlankSymbol,
	}
}

func parseAlignTropiCoolSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseAlignTropiCoolSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseAlignTropiCoolSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonAlignTropiCoolSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseAlignTropiCoolSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: AlignTropiCoolSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
