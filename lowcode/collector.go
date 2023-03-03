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

// MultiLevelReelsConfig - configuration for Collecotr
type CollecotrConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MaxVal               int `yaml:"maxVal"`
}

type Collecotr struct {
	*BasicComponent
	Config       *CollecotrConfig
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
}

// Init -
func (collector *Collecotr) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &CollecotrConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	collector.Config = cfg

	return nil
}

// OnNewGame -
func (collector *Collecotr) OnNewGame(gameProp *GameProperty) error {
	collector.Val = 0

	return nil
}

// OnNewStep -
func (collector *Collecotr) OnNewStep(gameProp *GameProperty) error {

	collector.BasicComponent.OnNewStep()

	collector.NewCollector = 0

	return nil
}

// playgame
func (collector *Collecotr) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	gameProp.SetStrVal(GamePropNextComponent, collector.Config.DefaultNextComponent)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (collector *Collecotr) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	if collector.NewCollector <= 0 {
		fmt.Printf("%v dose not collect new value, the collector value is %v", collector.Name, collector.Val)
	} else {
		fmt.Printf("%v collect %v. the collector value is %v", collector.Name, collector.NewCollector, collector.Val)
	}

	return nil
}

func NewCollector(name string) IComponent {
	collector := &Collecotr{
		BasicComponent: NewBasicComponent(name),
	}

	return collector
}
