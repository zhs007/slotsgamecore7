package lowcode

import (
	"context"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &GenGigaSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenGigaSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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
				zap.String("symbol", genGigaSymbol.Config.Symbol),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		genGigaSymbol.Config.SymbolCode = sc
	}

	genGigaSymbol.onInit(&genGigaSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genGigaSymbol *GenGigaSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceReelWithMask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	cd.UsedScenes = nil

	gs := genGigaSymbol.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	w := gs.Width - genGigaSymbol.Config.GigaWidth + 1
	h := gs.Height - genGigaSymbol.Config.GigaHeight + 1

	if w <= 0 || h <= 0 {
		goutils.Error("GenGigaSymbol.OnPlayGame:Random:w",
			zap.Int("w", w),
			zap.Int("h", h),
			zap.Error(ErrIvalidComponentConfig))

		return "", ErrIvalidComponentConfig
	}

	if genGigaSymbol.Config.Number == 1 {
		x, err := plugin.Random(context.Background(), w)
		if err != nil {
			goutils.Error("GenGigaSymbol.OnPlayGame:Random:w",
				zap.Error(err))

			return "", err
		}

		y, err := plugin.Random(context.Background(), h)
		if err != nil {
			goutils.Error("GenGigaSymbol.OnPlayGame:Random:h",
				zap.Error(err))

			return "", err
		}

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
	} else {
		// 需要考虑相互不覆盖

		ngs = gs.CloneEx(gameProp.PoolScene)

		pos := make([]int, 0, w*h*2)
		npos := make([]int, 0, w*h*2)

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				pos = append(pos, x, y)
			}
		}

		for i := 0; i < genGigaSymbol.Config.Number; i++ {
			p, err := plugin.Random(context.Background(), len(pos)/2)
			if err != nil {
				goutils.Error("GenGigaSymbol.OnPlayGame:Random",
					zap.Error(err))

				return "", err
			}

			x := pos[p*2]
			y := pos[p*2+1]

			if i < genGigaSymbol.Config.Number-1 {
				for j := 0; j < len(pos)/2; j++ {
					if !(pos[j*2] >= x && pos[j*2+1] >= y && pos[j*2] < x+genGigaSymbol.Config.GigaWidth && pos[j*2+1] < y+genGigaSymbol.Config.GigaHeight) {
						npos = append(npos, pos[j*2], pos[j*2+1])
					}
				}

				tpos := pos
				pos = npos
				npos = tpos
			}
		}
	}

	if ngs == gs {
		nc := genGigaSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	genGigaSymbol.AddScene(gameProp, curpr, ngs, cd)

	nc := genGigaSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (genGigaSymbol *GenGigaSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after genGigaSymbol", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
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
	Type       string `json:"type"`
	GigaWidth  int    `json:"gigaWidth"`
	GigaHeight int    `json:"gigaHeight"`
	Number     int    `json:"number"`
	Symbol     string `json:"symbol"`
}

func (jcfg *jsonGenGigaSymbol) build() *GenGigaSymbolConfig {
	cfg := &GenGigaSymbolConfig{
		StrType:    jcfg.Type,
		GigaWidth:  jcfg.GigaWidth,
		GigaHeight: jcfg.GigaHeight,
		Number:     jcfg.Number,
		Symbol:     jcfg.Symbol,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseGenGigaSymbol(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenGigaSymbol:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenGigaSymbol:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonGenGigaSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenGigaSymbol:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenGigaSymbolTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
