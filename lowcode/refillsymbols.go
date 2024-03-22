package lowcode

import (
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RefillSymbolsTypeName = "refillSymbols"

// RefillSymbolsConfig - configuration for RefillSymbols
type RefillSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
}

// SetLinkComponent
func (cfg *RefillSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RefillSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RefillSymbolsConfig `json:"config"`
}

// Init -
func (refillSymbols *RefillSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RefillSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RefillSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return refillSymbols.InitEx(cfg, pool)
}

// InitEx -
func (refillSymbols *RefillSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	refillSymbols.Config = cfg.(*RefillSymbolsConfig)
	refillSymbols.Config.ComponentType = RefillSymbolsTypeName

	refillSymbols.onInit(&refillSymbols.Config.BasicComponentConfig)

	return nil
}

func (refillSymbols *RefillSymbols) getSymbol(rd *sgc7game.ReelsData, x int, index int) int {
	index--

	for ; index < 0; index += len(rd.Reels[x]) {
	}

	return rd.Reels[x][index]
}

// playgame
func (refillSymbols *RefillSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// refillSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	bcd.UsedScenes = nil

	gs := refillSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	for x := 0; x < gs.Width; x++ {
		for y := gs.Width - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

				ngs.Arr[x][y] = refillSymbols.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--
			}
		}
	}

	if ngs == gs {
		nc := refillSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	refillSymbols.AddScene(gameProp, curpr, ngs, bcd)

	nc := refillSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (refillSymbols *RefillSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after refillSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (refillSymbols *RefillSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (refillSymbols *RefillSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	return 0, nil
// }

// EachUsedResults -
func (refillSymbols *RefillSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewRefillSymbols(name string) IComponent {
	return &RefillSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "configuration": {},
type jsonRefillSymbols struct {
}

func (jcfg *jsonRefillSymbols) build() *RefillSymbolsConfig {
	cfg := &RefillSymbolsConfig{}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseRefillSymbols(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRefillSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RefillSymbolsTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
