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

// SymbolExpanderTypeName is the registered component type name used in
// configuration and component registration for the symbol expander.
const SymbolExpanderTypeName = "symbolExpander"

// SymbolExpanderConfig holds configuration for the SymbolExpander component.
//
// Fields populated from YAML/JSON:
//   - Symbols: list of symbol names that should be treated as "expandable".
//   - IgnoreSymbols: list of symbol names that should be ignored when deciding
//     expansion boundaries. Typically a superset of Symbols is acceptable.
//
// Internal fields (SymbolCodes, IgnoreSymbolCodes) are filled during InitEx
// by mapping symbol names to integer codes from the game's paytables.
type SymbolExpanderConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string            `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int               `yaml:"-" json:"-"`
	IgnoreSymbols        []string            `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes    []int               `yaml:"-" json:"-"`
	JumpToComponent      string              `yaml:"jumpToComponent" json:"jumpToComponent"` // jump to
	MapAwards            map[string][]*Award `yaml:"controllers" json:"controllers"`
}

// SetLinkComponent
// SetLinkComponent implements BasicComponentConfig style linking for
// SymbolExpanderConfig. Supported links are:
//   - "next": set DefaultNextComponent
//   - "jump": set JumpToComponent
func (cfg *SymbolExpanderConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

// SymbolExpander is a lowcode component that expands specific symbols along a
// column (reel) in a scene. When expansion happens it may trigger configured
// controllers (MapAwards) and optionally jump to another component.
type SymbolExpander struct {
	*BasicComponent `json:"-"`
	Config          *SymbolExpanderConfig `json:"config"`
}

// Init loads YAML configuration from a file and delegates to InitEx. The
// configuration file is expected to match SymbolExpanderConfig.
func (se *SymbolExpander) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolExpander.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolExpanderConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolExpander.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return se.InitEx(cfg, pool)
}

// InitEx initializes the SymbolExpander from an already-unmarshaled config
// object. It performs the following steps:
//   - asserts cfg is a *SymbolExpanderConfig
//   - validates the provided GamePropertyPool and paytables
//   - converts Symbol/IgnoreSymbol names to internal integer codes
//   - initializes award controllers if present
//
// Returns an error when configuration or pool data is invalid.
func (se *SymbolExpander) InitEx(cfg any, pool *GamePropertyPool) error {
	// type assertion with safety
	scfg, ok := cfg.(*SymbolExpanderConfig)
	if !ok || scfg == nil {
		goutils.Error("SymbolExpander.InitEx:invalid config",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	if pool == nil || pool.DefaultPaytables == nil || pool.DefaultPaytables.MapSymbols == nil {
		goutils.Error("SymbolExpander.InitEx:invalid pool",
			goutils.Err(ErrInvalidGameData))

		return ErrInvalidGameData
	}

	se.Config = scfg
	se.Config.ComponentType = SymbolExpanderTypeName

	// build symbol codes: map symbol names to integer codes
	for _, s := range se.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolExpander.InitEx:Symbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		se.Config.SymbolCodes = append(se.Config.SymbolCodes, sc)
	}

	// build ignore symbol codes: map ignore symbol names to integer codes
	for _, s := range se.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolExpander.InitEx:IgnoreSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		se.Config.IgnoreSymbolCodes = append(se.Config.IgnoreSymbolCodes, sc)
	}

	// initialize awards defensively: call Award.Init() for each non-nil award
	if se.Config.MapAwards != nil {
		for k, award := range se.Config.MapAwards {
			if award == nil {
				// skip empty award list but log for context
				goutils.Debug("SymbolExpander.InitEx:empty award list",
					slog.String("controller", k))
				continue
			}

			for _, a := range award {
				if a == nil {
					goutils.Debug("SymbolExpander.InitEx:nil award in list",
						slog.String("controller", k))
					continue
				}

				a.Init()
			}
		}
	}

	// call lifecycle hook defined on BasicComponent
	se.onInit(&se.Config.BasicComponentConfig)

	return nil
}

// ProcControllers executes controllers configured for the provided string
// value (strVal). It looks up awards in the component's MapAwards and calls
// gameProp.procAwards to apply them.
func (se *SymbolExpander) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, pr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if se.Config == nil || se.Config.MapAwards == nil {
		return
	}

	awards, isok := se.Config.MapAwards[strVal]
	if !isok || awards == nil {
		return
	}

	gameProp.procAwards(plugin, awards, pr, gp)
}

