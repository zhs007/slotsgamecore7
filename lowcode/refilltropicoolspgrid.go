package lowcode

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const RefillTropiCoolSPGridTypeName = "refillTropiCoolSPGrid"

// RefillTropiCoolSPGridConfig - configuration for RefillTropiCoolSPGrid
type RefillTropiCoolSPGridConfig struct {
	BasicComponentConfig  `yaml:",inline" json:",inline"`
	MaxNumber             int                   `yaml:"maxNumber" json:"maxNumber"`
	SPGrid                string                `yaml:"spGrid" json:"spGrid"`
	BlankSymbol           string                `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode       int                   `yaml:"-" json:"-"`
	GigaSymbols           []string              `yaml:"gigSymbols" json:"gigSymbols"`
	GigaSymbolCodes       []int                 `yaml:"-" json:"-"`
	TargetGigaSymbolCodes map[int](map[int]int) `yaml:"-" json:"-"` // key: symbolCode, value: size->symbolCode
	SPSymbols             []string              `yaml:"spSymbols" json:"spSymbols"`
	SPSymbolCodes         []int                 `yaml:"-" json:"-"`
	Weight                string                `yaml:"weight" json:"weight"`
	WeightVM              *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	GigaWeight            string                `yaml:"gigaWeight" json:"gigaWeight"`
	GigaWeightVM          *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	EmptySymbol           string                `yaml:"emptySymbol" json:"emptySymbol"`
	EmptySymbolCode       int                   `yaml:"-" json:"-"`
	MapControls           map[string][]*Award   `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *RefillTropiCoolSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RefillTropiCoolSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *RefillTropiCoolSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *RefillTropiCoolSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RefillTropiCoolSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RefillTropiCoolSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillTropiCoolSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *RefillTropiCoolSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*RefillTropiCoolSPGridConfig)
	gen.Config.ComponentType = RefillTropiCoolSPGridTypeName

	// process symbol codes
	if gen.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.BlankSymbol]
		if !isok {
			goutils.Error("RefillTropiCoolSPGrid.InitEx:BlankSymbol",
				slog.String("BlankSymbol", gen.Config.BlankSymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.BlankSymbolCode = sc
	} else {
		gen.Config.BlankSymbolCode = -1
	}

	for _, v := range gen.Config.GigaSymbols {
		code, isok := pool.Config.GetDefaultPaytables().MapSymbols[v]
		if !isok {
			goutils.Error("RefillTropiCoolSPGrid.InitEx:GigaSymbols",
				slog.String("GigaSymbols", v),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.GigaSymbolCodes = append(gen.Config.GigaSymbolCodes, code)
	}

	if gen.Config.EmptySymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.EmptySymbol]
		if !isok {
			goutils.Error("RefillTropiCoolSPGrid.InitEx:EmptySymbol",
				slog.String("EmptySymbol", gen.Config.EmptySymbol),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.EmptySymbolCode = sc
	} else {
		gen.Config.EmptySymbolCode = -1
	}

	for _, v := range gen.Config.SPSymbols {
		code, isok := pool.Config.GetDefaultPaytables().MapSymbols[v]
		if !isok {
			goutils.Error("RefillTropiCoolSPGrid.InitEx:GigaSymbols",
				slog.String("GigaSymbols", v),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.SPSymbolCodes = append(gen.Config.SPSymbolCodes, code)
	}

	// weights
	if gen.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(gen.Config.Weight, true)
		if err != nil {
			goutils.Error("RefillTropiCoolSPGrid.Init:LoadStrWeights",
				slog.String("Weight", gen.Config.Weight),
				goutils.Err(err))

			return err
		}

		gen.Config.WeightVM = vw2
	}

	if gen.Config.GigaWeight != "" {
		vw2, err := pool.LoadIntWeights(gen.Config.GigaWeight, true)
		if err != nil {
			goutils.Error("RefillTropiCoolSPGrid.Init:LoadStrWeights",
				slog.String("Weight", gen.Config.GigaWeight),
				goutils.Err(err))

			return err
		}

		gen.Config.GigaWeightVM = vw2

		gen.Config.TargetGigaSymbolCodes = make(map[int](map[int]int))

		for _, v := range vw2.Vals {
			s := v.Int()

			if s == 0 {
				continue
			}

			for i, sc := range gen.Config.GigaSymbolCodes {
				if gen.Config.TargetGigaSymbolCodes[sc] == nil {
					gen.Config.TargetGigaSymbolCodes[sc] = make(map[int]int)
				}

				cs := fmt.Sprintf("%v_%v", gen.Config.GigaSymbols[i], s)
				csc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
				if isok {
					gen.Config.TargetGigaSymbolCodes[sc][s] = csc
				}
			}
		}
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *RefillTropiCoolSPGrid) getWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("RefillTropiCoolSPGrid.getWeight:LoadIntWeights",
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return gen.Config.WeightVM
}

func (gen *RefillTropiCoolSPGrid) getGigaWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVGigaWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("RefillTropiCoolSPGrid.getGigaWeight:LoadIntWeights",
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return gen.Config.GigaWeightVM
}

func (gen *RefillTropiCoolSPGrid) setGiga(gs *sgc7game.GameScene, x int, y int, s int, c int) {
	if x+s > gs.Width || y-s+1 < 0 {
		return
	}

	for ix := x; ix < x+s; ix++ {
		for iy := y; iy > y-s; iy-- {
			if goutils.IndexOfIntSlice(gen.Config.SPSymbolCodes, gs.Arr[ix][iy], 0) >= 0 {

				return
			}
		}
	}

	for ix := x; ix < x+s; ix++ {
		for iy := y; iy > y-s; iy-- {
			gs.Arr[ix][iy] = c
		}
	}
}

// OnPlayGame - generate SPGrid and optionally apply refill (gravity to bottom)
func (gen *RefillTropiCoolSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	bcd := icd.(*BasicComponentData)

	stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.SPGrid]
	if !isok {
		goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:MapSPGridStack",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	spgrid := stackSPGrid.Stack.GetTopSPGridEx(gen.Config.SPGrid, curpr, prs)
	if spgrid == nil {
		goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:GetTopSPGridEx",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	startX := -1
	for x := 0; x < spgrid.Width; x++ {
		isnone := true
		for y := 0; y < spgrid.Height; y++ {
			if spgrid.Arr[x][y] != -1 {
				isnone = false

				break
			}
		}

		if isnone {
			startX = x

			break
		}
	}

	if startX == -1 {
		// nothing to refill
		nc := gen.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	newspgrid := spgrid.CloneEx(gameProp.PoolScene)

	vw := gen.getWeight(gameProp, bcd)
	if vw == nil {
		goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:getWeight",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	for x := startX; x < newspgrid.Width; x++ {
		arr := newspgrid.Arr[x]

		for y := range arr {
			cv, err := vw.RandVal(plugin)
			if err != nil {
				goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
			}

			newspgrid.Arr[x][y] = cv.Int()
		}

		if newspgrid.Arr[x][newspgrid.Height-1] == gen.Config.EmptySymbolCode {
			newspgrid.Arr[x][newspgrid.Height-1] = gen.Config.BlankSymbolCode
		} else if newspgrid.Arr[x][newspgrid.Height-2] == gen.Config.EmptySymbolCode {
			if newspgrid.Arr[x][0] == gen.Config.BlankSymbolCode {
				newspgrid.Arr[x][newspgrid.Height-2] = gen.Config.BlankSymbolCode
				newspgrid.Arr[x][0] = -1
			} else if newspgrid.Arr[x][0] != gen.Config.EmptySymbolCode {
				newspgrid.Arr[x][newspgrid.Height-2] = newspgrid.Arr[x][0]
				newspgrid.Arr[x][0] = -1
			} else {
				newspgrid.Arr[x][0] = -1
			}
		}
	}

	if gen.Config.GigaWeightVM != nil {
		gw := gen.getGigaWeight(gameProp, bcd)

		for x := startX; x < newspgrid.Width; x++ {
			for y := newspgrid.Height - 1; y > 0; y-- {
				if goutils.IndexOfIntSlice(gen.Config.GigaSymbolCodes, newspgrid.Arr[x][y], 0) >= 0 {
					cv, err := gw.RandVal(plugin)
					if err != nil {
						goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:RandVal",
							goutils.Err(err))

						return "", err
					}

					s := cv.Int()

					gen.setGiga(newspgrid, x, y, s, gen.Config.TargetGigaSymbolCodes[newspgrid.Arr[x][y]][s])
				}
				cv, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("RefillTropiCoolSPGrid.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				newspgrid.Arr[x][y] = cv.Int()
			}

			if newspgrid.Arr[x][newspgrid.Height-1] == gen.Config.EmptySymbolCode {
				newspgrid.Arr[x][newspgrid.Height-1] = gen.Config.BlankSymbolCode
			} else if newspgrid.Arr[x][newspgrid.Height-2] == gen.Config.EmptySymbolCode {
				if newspgrid.Arr[x][0] == gen.Config.BlankSymbolCode {
					newspgrid.Arr[x][newspgrid.Height-2] = gen.Config.BlankSymbolCode
					newspgrid.Arr[x][0] = gen.Config.EmptySymbolCode
				} else if newspgrid.Arr[x][1] != gen.Config.EmptySymbolCode {
					newspgrid.Arr[x][newspgrid.Height-2] = newspgrid.Arr[x][0]
					newspgrid.Arr[x][0] = gen.Config.EmptySymbolCode
				}
			}
		}
	}

	gen.AddSPGrid(gen.Config.SPGrid, gameProp, curpr, newspgrid, bcd)

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *RefillTropiCoolSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewRefillTropiCoolSPGrid(name string) IComponent {
	return &RefillTropiCoolSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// json representation used by editor
type jsonRefillTropiCoolSPGrid struct {
	MaxNumber   int      `json:"maxNumber"`
	SPGrid      string   `json:"spGrid"`
	BlankSymbol string   `json:"BlankSymbol"`
	GigaSymbols []string `json:"gigsSymbols"`
	SPSymbols   []string `json:"spSymbols"`
	Weight      string   `json:"weight"`
	GigaWeight  string   `json:"gigaWeight"`
	EmptySymbol string   `json:"emptySymbol"`
}

func (j *jsonRefillTropiCoolSPGrid) build() *RefillTropiCoolSPGridConfig {
	return &RefillTropiCoolSPGridConfig{
		MaxNumber:   j.MaxNumber,
		SPGrid:      j.SPGrid,
		BlankSymbol: j.BlankSymbol,
		GigaSymbols: slices.Clone(j.GigaSymbols),
		SPSymbols:   slices.Clone(j.SPSymbols),
		Weight:      j.Weight,
		GigaWeight:  j.GigaWeight,
		EmptySymbol: j.EmptySymbol,
	}
}

func parseRefillTropiCoolSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillTropiCoolSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillTropiCoolSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRefillTropiCoolSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillTropiCoolSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RefillTropiCoolSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
