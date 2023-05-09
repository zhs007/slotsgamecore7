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

// SymbolValConfig - configuration for SymbolMulti feature
type SymbolValConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Symbol               string                   `yaml:"symbol"`
	WeightVal            string                   `yaml:"weightVal"`
	DefaultVal           int                      `yaml:"defaultVal"`
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
}

type SymbolVal struct {
	*BasicComponent
	Config            *SymbolValConfig
	SymbolCode        int
	WeightVal         *sgc7game.ValWeights2
	OtherSceneFeature *OtherSceneFeature
}

// Init -
func (symbolVal *SymbolVal) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolVal.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolValConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolVal.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	symbolVal.Config = cfg

	if symbolVal.Config.WeightVal != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(symbolVal.Config.WeightVal), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolVal.Init:LoadValWeights2FromExcel",
				zap.String("Weight", symbolVal.Config.WeightVal),
				zap.Error(err))

			return err
		}

		symbolVal.WeightVal = vw2
	}

	symbolVal.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	if cfg.OtherSceneFeature != nil {
		symbolVal.OtherSceneFeature = NewOtherSceneFeature(cfg.OtherSceneFeature)
	}

	symbolVal.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (symbolVal *SymbolVal) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolVal.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[symbolVal.Name].(*BasicComponentData)

	gs := symbolVal.GetTargetScene(gameProp, curpr, cd, "")

	if gs.HasSymbol(symbolVal.SymbolCode) {
		os, err := sgc7game.NewGameScene(gs.Width, gs.Height)
		if err != nil {
			goutils.Error("SymbolVal.OnPlayGame:NewGameScene",
				zap.Error(err))

			return err
		}

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if s == symbolVal.SymbolCode {
					cv, err := symbolVal.WeightVal.RandVal(plugin)
					if err != nil {
						goutils.Error("SymbolVal.OnPlayGame:WeightVal.RandVal",
							zap.Error(err))

						return err
					}

					os.Arr[x][y] = cv.Int()
				} else {
					os.Arr[x][y] = symbolVal.Config.DefaultVal
				}
			}
		}

		symbolVal.AddOtherScene(gameProp, curpr, os, cd)

		if symbolVal.OtherSceneFeature != nil {
			gameProp.procOtherSceneFeature(symbolVal.OtherSceneFeature, curpr, os)
		}
	}

	symbolVal.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(symbolVal.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolVal *SymbolVal) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[symbolVal.Name].(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("The value of the symbols", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// OnStats
func (symbolVal *SymbolVal) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewSymbolVal(name string) IComponent {
	return &SymbolVal{
		BasicComponent: NewBasicComponent(name),
	}
}
