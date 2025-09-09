package lowcode

import (
	"fmt"
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

// WinResultModifierExTypeName is the component type name used in configuration
// and registration for the WinResultModifierEx component.
const WinResultModifierExTypeName = "winResultModifierEx"

// WinResultModifierExData holds runtime state for a WinResultModifierEx component.
//
// It embeds BasicComponentData and tracks accumulated wins and the applied
// win multiplier for the current step.
type WinResultModifierExData struct {
	BasicComponentData
	Wins     int
	WinMulti int
}

// OnNewGame initializes per-game state for WinResultModifierExData.
//
// OnNewGame forwards to BasicComponentData.OnNewGame and performs any
// WinResultModifierEx-specific initialization required at the start of a new game.
func (winResultModifierDataEx *WinResultModifierExData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultModifierDataEx.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep resets step-specific fields on WinResultModifierExData.
//
// onNewStep resets accumulated wins and the step multiplier to their default
// values. This is an internal helper used at the beginning of each step.
func (winResultModifierDataEx *WinResultModifierExData) onNewStep() {
	winResultModifierDataEx.Wins = 0
	winResultModifierDataEx.WinMulti = 1
}

// Clone creates a deep copy of WinResultModifierExData and returns it as
// the IComponentData interface.
//
// Clone is used when component data must be duplicated (for example, when
// creating snapshots or copying state between game contexts).
func (winResultModifierDataEx *WinResultModifierExData) Clone() IComponentData {
	target := &WinResultModifierExData{
		BasicComponentData: winResultModifierDataEx.CloneBasicComponentData(),
		Wins:               winResultModifierDataEx.Wins,
		WinMulti:           winResultModifierDataEx.WinMulti,
	}

	return target
}

// BuildPBComponentData converts the component data into its protobuf
// representation used for telemetry or RPC.
//
// BuildPBComponentData returns a proto.Message representing the current
// WinResultModifierExData.
func (winResultModifierDataEx *WinResultModifierExData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultModifierData{
		BasicComponentData: winResultModifierDataEx.BuildPBBasicComponentData(),
		Wins:               int32(winResultModifierDataEx.Wins),
		WinMulti:           int32(winResultModifierDataEx.WinMulti),
	}

	return pbcd
}

// GetValEx retrieves a named integer value from the component data.
//
// For WinResultModifierExData the supported key is "CVWins" which returns the
// accumulated wins. The function returns the value and a boolean indicating
// whether the key was found.
func (winResultModifierDataEx *WinResultModifierExData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultModifierDataEx.Wins, true
	}

	return 0, false
}

// WinResultModifierExConfig - configuration for WinResultModifierEx
//
// WinResultModifierExConfig describes how a WinResultModifierEx component is
// configured via YAML/JSON. It includes the component base config, the
// type string, source components that will be inspected, and a mapping from
// target symbol names to multiplier values.
type WinResultModifierExConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`                         // type
	Type                 WinResultModifierType `yaml:"-" json:"-"`                               // type
	SourceComponents     []string              `yaml:"sourceComponents" json:"sourceComponents"` // target components
	MapTargetSymbols     map[string]int        `yaml:"mapTargetSymbols" json:"mapTargetSymbols"` // mapTargetSymbols
	MapTargetSymbolCodes map[int]int           `yaml:"-" json:"-"`                               // MapTargetSymbolCodes
}

// SetLinkComponent
// SetLinkComponent sets link-based connections for the component.
//
// The function supports the "next" link which sets the default next
// component name. This mirrors the behavior used by other components and is
// called by the builder when wiring components together.
func (cfg *WinResultModifierExConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// WinResultModifierEx is a component that modifies win results based on
// configured symbol mappings and multiplier logic.
//
// It embeds BasicComponent and holds a typed configuration pointer.
type WinResultModifierEx struct {
	*BasicComponent `json:"-"`
	Config          *WinResultModifierExConfig `json:"config"`
}

// Init loads YAML configuration from the given filename and initializes the
// component. It is a convenience wrapper around InitEx that reads the file
// content and unmarshals it into the component config structure.
func (winResultModifierEx *WinResultModifierEx) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WinResultModifierEx.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WinResultModifierExConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultModifierEx.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return winResultModifierEx.InitEx(cfg, pool)
}

// InitEx initializes the component from an already-unmarshaled configuration
// object (typically *WinResultModifierExConfig). It validates configuration
// values and resolves symbol codes from the paytable pool.
func (winResultModifierEx *WinResultModifierEx) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultModifierEx.Config = cfg.(*WinResultModifierExConfig)
	winResultModifierEx.Config.ComponentType = WinResultModifierExTypeName

	winResultModifierEx.Config.Type = parseWinResultModifierType(winResultModifierEx.Config.StrType)
	if !winResultModifierEx.Config.Type.isValidInWinResultModifierEx() {

		goutils.Error("WinResultModifierEx.InitEx:isValidInWinResultModifierEx",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	winResultModifierEx.Config.MapTargetSymbolCodes = make(map[int]int)

	for k, v := range winResultModifierEx.Config.MapTargetSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[k]
		if !isok {
			goutils.Error("WinResultModifierEx.InitEx:MapTargetSymbols.Symbol",
				slog.String("symbol", k),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		winResultModifierEx.Config.MapTargetSymbolCodes[sc] = v
	}

	winResultModifierEx.onInit(&winResultModifierEx.Config.BasicComponentConfig)

	return nil
}

// OnPlayGame applies the configured win-modification logic to the current
// PlayResult. It inspects results produced by configured source components
// and adjusts CoinWin, CashWin and OtherMul according to the mapping and
// modifier type.
//
// OnPlayGame returns the name of the next component (or empty) and an error
// indicating whether any modification was applied. If no modification took
// place it returns ErrComponentDoNothing.
func (winResultModifierEx *WinResultModifierEx) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	std, ok := icd.(*WinResultModifierExData)
	if !ok {
		goutils.Error("WinResultModifierEx.OnPlayGame:invalid icd type",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	std.onNewStep()

	gs := winResultModifierEx.GetTargetScene3(gameProp, curpr, prs, 0)
	isproced := false

	for _, cn := range winResultModifierEx.Config.SourceComponents {
		// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
		if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
			continue
		}

		ccd := gameProp.GetComponentDataWithName(cn)
		// ccd := gameProp.MapComponentData[cn]
		lst := ccd.GetResults()
		for _, ri := range lst {
			mul := CalcSymbolsInResultEx(gs, winResultModifierEx.Config.MapTargetSymbolCodes, curpr.Results[ri], winResultModifierEx.Config.Type)

			if mul > 1 {
				if winResultModifierEx.Config.Type == WRMTypeSymbolMultiOnWays {
					// protect against division by zero
					if curpr.Results[ri].Mul <= 0 {
						goutils.Error("WinResultModifierEx.OnPlayGame:curpr.Results[ri].Mul <= 0",
							goutils.Err(ErrInvalidComponentConfig))

						return "", ErrInvalidComponentConfig
					}

					curpr.Results[ri].OtherMul = mul

					curpr.Results[ri].CoinWin = curpr.Results[ri].CoinWin / curpr.Results[ri].Mul * mul
					curpr.Results[ri].CashWin = curpr.Results[ri].CashWin / curpr.Results[ri].Mul * mul
				} else {
					curpr.Results[ri].CashWin *= mul
					curpr.Results[ri].CoinWin *= mul
					curpr.Results[ri].OtherMul *= mul
				}

				std.Wins += curpr.Results[ri].CoinWin

				isproced = true
			}

		}
	}

	if !isproced {
		nc := winResultModifierEx.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := winResultModifierEx.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame outputs component debug information to the asciigame system.
//
// OnAsciiGame is intended for human-readable debugging of component behavior
// and prints the current multiplier and accumulated wins.
func (winResultModifierEx *WinResultModifierEx) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std, ok := icd.(*WinResultModifierExData)
	if !ok {
		goutils.Error("WinResultModifierEx.OnAsciiGame:invalid icd type",
			goutils.Err(ErrInvalidComponentData))

		return ErrInvalidComponentData
	}

	fmt.Printf("WinResultModifierEx x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// NewComponentData creates a fresh instance of WinResultModifierExData used
// to track per-game and per-step state for this component.
func (winResultModifierEx *WinResultModifierEx) NewComponentData() IComponentData {
	return &WinResultModifierExData{}
}

// NewWinResultModifierEx creates a new WinResultModifierEx instance with the
// provided name and default priority.
func NewWinResultModifierEx(name string) IComponent {
	return &WinResultModifierEx{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "addSymbolMulti",
// "mapTargetSymbols": [
//
//	[
//		"WL2",
//		2
//	],
//	[
//		"WL3",
//		3
//	],
//	[
//		"WL5",
//		5
//	]
//
// ],
// "sourceComponent": [
//
//	"bg-wins"
//
// ]
// jsonWinResultModifierEx is an intermediate structure for parsing compact
// JSON configuration when the component is defined inline in a script or
// when using the ast-based parser. It mirrors the minimal fields expected in
// JSON source form and is converted to a WinResultModifierExConfig by
// build().
//
// The JSON shape is expected as:
//
//	{
//	  "type": "addSymbolMulti",
//	  "sourceComponent": [ ... ],
//	  "mapTargetSymbols": [ ["SYM", 2], ["SYM2", 3] ]
//	}
type jsonWinResultModifierEx struct {
	Type             string   `json:"type"`             // type
	SourceComponents []string `json:"sourceComponent"`  // source components
	MapTargetSymbols [][]any  `json:"mapTargetSymbols"` // mapTargetSymbols
}

func (jcfg *jsonWinResultModifierEx) build() *WinResultModifierExConfig {
	cfg := &WinResultModifierExConfig{
		StrType:          jcfg.Type,
		SourceComponents: jcfg.SourceComponents,
		MapTargetSymbols: make(map[string]int),
	}

	for _, arr := range jcfg.MapTargetSymbols {
		if len(arr) < 2 {
			continue
		}

		// key should be string
		key, ok := arr[0].(string)
		if !ok {
			continue
		}

		// value can be float64 (default for numbers in sonic.Unmarshal) or int
		switch v := arr[1].(type) {
		case float64:
			cfg.MapTargetSymbols[key] = int(v)
		case int:
			cfg.MapTargetSymbols[key] = v
		case int64:
			cfg.MapTargetSymbols[key] = int(v)
		default:
			// try to handle numeric-like strings
			if s, ok := arr[1].(string); ok {
				// best-effort parse
				var vi int
				_, err := fmt.Sscanf(s, "%d", &vi)
				if err == nil {
					cfg.MapTargetSymbols[key] = vi
				}
			}
		}
	}

	return cfg
}

func parseWinResultModifierEx(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	// parseWinResultModifierEx parses an AST cell produced by the low-code
	// parser and registers a WinResultModifierEx component into the
	// provided BetConfig. It extracts the inline JSON-like config, maps it
	// into the typed config, and appends a ComponentConfig entry.
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultModifierEx:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultModifierEx:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWinResultModifierEx{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultModifierEx:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultModifierExTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
