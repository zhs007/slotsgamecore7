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

const ReplaceSymbolTypeName = "replaceSymbol"

// ReplaceSymbolConfig - configuration for ReplaceSymbol
type ReplaceSymbolConfig struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	Symbols                  []string       `yaml:"symbols" json:"symbols"`
	Chg2SymbolInReels        []string       `yaml:"chg2SymbolInReels" json:"chg2SymbolInReels"`
	MapChg2SymbolInReels     map[int]string `yaml:"mapChg2SymbolInReels" json:"mapChg2SymbolInReels"`
	Mask                     string         `yaml:"mask" json:"mask"`
	SymbolCodes              []int          `yaml:"-" json:"-"`
	MapChg2SymbolCodeInReels map[int]int    `yaml:"-" json:"-"`
}

type ReplaceSymbol struct {
	*BasicComponent `json:"-"`
	Config          *ReplaceSymbolConfig `json:"config"`
}

// Init -
func (replaceSymbol *ReplaceSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReplaceSymbol.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReplaceSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReplaceSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return replaceSymbol.InitEx(cfg, pool)
}

// InitEx -
func (replaceSymbol *ReplaceSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	replaceSymbol.Config = cfg.(*ReplaceSymbolConfig)
	replaceSymbol.Config.ComponentType = ReplaceSymbolTypeName

	for _, v := range replaceSymbol.Config.Symbols {
		replaceSymbol.Config.SymbolCodes = append(replaceSymbol.Config.SymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	replaceSymbol.Config.MapChg2SymbolCodeInReels = make(map[int]int)

	for i, v := range replaceSymbol.Config.Chg2SymbolInReels {
		replaceSymbol.Config.MapChg2SymbolCodeInReels[i] = pool.DefaultPaytables.MapSymbols[v]
	}

	for k, v := range replaceSymbol.Config.MapChg2SymbolInReels {
		replaceSymbol.Config.MapChg2SymbolCodeInReels[k] = pool.DefaultPaytables.MapSymbols[v]
	}

	replaceSymbol.onInit(&replaceSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (replaceSymbol *ReplaceSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := replaceSymbol.GetTargetScene3(gameProp, curpr, prs, 0)

	if !gs.HasSymbols(replaceSymbol.Config.SymbolCodes) {
		// replaceSymbol.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, cd)
	} else {
		// sc2 := gs.Clone()
		sc2 := gs.CloneEx(gameProp.PoolScene)

		if replaceSymbol.Config.Mask != "" {
			md := gameProp.GetCurComponentDataWithName(replaceSymbol.Config.Mask).(*MaskData)
			// md := gameProp.MapComponentData[replaceSymbol.Config.Mask].(*MaskData)
			if md != nil {
				for x, arr := range sc2.Arr {
					if md.Vals[x] {
						destSymbol, isok := replaceSymbol.Config.MapChg2SymbolCodeInReels[x]
						if isok {
							for y, s := range arr {
								if goutils.IndexOfIntSlice(replaceSymbol.Config.SymbolCodes, s, 0) >= 0 {
									sc2.Arr[x][y] = destSymbol
								}
							}
						}
					}
				}
			}
		} else {
			for x, arr := range sc2.Arr {
				for y, s := range arr {
					destSymbol, isok := replaceSymbol.Config.MapChg2SymbolCodeInReels[x]
					if isok {
						if goutils.IndexOfIntSlice(replaceSymbol.Config.SymbolCodes, s, 0) >= 0 {
							sc2.Arr[x][y] = destSymbol
						}
					}
				}
			}
		}

		replaceSymbol.AddScene(gameProp, curpr, sc2, cd)
	}

	nc := replaceSymbol.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(replaceSymbol.Name, cd)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (replaceSymbol *ReplaceSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (replaceSymbol *ReplaceSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewReplaceSymbol(name string) IComponent {
	replaceSymbol := &ReplaceSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}

	return replaceSymbol
}
