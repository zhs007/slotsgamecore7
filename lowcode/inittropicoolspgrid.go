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

const InitTropiCoolSPGridTypeName = "initTropiCoolSPGrid"

type InitTropiCoolSPGridData struct {
	BasicComponentData
	isSpBounsAdded bool
	vw             *sgc7game.ValWeights2
	vwNoSpBonus    *sgc7game.ValWeights2
	gigaData       []*gigaData
}

func (itcdpg *InitTropiCoolSPGridData) alignStep(sx int) {
	for _, v := range itcdpg.gigaData {
		if v.X > sx {
			v.X--
		}
	}
}

func (itcdpg *InitTropiCoolSPGridData) getGigaData(x, y int) *gigaData {
	for _, v := range itcdpg.gigaData {
		if v.X <= x && v.Y <= y && v.X+v.Width-1 >= x && v.Y+v.Height-1 >= y {
			return v
		}
	}

	return nil
}

func (itcdpg *InitTropiCoolSPGridData) rmGigaData(gd *gigaData) error {
	i := slices.Index(itcdpg.gigaData, gd)
	if i < 0 {
		goutils.Error("InitTropiCoolSPGridData.rmGigaData:Index",
			goutils.Err(ErrInvalidComponentData))

		return ErrInvalidComponentData
	}

	itcdpg.gigaData = slices.Delete(itcdpg.gigaData, i, i+1)

	return nil
}

func (itcdpg *InitTropiCoolSPGridData) splitGigaData(gs *sgc7game.GameScene, gd *gigaData) error {
	i := slices.Index(itcdpg.gigaData, gd)
	if i < 0 {
		goutils.Error("InitTropiCoolSPGridData.rmGigaData:Index",
			goutils.Err(ErrInvalidComponentData))

		return ErrInvalidComponentData
	}

	for ix := gd.X; ix < gd.X+gd.Width; ix++ {
		for iy := gd.Y; iy < gd.Y+gd.Height; iy++ {
			gs.Arr[ix][iy] = gd.SymbolCode
		}
	}

	itcdpg.gigaData = slices.Delete(itcdpg.gigaData, i, i+1)

	return nil
}

func (itcdpg *InitTropiCoolSPGridData) OnNewGame(gameProp *GameProperty, component IComponent) {
	itcdpg.BasicComponentData.OnNewGame(gameProp, component)

	itcdpg.isSpBounsAdded = false
	itcdpg.gigaData = nil
}

func (itcdpg *InitTropiCoolSPGridData) onNewStep() {
}

func (itcdpg *InitTropiCoolSPGridData) Clone() IComponentData {
	target := &InitTropiCoolSPGridData{
		BasicComponentData: itcdpg.CloneBasicComponentData(),
		isSpBounsAdded:     itcdpg.isSpBounsAdded,
	}

	return target
}

