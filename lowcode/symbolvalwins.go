package lowcode

import (
	"log/slog"
	"os"
	"slices"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

// SymbolValWinsTypeName is the component type name used in configuration
// and registration for the SymbolValWins component.
const SymbolValWinsTypeName = "symbolValWins"

// SymbolValWinsType enumerates the variant types for SymbolValWins behavior.
// It controls how collector symbols are interpreted when computing wins.
type SymbolValWinsType int

const (
	svwTypeNormal        SymbolValWinsType = 0
	svwTypeCollector     SymbolValWinsType = 1
	svwTypeReelCollector SymbolValWinsType = 2
)

// parseSymbolValWinsType converts a textual type (from configuration)
// into the corresponding SymbolValWinsType enum value.
func parseSymbolValWinsType(strType string) SymbolValWinsType {
	switch strType {
	case "collector":
		return svwTypeCollector
	case "reelcollector":
		return svwTypeReelCollector
	}

	return svwTypeNormal
}

const (
	SVWDVWins      string = "wins"      // 中奖的数值，线注的倍数
	SVWDVSymbolNum string = "symbolNum" // 符号数量
)

// SymbolValWinsData stores runtime state for a SymbolValWins component.
// It extends BasicComponentData and records the number of collected symbols
// and the total coin wins produced by the component during a step/session.
type SymbolValWinsData struct {
	BasicComponentData
	// SymbolNum is the count of symbol positions that contributed value
	SymbolNum int
	// Wins is the total coin value accumulated by this component
	Wins      int
}

// OnNewGame is called when a new game session starts.
// It forwards to BasicComponentData.OnNewGame to initialize base fields.
func (symbolValWinsData *SymbolValWinsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolValWinsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep resets the per-step fields of SymbolValWinsData. This is
// invoked when a new step (spin/round) begins so previous results don't
// leak into the new computation.
func (symbolValWinsData *SymbolValWinsData) onNewStep() {
	symbolValWinsData.UsedResults = nil
	symbolValWinsData.SymbolNum = 0
	symbolValWinsData.Wins = 0
}

// Clone produces a deep-ish copy of the component data. It clones the
// embedded BasicComponentData and copies primitive fields.
func (symbolValWinsData *SymbolValWinsData) Clone() IComponentData {
	target := &SymbolValWinsData{
		BasicComponentData: symbolValWinsData.CloneBasicComponentData(),
		SymbolNum:          symbolValWinsData.SymbolNum,
		Wins:               symbolValWinsData.Wins,
	}

	return target
}

// BuildPBComponentData serializes runtime data into the protobuf message
// used for debugging or telemetry. The returned message is of type
// sgc7pb.SymbolValWinsData and includes base component data as well as
// SymbolNum and Wins.
func (symbolValWinsData *SymbolValWinsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolValWinsData{
		BasicComponentData: symbolValWinsData.BuildPBBasicComponentData(),
	}

	pbcd.SymbolNum = int32(symbolValWinsData.SymbolNum)
	pbcd.Wins = int32(symbolValWinsData.Wins)

	return pbcd
}

// GetValEx returns extra integer values from the component data by key.
// Supported keys:
//   - SVWDVSymbolNum -> number of contributing symbols
//   - SVWDVWins -> total coin wins produced
//   - CVResultNum, CVWinResultNum -> number of used results
func (symbolValWinsData *SymbolValWinsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	switch key {
	case SVWDVSymbolNum:
		return symbolValWinsData.SymbolNum, true
	case SVWDVWins:
		return symbolValWinsData.Wins, true
	case CVResultNum, CVWinResultNum:
		return len(symbolValWinsData.UsedResults), true
	}

	return 0, false
}

// SymbolValWinsConfig describes the YAML/JSON configuration for a
// SymbolValWins component. It embeds BasicComponentConfig and adds
// fields to select symbol sets, coin symbols and behavior type.
type SymbolValWinsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetTypeString        string            `yaml:"betType" json:"betType"`   // bet or totalBet or noPay
	BetType              BetType           `yaml:"-" json:"-"`               // bet or totalBet or noPay
	WinMulti             int               `yaml:"winMulti" json:"winMulti"` // bet or totalBet
	Symbols              []string          `yaml:"symbols" json:"symbols"`   // like collect
	SymbolCodes          []int             `yaml:"-" json:"-"`               //
	StrType              string            `yaml:"type" json:"type"`
	Type                 SymbolValWinsType `yaml:"-" json:"-"`
	CoinSymbols          []string          `yaml:"coinSymbols" json:"coinSymbols"` // coin symbols
	CoinSymbolCodes      []int             `yaml:"-" json:"-"`                     // coin symbols
}

// SetLinkComponent
func (cfg *SymbolValWinsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type SymbolValWins struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValWinsConfig `json:"config"`
}

// SymbolValWins is a component that creates coin-value results based on an
// "other scene" (value map) and optionally collector symbols on the
// primary scene. It supports normal, collector and reelcollector modes.

// Init reads YAML config from file and initializes the component.
// It unmarshals into SymbolValWinsConfig and delegates to InitEx.
func (svw *SymbolValWins) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolValWins.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolValWinsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolValWins.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return svw.InitEx(cfg, pool)
}

