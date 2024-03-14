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

const SymbolValTypeName = "symbolVal"

// const (
// 	SVCVWeightVal string = "weightVal" // 可以修改配置项里的 weightVal
// )

// SymbolValConfig - configuration for SymbolMulti feature
type SymbolValConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string                   `yaml:"symbol" json:"symbol"`
	WeightVal            string                   `yaml:"weightVal" json:"weightVal"`
	DefaultVal           int                      `yaml:"defaultVal" json:"defaultVal"`
	OtherSceneFeature    *OtherSceneFeatureConfig `yaml:"otherSceneFeature" json:"otherSceneFeature"`
	EmptyOtherSceneVal   int                      `yaml:"emptyOtherSceneVal" json:"emptyOtherSceneVal"` // 如果配置了otherscene，那么当otherscene里的某个位置为这个值时，才新赋值
}

type SymbolVal struct {
	*BasicComponent   `json:"-"`
	Config            *SymbolValConfig      `json:"config"`
	SymbolCode        int                   `json:"-"`
	WeightVal         *sgc7game.ValWeights2 `json:"-"`
	OtherSceneFeature *OtherSceneFeature    `json:"-"`
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

	return symbolVal.InitEx(cfg, pool)
}

// InitEx -
func (symbolVal *SymbolVal) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolVal.Config = cfg.(*SymbolValConfig)
	symbolVal.Config.ComponentType = SymbolValTypeName

	if symbolVal.Config.WeightVal != "" {
		vw2, err := pool.LoadIntWeights(symbolVal.Config.WeightVal, symbolVal.Config.UseFileMapping)
		if err != nil {
			goutils.Error("SymbolVal.Init:LoadValWeights",
				zap.String("Weight", symbolVal.Config.WeightVal),
				zap.Error(err))

			return err
		}

		symbolVal.WeightVal = vw2
	}

	symbolVal.SymbolCode = pool.DefaultPaytables.MapSymbols[symbolVal.Config.Symbol]

	if symbolVal.Config.OtherSceneFeature != nil {
		symbolVal.OtherSceneFeature = NewOtherSceneFeature(symbolVal.Config.OtherSceneFeature)
	}

	symbolVal.onInit(&symbolVal.Config.BasicComponentConfig)

	return nil
}

func (symbolVal *SymbolVal) GetWeightVal(gameProp *GameProperty, basicCD *BasicComponentData) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVWeightVal)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, symbolVal.Config.UseFileMapping)

		return vw2
	}

	return symbolVal.WeightVal
}

// playgame
func (symbolVal *SymbolVal) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := symbolVal.GetTargetScene3(gameProp, curpr, prs, cd, symbolVal.Name, "", 0)

	if gs.HasSymbol(symbolVal.SymbolCode) {
		vw := symbolVal.GetWeightVal(gameProp, cd)

		os1 := symbolVal.GetTargetOtherScene3(gameProp, curpr, prs, 0)
		if os1 == nil {
			os := gameProp.PoolScene.New(gs.Width, gs.Height, false)

			for x, arr := range gs.Arr {
				for y, s := range arr {
					if s == symbolVal.SymbolCode {
						cv, err := vw.RandVal(plugin)
						if err != nil {
							goutils.Error("SymbolVal.OnPlayGame:WeightVal.RandVal",
								zap.Error(err))

							return "", err
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
		} else {
			os := os1.CloneEx(gameProp.PoolScene)

			for x, arr := range gs.Arr {
				for y, s := range arr {
					if os.Arr[x][y] != symbolVal.Config.EmptyOtherSceneVal {
						continue
					}

					if s == symbolVal.SymbolCode {
						cv, err := vw.RandVal(plugin)
						if err != nil {
							goutils.Error("SymbolVal.OnPlayGame:WeightVal.RandVal",
								zap.Error(err))

							return "", err
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
	} else {
		symbolVal.ClearOtherScene(gameProp)
	}

	nc := symbolVal.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(symbolVal.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (symbolVal *SymbolVal) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

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
		BasicComponent: NewBasicComponent(name, 1),
	}
}
