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

// SymbolMultiConfig - configuration for SymbolMulti feature
type SymbolMultiConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Symbol               string `yaml:"symbol"`
	WeightMulti          string `yaml:"weightMulti"`
}

type SymbolMulti struct {
	*BasicComponent
	Config      *SymbolMultiConfig
	SymbolCode  int
	WeightMulti *sgc7game.ValWeights2
}

// Init -
func (symbolMulti *SymbolMulti) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolMulti.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolMultiConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolMulti.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	symbolMulti.Config = cfg

	if symbolMulti.Config.WeightMulti != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(symbolMulti.Config.WeightMulti, "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolMulti.Init:LoadValWeights2FromExcel",
				zap.String("Weight", symbolMulti.Config.WeightMulti),
				zap.Error(err))

			return err
		}

		symbolMulti.WeightMulti = vw2
	}

	symbolMulti.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	symbolMulti.onInit(&cfg.BasicComponentConfig)

	return nil
}

// // OnNewGame -
// func (symbolMulti *SymbolMulti) OnNewGame(gameProp *GameProperty) error {
// 	return nil
// }

// // OnNewStep -
// func (symbolMulti *SymbolMulti) OnNewStep(gameProp *GameProperty) error {
// 	symbolMulti.BasicComponent.OnNewStep()

// 	return nil
// }

// playgame
func (symbolMulti *SymbolMulti) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[symbolMulti.Name].(*BasicComponentData)

	gs := symbolMulti.GetTargetScene(gameProp, curpr, cd)

	os, err := sgc7game.NewGameScene(gs.Width, gs.Height)
	if err != nil {
		goutils.Error("SymbolMulti.OnPlayGame:NewGameScene",
			zap.Error(err))

		return err
	}

	for x, arr := range gs.Arr {
		for y, s := range arr {
			if s == symbolMulti.SymbolCode {
				cv, err := symbolMulti.WeightMulti.RandVal(plugin)
				if err != nil {
					goutils.Error("SymbolMulti.OnPlayGame:WeightMulti.RandVal",
						zap.Error(err))

					return err
				}

				os.Arr[x][y] = cv.Int()
			} else {
				os.Arr[x][y] = 1
			}
		}
	}

	symbolMulti.AddOtherScene(gameProp, curpr, os, cd)

	symbolMulti.onStepEnd(gameProp, curpr, gp)

	gp.AddComponentData(symbolMulti.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolMulti *SymbolMulti) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[symbolMulti.Name].(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("The multi of the symbols", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// OnStats
func (symbolMulti *SymbolMulti) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewSymbolMulti(name string) IComponent {
	return &SymbolMulti{
		BasicComponent: NewBasicComponent(name),
	}
}
