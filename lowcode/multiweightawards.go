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

const MultiWeightAwardsTypeName = "multiWeightAwards"

type MultiWeightAwardsData struct {
	BasicComponentData
	HasGot []bool
}

// OnNewGame -
func (multiWeightAwardsData *MultiWeightAwardsData) OnNewGame() {
	multiWeightAwardsData.BasicComponentData.OnNewGame()

	multiWeightAwardsData.HasGot = nil
}

// OnNewStep -
func (multiWeightAwardsData *MultiWeightAwardsData) OnNewStep() {
	multiWeightAwardsData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (multiWeightAwardsData *MultiWeightAwardsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MultiWeightAwardsData{
		BasicComponentData: multiWeightAwardsData.BuildPBBasicComponentData(),
	}

	pbcd.HasGot = append(pbcd.HasGot, multiWeightAwardsData.HasGot...)

	// for _, s := range multiWeightAwardsData.HasGot {
	// 	pbcd.HasGot = append(pbcd.HasGot, s)
	// }

	return pbcd
}

type MultiWeightAwardsNode struct {
	Awards []*Award              `yaml:"awards" json:"awards"` // 新的奖励系统
	Weight string                `yaml:"weight" json:"weight"` //
	VW     *sgc7game.ValWeights2 `yaml:"-" json:"-"`           //
}

// MultiWeightAwardsConfig - configuration for MultiWeightAwards feature
type MultiWeightAwardsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Nodes                []*MultiWeightAwardsNode `yaml:"nodes" json:"nodes"`
}

type MultiWeightAwards struct {
	*BasicComponent `json:"-"`
	Config          *MultiWeightAwardsConfig `json:"config"`
}

// Init -
func (multiWeightAwards *MultiWeightAwards) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiWeightAwards.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiWeightAwardsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiWeightAwards.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return multiWeightAwards.InitEx(cfg, pool)
}

// InitEx -
func (multiWeightAwards *MultiWeightAwards) InitEx(cfg any, pool *GamePropertyPool) error {
	multiWeightAwards.Config = cfg.(*MultiWeightAwardsConfig)
	multiWeightAwards.Config.ComponentType = MultiWeightAwardsTypeName

	for _, v := range multiWeightAwards.Config.Nodes {
		vw2, err := pool.LoadIntWeights(v.Weight, multiWeightAwards.Config.UseFileMapping)
		if err != nil {
			goutils.Error("MultiWeightAwards.Init:LoadIntWeights",
				zap.String("Weight", v.Weight),
				zap.Error(err))

			return err
		}

		v.VW = vw2
	}

	multiWeightAwards.onInit(&multiWeightAwards.Config.BasicComponentConfig)

	return nil
}

// playgame
func (multiWeightAwards *MultiWeightAwards) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	multiWeightAwards.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	mwad := gameProp.MapComponentData[multiWeightAwards.Name].(*MultiWeightAwardsData)

	mwad.HasGot = nil

	for _, v := range multiWeightAwards.Config.Nodes {
		cv, err := v.VW.RandVal(plugin)
		if err != nil {
			goutils.Error("MultiWeightAwards.OnPlayGame:RandVal",
				zap.Error(err))

			return err
		}

		if cv.Int() != 0 {
			if len(v.Awards) > 0 {
				gameProp.procAwards(plugin, v.Awards, curpr, gp)
			}

			mwad.HasGot = append(mwad.HasGot, true)
		} else {
			mwad.HasGot = append(mwad.HasGot, false)
		}
	}

	multiWeightAwards.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiWeightAwards *MultiWeightAwards) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[multiWeightAwards.Name].(*MultiWeightAwardsData)

	if len(cd.HasGot) > 0 {
		fmt.Printf("MultiWeightAwards result is %v\n", cd.HasGot)
	}

	return nil
}

// OnStats
func (multiWeightAwards *MultiWeightAwards) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// // OnStatsWithPB -
// func (multiWeightAwards *MultiWeightAwards) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.MultiWeightAwardsData)
// 	if !isok {
// 		goutils.Error("MultiWeightAwards.OnStatsWithPB",
// 			zap.Error(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	return multiWeightAwards.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// NewComponentData -
func (multiWeightAwards *MultiWeightAwards) NewComponentData() IComponentData {
	return &MultiWeightAwardsData{}
}

// // EachUsedResults -
// func (multiWeightAwards *MultiWeightAwards) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
// 	pbcd := &sgc7pb.SymbolCollectionData{}

// 	err := pbComponentData.UnmarshalTo(pbcd)
// 	if err != nil {
// 		goutils.Error("MultiWeightAwards.EachUsedResults:UnmarshalTo",
// 			zap.Error(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

func NewMultiWeightAwards(name string) IComponent {
	return &MultiWeightAwards{
		BasicComponent: NewBasicComponent(name),
	}
}