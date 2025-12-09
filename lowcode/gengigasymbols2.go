package lowcode

import (
	"fmt"
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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const GenGigaSymbols2TypeName = "genGigaSymbols2"

type gigaData struct {
	SymbolCode    int
	CurSymbolCode int
	Width         int
	Height        int
	X             int
	Y             int
	od            [][]int
}

func (gd *gigaData) chgSymbol(gs *sgc7game.GameScene, newSymbolCode int, cfg *GenGigaSymbols2Config) {
	gd.CurSymbolCode = cfg.GigaSymbolCodes[newSymbolCode][gd.Width-1]
	gd.SymbolCode = newSymbolCode

	for tx := gd.X; tx < gd.X+gd.Width; tx++ {
		for ty := gd.Y; ty < gd.Y+gd.Height; ty++ {
			gs.Arr[tx][ty] = gd.CurSymbolCode
		}
	}
}

func (gd *gigaData) inPos(posData *PosData) bool {
	for i := 0; i < posData.Len(); i++ {
		x := posData.pos[i*2]
		y := posData.pos[i*2+1]

		if gd.X <= x && gd.Y <= y && gd.X+gd.Width-1 >= x && gd.Y+gd.Height-1 >= y {
			return true
		}
	}

	return false
}

func (gd *gigaData) getBottom() int {
	return gd.Y + gd.Height - 1
}

func (gd *gigaData) newOD() {
	gd.od = make([][]int, gd.Width)
	for x := 0; x < gd.Width; x++ {
		gd.od[x] = make([]int, gd.Height)
	}
}

func (gd *gigaData) checkWithBottomY(gs *sgc7game.GameScene, x, bottomY int) int {
	ny := bottomY

	for tx := 0; tx < gd.Width; tx++ {
		if x+tx >= gs.Width {
			return -1
		}

		for ty := 0; ty < gd.Height; ty++ {
			if bottomY-ty < 0 {
				return -1
			}

			if gs.Arr[x+tx][bottomY-ty] >= 0 {
				return gd.checkWithBottomY(gs, x, bottomY-1)
			}
		}
	}

	return ny
}

func (gd *gigaData) putInWithBottomY(gs *sgc7game.GameScene, x, bottomY int) error {
	for tx := 0; tx < gd.Width; tx++ {
		if x+tx >= gs.Width {
			goutils.Error("gigaData.putInWithBottomY:out of range x",
				slog.Int("x", x),
				slog.Int("tx", tx),
				slog.Int("gsWidth", gs.Width),
				goutils.Err(ErrInvalidComponentData))

			return ErrInvalidComponentData
		}

		for ty := 0; ty < gd.Height; ty++ {
			if bottomY-ty < 0 {
				goutils.Error("gigaData.putInWithBottomY:out of range y",
					slog.Int("bottomY", bottomY),
					slog.Int("ty", ty),
					slog.Int("gsHeight", gs.Height),
					goutils.Err(ErrInvalidComponentData))

				return ErrInvalidComponentData
			}

			gs.Arr[x+tx][bottomY-ty] = gd.CurSymbolCode
		}
	}

	return nil
}

type GenGigaSymbols2Data struct {
	BasicComponentData
	gigaData []*gigaData
	cfg      *GenGigaSymbols2Config
}

func (genGigaSymbols2Data *GenGigaSymbols2Data) removeSymbol(x, y int) {
	for i, v := range genGigaSymbols2Data.gigaData {
		if v.X <= x && v.Y <= y && v.X+v.Width-1 >= x && v.Y+v.Height-1 >= y {
			genGigaSymbols2Data.gigaData = append(genGigaSymbols2Data.gigaData[:i], genGigaSymbols2Data.gigaData[i+1:]...)

			return
		}
	}
}

func (genGigaSymbols2Data *GenGigaSymbols2Data) calcDropdown(gs *sgc7game.GameScene, gd *gigaData) int {
	cy := gd.Y
	for gy := gd.Y + gd.Height; gy < gs.Height; gy++ {
		cn := 0
		for gx := gd.X; gx <= gd.X+gd.Width-1; gx++ {
			if gs.Arr[gx][gy] != -1 {
				cn++
			}
		}

		if cn == gd.Width {
			return cy
		}

		for gx := gd.X; gx < gd.X+gd.Width-1; gx++ {
			if slices.Contains(genGigaSymbols2Data.cfg.SpSymbolCodes, gs.Arr[gx][gy]) && genGigaSymbols2Data.getGigaData(gx, gy) != nil {
				return cy
			}
		}

		cy++
	}

	return cy
}

func (genGigaSymbols2Data *GenGigaSymbols2Data) getGigaData(x, y int) *gigaData {
	for _, v := range genGigaSymbols2Data.gigaData {
		if v.X <= x && v.Y <= y && v.X+v.Width-1 >= x && v.Y+v.Height-1 >= y {
			return v
		}
	}

	return nil
}

func (genGigaSymbols2Data *GenGigaSymbols2Data) sortGigaData() {
	sort.Slice(genGigaSymbols2Data.gigaData, func(i, j int) bool {
		return genGigaSymbols2Data.gigaData[i].getBottom() > genGigaSymbols2Data.gigaData[j].getBottom()
	})
}

// OnNewGame -
func (genGigaSymbols2Data *GenGigaSymbols2Data) OnNewGame(gameProp *GameProperty, component IComponent) {
	genGigaSymbols2Data.BasicComponentData.OnNewGame(gameProp, component)

	genGigaSymbols2Data.gigaData = nil
}

// OnNewStep -
func (genGigaSymbols2Data *GenGigaSymbols2Data) OnNewStep() {
}

// Clone
func (genGigaSymbols2Data *GenGigaSymbols2Data) Clone() IComponentData {
	target := &GenGigaSymbols2Data{
		BasicComponentData: genGigaSymbols2Data.CloneBasicComponentData(),
	}

	for _, v := range genGigaSymbols2Data.gigaData {
		target.gigaData = append(target.gigaData, &gigaData{
			SymbolCode:    v.SymbolCode,
			CurSymbolCode: v.CurSymbolCode,
			Width:         v.Width,
			Height:        v.Height,
			X:             v.X,
			Y:             v.Y,
		})
	}

	return target
}

// BuildPBComponentData
func (genGigaSymbols2Data *GenGigaSymbols2Data) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: genGigaSymbols2Data.BuildPBBasicComponentData(),
	}
}

