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

const WeightChgSymbolTypeName = "weightChgSymbol"

const (
	WCSCVMapChgWeight string = "mapChgWeight" // 可以修改配置项里的mapChgWeight，这里因为是个map，所以要当成 mapChgWeight:S 这样传递
)

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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &WeightChgSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightChgSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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
				zap.String("Weight", fn),
				zap.Error(err))

			return err
		}

		weightChgSymbol.Config.MapChgWeightVW[sc] = vw2
	}

	weightChgSymbol.onInit(&weightChgSymbol.Config.BasicComponentConfig)

	return nil
}

func (weightChgSymbol *WeightChgSymbol) getChgWeight(gameProp *GameProperty, basicCD *BasicComponentData, symbol int) *sgc7game.ValWeights2 {
	str := basicCD.GetConfigVal(WCSCVMapChgWeight + ":" + gameProp.Pool.Config.GetDefaultPaytables().GetStringFromInt(symbol))
	if str != "" {
		vw2, _ := gameProp.Pool.LoadSymbolWeights(str, "val", "weight", gameProp.Pool.DefaultPaytables, weightChgSymbol.Config.UseFileMapping)

		return vw2
	}

	return weightChgSymbol.Config.MapChgWeightVW[symbol]
}

// playgame
func (weightChgSymbol *WeightChgSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	weightChgSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[weightChgSymbol.Name].(*BasicComponentData)

	gs := weightChgSymbol.GetTargetScene3(gameProp, curpr, prs, cd, weightChgSymbol.Name, "", 0)

	cgs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range cgs.Arr {
		for y, s := range arr {
			vw := weightChgSymbol.getChgWeight(gameProp, cd, s)
			// vw, isok := weightChgSymbol.Config.MapChgWeightVW[s]
			if vw != nil {
				cr, err := vw.RandVal(plugin)
				if err != nil {
					goutils.Error("WeightChgSymbol.OnPlayGame:RandVal",
						zap.Error(err))

					return err
				}

				cgs.Arr[x][y] = cr.Int()
			}
		}
	}

	weightChgSymbol.AddScene(gameProp, curpr, cgs, cd)

	weightChgSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (weightChgSymbol *WeightChgSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[weightChgSymbol.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("After WeightChgSymbol", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (weightChgSymbol *WeightChgSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewWeightChgSymbol(name string) IComponent {
	return &WeightChgSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
