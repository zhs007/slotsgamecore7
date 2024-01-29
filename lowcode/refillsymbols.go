package lowcode

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &RefillSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RefillSymbols.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	refillSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[refillSymbols.Name].(*BasicComponentData)

	gs := refillSymbols.GetTargetScene3(gameProp, curpr, prs, cd, refillSymbols.Name, "", 0)
	ngs := gs

	for x := 0; x < gs.Width; x++ {
		for y := gs.Width - 1; y >= 0; y-- {
			if ngs.Arr[x][y] == -1 {
				if ngs == gs {
					ngs = gs.Clone()
				}

				cr := gameProp.Pool.Config.MapReels[ngs.ReelName]

				ngs.Arr[x][y] = refillSymbols.getSymbol(cr, x, ngs.Indexes[x])
				ngs.Indexes[x]--
			}
		}
	}

	if ngs == gs {
		refillSymbols.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	refillSymbols.AddScene(gameProp, curpr, ngs, cd)

	refillSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (refillSymbols *RefillSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[refillSymbols.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after refillSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (refillSymbols *RefillSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (refillSymbols *RefillSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	return 0, nil
}

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

	cfg.UseSceneV3 = true

	return cfg
}

func parseRefillSymbols(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRefillSymbols:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRefillSymbols:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonRefillSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRefillSymbols:Unmarshal",
			zap.Error(err))

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
