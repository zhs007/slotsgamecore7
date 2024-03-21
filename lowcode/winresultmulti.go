package lowcode

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WinResultMultiTypeName = "winResultMulti"

const (
	WRMCVWinMulti string = "winMulti" // 可以修改配置项里的winMulti
)

type WinResultMultiData struct {
	BasicComponentData
	Wins     int
	WinMulti int
}

// OnNewGame -
func (winResultMultiData *WinResultMultiData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultMultiData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (winResultMultiData *WinResultMultiData) onNewStep() {
	// winResultMultiData.BasicComponentData.OnNewStep(gameProp, component)

	winResultMultiData.Wins = 0
	winResultMultiData.WinMulti = 1
}

// BuildPBComponentData
func (winResultMultiData *WinResultMultiData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultMultiData{
		BasicComponentData: winResultMultiData.BuildPBBasicComponentData(),
		Wins:               int32(winResultMultiData.Wins),
		WinMulti:           int32(winResultMultiData.WinMulti),
	}

	return pbcd
}

// GetVal -
func (winResultMultiData *WinResultMultiData) GetVal(key string) int {
	if key == CVWins {
		return winResultMultiData.Wins
	}

	return 0
}

// SetVal -
func (winResultMultiData *WinResultMultiData) SetVal(key string, val int) {
	if key == CVWins {
		winResultMultiData.Wins = val
	}
}

// WinResultMultiConfig - configuration for WinResultMulti
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WinResultMultiConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	TargetComponents     []string `yaml:"targetComponents" json:"targetComponents"` // target components
	WinMulti             int      `yaml:"winMulti" json:"winMulti"`                 // winMulti，最后的中奖倍数，默认为1
}

// SetLinkComponent
func (cfg *WinResultMultiConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WinResultMulti struct {
	*BasicComponent `json:"-"`
	Config          *WinResultMultiConfig `json:"config"`
}

// Init -
func (winResultMulti *WinResultMulti) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WinResultMulti.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WinResultMultiConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultMulti.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return winResultMulti.InitEx(cfg, pool)
}

// InitEx -
func (winResultMulti *WinResultMulti) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultMulti.Config = cfg.(*WinResultMultiConfig)
	winResultMulti.Config.ComponentType = WinResultMultiTypeName

	if winResultMulti.Config.WinMulti <= 0 {
		winResultMulti.Config.WinMulti = 1
	}

	winResultMulti.onInit(&winResultMulti.Config.BasicComponentConfig)

	return nil
}

// playgame
func (winResultMulti *WinResultMulti) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// winResultMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	std := icd.(*WinResultMultiData)
	std.onNewStep()

	winMulti := winResultMulti.GetWinMulti(&std.BasicComponentData)

	std.WinMulti = winMulti
	std.Wins = 0

	if winMulti == 1 {
		nc := winResultMulti.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	for _, cn := range winResultMulti.Config.TargetComponents {
		ccd := gameProp.GetComponentDataWithName(cn)
		// ccd := gameProp.MapComponentData[cn]
		lst := ccd.GetResults()
		for _, ri := range lst {
			curpr.Results[ri].CashWin *= winMulti
			curpr.Results[ri].CoinWin *= winMulti
			curpr.Results[ri].OtherMul *= winMulti

			std.Wins += curpr.Results[ri].CoinWin
		}
	}

	nc := winResultMulti.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (winResultMulti *WinResultMulti) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WinResultMultiData)

	fmt.Printf("winResultMulti x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// // OnStatsWithPB -
// func (winResultMulti *WinResultMulti) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.WinResultMultiData)
// 	if !isok {
// 		goutils.Error("WinResultMulti.OnStatsWithPB",
// 			zap.Error(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return winResultMulti.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (winResultMulti *WinResultMulti) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (winResultMulti *WinResultMulti) NewComponentData() IComponentData {
	return &WinResultMultiData{}
}

func (winResultMulti *WinResultMulti) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(WTCVWinMulti)
	if isok {
		return winMulti
	}

	return winResultMulti.Config.WinMulti
}

func NewWinResultMulti(name string) IComponent {
	return &WinResultMulti{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"triggerType": "lines",
//		"betType": "bet",
//		"checkWinType": "left2right",
//		"symbols": [
//			"WL",
//			"A",
//			"B",
//			"C",
//			"D",
//			"E",
//			"F",
//			"G",
//			"H",
//			"J",
//			"K",
//			"L"
//		],
//		"wildSymbols": [
//			"WL"
//		]
//	},
type jsonWinResultMulti struct {
	TargetComponents []string `json:"targetComponents"` // target components
	WinMulti         int      `json:"winMulti"`
}

func (jwt *jsonWinResultMulti) build() *WinResultMultiConfig {
	cfg := &WinResultMultiConfig{
		TargetComponents: jwt.TargetComponents,
		WinMulti:         jwt.WinMulti,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWinResultMulti(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultMulti:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultMulti:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonWinResultMulti{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultMulti:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultMultiTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
