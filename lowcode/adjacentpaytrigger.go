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

const AdjacentPayTriggerTypeName = "adjacentPayTrigger"

type AdjacentPayTriggerData struct {
	BasicComponentData
	PosComponentData
	NextComponent     string
	SymbolNum         int
	WildNum           int
	RespinNum         int
	Wins              int
	WinMulti          int
	AvgSymbolValMulti int // 平均的symbolVal倍数，用整数来表达浮点数，100是1倍
}

// OnNewGame -
func (adjacentPayTriggerData *AdjacentPayTriggerData) OnNewGame(gameProp *GameProperty, component IComponent) {
	adjacentPayTriggerData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (adjacentPayTriggerData *AdjacentPayTriggerData) onNewStep() {
	adjacentPayTriggerData.UsedResults = nil
	adjacentPayTriggerData.NextComponent = ""
	adjacentPayTriggerData.SymbolNum = 0
	adjacentPayTriggerData.WildNum = 0
	adjacentPayTriggerData.RespinNum = 0
	adjacentPayTriggerData.Wins = 0
	adjacentPayTriggerData.WinMulti = 1

	if !gIsReleaseMode {
		adjacentPayTriggerData.PosComponentData.ClearPos()
	}
}

// Clone
func (adjacentPayTriggerData *AdjacentPayTriggerData) Clone() IComponentData {
	target := &AdjacentPayTriggerData{
		BasicComponentData: adjacentPayTriggerData.CloneBasicComponentData(),
		NextComponent:      adjacentPayTriggerData.NextComponent,
		SymbolNum:          adjacentPayTriggerData.SymbolNum,
		WildNum:            adjacentPayTriggerData.WildNum,
		RespinNum:          adjacentPayTriggerData.RespinNum,
		Wins:               adjacentPayTriggerData.Wins,
		WinMulti:           adjacentPayTriggerData.WinMulti,
	}

	if !gIsReleaseMode {
		target.PosComponentData = adjacentPayTriggerData.PosComponentData.Clone()
	}

	return target
}

// BuildPBComponentData
func (adjacentPayTriggerData *AdjacentPayTriggerData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.AdjacentPayTriggerData{
		BasicComponentData: adjacentPayTriggerData.BuildPBBasicComponentData(),
		NextComponent:      adjacentPayTriggerData.NextComponent,
		SymbolNum:          int32(adjacentPayTriggerData.SymbolNum),
		WildNum:            int32(adjacentPayTriggerData.WildNum),
		RespinNum:          int32(adjacentPayTriggerData.RespinNum),
		Wins:               int32(adjacentPayTriggerData.Wins),
		WinMulti:           int32(adjacentPayTriggerData.WinMulti),
	}

	return pbcd
}

// GetValEx -
func (adjacentPayTriggerData *AdjacentPayTriggerData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolNum {
		return adjacentPayTriggerData.SymbolNum, true
	} else if key == CVWildNum {
		return adjacentPayTriggerData.WildNum, true
	} else if key == CVRespinNum {
		return adjacentPayTriggerData.RespinNum, true
	} else if key == CVWins {
		return adjacentPayTriggerData.Wins, true
	} else if key == CVAvgSymbolValMulti {
		if adjacentPayTriggerData.AvgSymbolValMulti == 0 {
			return 100, true
		}

		return adjacentPayTriggerData.AvgSymbolValMulti, true
	} else if key == CVResultNum || key == CVWinResultNum {
		return len(adjacentPayTriggerData.UsedResults), true
	}

	return 0, false
}

// GetPos -
func (adjacentPayTriggerData *AdjacentPayTriggerData) GetPos() []int {
	return adjacentPayTriggerData.Pos
}

// AddPos -
func (adjacentPayTriggerData *AdjacentPayTriggerData) AddPos(x, y int) {
	adjacentPayTriggerData.PosComponentData.Add(x, y)
}

// AdjacentPayTriggerConfig - configuration for AdjacentPayTrigger
// 需要特别注意，当判断scatter时，symbols里的符号会当作同一个符号来处理
type AdjacentPayTriggerConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string            `yaml:"symbols" json:"symbols"`                       // like scatter
	SymbolCodes          []int               `yaml:"-" json:"-"`                                   // like scatter
	Type                 string              `yaml:"type" json:"type"`                             // like scatters
	TriggerType          SymbolTriggerType   `yaml:"-" json:"-"`                                   // SymbolTriggerType
	OSMulTypeString      string              `yaml:"symbolValsMulti" json:"symbolValsMulti"`       // OtherSceneMultiType
	OSMulType            OtherSceneMultiType `yaml:"-" json:"-"`                                   // OtherSceneMultiType
	BetTypeString        string              `yaml:"betType" json:"betType"`                       // bet or totalBet or noPay
	BetType              BetType             `yaml:"-" json:"-"`                                   // bet or totalBet or noPay
	MinNum               int                 `yaml:"minNum" json:"minNum"`                         // like 3，countscatter 或 countscatterInArea 或 checkLines 或 checkWays 时生效
	WildSymbols          []string            `yaml:"wildSymbols" json:"wildSymbols"`               // wild etc
	WildSymbolCodes      []int               `yaml:"-" json:"-"`                                   // wild symbolCode
	WinMulti             int                 `yaml:"winMulti" json:"winMulti"`                     // winMulti，最后的中奖倍数，默认为1
	JumpToComponent      string              `yaml:"jumpToComponent" json:"jumpToComponent"`       // jump to
	ForceToNext          bool                `yaml:"forceToNext" json:"forceToNext"`               // 如果触发，默认跳转jump to，这里可以强制走next分支
	Awards               []*Award            `yaml:"awards" json:"awards"`                         // 新的奖励系统
	IsReverse            bool                `yaml:"isReverse" json:"isReverse"`                   // 如果isReverse，表示判定为否才触发
	PiggyBankComponent   string              `yaml:"piggyBankComponent" json:"piggyBankComponent"` // piggyBank component
}

