package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const SymbolVal2TypeName = "symbolVal2"

// SymbolVal2Config - configuration for SymbolVal2 feature
type SymbolVal2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string                   `yaml:"symbol" json:"symbol"`
	WeightSet            string                   `yaml:"weightSet" json:"weightSet"`
	WeightsVal           []string                 `yaml:"weightsVal" json:"weightsVal"`
	DefaultVal           int                      `yaml:"defaultVal" json:"defaultVal"`
	RNGSet               string                   `yaml:"RNGSet" json:"RNGSet"`
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature" json:"otherSceneFeature"`
}

type SymbolVal2 struct {
	*BasicComponent   `json:"-"`
	Config            *SymbolVal2Config       `json:"config"`
	SymbolCode        int                     `json:"-"`
	WeightsVal        []*sgc7game.ValWeights2 `json:"-"`
	WeightSet         *sgc7game.ValWeights2   `json:"-"`
	OtherSceneFeature *OtherSceneFeature      `json:"-"`
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

	return symbolVal2.InitEx(cfg, pool)
}

// InitEx -
func (symbolVal2 *SymbolVal2) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolVal2.Config = cfg.(*SymbolVal2Config)
	symbolVal2.Config.ComponentType = SymbolVal2TypeName

	if symbolVal2.Config.WeightSet != "" {
		vw2, err := pool.LoadIntWeights(symbolVal2.Config.WeightSet, symbolVal2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("SymbolVal2.Init:LoadValWeights",
				zap.String("Weight", symbolVal2.Config.WeightSet),
				zap.Error(err))

			return err
		}

		symbolVal2.WeightSet = vw2
	}

	for _, v := range symbolVal2.Config.WeightsVal {
		vw2, err := pool.LoadIntWeights(v, symbolVal2.Config.UseFileMapping)
		if err != nil {
			goutils.Error("SymbolVal2.Init:LoadValWeights",
				zap.String("Weight", v),
				zap.Error(err))

			return err
		}

		symbolVal2.WeightsVal = append(symbolVal2.WeightsVal, vw2)
	}

	symbolVal2.SymbolCode = pool.DefaultPaytables.MapSymbols[symbolVal2.Config.Symbol]

	if symbolVal2.Config.OtherSceneFeature != nil {
		symbolVal2.OtherSceneFeature = NewOtherSceneFeature(symbolVal2.Config.OtherSceneFeature)
	}

	symbolVal2.onInit(&symbolVal2.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolVal2 *SymbolVal2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := symbolVal2.GetTargetScene3(gameProp, curpr, prs, 0)

	if gs.HasSymbol(symbolVal2.SymbolCode) {
		os := gameProp.PoolScene.New(gs.Width, gs.Height)
		// os, err := sgc7game.NewGameScene(gs.Width, gs.Height)
		// if err != nil {
		// 	goutils.Error("SymbolVal2.OnPlayGame:NewGameScene",
		// 		zap.Error(err))

		// 	return err
		// }

		setIndex := -1
		if symbolVal2.Config.RNGSet != "" {
			rng := gameProp.GetTagInt(symbolVal2.Config.RNGSet)
			setIndex = rng
		} else {
			rv, err := symbolVal2.WeightSet.RandVal(plugin)
			if err != nil {
				goutils.Error("SymbolVal2.OnPlayGame:RandVal",
					zap.Error(err))

				return "", err
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

						return "", err
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
	} else {
		symbolVal2.ClearOtherScene(gameProp)
	}

	nc := symbolVal2.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(symbolVal2.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (symbolVal2 *SymbolVal2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after SymbolVal2", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// // OnStats
// func (symbolVal2 *SymbolVal2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewSymbolVal2(name string) IComponent {
	return &SymbolVal2{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
