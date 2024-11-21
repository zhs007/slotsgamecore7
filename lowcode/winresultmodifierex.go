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

const WinResultModifierExTypeName = "winResultModifierEx"

type WinResultModifierExData struct {
	BasicComponentData
	Wins     int
	WinMulti int
}

// OnNewGame -
func (winResultModifierDataEx *WinResultModifierExData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultModifierDataEx.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (winResultModifierDataEx *WinResultModifierExData) onNewStep() {
	winResultModifierDataEx.Wins = 0
	winResultModifierDataEx.WinMulti = 1
}

// Clone
func (winResultModifierDataEx *WinResultModifierExData) Clone() IComponentData {
	target := &WinResultModifierData{
		BasicComponentData: winResultModifierDataEx.CloneBasicComponentData(),
		Wins:               winResultModifierDataEx.Wins,
		WinMulti:           winResultModifierDataEx.WinMulti,
	}

	return target
}

// BuildPBComponentData
func (winResultModifierDataEx *WinResultModifierExData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultModifierData{
		BasicComponentData: winResultModifierDataEx.BuildPBBasicComponentData(),
		Wins:               int32(winResultModifierDataEx.Wins),
		WinMulti:           int32(winResultModifierDataEx.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (winResultModifierDataEx *WinResultModifierExData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultModifierDataEx.Wins, true
	}

	return 0, false
}

// WinResultModifierExConfig - configuration for WinResultModifierEx
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WinResultModifierExConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`                         // type
	Type                 WinResultModifierType `yaml:"-" json:"-"`                               // type
	SourceComponents     []string              `yaml:"sourceComponents" json:"sourceComponents"` // target components
	MapTargetSymbols     map[string]int        `yaml:"mapTargetSymbols" json:"mapTargetSymbols"` // mapTargetSymbols
	MapTargetSymbolCodes map[int]int           `yaml:"-" json:"-"`                               // MapTargetSymbolCodes
}

// SetLinkComponent
func (cfg *WinResultModifierExConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WinResultModifierEx struct {
	*BasicComponent `json:"-"`
	Config          *WinResultModifierExConfig `json:"config"`
}

// Init -
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

// InitEx -
func (winResultModifierEx *WinResultModifierEx) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultModifierEx.Config = cfg.(*WinResultModifierExConfig)
	winResultModifierEx.Config.ComponentType = WinResultModifierExTypeName

	winResultModifierEx.Config.Type = parseWinResultModifierType(winResultModifierEx.Config.StrType)
	if winResultModifierEx.Config.Type == WRMTypeExistSymbol {
		goutils.Error("WinResultModifierEx.InitEx:WRMTypeExistSymbol",
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

// playgame
func (winResultModifierEx *WinResultModifierEx) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	std := icd.(*WinResultModifierExData)
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

// OnAsciiGame - outpur to asciigame
func (winResultModifierEx *WinResultModifierEx) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WinResultModifierExData)

	fmt.Printf("WinResultModifierEx x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// NewComponentData -
func (winResultModifierEx *WinResultModifierEx) NewComponentData() IComponentData {
	return &WinResultModifierExData{}
}

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
		cfg.MapTargetSymbols[arr[0].(string)] = int(arr[1].(float64))
	}

	return cfg
}

func parseWinResultModifierEx(gamecfg *BetConfig, cell *ast.Node) (string, error) {
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

	data := &jsonWinResultModifier{}

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
