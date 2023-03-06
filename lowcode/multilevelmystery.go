package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// MultiLevelMysteryLevelConfig - configuration for MultiLevelMystery's Level
type MultiLevelMysteryLevelConfig struct {
	MysteryWeight string `yaml:"mysteryWeight"`
	Collector     string `yaml:"collector"`
	CollectorVal  int    `yaml:"collectorVal"`
}

// MultiLevelMysteryConfig - configuration for MultiLevelMystery
type MultiLevelMysteryConfig struct {
	BasicComponentConfig   `yaml:",inline"`
	TargetScene            string                          `yaml:"targetScene"` // basicReels.init
	Mystery                string                          `yaml:"mystery"`
	Levels                 []*MultiLevelMysteryLevelConfig `yaml:"levels"`
	MysteryTriggerFeatures []*MysteryTriggerFeatureConfig  `yaml:"mysteryTriggerFeatures"`
}

type MultiLevelMystery struct {
	*BasicComponent
	Config                   *MultiLevelMysteryConfig
	MysterySymbol            int
	MapMysteryTriggerFeature map[int]*MysteryTriggerFeatureConfig
	LevelMysteryWeights      []*sgc7game.ValWeights2
	CurLevel                 int
}

// maskOtherScene -
func (multiLevelMystery *MultiLevelMystery) maskOtherScene(gs *sgc7game.GameScene, symbolCode int) *sgc7game.GameScene {
	cgs := gs.Clone()

	for x, arr := range cgs.Arr {
		for y, v := range arr {
			if v != symbolCode {
				cgs.Arr[x][y] = -1
			} else {
				cgs.Arr[x][y] = 1
			}
		}
	}

	return cgs
}

// Init -
func (multiLevelMystery *MultiLevelMystery) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiLevelMystery.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiLevelMysteryConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiLevelMystery.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiLevelMystery.Config = cfg

	for _, v := range cfg.Levels {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(v.MysteryWeight, "val", "weight", gameProp.CurPaytables)
		if err != nil {
			goutils.Error("MultiLevelMystery.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", v.MysteryWeight),
				zap.Error(err))

			return err
		}

		multiLevelMystery.LevelMysteryWeights = append(multiLevelMystery.LevelMysteryWeights, vw2)
	}

	multiLevelMystery.MysterySymbol = gameProp.CurPaytables.MapSymbols[multiLevelMystery.Config.Mystery]

	for _, v := range cfg.MysteryTriggerFeatures {
		symbolCode := gameProp.CurPaytables.MapSymbols[v.Symbol]

		multiLevelMystery.MapMysteryTriggerFeature[symbolCode] = v
	}

	multiLevelMystery.BasicComponent.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (multiLevelMystery *MultiLevelMystery) OnNewGame(gameProp *GameProperty) error {
	multiLevelMystery.CurLevel = 0

	return nil
}

// OnNewStep -
func (multiLevelMystery *MultiLevelMystery) OnNewStep(gameProp *GameProperty) error {
	multiLevelMystery.BasicComponent.OnNewStep()

	for i, v := range multiLevelMystery.Config.Levels {
		if multiLevelMystery.CurLevel > i {
			collecotr, isok := gameProp.MapCollectors[v.Collector]
			if isok {
				if collecotr.Val >= v.CollectorVal {
					multiLevelMystery.CurLevel = i
				}
			}
		}
	}

	return nil
}

// playgame
func (multiLevelMystery *MultiLevelMystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	gs := gameProp.GetScene(curpr, multiLevelMystery.Config.TargetScene)
	if gs.HasSymbol(multiLevelMystery.MysterySymbol) {
		curm, err := multiLevelMystery.LevelMysteryWeights[multiLevelMystery.CurLevel].RandVal(plugin)
		if err != nil {
			goutils.Error("MultiLevelMystery.OnPlayGame:RandVal",
				zap.Error(err))

			return err
		}

		curmcode := curm.Int()

		gameProp.SetVal(GamePropCurMystery, curm.Int())

		sc2 := gs.Clone()
		sc2.ReplaceSymbol(multiLevelMystery.MysterySymbol, curm.Int())

		multiLevelMystery.AddScene(gameProp, curpr, sc2)

		v, isok := multiLevelMystery.MapMysteryTriggerFeature[curmcode]
		if isok {
			if v.RespinFirstComponent != "" {
				os := multiLevelMystery.maskOtherScene(sc2, curmcode)

				gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

				return nil
			}
		}
	}

	multiLevelMystery.onStepEnd(gameProp, curpr, gp)

	multiLevelMystery.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelMystery *MultiLevelMystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(multiLevelMystery.UsedScenes) > 0 {
		fmt.Printf("mystery is %v\n", gameProp.GetStrVal(GamePropCurMystery))
		asciigame.OutputScene("after symbols", pr.Scenes[multiLevelMystery.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (multiLevelMystery *MultiLevelMystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMultiLevelMystery(name string) IComponent {
	multiLevelMystery := &MultiLevelMystery{
		BasicComponent:           NewBasicComponent(name),
		MapMysteryTriggerFeature: make(map[int]*MysteryTriggerFeatureConfig),
	}

	return multiLevelMystery
}
