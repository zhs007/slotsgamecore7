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

// CollectorConfig - configuration for Collector
type CollectorConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MaxVal               int `yaml:"maxVal"`
}

type Collector struct {
	*BasicComponent
	Config *CollectorConfig
}

// Init -
func (collector *Collector) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &CollectorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	collector.Config = cfg

	collector.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (collector *Collector) OnNewGame(gameProp *GameProperty) error {
	cd, isok := gameProp.MapCollectors[collector.Name]
	if isok {
		cd.onNewGame()
	}

	return nil
}

// OnNewStep -
func (collector *Collector) OnNewStep(gameProp *GameProperty) error {
	collector.BasicComponent.OnNewStep()

	cd, isok := gameProp.MapCollectors[collector.Name]
	if isok {
		cd.onNewStep()
	}

	return nil
}

// playgame
func (collector *Collector) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	gameProp.SetStrVal(GamePropNextComponent, collector.Config.DefaultNextComponent)

	collector.onStepEnd(gameProp, curpr, gp)

	collector.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (collector *Collector) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd, isok := gameProp.MapCollectors[collector.Name]
	if isok {
		if cd.NewCollector <= 0 {
			fmt.Printf("%v dose not collect new value, the collector value is %v", collector.Name, cd.Val)
		} else {
			fmt.Printf("%v collect %v. the collector value is %v", collector.Name, cd.NewCollector, cd.Val)
		}
	}

	return nil
}

// OnStats
func (collector *Collector) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewCollector(name string) IComponent {
	collector := &Collector{
		BasicComponent: NewBasicComponent(name),
	}

	return collector
}
