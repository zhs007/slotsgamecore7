package lowcode

import (
	"fmt"
	"log/slog"
	"os"
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

const PiggyBankTypeName = "piggyBank"

type PiggyBankType int

const (
	PiggyBankTypeNone                   PiggyBankType = 0
	PiggyBankTypeSumSymbolVals          PiggyBankType = 1
	PiggyBankTypeAddSumSymbolVals       PiggyBankType = 2
	PiggyBankTypeSumEmptyWithSymbolVals PiggyBankType = 3
)

func parsePiggyBankType(str string) PiggyBankType {
	str = strings.ToLower(strings.ReplaceAll(str, " ", ""))
	if str == "winmulti=sum(symbolvals)" {
		return PiggyBankTypeSumSymbolVals
	} else if str == "winmulti+=sum(symbolvals)" {
		return PiggyBankTypeAddSumSymbolVals
	} else if str == "winmulti=sum(emptysymbol(symbolvals))" {
		return PiggyBankTypeSumEmptyWithSymbolVals
	}

	return PiggyBankTypeNone
}

type PiggyBankData struct {
	BasicComponentData
	SavedMoney int
	Wins       int
	WinMulti   int
}

// OnNewGame -
func (piggyBankData *PiggyBankData) OnNewGame(gameProp *GameProperty, component IComponent) {
	piggyBankData.BasicComponentData.OnNewGame(gameProp, component)

	piggyBankData.Wins = 0

	piggyBank, isok := component.(*PiggyBank)
	if isok {
		piggyBankData.WinMulti = piggyBank.Config.WinMulti
		piggyBankData.SetConfigIntVal(CCVWinMulti, piggyBank.Config.WinMulti)
	} else {
		piggyBankData.WinMulti = 1
		piggyBankData.SetConfigIntVal(CCVWinMulti, 1)
	}
}

// Clone
func (piggyBankData *PiggyBankData) Clone() IComponentData {
	target := &PiggyBankData{
		BasicComponentData: piggyBankData.CloneBasicComponentData(),
		Wins:               piggyBankData.Wins,
		WinMulti:           piggyBankData.WinMulti,
		SavedMoney:         piggyBankData.SavedMoney,
	}

	return target
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
func (piggyBankData *PiggyBankData) GetVal(key string) (int, bool) {
	if key == CVWins {
		return piggyBankData.Wins, true
	} else if key == CVWinMulti {
		return piggyBankData.WinMulti, true
	}

	return 0, false
}

// PiggyBankConfig - configuration for PiggyBank
type PiggyBankConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	WinMulti             int           `yaml:"winMulti" json:"winMulti"` // winMulti，最后的中奖倍数，默认为1
	StrType              string        `yaml:"type" json:"type"`         // 如何初始化winmulti
	Type                 PiggyBankType `yaml:"-" json:"-"`               // 如何初始化winmulti
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &PiggyBankConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("PiggyBank.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return piggyBank.InitEx(cfg, pool)
}

// InitEx -
func (piggyBank *PiggyBank) InitEx(cfg any, pool *GamePropertyPool) error {
	piggyBank.Config = cfg.(*PiggyBankConfig)
	piggyBank.Config.ComponentType = PiggyBankTypeName

	piggyBank.Config.Type = parsePiggyBankType(piggyBank.Config.StrType)

	if piggyBank.Config.WinMulti < 0 {
		piggyBank.Config.WinMulti = 0
	}

	piggyBank.onInit(&piggyBank.Config.BasicComponentConfig)

	return nil
}

// playgame
func (piggyBank *PiggyBank) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*PiggyBankData)
	var winMulti int

	if piggyBank.Config.Type == PiggyBankTypeSumSymbolVals {
		os := piggyBank.GetTargetOtherScene3(gameProp, curpr, prs, 0)
		if os != nil {
			winMulti = 0

			for _, arr := range os.Arr {
				for _, v := range arr {
					winMulti += v
				}
			}
		}

		if winMulti == 0 {
			winMulti = 1
		}
	} else if piggyBank.Config.Type == PiggyBankTypeAddSumSymbolVals {
		initWinMulti := 0

		os := piggyBank.GetTargetOtherScene3(gameProp, curpr, prs, 0)
		if os != nil {
			for _, arr := range os.Arr {
				for _, v := range arr {
					winMulti += v
				}
			}
		}

		if initWinMulti == 0 {
			initWinMulti = 1
		}

		winMulti += initWinMulti
	} else if piggyBank.Config.Type == PiggyBankTypeSumEmptyWithSymbolVals {
		gs := piggyBank.GetTargetScene3(gameProp, curpr, prs, 0)

		os := piggyBank.GetTargetOtherScene3(gameProp, curpr, prs, 0)
		if os != nil {
			winMulti = 0

			for x, arr := range os.Arr {
				for y, v := range arr {
					if gs.Arr[x][y] < 0 {
						winMulti += v
					}
				}
			}
		}

		if winMulti == 0 {
			winMulti = 1
		}
	} else {
		winMulti = piggyBank.GetWinMulti(&cd.BasicComponentData)
	}

	cd.WinMulti = winMulti
	sm, isok := cd.GetConfigIntVal(CCVSavedMoney)
	if !isok {
		nc := piggyBank.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	cd.Wins = sm * winMulti

	bet := gameProp.GetBet3(stake, BTypeBet)

	ret := &sgc7game.Result{
		Symbol:    -1,
		Type:      sgc7game.RTCoins,
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

// NewComponentData -
func (piggyBank *PiggyBank) NewComponentData() IComponentData {
	return &PiggyBankData{}
}

func (piggyBank *PiggyBank) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		if winMulti <= 0 {
			return 1
		}

		return winMulti
	}

	if piggyBank.Config.WinMulti <= 0 {
		return 1
	}

	return piggyBank.Config.WinMulti
}

func NewPiggyBank(name string) IComponent {
	return &PiggyBank{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": "winMulti=sum(emptySymbol(symbolVals))",
// "winMulti": 0
type jsonPiggyBank struct {
	WinMulti int    `json:"winMulti"`
	StrType  string `yaml:"type" json:"type"` // 如何初始化winmulti
}

func (jwt *jsonPiggyBank) build() *PiggyBankConfig {
	cfg := &PiggyBankConfig{
		WinMulti: jwt.WinMulti,
		StrType:  jwt.StrType,
	}

	return cfg
}

func parsePiggyBank(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parsePiggyBank:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parsePiggyBank:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonPiggyBank{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parsePiggyBank:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: PiggyBankTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
