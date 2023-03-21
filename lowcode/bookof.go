package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// BookOfConfig - configuration for BookOf feature
type BookOfConfig struct {
	BasicComponentConfig `yaml:",inline"`
	BetType              string `yaml:"betType"` // bet or totalBet
	ForceTrigger         bool   `yaml:"forceTrigger"`
	WeightTrigger        string `yaml:"weightTrigger"`
	WeightSymbolNum      string `yaml:"weightSymbolNum"`
	WeightSymbol         string `yaml:"weightSymbol"`
	ForceSymbolNum       int    `yaml:"forceSymbolNum"`
	SymbolRNG            string `yaml:"symbolRNG"` // 只在ForceSymbolNum为1时有效
}

type BookOf struct {
	*BasicComponent
	Config          *BookOfConfig
	Symbols         []int
	WeightTrigger   *sgc7game.ValWeights2
	WeightSymbolNum *sgc7game.ValWeights2
	WeightSymbol    *sgc7game.ValWeights2
}

// Init -
func (bookof *BookOf) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BookOf.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BookOfConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BookOf.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	bookof.Config = cfg

	if bookof.Config.WeightTrigger != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(bookof.Config.WeightTrigger, "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("BookOf.Init:LoadValWeights2FromExcel",
				zap.String("Weight", bookof.Config.WeightTrigger),
				zap.Error(err))

			return err
		}

		bookof.WeightTrigger = vw2
	}

	if bookof.Config.WeightSymbolNum != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(bookof.Config.WeightSymbolNum, "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("BookOf.Init:LoadValWeights2FromExcel",
				zap.String("Weight", bookof.Config.WeightSymbolNum),
				zap.Error(err))

			return err
		}

		bookof.WeightSymbolNum = vw2
	}

	if bookof.Config.WeightSymbol != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(bookof.Config.WeightSymbol, "val", "weight", gameProp.CurPaytables)
		if err != nil {
			goutils.Error("BookOf.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("Weight", bookof.Config.WeightSymbol),
				zap.Error(err))

			return err
		}

		bookof.WeightSymbol = vw2
	}

	bookof.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (bookof *BookOf) OnNewGame(gameProp *GameProperty) error {
	return nil
}

// OnNewStep -
func (bookof *BookOf) OnNewStep(gameProp *GameProperty) error {
	bookof.BasicComponent.OnNewStep()

	bookof.Symbols = nil

	return nil
}

// playgame
func (bookof *BookOf) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	isTrigger := bookof.Config.ForceTrigger

	if !isTrigger && bookof.WeightTrigger != nil {
		iv, err := bookof.WeightTrigger.RandVal(plugin)
		if err != nil {
			goutils.Error("bookof.OnPlayGame:WeightTrigger.RandVal",
				zap.Error(err))

			return err
		}

		if iv.Int() != 0 {
			isTrigger = true
		}
	}

	if isTrigger {
		symbolNum := bookof.Config.ForceSymbolNum

		if symbolNum <= 0 && bookof.WeightSymbolNum != nil {
			iv, err := bookof.WeightSymbolNum.RandVal(plugin)
			if err != nil {
				goutils.Error("bookof.OnPlayGame:WeightSymbolNum.RandVal",
					zap.Error(err))

				return err
			}

			symbolNum = iv.Int()
		}

		gs := bookof.GetTargetScene(gameProp, curpr)

		if bookof.Config.ForceSymbolNum == 1 && bookof.Config.SymbolRNG != "" {
			rng := gameProp.MapInt[bookof.Config.SymbolRNG]
			cs := bookof.WeightSymbol.Vals[rng]

			bookof.Symbols = append(bookof.Symbols, cs.Int())

			ngs, err := bookof.procBookOfScene(gs, cs.Int())
			if err != nil {
				goutils.Error("bookof.OnPlayGame:procBookOfScene",
					zap.Error(err))

				return err
			}

			bookof.AddScene(gameProp, curpr, ngs)

			scr := sgc7game.CalcScatter3(gs, gameProp.CurPaytables, cs.Int(), GetBet(stake, bookof.Config.BetType), 1, func(scatter int, cursymbol int) bool {
				return cursymbol == cs.Int()
			}, true)
			if scr != nil {
				bookof.AddResult(curpr, scr)
			}
		} else {
			curWeight := bookof.WeightSymbol.Clone()

			for i := 0; i < symbolNum; i++ {
				cs, err := curWeight.RandVal(plugin)
				if err != nil {
					goutils.Error("bookof.OnPlayGame:curWeight.RandVal",
						zap.Error(err))

					return err
				}

				bookof.Symbols = append(bookof.Symbols, cs.Int())

				err = curWeight.RemoveVal(cs)
				if err != nil {
					goutils.Error("bookof.OnPlayGame:curWeight.RemoveVal",
						zap.Error(err))

					return err
				}

				ngs, err := bookof.procBookOfScene(gs, cs.Int())
				if err != nil {
					goutils.Error("bookof.OnPlayGame:procBookOfScene",
						zap.Error(err))

					return err
				}

				bookof.AddScene(gameProp, curpr, ngs)

				scr := sgc7game.CalcScatter3(gs, gameProp.CurPaytables, cs.Int(), GetBet(stake, bookof.Config.BetType), 1, func(scatter int, cursymbol int) bool {
					return cursymbol == cs.Int()
				}, true)
				if scr != nil {
					bookof.AddResult(curpr, scr)
				}
			}
		}

		bookof.onStepEnd(gameProp, curpr, gp)

		bookof.BuildPBComponent(gp)
	} else {
		bookof.onStepEnd(gameProp, curpr, gp)
	}

	return nil
}

// procBookOfScene - outpur to asciigame
func (bookof *BookOf) procBookOfScene(gs *sgc7game.GameScene, symbol int) (*sgc7game.GameScene, error) {
	ngs := gs.Clone()

	for x, arr := range gs.Arr {
		hass := false
		for _, s := range arr {
			if s == symbol {
				hass = true

				break
			}
		}

		if hass {
			for y := range arr {
				ngs.Arr[x][y] = symbol
			}
		}
	}

	return ngs, nil
}

// OnAsciiGame - outpur to asciigame
func (bookof *BookOf) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(bookof.Symbols) > 0 {
		strsymbols := ""
		for i, v := range bookof.Symbols {
			if i > 0 {
				strsymbols += ", "
			}

			strsymbols += gameProp.CurPaytables.GetStringFromInt(v)
		}

		fmt.Printf("The BookOf Symbols is %v\n", strsymbols)

		for i, si := range bookof.UsedScenes {
			asciigame.OutputScene(fmt.Sprintf("The symbols for BookOf - %v ", i+1), pr.Scenes[si], mapSymbolColor)
		}

		asciigame.OutputResults(fmt.Sprintf("%v wins", bookof.Name), pr, func(i int, ret *sgc7game.Result) bool {
			return goutils.IndexOfIntSlice(bookof.UsedResults, i, 0) >= 0
		}, mapSymbolColor)
	}

	return nil
}

// OnStats
func (bookof *BookOf) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewBookOf(name string) IComponent {
	return &BookOf{
		BasicComponent: NewBasicComponent(name),
	}
}
