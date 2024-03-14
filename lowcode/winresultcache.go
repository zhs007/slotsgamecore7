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
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

// 发现这个组件没用了，先标记一下
const WinResultCacheTypeName = "winResultCache"

type WinResultCacheData struct {
	BasicComponentData
	WinResultNum int
	Wins         int
	WinMulti     int
}

// OnNewGame -
func (winResultCacheData *WinResultCacheData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultCacheData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (winResultCacheData *WinResultCacheData) OnNewStep(gameProp *GameProperty, component IComponent) {
	winResultCacheData.BasicComponentData.OnNewStep(gameProp, component)

	winResultCacheData.Wins = 0
	winResultCacheData.WinMulti = 1
}

// BuildPBComponentData
func (winResultCacheData *WinResultCacheData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultCacheData{
		BasicComponentData: winResultCacheData.BuildPBBasicComponentData(),
		Wins:               int32(winResultCacheData.Wins),
		WinMulti:           int32(winResultCacheData.WinMulti),
		WinResultNum:       int32(winResultCacheData.WinResultNum),
	}

	return pbcd
}

// GetVal -
func (winResultCacheData *WinResultCacheData) GetVal(key string) int {
	if key == CVWins {
		return winResultCacheData.Wins
	}

	return 0
}

// SetVal -
func (winResultCacheData *WinResultCacheData) SetVal(key string, val int) {
	if key == CVWins {
		winResultCacheData.Wins = val
	}
}

// WinResultCacheConfig - configuration for WinResultCache
type WinResultCacheConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	WinMulti             int `yaml:"winMulti" json:"winMulti"` // winMulti，最后的中奖倍数，默认为1
}

// SetLinkComponent
func (cfg *WinResultCacheConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WinResultCache struct {
	*BasicComponent `json:"-"`
	Config          *WinResultCacheConfig `json:"config"`
}

// Init -
func (winResultCache *WinResultCache) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WinResultCache.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WinResultCacheConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultCache.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return winResultCache.InitEx(cfg, pool)
}

// InitEx -
func (winResultCache *WinResultCache) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultCache.Config = cfg.(*WinResultCacheConfig)
	winResultCache.Config.ComponentType = WinResultCacheTypeName

	if winResultCache.Config.WinMulti <= 0 {
		winResultCache.Config.WinMulti = 1
	}

	winResultCache.onInit(&winResultCache.Config.BasicComponentConfig)

	return nil
}

// playgame
func (winResultCache *WinResultCache) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// winResultMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	// cd := icd.(*WinResultCacheData)

	// cd.OnNewStep(gameProp, piggyBank)

	// winMulti := piggyBank.GetWinMulti(&cd.BasicComponentData)

	// cd.WinMulti = winMulti
	// sm, isok := cd.GetConfigIntVal(CCVSavedMoney)
	// if !isok {
	// 	nc := piggyBank.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, ErrComponentDoNothing
	// }

	// cd.Wins = sm * winMulti

	// bet := gameProp.GetBet2(stake, BTypeBet)

	// ret := &sgc7game.Result{
	// 	Symbol:    -1,
	// 	Type:      sgc7game.RTSymbolVal,
	// 	LineIndex: -1,
	// 	CoinWin:   cd.Wins,
	// 	CashWin:   cd.Wins * bet,
	// }

	// piggyBank.AddResult(curpr, ret, &cd.BasicComponentData)

	nc := winResultCache.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (winResultCache *WinResultCache) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WinResultCacheData)

	fmt.Printf("winResultCache x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// OnStatsWithPB -
func (winResultCache *WinResultCache) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.WinResultCacheData)
	if !isok {
		goutils.Error("WinResultCache.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return winResultCache.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// OnStats
func (winResultCache *WinResultCache) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (winResultCache *WinResultCache) NewComponentData() IComponentData {
	return &WinResultCacheData{}
}

func (winResultCache *WinResultCache) getWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(WTCVWinMulti)
	if isok {
		return winMulti
	}

	return winResultCache.Config.WinMulti
}

func NewWinResultCache(name string) IComponent {
	return &WinResultCache{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonWinResultCache struct {
	WinMulti int `json:"winMulti"`
}

func (jwt *jsonWinResultCache) build() *WinResultCacheConfig {
	cfg := &WinResultCacheConfig{
		WinMulti: jwt.WinMulti,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseWinResultCache(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultCache:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultCache:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonWinResultCache{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultCache:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultCacheTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
