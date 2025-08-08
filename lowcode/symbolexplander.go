package lowcode

import (
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

const SymbolExplanderTypeName = "symbolExplander"

// SymbolExplanderConfig - configuration for SymbolExplander
type SymbolExplanderConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string            `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int               `yaml:"-" json:"-"`
	JumpToComponent      string              `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
	MapAwards            map[string][]*Award `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
func (cfg *SymbolExplanderConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type SymbolExplander struct {
	*BasicComponent `json:"-"`
	Config          *SymbolExplanderConfig `json:"config"`
}

// Init -
func (symbolExplander *SymbolExplander) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolExplander.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolExplanderConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolExplander.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return symbolExplander.InitEx(cfg, pool)
}

// InitEx -
func (symbolExplander *SymbolExplander) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolExplander.Config = cfg.(*SymbolExplanderConfig)
	symbolExplander.Config.ComponentType = SymbolExplanderTypeName

	for _, s := range symbolExplander.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolExplander.InitEx:Symbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		symbolExplander.Config.SymbolCodes = append(symbolExplander.Config.SymbolCodes, sc)
	}

	for _, award := range symbolExplander.Config.MapAwards {
		for _, a := range award {
			a.Init()
		}
	}

	symbolExplander.onInit(&symbolExplander.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (symbolExplander *SymbolExplander) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	awards, isok := symbolExplander.Config.MapAwards[strVal]
	if isok {
		gameProp.procAwards(plugin, awards, curpr, gp)
	}
}

// playgame
func (symbolExplander *SymbolExplander) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*BasicComponentData)

	cd.UsedScenes = nil
	lstSymbolCores := make([]int, 0, len(symbolExplander.Config.SymbolCodes))

	gs := symbolExplander.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	for x, arr := range gs.Arr {
		for _, s := range arr {
			if goutils.IndexOfIntSlice(symbolExplander.Config.SymbolCodes, s, 0) >= 0 {
				if !slices.Contains(lstSymbolCores, s) {
					lstSymbolCores = append(lstSymbolCores, s)
				}

				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				for ty := range ngs.Height {
					ngs.Arr[x][ty] = s
				}
			}
		}
	}

	if len(lstSymbolCores) > 0 {
		symbolExplander.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")

		for _, sc := range lstSymbolCores {
			symbolExplander.ProcControllers(gameProp, plugin, curpr, gp, 0, gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))
		}
	}

	if ngs == gs {
		nc := symbolExplander.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	symbolExplander.AddScene(gameProp, curpr, ngs, cd)

	nc := symbolExplander.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (symbolExplander *SymbolExplander) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after symbolExplander", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

func NewSymbolExplander(name string) IComponent {
	return &SymbolExplander{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "symbols": [
//
//	"CO",
//	"CH",
//	"CN",
//	"CB"
//
// ]
type jsonSymbolExplander struct {
	Symbols []string `json:"symbols"`
}

func (jcfg *jsonSymbolExplander) build() *SymbolExplanderConfig {
	cfg := &SymbolExplanderConfig{
		Symbols: slices.Clone(jcfg.Symbols),
	}

	return cfg
}

func parseSymbolExplander(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolExplander:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolExplander:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolExplander{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolExplander:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseSymbolValsSP:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolExplanderTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