// InitEx initializes the component from an already-parsed config object.
// It resolves symbol names to codes using the provided GamePropertyPool
// and performs basic validation.
func (svw *SymbolValWins) InitEx(cfg any, pool *GamePropertyPool) error {
	svw.Config = cfg.(*SymbolValWinsConfig)
	svw.Config.ComponentType = SymbolValWinsTypeName

	svw.Config.BetType = ParseBetType(svw.Config.BetTypeString)
	svw.Config.Type = parseSymbolValWinsType(svw.Config.StrType)

	for _, s := range svw.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolValWins.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		svw.Config.SymbolCodes = append(svw.Config.SymbolCodes, sc)
	}

	for _, s := range svw.Config.CoinSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolValWins.InitEx:CoinSymbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		svw.Config.CoinSymbolCodes = append(svw.Config.CoinSymbolCodes, sc)
	}

	svw.onInit(&svw.Config.BasicComponentConfig)

	return nil
}

// OnPlayGame executes the component logic for a single play step.
// Behavior summary:
//   - locate the target scene (gs) and other scene (os)
//   - depending on Type (normal/collector/reelcollector) determine multiplier
//     and collector positions
//   - scan the other scene for coin values (optionally filtered by coin symbols)
//   - if any values found, create RTCoins results (one per multiplier) with
//     appropriate Pos, CoinWin and CashWin and add them to PlayResult
// It returns the next component name and an error (ErrComponentDoNothing
// if nothing was produced).
func (svw *SymbolValWins) OnPlayGame(gameProp *GameProperty, pr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	svwd, ok := icd.(*SymbolValWinsData)
	if !ok || svwd == nil {
		goutils.Error("SymbolValWins.OnPlayGame:SymbolValWinsData",
			goutils.Err(ErrInvalidComponentData))

		return "", ErrInvalidComponentData
	}

	svwd.onNewStep()

	gs := svw.GetTargetScene3(gameProp, pr, prs, 0)
	if gs == nil {
		goutils.Error("SymbolValWins.OnPlayGame:GetTargetScene3",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	os := svw.GetTargetOtherScene3(gameProp, pr, prs, 0)

	if os != nil {
	// collectorpos stores pairs of (x,y) coordinates for collector symbols.
	// mul is the number of collectors (effectively how many separate
	// RTCoins results we will produce when collectors are present).
	collectorpos := []int{}
	mul := 0
		switch svw.Config.Type {
		case svwTypeCollector:
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(svw.Config.SymbolCodes, s, 0) >= 0 {
						mul++

						collectorpos = append(collectorpos, x, y)
					}
				}
			}
		case svwTypeReelCollector:
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(svw.Config.SymbolCodes, s, 0) >= 0 {
						mul++

						collectorpos = append(collectorpos, x, y)

						break
					}
				}
			}
		default:
			mul = 1
		}

	// totalvals accumulates the sum of coin values from the other scene.
	// pos collects the positions (x,y) that contribute; stored as a flat
	// slice [x0,y0,x1,y1,...]. Preallocate using os.Width/os.Height
	// which are provided by the scene for convenience.
	totalvals := 0

	pos := make([]int, 0, os.Width*os.Height*2)

		if len(svw.Config.CoinSymbolCodes) > 0 {
			for x, arr := range os.Arr {
				for y, v := range arr {
					if v > 0 && slices.Contains(svw.Config.CoinSymbolCodes, gs.Arr[x][y]) {
						totalvals += v
						pos = append(pos, x, y)

						svwd.SymbolNum++
					}
				}
			}
		} else {
			for x, arr := range os.Arr {
				for y, v := range arr {
					if v > 0 {
						totalvals += v
						pos = append(pos, x, y)

						svwd.SymbolNum++
					}
				}
			}
		}

	// If we found any coin values and we have at least one multiplier
	// (normal mode sets mul=1), construct Result entries and append
	// them to the PlayResult.
	if totalvals > 0 && mul > 0 {
			bet := gameProp.GetBet3(stake, svw.Config.BetType)
			othermul := svw.GetWinMulti(&svwd.BasicComponentData)

			for i := 0; i < mul; i++ {
				// build a new position list for this result. For collector
				// variants, the collector coordinate precedes the coin pos list.
				newpos := make([]int, 0, len(pos)+2)

				if svw.isCollectorType() {
					newpos = append(newpos, collectorpos[i*2], collectorpos[i*2+1])
				}

				newpos = append(newpos, pos...)

				// Construct the result object describing coin wins. SymbolNums
				// is the count of contributing coin positions.
				ret := &sgc7game.Result{
					Type:       sgc7game.RTCoins,
					LineIndex:  -1,
					Pos:        newpos,
					SymbolNums: len(pos) / 2,
					Mul:        1,
				}

				if svw.isCollectorType() {
					// For collector results, set the symbol to the collector's
					// symbol code so consumers can render it.
					ret.Symbol = gs.Arr[newpos[0]][newpos[1]]
				}

				ret.CoinWin = totalvals * othermul
				ret.CashWin = ret.CoinWin * bet
				ret.OtherMul = othermul

				svwd.Wins += ret.CoinWin

				svw.AddResult(pr, ret, &svwd.BasicComponentData)
			}

			nc := svw.onStepEnd(gameProp, pr, gp, "")

			return nc, nil
		}
	}

	nc := svw.onStepEnd(gameProp, pr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (svw *SymbolValWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd, ok := icd.(*SymbolValWinsData)
	if !ok || cd == nil {
		goutils.Error("SymbolValWins.OnAsciiGame:SymbolValWinsData",
			goutils.Err(ErrInvalidComponentData))

		return ErrInvalidComponentData
	}

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// NewComponentData -
// NewComponentData creates a fresh SymbolValWinsData instance for use by
// the runtime when a new player/session/component instance is created.
func (svw *SymbolValWins) NewComponentData() IComponentData {
	return &SymbolValWinsData{}
}

// NewStats2 -
// NewStats2 builds a stats2.Feature that describes which metrics this
// component will emit. It returns a feature that reports wins and an
// integer value (used for win multipliers).
func (svw *SymbolValWins) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins, stats2.OptIntVal})
}

