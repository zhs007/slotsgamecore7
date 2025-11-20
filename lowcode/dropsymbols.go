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

const DropSymbolsTypeName = "dropSymbols"

// DropSymbolsConfig - configuration for DropSymbols (placeholder)
type DropSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelIndex            int    `yaml:"reelIndex" json:"reelIndex"`
	Number               int    `yaml:"number" json:"number"`
	Symbol               string `yaml:"symbol" json:"symbol"`
	SymbolCode           int    `yaml:"-" json:"-"`
	JumpToComponent      string `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
}

func (cfg *DropSymbolsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

// DropSymbols - placeholder component
type DropSymbols struct {
	*BasicComponent `json:"-"`
	Config          *DropSymbolsConfig `json:"config"`
}

// Init - read yaml file
func (dropSymbols *DropSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("DropSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &DropSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("DropSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return dropSymbols.InitEx(cfg, pool)
}

// InitEx - initialize from config
func (dropSymbols *DropSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	dropSymbols.Config = cfg.(*DropSymbolsConfig)
	dropSymbols.Config.ComponentType = DropSymbolsTypeName

	symbolCode, isok := pool.DefaultPaytables.MapSymbols[dropSymbols.Config.Symbol]
	if isok {
		dropSymbols.Config.SymbolCode = symbolCode
	} else {
		goutils.Error("DropSymbols.InitEx:Symbol",
			slog.String("Symbol", dropSymbols.Config.Symbol),
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	dropSymbols.onInit(&dropSymbols.Config.BasicComponentConfig)

	return nil
}

// getReelIndex - helper
func (dropSymbols *DropSymbols) getReelIndex(_ *GameProperty, cd *BasicComponentData) int {
	v, isok := cd.GetConfigIntVal(CCVReelIndex)
	if isok {
		return v
	}

	return dropSymbols.Config.ReelIndex
}

// OnPlayGame - placeholder behavior: generate symbols and add to output symbol collection
func (dropSymbols *DropSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {
	bcd := icd.(*BasicComponentData)

	maxNumber := dropSymbols.Config.Number

	gs := dropSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("DropSymbols.OnPlayGame",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	reelindex := dropSymbols.getReelIndex(gameProp, bcd)
	if reelindex < 0 || reelindex >= gs.Width {
		// 这里表示不需要处理
		nc := dropSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	ngs := gs

	for y := ngs.Height - 1; y >= 0; y-- {
		if gs.Arr[reelindex][y] < 0 {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			ngs.Arr[reelindex][y] = dropSymbols.Config.SymbolCode
			maxNumber--

			if maxNumber <= 0 {
				break
			}
		}
	}

	if ngs != nil {
		dropSymbols.AddScene(gameProp, curpr, ngs, bcd)

		nc := dropSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := dropSymbols.onStepEnd(gameProp, curpr, gp, dropSymbols.Config.JumpToComponent)

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - output to asciigame (no-op placeholder)
func (dropSymbols *DropSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewDropSymbols(name string) IComponent {
	return &DropSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "ReelIndex": 0,
// "number": 1,
// "Symbol": "SC"
type jsonDropSymbols struct {
	ReelIndex int    `json:"ReelIndex"`
	Number    int    `json:"number"`
	Symbol    string `json:"Symbol"`
}

func (jcfg *jsonDropSymbols) build() *DropSymbolsConfig {
	cfg := &DropSymbolsConfig{
		ReelIndex: jcfg.ReelIndex,
		Number:    jcfg.Number,
		Symbol:    jcfg.Symbol,
	}

	return cfg
}

func parseDropSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseDropSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseDropSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonDropSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseDropSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: DropSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
