package lowcode

import (
	"os"

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

const SymbolValWinsTypeName = "symbolValWins"

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

// OnNewStep -
func (symbolValWinsData *SymbolValWinsData) OnNewStep(gameProp *GameProperty, component IComponent) {
	symbolValWinsData.BasicComponentData.OnNewStep(gameProp, component)

	// symbolValWinsData.SymbolNum = 0
	// symbolValWinsData.Wins = 0
	// symbolValWinsData.CollectorNum = 0
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

// GetVal -
func (symbolValWinsData *SymbolValWinsData) GetVal(key string) int {
	if key == SVWDVSymbolNum {
		return symbolValWinsData.SymbolNum
	} else if key == SVWDVWins {
		return symbolValWinsData.Wins
	}
	// } else if key == SVWDVCollectorNum {
	// 	return symbolValWinsData.CollectorNum
	// }

	return 0
}

// SetVal -
func (symbolValWinsData *SymbolValWinsData) SetVal(key string, val int) {
	if key == STDVWins {
		symbolValWinsData.Wins = val
	} else if key == SVWDVSymbolNum {
		symbolValWinsData.SymbolNum = val
	}
}

// SymbolValWinsConfig - configuration for SymbolValWins
type SymbolValWinsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	BetTypeString        string  `yaml:"betType" json:"betType"`   // bet or totalBet or noPay
	BetType              BetType `yaml:"-" json:"-"`               // bet or totalBet or noPay
	WinMulti             int     `yaml:"winMulti" json:"winMulti"` // bet or totalBet
	// TriggerSymbol        string `yaml:"triggerSymbol" json:"triggerSymbol"` // like collect
	// TriggerSymbolCode    int    `json:"-"`                                  //
	// Type                    string `yaml:"type" json:"type"`                                       // like scatters
	// MinNum                  int  `yaml:"minNum" json:"minNum"`                                   // like 3
	// IsTriggerSymbolNumMulti bool `yaml:"isTriggerSymbolNumMulti" json:"isTriggerSymbolNumMulti"` // totalwins = totalvals * triggetSymbol's num
}

type SymbolValWins struct {
	*BasicComponent `json:"-"`
	Config          *SymbolValWinsConfig `json:"config"`
	// TriggerSymbolCode int                  `json:"-"`
}

// Init -
func (symbolValWins *SymbolValWins) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolValWins.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolValWinsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolValWins.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolValWins.InitEx(cfg, pool)
}

// InitEx -
func (symbolValWins *SymbolValWins) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolValWins.Config = cfg.(*SymbolValWinsConfig)
	symbolValWins.Config.ComponentType = SymbolValWinsTypeName

	symbolValWins.Config.BetType = ParseBetType(symbolValWins.Config.BetTypeString)

	// if symbolValWins.Config.TriggerSymbol != "" {
	// 	symbolValWins.TriggerSymbolCode = pool.DefaultPaytables.MapSymbols[symbolValWins.Config.TriggerSymbol]
	// } else {
	// 	symbolValWins.TriggerSymbolCode = -1
	// }

	symbolValWins.onInit(&symbolValWins.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolValWins *SymbolValWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolValWins.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	svwd := icd.(*SymbolValWinsData)

	// gs := symbolValWins.GetTargetScene3(gameProp, curpr, prs, &svwd.BasicComponentData, symbolValWins.Name, "", 0)
	// isTrigger := true
	// symbolnum := 0

	svwd.SymbolNum = 0
	svwd.Wins = 0
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

	os := symbolValWins.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	if os != nil {
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

		if totalvals > 0 {
			ret := &sgc7game.Result{
				Symbol:     -1, //gs.Arr[pos[0]][pos[1]],
				Type:       sgc7game.RTSymbolVal,
				LineIndex:  -1,
				Pos:        pos,
				SymbolNums: len(pos) / 2,
			}

			bet := gameProp.GetBet2(stake, symbolValWins.Config.BetType)

			mul := symbolValWins.GetWinMulti(&svwd.BasicComponentData) //1 //gameProp.GetVal(GamePropGameCoinMulti) * gameProp.GetVal(GamePropStepCoinMulti)

			// if symbolValWins.Config.IsTriggerSymbolNumMulti {
			// 	ret.CoinWin = totalvals * symbolnum * mul
			// 	ret.CashWin = ret.CoinWin * bet
			// } else {
			ret.CoinWin = totalvals * mul
			ret.CashWin = ret.CoinWin * bet
			// }

			svwd.Wins = ret.CoinWin

			symbolValWins.AddResult(curpr, ret, &svwd.BasicComponentData)
		}
	}
	// }

	nc := symbolValWins.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(symbolValWins.Name, cd)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (symbolValWins *SymbolValWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*SymbolValWinsData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// OnStatsWithPB -
func (symbolValWins *SymbolValWins) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	return 0, nil
	// pbcd, isok := pbComponentData.(*sgc7pb.SymbolValWinsData)
	// if !isok {
	// 	goutils.Error("SymbolValWins.OnStatsWithPB",
	// 		zap.Error(ErrIvalidProto))

	// 	return 0, ErrIvalidProto
	// }

	// return symbolValWins.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// OnStats
func (symbolValWins *SymbolValWins) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
	// wins := int64(0)
	// isTrigger := false

	// for _, v := range lst {
	// 	gp, isok := v.CurGameModParams.(*GameParams)
	// 	if isok {
	// 		curComponent, isok := gp.MapComponentMsgs[symbolValWins.Name]
	// 		if isok {
	// 			curwins, err := symbolValWins.OnStatsWithPB(feature, curComponent, v)
	// 			if err != nil {
	// 				goutils.Error("SymbolValWins.OnStats",
	// 					zap.Error(err))

	// 				continue
	// 			}

	// 			isTrigger = true
	// 			wins += curwins
	// 		}
	// 	}
	// }

	// feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

	// if feature.Parent != nil {
	// 	totalwins := int64(0)

	// 	for _, v := range lst {
	// 		totalwins += v.CashWin
	// 	}

	// 	feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
	// }

	// return isTrigger, stake.CashBet, wins
}

// NewComponentData -
func (symbolValWins *SymbolValWins) NewComponentData() IComponentData {
	return &SymbolValWinsData{}
}

func (symbolValWins *SymbolValWins) GetWinMulti(basicCD *BasicComponentData) int {
	winMulti, isok := basicCD.GetConfigIntVal(STCVWinMulti)
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
