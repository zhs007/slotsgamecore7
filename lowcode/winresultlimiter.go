package lowcode

import (
	"fmt"
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
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WinResultLimiterTypeName = "winResultLimiter"

type WinResultLimiterType int

const (
	WRLTypeMaxOnLine WinResultLimiterType = 0
)

func parseWinResultLimiterType(str string) WinResultLimiterType {
	s := strings.TrimSpace(strings.ToLower(str))
	if s == "maxonline" {
		return WRLTypeMaxOnLine
	}

	// default to max on line if unknown (backwards compatible)
	return WRLTypeMaxOnLine
}

// WinResultLimiterData holds runtime data for the WinResultLimiter component.
// It embeds BasicComponentData and tracks the total wins calculated by the limiter.
type WinResultLimiterData struct {
	BasicComponentData
	Wins int
}

// OnNewGame calls BasicComponentData.OnNewGame to initialize component data at the start of a new game.
func (winResultLimiterData *WinResultLimiterData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultLimiterData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep resets the per-step counters for the component data.
func (winResultLimiterData *WinResultLimiterData) onNewStep() {
	winResultLimiterData.Wins = 0
}

// Clone creates a deep copy of the component data and returns the clone.
func (winResultLimiterData *WinResultLimiterData) Clone() IComponentData {
	target := &WinResultLimiterData{
		BasicComponentData: winResultLimiterData.CloneBasicComponentData(),
		Wins:               winResultLimiterData.Wins,
	}

	return target
}

// BuildPBComponentData builds a protobuf message representing the component data.
func (winResultLimiterData *WinResultLimiterData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultLimiterData{
		BasicComponentData: winResultLimiterData.BuildPBBasicComponentData(),
		Wins:               int32(winResultLimiterData.Wins),
	}

	return pbcd
}

// GetValEx returns named values exposed by this component data.
// Supported keys: CVWins -> returns the accumulated wins for the limiter.
func (winResultLimiterData *WinResultLimiterData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultLimiterData.Wins, true
	}

	return 0, false
}

// WinResultLimiterConfig is the configuration for the WinResultLimiter component.
// It is parsed from YAML/JSON and controls limiter behavior.
type WinResultLimiterConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string               `yaml:"type" json:"type"`                   // type
	Type                 WinResultLimiterType `yaml:"-" json:"-"`                         // type
	SrcComponents        []string             `yaml:"srcComponents" json:"srcComponents"` // srcComponents
}

// SetLinkComponent sets a link between components such as the "next" component.
// If link is "next", the DefaultNextComponent field is updated.
func (cfg *WinResultLimiterConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// WinResultLimiter is a component that limits win results according to configured rules.
type WinResultLimiter struct {
	*BasicComponent `json:"-"`
	Config          *WinResultLimiterConfig `json:"config"`
}

// Init loads a YAML configuration file and initializes the component.
// Init reads the file from fn and unmarshals it into a WinResultLimiterConfig before calling InitEx.
func (w *WinResultLimiter) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WinResultLimiter.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WinResultLimiterConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultLimiter.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return w.InitEx(cfg, pool)
}

