package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const RollSymbolTypeName = "rollSymbol"

const (
	RSDVSymbol string = "symbol" // roll a symbol
)

type RollSymbolData struct {
	BasicComponentData
	SymbolCode int
}

// OnNewGame -
func (rollSymbolData *RollSymbolData) OnNewGame() {
	rollSymbolData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (rollSymbolData *RollSymbolData) OnNewStep() {
	rollSymbolData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (rollSymbolData *RollSymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RollSymbolData{
		BasicComponentData: rollSymbolData.BuildPBBasicComponentData(),
		SymbolCode:         int32(rollSymbolData.SymbolCode),
	}

	return pbcd
}

// GetVal -
func (rollSymbolData *RollSymbolData) GetVal(key string) int {
	return 0
}

// SetVal -
func (rollSymbolData *RollSymbolData) SetVal(key string, val int) {
}

// RollSymbolConfig - configuration for RollSymbol
type RollSymbolConfig struct {
	BasicComponentConfig   `yaml:",inline" json:",inline"`
	Weight                 string                `yaml:"weight" json:"weight"`
	WeightVW               *sgc7game.ValWeights2 `json:"-"`
	SrcSymbolCollection    string                `yaml:"srcSymbolCollection" json:"srcSymbolCollection"`
	IgnoreSymbolCollection string                `yaml:"ignoreSymbolCollection" json:"ignoreSymbolCollection"`
	TargetSymbolCollection string                `yaml:"targetSymbolCollection" json:"targetSymbolCollection"`
}

type RollSymbol struct {
	*BasicComponent `json:"-"`
	Config          *RollSymbolConfig `json:"config"`
}

// Init -
func (rollSymbol *RollSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RollSymbol.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &RollSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightBranch.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return rollSymbol.InitEx(cfg, pool)
}

// InitEx -
func (rollSymbol *RollSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	rollSymbol.Config = cfg.(*RollSymbolConfig)
	rollSymbol.Config.ComponentType = RollSymbolTypeName

	if rollSymbol.Config.Weight != "" {
		vw2, err := pool.LoadSymbolWeights(rollSymbol.Config.Weight, "val", "weight", pool.DefaultPaytables, rollSymbol.Config.UseFileMapping)
		if err != nil {
			goutils.Error("RollSymbol.Init:LoadStrWeights",
				zap.String("Weight", rollSymbol.Config.Weight),
				zap.Error(err))

			return err
		}

		rollSymbol.Config.WeightVW = vw2
	} else {
		goutils.Error("RollSymbol.InitEx:Weight",
			zap.Error(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	rollSymbol.onInit(&rollSymbol.Config.BasicComponentConfig)

	return nil
}

func (rollSymbol *RollSymbol) getValWeight(gameProp *GameProperty) *sgc7game.ValWeights2 {
	if rollSymbol.Config.SrcSymbolCollection == "" && rollSymbol.Config.IgnoreSymbolCollection == "" {
		return rollSymbol.Config.WeightVW
	}

	var vw *sgc7game.ValWeights2

	if rollSymbol.Config.SrcSymbolCollection != "" {
		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.SrcSymbolCollection)

		vw = rollSymbol.Config.WeightVW.CloneWithIntArray(symbols)
	}

	if vw == nil {
		vw = rollSymbol.Config.WeightVW.Clone()
	}

	if rollSymbol.Config.IgnoreSymbolCollection != "" {
		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.IgnoreSymbolCollection)

		if len(symbols) > 0 {
			vw = vw.CloneWithoutIntArray(symbols)
		}
	}

	if len(vw.Vals) == 0 {
		return nil
	}

	return vw
}

// playgame
func (rollSymbol *RollSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	rollSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	rsd := gameProp.MapComponentData[rollSymbol.Name].(*RollSymbolData)

	vw := rollSymbol.getValWeight(gameProp)
	if vw == nil {
		rollSymbol.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	cr, err := vw.RandVal(plugin)
	if err != nil {
		goutils.Error("RollSymbol.OnPlayGame:RandVal",
			zap.Error(err))

		return err
	}

	rsd.SymbolCode = cr.Int()

	if rollSymbol.Config.TargetSymbolCollection != "" {
		gameProp.AddComponentSymbol(rollSymbol.Config.TargetSymbolCollection, rsd.SymbolCode)
	}

	rollSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (rollSymbol *RollSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	rsd := gameProp.MapComponentData[rollSymbol.Name].(*RollSymbolData)

	fmt.Printf("rollSymbol %v, got %v \n", rollSymbol.GetName(), gameProp.Pool.DefaultPaytables.GetStringFromInt(rsd.SymbolCode))

	return nil
}

// OnStats
func (rollSymbol *RollSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewRollSymbol(name string) IComponent {
	return &RollSymbol{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
