package lowcode

import (
	"context"
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

const AddSymbolsTypeName = "addSymbols"

// AddSymbolsConfig - configuration for AddSymbols
type AddSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string                `yaml:"symbol" json:"symbol"`
	SymbolCode           int                   `yaml:"-" json:"-"`
	SymbolNum            int                   `yaml:"symbolNum" json:"symbolNum"`
	SymbolNumWeight      string                `yaml:"symbolNumWeight" json:"symbolNumWeight"`
	SymbolNumWeightVW    *sgc7game.ValWeights2 `yaml:"-" json:"-"`
	IgnoreSymbols        []string              `yaml:"ignoreSymbols" json:"ignoreSymbols"`
	IgnoreSymbolCodes    []int                 `yaml:"-" json:"-"`
}

// SetLinkComponent
func (cfg *AddSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type AddSymbols struct {
	*BasicComponent `json:"-"`
	Config          *AddSymbolsConfig `json:"config"`
}

// Init -
func (addSymbols *AddSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("AddSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &AddSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("AddSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return addSymbols.InitEx(cfg, pool)
}

// InitEx -
func (addSymbols *AddSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	addSymbols.Config = cfg.(*AddSymbolsConfig)
	addSymbols.Config.ComponentType = AddSymbolsTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[addSymbols.Config.Symbol]
	if !isok {
		goutils.Error("AddSymbols.InitEx:Symbol",
			slog.String("symbol", addSymbols.Config.Symbol),
			goutils.Err(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	addSymbols.Config.SymbolCode = sc

	if addSymbols.Config.SymbolNumWeight != "" {
		vw2, err := pool.LoadIntWeights(addSymbols.Config.SymbolNumWeight, addSymbols.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightReels.Init:LoadIntWeights",
				slog.String("SymbolNumWeight", addSymbols.Config.SymbolNumWeight),
				goutils.Err(err))

			return err
		}

		addSymbols.Config.SymbolNumWeightVW = vw2
	}

	for _, v := range addSymbols.Config.IgnoreSymbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("AddSymbols.InitEx:IgnoreSymbols",
				slog.String("symbol", v),
				goutils.Err(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		addSymbols.Config.IgnoreSymbolCodes = append(addSymbols.Config.IgnoreSymbolCodes, sc)
	}

	addSymbols.onInit(&addSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (addSymbols *AddSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceReelWithMask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	cd.UsedScenes = nil

	gs := addSymbols.GetTargetScene3(gameProp, curpr, prs, 0)

	num := addSymbols.Config.SymbolNum

	if addSymbols.Config.SymbolNumWeightVW != nil {
		cv, err := addSymbols.Config.SymbolNumWeightVW.RandVal(plugin)
		if err != nil {
			goutils.Error("AddSymbols.OnPlayGame:SymbolNumWeightVW",
				goutils.Err(err))

			return "", err
		}

		num = cv.Int()
	}

	if num <= 0 {
		nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	pos := make([]int, 0, gs.Width*gs.Height*2)
	for x, arr := range gs.Arr {
		for y, s := range arr {
			if goutils.IndexOfIntSlice(addSymbols.Config.IgnoreSymbolCodes, s, 0) < 0 {
				pos = append(pos, x, y)
			}
		}
	}

	if len(pos) <= 0 {
		nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	for i := 0; i < num; i++ {
		cr, err := plugin.Random(context.Background(), len(pos)/2)
		if err != nil {
			goutils.Error("AddSymbols.OnPlayGame:Random",
				goutils.Err(err))

			return "", err
		}

		ngs.Arr[pos[cr*2]][pos[cr*2+1]] = addSymbols.Config.SymbolCode

		pos = append(pos[:cr*2], pos[(cr+1)*2:]...)

		if len(pos) <= 0 {
			break
		}
	}

	addSymbols.AddScene(gameProp, curpr, ngs, cd)

	nc := addSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (addSymbols *AddSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after addSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (addSymbols *AddSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewAddSymbols(name string) IComponent {
	return &AddSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "symbol": "WL",
// "symbolNumType": "number",
// "symbolNum": 1,
// "ignoreSymbols": [
//
//	"WL",
//	"SC"
//
// ]
type jsonAddSymbols struct {
	Symbol          string   `json:"symbol"`
	SymbolNum       int      `json:"symbolNum"`
	SymbolNumWeight string   `json:"symbolNumWeight"`
	IgnoreSymbols   []string `json:"ignoreSymbols"`
}

func (jcfg *jsonAddSymbols) build() *AddSymbolsConfig {
	cfg := &AddSymbolsConfig{
		Symbol:          jcfg.Symbol,
		SymbolNum:       jcfg.SymbolNum,
		SymbolNumWeight: jcfg.SymbolNumWeight,
		IgnoreSymbols:   jcfg.IgnoreSymbols,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseAddSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseAddSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseAddSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonAddSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseAddSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: AddSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
