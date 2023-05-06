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

// ReplaceSymbolConfig - configuration for ReplaceSymbol
type ReplaceSymbolConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Symbols              []string `yaml:"symbols"`
	Chg2SymbolInReels    []string `yaml:"chg2SymbolInReels"`
	Mask                 string   `yaml:"mask"`
}

type ReplaceSymbol struct {
	*BasicComponent
	Config                *ReplaceSymbolConfig
	SymbolCodes           []int
	Chg2SymbolCodeInReels []int
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

	replaceSymbol.Config = cfg

	for _, v := range cfg.Symbols {
		replaceSymbol.SymbolCodes = append(replaceSymbol.SymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range cfg.Chg2SymbolInReels {
		replaceSymbol.Chg2SymbolCodeInReels = append(replaceSymbol.Chg2SymbolCodeInReels, pool.DefaultPaytables.MapSymbols[v])
	}

	replaceSymbol.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (replaceSymbol *ReplaceSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[replaceSymbol.Name].(*BasicComponentData)

	gs := replaceSymbol.GetTargetScene(gameProp, curpr, cd, "")

	if !gs.HasSymbols(replaceSymbol.SymbolCodes) {
		replaceSymbol.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, cd)
	} else {
		sc2 := gs.Clone()

		if replaceSymbol.Config.Mask != "" {
			md := gameProp.MapComponentData[replaceSymbol.Config.Mask].(*MaskData)
			if md != nil {
				for x, arr := range sc2.Arr {
					if md.Vals[x] {
						for y, s := range arr {
							if goutils.IndexOfIntSlice(replaceSymbol.SymbolCodes, s, 0) >= 0 {
								sc2.Arr[x][y] = replaceSymbol.Chg2SymbolCodeInReels[x]
							}
						}
					}
				}
			}
		} else {
			for x, arr := range sc2.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(replaceSymbol.SymbolCodes, s, 0) >= 0 {
						sc2.Arr[x][y] = replaceSymbol.Chg2SymbolCodeInReels[x]
					}
				}
			}
		}

		replaceSymbol.AddScene(gameProp, curpr, sc2, cd)
	}

	replaceSymbol.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(replaceSymbol.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (replaceSymbol *ReplaceSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[replaceSymbol.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (replaceSymbol *ReplaceSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewReplaceSymbol(name string) IComponent {
	replaceSymbol := &ReplaceSymbol{
		BasicComponent: NewBasicComponent(name),
	}

	return replaceSymbol
}
