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

// 已弃用，待清理

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

// Clone
func (winResultCacheData *WinResultCacheData) Clone() IComponentData {
	target := &WinResultCacheData{
		BasicComponentData: winResultCacheData.CloneBasicComponentData(),
		Wins:               winResultCacheData.Wins,
		WinMulti:           winResultCacheData.WinMulti,
		WinResultNum:       winResultCacheData.WinResultNum,
	}

	return target
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

// GetValEx -
func (winResultCacheData *WinResultCacheData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultCacheData.Wins, true
	}

	return 0, false
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WinResultCacheConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WinResultCache.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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

	nc := winResultCache.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (winResultCache *WinResultCache) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*WinResultCacheData)

	fmt.Printf("winResultCache x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// NewComponentData -
func (winResultCache *WinResultCache) NewComponentData() IComponentData {
	return &WinResultCacheData{}
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

	return cfg
}

func parseWinResultCache(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseWinResultCache:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseWinResultCache:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonWinResultCache{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseWinResultCache:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: WinResultCacheTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
