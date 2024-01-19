package lowcode

import (
	"context"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"github.com/zhs007/slotsgamecore7/stats2"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const SymbolModifierTypeName = "symbolModifier"

// SymbolModifierConfig - configuration for SymbolModifier feature
type SymbolModifierConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int    `yaml:"-" json:"-"`
	TargetSymbols        []string `yaml:"targetSymbols" json:"targetSymbols"`
	TargetSymbolCodes    []int    `yaml:"-" json:"-"`
	Triggers             []string `yaml:"triggers" json:"triggers"` // 替换完图标后需要保证所有trigger返回true
	MinNum               int      `yaml:"minNum" json:"minNum"`     // 至少换几个，如果小于等于0，就表示要全部换
	PosArea              []int    `yaml:"posArea" json:"posArea"`   // 只在countscatterInArea时生效，[minx,maxx,miny,maxy]，当x，y分别符合双闭区间才合法
}

type SymbolModifier struct {
	*BasicComponent `json:"-"`
	Config          *SymbolModifierConfig `json:"config"`
}

// Init -
func (symbolModifier *SymbolModifier) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolModifier.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolModifierConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolModifier.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolModifier.InitEx(cfg, pool)
}

// InitEx -
func (symbolModifier *SymbolModifier) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolModifier.Config = cfg.(*SymbolModifierConfig)
	symbolModifier.Config.ComponentType = SymbolValTypeName

	for _, s := range symbolModifier.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolModifier.InitEx:Symbols",
				zap.String("symbol", s),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		symbolModifier.Config.SymbolCodes = append(symbolModifier.Config.SymbolCodes, sc)
	}

	for _, s := range symbolModifier.Config.TargetSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("SymbolModifier.InitEx:TargetSymbols",
				zap.String("symbol", s),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		symbolModifier.Config.TargetSymbolCodes = append(symbolModifier.Config.TargetSymbolCodes, sc)
	}

	symbolModifier.onInit(&symbolModifier.Config.BasicComponentConfig)

	return nil
}

// getSymbols
func (symbolModifier *SymbolModifier) getSymbols(gs *sgc7game.GameScene) []int {
	lst := []int{}

	if len(symbolModifier.Config.PosArea) == 4 {
		for x, arr := range gs.Arr {
			for y, s := range arr {
				if IsInPosArea(x, y, symbolModifier.Config.PosArea) {
					if goutils.IndexOfIntSlice(symbolModifier.Config.SymbolCodes, s, 0) >= 0 {
						lst = append(lst, x, y)
					}
				}
			}
		}
	} else {
		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(symbolModifier.Config.SymbolCodes, s, 0) >= 0 {
					lst = append(lst, x, y)
				}
			}
		}
	}

	return lst
}

// procSymbolsRandPos
func (symbolModifier *SymbolModifier) canModify(gameProp *GameProperty, gs *sgc7game.GameScene, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	for _, st := range symbolModifier.Config.Triggers {
		if !gameProp.CanTrigger(st, gs, curpr, stake) {
			return false
		}
	}

	return true
}

// procSymbolsRandPos
func (symbolModifier *SymbolModifier) chgSymbols(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, x, y int, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	srcs := gs.Arr[x][y]

	lst := make([]int, len(symbolModifier.Config.TargetSymbolCodes))
	copy(lst, symbolModifier.Config.TargetSymbolCodes)

	for {
		cr, err := plugin.Random(context.Background(), len(lst))
		if err != nil {
			goutils.Error("SymbolModifier.chgSymbols:random symbols",
				zap.Error(err))

			break
		}

		gs.Arr[x][y] = lst[cr]

		if symbolModifier.canModify(gameProp, gs, curpr, stake) {
			return true
		}

		if len(lst) == 1 {
			break
		}

		lst = append(lst[0:cr], lst[cr+1:]...)
	}

	gs.Arr[x][y] = srcs

	return false
}

