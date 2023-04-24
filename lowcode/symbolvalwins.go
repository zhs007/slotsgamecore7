package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// SymbolValWinsConfig - configuration for SymbolValWins
type SymbolValWinsConfig struct {
	BasicComponentConfig    `yaml:",inline"`
	BetType                 string `yaml:"betType"`                 // bet or totalBet
	TriggerSymbol           string `yaml:"triggerSymbol"`           // like collect
	Type                    string `yaml:"type"`                    // like scatters
	MinNum                  int    `yaml:"minNum"`                  // like 3
	IsTriggerSymbolNumMulti bool   `yaml:"isTriggerSymbolNumMulti"` // totalwins = totalvals * triggetSymbol's num
}

type SymbolValWins struct {
	*BasicComponent
	Config            *SymbolValWinsConfig
	TriggerSymbolCode int
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

	symbolValWins.Config = cfg

	if cfg.TriggerSymbol != "" {
		symbolValWins.TriggerSymbolCode = pool.DefaultPaytables.MapSymbols[cfg.TriggerSymbol]
	} else {
		symbolValWins.TriggerSymbolCode = -1
	}

	symbolValWins.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (symbolValWins *SymbolValWins) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[symbolValWins.Name].(*BasicComponentData)

	gs := symbolValWins.GetTargetScene(gameProp, curpr, cd)
	isTrigger := true
	symbolnum := 0

	if symbolValWins.TriggerSymbolCode >= 0 {
		isTrigger = false

		if symbolValWins.Config.Type == WinTypeCountScatter {
			ret := sgc7game.CalcScatterEx(gs, symbolValWins.TriggerSymbolCode, symbolValWins.Config.MinNum, func(scatter int, cursymbol int) bool {
				return cursymbol == scatter
			})

			if ret != nil {
				isTrigger = true

				symbolnum = ret.SymbolNums
			}
		}
	}

	if isTrigger {
		os := symbolValWins.GetTargetOtherScene(gameProp, curpr, cd)

		if os != nil {
			totalvals := 0
			pos := make([]int, 0, len(os.Arr)*len(os.Arr[0])*2)

			for x := 0; x < len(os.Arr); x++ {
				for y := 0; y < len(os.Arr[x]); y++ {
					if os.Arr[x][y] > 0 {
						totalvals += os.Arr[x][y]
						pos = append(pos, x, y)
					}
				}
			}

			if totalvals > 0 {
				ret := &sgc7game.Result{
					Symbol:     gs.Arr[pos[0]][pos[1]],
					Type:       sgc7game.RTSymbolVal,
					LineIndex:  -1,
					Pos:        pos,
					SymbolNums: len(pos) / 2,
				}

				bet := GetBet(stake, symbolValWins.Config.BetType)

				if symbolValWins.Config.IsTriggerSymbolNumMulti {
					ret.CoinWin = totalvals * symbolnum
					ret.CashWin = ret.CoinWin * bet
				} else {
					ret.CoinWin = totalvals
					ret.CashWin = ret.CoinWin * bet
				}

				symbolValWins.AddResult(curpr, ret, cd)
			}
		}
	}

	symbolValWins.onStepEnd(gameProp, curpr, gp, "")

	gp.AddComponentData(symbolValWins.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolValWins *SymbolValWins) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[symbolValWins.Name].(*BasicComponentData)

	asciigame.OutputResults("wins", pr, func(i int, ret *sgc7game.Result) bool {
		return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
	}, mapSymbolColor)

	return nil
}

// OnStats
func (symbolValWins *SymbolValWins) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	wins := int64(0)
	isTrigger := false

	for _, v := range lst {
		gp, isok := v.CurGameModParams.(*GameParams)
		if isok {
			curComponent, isok := gp.MapComponents[symbolValWins.Name]
			if isok {
				curwins, err := symbolValWins.OnStatsWithPB(feature, curComponent, v)
				if err != nil {
					goutils.Error("SymbolValWins.OnStats",
						zap.Error(err))

					continue
				}

				isTrigger = true
				wins += curwins
			}
		}
	}

	feature.CurWins.AddWin(int(wins) * 100 / int(stake.CashBet))

	if feature.Parent != nil {
		totalwins := int64(0)

		for _, v := range lst {
			totalwins += v.CashWin
		}

		feature.AllWins.AddWin(int(totalwins) * 100 / int(stake.CashBet))
	}

	return isTrigger, stake.CashBet, wins
}

func NewSymbolValWins(name string) IComponent {
	return &SymbolValWins{
		BasicComponent: NewBasicComponent(name),
	}
}
