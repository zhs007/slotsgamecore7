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

// MysteryTriggerFeatureConfig - configuration for mystery trigger feature
type MysteryTriggerFeatureConfig struct {
	Symbol               string `yaml:"symbol"`               // like LIGHTNING
	RespinFirstComponent string `yaml:"respinFirstComponent"` // like lightning
}

// MysteryConfig - configuration for Mystery
type MysteryConfig struct {
	BasicComponentConfig   `yaml:",inline"`
	TargetScene            string                         `yaml:"targetScene"` // basicReels.init
	MysteryWeight          string                         `yaml:"mysteryWeight"`
	Mystery                string                         `yaml:"mystery"`
	MysteryTriggerFeatures []*MysteryTriggerFeatureConfig `yaml:"mysteryTriggerFeatures"`
}

type Mystery struct {
	*BasicComponent
	Config                   *MysteryConfig
	MysteryWeights           *sgc7game.ValWeights2
	MysterySymbol            int
	MapMysteryTriggerFeature map[int]*MysteryTriggerFeatureConfig
}

// maskOtherScene -
func (mystery *Mystery) maskOtherScene(gs *sgc7game.GameScene, symbolCode int) *sgc7game.GameScene {
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
func (mystery *Mystery) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MysteryConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	mystery.Config = cfg

	if mystery.Config.MysteryWeight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(mystery.Config.MysteryWeight, "val", "weight", gameProp.CurPaytables)
		if err != nil {
			goutils.Error("BasicReels.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", mystery.Config.MysteryWeight),
				zap.Error(err))

			return err
		}

		mystery.MysteryWeights = vw2
	}

	mystery.MysterySymbol = gameProp.CurPaytables.MapSymbols[mystery.Config.Mystery]

	for _, v := range cfg.MysteryTriggerFeatures {
		symbolCode := gameProp.CurPaytables.MapSymbols[v.Symbol]

		mystery.MapMysteryTriggerFeature[symbolCode] = v
	}

	mystery.BasicComponent.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (mystery *Mystery) OnNewGame(gameProp *GameProperty) error {

	return nil
}

// OnNewStep -
func (mystery *Mystery) OnNewStep(gameProp *GameProperty) error {

	mystery.BasicComponent.OnNewStep()

	return nil
}

// playgame
func (mystery *Mystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	if mystery.MysteryWeights != nil {
		gs := gameProp.GetScene(curpr, mystery.Config.TargetScene)
		if gs.HasSymbol(mystery.MysterySymbol) {
			curm, err := mystery.MysteryWeights.RandVal(plugin)
			if err != nil {
				goutils.Error("BasicReels.OnPlayGame:RandVal",
					zap.Error(err))

				return err
			}

			curmcode := curm.Int()

			gameProp.SetVal(GamePropCurMystery, curm.Int())

			sc2 := gs.Clone()
			sc2.ReplaceSymbol(mystery.MysterySymbol, curm.Int())

			mystery.AddScene(gameProp, curpr, sc2)

			v, isok := mystery.MapMysteryTriggerFeature[curmcode]
			if isok {
				if v.RespinFirstComponent != "" {
					os := mystery.maskOtherScene(sc2, curmcode)

					gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

					return nil
				}
			}
		}
	}

	mystery.onStepEnd(gameProp, curpr, gp)

	mystery.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (mystery *Mystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(mystery.UsedScenes) > 0 {
		if mystery.MysteryWeights != nil {
			fmt.Printf("mystery is %v\n", gameProp.GetStrVal(GamePropCurMystery))
			asciigame.OutputScene("after symbols", pr.Scenes[mystery.UsedScenes[0]], mapSymbolColor)
		}
	}

	return nil
}

// OnStats
func (mystery *Mystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMystery(name string) IComponent {
	mystery := &Mystery{
		BasicComponent:           NewBasicComponent(name),
		MapMysteryTriggerFeature: make(map[int]*MysteryTriggerFeatureConfig),
	}

	return mystery
}