// InitTropiCoolSPGridConfig - configuration for InitTropiCoolSPGrid
type InitTropiCoolSPGridConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MaxNumber            int                   `yaml:"maxNumber" json:"maxNumber"`
	SPGrid               string                `yaml:"spGrid" json:"spGrid"`
	BlankSymbol          string                `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int                   `yaml:"-" json:"-"`
	GigaSymbols          []string              `yaml:"gigSymbols" json:"gigSymbols"`
	GigaSymbolCodes      []int                 `yaml:"-" json:"-"`
	SPSymbols            []string              `yaml:"spSymbols" json:"spSymbols"`
	SPSymbolCodes        []int                 `yaml:"-" json:"-"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVM             *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	GigaWeight           string                `yaml:"gigaWeight" json:"gigaWeight"`
	GigaWeightVM         *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	EmptySymbol          string                `yaml:"emptySymbol" json:"emptySymbol"`
	EmptySymbolCode      int                   `yaml:"-" json:"-"`
	SpBonusSymbol        string                `yaml:"spBonusSymbol" json:"spBonusSymbol"`
	SpBonusSymbolCode    int                   `yaml:"-" json:"-"`
	SpBonusSymbolCode2   int                   `yaml:"-" json:"-"`
	MapGigaSymbolCodes   map[int][]int         `yaml:"-" json:"-"`
	MapControls          map[string][]*Award   `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *InitTropiCoolSPGridConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type InitTropiCoolSPGrid struct {
	*BasicComponent `json:"-"`
	Config          *InitTropiCoolSPGridConfig `json:"config"`
}

// Init - load from file
func (gen *InitTropiCoolSPGrid) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("InitTropiCoolSPGrid.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &InitTropiCoolSPGridConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("InitTropiCoolSPGrid.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *InitTropiCoolSPGrid) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*InitTropiCoolSPGridConfig)
	gen.Config.ComponentType = InitTropiCoolSPGridTypeName

	// process symbol codes
	if gen.Config.BlankSymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.BlankSymbol]
		if !isok {
			goutils.Error("InitTropiCoolSPGrid.InitEx:BlankSymbol",
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
			goutils.Error("InitTropiCoolSPGrid.InitEx:GigaSymbols",
				slog.String("GigaSymbols", v),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.GigaSymbolCodes = append(gen.Config.GigaSymbolCodes, code)
	}

	if gen.Config.EmptySymbol != "" {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.EmptySymbol]
		if !isok {
			goutils.Error("InitTropiCoolSPGrid.InitEx:EmptySymbol",
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
			goutils.Error("InitTropiCoolSPGrid.InitEx:GigaSymbols",
				slog.String("GigaSymbols", v),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.SPSymbolCodes = append(gen.Config.SPSymbolCodes, code)
	}

	code, isok := pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.SpBonusSymbol]
	if !isok {
		goutils.Error("InitTropiCoolSPGrid.InitEx:SpBonusSymbol",
			slog.String("SpBonusSymbol", gen.Config.SpBonusSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	gen.Config.SpBonusSymbolCode = code

	code, isok = pool.Config.GetDefaultPaytables().MapSymbols[gen.Config.SpBonusSymbol+"_2"]
	if !isok {
		goutils.Error("InitTropiCoolSPGrid.InitEx:SpBonusSymbol",
			slog.String("SpBonusSymbol2", gen.Config.SpBonusSymbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	gen.Config.SpBonusSymbolCode2 = code

	// weights
	if gen.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(gen.Config.Weight, true)
		if err != nil {
			goutils.Error("InitTropiCoolSPGrid.Init:LoadStrWeights",
				slog.String("Weight", gen.Config.Weight),
				goutils.Err(err))

			return err
		}

		gen.Config.WeightVM = vw2
	}

	if gen.Config.GigaWeight != "" {
		vw2, err := pool.LoadIntWeights(gen.Config.GigaWeight, true)
		if err != nil {
			goutils.Error("InitTropiCoolSPGrid.Init:LoadStrWeights",
				slog.String("Weight", gen.Config.GigaWeight),
				goutils.Err(err))

			return err
		}

		gen.Config.GigaWeightVM = vw2

		gen.Config.MapGigaSymbolCodes = make(map[int][]int)
		for _, sym := range gen.Config.GigaSymbols {
			sc, isok := pool.DefaultPaytables.MapSymbols[sym]
			if !isok {
				goutils.Error("GenGigaSymbols2.InitEx:Symbols2",
					slog.String("symbol", sym),
					goutils.Err(ErrInvalidSymbol))

				return ErrInvalidSymbol
			}

			gen.Config.MapGigaSymbolCodes[sc] = make([]int, pool.Config.Width-1)
			gen.Config.MapGigaSymbolCodes[sc][0] = -1

			for gi := 2; gi < pool.Config.Width; gi++ {
				str := fmt.Sprintf("%v_%d", sym, gi)
				gsc, isok := pool.DefaultPaytables.MapSymbols[str]
				if isok {
					gen.Config.MapGigaSymbolCodes[sc][gi-1] = gsc
				} else {
					gen.Config.MapGigaSymbolCodes[sc][gi-1] = -1
				}
			}
		}

		// gen.Config.TargetGigaSymbolCodes = make(map[int](map[int]int))

		// for _, v := range vw2.Vals {
		// 	s := v.Int()

		// 	if s == 0 {
		// 		continue
		// 	}

		// 	for i, sc := range gen.Config.GigaSymbolCodes {
		// 		if gen.Config.TargetGigaSymbolCodes[sc] == nil {
		// 			gen.Config.TargetGigaSymbolCodes[sc] = make(map[int]int)
		// 		}

		// 		cs := fmt.Sprintf("%v_%v", gen.Config.GigaSymbols[i], s)
		// 		csc, isok := pool.Config.GetDefaultPaytables().MapSymbols[cs]
		// 		if isok {
		// 			gen.Config.TargetGigaSymbolCodes[sc][s] = csc
		// 		}
		// 	}
		// }
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *InitTropiCoolSPGrid) getWeight(gameProp *GameProperty, cd *InitTropiCoolSPGridData) (*sgc7game.ValWeights2, *sgc7game.ValWeights2) {
	str := cd.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("InitTropiCoolSPGrid.getWeight:LoadIntWeights",
				goutils.Err(err))

			return nil, nil
		}

		if cd.vw != vw2 {
			nospbonusvw, err := vw2.CloneExcludeVal(sgc7game.NewIntValEx(gen.Config.SpBonusSymbolCode))
			if err != nil {
				goutils.Error("InitTropiCoolSPGrid.getWeight:CloneExcludeVal",
					goutils.Err(err))

				return nil, nil
			}

			cd.vw = vw2
			cd.vwNoSpBonus = nospbonusvw
		}

		return vw2, cd.vwNoSpBonus
	}

	if cd.vw != gen.Config.WeightVM {
		nospbonusvw, err := gen.Config.WeightVM.CloneExcludeVal(sgc7game.NewIntValEx(gen.Config.SpBonusSymbolCode))
		if err != nil {
			goutils.Error("InitTropiCoolSPGrid.getWeight:CloneExcludeVal",
				goutils.Err(err))

			return nil, nil
		}

		cd.vw = gen.Config.WeightVM
		cd.vwNoSpBonus = nospbonusvw
	}

	return gen.Config.WeightVM, cd.vwNoSpBonus
}

func (gen *InitTropiCoolSPGrid) getGigaWeight(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVGigaWeight)
	if str != "" {
		vw2, err := gameProp.Pool.LoadIntWeights(str, true)
		if err != nil {
			goutils.Error("InitTropiCoolSPGrid.getGigaWeight:LoadIntWeights",
				goutils.Err(err))

			return nil
		}

		return vw2
	}

	return gen.Config.GigaWeightVM
}

func (gen *InitTropiCoolSPGrid) setGiga(gs *sgc7game.GameScene, x int, y int, s int, srcsc int, c int, cd *InitTropiCoolSPGridData) {
	if x+s-1 >= gs.Width || y+s-1 >= gs.Height {
		return
	}

	for ix := x; ix < x+s; ix++ {
		for iy := y; iy < y+s; iy++ {
			if goutils.IndexOfIntSlice(gen.Config.SPSymbolCodes, gs.Arr[ix][iy], 0) >= 0 {

				return
			}
		}
	}

	for ix := x; ix < x+s; ix++ {
		for iy := y; iy < y+s; iy++ {
			gs.Arr[ix][iy] = c
		}
	}

	cd.gigaData = append(cd.gigaData, &gigaData{
		X:             x,
		Y:             y,
		Width:         s,
		Height:        s,
		SymbolCode:    srcsc,
		CurSymbolCode: c,
	})
}

// OnPlayGame - minimal implementation: does nothing but advance
func (gen *InitTropiCoolSPGrid) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// Initialize an SPGrid for this component on the current play result
	cd := icd.(*InitTropiCoolSPGridData)

	stackSPGrid, isok := gameProp.MapSPGridStack[gen.Config.SPGrid]
	if !isok {
		goutils.Error("InitTropiCoolSPGrid.OnPlayGame:MapSPGridStack",
			slog.String("SPGrid", gen.Config.SPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	gs := gameProp.PoolScene.New2(stackSPGrid.Width, stackSPGrid.Height, gen.Config.EmptySymbolCode)
	if gs == nil {
		goutils.Error("InitTropiCoolSPGrid.OnPlayGame:New2",
			slog.Int("Width", stackSPGrid.Width),
			slog.Int("Height", stackSPGrid.Height),
			slog.Int("EmptySymbolCode", gen.Config.EmptySymbolCode),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	vw, vwsp := gen.getWeight(gameProp, cd)
	if vw == nil {
		goutils.Error("InitTropiCoolSPGrid.OnPlayGame:getWeight",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	for x, arr := range gs.Arr {
		for y := range arr {
			if cd.isSpBounsAdded {
				cv, err := vwsp.RandVal(plugin)
				if err != nil {
					goutils.Error("InitTropiCoolSPGrid.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				gs.Arr[x][y] = cv.Int()
			} else {
				cv, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("InitTropiCoolSPGrid.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				if cv.Int() == gen.Config.SpBonusSymbolCode {
					cd.isSpBounsAdded = true

					gs.Arr[x][len(arr)-1] = gen.Config.SpBonusSymbolCode2
					gs.Arr[x][len(arr)-2] = gen.Config.SpBonusSymbolCode2
					gs.Arr[x][len(arr)-3] = gen.Config.EmptySymbolCode

					break
				}

				gs.Arr[x][y] = cv.Int()
			}
		}

		if gs.Arr[x][gs.Height-1] == gen.Config.EmptySymbolCode {
			gs.Arr[x][gs.Height-1] = gen.Config.BlankSymbolCode
		}

		if gs.Arr[x][gs.Height-2] == gen.Config.EmptySymbolCode {
			if gs.Arr[x][0] != gen.Config.EmptySymbolCode {
				gs.Arr[x][gs.Height-2] = gs.Arr[x][0]
				gs.Arr[x][0] = -1
			} else {
				gs.Arr[x][gs.Height-2] = -1
				gs.Arr[x][0] = -1
			}
		}

		if gs.Arr[x][gs.Height-3] == gen.Config.EmptySymbolCode {
			gs.Arr[x][0] = -1
		}
	}

	if gen.Config.GigaWeightVM != nil {
		gw := gen.getGigaWeight(gameProp, &cd.BasicComponentData)

		for x, _ := range gs.Arr {
			for y := 0; y < gs.Height; y++ {
				if goutils.IndexOfIntSlice(gen.Config.GigaSymbolCodes, gs.Arr[x][y], 0) >= 0 {
					cv, err := gw.RandVal(plugin)
					if err != nil {
						goutils.Error("InitTropiCoolSPGrid.OnPlayGame:RandVal",
							goutils.Err(err))

						return "", err
					}

					s := cv.Int()

					if s > 1 {
						gen.setGiga(gs, x, y, s, gs.Arr[x][y], gen.Config.MapGigaSymbolCodes[gs.Arr[x][y]][s-1], cd)
					}
				}
			}
		}
	}

	gen.AddSPGrid(gen.Config.SPGrid, gameProp, curpr, gs, &cd.BasicComponentData)

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - output to asciigame (no-op)
func (gen *InitTropiCoolSPGrid) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func (gen *InitTropiCoolSPGrid) NewComponentData() IComponentData {
	cd := &InitTropiCoolSPGridData{}

	return cd
}

func NewInitTropiCoolSPGrid(name string) IComponent {
	return &InitTropiCoolSPGrid{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "maxNumber": 0,
// "spGrid": "bg-spgrid",
// "BlankSymbol": "BN",
// "gigsSymbols": [
//
//	"WL",
//	"LW3",
//	"MY",
//	"RS"
//
// ],
// "spSymbols": [
//
//	"SC"
//
// ],
// "weight": "bgspgridsymsweight",
// "gigaWeight": "bgspgridgigaweight",
// "emptySymbol": "EM"
// "spBonusSymbol": "EL"
type jsonInitTropiCoolSPGrid struct {
	MaxNumber     int      `json:"maxNumber"`
	SPGrid        string   `json:"spGrid"`
	BlankSymbol   string   `json:"BlankSymbol"`
	GigaSymbols   []string `json:"gigsSymbols"`
	SPSymbols     []string `json:"spSymbols"`
	Weight        string   `json:"weight"`
	GigaWeight    string   `json:"gigaWeight"`
	EmptySymbol   string   `json:"emptySymbol"`
	SpBonusSymbol string   `json:"spBonusSymbol"`
}

func (j *jsonInitTropiCoolSPGrid) build() *InitTropiCoolSPGridConfig {
	return &InitTropiCoolSPGridConfig{
		MaxNumber:     j.MaxNumber,
		SPGrid:        j.SPGrid,
		BlankSymbol:   j.BlankSymbol,
		GigaSymbols:   slices.Clone(j.GigaSymbols),
		SPSymbols:     slices.Clone(j.SPSymbols),
		Weight:        j.Weight,
		GigaWeight:    j.GigaWeight,
		EmptySymbol:   j.EmptySymbol,
		SpBonusSymbol: j.SpBonusSymbol,
	}
}

func parseInitTropiCoolSPGrid(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseInitTropiCoolSPGrid:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseInitTropiCoolSPGrid:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonInitTropiCoolSPGrid{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseInitTropiCoolSPGrid:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: InitTropiCoolSPGridTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
