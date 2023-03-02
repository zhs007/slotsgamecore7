package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// LightningConfig - configuration for Lightning
type LightningConfig struct {
	BasicComponentConfig `yaml:",inline"`
}

type Lightning struct {
	*BasicComponent
	Config      *LightningConfig
	UsedScenes  []int
	UsedResults []int
}

// AddScene -
func (lightning *Lightning) AddScene(curpr *sgc7game.PlayResult, sc *sgc7game.GameScene) {
	lightning.UsedScenes = append(lightning.UsedScenes, len(curpr.Scenes))

	curpr.Scenes = append(curpr.Scenes, sc)
}

// AddResult -
func (lightning *Lightning) AddResult(curpr *sgc7game.PlayResult, ret *sgc7game.Result) {
	curpr.CashWin += int64(ret.CashWin)
	curpr.CoinWin += ret.CoinWin

	lightning.UsedResults = append(lightning.UsedResults, len(curpr.Results))

	curpr.Results = append(curpr.Results, ret)
}

// Init -
func (lightning *Lightning) Init(fn string, gameProp *GameProperty) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &LightningConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	lightning.Config = cfg

	return nil
}

// playgame
func (lightning *Lightning) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	gameProp.SetStrVal(GamePropNextComponent, lightning.Config.DefaultNextComponent)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (lightning *Lightning) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	return nil
}

func NewLightning(name string) IComponent {
	return &Lightning{
		BasicComponent: NewBasicComponent(name),
	}
}