// SetLinkComponent
func (cfg *AdjacentPayTriggerConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type AdjacentPayTrigger struct {
	*BasicComponent `json:"-"`
	Config          *AdjacentPayTriggerConfig `json:"config"`
}

// Init -
func (adjacentPayTrigger *AdjacentPayTrigger) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("AdjacentPayTrigger.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &AdjacentPayTriggerConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("AdjacentPayTrigger.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return adjacentPayTrigger.InitEx(cfg, pool)
}

// InitEx -
func (adjacentPayTrigger *AdjacentPayTrigger) InitEx(cfg any, pool *GamePropertyPool) error {
	adjacentPayTrigger.Config = cfg.(*AdjacentPayTriggerConfig)
	adjacentPayTrigger.Config.ComponentType = AdjacentPayTriggerTypeName

	adjacentPayTrigger.Config.OSMulType = ParseOtherSceneMultiType(adjacentPayTrigger.Config.OSMulTypeString)

	for _, s := range adjacentPayTrigger.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("AdjacentPayTrigger.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))
		}

		adjacentPayTrigger.Config.SymbolCodes = append(adjacentPayTrigger.Config.SymbolCodes, sc)
	}

	for _, s := range adjacentPayTrigger.Config.WildSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("AdjacentPayTrigger.InitEx:WildSymbols",
				slog.String("symbol", s),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		adjacentPayTrigger.Config.WildSymbolCodes = append(adjacentPayTrigger.Config.WildSymbolCodes, sc)
	}

	stt := ParseSymbolTriggerType(adjacentPayTrigger.Config.Type)
	if stt == STTypeUnknow {
		goutils.Error("AdjacentPayTrigger.InitEx:ParseSymbolTriggerType",
			slog.String("SymbolTriggerType", adjacentPayTrigger.Config.Type),
			goutils.Err(ErrInvalidSymbolTriggerType))

		return ErrInvalidSymbolTriggerType
	}

	adjacentPayTrigger.Config.TriggerType = stt

	adjacentPayTrigger.Config.BetType = ParseBetType(adjacentPayTrigger.Config.BetTypeString)

	for _, award := range adjacentPayTrigger.Config.Awards {
		award.Init()
	}

	if adjacentPayTrigger.Config.WinMulti <= 0 {
		adjacentPayTrigger.Config.WinMulti = 1
	}

	adjacentPayTrigger.onInit(&adjacentPayTrigger.Config.BasicComponentConfig)

	return nil
}

func (adjacentPayTrigger *AdjacentPayTrigger) calcSymbolValMulti(ret *sgc7game.Result, os *sgc7game.GameScene, funcCalcMulti sgc7game.FuncCalcMulti) int {
	mul := 1

	for i := 0; i < len(ret.Pos)/2; i++ {
		x := ret.Pos[i*2]
		y := ret.Pos[i*2+1]

		mul = funcCalcMulti(mul, os.Arr[x][y])
	}

	return mul
}

