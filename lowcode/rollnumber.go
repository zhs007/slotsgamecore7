package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const RollNumberTypeName = "rollNumber"

type RollNumberData struct {
	BasicComponentData
	Number int
}

// OnNewGame -
func (rollNumberData *RollNumberData) OnNewGame(gameProp *GameProperty, component IComponent) {
	rollNumberData.BasicComponentData.OnNewGame(gameProp, component)
}

// // OnNewStep -
// func (rollSymbolData *RollSymbolData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	rollSymbolData.BasicComponentData.OnNewStep(gameProp, component)
// }

// Clone
func (rollNumberData *RollNumberData) Clone() IComponentData {
	target := &RollNumberData{
		BasicComponentData: rollNumberData.CloneBasicComponentData(),
		Number:             rollNumberData.Number,
	}

	// target.SymbolCodes = make([]int, len(rollSymbolData.SymbolCodes))
	// copy(target.SymbolCodes, rollSymbolData.SymbolCodes)

	return target
}

// BuildPBComponentData
func (rollNumberData *RollNumberData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RollNumberData{
		BasicComponentData: rollNumberData.BuildPBBasicComponentData(),
		Number:             int32(rollNumberData.Number),
		// SymbolCode:         int32(rollSymbolData.SymbolCode),
	}

	// for _, v := range rollSymbolData.SymbolCodes {
	// 	pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(v))
	// }

	return pbcd
}

// GetVal -
func (rollNumberData *RollNumberData) GetVal(key string) (int, bool) {
	if key == CVNumber || key == CVOutputInt {
		return rollNumberData.Number, true
	}

	return 0, false
}

// // SetVal -
// func (rollNumberData *RollNumberData) SetVal(key string, val int) {
// 	if key == CVNumber {
// 		rollNumberData.Number = val
// 	}
// }

