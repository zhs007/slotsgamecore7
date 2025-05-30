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

const GenPositionCollectionTypeName = "genPositionCollection"

// GenPositionCollectionConfig - configuration for GenPositionCollection
type GenPositionCollectionConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int    `yaml:"-" json:"-"`
	OutputToComponent    string   `yaml:"outputToComponent" json:"outputToComponent"`
}

// SetLinkComponent
func (cfg *GenPositionCollectionConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenPositionCollection struct {
	*BasicComponent `json:"-"`
	Config          *GenPositionCollectionConfig `json:"config"`
}

// Init -
func (genPositionCollection *GenPositionCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenPositionCollection.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenPositionCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenPositionCollection.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genPositionCollection.InitEx(cfg, pool)
}

// InitEx -
func (genPositionCollection *GenPositionCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	genPositionCollection.Config = cfg.(*GenPositionCollectionConfig)
	genPositionCollection.Config.ComponentType = GenPositionCollectionTypeName

	for _, s := range genPositionCollection.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenPositionCollection.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))
		}

		genPositionCollection.Config.SymbolCodes = append(genPositionCollection.Config.SymbolCodes, sc)
	}

	genPositionCollection.onInit(&genPositionCollection.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genPositionCollection *GenPositionCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// cd := icd.(*BasicComponentData)

	gs := genPositionCollection.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs != nil {
		pccd := gameProp.GetComponentDataWithName(genPositionCollection.Config.OutputToComponent)
		if pccd == nil {
			goutils.Error("GenPositionCollection.procReels:GetComponentDataWithName",
				slog.String("OutputToComponent", genPositionCollection.Config.OutputToComponent),
				goutils.Err(ErrNoComponent))

			return "", ErrNoComponent
		}

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(genPositionCollection.Config.SymbolCodes, s, 0) >= 0 {
					pccd.AddPos(x, y)
				}
			}
		}

		nc := genPositionCollection.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := genPositionCollection.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genPositionCollection *GenPositionCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewGenPositionCollection(name string) IComponent {
	return &GenPositionCollection{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "symbols": [
// 	"WL2",
// 	"WL3",
// 	"WL5"
// ],
// "outputToComponent": "fg-wlpos"

type jsonGenPositionCollection struct {
	Symbols           []string `json:"symbols"`
	OutputToComponent string   `json:"outputToComponent"`
}

func (jcfg *jsonGenPositionCollection) build() *GenPositionCollectionConfig {
	cfg := &GenPositionCollectionConfig{
		Symbols:           jcfg.Symbols,
		OutputToComponent: jcfg.OutputToComponent,
	}

	return cfg
}

func parseGenPositionCollection(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenPositionCollection:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenPositionCollection:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenPositionCollection{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenPositionCollection:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenPositionCollectionTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