// GenGigaSymbols2Config - configuration for GenGigaSymbols2
type GenGigaSymbols2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MaxNumber            int                   `yaml:"maxNumber" json:"maxNumber"`
	Symbols              []string              `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int                 `yaml:"-" json:"-"`
	SpSymbols            []string              `yaml:"spSymbols" json:"spSymbols"`
	SpSymbolCodes        []int                 `yaml:"-" json:"-"`
	Weight               string                `yaml:"weight" json:"weight"`
	GigaSymbolCodes      map[int][]int         `yaml:"-" json:"-"`
	WeightVW             *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	BrokenSymbols        []string              `yaml:"brokenSymbols" json:"brokenSymbols"`
	BrokenSymbolCodes    []int                 `yaml:"-" json:"-"`
}

func (cfg *GenGigaSymbols2Config) isBroken(sc int) bool {
	return slices.Contains(cfg.BrokenSymbolCodes, sc)
}

// SetLinkComponent
func (cfg *GenGigaSymbols2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenGigaSymbols2 struct {
	*BasicComponent `json:"-"`
	Config          *GenGigaSymbols2Config `json:"config"`
}

// Init -
func (genGigaSymbols2 *GenGigaSymbols2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenGigaSymbols2.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenGigaSymbols2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenGigaSymbols2.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genGigaSymbols2.InitEx(cfg, pool)
}

// InitEx -
func (genGigaSymbols2 *GenGigaSymbols2) InitEx(cfg any, pool *GamePropertyPool) error {
	genGigaSymbols2.Config = cfg.(*GenGigaSymbols2Config)
	genGigaSymbols2.Config.ComponentType = GenGigaSymbols2TypeName

	for _, v := range genGigaSymbols2.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("GenGigaSymbols2.InitEx:Symbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbols2.Config.SymbolCodes = append(genGigaSymbols2.Config.SymbolCodes, sc)
	}

	for _, v := range genGigaSymbols2.Config.SpSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("GenGigaSymbols2.InitEx:SpSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbols2.Config.SpSymbolCodes = append(genGigaSymbols2.Config.SpSymbolCodes, sc)
	}

	weightVM, err := pool.LoadIntWeights(genGigaSymbols2.Config.Weight, true)
	if err != nil {
		goutils.Error("GenGigaSymbols2.InitEx:LoadIntWeights",
			slog.String("weight", genGigaSymbols2.Config.Weight),
			goutils.Err(err))

		return err
	}

	genGigaSymbols2.Config.WeightVW = weightVM

	genGigaSymbols2.Config.GigaSymbolCodes = make(map[int][]int)
	for _, sym := range genGigaSymbols2.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[sym]
		if !isok {
			goutils.Error("GenGigaSymbols2.InitEx:Symbols2",
				slog.String("symbol", sym),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbols2.Config.GigaSymbolCodes[sc] = make([]int, pool.Config.Width-1)
		genGigaSymbols2.Config.GigaSymbolCodes[sc][0] = -1

		for gi := 2; gi < pool.Config.Width; gi++ {
			str := fmt.Sprintf("%v_%d", sym, gi)
			gsc, isok := pool.DefaultPaytables.MapSymbols[str]
			if isok {
				genGigaSymbols2.Config.GigaSymbolCodes[sc][gi-1] = gsc
			} else {
				genGigaSymbols2.Config.GigaSymbolCodes[sc][gi-1] = -1
			}
		}
	}

	for _, bs := range genGigaSymbols2.Config.BrokenSymbols {
		sc, isok := pool.Config.GetDefaultPaytables().MapSymbols[bs]
		if !isok {
			goutils.Error("DropDownTropiCoolSPGrid.InitEx:BrokenSymbols",
				slog.String("BrokenSymbol", bs),
				goutils.Err(ErrInvalidComponentConfig))

			return ErrInvalidComponentConfig
		}

		genGigaSymbols2.Config.BrokenSymbolCodes = append(genGigaSymbols2.Config.BrokenSymbolCodes, sc)
	}

	genGigaSymbols2.onInit(&genGigaSymbols2.Config.BasicComponentConfig)

	return nil
}

func (genGigaSymbols2 *GenGigaSymbols2) getWeight(gameProp *GameProperty, bcd *BasicComponentData) *sgc7game.ValWeights2 {
	str := bcd.GetConfigVal(CCVWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return genGigaSymbols2.Config.WeightVW
}

func (genGigaSymbols2 *GenGigaSymbols2) getMaxNumber(gameProp *GameProperty, bcd *BasicComponentData) int {
	v, isok := bcd.GetConfigIntVal(CCVMaxNumber)
	if isok {
		return v
	}

	return genGigaSymbols2.Config.MaxNumber
}

func (genGigaSymbols2 *GenGigaSymbols2) isCanGiga(x, y int, giga int, gs *sgc7game.GameScene) bool {
	for gx := 0; gx <= giga-1; gx++ {
		if x+gx >= gs.Width {
			return false
		}

		for gy := 0; gy <= giga-1; gy++ {
			if gx == 0 && gy == 0 {
				continue
			}

			if y+gy >= gs.Height {
				return false
			}

			if !slices.Contains(genGigaSymbols2.Config.SymbolCodes, gs.Arr[x+gx][y+gy]) {
				return false
			}
		}
	}

	return true
}

func (genGigaSymbols2 *GenGigaSymbols2) procGiga(gs *sgc7game.GameScene, vw *sgc7game.ValWeights2, gameProp *GameProperty,
	plugin sgc7plugin.IPlugin, cd *GenGigaSymbols2Data, maxNumber int) (*sgc7game.GameScene, error) {

	ngs := gs

	for x := 0; x < gs.Width-1; x++ {
		for y := 0; y < gs.Height-1; y++ {

			symbolCode := ngs.Arr[x][y]

			if slices.Contains(genGigaSymbols2.Config.SymbolCodes, symbolCode) {
				cr, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("GenGigaSymbols2.procGiga:RandVal",
						goutils.Err(err))

					return nil, err
				}

				gigasize := cr.Int()

				if gigasize > 1 && genGigaSymbols2.Config.GigaSymbolCodes[symbolCode][gigasize-1] != -1 {

					if genGigaSymbols2.isCanGiga(x, y, gigasize, ngs) {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						for gx := 0; gx < gigasize; gx++ {
							for gy := 0; gy < gigasize; gy++ {
								ngs.Arr[x+gx][y+gy] = genGigaSymbols2.Config.GigaSymbolCodes[symbolCode][gigasize-1]
							}
						}

						gd := &gigaData{
							SymbolCode:    symbolCode,
							CurSymbolCode: genGigaSymbols2.Config.GigaSymbolCodes[symbolCode][gigasize-1],
							Width:         cr.Int(),
							Height:        cr.Int(),
							X:             x,
							Y:             y,
						}

						cd.gigaData = append(cd.gigaData, gd)

						if maxNumber > 0 {
							if len(cd.gigaData) >= maxNumber {
								return ngs, nil
							}
						}
					}
				}

			}
		}
	}

	return ngs, nil
}

// playgame placeholder
func (genGigaSymbols2 *GenGigaSymbols2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*GenGigaSymbols2Data)

	cd.OnNewStep()

	gs := genGigaSymbols2.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("GenGigaSymbols2.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	vw := genGigaSymbols2.getWeight(gameProp, &cd.BasicComponentData)
	maxNumber := genGigaSymbols2.getMaxNumber(gameProp, &cd.BasicComponentData)

	ngs, err := genGigaSymbols2.procGiga(gs, vw, gameProp, plugin, cd, maxNumber)
	if err != nil {
		goutils.Error("GenGigaSymbols2.OnPlayGame:procGiga",
			goutils.Err(err))

		return "", err
	}

	if ngs != gs {
		genGigaSymbols2.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

		genGigaSymbols2.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

		nc := genGigaSymbols2.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	// Placeholder: Do nothing
	nc := genGigaSymbols2.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genGigaSymbols2 *GenGigaSymbols2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	// Placeholder: nothing to output
	return nil
}

// NewComponentData -
func (genGigaSymbols2 *GenGigaSymbols2) NewComponentData() IComponentData {
	return &GenGigaSymbols2Data{
		cfg: genGigaSymbols2.Config,
	}
}

func NewGenGigaSymbols2(name string) IComponent {
	return &GenGigaSymbols2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "maxNumber": 0,
// "symbols": [
//
//	"WL",
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
// ],
// "spSymbols": [
//
//	"SC",
//	"EL",
//	"LW2",
//	"LW3",
//	"MY",
//	"RS",
//	"CS",
//	"C3",
//	"X3",
//	"B1",
//	"B2",
//	"WL_2",
//	"WL_3",
//	"LW2_2",
//	"LW2_3",
//	"LW3_2",
//	"LW3_3",
//	"MY_2",
//	"MY_3",
//	"RS_2",
//	"RS_3",
//	"EL_2"
//
// ],
// "weight": "bgspgridgigaweight"
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
type jsonGenGigaSymbols2 struct {
	MaxNumber     int      `json:"maxNumber"`
	Symbols       []string `json:"symbols"`
	SpSymbols     []string `json:"spSymbols"`
	Weight        string   `json:"weight"`
	BrokenSymbols []string `json:"brokenSymbols"`
}

func (jcfg *jsonGenGigaSymbols2) build() *GenGigaSymbols2Config {
	cfg := &GenGigaSymbols2Config{
		MaxNumber:     jcfg.MaxNumber,
		Symbols:       jcfg.Symbols,
		SpSymbols:     jcfg.SpSymbols,
		Weight:        jcfg.Weight,
		BrokenSymbols: slices.Clone(jcfg.BrokenSymbols),
	}

	return cfg
}

func parseGenGigaSymbols2(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenGigaSymbols2:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenGigaSymbols2:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenGigaSymbols2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenGigaSymbols2:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenGigaSymbols2TypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
