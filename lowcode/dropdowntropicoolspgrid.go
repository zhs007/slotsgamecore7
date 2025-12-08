package lowcode

import (
	"log/slog"
	"os"
	"slices"
	"sort"

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
	SPGrid               string   `yaml:"spGrid" json:"spGrid"`
	BlankSymbol          string   `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int      `yaml:"-" json:"-"`
	InitTropiCoolSPGrid  string   `yaml:"initTropiCoolSPGrid" json:"initTropiCoolSPGrid"`
	GenGigaSymbols2      string   `yaml:"genGigaSymbols2" json:"genGigaSymbols2"`
	BrokenSymbols        []string `yaml:"brokenSymbols" json:"brokenSymbols"`
	BrokenSymbolCodes    []int    `yaml:"-" json:"-"`
}

func (cfg *DropDownTropiCoolSPGridConfig) isBroken(sc int) bool {
	return slices.Contains(cfg.BrokenSymbolCodes, sc)
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

	for _, bs := range gen.Config.BrokenSymbols {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[bs]
		if !isok {
			goutils.Error("DropDownTropiCoolSPGrid.InitEx:BrokenSymbols",
				slog.String("BrokenSymbol", bs),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		gen.Config.BrokenSymbolCodes = append(gen.Config.BrokenSymbolCodes, sc)
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *DropDownTropiCoolSPGrid) getInitTropiCoolSPGridData(gameProp *GameProperty) (*InitTropiCoolSPGridData, error) {
	gigaicd := gameProp.GetComponentDataWithName(gen.Config.InitTropiCoolSPGrid)
	if gigaicd == nil {
		goutils.Error("DropDownTropiCoolSPGrid.getInitTropiCoolSPGridData:GetComponentDataWithName",
			slog.String("InitTropiCoolSPGrid", gen.Config.InitTropiCoolSPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	itccd, isok := gigaicd.(*InitTropiCoolSPGridData)
	if !isok {
		goutils.Error("DropDownTropiCoolSPGrid.getInitTropiCoolSPGridData:InitTropiCoolSPGridData",
			slog.String("InitTropiCoolSPGrid", gen.Config.InitTropiCoolSPGrid),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	return itccd, nil
}

func (gen *DropDownTropiCoolSPGrid) getGenGigaSymbols2Data(gameProp *GameProperty) (*GenGigaSymbols2Data, error) {
	gigaicd := gameProp.GetComponentDataWithName(gen.Config.GenGigaSymbols2)
	if gigaicd == nil {
		goutils.Error("DropDownTropiCoolSPGrid.getInitTropiCoolSPGridData:GetComponentDataWithName",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	ggcd, isok := gigaicd.(*GenGigaSymbols2Data)
	if !isok {
		goutils.Error("DropDownTropiCoolSPGrid.getInitTropiCoolSPGridData:GenGigaSymbols2Data",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	return ggcd, nil
}

func (gen *DropDownTropiCoolSPGrid) getGiga(spgrid *sgc7game.GameScene, x int, iicd *InitTropiCoolSPGridData) *gigaData {
	return iicd.getGigaData(x, spgrid.Height-1)
}

func (gen *DropDownTropiCoolSPGrid) getSPGridSymbol(spgrid *sgc7game.GameScene, x int, iicd *InitTropiCoolSPGridData) int {
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

func (gen *DropDownTropiCoolSPGrid) brokenGigaSymbols(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) {
	if gigadata.Y+gigadata.Height >= tgs.Height-1 {
		return
	}

	for y := gigadata.Y + gigadata.Height; y < tgs.Height; y++ {
		for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
			if gen.Config.isBroken(tgs.Arr[x][y]) {
				tgs.Arr[x][y] = -2
			} else {
				newgigadata := ggcd.getGigaData(x, y)
				if newgigadata != nil {
					gen.brokenGigaSymbols(tgs, newgigadata, ggcd)
				}
			}
		}
	}
}

func (gen *DropDownTropiCoolSPGrid) dropdown(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) {
	if gigadata.Y+gigadata.Height >= tgs.Height-1 {
		return
	}

	for y := gigadata.Y + gigadata.Height; y < tgs.Height; y++ {
		for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
			if gen.Config.isBroken(tgs.Arr[x][y]) {
				tgs.Arr[x][y] = -1
			} else {
				newgigadata := ggcd.getGigaData(x, y)
				if newgigadata != nil {
					gen.brokenGigaSymbols(tgs, newgigadata, ggcd)
				}
			}
		}
	}
}

func (gen *DropDownTropiCoolSPGrid) canDownSymbol(tgs *sgc7game.GameScene, x, y int, ggcd *GenGigaSymbols2Data) bool {
	if y+1 >= tgs.Height {
		return false
	}

	ny := y + 1
	if tgs.Arr[x][ny] != -1 && !gen.Config.isBroken(tgs.Arr[x][ny]) {
		newgigadata := ggcd.getGigaData(x, ny)
		if newgigadata != nil {
			return gen.canDownGiga(tgs, newgigadata, ggcd)
		}

		return gen.canDownSymbol(tgs, x, ny, ggcd)
	}

	return true
}

func (gen *DropDownTropiCoolSPGrid) canDownGiga(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) bool {
	if gigadata.Y+gigadata.Height >= tgs.Height-1 {
		return false
	}

	y := gigadata.Y + gigadata.Height
	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		if tgs.Arr[x][y] != -1 && !gen.Config.isBroken(tgs.Arr[x][y]) {

			newgigadata := ggcd.getGigaData(x, y)
			if newgigadata != nil {
				candrop := gen.canDownGiga(tgs, newgigadata, ggcd)
				if !candrop {
					return false
				}
			} else {
				candrop := gen.canDownSymbol(tgs, x, y, ggcd)
				if !candrop {
					return false
				}
			}
		}
	}

	return true
}

func (gen *DropDownTropiCoolSPGrid) downSymbol(tgs *sgc7game.GameScene, x, y int, ggcd *GenGigaSymbols2Data) {
	ny := y + 1
	if tgs.Arr[x][ny] != -1 && !gen.Config.isBroken(tgs.Arr[x][ny]) {
		newgigadata := ggcd.getGigaData(x, ny)
		if newgigadata != nil {
			gen.downGiga(tgs, newgigadata, ggcd)
		} else {
			gen.canDownSymbol(tgs, x, ny, ggcd)
		}
	} else {
		tgs.Arr[x][ny] = tgs.Arr[x][y]
		tgs.Arr[x][y] = -1
	}
}

func (gen *DropDownTropiCoolSPGrid) downGiga(tgs *sgc7game.GameScene, gigadata *gigaData, ggcd *GenGigaSymbols2Data) {
	y := gigadata.Y + gigadata.Height
	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		if tgs.Arr[x][y] != -1 && !gen.Config.isBroken(tgs.Arr[x][y]) {
			newgigadata := ggcd.getGigaData(x, y)
			if newgigadata != nil {
				gen.downGiga(tgs, newgigadata, ggcd)
			} else {
				gen.downSymbol(tgs, x, y, ggcd)
			}
		}
	}

	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		for y := gigadata.Y; y < gigadata.Y+gigadata.Height; y++ {
			tgs.Arr[x][y] = -1
		}
	}

	gigadata.Y++

	for x := gigadata.X; x < gigadata.X+gigadata.Width; x++ {
		for y := gigadata.Y; y < gigadata.Y+gigadata.Height; y++ {
			tgs.Arr[x][y] = gigadata.CurSymbolCode
		}
	}
}

func (gen *DropDownTropiCoolSPGrid) dropdownSPGigaList(gameProp *GameProperty, gs *sgc7game.GameScene, lst []*gigaData, ggcd *GenGigaSymbols2Data) *sgc7game.GameScene {
	ngs := gs.CloneEx(gameProp.PoolScene)

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].getBottom() > lst[j].getBottom()
	})

	for _, gigadata := range lst {
		for {
			candrop := gen.canDownGiga(ngs, gigadata, ggcd)
			if !candrop {
				break
			}

			gen.downGiga(ngs, gigadata, ggcd)
		}
	}

	// for {
	// 	for x := range ngs.Arr {
	// 		for y := len(ngs.Arr[x]) - 1; y >= 0; {
	// 			if ngs.Arr[x][y] < 0 {
	// 				hass := false
	// 				for y1 := y - 1; y1 >= 0; y1-- {
	// 					// if giga then break
	// 					if ggcd.getGigaData(x, y1) != nil {
	// 						break
	// 					}

	// 					if ngs.Arr[x][y1] >= 0 {
	// 						ngs.Arr[x][y] = ngs.Arr[x][y1]
	// 						ngs.Arr[x][y1] = -1

	// 						hass = true
	// 						y--
	// 						break
	// 					}
	// 				}

	// 				if !hass {
	// 					break
	// 				}
	// 			} else {
	// 				y--
	// 			}
	// 		}
	// 	}

	// 	isdrop := dropDownSymbols.dropdownGiga(gameProp, ngs, gigacd)
	// 	if !isdrop {
	// 		break
	// 	}
	// }

	return ngs
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

	iicd, err := gen.getInitTropiCoolSPGridData(gameProp)
	if err != nil {
		goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:getInitTropiCoolSPGridData",
			slog.String("InitTropiCoolSPGrid", gen.Config.InitTropiCoolSPGrid),
			goutils.Err(err))

		return "", err
	}

	ggcd, err := gen.getGenGigaSymbols2Data(gameProp)
	if err != nil {
		goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:getGenGigaSymbols2Data",
			slog.String("GenGigaSymbols2", gen.Config.GenGigaSymbols2),
			goutils.Err(err))

		return "", err
	}

	newspgrid := spgrid.CloneEx(gameProp.PoolScene)

	gs := gameProp.SceneStack.GetTopSceneEx(curpr, prs)
	ngs := gs.CloneEx(gameProp.PoolScene)

	newgigadatalist := []*gigaData{}

	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				gigadata := gen.getGiga(newspgrid, x, iicd)
				if gigadata != nil {
					ny := gigadata.checkWithBottomY(ngs, x, y)
					if ny >= 0 {
						err = gigadata.putInWithBottomY(ngs, x, ny)
						if err != nil {
							goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:putInWithBottomY",
								slog.Int("x", x),
								slog.Int("y", y),
								goutils.Err(err))

							return "", err
						}

						newgigadata := &gigaData{
							X:             x,
							Y:             ny - gigadata.Height + 1,
							Width:         gigadata.Width,
							Height:        gigadata.Height,
							SymbolCode:    gigadata.SymbolCode,
							CurSymbolCode: gigadata.CurSymbolCode,
						}

						newgigadatalist = append(newgigadatalist, newgigadata)

						ggcd.gigaData = append(ggcd.gigaData, newgigadata)
						iicd.rmGigaData(gigadata)

						continue
					}

					err = iicd.splitGigaData(newspgrid, gigadata)
					if err != nil {
						goutils.Error("DropDownTropiCoolSPGrid.OnPlayGame:splitGigaData",
							slog.Int("x", x),
							slog.Int("y", y),
							goutils.Err(err))

						return "", err
					}
				}

				sym := gen.getSPGridSymbol(newspgrid, x, iicd)
				if sym == -1 {
					break
				}

				ngs.Arr[x][y] = sym
			}
		}
	}

	gen.AddScene(gameProp, curpr, ngs, bcd)

	if len(newgigadatalist) > 0 {
		ngs3 := gen.dropdownSPGigaList(gameProp, ngs, newgigadatalist, ggcd)

		gen.AddScene(gameProp, curpr, ngs3, bcd)
	}

	ngs2 := ngs
	for x := 0; x < gs.Width; x++ {
		for y := gs.Height - 1; y >= 0; y-- {
			if ngs2.Arr[x][y] == gen.Config.BlankSymbolCode {
				if ngs2 == ngs {
					ngs2 = ngs.CloneEx(gameProp.PoolScene)
				}

				ngs2.Arr[x][y] = -1
			}
		}
	}

	if ngs2 != ngs {
		for _, arr := range ngs2.Arr {
			for y := len(arr) - 1; y >= 0; {
				if arr[y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						if arr[y1] != -1 {
							arr[y] = arr[y1]
							arr[y1] = -1

							hass = true
							y--
							break
						}
					}

					if !hass {
						break
					}
				} else {
					y--
				}
			}
		}

		gen.AddScene(gameProp, curpr, ngs2, bcd)
	}

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
// "initTropiCoolSPGrid": "bg-spgrid-init"
// "genGigaSymbols2": "bg-gengiga"
// "brokenSymbols": [
//
//	"H1",
//	"H2",
//	"H3",
//	"H4",
//	"H5",
//	"L1",
//	"L2",
//	"L3",
//	"L4"
//
// ]
type jsonDropDownTropiCoolSPGrid struct {
	SPGrid              string   `json:"spGrid"`
	BlankSymbol         string   `json:"BlankSymbol"`
	InitTropiCoolSPGrid string   `json:"initTropiCoolSPGrid"`
	GenGigaSymbols2     string   `json:"genGigaSymbols2"`
	BrokenSymbols       []string `json:"brokenSymbols"`
}

func (j *jsonDropDownTropiCoolSPGrid) build() *DropDownTropiCoolSPGridConfig {
	return &DropDownTropiCoolSPGridConfig{
		SPGrid:              j.SPGrid,
		BlankSymbol:         j.BlankSymbol,
		InitTropiCoolSPGrid: j.InitTropiCoolSPGrid,
		GenGigaSymbols2:     j.GenGigaSymbols2,
		BrokenSymbols:       slices.Clone(j.BrokenSymbols),
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
