package lowcode

import (
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

const SymbolValWinsTypeName = "symbolValWins"

type SymbolValWinsType int

const (
	SVWTypeNormal    SymbolValWinsType = 0
	SVWTypeCollector SymbolValWinsType = 1
)

func parseSymbolValWinsType(strType string) SymbolValWinsType {
	if strType == "collector" {
		return SVWTypeCollector
	}

	return SVWTypeNormal
}

const (
	SVWDVWins      string = "wins"      // 中奖的数值，线注的倍数
	SVWDVSymbolNum string = "symbolNum" // 符号数量
	// SVWDVCollectorNum string = "collectorNum" // 收集器数量
)

type SymbolValWinsData struct {
	BasicComponentData
	SymbolNum int
	Wins      int
	// CollectorNum int
}

// OnNewGame -
func (symbolValWinsData *SymbolValWinsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	symbolValWinsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (symbolValWinsData *SymbolValWinsData) onNewStep() {
	// symbolValWinsData.BasicComponentData.OnNewStep(gameProp, component)

	symbolValWinsData.UsedResults = nil
	symbolValWinsData.SymbolNum = 0
	symbolValWinsData.Wins = 0

	// symbolValWinsData.SymbolNum = 0
	// symbolValWinsData.Wins = 0
	// symbolValWinsData.CollectorNum = 0
}

// Clone
func (symbolValWinsData *SymbolValWinsData) Clone() IComponentData {
	target := &SymbolValWinsData{
		BasicComponentData: symbolValWinsData.CloneBasicComponentData(),
		SymbolNum:          symbolValWinsData.SymbolNum,
		Wins:               symbolValWinsData.Wins,
	}

	return target
}

// BuildPBComponentData
func (symbolValWinsData *SymbolValWinsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolValWinsData{
		BasicComponentData: symbolValWinsData.BuildPBBasicComponentData(),
	}

	if !gIsReleaseMode {
		pbcd.SymbolNum = int32(symbolValWinsData.SymbolNum)
		pbcd.Wins = int32(symbolValWinsData.Wins)
		// pbcd.CollectorNum = int32(symbolValWinsData.CollectorNum)
	}

	return pbcd
}

// GetValEx -
func (symbolValWinsData *SymbolValWinsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == SVWDVSymbolNum {
		return symbolValWinsData.SymbolNum, true
	} else if key == SVWDVWins {
		return symbolValWinsData.Wins, true
	} else if key == CVResultNum {
		return len(symbolValWinsData.UsedResults), true
	}

	return 0, false
}

// // SetVal -
// func (symbolValWinsData *SymbolValWinsData) SetVal(key string, val int) {
// 	if key == CVWins {
// 		symbolValWinsData.Wins = val
// 	} else if key == CVSymbolNum {
// 		symbolValWinsData.SymbolNum = val
// 	}
// }

// SymbolValWinsConfig - configuration for SymbolValWins
type SymbolValWinsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetTypeString        string            `yaml:"betType" json:"betType"`   // bet or totalBet or noPay
	BetType              BetType           `yaml:"-" json:"-"`               // bet or totalBet or noPay
	WinMulti             int               `yaml:"winMulti" json:"winMulti"` // bet or totalBet
	Symbols              []string          `yaml:"symbols" json:"symbols"`   // like collect
	SymbolCodes          []int             `json:"-"`                        //
	StrType              string            `yaml:"type" json:"type"`
	Type                 SymbolValWinsType `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *SymbolValWinsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type SymbolValWins struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValWinsConfig `json:"config"`
}

// Init -
func (symbolValWins *SymbolValWins) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolValWins.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SymbolValWinsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolValWins.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return symbolValWins.InitEx(cfg, pool)
}

// InitEx -
func (symbolValWins *SymbolValWins) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolValWins.Config = cfg.(*SymbolValWinsConfig)
	symbolValWins.Config.ComponentType = SymbolValWinsTypeName

	symbolValWins.Config.BetType = ParseBetType(symbolValWins.Config.BetTypeString)
	symbolValWins.Config.Type = parseSymbolValWinsType(symbolValWins.Config.StrType)

	for _, s := range symbolValWins.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolValWins.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		symbolValWins.Config.SymbolCodes = append(symbolValWins.Config.SymbolCodes, sc)
	}

	symbolValWins.onInit(&symbolValWins.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolValWins *SymbolValWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolValWins.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	svwd := icd.(*SymbolValWinsData)
	svwd.onNewStep()

	// gs := symbolValWins.GetTargetScene3(gameProp, curpr, prs, &svwd.BasicComponentData, symbolValWins.Name, "", 0)
	// isTrigger := true
	// symbolnum := 0

	// svwd.SymbolNum = 0
	// svwd.Wins = 0
	// svwd.CollectorNum = 0

	// if symbolValWins.TriggerSymbolCode >= 0 {
	// 	isTrigger = false

	// 	// if symbolValWins.Config.Type == WinTypeCountScatter {
	// 	ret := sgc7game.CalcScatterEx(gs, symbolValWins.TriggerSymbolCode, symbolValWins.Config.MinNum, func(scatter int, cursymbol int) bool {
	// 		return cursymbol == scatter
	// 	})

	// 	if ret != nil {
	// 		isTrigger = true

	// 		symbolnum = ret.SymbolNums
	// 	}
	// 	// }
	// }

	// if isTrigger {
	// svwd.CollectorNum = symbolnum

	var gs *sgc7game.GameScene
	os := symbolValWins.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	if os != nil {
		collectorpos := []int{}
		mul := 0
		if symbolValWins.Config.Type == SVWTypeCollector {
			gs = symbolValWins.GetTargetScene3(gameProp, curpr, prs, 0)
			if gs != nil {
				for x, arr := range gs.Arr {
					for y, s := range arr {
						if goutils.IndexOfIntSlice(symbolValWins.Config.SymbolCodes, s, 0) >= 0 {
							mul++

							collectorpos = append(collectorpos, x, y)
						}
					}
				}
			}
		} else {
			mul = 1
		}

		totalvals := 0
		pos := make([]int, 0, len(os.Arr)*len(os.Arr[0])*2)

		for x := 0; x < len(os.Arr); x++ {
			for y := 0; y < len(os.Arr[x]); y++ {
				if os.Arr[x][y] > 0 {
					totalvals += os.Arr[x][y]
					pos = append(pos, x, y)

					svwd.SymbolNum++
				}
			}
		}

		if totalvals > 0 && mul > 0 {
			bet := gameProp.GetBet3(stake, symbolValWins.Config.BetType)
			othermul := symbolValWins.GetWinMulti(&svwd.BasicComponentData)

			for i := 0; i < mul; i++ {
				newpos := make([]int, 0, len(pos)+2)

				if symbolValWins.Config.Type == SVWTypeCollector {
					newpos = append(newpos, collectorpos[i*2], collectorpos[i*2+1])
				}

				newpos = append(newpos, pos...)

				ret := &sgc7game.Result{
					// Symbol:     gs.Arr[newpos[0]][newpos[1]],
					Type:       sgc7game.RTCoins,
					LineIndex:  -1,
					Pos:        newpos,
					SymbolNums: len(pos) / 2,
					Mul:        1,
				}

				if gs != nil {
					ret.Symbol = gs.Arr[newpos[0]][newpos[1]]
				}

				ret.CoinWin = totalvals * othermul
				ret.CashWin = ret.CoinWin * bet
				ret.OtherMul = othermul

				svwd.Wins += ret.CoinWin

				symbolValWins.AddResult(curpr, ret, &svwd.BasicComponentData)
			}

			// ret := &sgc7game.Result{
			// 	Symbol:     -1, //gs.Arr[pos[0]][pos[1]],
			// 	Type:       sgc7game.RTSymbolVal,
			// 	LineIndex:  -1,
			// 	Pos:        pos,
			// 	SymbolNums: len(pos) / 2,
			// 	Mul:        mul,
			// }

			// bet := gameProp.GetBet2(stake, symbolValWins.Config.BetType)

			// othermul := symbolValWins.GetWinMulti(&svwd.BasicComponentData)

			// ret.CoinWin = totalvals * mul * othermul
			// ret.CashWin = ret.CoinWin * bet
			// ret.OtherMul = othermul

			// svwd.Wins = ret.CoinWin

			// symbolValWins.AddResult(curpr, ret, &svwd.BasicComponentData)

			nc := symbolValWins.onStepEnd(gameProp, curpr, gp, "")

			return nc, nil
		}
	}
	// }

	nc := symbolValWins.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (symbolValWins *SymbolValWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*SymbolValWinsData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// // OnStatsWithPB -
// func (symbolValWins *SymbolValWins) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	return 0, nil
// 	// pbcd, isok := pbComponentData.(*sgc7pb.SymbolValWinsData)
// 	// if !isok {
// 	// 	goutils.Error("SymbolValWins.OnStatsWithPB",
// 	// 		goutils.Err(ErrIvalidProto))

// 	// 	return 0, ErrIvalidProto
// 	// }

// 	// return symbolValWins.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// // OnStats
// func (symbolValWins *SymbolValWins) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// 	// wins := int64(0)
// 	// isTrigger := false

// 	// for _, v := range lst {
// 	// 	gp, isok := v.CurGameModParams.(*GameParams)
// 	// 	if isok {
// 	// 		curComponent, isok := gp.MapComponentMsgs[symbolValWins.Name]
// 	// 		if isok {
// 	// 			curwins, err := symbolValWins.OnStatsWithPB(feature, curComponent, v)
// 	// 			if err != nil {
// 	// 				goutils.Error("SymbolValWins.OnStats",
// 	// 					goutils.Err(err))

// 	// 				continue
// 	// 			}

// 	// 			isTrigger = true
// 	// 			wins += curwins
// 	// 		}
// 	// 	}
// 	// }

// 	// feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

// 	// if feature.Parent != nil {
// 	// 	totalwins := int64(0)

// 	// 	for _, v := range lst {
// 	// 		totalwins += v.CashWin
// 	// 	}

// 	// 	feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
// 	// }

// 	// return isTrigger, stake.CashBet, wins
// }

// NewComponentData -
func (symbolValWins *SymbolValWins) NewComponentData() IComponentData {
	return &SymbolValWinsData{}
}

// NewStats2 -
func (symbolValWins *SymbolValWins) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, stats2.Options{stats2.OptWins})
}

// OnStats2
func (symbolValWins *SymbolValWins) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult, isOnStepEnd bool) {
	symbolValWins.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr, isOnStepEnd)

	svwd := icd.(*SymbolValWinsData)

	s2.ProcStatsWins(symbolValWins.Name, int64(svwd.Wins))
}

func (symbolValWins *SymbolValWins) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(CCVWinMulti)
	if isok {
		return winMulti
	}

	return symbolValWins.Config.WinMulti
}

func NewSymbolValWins(name string) IComponent {
	return &SymbolValWins{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonSymbolValWins struct {
	BetType  string   `json:"betType"`  // bet or totalBet or noPay
	WinMulti int      `json:"winMulti"` // bet or totalBet
	Symbols  []string `json:"symbols"`  // like collect
	Type     string   `yaml:"type" json:"type"`
}

func (jcfg *jsonSymbolValWins) build() *SymbolValWinsConfig {
	cfg := &SymbolValWinsConfig{
		BetTypeString: jcfg.BetType,
		WinMulti:      jcfg.WinMulti,
		Symbols:       jcfg.Symbols,
		StrType:       jcfg.Type,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseSymbolValWins(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolValWins:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolValWins:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSymbolValWins{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolValWins:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolValWinsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