// InitEx initializes the component from a parsed configuration object.
// InitEx validates the cfg type and sets internal fields required for operation.
func (winResultLimiter *WinResultLimiter) InitEx(cfg any, pool *GamePropertyPool) error {
	c, ok := cfg.(*WinResultLimiterConfig)
	if !ok {
		goutils.Error("WinResultLimiter.InitEx:InvalidCfg",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	winResultLimiter.Config = c
	// set correct component type for limiter
	winResultLimiter.Config.ComponentType = WinResultLimiterTypeName

	winResultLimiter.Config.Type = parseWinResultLimiterType(winResultLimiter.Config.StrType)

	winResultLimiter.onInit(&winResultLimiter.Config.BasicComponentConfig)

	return nil
}

// onMaxOnLine enforces the max-on-line rule: for each line, only the highest win remains.
// Other wins on the same line are zeroed out and the kept win is accumulated into cd.Wins.
func (winResultLimiter *WinResultLimiter) onMaxOnLine(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd *WinResultLimiterData) (string, error) {

	mapLinesWin := make(map[int][]int)

	for _, cn := range winResultLimiter.Config.SrcComponents {
		// If previous components haven't executed, their data may not be cleaned up,
		// so skip this source component.
		if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
			continue
		}

		ccd := gameProp.GetComponentDataWithName(cn)
		if ccd == nil {
			goutils.Error("WinResultLimiter.onMaxOnLine:MissingComponentData",
				slog.String("component", cn),
				goutils.Err(ErrInvalidComponentData))

			continue
		}

		lst := ccd.GetResults()
		for _, ri := range lst {
			if ri < 0 || ri >= len(curpr.Results) {
				goutils.Error("WinResultLimiter.onMaxOnLine:ResultIndexOutOfRange",
					slog.Int("ri", ri),
					slog.Int("results_len", len(curpr.Results)))

				continue
			}

			curline := curpr.Results[ri].LineIndex
			mapLinesWin[curline] = append(mapLinesWin[curline], ri)
		}
	}

	if len(mapLinesWin) <= 0 {
		nc := winResultLimiter.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	cd.Wins = 0

	for _, lst := range mapLinesWin {
		if len(lst) <= 1 {
			continue
		}

		maxwin := 0
		maxwi := -1
		for _, ri := range lst {
			if curpr.Results[ri].CoinWin > maxwin {
				maxwin = curpr.Results[ri].CoinWin
				maxwi = ri
			}
		}

		for _, ri := range lst {
			if ri != maxwi {
				curpr.Results[ri].CoinWin = 0
				curpr.Results[ri].CashWin = 0
			} else {
				cd.Wins += curpr.Results[ri].CoinWin
			}
		}
	}

	if cd.Wins == 0 {
		nc := winResultLimiter.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := winResultLimiter.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnPlayGame is the runtime entry used during play to apply limiter logic to a PlayResult.
// OnPlayGame dispatches to the configured limiter behavior for the current step.
func (winResultLimiter *WinResultLimiter) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*WinResultLimiterData)
	cd.onNewStep()

	if winResultLimiter.Config.Type == WRLTypeMaxOnLine {
		return winResultLimiter.onMaxOnLine(gameProp, curpr, gp, cd)
	}

	goutils.Error("WinResultLimiter.OnPlayGame:InvalidType",
		slog.String("Type", winResultLimiter.Config.StrType),
		goutils.Err(ErrInvalidComponentConfig))

	return "", ErrInvalidComponentConfig
}

// OnAsciiGame outputs component information for ASCII game mode.
// OnAsciiGame prints the ending wins to stdout for legacy ASCII output.
func (winResultModifier *WinResultLimiter) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	std := icd.(*WinResultLimiterData)

	// This uses fmt.Printf because it's a special ASCII game mode that expects stdout output.
	fmt.Printf("WinResultLimiter, ending wins = %v \n", std.Wins)

	return nil
}

// NewComponentData creates and returns a fresh WinResultLimiterData instance.
func (winResultModifier *WinResultLimiter) NewComponentData() IComponentData {
	return &WinResultLimiterData{}
}

func NewWinResultLimiter(name string) IComponent {
	return &WinResultLimiter{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "maxOnLine",
// "srcComponents": [
//
//	"fg-wins",
//	"fg-wins-h3",
//	"fg-wins-h4"
//
// ]
type jsonWinResultLimiter struct {
	Type          string   `json:"type"`          // type
	SrcComponents []string `json:"srcComponents"` // srcComponents
}

func (jwt *jsonWinResultLimiter) build() *WinResultLimiterConfig {
	cfg := &WinResultLimiterConfig{
		StrType:       jwt.Type,
		SrcComponents: slices.Clone(jwt.SrcComponents),
	}

	return cfg
}

func parseWinResultLimiter(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultLimiter:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultLimiter:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWinResultLimiter{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultLimiter:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultLimiterTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
