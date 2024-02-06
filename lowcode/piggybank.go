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

const PiggyBankTypeName = "piggyBank"

type PiggyBankData struct {
	BasicComponentData
	SavedMoney int
	Wins       int
	WinMulti   int
}

// OnNewGame -
func (piggyBankData *PiggyBankData) OnNewGame(gameProp *GameProperty, component IComponent) {
	piggyBankData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (piggyBankData *PiggyBankData) OnNewStep(gameProp *GameProperty, component IComponent) {
	piggyBankData.BasicComponentData.OnNewStep(gameProp, component)

	piggyBankData.Wins = 0
	piggyBankData.WinMulti = 1
}

// BuildPBComponentData
func (piggyBankData *PiggyBankData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.PiggyBankData{
		BasicComponentData: piggyBankData.BuildPBBasicComponentData(),
		Wins:               int32(piggyBankData.Wins),
		WinMulti:           int32(piggyBankData.WinMulti),
		SavedMoney:         int32(piggyBankData.SavedMoney),
	}

	return pbcd
}

// GetVal -
func (piggyBankData *PiggyBankData) GetVal(key string) int {
	if key == CVWins {
		return piggyBankData.Wins
	}

	return 0
}

// SetVal -
func (piggyBankData *PiggyBankData) SetVal(key string, val int) {
	if key == CVWins {
		piggyBankData.Wins = val
	}
}

// PiggyBankConfig - configuration for PiggyBank
type PiggyBankConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	WinMulti             int `yaml:"winMulti" json:"winMulti"` // winMulti，最后的中奖倍数，默认为1
}

// SetLinkComponent
func (cfg *PiggyBankConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type PiggyBank struct {
	*BasicComponent `json:"-"`
	Config          *PiggyBankConfig `json:"config"`
}

// Init -
func (piggyBank *PiggyBank) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("PiggyBank.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &PiggyBankConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("PiggyBank.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return piggyBank.InitEx(cfg, pool)
}

// InitEx -
func (piggyBank *PiggyBank) InitEx(cfg any, pool *GamePropertyPool) error {
	piggyBank.Config = cfg.(*PiggyBankConfig)
	piggyBank.Config.ComponentType = PiggyBankTypeName

	if piggyBank.Config.WinMulti <= 0 {
		piggyBank.Config.WinMulti = 1
	}

	piggyBank.onInit(&piggyBank.Config.BasicComponentConfig)

	return nil
}

// playgame
func (piggyBank *PiggyBank) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// winResultMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*PiggyBankData)

	cd.OnNewStep(gameProp, piggyBank)

	winMulti := piggyBank.GetWinMulti(&cd.BasicComponentData)

	cd.WinMulti = winMulti
	sm, isok := cd.GetConfigIntVal(CCVSavedMoney)
	if !isok {
		nc := piggyBank.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	cd.Wins = sm * winMulti

	bet := gameProp.GetBet2(stake, BTypeBet)

	ret := &sgc7game.Result{
		Symbol:    -1,
		Type:      sgc7game.RTSymbolVal,
		LineIndex: -1,
		CoinWin:   cd.Wins,
		CashWin:   cd.Wins * bet,
	}

	piggyBank.AddResult(curpr, ret, &cd.BasicComponentData)

	nc := piggyBank.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (piggyBank *PiggyBank) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	std := icd.(*PiggyBankData)

	fmt.Printf("PiggyBank x %v, ending wins = %v \n", std.WinMulti, std.Wins)

	return nil
}

// OnStatsWithPB -
func (piggyBank *PiggyBank) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.PiggyBankData)
	if !isok {
		goutils.Error("PiggyBank.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return piggyBank.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// OnStats
func (piggyBank *PiggyBank) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (piggyBank *PiggyBank) NewComponentData() IComponentData {
	return &PiggyBankData{}
}

func (piggyBank *PiggyBank) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(WTCVWinMulti)
	if isok {
		return winMulti
	}

	return piggyBank.Config.WinMulti
}

func NewPiggyBank(name string) IComponent {
	return &PiggyBank{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonPiggyBank struct {
	WinMulti int `json:"winMulti"`
}

func (jwt *jsonPiggyBank) build() *PiggyBankConfig {
	cfg := &PiggyBankConfig{
		WinMulti: jwt.WinMulti,
	}

	cfg.UseSceneV3 = true

	return cfg
}

func parsePiggyBank(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parsePiggyBank:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parsePiggyBank:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonPiggyBank{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parsePiggyBank:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: PiggyBankTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