// procWins
func (adjacentPayTrigger *AdjacentPayTrigger) procWins(gameProp *GameProperty, curpr *sgc7game.PlayResult, std *AdjacentPayTriggerData, lst []*sgc7game.Result, os *sgc7game.GameScene) (int, error) {
	if adjacentPayTrigger.Config.BetType == BTypeNoPay {
		for _, v := range lst {
			v.CoinWin = 0
			v.CashWin = 0

			adjacentPayTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds

			if !gIsReleaseMode {
				std.MergePosList(v.Pos)
			}
		}

		return 0, nil
	}

	std.WinMulti = adjacentPayTrigger.GetWinMulti(&std.BasicComponentData)

	if adjacentPayTrigger.Config.OSMulType == OSMTNone || os == nil {
		for _, v := range lst {
			v.OtherMul = std.WinMulti

			v.CoinWin *= std.WinMulti
			v.CashWin *= std.WinMulti

			std.Wins += v.CoinWin

			adjacentPayTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds

			if !gIsReleaseMode {
				std.MergePosList(v.Pos)
			}
		}
	} else {
		funcCalcMulti := GetSymbolValMultiFunc(adjacentPayTrigger.Config.OSMulType)

		if !gIsReleaseMode {
			std.AvgSymbolValMulti = 0
		}

		for _, v := range lst {
			svm := adjacentPayTrigger.calcSymbolValMulti(v, os, funcCalcMulti)

			if !gIsReleaseMode {
				std.AvgSymbolValMulti += svm
			}

			v.OtherMul = std.WinMulti * svm

			v.CoinWin *= v.OtherMul
			v.CashWin *= v.OtherMul

			std.Wins += v.CoinWin

			adjacentPayTrigger.AddResult(curpr, v, &std.BasicComponentData)

			std.SymbolNum += v.SymbolNums
			std.WildNum += v.Wilds

			if !gIsReleaseMode {
				std.MergePosList(v.Pos)
			}
		}

		if !gIsReleaseMode {
			std.AvgSymbolValMulti = std.AvgSymbolValMulti * 100 / len(lst)
		}
	}

	if std.Wins > 0 {
		if adjacentPayTrigger.Config.PiggyBankComponent != "" {
			cd := gameProp.GetCurComponentDataWithName(adjacentPayTrigger.Config.PiggyBankComponent)
			if cd == nil {
				goutils.Error("AdjacentPayTrigger.procWins:GetCurComponentDataWithName",
					slog.String("PiggyBankComponent", adjacentPayTrigger.Config.PiggyBankComponent),
					goutils.Err(ErrInvalidComponent))

				return 0, ErrInvalidComponent
			}

			for _, v := range lst {
				v.IsNoPayNow = true
			}

			cd.ChgConfigIntVal(CCVSavedMoney, std.Wins)

			gameProp.UseComponent(adjacentPayTrigger.Config.PiggyBankComponent)
		}
	}

	return std.Wins, nil
}

// OnProcControllers -
func (adjacentPayTrigger *AdjacentPayTrigger) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(adjacentPayTrigger.Config.Awards) > 0 {
		gameProp.procAwards(plugin, adjacentPayTrigger.Config.Awards, curpr, gp)
	}
}

// playgame
func (adjacentPayTrigger *AdjacentPayTrigger) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	std := cd.(*AdjacentPayTriggerData)
	std.onNewStep()

	gs := adjacentPayTrigger.GetTargetScene3(gameProp, curpr, prs, 0)
	os := adjacentPayTrigger.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	isTrigger, lst := adjacentPayTrigger.CanTriggerWithScene(gameProp, gs, curpr, stake, cd)

	if isTrigger {
		adjacentPayTrigger.procWins(gameProp, curpr, std, lst, os)

		adjacentPayTrigger.ProcControllers(gameProp, plugin, curpr, gp, -1, "")
		// if len(adjacentPayTrigger.Config.Awards) > 0 {
		// 	gameProp.procAwards(plugin, adjacentPayTrigger.Config.Awards, curpr, gp)
		// }

		if adjacentPayTrigger.Config.JumpToComponent != "" {
			std.NextComponent = adjacentPayTrigger.Config.JumpToComponent

			nc := adjacentPayTrigger.onStepEnd(gameProp, curpr, gp, std.NextComponent)

			return nc, nil
		}

		nc := adjacentPayTrigger.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := adjacentPayTrigger.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (adjacentPayTrigger *AdjacentPayTrigger) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	std := cd.(*AdjacentPayTriggerData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(std.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	if std.NextComponent != "" {
		fmt.Printf("%v triggered, jump to %v \n", adjacentPayTrigger.Name, std.NextComponent)
	}

	return nil
}

// NewComponentData -
func (adjacentPayTrigger *AdjacentPayTrigger) NewComponentData() IComponentData {
	return &AdjacentPayTriggerData{}
}

func (adjacentPayTrigger *AdjacentPayTrigger) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return adjacentPayTrigger.Config.WinMulti
}

// NewStats2 -
func (adjacentPayTrigger *AdjacentPayTrigger) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (adjacentPayTrigger *AdjacentPayTrigger) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	adjacentPayTrigger.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	cd := icd.(*AdjacentPayTriggerData)

	s2.ProcStatsWins(adjacentPayTrigger.Name, int64(cd.Wins))
}

