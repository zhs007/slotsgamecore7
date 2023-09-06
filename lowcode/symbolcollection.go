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
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

type SymbolCollectionData struct {
	BasicComponentData
	SymbolCodes []int
}

// OnNewGame -
func (symbolCollectionData *SymbolCollectionData) OnNewGame() {
	symbolCollectionData.BasicComponentData.OnNewGame()

	symbolCollectionData.SymbolCodes = nil
}

// OnNewStep -
func (symbolCollectionData *SymbolCollectionData) OnNewStep() {
	symbolCollectionData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (symbolCollectionData *SymbolCollectionData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SymbolCollectionData{
		BasicComponentData: symbolCollectionData.BuildPBBasicComponentData(),
	}

	for _, s := range symbolCollectionData.SymbolCodes {
		pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(s))
	}

	return pbcd
}

// SymbolCollectionConfig - configuration for SymbolCollection feature
type SymbolCollectionConfig struct {
	BasicComponentConfig `yaml:",inline"`
	WeightVal            string `yaml:"weightVal"`
}

// SymbolCollection - 也是一个非常特殊的组件，用来处理symbol随机集合的，譬如一共8个符号，按恒定权重先后选出3个符号
type SymbolCollection struct {
	*BasicComponent
	Config    *SymbolCollectionConfig
	WeightVal *sgc7game.ValWeights2
}

// Init -
func (symbolCollection *SymbolCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SymbolCollection.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &SymbolCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SymbolCollection.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return symbolCollection.InitEx(cfg, pool)
}

// InitEx -
func (symbolCollection *SymbolCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	symbolCollection.Config = cfg.(*SymbolCollectionConfig)

	if symbolCollection.Config.WeightVal != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(symbolCollection.Config.WeightVal, symbolCollection.Config.UseFileMapping), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("SymbolCollection.Init:LoadValWeights2FromExcel",
				zap.String("Weight", symbolCollection.Config.WeightVal),
				zap.Error(err))

			return err
		}

		symbolCollection.WeightVal = vw2
	}

	symbolCollection.onInit(&symbolCollection.Config.BasicComponentConfig)

	return nil
}

// Push -
func (symbolCollection *SymbolCollection) Push(plugin sgc7plugin.IPlugin, gameProp *GameProperty, gp *GameParams) error {
	cd := gameProp.MapComponentData[symbolCollection.Name].(*SymbolCollectionData)

	// 这样分开写，效率稍高一点点
	if len(cd.SymbolCodes) == 0 {
		cr, err := symbolCollection.WeightVal.RandVal(plugin)
		if err != nil {
			goutils.Error("SymbolCollection.Push:RandVal",
				zap.Error(err))

			return err
		}

		cd.SymbolCodes = append(cd.SymbolCodes, cr.Int())
	} else if len(cd.SymbolCodes) != len(symbolCollection.WeightVal.Vals) {
		vals := []sgc7game.IVal{}
		weights := []int{}

		for i, v := range symbolCollection.WeightVal.Vals {
			if goutils.IndexOfIntSlice(cd.SymbolCodes, v.Int(), 0) < 0 {
				vals = append(vals, v)
				weights = append(weights, symbolCollection.WeightVal.Weights[i])
			}
		}

		vw2, err := sgc7game.NewValWeights2(vals, weights)
		if err != nil {
			goutils.Error("SymbolCollection.Push:NewValWeights2",
				zap.Error(err))

			return err
		}

		cr, err := vw2.RandVal(plugin)
		if err != nil {
			goutils.Error("SymbolCollection.Push:RandVal",
				zap.Error(err))

			return err
		}

		cd.SymbolCodes = append(cd.SymbolCodes, cr.Int())
	}

	return nil
}

// OnNewGame -
func (symbolCollection *SymbolCollection) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[symbolCollection.Name]

	cd.OnNewGame()

	return nil
}

// playgame
func (symbolCollection *SymbolCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	symbolCollection.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	symbolCollection.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (symbolCollection *SymbolCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[symbolCollection.Name].(*SymbolCollectionData)

	if len(cd.SymbolCodes) > 0 {
		fmt.Printf("Symbols is %v\n", cd.SymbolCodes)
	}

	return nil
}

// OnStats
func (symbolCollection *SymbolCollection) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (symbolCollection *SymbolCollection) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.SymbolCollectionData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("SymbolCollection.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return symbolCollection.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (symbolCollection *SymbolCollection) NewComponentData() IComponentData {
	return &SymbolCollectionData{}
}

// EachUsedResults -
func (symbolCollection *SymbolCollection) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.SymbolCollectionData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("SymbolCollection.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewSymbolCollection(name string) IComponent {
	return &SymbolCollection{
		BasicComponent: NewBasicComponent(name),
	}
}
