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

const WinResultModifierTypeName = "winResultModifier"

type WinResultModifierType int

const (
	WRMTypeExistSymbol WinResultModifierType = 0
)

func parseWinResultModifierType(str string) WinResultModifierType {
	if str == "existSymbol" {
		return WRMTypeExistSymbol
	}

	return WRMTypeExistSymbol
}

type WinResultModifierData struct {
	BasicComponentData
	Wins     int
	WinMulti int
}

// OnNewGame -
func (winResultModifierData *WinResultModifierData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultModifierData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (winResultModifierData *WinResultModifierData) onNewStep() {
	// winResultMultiData.BasicComponentData.OnNewStep(gameProp, component)

	winResultModifierData.Wins = 0
	winResultModifierData.WinMulti = 1
}

// Clone
func (winResultModifierData *WinResultModifierData) Clone() IComponentData {
	target := &WinResultModifierData{
		BasicComponentData: winResultModifierData.CloneBasicComponentData(),
		Wins:               winResultModifierData.Wins,
		WinMulti:           winResultModifierData.WinMulti,
	}

	return target
}

// BuildPBComponentData
func (winResultModifierData *WinResultModifierData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultModifierData{
		BasicComponentData: winResultModifierData.BuildPBBasicComponentData(),
		Wins:               int32(winResultModifierData.Wins),
		WinMulti:           int32(winResultModifierData.WinMulti),
	}

	return pbcd
}

// GetVal -
func (winResultModifierData *WinResultModifierData) GetVal(key string) (int, bool) {
	if key == CVWins {
		return winResultModifierData.Wins, true
	}

	return 0, false
}

// SetVal -
func (winResultModifierData *WinResultModifierData) SetVal(key string, val int) {
	if key == CVWins {
		winResultModifierData.Wins = val
	}
}

// WinResultModifierConfig - configuration for WinResultModifier
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WinResultModifierConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`                         // type
	Type                 WinResultModifierType `yaml:"-" json:"-"`                               // type
	SourceComponents     []string              `yaml:"sourceComponents" json:"sourceComponents"` // target components
	WinMulti             int                   `yaml:"winMulti" json:"winMulti"`                 // winMulti，最后的中奖倍数，默认为1
	TargetSymbols        []string              `yaml:"sourceComponents" json:"targetSymbols"`    // targetSymbols
	TargetSymbolCodes    []int                 `yaml:"-" json:"-"`                               // target SymbolCodes
}

// SetLinkComponent
func (cfg *WinResultModifierConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WinResultModifier struct {
	*BasicComponent `json:"-"`
	Config          *WinResultModifierConfig `json:"config"`
}

// Init -
func (winResultModifier *WinResultModifier) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WinResultModifier.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WinResultModifierConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultModifier.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return winResultModifier.InitEx(cfg, pool)
}

// InitEx -
func (winResultModifier *WinResultModifier) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultModifier.Config = cfg.(*WinResultModifierConfig)
	winResultModifier.Config.ComponentType = WinResultModifierTypeName

	winResultModifier.Config.Type = parseWinResultModifierType(winResultModifier.Config.StrType)

	if winResultModifier.Config.WinMulti <= 0 {
		winResultModifier.Config.WinMulti = 1
	}

	for _, s := range winResultModifier.Config.TargetSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("WinResultModifier.InitEx:TargetSymbols.Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		winResultModifier.Config.TargetSymbolCodes = append(winResultModifier.Config.TargetSymbolCodes, sc)
	}

	winResultModifier.onInit(&winResultModifier.Config.BasicComponentConfig)

	return nil
}

// playgame
func (winResultModifier *WinResultModifier) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// winResultMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := icd.(*WinResultModifierData)
	std.onNewStep()

	winMulti := winResultModifier.GetWinMulti(&std.BasicComponentData)

	std.WinMulti = winMulti
	std.Wins = 0

	if winMulti == 1 {
		nc := winResultModifier.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	gs := winResultModifier.GetTargetScene3(gameProp, curpr, prs, 0)

	for _, cn := range winResultModifier.Config.SourceComponents {
		// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
		if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
			continue
		}

		ccd := gameProp.GetComponentDataWithName(cn)
		// ccd := gameProp.MapComponentData[cn]
		lst := ccd.GetResults()
		for _, ri := range lst {
			if HasSymbolsInResult(gs, winResultModifier.Config.TargetSymbolCodes, curpr.Results[ri]) {
				curpr.Results[ri].CashWin *= winMulti
				curpr.Results[ri].CoinWin *= winMulti
				curpr.Results[ri].OtherMul *= winMulti

				std.Wins += curpr.Results[ri].CoinWin
			}
		}
	}

	nc := winResultModifier.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (winResultModifier *WinResultModifier) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WinResultModifierData)

	fmt.Printf("WinResultModifier x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// // OnStatsWithPB -
// func (winResultMulti *WinResultMulti) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.WinResultMultiData)
// 	if !isok {
// 		goutils.Error("WinResultMulti.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return winResultMulti.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (winResultMulti *WinResultMulti) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (winResultModifier *WinResultModifier) NewComponentData() IComponentData {
	return &WinResultModifierData{}
}

func (winResultModifier *WinResultModifier) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		if winMulti <= 0 {
			return 1
		}

		return winMulti
	}

	if winResultModifier.Config.WinMulti <= 0 {
		return 1
	}

	return winResultModifier.Config.WinMulti
}

func NewWinResultModifier(name string) IComponent {
	return &WinResultModifier{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"winMulti": 2,
//		"type": "existSymbol",
//		"sourceComponent": [
//			"fg-payfg"
//		],
//		"targetSymbols": ["RW2"]
//	},
type jsonWinResultModifier struct {
	Type             string   `json:"type"`            // type
	SourceComponents []string `json:"sourceComponent"` // source components
	WinMulti         int      `json:"winMulti"`        // winMulti
	TargetSymbols    []string `json:"targetSymbols"`   // targetSymbols
}

func (jwt *jsonWinResultModifier) build() *WinResultModifierConfig {
	cfg := &WinResultModifierConfig{
		StrType:          jwt.Type,
		SourceComponents: jwt.SourceComponents,
		WinMulti:         jwt.WinMulti,
		TargetSymbols:    jwt.TargetSymbols,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWinResultModifier(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultModifier:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultModifier:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWinResultModifier{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultModifier:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultModifierTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
