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
	"gopkg.in/yaml.v2"
)

const ChgSymbolsTypeName = "chgSymbols"

// ChgSymbolsConfig - configuration for ChgSymbols
type ChgSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbols              []string              `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int                 `yaml:"-" json:"-"`
	BlankSymbol          string                `yaml:"blankSymbol" json:"blankSymbol"`
	BlankSymbolCode      int                   `yaml:"-" json:"-"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVW2            *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	Controllers          []*Award              `yaml:"controllers" json:"controllers"`
	JumpToComponent      string                `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *ChgSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
}

type ChgSymbols struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolsConfig `json:"config"`
}

// Init -
func (chgSymbols *ChgSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbols.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbols *ChgSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbols.Config = cfg.(*ChgSymbolsConfig)
	chgSymbols.Config.ComponentType = ChgSymbolsTypeName

	for _, s := range chgSymbols.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("ChgSymbols.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		chgSymbols.Config.SymbolCodes = append(chgSymbols.Config.SymbolCodes, sc)
	}

	blankSymbolCode, isok := pool.DefaultPaytables.MapSymbols[chgSymbols.Config.BlankSymbol]
	if isok {
		chgSymbols.Config.BlankSymbolCode = blankSymbolCode
	} else {
		chgSymbols.Config.BlankSymbolCode = -1
	}

	if chgSymbols.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(chgSymbols.Config.Weight, chgSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols.InitEx:LoadIntWeights",
				slog.String("Weight", chgSymbols.Config.Weight),
				goutils.Err(err))

			return err
		}

		chgSymbols.Config.WeightVW2 = vw2
	}

	for _, award := range chgSymbols.Config.Controllers {
		award.Init()
	}

	chgSymbols.onInit(&chgSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (chgSymbols *ChgSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := chgSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs != nil {
		ngs := gs

		for x, arr := range gs.Arr {
			for y, s := range arr {
				if goutils.IndexOfIntSlice(chgSymbols.Config.SymbolCodes, s, 0) >= 0 {
					curs, err := chgSymbols.Config.WeightVW2.RandVal(plugin)
					if err != nil {
						goutils.Error("ChgSymbols.OnPlayGame:RandVal",
							goutils.Err(err))

						return "", err
					}

					cursc := curs.Int()
					if cursc != chgSymbols.Config.BlankSymbolCode {
						if ngs == gs {
							ngs = gs.CloneEx(gameProp.PoolScene)
						}

						ngs.Arr[x][y] = cursc
					}
				}
			}
		}

		if ngs == gs {
			nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		chgSymbols.AddScene(gameProp, curpr, ngs, cd)

		if len(chgSymbols.Config.Controllers) > 0 {
			gameProp.procAwards(plugin, chgSymbols.Config.Controllers, curpr, gp)
		}

		nc := chgSymbols.onStepEnd(gameProp, curpr, gp, chgSymbols.Config.JumpToComponent)

		return nc, nil
	}

	nc := chgSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (chgSymbols *ChgSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after ChgSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (chgSymbols *ChgSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// GetAllLinkComponents - get all link components
func (chgSymbols *ChgSymbols) GetAllLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (chgSymbols *ChgSymbols) GetNextLinkComponents() []string {
	return []string{chgSymbols.Config.DefaultNextComponent, chgSymbols.Config.JumpToComponent}
}

func NewChgSymbols(name string) IComponent {
	return &ChgSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "symbols": [
//
//	"E"
//
// ],
// "blankSymbol": "BN",
// "weight": "bgchgsymweight"
type jsonChgSymbols struct {
	Symbols     []string `json:"symbols"`
	BlankSymbol string   `yaml:"blankSymbol" json:"blankSymbol"`
	Weight      string   `yaml:"weight" json:"weight"`
}

func (jcfg *jsonChgSymbols) build() *ChgSymbolsConfig {
	cfg := &ChgSymbolsConfig{
		Symbols:     jcfg.Symbols,
		BlankSymbol: jcfg.BlankSymbol,
		Weight:      jcfg.Weight,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseChgSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseClusterTrigger:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Controllers = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ChgSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
