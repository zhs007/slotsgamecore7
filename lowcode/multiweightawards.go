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

const MultiWeightAwardsTypeName = "multiWeightAwards"

type MultiWeightAwardsData struct {
	BasicComponentData
	HasGot []bool
}

// OnNewGame -
func (multiWeightAwardsData *MultiWeightAwardsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	multiWeightAwardsData.BasicComponentData.OnNewGame(gameProp, component)

	multiWeightAwardsData.HasGot = nil
}

// // OnNewStep -
// func (multiWeightAwardsData *MultiWeightAwardsData) OnNewStep(gameProp *GameProperty, component IComponent) {
// 	multiWeightAwardsData.BasicComponentData.OnNewStep(gameProp, component)
// }

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
	InitMask             string                   `yaml:"initMask" json:"initMask"`                   // 用这个来初始化，true表示需要开奖
	ReverseInitMask      bool                     `yaml:"reverseInitMask" json:"reverseInitMask"`     // reverse the target mask
	TargetMask           string                   `yaml:"targetMask" json:"targetMask"`               // 用这个来初始化，true表示需要开奖
	ReverseTargetMask    bool                     `yaml:"reverseTargetMask" json:"reverseTargetMask"` // reverse the target mask
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MultiWeightAwardsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiWeightAwards.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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
				slog.String("Weight", v.Weight),
				goutils.Err(err))

			return err
		}

		v.VW = vw2

		for _, award := range v.Awards {
			award.Init()
		}
	}

	multiWeightAwards.onInit(&multiWeightAwards.Config.BasicComponentConfig)

	return nil
}

// InitEx -
func (multiWeightAwards *MultiWeightAwards) buildMask(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, cd *MultiWeightAwardsData) error {
	if multiWeightAwards.Config.ReverseTargetMask {
		mask := make([]bool, len(cd.HasGot))

		for i := 0; i < len(cd.HasGot); i++ {
			mask[i] = !cd.HasGot[i]
		}

		return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, multiWeightAwards.Config.TargetMask, mask, true)
	}

	return gameProp.Pool.SetMask(plugin, gameProp, curpr, gp, multiWeightAwards.Config.TargetMask, cd.HasGot, true)
}

// playgame
func (multiWeightAwards *MultiWeightAwards) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// multiWeightAwards.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	mwad := cd.(*MultiWeightAwardsData)

	mwad.HasGot = nil

	if multiWeightAwards.Config.InitMask != "" {
		mask, err := gameProp.Pool.GetMask(multiWeightAwards.Config.InitMask, gameProp)
		if err != nil {
			goutils.Error("MultiWeightAwards.OnPlayGame:GetMask",
				slog.String("mask", multiWeightAwards.Config.InitMask),
				goutils.Err(err))

			return "", err
		}

		for i, maskv := range mask {
			if maskv == !multiWeightAwards.Config.ReverseInitMask {
				v := multiWeightAwards.Config.Nodes[i]
				cv, err := v.VW.RandVal(plugin)
				if err != nil {
					goutils.Error("MultiWeightAwards.OnPlayGame:RandVal",
						goutils.Err(err))

					return "", err
				}

				if cv.Int() != 0 {
					if len(v.Awards) > 0 {
						gameProp.procAwards(plugin, v.Awards, curpr, gp)
					}

					mwad.HasGot = append(mwad.HasGot, true)
				} else {
					mwad.HasGot = append(mwad.HasGot, false)
				}
			} else {
				mwad.HasGot = append(mwad.HasGot, false)
			}
		}
	} else {
		for _, v := range multiWeightAwards.Config.Nodes {
			cv, err := v.VW.RandVal(plugin)
			if err != nil {
				goutils.Error("MultiWeightAwards.OnPlayGame:RandVal",
					goutils.Err(err))

				return "", err
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
	}

	if multiWeightAwards.Config.TargetMask != "" {
		err := multiWeightAwards.buildMask(plugin, gameProp, curpr, gp, mwad)
		if err != nil {
			goutils.Error("MultiWeightAwards.OnPlayGame:buildMask",
				slog.String("mask", multiWeightAwards.Config.TargetMask),
				goutils.Err(err))

			return "", err
		}
	}

	nc := multiWeightAwards.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (multiWeightAwards *MultiWeightAwards) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	mcd := cd.(*MultiWeightAwardsData)

	if len(mcd.HasGot) > 0 {
		fmt.Printf("MultiWeightAwards result is %v\n", mcd.HasGot)
	}

	return nil
}

// // OnStats
// func (multiWeightAwards *MultiWeightAwards) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (multiWeightAwards *MultiWeightAwards) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.MultiWeightAwardsData)
// 	if !isok {
// 		goutils.Error("MultiWeightAwards.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

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
// 			goutils.Err(err))

// 		return
// 	}

// 	for _, v := range pbcd.BasicComponentData.UsedResults {
// 		oneach(pr.Results[v])
// 	}
// }

func NewMultiWeightAwards(name string) IComponent {
	return &MultiWeightAwards{
		BasicComponent: NewBasicComponent(name, 0),
	}
}
