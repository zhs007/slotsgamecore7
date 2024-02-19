package lowcode

import (
	"fmt"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
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
	SymbolCodes []int
}

// OnNewGame -
func (rollSymbolData *RollSymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	rollSymbolData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (rollSymbolData *RollSymbolData) OnNewStep(gameProp *GameProperty, component IComponent) {
	rollSymbolData.BasicComponentData.OnNewStep(gameProp, component)
}

// BuildPBComponentData
func (rollSymbolData *RollSymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RollSymbolData{
		BasicComponentData: rollSymbolData.BuildPBBasicComponentData(),
		// SymbolCode:         int32(rollSymbolData.SymbolCode),
	}

	for _, v := range rollSymbolData.SymbolCodes {
		pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(v))
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
	SymbolNum              int                   `yaml:"symbolNum" json:"symbolNum"`
	Weight                 string                `yaml:"weight" json:"weight"`
	WeightVW               *sgc7game.ValWeights2 `json:"-"`
	SrcSymbolCollection    string                `yaml:"srcSymbolCollection" json:"srcSymbolCollection"`
	IgnoreSymbolCollection string                `yaml:"ignoreSymbolCollection" json:"ignoreSymbolCollection"`
	TargetSymbolCollection string                `yaml:"targetSymbolCollection" json:"targetSymbolCollection"`
}

// SetLinkComponent
func (cfg *RollSymbolConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
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

func (rollSymbol *RollSymbol) GetSymbolNum(basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVSymbolNum)
	if isok {
		return v
	}

	return rollSymbol.Config.SymbolNum
}

// playgame
func (rollSymbol *RollSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// rollSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	rsd := icd.(*RollSymbolData)

	rsd.SymbolCodes = nil

	sn := rollSymbol.GetSymbolNum(&rsd.BasicComponentData)
	for i := 0; i < sn; i++ {
		vw := rollSymbol.getValWeight(gameProp)
		if vw == nil {
			nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("RollSymbol.OnPlayGame:RandVal",
				zap.Error(err))

			return "", err
		}

		sc := cr.Int()

		rsd.SymbolCodes = append(rsd.SymbolCodes, sc)

		if rollSymbol.Config.TargetSymbolCollection != "" {
			gameProp.AddComponentSymbol(rollSymbol.Config.TargetSymbolCollection, sc)
		}
	}

	if len(rsd.SymbolCodes) == 0 {
		nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (rollSymbol *RollSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	rsd := icd.(*RollSymbolData)

	fmt.Printf("rollSymbol %v, got ", rollSymbol.GetName())

	for _, v := range rsd.SymbolCodes {
		fmt.Printf("%v ", gameProp.Pool.DefaultPaytables.GetStringFromInt(v))
	}

	fmt.Print("\n")

	return nil
}

// OnStats
func (rollSymbol *RollSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (rollSymbol *RollSymbol) NewComponentData() IComponentData {
	return &RollSymbolData{}
}

func NewRollSymbol(name string) IComponent {
	return &RollSymbol{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"weight": "fgbookofsymbol",
//		"symbolNum": 3,
//		"ignoreSymbolCollection": "fg-syms",
//		"targetSymbolCollection": "fg-syms"
//	},
type jsonRollSymbol struct {
	Weight                 string `json:"weight"`
	SymbolNum              int    `json:"symbolNum"`
	SrcSymbolCollection    string `json:"srcSymbolCollection"`
	IgnoreSymbolCollection string `json:"ignoreSymbolCollection"`
	TargetSymbolCollection string `json:"targetSymbolCollection"`
}

func (jcfg *jsonRollSymbol) build() *RollSymbolConfig {
	cfg := &RollSymbolConfig{
		Weight:                 jcfg.Weight,
		SymbolNum:              jcfg.SymbolNum,
		SrcSymbolCollection:    jcfg.SrcSymbolCollection,
		IgnoreSymbolCollection: jcfg.IgnoreSymbolCollection,
		TargetSymbolCollection: jcfg.TargetSymbolCollection,
	}

	cfg.UseSceneV3 = true

	return cfg
}

func parseRollSymbol(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRollSymbol:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRollSymbol:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonRollSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRollSymbol:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RollSymbolTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
