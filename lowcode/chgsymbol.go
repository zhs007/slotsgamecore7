package lowcode

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const ChgSymbolTypeName = "chgSymbol"

// ChgSymbolNodeConfig -
type ChgSymbolNodeConfig struct {
	X          int    `yaml:"x" json:"x"`
	Y          int    `yaml:"y" json:"y"`
	Symbol     string `yaml:"symbol" json:"symbol"`
	SymbolCode int    `yaml:"symbolCode" json:"symbolCode"`
}

// ChgSymbolConfig - configuration for ChgSymbol feature
type ChgSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Nodes                []*ChgSymbolNodeConfig `yaml:"nodes" json:"nodes"`
}

type ChgSymbol struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolConfig `json:"config"`
}

// Init -
func (chgSymbol *ChgSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbol.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbol *ChgSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbol.Config = cfg.(*ChgSymbolConfig)
	chgSymbol.Config.ComponentType = ChgSymbolTypeName

	for _, v := range chgSymbol.Config.Nodes {
		v.SymbolCode = pool.DefaultPaytables.MapSymbols[v.Symbol]
	}

	chgSymbol.onInit(&chgSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (chgSymbol *ChgSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// chgSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	gs := chgSymbol.GetTargetScene3(gameProp, curpr, prs, 0)

	cgs := gs.CloneEx(gameProp.PoolScene)
	// cgs := gs.Clone()

	for _, v := range chgSymbol.Config.Nodes {
		cgs.Arr[v.X][v.Y] = v.SymbolCode
	}

	chgSymbol.AddScene(gameProp, curpr, cgs, bcd)

	nc := chgSymbol.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(chgSymbol.Name, cd)
	// symbolMulti.BuildPBComponent(gp)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (chgSymbol *ChgSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbol", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (chgSymbol *ChgSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewChgSymbol(name string) IComponent {
	return &ChgSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
