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

const WinResultModifierTypeName = "winResultModifier"

type WinResultModifierType int

func (wrmt WinResultModifierType) isNeedMultiply() bool {
	return wrmt != WRMTypeDivide
}

func (wrmt WinResultModifierType) isNeedGameScene() bool {
	return wrmt == WRMTypeExistSymbol ||
		wrmt == WRMTypeAddSymbolMulti ||
		wrmt == WRMTypeMulSymbolMulti ||
		wrmt == WRMTypeSymbolMultiOnWays
}

func (wrmt WinResultModifierType) isValidInWinResultModifierEx() bool {
	return wrmt == WRMTypeAddSymbolMulti ||
		wrmt == WRMTypeMulSymbolMulti ||
		wrmt == WRMTypeSymbolMultiOnWays
}

const (
	WRMTypeExistSymbol       WinResultModifierType = 0
	WRMTypeAddSymbolMulti    WinResultModifierType = 1
	WRMTypeMulSymbolMulti    WinResultModifierType = 2
	WRMTypeSymbolMultiOnWays WinResultModifierType = 3
	WRMTypeDivide            WinResultModifierType = 4
	WRMTypeMultiply          WinResultModifierType = 5
)

func parseWinResultModifierType(str string) WinResultModifierType {
	switch str {
	case "existsymbol":
		return WRMTypeExistSymbol
	case "addsymbolmulti":
		return WRMTypeAddSymbolMulti
	case "mulsymbolmulti":
		return WRMTypeMulSymbolMulti
	case "symbolmultionways":
		return WRMTypeSymbolMultiOnWays
	case "divide":
		return WRMTypeDivide
	case "multiply":
		return WRMTypeMultiply
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

// GetValEx -
func (winResultModifierData *WinResultModifierData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVWins {
		return winResultModifierData.Wins, true
	}

	return 0, false
}

// WinResultModifierConfig - configuration for WinResultModifier
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type WinResultModifierConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                `yaml:"type" json:"type"`                         // type
	Type                 WinResultModifierType `yaml:"-" json:"-"`                               // type
	SourceComponents     []string              `yaml:"sourceComponents" json:"sourceComponents"` // target components
	WinMulti             int                   `yaml:"winMulti" json:"winMulti"`                 // winMulti，最后的中奖倍数，默认为1
	WinDivisor           int                   `yaml:"winDivisor" json:"winDivisor"`             // winDivisor
	TargetSymbols        []string              `yaml:"targetSymbols" json:"targetSymbols"`       // targetSymbols
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

	// yaml只会自动生成，不需要考虑字符串大小写问题
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

	if winResultModifier.Config.WinDivisor <= 0 {
		winResultModifier.Config.WinDivisor = 1
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

	std := icd.(*WinResultModifierData)
	std.onNewStep()

	winMulti := winResultModifier.GetWinMulti(&std.BasicComponentData)

	std.WinMulti = winMulti

	if winMulti == 1 && winResultModifier.Config.Type.isNeedMultiply() {
		nc := winResultModifier.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	gs := winResultModifier.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil && winResultModifier.Config.Type.isNeedGameScene() {
		nc := winResultModifier.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	isproced := false

	if winResultModifier.Config.Type == WRMTypeExistSymbol {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				if HasSymbolsInResult(gs, winResultModifier.Config.TargetSymbolCodes, curpr.Results[ri]) {
					curpr.Results[ri].CashWin *= winMulti
					curpr.Results[ri].CoinWin *= winMulti
					curpr.Results[ri].OtherMul *= winMulti

					std.Wins += curpr.Results[ri].CoinWin

					isproced = true
				}
			}
		}
	} else if winResultModifier.Config.Type == WRMTypeAddSymbolMulti {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				num := CountSymbolsInResult(gs, winResultModifier.Config.TargetSymbolCodes, curpr.Results[ri])
				if num > 0 {
					curpr.Results[ri].CashWin *= winMulti * num
					curpr.Results[ri].CoinWin *= winMulti * num
					curpr.Results[ri].OtherMul *= winMulti * num

					std.Wins += curpr.Results[ri].CoinWin

					isproced = true
				}
			}
		}
	} else if winResultModifier.Config.Type == WRMTypeMulSymbolMulti {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				num := CountSymbolsInResult(gs, winResultModifier.Config.TargetSymbolCodes, curpr.Results[ri])
				if num > 0 {
					m := intPow(winMulti, num)
					curpr.Results[ri].CashWin *= m
					curpr.Results[ri].CoinWin *= m
					curpr.Results[ri].OtherMul *= m

					std.Wins += curpr.Results[ri].CoinWin

					isproced = true
				}
			}
		}
	} else if winResultModifier.Config.Type == WRMTypeSymbolMultiOnWays {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				mul := 1

				for x, arr := range gs.Arr {
					curmul := 0

					for i := 0; i < len(curpr.Results[ri].Pos)/2; i++ {
						if curpr.Results[ri].Pos[i*2] == x {
							if goutils.IndexOfIntSlice(winResultModifier.Config.TargetSymbolCodes, arr[curpr.Results[ri].Pos[i*2+1]], 0) >= 0 {
								curmul += winMulti
							} else {
								curmul++
							}
						}
					}

					if curmul > 0 {
						mul *= curmul
					}
				}

				if curpr.Results[ri].Mul <= 0 {
					goutils.Error("WinResultModifier.OnPlayGame:curpr.Results[ri].Mul <= 0",
						goutils.Err(ErrInvalidComponentConfig))

					return "", ErrInvalidComponentConfig
				}

				curpr.Results[ri].CoinWin = curpr.Results[ri].CoinWin / curpr.Results[ri].Mul * mul
				curpr.Results[ri].CashWin = curpr.Results[ri].CashWin / curpr.Results[ri].Mul * mul

				// 这个地方不确定能修改OtherMul，先注释掉，待查
				// curpr.Results[ri].OtherMul *= mul

				std.Wins += curpr.Results[ri].CoinWin

				isproced = true
			}
		}
	} else if winResultModifier.Config.Type == WRMTypeDivide {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				div := winResultModifier.Config.WinDivisor

				if div <= 0 {
					div = 1
				}

				curpr.Results[ri].CoinWin = curpr.Results[ri].CoinWin / div
				curpr.Results[ri].CashWin = curpr.Results[ri].CashWin / div

				std.Wins += curpr.Results[ri].CoinWin

				isproced = true
			}
		}
	} else if winResultModifier.Config.Type == WRMTypeMultiply {
		for _, cn := range winResultModifier.Config.SourceComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				mul := winMulti

				if mul <= 0 {
					mul = 1
				}

				curpr.Results[ri].CoinWin = curpr.Results[ri].CoinWin * mul
				curpr.Results[ri].CashWin = curpr.Results[ri].CashWin * mul
				curpr.Results[ri].OtherMul *= mul

				std.Wins += curpr.Results[ri].CoinWin

				isproced = true
			}
		}
	}

	if !isproced {
		nc := winResultModifier.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
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

// "type": "divide",
// "winMulti": 1,
// "targetSymbols": [],
// "winDivisor": 10,
// "sourceComponent": [
//
//	"bg-pay"
//
// ]
type jsonWinResultModifier struct {
	Type             string   `json:"type"`            // type
	SourceComponents []string `json:"sourceComponent"` // source components
	WinMulti         int      `json:"winMulti"`        // winMulti
	TargetSymbols    []string `json:"targetSymbols"`   // targetSymbols
	WinDivisor       int      `json:"winDivisor"`      // winDivisor
}

func (jwt *jsonWinResultModifier) build() *WinResultModifierConfig {
	cfg := &WinResultModifierConfig{
		StrType:          strings.ToLower(jwt.Type),
		SourceComponents: jwt.SourceComponents,
		WinMulti:         jwt.WinMulti,
		TargetSymbols:    jwt.TargetSymbols,
		WinDivisor:       jwt.WinDivisor,
	}

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