// OnPlayGame performs the expansion logic on the target scene.
//
// Behavior summary:
//   - finds the target scene via GetTargetScene3
//   - for each column (axis) checks if the column contains any expandable
//     symbol codes (SymbolCodes)
//   - determines the first row (starty) from top that is neither expandable
//     nor ignored; from that row downwards replaces symbols with the
//     expandable symbol code (sc)
//   - clones the scene lazily (only when a change is required)
//   - records which symbol cores were expanded and triggers controllers
//   - adds the new scene to PlayResult and returns the next component to run
func (se *SymbolExpander) OnPlayGame(gameProp *GameProperty, pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// validate component data type to avoid panics on wrong usage
	cd, ok := icd.(*BasicComponentData)
	if !ok || cd == nil {
		goutils.Error("SymbolExpander.OnPlayGame:invalid icd type",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	cd.UsedScenes = nil
	lstSymbolCores := make([]int, 0, len(se.Config.SymbolCodes))

	gs := se.GetTargetScene3(gameProp, pr, prs, 0)
	if gs == nil {
		nc := se.onStepEnd(gameProp, pr, gp, "")

		return nc, ErrComponentDoNothing
	}

	ngs := gs

	for x, arr := range gs.Arr {
		// First: determine whether this column contains any expandable symbol
		sc := -1
		for _, s := range arr {
			if slices.Contains(se.Config.SymbolCodes, s) {
				sc = s
				break
			}
		}

		if sc == -1 {
			// no expandable symbol in this column
			continue
		}

		// Second: find the first row index from top which is neither an
		// expandable symbol nor an ignored symbol. That index is the start of
		// the expansion region. If every cell is expandable/ignored, there is
		// nothing to expand.
		starty := -1
		for y, s := range arr {
			if !slices.Contains(se.Config.SymbolCodes, s) && !slices.Contains(se.Config.IgnoreSymbolCodes, s) {
				starty = y

				break
			}
		}

		if starty != -1 {
			// column needs expansion

			// record symbol core for controller triggers (avoid duplicates)
			if !slices.Contains(lstSymbolCores, sc) {
				lstSymbolCores = append(lstSymbolCores, sc)
			}

			// lazy clone: clone the scene only once when the first change is made
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			// apply expansion: from starty downwards, replace cells that are
			// neither expandable nor ignored with the core symbol sc
			for ty := starty; ty < len(arr); ty++ {
				if ty == starty {
					ngs.Arr[x][ty] = sc

					continue
				}
				// skip cells that are already expandable or explicitly ignored
				s := arr[ty]

				if !slices.Contains(se.Config.SymbolCodes, s) && !slices.Contains(se.Config.IgnoreSymbolCodes, s) {
					ngs.Arr[x][ty] = sc
				}
			}
		}
	}

	if len(lstSymbolCores) > 0 {
		se.ProcControllers(gameProp, plugin, pr, gp, 0, "<trigger>")

		for _, sc := range lstSymbolCores {
			se.ProcControllers(gameProp, plugin, pr, gp, 0, gameProp.Pool.DefaultPaytables.GetStringFromInt(sc))
		}
	}

	if ngs == gs {
		nc := se.onStepEnd(gameProp, pr, gp, "")

		return nc, ErrComponentDoNothing
	}

	se.AddScene(gameProp, pr, ngs, cd)

	nc := se.onStepEnd(gameProp, pr, gp, se.Config.JumpToComponent)

	return nc, nil
}

// OnAsciiGame outputs an ASCII representation of the scene after the
// expander runs. Only the first used scene in component data is printed. This
// is primarily a debugging helper and will silently return if the provided
// component data is of an unexpected type.
func (se *SymbolExpander) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd, ok := icd.(*BasicComponentData)
	if !ok || cd == nil {
		// nothing to output
		return nil
	}

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after symbolExpander", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// NewSymbolExpander constructs a new SymbolExpander component instance with
// the provided name.
func NewSymbolExpander(name string) IComponent {
	return &SymbolExpander{
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
// jsonSymbolExpander represents the minimal JSON shape used when parsing
// an inline JSON block for the expander in lowcode scripts.
type jsonSymbolExpander struct {
	Symbols       []string `json:"symbols"`
	IgnoreSymbols []string `json:"ignoreSymbols"`
}

// build converts the JSON helper type into a full SymbolExpanderConfig. It
// clones slice values to avoid accidental sharing of backing arrays.
func (jcfg *jsonSymbolExpander) build() *SymbolExpanderConfig {
	cfg := &SymbolExpanderConfig{
		Symbols:       slices.Clone(jcfg.Symbols),
		IgnoreSymbols: slices.Clone(jcfg.IgnoreSymbols),
	}

	return cfg
}

// parseSymbolExpander parses a lowcode JSON/AST cell and registers the
// resulting SymbolExpanderConfig into the provided BetConfig. It returns the
// assigned label name for the component.
//
// The function expects the AST cell to provide a JSON object and optional
// controller blocks. The controllers are parsed and attached to cfgd.MapAwards.
func parseSymbolExpander(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolExpander:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolExpander:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolExpander{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolExpander:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapAwards, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseSymbolExpander:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapAwards = mapAwards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolExpanderTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