// OnStats2
func (svw *SymbolValWins) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	svw.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	svwd, ok := icd.(*SymbolValWinsData)
	if !ok || svwd == nil {
		goutils.Error("SymbolValWins.OnStats2:SymbolValWinsData",
			goutils.Err(ErrInvalidComponentData))

		return
	}

	s2.ProcStatsWins(svw.Name, int64(svwd.Wins))

	multi := svw.GetWinMulti(&svwd.BasicComponentData)

	s2.ProcStatsIntVal(svw.GetName(), multi)
}

// GetWinMulti returns the effective multiplier to apply to coin values.
// It first consults runtime overrides in BasicComponentData, and falls
// back to the configured WinMulti.
func (svw *SymbolValWins) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return svw.Config.WinMulti
}

// isCollectorType reports whether the config type is a collector variant
func (svw *SymbolValWins) isCollectorType() bool {
	return svw.Config.Type == svwTypeCollector || svw.Config.Type == svwTypeReelCollector
}

func NewSymbolValWins(name string) IComponent {
	return &SymbolValWins{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "betType": "bet",
// "winMulti": 1,
// "type": "normal",
// "coinSymbols": [
//
//	"CA"
//
// ]
type jsonSymbolValWins struct {
	BetType     string   `json:"betType"`  // bet or totalBet or noPay
	WinMulti    int      `json:"winMulti"` // bet or totalBet
	Symbols     []string `json:"symbols"`  // like collect
	Type        string   `yaml:"type" json:"type"`
	CoinSymbols []string `json:"coinSymbols"` // coin symbols
}

func (jcfg *jsonSymbolValWins) build() *SymbolValWinsConfig {
	cfg := &SymbolValWinsConfig{
		BetTypeString: jcfg.BetType,
		WinMulti:      jcfg.WinMulti,
		Symbols:       jcfg.Symbols,
		StrType:       strings.ToLower(jcfg.Type),
		CoinSymbols:   jcfg.CoinSymbols,
	}

	return cfg
}

// jsonSymbolValWins is an intermediate structure used to parse JSON/YAML
// configuration embedded in spreadsheets or design tools. The build
// method converts it into the internal SymbolValWinsConfig.


// parseSymbolValWins parses a configuration cell (AST node) from a design
// spreadsheet into a SymbolValWinsConfig and registers it with the
// provided BetConfig. It returns the component label/name or an error.
func parseSymbolValWins(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolValWins:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolValWins:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolValWins{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolValWins:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolValWinsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
