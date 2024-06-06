package lowcode

import (
	"context"
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

const GenGigaSymbolTypeName = "genGigaSymbol"

type GenGigaSymbolType int

const (
	GGSTypeOverwrite GenGigaSymbolType = 0
	GGSTypeExpand    GenGigaSymbolType = 1
)

func parseGenGigaSymbolType(strType string) GenGigaSymbolType {
	if strType == "overwrite" {
		return GGSTypeOverwrite
	}

	return GGSTypeExpand
}

type GenGigaSymbolData struct {
	BasicComponentData
	Pos []int
}

// OnNewGame -
func (genGigaSymbolData *GenGigaSymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	genGigaSymbolData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (genGigaSymbolData *GenGigaSymbolData) OnNewStep(gameProp *GameProperty, component IComponent) {
	// positionCollectionData.BasicComponentData.OnNewStep(gameProp, component)
	genGigaSymbolData.UsedScenes = nil
	genGigaSymbolData.Pos = nil
}

// Clone
func (genGigaSymbolData *GenGigaSymbolData) Clone() IComponentData {
	target := &GenGigaSymbolData{
		BasicComponentData: genGigaSymbolData.CloneBasicComponentData(),
	}

	target.Pos = make([]int, len(genGigaSymbolData.Pos))
	copy(target.Pos, genGigaSymbolData.Pos)

	return target
}

// BuildPBComponentData
func (genGigaSymbolData *GenGigaSymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.GenGigaSymbolData{
		BasicComponentData: genGigaSymbolData.BuildPBBasicComponentData(),
	}

	for _, s := range genGigaSymbolData.Pos {
		pbcd.Pos = append(pbcd.Pos, int32(s))
	}

	return pbcd
}

// GetPos -
func (genGigaSymbolData *GenGigaSymbolData) GetPos() []int {
	return genGigaSymbolData.Pos
}

// HasPos -
func (genGigaSymbolData *GenGigaSymbolData) HasPos(x int, y int) bool {
	return goutils.IndexOfInt2Slice(genGigaSymbolData.Pos, x, y, 0) >= 0
}

// AddPos -
func (genGigaSymbolData *GenGigaSymbolData) AddPos(x int, y int) {
	genGigaSymbolData.Pos = append(genGigaSymbolData.Pos, x, y)
}

// GenGigaSymbolConfig - configuration for GenGigaSymbol
type GenGigaSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string            `yaml:"symbol" json:"type"`
	Type                 GenGigaSymbolType `yaml:"-" json:"-"`
	GigaWidth            int               `yaml:"gigaWidth" json:"gigaWidth"`
	GigaHeight           int               `yaml:"gigaHeight" json:"gigaHeight"`
	Number               int               `yaml:"number" json:"number"`
	Symbol               string            `yaml:"symbol" json:"symbol"`
	SymbolCode           int               `yaml:"-" json:"-"`
	ExcludeSymbols       []string          `yaml:"excludeSymbols" json:"excludeSymbols"`
	ExcludeSymbolCodes   []int             `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *GenGigaSymbolConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenGigaSymbol struct {
	*BasicComponent `json:"-"`
	Config          *GenGigaSymbolConfig `json:"config"`
}

// Init -
func (genGigaSymbol *GenGigaSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenGigaSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenGigaSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenGigaSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genGigaSymbol.InitEx(cfg, pool)
}

// InitEx -
func (genGigaSymbol *GenGigaSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	genGigaSymbol.Config = cfg.(*GenGigaSymbolConfig)
	genGigaSymbol.Config.ComponentType = GenGigaSymbolTypeName

	genGigaSymbol.Config.Type = parseGenGigaSymbolType(genGigaSymbol.Config.StrType)

	if genGigaSymbol.Config.Type == GGSTypeOverwrite {
		sc, isok := pool.DefaultPaytables.MapSymbols[genGigaSymbol.Config.Symbol]
		if !isok {
			goutils.Error("GenGigaSymbol.InitEx:Symbol",
				slog.String("symbol", genGigaSymbol.Config.Symbol),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbol.Config.SymbolCode = sc
	}

	genGigaSymbol.Config.ExcludeSymbolCodes = nil
	for _, v := range genGigaSymbol.Config.ExcludeSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("GenGigaSymbol.InitEx:ExcludeSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbol.Config.ExcludeSymbolCodes = append(genGigaSymbol.Config.ExcludeSymbolCodes, sc)
	}

	genGigaSymbol.onInit(&genGigaSymbol.Config.BasicComponentConfig)

	return nil
}

func (genGigaSymbol *GenGigaSymbol) isValidPos(gs *sgc7game.GameScene, x, y int) bool {
	for tx := 0; tx < genGigaSymbol.Config.GigaWidth; tx++ {
		for ty := 0; ty < genGigaSymbol.Config.GigaHeight; ty++ {
			if goutils.IndexOfIntSlice(genGigaSymbol.Config.ExcludeSymbolCodes, gs.Arr[x+tx][y+ty], 0) >= 0 {
				return false
			}
		}
	}

	return true
}

func (genGigaSymbol *GenGigaSymbol) genValidPos(gs *sgc7game.GameScene) ([]int, error) {
	lstPos := []int{}

	if len(genGigaSymbol.Config.ExcludeSymbolCodes) == 0 {
		cx := genGigaSymbol.Config.GigaWidth / 2
		cy := genGigaSymbol.Config.GigaHeight / 2

		for tx := cx; tx < gs.Width-(genGigaSymbol.Config.GigaWidth-cx); tx++ {
			for ty := cy; ty < gs.Height-(genGigaSymbol.Config.GigaHeight-cy); ty++ {
				lstPos = append(lstPos, tx, ty)
			}
		}

		return lstPos, nil
	}

	cx := genGigaSymbol.Config.GigaWidth / 2
	cy := genGigaSymbol.Config.GigaHeight / 2

	for tx := cx; tx < gs.Width-(genGigaSymbol.Config.GigaWidth-cx); tx++ {
		for ty := cy; ty < gs.Height-(genGigaSymbol.Config.GigaHeight-cy); ty++ {
			if genGigaSymbol.isValidPos(gs, tx, ty) {
				lstPos = append(lstPos, tx, ty)
			}
		}
	}

	return lstPos, nil
}

// playgame
func (genGigaSymbol *GenGigaSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceReelWithMask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	ggsd := icd.(*GenGigaSymbolData)

	ggsd.OnNewStep(gameProp, genGigaSymbol)

	gs := genGigaSymbol.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("GenGigaSymbol.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	ngs := gs

	lstpos, err := genGigaSymbol.genValidPos(gs)
	if err != nil {
		goutils.Error("GenGigaSymbol.OnPlayGame:genValidPos",
			goutils.Err(err))

		return "", err
	}

	if len(lstpos) == 0 {
		goutils.Error("GenGigaSymbol.OnPlayGame:lstpos",
			goutils.Err(ErrInvalidPosition))

		return "", ErrInvalidPosition
	}

	if genGigaSymbol.Config.Number == 1 {
		i, err := plugin.Random(context.Background(), len(lstpos)/2)
		if err != nil {
			goutils.Error("GenGigaSymbol.OnPlayGame:Random",
				goutils.Err(err))

			return "", err
		}

		x := lstpos[i*2]
		y := lstpos[i*2+1]

		ngs = gs.CloneEx(gameProp.PoolScene)

		s := genGigaSymbol.Config.SymbolCode
		if genGigaSymbol.Config.Type == GGSTypeExpand {
			s = ngs.Arr[x+genGigaSymbol.Config.GigaWidth/2][y+genGigaSymbol.Config.GigaHeight/2]
		}

		for tx := 0; tx < genGigaSymbol.Config.GigaWidth; tx++ {
			for ty := 0; ty < genGigaSymbol.Config.GigaHeight; ty++ {
				ngs.Arr[x+tx][y+ty] = s
			}
		}

		ggsd.AddPos(x+genGigaSymbol.Config.GigaWidth/2, y+genGigaSymbol.Config.GigaHeight/2)
	} else {
		// 需要考虑相互不覆盖

		// ngs = gs.CloneEx(gameProp.PoolScene)

		// pos := make([]int, 0, w*h*2)
		// npos := make([]int, 0, w*h*2)

		// for x := 0; x < w; x++ {
		// 	for y := 0; y < h; y++ {
		// 		pos = append(pos, x, y)
		// 	}
		// }

		// for i := 0; i < genGigaSymbol.Config.Number; i++ {
		// 	p, err := plugin.Random(context.Background(), len(pos)/2)
		// 	if err != nil {
		// 		goutils.Error("GenGigaSymbol.OnPlayGame:Random",
		// 			goutils.Err(err))

		// 		return "", err
		// 	}

		// 	x := pos[p*2]
		// 	y := pos[p*2+1]

		// 	if i < genGigaSymbol.Config.Number-1 {
		// 		for j := 0; j < len(pos)/2; j++ {
		// 			if !(pos[j*2] >= x && pos[j*2+1] >= y && pos[j*2] < x+genGigaSymbol.Config.GigaWidth && pos[j*2+1] < y+genGigaSymbol.Config.GigaHeight) {
		// 				npos = append(npos, pos[j*2], pos[j*2+1])
		// 			}
		// 		}

		// 		tpos := pos
		// 		pos = npos
		// 		npos = tpos
		// 	}
		// }
	}

	if ngs == gs {
		nc := genGigaSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	genGigaSymbol.AddScene(gameProp, curpr, ngs, &ggsd.BasicComponentData)

	nc := genGigaSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// NewComponentData -
func (genGigaSymbol *GenGigaSymbol) NewComponentData() IComponentData {
	return &GenGigaSymbolData{}
}

// OnAsciiGame - outpur to asciigame
func (genGigaSymbol *GenGigaSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	ggsd := icd.(*GenGigaSymbolData)

	if len(ggsd.UsedScenes) > 0 {
		asciigame.OutputScene("after genGigaSymbol", pr.Scenes[ggsd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (genGigaSymbol *GenGigaSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewGenGigaSymbol(name string) IComponent {
	return &GenGigaSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "expand",
// "gigaWidth": 3,
// "gigaHeight": 3,
// "number": 1
type jsonGenGigaSymbol struct {
	Type           string   `json:"type"`
	GigaWidth      int      `json:"gigaWidth"`
	GigaHeight     int      `json:"gigaHeight"`
	Number         int      `json:"number"`
	Symbol         string   `json:"symbol"`
	ExcludeSymbols []string `json:"excludeSymbols"`
}

func (jcfg *jsonGenGigaSymbol) build() *GenGigaSymbolConfig {
	cfg := &GenGigaSymbolConfig{
		StrType:        jcfg.Type,
		GigaWidth:      jcfg.GigaWidth,
		GigaHeight:     jcfg.GigaHeight,
		Number:         jcfg.Number,
		Symbol:         jcfg.Symbol,
		ExcludeSymbols: jcfg.ExcludeSymbols,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseGenGigaSymbol(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenGigaSymbol:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenGigaSymbol:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenGigaSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenGigaSymbol:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenGigaSymbolTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
