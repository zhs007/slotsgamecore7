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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const GenTropiCoolSPSymbolsTypeName = "genTropiCoolSPSymbols"

type TropiCoolSPData struct {
	SymbolCode    int
	CurSymbolCode int
	Height        int
	X             int
	Y             int
}

type GenTropiCoolSPSymbolsData struct {
	BasicComponentData
	tropiCoolSPData []*TropiCoolSPData
	cfg             *GenTropiCoolSPSymbolsConfig
}

// OnNewGame -
func (gcd *GenTropiCoolSPSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	gcd.BasicComponentData.OnNewGame(gameProp, component)

	gcd.tropiCoolSPData = nil
}

// OnNewStep -
func (gcd *GenTropiCoolSPSymbolsData) onNewStep() {
	gcd.UsedScenes = nil
}

// Clone
func (gcd *GenTropiCoolSPSymbolsData) Clone() IComponentData {
	target := &GenTropiCoolSPSymbolsData{
		BasicComponentData: gcd.CloneBasicComponentData(),
	}

	for _, v := range gcd.tropiCoolSPData {
		target.tropiCoolSPData = append(target.tropiCoolSPData, &TropiCoolSPData{
			SymbolCode:    v.SymbolCode,
			CurSymbolCode: v.CurSymbolCode,
			Height:        v.Height,
			X:             v.X,
			Y:             v.Y,
		})
	}

	return target
}

// BuildPBComponentData
func (gcd *GenTropiCoolSPSymbolsData) BuildPBComponentData() proto.Message {
	return &sgc7pb.BasicComponentData{
		BasicComponentData: gcd.BuildPBBasicComponentData(),
	}
}

