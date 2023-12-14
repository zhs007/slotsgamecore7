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

const SymbolMultiTypeName = "symbolMulti"

// SymbolMultiConfig - configuration for SymbolMulti feature
type SymbolMultiConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string                   `yaml:"symbol" json:"-"`                      // 弃用，用symbols
	Symbols              []string                 `yaml:"symbols" json:"symbols"`               // 这些符号可以有倍数
	WeightMulti          string                   `yaml:"weightMulti" json:"weightMulti"`       // 倍数权重
	StaticMulti          int                      `yaml:"staticMulti" json:"staticMulti"`       // 恒定倍数
	MapWeightMulti       map[string]string        `yaml:"mapWeightMulti" json:"mapWeightMulti"` // 可以配置多套权重
	ValUsed              string                   `yaml:"valUsed" json:"valUsed"`               // 用这个值来确定使用的权重
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature" json:"otherSceneFeature"`
}

type SymbolMulti struct {
	*BasicComponent   `json:"-"`
	Config            *SymbolMultiConfig               `json:"config"`
	SymbolCodes       []int                            `json:"-"`
	WeightMulti       *sgc7game.ValWeights2            `json:"-"`
	MapWeightMulti    map[string]*sgc7game.ValWeights2 `json:"-"`
	OtherSceneFeature *OtherSceneFeature               `json:"-"`
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

	return symbolMulti.InitEx(cfg, pool)
}

// InitEx -
func (symbolMulti *SymbolMulti) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolMulti.Config = cfg.(*SymbolMultiConfig)
	symbolMulti.Config.ComponentType = SymbolMultiTypeName

	if len(symbolMulti.Config.MapWeightMulti) > 0 {
		symbolMulti.MapWeightMulti = make(map[string]*sgc7game.ValWeights2, 0)

		for k, v := range symbolMulti.Config.MapWeightMulti {
			vw2, err := pool.LoadIntWeights(v, symbolMulti.Config.UseFileMapping)
			if err != nil {
				goutils.Error("SymbolMulti.Init:LoadValWeights",
					zap.String("Weight", v),
					zap.Error(err))

				return err
			}

			symbolMulti.MapWeightMulti[k] = vw2
		}
	} else if symbolMulti.Config.WeightMulti != "" {
		vw2, err := pool.LoadIntWeights(symbolMulti.Config.WeightMulti, symbolMulti.Config.UseFileMapping)
		if err != nil {
			goutils.Error("SymbolMulti.Init:LoadValWeights",
				zap.String("Weight", symbolMulti.Config.WeightMulti),
				zap.Error(err))

			return err
		}

		symbolMulti.WeightMulti = vw2
	}

	if len(symbolMulti.Config.Symbols) > 0 {
		for _, v := range symbolMulti.Config.Symbols {
			symbolMulti.SymbolCodes = append(symbolMulti.SymbolCodes, pool.DefaultPaytables.MapSymbols[v])
		}
	} else {
		symbolMulti.SymbolCodes = append(symbolMulti.SymbolCodes, pool.DefaultPaytables.MapSymbols[symbolMulti.Config.Symbol])
	}

	if symbolMulti.Config.OtherSceneFeature != nil {
		symbolMulti.OtherSceneFeature = NewOtherSceneFeature(symbolMulti.Config.OtherSceneFeature)
	}

	symbolMulti.onInit(&symbolMulti.Config.BasicComponentConfig)

	return nil
}

// playgame
func (symbolMulti *SymbolMulti) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolMulti.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[symbolMulti.Name].(*BasicComponentData)

	gs := symbolMulti.GetTargetScene2(gameProp, curpr, cd, symbolMulti.Name, "")

	if gs.HasSymbols(symbolMulti.SymbolCodes) {
		os := gameProp.PoolScene.New(gs.Width, gs.Height, false)

		if symbolMulti.WeightMulti == nil {
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(symbolMulti.SymbolCodes, s, 0) >= 0 {
						os.Arr[x][y] = symbolMulti.Config.StaticMulti
					} else {
						os.Arr[x][y] = 1
					}
				}
			}
		} else {
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
		}

		symbolMulti.AddOtherScene(gameProp, curpr, os, cd)

		if symbolMulti.OtherSceneFeature != nil {
			gameProp.procOtherSceneFeature(symbolMulti.OtherSceneFeature, curpr, os)
		}
	} else {
		symbolMulti.ClearOtherScene(gameProp)
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
