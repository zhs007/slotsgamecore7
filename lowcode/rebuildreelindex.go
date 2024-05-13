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

const RebuildReelIndexTypeName = "rebuildReelIndex"

type RebuildReelIndexType int

const (
	RebuildReelIndexTypeCircle RebuildReelIndexType = 0 // circle
	RebuildReelIndexTypeRandom RebuildReelIndexType = 1 // random
)

func parseRebuildReelIndexType(str string) RebuildReelIndexType {
	if str == "random" {
		return RebuildReelIndexTypeRandom
	}

	return RebuildReelIndexTypeCircle
}

// RebuildReelIndexConfig - configuration for RebuildReelIndex feature
type RebuildReelIndexConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string               `yaml:"type" json:"type"` // type
	Type                 RebuildReelIndexType `yaml:"-" json:"-"`       // type
}

// SetLinkComponent
func (cfg *RebuildReelIndexConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RebuildReelIndex struct {
	*BasicComponent `json:"-"`
	Config          *RebuildReelIndexConfig `json:"config"`
}

// Init -
func (rebuildReelIndex *RebuildReelIndex) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RebuildReelIndex.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RebuildReelIndexConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RebuildReelIndex.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return rebuildReelIndex.InitEx(cfg, pool)
}

// InitEx -
func (rebuildReelIndex *RebuildReelIndex) InitEx(cfg any, pool *GamePropertyPool) error {
	rebuildReelIndex.Config = cfg.(*RebuildReelIndexConfig)
	rebuildReelIndex.Config.ComponentType = RebuildReelIndexTypeName

	rebuildReelIndex.Config.Type = parseRebuildReelIndexType(rebuildReelIndex.Config.StrType)

	rebuildReelIndex.onInit(&rebuildReelIndex.Config.BasicComponentConfig)

	return nil
}

func (rebuildReelIndex *RebuildReelIndex) procCircle(gameProp *GameProperty, gs *sgc7game.GameScene, plugin sgc7plugin.IPlugin) (*sgc7game.GameScene, error) {
	cr, err := plugin.Random(context.Background(), gs.Width)
	if err != nil {
		goutils.Error("RebuildReelIndex.procCircle:Random",
			goutils.Err(err))

		return nil, err
	}

	if cr == 0 {
		return gs, nil
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range ngs.Arr {
		for y := range arr {
			ngs.Arr[x][y] = gs.Arr[cr][y]
		}

		cr++
		if cr >= gs.Width {
			cr = 0
		}
	}

	return ngs, nil
}

func (rebuildReelIndex *RebuildReelIndex) procRandom(gameProp *GameProperty, gs *sgc7game.GameScene, plugin sgc7plugin.IPlugin) (*sgc7game.GameScene, error) {
	arr := GenInitialArr(gs.Width)
	arr1, err := Shuffle(arr, plugin)
	if err != nil {
		goutils.Error("RebuildReelIndex.procRandom:Shuffle",
			goutils.Err(err))

		return nil, err
	}

	if IsInitialArr(arr1) {
		return gs, nil
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range ngs.Arr {
		for y := range arr {
			ngs.Arr[x][y] = gs.Arr[arr1[x]][y]
		}
	}

	return ngs, nil
}

// playgame
func (rebuildReelIndex *RebuildReelIndex) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// reelModifier.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil

	gs := rebuildReelIndex.GetTargetScene3(gameProp, curpr, prs, 0)
	var ngs *sgc7game.GameScene
	if rebuildReelIndex.Config.Type == RebuildReelIndexTypeCircle {
		gs1, err := rebuildReelIndex.procCircle(gameProp, gs, plugin)
		if err != nil {
			goutils.Error("RebuildReelIndex.OnPlayGame:procCircle",
				goutils.Err(err))

			return "", err
		}

		ngs = gs1
	} else {
		gs1, err := rebuildReelIndex.procRandom(gameProp, gs, plugin)
		if err != nil {
			goutils.Error("RebuildReelIndex.OnPlayGame:procRandom",
				goutils.Err(err))

			return "", err
		}

		ngs = gs1
	}

	if ngs == gs {
		nc := rebuildReelIndex.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	rebuildReelIndex.AddScene(gameProp, curpr, ngs, bcd)

	nc := rebuildReelIndex.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (rebuildReelIndex *RebuildReelIndex) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("rebuildReelIndex symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
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

func NewRebuildReelIndex(name string) IComponent {
	return &RebuildReelIndex{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"type": "cycle"
//	},
type jsonRebuildReelIndex struct {
	StrType string `json:"type"` // type
}

func (jcfg *jsonRebuildReelIndex) build() *RebuildReelIndexConfig {
	cfg := &RebuildReelIndexConfig{
		StrType: jcfg.StrType,
	}

	return cfg
}

func parseRebuildReelIndex(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRebuildReelIndex:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRebuildReelIndex:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRebuildReelIndex{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRebuildReelIndex:Unmarshal",
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
		Type: RebuildReelIndexTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
