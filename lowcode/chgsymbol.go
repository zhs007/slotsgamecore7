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

// ChgSymbolNodeConfig -
type ChgSymbolNodeConfig struct {
	X          int    `yaml:"x"`
	Y          int    `yaml:"y"`
	Symbol     string `yaml:"symbol"`
	SymbolCode int    `yaml:"symbolCode"`
}

// ChgSymbolConfig - configuration for ChgSymbol feature
type ChgSymbolConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Nodes                []*ChgSymbolNodeConfig `yaml:"nodes"`
}

type ChgSymbol struct {
	*BasicComponent
	Config *ChgSymbolConfig
}

// Init -
func (chgSymbol *ChgSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbol.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ChgSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	chgSymbol.Config = cfg

	for _, v := range cfg.Nodes {
		v.SymbolCode = pool.DefaultPaytables.MapSymbols[v.Symbol]
	}

	chgSymbol.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (chgSymbol *ChgSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	chgSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[chgSymbol.Name].(*BasicComponentData)

	gs := chgSymbol.GetTargetScene(gameProp, curpr, cd, "")

	cgs := gs.Clone()

	for _, v := range chgSymbol.Config.Nodes {
		cgs.Arr[v.X][v.Y] = v.SymbolCode
	}

	chgSymbol.AddScene(gameProp, curpr, cgs, cd)

	chgSymbol.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(chgSymbol.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (chgSymbol *ChgSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[chgSymbol.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("The value of the symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (chgSymbol *ChgSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewChgSymbol(name string) IComponent {
	return &ChgSymbol{
		BasicComponent: NewBasicComponent(name),
	}
}
