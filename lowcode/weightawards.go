package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const WeightAwardsTypeName = "weightAwards"

type WeightAwardsData struct {
	BasicComponentData
	GotIndex []int
}

// OnNewGame -
func (weightAwardsData *WeightAwardsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	weightAwardsData.BasicComponentData.OnNewGame(gameProp, component)

	weightAwardsData.GotIndex = nil
}

// // OnNewStep -
// func (weightAwardsData *WeightAwardsData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	weightAwardsData.BasicComponentData.OnNewStep(gameProp, component)
// }

// Clone
func (weightAwardsData *WeightAwardsData) Clone() IComponentData {
	target := &WeightAwardsData{
		BasicComponentData: weightAwardsData.CloneBasicComponentData(),
	}

	target.GotIndex = make([]int, len(weightAwardsData.GotIndex))
	copy(target.GotIndex, weightAwardsData.GotIndex)

	return target
}

// BuildPBComponentData
func (weightAwardsData *WeightAwardsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.WeightAwardsData{
		BasicComponentData: weightAwardsData.BuildPBBasicComponentData(),
	}

	for _, v := range weightAwardsData.GotIndex {
		pbcd.GotIndex = append(pbcd.GotIndex, int32(v))
	}

	return pbcd
}

// WeightAwardsConfig - configuration for WeightAwards feature
type WeightAwardsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	AwardWeight          string                `yaml:"awardWeight" json:"awardWeight"`
	AwardWeightVW        *sgc7game.ValWeights2 `json:"-"`
	Awards               [][]*Award            `yaml:"awards" json:"awards"`                       // 新的奖励系统
	Nums                 int                   `yaml:"nums" json:"nums"`                           // how many arards are given
	TargetMask           string                `yaml:"targetMask" json:"targetMask"`               // output for the mask
	ReverseTargetMask    bool                  `yaml:"reverseTargetMask" json:"reverseTargetMask"` // reverse the target mask
}

type WeightAwards struct {
	*BasicComponent `json:"-"`
	Config          *WeightAwardsConfig `json:"config"`
}

// Init -
func (weightAwards *WeightAwards) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("WeightAwards.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &WeightAwardsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("WeightAwards.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return weightAwards.InitEx(cfg, pool)
}

// InitEx -
func (weightAwards *WeightAwards) InitEx(cfg any, pool *GamePropertyPool) error {
	weightAwards.Config = cfg.(*WeightAwardsConfig)
	weightAwards.Config.ComponentType = WeightAwardsTypeName

	for _, lst := range weightAwards.Config.Awards {
		for _, v := range lst {
			v.Init()
		}
	}

	if weightAwards.Config.AwardWeight != "" {
		vw2, err := pool.LoadIntWeights(weightAwards.Config.AwardWeight, weightAwards.Config.UseFileMapping)
		if err != nil {
			goutils.Error("WeightReels.Init:LoadIntWeights",
				slog.String("AwardWeight", weightAwards.Config.AwardWeight),
				goutils.Err(err))

			return err
		}

		weightAwards.Config.AwardWeightVW = vw2
	}

	weightAwards.onInit(&weightAwards.Config.BasicComponentConfig)

	return nil
}

// InitEx -
func (weightAwards *WeightAwards) buildMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd *WeightAwardsData) error {
	mask := make([]bool, len(weightAwards.Config.Awards))

	if weightAwards.Config.ReverseTargetMask {
		for i := range mask {
			mask[i] = true
		}

		for i := 0; i < len(cd.GotIndex); i++ {
			mask[cd.GotIndex[i]] = false
		}
	} else {
		for i := 0; i < len(cd.GotIndex); i++ {
			mask[cd.GotIndex[i]] = true
		}
	}

	return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, weightAwards.Config.TargetMask, mask, true)
}

// playgame
func (weightAwards *WeightAwards) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// weightAwards.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	mwad := icd.(*WeightAwardsData)

	mwad.GotIndex = nil

	if weightAwards.Config.Nums > 1 {
		vw := weightAwards.Config.AwardWeightVW.Clone()
		for i := 0; i < weightAwards.Config.Nums; i++ {
			cv, _, err := vw.RandValEx(plugin)
			if err != nil {
				goutils.Error("WeightAwards.OnPlayGame:RandValEx",
					goutils.Err(err))

				return "", err
			}

			if cv.Int() >= len(weightAwards.Config.Awards) {
				goutils.Error("WeightAwards.OnPlayGame",
					slog.Int("val", cv.Int()),
					goutils.Err(ErrInvalidWeightVal))

				return "", ErrInvalidWeightVal
			}

			gameProp.procAwards(plugin, weightAwards.Config.Awards[cv.Int()], curpr, gp)

			vw.RemoveVal(cv)

			mwad.GotIndex = append(mwad.GotIndex, cv.Int())
		}
	} else if weightAwards.Config.Nums == 1 {
		vw := weightAwards.Config.AwardWeightVW

		cv, i, err := vw.RandValEx(plugin)
		if err != nil {
			goutils.Error("WeightAwards.OnPlayGame:RandValEx",
				goutils.Err(err))

			return "", err
		}

		if cv.Int() >= len(weightAwards.Config.Awards) {
			goutils.Error("WeightAwards.OnPlayGame",
				slog.Int("val", cv.Int()),
				goutils.Err(ErrInvalidWeightVal))

			return "", ErrInvalidWeightVal
		}

		gameProp.procAwards(plugin, weightAwards.Config.Awards[i], curpr, gp)

		mwad.GotIndex = append(mwad.GotIndex, i)
	}

	if weightAwards.Config.TargetMask != "" {
		err := weightAwards.buildMask(plugin, gameProp, curpr, gp, mwad)
		if err != nil {
			goutils.Error("WeightAwards.OnPlayGame:buildMask",
				goutils.Err(err))

			return "", err
		}
	}

	nc := weightAwards.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (weightAwards *WeightAwards) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*WeightAwardsData)

	if len(cd.GotIndex) > 0 {
		fmt.Printf("WeightAwards result is %v\n", cd.GotIndex)
	}

	return nil
}

// // OnStats
// func (weightAwards *WeightAwards) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (multiWeightAwards *MultiWeightAwards) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.MultiWeightAwardsData)
// 	if !isok {
// 		goutils.Error("MultiWeightAwards.OnStatsWithPB",
// 			goutils.Err(ErrInvalidProto))

// 		return 0, ErrInvalidProto
// 	}

// 	return multiWeightAwards.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
// }

// NewComponentData -
func (weightAwards *WeightAwards) NewComponentData() IComponentData {
	return &WeightAwardsData{}
}

// // EachUsedResults -
// func (multiWeightAwards *MultiWeightAwards) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
// 	pbcd := &sgc7pb.SymbolCollectionData{}

// 	err := pbComponentData.UnmarshalTo(pbcd)
// 	if err != nil {
// 		goutils.Error("MultiWeightAwards.EachUsedResults:UnmarshalTo",
// 			goutils.Err(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

func NewWeightAwards(name string) IComponent {
	return &WeightAwards{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
