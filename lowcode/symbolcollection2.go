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

const SymbolCollection2TypeName = "symbolCollection2"

type SymbolCollection2Data struct {
	BasicComponentData
	SymbolCodes []int
}

// OnNewGame -
func (symbolCollection2Data *SymbolCollection2Data) OnNewGame() {
	symbolCollection2Data.BasicComponentData.OnNewGame()

	symbolCollection2Data.SymbolCodes = nil
}

// OnNewStep -
func (symbolCollection2Data *SymbolCollection2Data) OnNewStep() {
	symbolCollection2Data.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (symbolCollection2Data *SymbolCollection2Data) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolCollection2Data{
		BasicComponentData: symbolCollection2Data.BuildPBBasicComponentData(),
	}

	for _, s := range symbolCollection2Data.SymbolCodes {
		pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(s))
	}

	return pbcd
}

// SymbolCollection2Config - configuration for SymbolCollection2 feature
type SymbolCollection2Config struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MaxSymbolNum         int      `yaml:"maxSymbolNum" json:"maxSymbolNum"` // 0表示不限制
	InitSymbols          []string `yaml:"initSymbols" json:"initSymbols"`   // 初始化symbols
	InitSymbolCodes      []int    `yaml:"-" json:"-"`                       // 初始化symbols
}

// SetLinkComponent
func (cfg *SymbolCollection2Config) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// SymbolCollection2 - 也是一个非常特殊的组件，symbol集合
type SymbolCollection2 struct {
	*BasicComponent `json:"-"`
	Config          *SymbolCollection2Config `json:"config"`
}

// Init -
func (symbolCollection2 *SymbolCollection2) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolCollection2.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolCollection2Config{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolCollection2.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolCollection2.InitEx(cfg, pool)
}

// InitEx -
func (symbolCollection2 *SymbolCollection2) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolCollection2.Config = cfg.(*SymbolCollection2Config)
	symbolCollection2.Config.ComponentType = SymbolCollection2TypeName

	for _, v := range symbolCollection2.Config.InitSymbols {
		symbolCollection2.Config.InitSymbolCodes = append(symbolCollection2.Config.InitSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	symbolCollection2.onInit(&symbolCollection2.Config.BasicComponentConfig)

	return nil
}

// // Push -
// func (symbolCollection2 *SymbolCollection2) Push(plugin sgc7plugin.IPlugin, gameProp *GameProperty, gp *GameParams) error {
// 	cd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

// 	// 这样分开写，效率稍高一点点
// 	if len(cd.SymbolCodes) == 0 {
// 		cr, err := symbolCollection.WeightVal.RandVal(plugin)
// 		if err != nil {
// 			goutils.Error("SymbolCollection2.Push:RandVal",
// 				zap.Error(err))

// 			return err
// 		}

// 		cd.SymbolCodes = append(cd.SymbolCodes, cr.Int())
// 	} else if len(cd.SymbolCodes) != len(symbolCollection.WeightVal.Vals) {
// 		vals := []sgc7game.IVal{}
// 		weights := []int{}

// 		for i, v := range symbolCollection.WeightVal.Vals {
// 			if goutils.IndexOfIntSlice(cd.SymbolCodes, v.Int(), 0) < 0 {
// 				vals = append(vals, v)
// 				weights = append(weights, symbolCollection.WeightVal.Weights[i])
// 			}
// 		}

// 		vw2, err := sgc7game.NewValWeights2(vals, weights)
// 		if err != nil {
// 			goutils.Error("SymbolCollection2.Push:NewValWeights2",
// 				zap.Error(err))

// 			return err
// 		}

// 		cr, err := vw2.RandVal(plugin)
// 		if err != nil {
// 			goutils.Error("SymbolCollection2.Push:RandVal",
// 				zap.Error(err))

// 			return err
// 		}

// 		cd.SymbolCodes = append(cd.SymbolCodes, cr.Int())
// 	}

// 	return nil
// }

// // OnNewGame -
// func (symbolCollection *SymbolCollection) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[symbolCollection.Name]

// 	cd.OnNewGame()

// 	return nil
// }

// playgame
func (symbolCollection2 *SymbolCollection2) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolCollection2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	symbolCollection2.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolCollection2 *SymbolCollection2) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

	if len(cd.SymbolCodes) > 0 {
		fmt.Printf("Symbols is %v\n", cd.SymbolCodes)
	}

	return nil
}

// OnStats
func (symbolCollection2 *SymbolCollection2) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// NewComponentData -
func (symbolCollection2 *SymbolCollection2) NewComponentData() IComponentData {
	return &SymbolCollection2Data{}
}

// // EachUsedResults -
// func (symbolCollection2 *SymbolCollection2) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
// 	pbcd := &sgc7pb.SymbolCollectionData{}

// 	err := pbComponentData.UnmarshalTo(pbcd)
// 	if err != nil {
// 		goutils.Error("SymbolCollection.EachUsedResults:UnmarshalTo",
// 			zap.Error(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

// GetSymbols -
func (symbolCollection2 *SymbolCollection2) GetSymbols(gameProp *GameProperty) []int {
	scd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

	return scd.SymbolCodes
}

// AddSymbol -
func (symbolCollection2 *SymbolCollection2) AddSymbol(gameProp *GameProperty, symbolCode int) {
	scd := gameProp.MapComponentData[symbolCollection2.Name].(*SymbolCollection2Data)

	if symbolCollection2.Config.MaxSymbolNum <= 0 || len(scd.SymbolCodes) < symbolCollection2.Config.MaxSymbolNum {
		scd.SymbolCodes = append(scd.SymbolCodes, symbolCode)
	}
}

func NewSymbolCollection2(name string) IComponent {
	return &SymbolCollection2{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"isWinBreak": "false"
//	},
type jsonSymbolCollection2 struct {
}

func (jr *jsonSymbolCollection2) build() *SymbolCollection2Config {
	cfg := &SymbolCollection2Config{}

	cfg.UseSceneV3 = true

	return cfg
}

func parseSymbolCollection2(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSymbolCollection2:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolCollection2:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonSymbolCollection2{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSymbolCollection2:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SymbolCollection2TypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
