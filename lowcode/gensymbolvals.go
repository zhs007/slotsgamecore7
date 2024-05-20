package lowcode

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const GenSymbolValsTypeName = "genSymbolVals"

// GenSymbolValsConfig - configuration for GenSymbolVals
type GenSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	DefaultVal           int `yaml:"defaultVal" json:"defaultVal"`
}

// SetLinkComponent
func (cfg *GenSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolValsConfig `json:"config"`
}

// Init -
func (genSymbolVals *GenSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolVals *GenSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolVals.Config = cfg.(*GenSymbolValsConfig)
	genSymbolVals.Config.ComponentType = GenSymbolValsTypeName

	genSymbolVals.onInit(&genSymbolVals.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genSymbolVals *GenSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	// gs := genSymbolVals.GetTargetScene3(gameProp, curpr, prs, 0)
	// if gs != nil {
	os := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolVals.Config.DefaultVal)

	// if genSymbolValsWithSymbol.Config.IsUseSource {
	// 	os = genSymbolValsWithSymbol.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	// }

	// nos := os

	// if genSymbolValsWithSymbol.Config.Type == GSVWSTypeNormal {
	// 	for x, arr := range gs.Arr {
	// 		for y, s := range arr {
	// 			if goutils.IndexOfIntSlice(genSymbolValsWithSymbol.Config.SymbolCodes, s, 0) >= 0 {
	// 				if nos == nil {
	// 					curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
	// 					if err != nil {
	// 						goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
	// 							goutils.Err(err))

	// 						return "", err
	// 					}

	// 					nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)

	// 					nos.Arr[x][y] = curv.Int()
	// 				} else if nos.Arr[x][y] == genSymbolValsWithSymbol.Config.DefaultVal {
	// 					if nos == os {
	// 						nos = os.CloneEx(gameProp.PoolScene)
	// 					}

	// 					curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
	// 					if err != nil {
	// 						goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
	// 							goutils.Err(err))

	// 						return "", err
	// 					}

	// 					nos.Arr[x][y] = curv.Int()
	// 				}
	// 			} else {
	// 				if nos != nil && nos.Arr[x][y] != genSymbolValsWithSymbol.Config.DefaultVal {
	// 					if nos == os {
	// 						nos = os.CloneEx(gameProp.PoolScene)
	// 					}

	// 					nos.Arr[x][y] = genSymbolValsWithSymbol.Config.DefaultVal
	// 				}
	// 			}
	// 		}
	// 	}
	// } else if genSymbolValsWithSymbol.Config.Type == GSVWSTypeNonClear {
	// 	for x, arr := range gs.Arr {
	// 		for y, s := range arr {
	// 			if goutils.IndexOfIntSlice(genSymbolValsWithSymbol.Config.SymbolCodes, s, 0) >= 0 {
	// 				if nos == nil {
	// 					curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
	// 					if err != nil {
	// 						goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
	// 							goutils.Err(err))

	// 						return "", err
	// 					}

	// 					nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)

	// 					nos.Arr[x][y] = curv.Int()
	// 				} else if nos.Arr[x][y] == genSymbolValsWithSymbol.Config.DefaultVal {
	// 					if nos == os {
	// 						nos = os.CloneEx(gameProp.PoolScene)
	// 					}

	// 					curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
	// 					if err != nil {
	// 						goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
	// 							goutils.Err(err))

	// 						return "", err
	// 					}

	// 					nos.Arr[x][y] = curv.Int()
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	// if nos == nil && genSymbolValsWithSymbol.Config.IsAlwaysGen {
	// 	nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)
	// }

	// if nos == os {
	// 	nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, ErrComponentDoNothing
	// }

	genSymbolVals.AddOtherScene(gameProp, curpr, os, cd)

	nc := genSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
	// }

	// nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

	// return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genSymbolVals *GenSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("GenSymbolVals", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// // OnStats
// func (genSymbolValsWithPos *GenSymbolValsWithSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewGenSymbolVals(name string) IComponent {
	return &GenSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonGenSymbolVals struct {
	DefaultVal int `json:"defaultVal"`
}

func (jcfg *jsonGenSymbolVals) build() *GenSymbolValsConfig {
	cfg := &GenSymbolValsConfig{
		DefaultVal: jcfg.DefaultVal,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseGenSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