// GetAllLinkComponents - get all link components
func (adjacentPayTrigger *AdjacentPayTrigger) GetAllLinkComponents() []string {
	return []string{adjacentPayTrigger.Config.DefaultNextComponent, adjacentPayTrigger.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (adjacentPayTrigger *AdjacentPayTrigger) GetNextLinkComponents() []string {
	return []string{adjacentPayTrigger.Config.DefaultNextComponent, adjacentPayTrigger.Config.JumpToComponent}
}

func (adjacentPayTrigger *AdjacentPayTrigger) getSymbols(gameProp *GameProperty) []int {
	s := gameProp.GetCurCallStackSymbol()
	if s >= 0 {
		return []int{s}
	}

	return adjacentPayTrigger.Config.SymbolCodes
}

// CanTriggerWithScene -
func (adjacentPayTrigger *AdjacentPayTrigger) CanTriggerWithScene(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake, icd IComponentData) (bool, []*sgc7game.Result) {
	isTrigger := false
	lst := []*sgc7game.Result{}

	if adjacentPayTrigger.Config.TriggerType == STTypeAdjacentPay {

		symbols := adjacentPayTrigger.getSymbols(gameProp)

		currets, err := sgc7game.CalcAdjacentPay(gs, gameProp.CurPaytables, gameProp.GetBet3(stake, adjacentPayTrigger.Config.BetType),
			func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(symbols, cursymbol, 0) >= 0
			}, func(cursymbol int) bool {
				return goutils.IndexOfIntSlice(adjacentPayTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int, startsymbol int) bool {
				if cursymbol == startsymbol {
					return true
				}

				return goutils.IndexOfIntSlice(adjacentPayTrigger.Config.WildSymbolCodes, cursymbol, 0) >= 0
			}, func(cursymbol int) int {
				return cursymbol
			})
		if err != nil {
			goutils.Error("AdjacentPayTrigger.CanTriggerWithScene:CalcAdjacentPay",
				goutils.Err(err))

			return false, nil
		}

		lst = append(lst, currets...)

		if len(lst) > 0 {
			isTrigger = true
		}
	}

	if adjacentPayTrigger.Config.IsReverse {
		isTrigger = !isTrigger
	}

	return isTrigger, lst
}

func NewAdjacentPayTrigger(name string) IComponent {
	return &AdjacentPayTrigger{
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
type jsonAdjacentPayTrigger struct {
	Symbols             []string `json:"symbols"`
	BetType             string   `json:"betType"`
	SymbolValsMulti     string   `json:"symbolValsMulti"`
	MinNum              int      `json:"minNum"`
	WildSymbols         []string `json:"wildSymbols"`
	WinMulti            int      `json:"winMulti"`
	PutMoneyInPiggyBank string   `json:"putMoneyInPiggyBank"`
}

func (jcfg *jsonAdjacentPayTrigger) build() *AdjacentPayTriggerConfig {
	cfg := &AdjacentPayTriggerConfig{
		Symbols:            jcfg.Symbols,
		Type:               "adjacentpay",
		BetTypeString:      jcfg.BetType,
		OSMulTypeString:    jcfg.SymbolValsMulti,
		MinNum:             jcfg.MinNum,
		WildSymbols:        jcfg.WildSymbols,
		WinMulti:           jcfg.WinMulti,
		PiggyBankComponent: jcfg.PutMoneyInPiggyBank,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseAdjacentPayTrigger(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseAdjacentPayTrigger:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseAdjacentPayTrigger:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonAdjacentPayTrigger{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseAdjacentPayTrigger:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseAdjacentPayTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: AdjacentPayTriggerTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
