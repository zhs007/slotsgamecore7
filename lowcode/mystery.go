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

type MysteryData struct {
	BasicComponentData
	CurMysteryCode int
}

// OnNewGame -
func (mysteryData *MysteryData) OnNewGame() {
	mysteryData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (mysteryData *MysteryData) OnNewStep() {
	mysteryData.BasicComponentData.OnNewStep()

	mysteryData.CurMysteryCode = -1
}

// BuildPBComponentData
func (mysteryData *MysteryData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MysteryData{
		BasicComponentData: mysteryData.BuildPBBasicComponentData(),
		CurMysteryCode:     int32(mysteryData.CurMysteryCode),
	}

	return pbcd
}

// MysteryTriggerFeatureConfig - configuration for mystery trigger feature
type MysteryTriggerFeatureConfig struct {
	Symbol               string `yaml:"symbol" json:"symbol"`                             // like LIGHTNING
	RespinFirstComponent string `yaml:"respinFirstComponent" json:"respinFirstComponent"` // like lightning
}

// MysteryConfig - configuration for Mystery
type MysteryConfig struct {
	BasicComponentConfig   `yaml:",inline" json:",inline"`
	MysteryRNG             string                         `yaml:"mysteryRNG" json:"mysteryRNG"` // 强制用已经使用的随机数结果做 Mystery
	MysteryWeight          string                         `yaml:"mysteryWeight" json:"mysteryWeight"`
	Mystery                string                         `yaml:"mystery" json:"-"`
	Mysterys               []string                       `yaml:"mysterys" json:"mysterys"`
	MysteryTriggerFeatures []*MysteryTriggerFeatureConfig `yaml:"mysteryTriggerFeatures" json:"mysteryTriggerFeatures"`
}

type Mystery struct {
	*BasicComponent          `json:"-"`
	Config                   *MysteryConfig                       `json:"config"`
	MysteryWeights           *sgc7game.ValWeights2                `json:"-"`
	MysterySymbols           []int                                `json:"-"`
	MapMysteryTriggerFeature map[int]*MysteryTriggerFeatureConfig `json:"-"`
}

// maskOtherScene -
func (mystery *Mystery) maskOtherScene(gameProp *GameProperty, gs *sgc7game.GameScene, symbolCode int) *sgc7game.GameScene {
	// cgs := gs.Clone()
	cgs := gs.CloneEx(gameProp.PoolScene)

	for x, arr := range cgs.Arr {
		for y, v := range arr {
			if v != symbolCode {
				cgs.Arr[x][y] = -1
			} else {
				cgs.Arr[x][y] = 1
			}
		}
	}

	return cgs
}

// Init -
func (mystery *Mystery) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Mystery.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MysteryConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Mystery.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return mystery.InitEx(cfg, pool)
}

// InitEx -
func (mystery *Mystery) InitEx(cfg any, pool *GamePropertyPool) error {
	mystery.Config = cfg.(*MysteryConfig)

	if mystery.Config.MysteryWeight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(pool.Config.GetPath(mystery.Config.MysteryWeight, mystery.Config.UseFileMapping), "val", "weight", pool.DefaultPaytables)
		if err != nil {
			goutils.Error("Mystery.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", mystery.Config.MysteryWeight),
				zap.Error(err))

			return err
		}

		mystery.MysteryWeights = vw2
	}

	if len(mystery.Config.Mysterys) > 0 {
		for _, v := range mystery.Config.Mysterys {
			mystery.MysterySymbols = append(mystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[v])
		}
	} else {
		mystery.MysterySymbols = append(mystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[mystery.Config.Mystery])
	}

	for _, v := range mystery.Config.MysteryTriggerFeatures {
		symbolCode := pool.DefaultPaytables.MapSymbols[v.Symbol]

		mystery.MapMysteryTriggerFeature[symbolCode] = v
	}

	mystery.onInit(&mystery.Config.BasicComponentConfig)

	return nil
}

// playgame
func (mystery *Mystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	mystery.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[mystery.Name].(*MysteryData)

	gs := mystery.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")
	if !gs.HasSymbols(mystery.MysterySymbols) {
		mystery.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, &cd.BasicComponentData)
	} else {
		if mystery.MysteryWeights != nil {
			if mystery.Config.MysteryRNG != "" {
				rng := gameProp.GetTagInt(mystery.Config.MysteryRNG)
				cs := mystery.MysteryWeights.Vals[rng]

				curmcode := cs.Int()
				cd.CurMysteryCode = curmcode

				// gameProp.SetVal(GamePropCurMystery, curmcode)

				// sc2 := gs.Clone()
				sc2 := gs.CloneEx(gameProp.PoolScene)
				for _, v := range mystery.MysterySymbols {
					sc2.ReplaceSymbol(v, curmcode)
				}

				mystery.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)

				v, isok := mystery.MapMysteryTriggerFeature[curmcode]
				if isok {
					if v.RespinFirstComponent != "" {
						os := mystery.maskOtherScene(gameProp, sc2, curmcode)

						gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

						return nil
					}
				}
			} else {
				curm, err := mystery.MysteryWeights.RandVal(plugin)
				if err != nil {
					goutils.Error("Mystery.OnPlayGame:RandVal",
						zap.Error(err))

					return err
				}

				curmcode := curm.Int()

				// gameProp.SetVal(GamePropCurMystery, curm.Int())

				sc2 := gs.CloneEx(gameProp.PoolScene)
				// sc2 := gs.Clone()
				for _, v := range mystery.MysterySymbols {
					sc2.ReplaceSymbol(v, curm.Int())
				}

				mystery.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)

				v, isok := mystery.MapMysteryTriggerFeature[curmcode]
				if isok {
					if v.RespinFirstComponent != "" {
						os := mystery.maskOtherScene(gameProp, sc2, curmcode)

						gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

						return nil
					}
				}
			}
		}
	}

	mystery.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(mystery.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (mystery *Mystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[mystery.Name].(*MysteryData)

	if len(cd.UsedScenes) > 0 {
		if mystery.MysteryWeights != nil {
			fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(cd.CurMysteryCode))
			asciigame.OutputScene("after symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
		}
	}

	return nil
}

// OnStats
func (mystery *Mystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (mystery *Mystery) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.MysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Mystery.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return mystery.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (mystery *Mystery) NewComponentData() IComponentData {
	return &MysteryData{}
}

// EachUsedResults -
func (mystery *Mystery) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.MysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Mystery.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewMystery(name string) IComponent {
	mystery := &Mystery{
		BasicComponent:           NewBasicComponent(name),
		MapMysteryTriggerFeature: make(map[int]*MysteryTriggerFeatureConfig),
	}

	return mystery
}