// procSymbolsRandPos
func (symbolModifier *SymbolModifier) procSymbolsRandPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, lst []int, minnum int, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	if minnum > len(lst)/2 {
		return false
	}

	ci, err := plugin.Random(context.Background(), len(lst)/2)
	if err != nil {
		goutils.Error("SymbolModifier.procSymbolsRandPos:random pos",
			zap.Error(err))

		return false
	}

	x := lst[ci*2]
	y := lst[ci*2+1]

	isok := symbolModifier.chgSymbols(gameProp, plugin, gs, x, y, curpr, stake)
	if isok {
		minnum--

		if minnum == 0 {
			return true
		}

		if ci == len(lst)/2-1 {
			lst = lst[0 : ci*2]
		} else {
			lst = append(lst[0:ci*2], lst[(ci+1)*2:]...)
		}

		return symbolModifier.procSymbolsRandPos(gameProp, plugin, gs, lst, minnum, curpr, stake)
	}

	if ci == len(lst)/2-1 {
		lst = lst[0 : ci*2]
	} else {
		lst = append(lst[0:ci*2], lst[(ci+1)*2:]...)
	}

	if minnum < len(lst)/2 {
		return false
	}

	return symbolModifier.procSymbolsRandPos(gameProp, plugin, gs, lst, minnum, curpr, stake)
}

// procSymbols
func (symbolModifier *SymbolModifier) procSymbols(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene, lst []int, minnum int, curpr *sgc7game.PlayResult, stake *sgc7game.Stake) bool {
	if minnum > len(lst)/2 {
		return false
	}

	ci := 0

	x := lst[ci*2]
	y := lst[ci*2+1]

	isok := symbolModifier.chgSymbols(gameProp, plugin, gs, x, y, curpr, stake)
	if isok {
		minnum--

		if minnum == 0 {
			return true
		}

		if ci == len(lst)/2-1 {
			lst = lst[0 : ci*2]
		} else {
			lst = append(lst[0:ci*2], lst[(ci+1)*2:]...)
		}

		return symbolModifier.procSymbolsRandPos(gameProp, plugin, gs, lst, minnum, curpr, stake)
	}

	if ci == len(lst)/2-1 {
		lst = lst[0 : ci*2]
	} else {
		lst = append(lst[0:ci*2], lst[(ci+1)*2:]...)
	}

	if minnum < len(lst)/2 {
		return false
	}

	return symbolModifier.procSymbolsRandPos(gameProp, plugin, gs, lst, minnum, curpr, stake)
}

// playgame
func (symbolModifier *SymbolModifier) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolModifier.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[symbolModifier.Name].(*BasicComponentData)

	if len(symbolModifier.Config.TargetSymbolCodes) > 0 {
		gs := symbolModifier.GetTargetScene3(gameProp, curpr, cd, symbolModifier.Name, "", 0)

		lst := symbolModifier.getSymbols(gs)
		if len(lst) > 0 {
			minnum := len(lst) / 2

			// 这个分支将决定是否随机选symbol
			if symbolModifier.Config.MinNum > 0 {
				minnum = symbolModifier.Config.MinNum

				if minnum <= len(lst)/2 {
					gs1 := gs.CloneEx(gameProp.PoolScene)

					isok := symbolModifier.procSymbolsRandPos(gameProp, plugin, gs1, lst, minnum, curpr, stake)
					if isok {
						symbolModifier.AddScene(gameProp, curpr, gs1, cd)
					}
				}
			} else {
				if minnum <= len(lst)/2 {
					gs1 := gs.CloneEx(gameProp.PoolScene)

					isok := symbolModifier.procSymbols(gameProp, plugin, gs1, lst, minnum, curpr, stake)
					if isok {
						symbolModifier.AddScene(gameProp, curpr, gs1, cd)
					}
				}
			}
		}
	}

	symbolModifier.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolModifier *SymbolModifier) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[symbolModifier.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("symbolModifier symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (symbolModifier *SymbolModifier) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewStats2 -
func (symbolModifier *SymbolModifier) NewStats2() *stats2.Stats {
	return stats2.NewStats(stats2.Options{stats2.OptStepTrigger})
}

// OnStats2
func (symbolModifier *SymbolModifier) OnStats2(icd IComponentData, s2 *Stats2) {
	s2.pushStepTrigger(symbolModifier.Name, true)
}

// // OnStats2Trigger
// func (symbolModifier *SymbolModifier) OnStats2Trigger(s2 *Stats2) {
// 	s2.pushTriggerStats(symbolModifier.Name, true)
// }

func NewSymbolModifier(name string) IComponent {
	return &SymbolModifier{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