// GenTropiCoolSPSymbolsConfig - placeholder configuration
type GenTropiCoolSPSymbolsConfig struct {
	BasicComponentConfig  `yaml:",inline" json:",inline"`
	SpBonusSymbol         string                `yaml:"spBonusSymbol" json:"spBonusSymbol"`
	SpBonusSymbolCode     int                   `yaml:"-" json:"-"`
	SpSymbols             []string              `yaml:"SpSymbols" json:"SpSymbols"`
	SpSymbolCodes         []int                 `yaml:"-" json:"-"`
	GenGigaSymbols2       string                `yaml:"genGigaSymbols2" json:"genGigaSymbols2"`
	Symbols               []string              `yaml:"symbols" json:"symbols"`
	SymbolCodes           []int                 `yaml:"-" json:"-"`
	GenTropiCoolSPSymbols string                `yaml:"genTropiCoolSPSymbols" json:"genTropiCoolSPSymbols"`
	SymbolWeight          string                `yaml:"symbolWeight" json:"symbolWeight"`
	SymbolWeightVW        *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	BlankSymbol           string                `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode       int                   `yaml:"-" json:"-"`
	MapSpSymbolCodes      map[int][]int         `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *GenTropiCoolSPSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

func (cfg *GenTropiCoolSPSymbolsConfig) getSpSymbolInfo(sc int) (int, int) {
	for k, arr := range cfg.MapSpSymbolCodes {
		i := slices.Index(arr, sc)
		if i >= 0 {
			return k, i + 1
		}
	}

	return -1, -1
}

type GenTropiCoolSPSymbols struct {
	*BasicComponent `json:"-"`
	Config          *GenTropiCoolSPSymbolsConfig `json:"config"`
}

// Init - load from file
func (gen *GenTropiCoolSPSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenTropiCoolSPSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenTropiCoolSPSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenTropiCoolSPSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gen.InitEx(cfg, pool)
}

// InitEx - initialize from config object
func (gen *GenTropiCoolSPSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	gen.Config = cfg.(*GenTropiCoolSPSymbolsConfig)
	gen.Config.ComponentType = GenTropiCoolSPSymbolsTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[gen.Config.SpBonusSymbol]
	if !isok {
		goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
			slog.String("SpBonusSymbol", gen.Config.SpBonusSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.SpBonusSymbolCode = sc

	for _, s := range gen.Config.SpSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
				slog.String("SpSymbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gen.Config.SpSymbolCodes = append(gen.Config.SpSymbolCodes, sc)
	}

	for _, s := range gen.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
				slog.String("Symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		gen.Config.SymbolCodes = append(gen.Config.SymbolCodes, sc)
	}

	sc, isok = pool.DefaultPaytables.MapSymbols[gen.Config.BlankSymbol]
	if !isok {
		goutils.Error("GenTropiCoolSPSymbols.InitEx:Symbol",
			slog.String("BlankSymbol", gen.Config.BlankSymbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	gen.Config.BlankSymbolCode = sc

	vw, err := pool.LoadIntWeights(gen.Config.SymbolWeight, true)
	if err != nil {
		goutils.Error("GenTropiCoolSPSymbols.InitEx:LoadIntWeights",
			slog.String("SymbolWeight", gen.Config.SymbolWeight),
			goutils.Err(err))

		return err
	}

	gen.Config.SymbolWeightVW = vw

	gen.Config.MapSpSymbolCodes = make(map[int][]int)
	for i, s := range gen.Config.SpSymbols {
		sc = gen.Config.SpSymbolCodes[i]
		gen.Config.MapSpSymbolCodes[sc] = make([]int, pool.Config.Height)

		for i := 2; i <= pool.Config.Height; i++ {
			gen.Config.MapSpSymbolCodes[sc][0] = sc

			nsc, isok := pool.DefaultPaytables.MapSymbols[fmt.Sprintf("%v_%v", s, i)]
			if isok {
				gen.Config.MapSpSymbolCodes[sc][i-1] = nsc
			} else {
				gen.Config.MapSpSymbolCodes[sc][i-1] = -1
			}
		}
	}

	gen.onInit(&gen.Config.BasicComponentConfig)

	return nil
}

func (gen *GenTropiCoolSPSymbols) getWeight(gameProp *GameProperty, bcd *BasicComponentData) *sgc7game.ValWeights2 {
	str := bcd.GetConfigVal(CCVSymbolWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return gen.Config.SymbolWeightVW
}

func (gen *GenTropiCoolSPSymbols) isCanSpSymbol(x, y int, h int, gs *sgc7game.GameScene) bool {

	for gy := 0; gy <= h-1; gy++ {
		if gy == 0 {
			continue
		}

		if y+gy >= gs.Height {
			return false
		}

		if !slices.Contains(gen.Config.SymbolCodes, gs.Arr[x][y+gy]) {
			return false
		}
	}

	return true
}

func (gen *GenTropiCoolSPSymbols) procSpSymbol(gs *sgc7game.GameScene, vw *sgc7game.ValWeights2, gameProp *GameProperty,
	plugin sgc7plugin.IPlugin, cd *GenTropiCoolSPSymbolsData) (*sgc7game.GameScene, error) {

	ngs := gs

	for x := 0; x < gs.Width-1; x++ {
		for y := 0; y < gs.Height-1; y++ {

			symbolCode := ngs.Arr[x][y]

			if slices.Contains(gen.Config.SymbolCodes, symbolCode) {
				cr, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("GenTropiCoolSPSymbols.procSpSymbol:RandVal",
						goutils.Err(err))

					return nil, err
				}

				spsc := cr.Int()

				if spsc == gen.Config.BlankSymbolCode {
					continue
				}

				msc, h := gen.Config.getSpSymbolInfo(spsc)
				if msc < 0 {
					goutils.Error("GenTropiCoolSPSymbols.procSpSymbol:getSpSymbolInfo",
						goutils.Err(ErrInvalidComponentData))

					return nil, ErrInvalidComponentData
				}

				if gen.isCanSpSymbol(x, y, h, ngs) {
					if ngs == gs {
						ngs = gs.CloneEx(gameProp.PoolScene)
					}

					for gy := 0; gy < h; gy++ {
						ngs.Arr[x][y+gy] = gen.Config.MapSpSymbolCodes[msc][h-1]
					}

					gd := &TropiCoolSPData{
						SymbolCode:    msc,
						CurSymbolCode: gen.Config.MapSpSymbolCodes[msc][h-1],
						Height:        cr.Int(),
						X:             x,
						Y:             y,
					}

					cd.tropiCoolSPData = append(cd.tropiCoolSPData, gd)
				}

			}
		}
	}

	return ngs, nil
}

// OnPlayGame - placeholder: do nothing
func (gen *GenTropiCoolSPSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	gcd, isok := icd.(*GenTropiCoolSPSymbolsData)
	if !isok {
		goutils.Error("GenTropiCoolSPSymbols.OnPlayGame:GenTropiCoolSPSymbolsData",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	vw := gen.getWeight(gameProp, &gcd.BasicComponentData)
	if vw == nil {
		goutils.Error("GenTropiCoolSPSymbols.OnPlayGame:weight is nil",
			slog.String("component", gen.GetName()),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	gs := gen.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("GenTropiCoolSPSymbols.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	gcd.onNewStep()

	ngs, err := gen.procSpSymbol(gs, vw, gameProp, plugin, gcd)
	if err != nil {
		goutils.Error("GenTropiCoolSPSymbols.OnPlayGame:procSpSymbol",
			goutils.Err(err))

		return "", err
	}

	if ngs != gs {
		gen.AddScene(gameProp, curpr, ngs, &gcd.BasicComponentData)

		gen.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

		nc := gen.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := gen.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - no-op
func (gen *GenTropiCoolSPSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData - return base component data
func (gen *GenTropiCoolSPSymbols) NewComponentData() IComponentData {
	return &GenTropiCoolSPSymbolsData{
		cfg: gen.Config,
	}
}

func NewGenTropiCoolSPSymbols(name string) IComponent {
	return &GenTropiCoolSPSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "spBonusSymbol": "EL",
// "spSymbols": [
//
//	"B1",
//	"B2"
//
// ],
// "srcGenGigaSymbols2": "bg-gengiga",
// "genGigaSymbols2": "bg-gengiga",
// "symbols": [
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
// ],
// "genTropiCoolSPSymbols": "bg-regenspsym",
// "symbolWeight": "spsymbolsweight",
// "blankSymbol": "BN"
type jsonGenTropiCoolSPSymbols struct {
	SpBonusSymbol         string   `json:"spBonusSymbol"`
	SpSymbols             []string `json:"spSymbols"`
	GenGigaSymbols2       string   `json:"genGigaSymbols2"`
	Symbols               []string `json:"symbols"`
	GenTropiCoolSPSymbols string   `json:"genTropiCoolSPSymbols"`
	SymbolWeight          string   `json:"symbolWeight"`
	BlankSymbol           string   `json:"blankSymbol"`
}

func (j *jsonGenTropiCoolSPSymbols) build() *GenTropiCoolSPSymbolsConfig {
	return &GenTropiCoolSPSymbolsConfig{
		SpBonusSymbol:         j.SpBonusSymbol,
		SpSymbols:             j.SpSymbols,
		GenGigaSymbols2:       j.GenGigaSymbols2,
		Symbols:               j.Symbols,
		GenTropiCoolSPSymbols: j.GenTropiCoolSPSymbols,
		SymbolWeight:          j.SymbolWeight,
		BlankSymbol:           j.BlankSymbol,
	}
}

func parseGenTropiCoolSPSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenTropiCoolSPSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenTropiCoolSPSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenTropiCoolSPSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenTropiCoolSPSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenTropiCoolSPSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
