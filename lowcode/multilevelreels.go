package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// MultiLevelReelsLevelConfig - configuration for MultiLevelReels's Level
type MultiLevelReelsLevelConfig struct {
	Reel           string `yaml:"reel"`
	ReelSetsWeight string `yaml:"reelSetWeight"`
	Collector      string `yaml:"collector"`
	CollectorVal   int    `yaml:"collectorVal"`
}

// MultiLevelReelsConfig - configuration for MultiLevelReels
type MultiLevelReelsConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Levels               []*MultiLevelReelsLevelConfig `yaml:"levels"`
}

type MultiLevelReels struct {
	*BasicComponent
	Config              *MultiLevelReelsConfig
	LevelReelSetWeights []*sgc7game.ValWeights2
	CurLevel            int
}

// Init -
func (multiLevelReels *MultiLevelReels) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiLevelReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiLevelReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiLevelReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiLevelReels.Config = cfg

	for _, v := range cfg.Levels {
		if v.ReelSetsWeight != "" {
			vw2, err := sgc7game.LoadValWeights2FromExcel(v.ReelSetsWeight, "val", "weight", sgc7game.NewStrVal)
			if err != nil {
				goutils.Error("MultiLevelReels.Init:LoadValWeights2FromExcel",
					zap.String("ReelSetsWeight", v.ReelSetsWeight),
					zap.Error(err))

				return err
			}

			multiLevelReels.LevelReelSetWeights = append(multiLevelReels.LevelReelSetWeights, vw2)
		}
	}

	if len(multiLevelReels.LevelReelSetWeights) > 0 && len(multiLevelReels.LevelReelSetWeights) != len(cfg.Levels) {
		goutils.Error("MultiLevelReels.Init:check levels",
			zap.Int("reelsetLength", len(multiLevelReels.LevelReelSetWeights)),
			zap.Int("levelLength", len(cfg.Levels)),
			zap.Error(ErrIvalidMultiLevelReelsConfig))

		return ErrIvalidMultiLevelReelsConfig
	}

	return nil
}

// playgame
func (multiLevelReels *MultiLevelReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	multiLevelReels.OnNewStep()

	if multiLevelReels.LevelReelSetWeights != nil {
		val, err := multiLevelReels.LevelReelSetWeights[multiLevelReels.CurLevel].RandVal(plugin)
		if err != nil {
			goutils.Error("MultiLevelReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Int("curLevel", multiLevelReels.CurLevel),
				zap.Error(err))

			return err
		}

		rd, isok := gameProp.Config.MapReels[val.String()]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", multiLevelReels.CurLevel),
				zap.String("reelset", val.String()),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
	} else {
		rd, isok := gameProp.Config.MapReels[multiLevelReels.Config.Levels[multiLevelReels.CurLevel].Reel]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", multiLevelReels.CurLevel),
				zap.String("reelset", multiLevelReels.Config.Levels[multiLevelReels.CurLevel].Reel),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
	}

	sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	if err != nil {
		goutils.Error("MultiLevelReels.OnPlayGame:NewGameScene",
			zap.Error(err))

		return err
	}

	sc.RandReelsWithReelData(gameProp.CurReels, plugin)

	multiLevelReels.AddScene(gameProp, curpr, sc, fmt.Sprintf("%v.init", multiLevelReels.Name))

	gameProp.SetStrVal(GamePropNextComponent, multiLevelReels.Config.DefaultNextComponent)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelReels *MultiLevelReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(multiLevelReels.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[multiLevelReels.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

func NewMultiLevelReels(name string) IComponent {
	multiLevelReels := &MultiLevelReels{
		BasicComponent: NewBasicComponent(name),
	}

	return multiLevelReels
}
