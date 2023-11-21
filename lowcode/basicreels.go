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

const BasicReelsTypeName = "basicReels"

// BasicReelsConfig - configuration for BasicReels
type BasicReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	ReelSetsWeight       string `yaml:"reelSetWeight" json:"reelSetWeight"`
	ReelSet              string `yaml:"reelSet" json:"reelSet"`
	IsExpandReel         bool   `yaml:"isExpandReel" json:"isExpandReel"`
}

type BasicReels struct {
	*BasicComponent `json:"-"`
	Config          *BasicReelsConfig     `json:"config"`
	ReelSetWeights  *sgc7game.ValWeights2 `json:"-"`
}

// Init -
func (basicReels *BasicReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &BasicReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return basicReels.InitEx(cfg, pool)
}

// InitEx -
func (basicReels *BasicReels) InitEx(cfg any, pool *GamePropertyPool) error {
	basicReels.Config = cfg.(*BasicReelsConfig)
	basicReels.Config.ComponentType = BasicReelsTypeName

	if basicReels.Config.ReelSetsWeight != "" {
		vw2, err := pool.LoadStrWeights(basicReels.Config.ReelSetsWeight, basicReels.Config.UseFileMapping)
		if err != nil {
			goutils.Error("BasicReels.Init:LoadValWeights",
				zap.String("ReelSetsWeight", basicReels.Config.ReelSetsWeight),
				zap.Error(err))

			return err
		}

		basicReels.ReelSetWeights = vw2
	}

	basicReels.onInit(&basicReels.Config.BasicComponentConfig)

	return nil
}

// playgame
func (basicReels *BasicReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	basicReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[basicReels.Name].(*BasicComponentData)

	reelname := ""
	if basicReels.ReelSetWeights != nil {
		val, si, err := basicReels.ReelSetWeights.RandValEx(plugin)
		if err != nil {
			goutils.Error("BasicReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Error(err))

			return err
		}

		basicReels.AddRNG(gameProp, si, cd)

		curreels := val.String()
		gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("BasicReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
		reelname = curreels
	} else {
		rd, isok := gameProp.Pool.Config.MapReels[basicReels.Config.ReelSet]
		if !isok {
			goutils.Error("BasicReels.OnPlayGame:MapReels",
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.TagStr(TagCurReels, basicReels.Config.ReelSet)

		gameProp.CurReels = rd
		reelname = basicReels.Config.ReelSet
	}

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), false)
	sc.ReelName = reelname
	// sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	// if err != nil {
	// 	goutils.Error("BasicReels.OnPlayGame:NewGameScene",
	// 		zap.Error(err))

	// 	return err
	// }

	if basicReels.Config.IsExpandReel {
		sc.RandExpandReelsWithReelData(gameProp.CurReels, plugin)
	} else {
		sc.RandReelsWithReelData(gameProp.CurReels, plugin)
	}

	basicReels.AddScene(gameProp, curpr, sc, cd)

	basicReels.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(basicReels.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (basicReels *BasicReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[basicReels.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (basicReels *BasicReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewBasicReels(name string) IComponent {
	basicReels := &BasicReels{
		BasicComponent: NewBasicComponent(name),
	}

	return basicReels
}