// RollNumberConfig - configuration for RollNumber
type RollNumberConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Weight               string                `yaml:"weight" json:"weight"`
	WeightVW             *sgc7game.ValWeights2 `json:"-"`
	Awards               []*Award              `yaml:"awards" json:"awards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *RollNumberConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type RollNumber struct {
	*BasicComponent `json:"-"`
	Config          *RollNumberConfig `json:"config"`
}

// Init -
func (rollNumber *RollNumber) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RollNumber.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RollNumberConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RollNumber.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return rollNumber.InitEx(cfg, pool)
}

// InitEx -
func (rollNumber *RollNumber) InitEx(cfg any, pool *GamePropertyPool) error {
	rollNumber.Config = cfg.(*RollNumberConfig)
	rollNumber.Config.ComponentType = RollNumberTypeName

	if rollNumber.Config.Weight != "" {
		vw2, err := pool.LoadSymbolWeights(rollNumber.Config.Weight, "val", "weight", pool.DefaultPaytables, rollNumber.Config.UseFileMapping)
		if err != nil {
			goutils.Error("RollNumber.Init:LoadStrWeights",
				slog.String("Weight", rollNumber.Config.Weight),
				goutils.Err(err))

			return err
		}

		rollNumber.Config.WeightVW = vw2
	} else {
		goutils.Error("RollNumber.InitEx:Weight",
			goutils.Err(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	for _, award := range rollNumber.Config.Awards {
		award.Init()
	}

	rollNumber.onInit(&rollNumber.Config.BasicComponentConfig)

	return nil
}

// func (rollSymbol *RollNumber) getValWeight(gameProp *GameProperty) *sgc7game.ValWeights2 {
// 	if rollSymbol.Config.SrcSymbolCollection == "" && rollSymbol.Config.IgnoreSymbolCollection == "" {
// 		return rollSymbol.Config.WeightVW
// 	}

// 	var vw *sgc7game.ValWeights2

// 	if rollSymbol.Config.SrcSymbolCollection != "" {
// 		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.SrcSymbolCollection)

// 		vw = rollSymbol.Config.WeightVW.CloneWithIntArray(symbols)
// 	}

// 	if vw == nil {
// 		vw = rollSymbol.Config.WeightVW.Clone()
// 	}

// 	if rollSymbol.Config.IgnoreSymbolCollection != "" {
// 		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.IgnoreSymbolCollection)

// 		if len(symbols) > 0 {
// 			vw = vw.CloneWithoutIntArray(symbols)
// 		}
// 	}

// 	if len(vw.Vals) == 0 {
// 		return nil
// 	}

// 	return vw
// }

// func (rollSymbol *RollNumber) getSymbolNum(gameProp *GameProperty, basicCD *BasicComponentData) int {
// 	v, isok := basicCD.GetConfigIntVal(CCVSymbolNum)
// 	if isok {
// 		return v
// 	}

// 	if rollSymbol.Config.SymbolNumComponent != "" {
// 		cd := gameProp.GetComponentDataWithName(rollSymbol.Config.SymbolNumComponent)
// 		if cd != nil {
// 			return cd.GetOutput()
// 		}
// 	}

// 	return rollSymbol.Config.SymbolNum
// }

// playgame
func (rollSymbol *RollNumber) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// rollSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	rnd := icd.(*RollNumberData)

	rnd.Number = 0

	cr, err := rollSymbol.Config.WeightVW.RandVal(plugin)
	if err != nil {
		goutils.Error("RollNumber.OnPlayGame:RandVal",
			goutils.Err(err))

		return "", err
	}

	rnd.Number = cr.Int()

	// sn := rollSymbol.getSymbolNum(gameProp, &rsd.BasicComponentData)
	// for i := 0; i < sn; i++ {
	// 	vw := rollSymbol.getValWeight(gameProp)
	// 	if vw == nil {
	// 		break
	// 	}

	// 	cr, err := vw.RandVal(plugin)
	// 	if err != nil {
	// 		goutils.Error("RollSymbol.OnPlayGame:RandVal",
	// 			goutils.Err(err))

	// 		return "", err
	// 	}

	// 	sc := cr.Int()

	// 	rsd.SymbolCodes = append(rsd.SymbolCodes, sc)

	// 	if rollSymbol.Config.TargetSymbolCollection != "" {
	// 		gameProp.AddComponentSymbol(rollSymbol.Config.TargetSymbolCollection, sc)
	// 	}
	// }

	// if len(rsd.SymbolCodes) == 0 {
	// 	nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

	// 	return nc, ErrComponentDoNothing
	// }

	if len(rollSymbol.Config.Awards) > 0 {
		gameProp.procAwards(plugin, rollSymbol.Config.Awards, curpr, gp)
	}

	nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (rollSymbol *RollNumber) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	rsd := icd.(*RollSymbolData)

	fmt.Printf("rollSymbol %v, got ", rollSymbol.GetName())

	for _, v := range rsd.SymbolCodes {
		fmt.Printf("%v ", gameProp.Pool.DefaultPaytables.GetStringFromInt(v))
	}

	fmt.Print("\n")

	return nil
}

// // OnStats
// func (rollSymbol *RollSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// NewComponentData -
func (rollSymbol *RollNumber) NewComponentData() IComponentData {
	return &RollNumberData{}
}

func NewRollNumber(name string) IComponent {
	return &RollNumber{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"weight": "fgbookofsymbol",
//		"symbolNum": 3,
//	    "symbolNumComponent": "bg-symnum",
//		"ignoreSymbolCollection": "fg-syms",
//		"targetSymbolCollection": "fg-syms"
//	},
type jsonRollNumber struct {
	Weight string `json:"weight"`
}

func (jcfg *jsonRollNumber) build() *RollNumberConfig {
	cfg := &RollNumberConfig{
		Weight: jcfg.Weight,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseRollNumber(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRollNumber:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRollNumber:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRollNumber{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRollNumber:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseRollNumber:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RollNumberTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
