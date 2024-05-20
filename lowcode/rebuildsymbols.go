package lowcode

import (
	"context"
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

const RebuildSymbolsTypeName = "rebuildSymbols"

type RebuildSymbolsType int

const (
	RebuildSymbolsTypeCircle RebuildSymbolsType = 0 // circle
	RebuildSymbolsTypeRandom RebuildSymbolsType = 1 // random
)

func parseRebuildSymbolsType(str string) RebuildSymbolsType {
	if str == "random" {
		return RebuildSymbolsTypeRandom
	}

	return RebuildSymbolsTypeCircle
}

// RebuildSymbolsConfig - configuration for RebuildSymbols feature
type RebuildSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string             `yaml:"type" json:"type"` // type
	Type                 RebuildSymbolsType `yaml:"-" json:"-"`       // type
	Symbols              []string           `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int              `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *RebuildSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RebuildSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RebuildSymbolsConfig `json:"config"`
}

// Init -
func (rebuildSymbols *RebuildSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RebuildSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RebuildSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RebuildSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return rebuildSymbols.InitEx(cfg, pool)
}

// InitEx -
func (rebuildSymbols *RebuildSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	rebuildSymbols.Config = cfg.(*RebuildSymbolsConfig)
	rebuildSymbols.Config.ComponentType = RebuildSymbolsTypeName

	rebuildSymbols.Config.Type = parseRebuildSymbolsType(rebuildSymbols.Config.StrType)

	for _, s := range rebuildSymbols.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("RebuildSymbols.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		rebuildSymbols.Config.SymbolCodes = append(rebuildSymbols.Config.SymbolCodes, sc)
	}

	rebuildSymbols.onInit(&rebuildSymbols.Config.BasicComponentConfig)

	return nil
}

func (rebuildSymbols *RebuildSymbols) procCircle(gameProp *GameProperty, gs *sgc7game.GameScene, plugin sgc7plugin.IPlugin) (*sgc7game.GameScene, error) {
	cr, err := plugin.Random(context.Background(), len(rebuildSymbols.Config.SymbolCodes))
	if err != nil {
		goutils.Error("RebuildSymbols.procCircle:Random",
			goutils.Err(err))

		return nil, err
	}

	if cr == 0 {
		return gs, nil
	}

	lst := make([]int, len(rebuildSymbols.Config.SymbolCodes))

	for i := 0; i < len(rebuildSymbols.Config.SymbolCodes); i++ {
		lst[i] = rebuildSymbols.Config.SymbolCodes[cr+i]
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range ngs.Arr {
		for y, v := range arr {
			srci := goutils.IndexOfIntSlice(rebuildSymbols.Config.SymbolCodes, v, 0)
			ngs.Arr[x][y] = lst[srci]
		}
	}

	return ngs, nil
}

func (rebuildSymbols *RebuildSymbols) procRandom(gameProp *GameProperty, gs *sgc7game.GameScene, plugin sgc7plugin.IPlugin) (*sgc7game.GameScene, error) {
	lst, err := Shuffle(rebuildSymbols.Config.SymbolCodes, plugin)
	if err != nil {
		goutils.Error("RebuildSymbols.procRandom:Shuffle",
			goutils.Err(err))

		return nil, err
	}

	if IsSameIntArr(rebuildSymbols.Config.SymbolCodes, lst) {
		return gs, nil
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range ngs.Arr {
		for y, v := range arr {
			srci := goutils.IndexOfIntSlice(rebuildSymbols.Config.SymbolCodes, v, 0)
			ngs.Arr[x][y] = lst[srci]
		}
	}

	return ngs, nil
}

// playgame
func (rebuildSymbols *RebuildSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// reelModifier.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil

	gs := rebuildSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	var ngs *sgc7game.GameScene
	if rebuildSymbols.Config.Type == RebuildSymbolsTypeCircle {
		gs1, err := rebuildSymbols.procCircle(gameProp, gs, plugin)
		if err != nil {
			goutils.Error("RebuildSymbols.OnPlayGame:procCircle",
				goutils.Err(err))

			return "", err
		}

		ngs = gs1
	} else {
		gs1, err := rebuildSymbols.procRandom(gameProp, gs, plugin)
		if err != nil {
			goutils.Error("RebuildSymbols.OnPlayGame:procRandom",
				goutils.Err(err))

			return "", err
		}

		ngs = gs1
	}

	if ngs == gs {
		nc := rebuildSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	rebuildSymbols.AddScene(gameProp, curpr, ngs, bcd)

	nc := rebuildSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (rebuildSymbols *RebuildSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("rebuildSymbols symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (reelModifier *ReelModifier) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // NewStats2 -
// func (reelModifier *ReelModifier) NewStats2(parent string) *stats2.Feature {
// 	return stats2.NewFeature(parent, nil)
// }

// // OnStats2
// func (reelModifier *ReelModifier) OnStats2(icd IComponentData, s2 *stats2.Cache) {
// 	// s2.PushStepTrigger(reelModifier.Name, true)
// 	s2.ProcStatsTrigger(reelModifier.Name)
// }

// // OnStats2Trigger
// func (reelModifier *ReelModifier) OnStats2Trigger(s2 *Stats2) {
// 	s2.pushTriggerStats(reelModifier.Name, true)
// }

func NewRebuildSymbols(name string) IComponent {
	return &RebuildSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"type": "cycle"
//	},
type jsonRebuildSymbols struct {
	StrType string   `json:"type"` // type
	Symbols []string `json:"symbols"`
}

func (jcfg *jsonRebuildSymbols) build() *RebuildSymbolsConfig {
	cfg := &RebuildSymbolsConfig{
		StrType: jcfg.StrType,
		Symbols: jcfg.Symbols,
	}

	return cfg
}

func parseRebuildSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRebuildSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRebuildSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRebuildSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRebuildSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	// if ctrls != nil {
	// 	awards, err := parseControllers(ctrls)
	// 	if err != nil {
	// 		goutils.Error("parseRebuildReelIndex:parseControllers",
	// 			goutils.Err(err))

	// 		return "", err
	// 	}

	// 	cfgd.Awards = awards
	// }

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RebuildSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
