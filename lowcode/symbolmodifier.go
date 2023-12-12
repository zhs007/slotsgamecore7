package lowcode

// import (
// 	"os"

// 	"github.com/zhs007/goutils"
// 	"github.com/zhs007/slotsgamecore7/asciigame"
// 	sgc7game "github.com/zhs007/slotsgamecore7/game"
// 	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
// 	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
// 	"go.uber.org/zap"
// 	"gopkg.in/yaml.v2"
// )

// const SymbolModifierTypeName = "symbolModifier"

// // SymbolModifierConfig - configuration for SymbolModifier feature
// type SymbolModifierConfig struct {
// 	BasicComponentConfig `yaml:",inline" json:",inline"`
// 	Symbols              []string `yaml:"symbols" json:"symbols"`
// 	SymbolCodes          []int    `yaml:"-" json:"-"`
// 	TargetSymbols        []string `yaml:"targetSymbols" json:"targetSymbols"`
// 	TargetSymbolCodes    []int    `yaml:"-" json:"-"`
// 	PosArea              []int    `yaml:"posArea" json:"posArea"` // [minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
// }

// type SymbolModifier struct {
// 	*BasicComponent `json:"-"`
// 	Config          *SymbolValConfig `json:"config"`
// }

// // Init -
// func (symbolVal *SymbolModifier) Init(fn string, pool *GamePropertyPool) error {
// 	data, err := os.ReadFile(fn)
// 	if err != nil {
// 		goutils.Error("SymbolModifier.Init:ReadFile",
// 			zap.String("fn", fn),
// 			zap.Error(err))

// 		return err
// 	}

// 	cfg := &SymbolValConfig{}

// 	err = yaml.Unmarshal(data, cfg)
// 	if err != nil {
// 		goutils.Error("SymbolModifier.Init:Unmarshal",
// 			zap.String("fn", fn),
// 			zap.Error(err))

// 		return err
// 	}

// 	return symbolVal.InitEx(cfg, pool)
// }

// // InitEx -
// func (symbolVal *SymbolModifier) InitEx(cfg any, pool *GamePropertyPool) error {
// 	symbolVal.Config = cfg.(*SymbolValConfig)
// 	symbolVal.Config.ComponentType = SymbolValTypeName

// 	if symbolVal.Config.WeightVal != "" {
// 		vw2, err := pool.LoadIntWeights(symbolVal.Config.WeightVal, symbolVal.Config.UseFileMapping)
// 		if err != nil {
// 			goutils.Error("SymbolModifier.Init:LoadValWeights",
// 				zap.String("Weight", symbolVal.Config.WeightVal),
// 				zap.Error(err))

// 			return err
// 		}

// 		symbolVal.WeightVal = vw2
// 	}

// 	symbolVal.SymbolCode = pool.DefaultPaytables.MapSymbols[symbolVal.Config.Symbol]

// 	if symbolVal.Config.OtherSceneFeature != nil {
// 		symbolVal.OtherSceneFeature = NewOtherSceneFeature(symbolVal.Config.OtherSceneFeature)
// 	}

// 	symbolVal.onInit(&symbolVal.Config.BasicComponentConfig)

// 	return nil
// }

// // playgame
// func (symbolVal *SymbolModifier) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
// 	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

// 	symbolVal.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

// 	cd := gameProp.MapComponentData[symbolVal.Name].(*BasicComponentData)

// 	gs := symbolVal.GetTargetScene2(gameProp, curpr, cd, symbolVal.Name, "")

// 	if gs.HasSymbol(symbolVal.SymbolCode) {
// 		os1 := symbolVal.GetTargetOtherScene2(gameProp, curpr, cd, symbolVal.Name, "")
// 		if os1 == nil {
// 			os := gameProp.PoolScene.New(gs.Width, gs.Height, false)

// 			for x, arr := range gs.Arr {
// 				for y, s := range arr {
// 					if s == symbolVal.SymbolCode {
// 						cv, err := symbolVal.WeightVal.RandVal(plugin)
// 						if err != nil {
// 							goutils.Error("SymbolModifier.OnPlayGame:WeightVal.RandVal",
// 								zap.Error(err))

// 							return err
// 						}

// 						os.Arr[x][y] = cv.Int()
// 					} else {
// 						os.Arr[x][y] = symbolVal.Config.DefaultVal
// 					}
// 				}
// 			}

// 			symbolVal.AddOtherScene(gameProp, curpr, os, cd)

// 			if symbolVal.OtherSceneFeature != nil {
// 				gameProp.procOtherSceneFeature(symbolVal.OtherSceneFeature, curpr, os)
// 			}
// 		} else {
// 			os := os1.CloneEx(gameProp.PoolScene)

// 			for x, arr := range gs.Arr {
// 				for y, s := range arr {
// 					if s != symbolVal.Config.EmptyOtherSceneVal {
// 						continue
// 					}

// 					if s == symbolVal.SymbolCode {
// 						cv, err := symbolVal.WeightVal.RandVal(plugin)
// 						if err != nil {
// 							goutils.Error("SymbolModifier.OnPlayGame:WeightVal.RandVal",
// 								zap.Error(err))

// 							return err
// 						}

// 						os.Arr[x][y] = cv.Int()
// 					} else {
// 						os.Arr[x][y] = symbolVal.Config.DefaultVal
// 					}
// 				}
// 			}

// 			symbolVal.AddOtherScene(gameProp, curpr, os, cd)

// 			if symbolVal.OtherSceneFeature != nil {
// 				gameProp.procOtherSceneFeature(symbolVal.OtherSceneFeature, curpr, os)
// 			}
// 		}
// 	}

// 	symbolVal.onStepEnd(gameProp, curpr, gp, "")

// 	// gp.AddComponentData(symbolVal.Name, cd)
// 	// symbolMulti.BuildPBComponent(gp)

// 	return nil
// }

// // OnAsciiGame - outpur to asciigame
// func (symbolVal *SymbolModifier) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

// 	cd := gameProp.MapComponentData[symbolVal.Name].(*BasicComponentData)

// 	if len(cd.UsedOtherScenes) > 0 {
// 		asciigame.OutputOtherScene("The value of the symbols", pr.OtherScenes[cd.UsedOtherScenes[0]])
// 	}

// 	return nil
// }

// // OnStats
// func (symbolVal *SymbolModifier) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// func NewSymbolModifier(name string) IComponent {
// 	return &SymbolModifier{
// 		BasicComponent: NewBasicComponent(name),
// 	}
// }
