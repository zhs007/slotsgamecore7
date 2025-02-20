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
	"github.com/zhs007/slotsgamecore7/stats2"
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

// Clone
func (winResultMultiData *WinResultMultiData) Clone() IComponentData {
	target := &WinResultMultiData{
		BasicComponentData: winResultMultiData.CloneBasicComponentData(),
		Wins:               winResultMultiData.Wins,
		WinMulti:           winResultMultiData.WinMulti,
	}

	return target
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

// GetValEx -
func (winResultMultiData *WinResultMultiData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultMultiData.Wins, true
	} else if key == CVWinMulti {
		return winResultMultiData.WinMulti, true
	}

	return 0, false
}

// // SetVal -
// func (winResultMultiData *WinResultMultiData) SetVal(key string, val int) {
// 	if key == CVWins {
// 		winResultMultiData.Wins = val
// 	}
// }

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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WinResultMultiConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultMulti.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
		// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
		if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
			continue
		}

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
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return winResultMulti.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (winResultMulti *WinResultMulti) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// OnStats2
func (winResultMulti *WinResultMulti) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	winResultMulti.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*WinResultMultiData)

	s2.ProcStatsIntVal(winResultMulti.GetName(), winResultMulti.GetWinMulti(&cd.BasicComponentData))
}

// NewStats2 -
func (winResultMulti *WinResultMulti) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptIntVal})
}

// NewComponentData -
func (winResultMulti *WinResultMulti) NewComponentData() IComponentData {
	return &WinResultMultiData{}
}

func (winResultMulti *WinResultMulti) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		if winMulti <= 0 {
			return 1
		}

		return winMulti
	}

	if winResultMulti.Config.WinMulti <= 0 {
		return 1
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

func parseWinResultMulti(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultMulti:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultMulti:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWinResultMulti{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultMulti:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultMultiTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
