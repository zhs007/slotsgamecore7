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

// 已弃用，待清理

const WeightChgSymbolTypeName = "weightChgSymbol"

// WeightChgSymbolConfig - configuration for WeightChgSymbol feature
type WeightChgSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MapChgWeight         map[string]string             `yaml:"mapChgWeight" json:"mapChgWeight"`
	MapChgWeightVW       map[int]*sgc7game.ValWeights2 `yaml:"-" json:"-"`
}

type WeightChgSymbol struct {
	*BasicComponent `json:"-"`
	Config          *WeightChgSymbolConfig `json:"config"`
}

// Init -
func (weightChgSymbol *WeightChgSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightChgSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WeightChgSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightChgSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return weightChgSymbol.InitEx(cfg, pool)
}

// InitEx -
func (weightChgSymbol *WeightChgSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	weightChgSymbol.Config = cfg.(*WeightChgSymbolConfig)
	weightChgSymbol.Config.ComponentType = WeightChgSymbolTypeName

	weightChgSymbol.Config.MapChgWeightVW = make(map[int]*sgc7game.ValWeights2)

	for s, fn := range weightChgSymbol.Config.MapChgWeight {
		sc := pool.DefaultPaytables.MapSymbols[s]

		vw2, err := pool.LoadSymbolWeights(fn, "val", "weight", pool.DefaultPaytables, weightChgSymbol.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightChgSymbol.InitEx:LoadIntWeights",
				slog.String("Weight", fn),
				goutils.Err(err))

			return err
		}

		weightChgSymbol.Config.MapChgWeightVW[sc] = vw2
	}

	weightChgSymbol.onInit(&weightChgSymbol.Config.BasicComponentConfig)

	return nil
}

func (weightChgSymbol *WeightChgSymbol) getChgWeight(gameProp *GameProperty, basicCD *BasicComponentData, symbol int) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(CCVMapChgWeight + ":" + gameProp.Pool.Config.GetDefaultPaytables().GetStringFromInt(symbol))
	if str != "" {
		vw2, _ := gameProp.Pool.LoadSymbolWeights(str, "val", "weight", gameProp.Pool.DefaultPaytables, weightChgSymbol.Config.UseFileMapping)

		return vw2
	}

	return weightChgSymbol.Config.MapChgWeightVW[symbol]
}

// playgame
func (weightChgSymbol *WeightChgSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*BasicComponentData)

	gs := weightChgSymbol.GetTargetScene3(gameProp, curpr, prs, 0)

	cgs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range cgs.Arr {
		for y, s := range arr {
			vw := weightChgSymbol.getChgWeight(gameProp, cd, s)
			if vw != nil {
				cr, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("WeightChgSymbol.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				cgs.Arr[x][y] = cr.Int()
			}
		}
	}

	weightChgSymbol.AddScene(gameProp, curpr, cgs, cd)

	nc := weightChgSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightChgSymbol *WeightChgSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("After WeightChgSymbol", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

func NewWeightChgSymbol(name string) IComponent {
	return &WeightChgSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
