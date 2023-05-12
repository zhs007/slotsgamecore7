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
	Symbol               string                   `yaml:"symbol"`
	Symbols              []string                 `yaml:"symbols"`
	WeightMulti          string                   `yaml:"weightMulti"`
	MapWeightMulti       map[string]string        `yaml:"mapWeightMulti"` // 可以配置多套权重
	ValUsed              string                   `yaml:"valUsed"`        // 用这个值来确定使用的权重
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature"`
}

type SymbolMulti struct {
	*BasicComponent
	Config            *SymbolMultiConfig
	SymbolCodes       []int
	WeightMulti       *sgc7game.ValWeights2
	MapWeightMulti    map[string]*sgc7game.ValWeights2
	OtherSceneFeature *OtherSceneFeature
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

	if len(cfg.MapWeightMulti) > 0 {
		symbolMulti.MapWeightMulti = make(map[string]*sgc7game.ValWeights2, 0)

		for k, v := range cfg.MapWeightMulti {
			vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(v, symbolMulti.Config.UseFileMapping), "val", "weight", sgc7game.NewIntVal[int])
			if err != nil {
				goutils.Error("SymbolMulti.Init:LoadValWeights2FromExcel",
					zap.String("Weight", v),
					zap.Error(err))

				return err
			}

			symbolMulti.MapWeightMulti[k] = vw2
		}
	} else if symbolMulti.Config.WeightMulti != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(symbolMulti.Config.WeightMulti, symbolMulti.Config.UseFileMapping), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolMulti.Init:LoadValWeights2FromExcel",
				zap.String("Weight", symbolMulti.Config.WeightMulti),
				zap.Error(err))

			return err
		}

		symbolMulti.WeightMulti = vw2
	}

	if len(cfg.Symbols) > 0 {
		for _, v := range cfg.Symbols {
			symbolMulti.SymbolCodes = append(symbolMulti.SymbolCodes, pool.DefaultPaytables.MapSymbols[v])
		}
	} else {
		symbolMulti.SymbolCodes = append(symbolMulti.SymbolCodes, pool.DefaultPaytables.MapSymbols[cfg.Symbol])
	}

	if cfg.OtherSceneFeature != nil {
		symbolMulti.OtherSceneFeature = NewOtherSceneFeature(cfg.OtherSceneFeature)
	}

	symbolMulti.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (symbolMulti *SymbolMulti) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[symbolMulti.Name].(*BasicComponentData)

	gs := symbolMulti.GetTargetScene(gameProp, curpr, cd, "")

	if gs.HasSymbols(symbolMulti.SymbolCodes) {
		// os := gameProp.Pool.PoolGameScene.New(gs.Width, gs.Height, false)
		os, err := sgc7game.NewGameScene(gs.Width, gs.Height)
		if err != nil {
			goutils.Error("SymbolMulti.OnPlayGame:NewGameScene",
				zap.Error(err))

			return err
		}

		vw2 := symbolMulti.WeightMulti
		if len(symbolMulti.MapWeightMulti) > 0 && symbolMulti.Config.ValUsed != "" {
			val := gameProp.GetTagGlobalStr(symbolMulti.Config.ValUsed)
			vw2 = symbolMulti.MapWeightMulti[val]
		}

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(symbolMulti.SymbolCodes, s, 0) >= 0 {
					cv, err := vw2.RandVal(plugin)
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

		if symbolMulti.OtherSceneFeature != nil {
			gameProp.procOtherSceneFeature(symbolMulti.OtherSceneFeature, curpr, os)
		}
	}

	symbolMulti.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(symbolMulti.Name, cd)
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
