package lowcode

import (
	"fmt"
	"log/slog"
	"os"
	"slices"

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
	if str == "maxonline" {
		return WRLTypeMaxOnLine
	}

	return WRLTypeMaxOnLine
}

type WinResultLimiterData struct {
	BasicComponentData
	Wins int
}

// OnNewGame -
func (winResultLimiterData *WinResultLimiterData) OnNewGame(gameProp *GameProperty, component IComponent) {
	winResultLimiterData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (winResultLimiterData *WinResultLimiterData) onNewStep() {
	winResultLimiterData.Wins = 0
}

// Clone
func (winResultLimiterData *WinResultLimiterData) Clone() IComponentData {
	target := &WinResultLimiterData{
		BasicComponentData: winResultLimiterData.CloneBasicComponentData(),
		Wins:               winResultLimiterData.Wins,
	}

	return target
}

// BuildPBComponentData
func (winResultLimiterData *WinResultLimiterData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WinResultLimiterData{
		BasicComponentData: winResultLimiterData.BuildPBBasicComponentData(),
		Wins:               int32(winResultLimiterData.Wins),
	}

	return pbcd
}

// GetValEx -
func (winResultLimiterData *WinResultLimiterData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultLimiterData.Wins, true
	}

	return 0, false
}

// WinResultLimiterConfig - configuration for WinResultLimiter
type WinResultLimiterConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string               `yaml:"type" json:"type"`                   // type
	Type                 WinResultLimiterType `yaml:"-" json:"-"`                         // type
	SrcComponents        []string             `yaml:"srcComponents" json:"srcComponents"` // srcComponents
}

// SetLinkComponent
func (cfg *WinResultLimiterConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type WinResultLimiter struct {
	*BasicComponent `json:"-"`
	Config          *WinResultLimiterConfig `json:"config"`
}

// Init -
func (winResultModifier *WinResultLimiter) Init(fn string, pool *GamePropertyPool) error {
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

	return winResultModifier.InitEx(cfg, pool)
}

// InitEx -
func (winResultLimiter *WinResultLimiter) InitEx(cfg any, pool *GamePropertyPool) error {
	winResultLimiter.Config = cfg.(*WinResultLimiterConfig)
	winResultLimiter.Config.ComponentType = WinResultModifierTypeName

	winResultLimiter.Config.Type = parseWinResultLimiterType(winResultLimiter.Config.StrType)

	winResultLimiter.onInit(&winResultLimiter.Config.BasicComponentConfig)

	return nil
}

// playgame
func (winResultLimiter *WinResultLimiter) onMaxOnLine(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd *WinResultLimiterData) (string, error) {

	mapLinesWin := make(map[int][]int)

	for _, cn := range winResultLimiter.Config.SrcComponents {
		// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
		if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
			continue
		}

		ccd := gameProp.GetComponentDataWithName(cn)
		lst := ccd.GetResults()
		for _, ri := range lst {
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

// playgame
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

// OnAsciiGame - outpur to asciigame
func (winResultModifier *WinResultLimiter) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	std := icd.(*WinResultLimiterData)

	fmt.Printf("WinResultLimiter, ending wins = %v \n", std.Wins)

	return nil
}

// NewComponentData -
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
