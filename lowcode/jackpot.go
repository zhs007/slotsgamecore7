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

const JackpotTypeName = "jackpot"

type JackpotData struct {
	BasicComponentData
	Wins     int
	WinMulti int
}

// OnNewGame -
func (jackpotData *JackpotData) OnNewGame(gameProp *GameProperty, component IComponent) {
	jackpotData.BasicComponentData.OnNewGame(gameProp, component)

	jackpotData.Wins = 0

	piggyBank, isok := component.(*Jackpot)
	if isok {
		jackpotData.WinMulti = piggyBank.Config.WinMulti
		jackpotData.SetConfigIntVal(CCVWinMulti, piggyBank.Config.WinMulti)
	} else {
		jackpotData.WinMulti = 1
		jackpotData.SetConfigIntVal(CCVWinMulti, 1)
	}
}

// Clone
func (jackpotData *JackpotData) Clone() IComponentData {
	target := &JackpotData{
		BasicComponentData: jackpotData.CloneBasicComponentData(),
		Wins:               jackpotData.Wins,
		WinMulti:           jackpotData.WinMulti,
	}

	return target
}

// BuildPBComponentData
func (jackpotData *JackpotData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.JackpotData{
		BasicComponentData: jackpotData.BuildPBBasicComponentData(),
		Wins:               int32(jackpotData.Wins),
		WinMulti:           int32(jackpotData.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (jackpotData *JackpotData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return jackpotData.Wins, true
	} else if key == CVWinMulti {
		return jackpotData.WinMulti, true
	}

	return 0, false
}

// // SetVal -
// func (jackpotData *JackpotData) SetVal(key string, val int) {
// 	if key == CVWins {
// 		jackpotData.Wins = val
// 	} else if key == CVWinMulti {
// 		jackpotData.WinMulti = val
// 	}
// }

// JackpotConfig - configuration for Jackpot
type JackpotConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetTypeString        string  `yaml:"betType" json:"betType"`   // bet or totalBet or noPay
	BetType              BetType `yaml:"-" json:"-"`               // bet or totalBet or noPay
	Wins                 int     `yaml:"wins" json:"wins"`         // wins
	WinMulti             int     `yaml:"winMulti" json:"winMulti"` // winMulti，最后的中奖倍数，默认为1
}

// SetLinkComponent
func (cfg *JackpotConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Jackpot struct {
	*BasicComponent `json:"-"`
	Config          *JackpotConfig `json:"config"`
}

// Init -
func (jackpot *Jackpot) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Jackpot.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &JackpotConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Jackpot.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return jackpot.InitEx(cfg, pool)
}

// InitEx -
func (jackpot *Jackpot) InitEx(cfg any, pool *GamePropertyPool) error {
	jackpot.Config = cfg.(*JackpotConfig)
	jackpot.Config.ComponentType = JackpotTypeName

	jackpot.Config.BetType = ParseBetType(jackpot.Config.BetTypeString)

	if jackpot.Config.WinMulti < 0 {
		jackpot.Config.WinMulti = 0
	}

	jackpot.onInit(&jackpot.Config.BasicComponentConfig)

	return nil
}

// playgame
func (jackpot *Jackpot) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// winResultMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*JackpotData)

	winMulti := jackpot.GetWinMulti(&cd.BasicComponentData)
	wins := jackpot.GetWins(&cd.BasicComponentData)

	cd.WinMulti = winMulti

	cd.Wins = wins * winMulti

	bet := gameProp.GetBet3(stake, jackpot.Config.BetType)

	ret := &sgc7game.Result{
		Symbol:    -1,
		Type:      sgc7game.RTBonus,
		LineIndex: -1,
		CoinWin:   cd.Wins,
		CashWin:   cd.Wins * bet,
	}

	jackpot.AddResult(curpr, ret, &cd.BasicComponentData)

	nc := jackpot.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (jackpot *Jackpot) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*JackpotData)

	fmt.Printf("Jackpot x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// NewComponentData -
func (jackpot *Jackpot) NewComponentData() IComponentData {
	return &JackpotData{}
}

func (jackpot *Jackpot) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		if winMulti <= 0 {
			return 1
		}

		return winMulti
	}

	if jackpot.Config.WinMulti <= 0 {
		return 1
	}

	return jackpot.Config.WinMulti
}

func (jackpot *Jackpot) GetWins(basicCD *BasicComponentData) int {
	wins, isok := basicCD.GetConfigIntVal(CCVWins)
	if isok {
		if wins < 0 {
			return 0
		}

		return wins
	}

	if jackpot.Config.Wins < 0 {
		return 0
	}

	return jackpot.Config.Wins
}

// NewStats2 -
func (jackpot *Jackpot) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (jackpot *Jackpot) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	jackpot.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*JackpotData)

	s2.ProcStatsWins(jackpot.Name, int64(cd.Wins))
}

func NewJackpot(name string) IComponent {
	return &Jackpot{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "betType": "bet",
// "wins": 1000
type jsonJackpot struct {
	WinMulti int    `json:"winMulti"`
	BetType  string `json:"betType"`
	Wins     int    `json:"wins"`
}

func (jwt *jsonJackpot) build() *JackpotConfig {
	cfg := &JackpotConfig{
		WinMulti:      jwt.WinMulti,
		BetTypeString: jwt.BetType,
		Wins:          jwt.Wins,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseJackpot(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseJackpot:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseJackpot:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonJackpot{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseJackpot:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: JackpotTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
