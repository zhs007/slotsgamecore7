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

// SymbolVal2Config - configuration for SymbolVal2 feature
type SymbolVal2Config struct {
	BasicComponentConfig `yaml:",inline"`
	Symbol               string                   `yaml:"symbol"`
	WeightSet            string                   `yaml:"weightSet"`
	WeightsVal           []string                 `yaml:"weightsVal"`
	DefaultVal           int                      `yaml:"defaultVal"`
	RNGSet               string                   `yaml:"RNGSet"`
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
}

type SymbolVal2 struct {
	*BasicComponent
	Config            *SymbolVal2Config
	SymbolCode        int
	WeightsVal        []*sgc7game.ValWeights2
	WeightSet         *sgc7game.ValWeights2
	OtherSceneFeature *OtherSceneFeature
}

// Init -
func (symbolVal2 *SymbolVal2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolVal2.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolVal2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolVal2.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	symbolVal2.Config = cfg

	if cfg.WeightSet != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(cfg.WeightSet), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolVal2.Init:LoadValWeights2FromExcel",
				zap.String("Weight", cfg.WeightSet),
				zap.Error(err))

			return err
		}

		symbolVal2.WeightSet = vw2
	}

	for _, v := range cfg.WeightsVal {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(v), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolVal2.Init:LoadValWeights2FromExcel",
				zap.String("Weight", v),
				zap.Error(err))

			return err
		}

		symbolVal2.WeightsVal = append(symbolVal2.WeightsVal, vw2)
	}

	symbolVal2.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	if cfg.OtherSceneFeature != nil {
		symbolVal2.OtherSceneFeature = NewOtherSceneFeature(cfg.OtherSceneFeature)
	}

	symbolVal2.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (symbolVal2 *SymbolVal2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[symbolVal2.Name].(*BasicComponentData)

	gs := symbolVal2.GetTargetScene(gameProp, curpr, cd)

	if gs.HasSymbol(symbolVal2.SymbolCode) {
		os, err := sgc7game.NewGameScene(gs.Width, gs.Height)
		if err != nil {
			goutils.Error("SymbolVal2.OnPlayGame:NewGameScene",
				zap.Error(err))

			return err
		}

		setIndex := -1
		if symbolVal2.Config.RNGSet != "" {
			rng := gameProp.GetTagInt(symbolVal2.Config.RNGSet)
			setIndex = rng
		} else {
			rv, err := symbolVal2.WeightSet.RandVal(plugin)
			if err != nil {
				goutils.Error("SymbolVal2.OnPlayGame:RandVal",
					zap.Error(err))

				return err
			}

			setIndex = rv.Int()
		}

		vw2 := symbolVal2.WeightsVal[setIndex]

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if s == symbolVal2.SymbolCode {
					cv, err := vw2.RandVal(plugin)
					if err != nil {
						goutils.Error("SymbolVal2.OnPlayGame:WeightVal.RandVal",
							zap.Error(err))

						return err
					}

					os.Arr[x][y] = cv.Int()
				} else {
					os.Arr[x][y] = symbolVal2.Config.DefaultVal
				}
			}
		}

		symbolVal2.AddOtherScene(gameProp, curpr, os, cd)

		if symbolVal2.OtherSceneFeature != nil {
			gameProp.procOtherSceneFeature(symbolVal2.OtherSceneFeature, curpr, os)
		}
	}

	symbolVal2.onStepEnd(gameProp, curpr, gp, "")

	gp.AddComponentData(symbolVal2.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolVal2 *SymbolVal2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[symbolVal2.Name].(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("The value of the symbols", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// OnStats
func (symbolVal2 *SymbolVal2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewSymbolVal2(name string) IComponent {
	return &SymbolVal2{
		BasicComponent: NewBasicComponent(name),
	}
}
