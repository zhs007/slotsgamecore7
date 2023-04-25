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

type ReelSetMysteryData struct {
	BasicComponentData
	CurMysteryCode int
}

// OnNewGame -
func (reelSetMysteryData *ReelSetMysteryData) OnNewGame() {
	reelSetMysteryData.BasicComponentData.OnNewGame()
}

// OnNewGame -
func (reelSetMysteryData *ReelSetMysteryData) OnNewStep() {
	reelSetMysteryData.BasicComponentData.OnNewStep()

	reelSetMysteryData.CurMysteryCode = -1
}

// BuildPBComponentData
func (reelSetMysteryData *ReelSetMysteryData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.ReelSetMysteryData{
		BasicComponentData: reelSetMysteryData.BuildPBBasicComponentData(),
		CurMysteryCode:     int32(reelSetMysteryData.CurMysteryCode),
	}

	return pbcd
}

// ReelSetMysteryConfig - configuration for ReelSetMystery
type ReelSetMysteryConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MysteryRNG           string            `yaml:"mysteryRNG"` // 强制用已经使用的随机数结果做 ReelSetMystery
	MysterySymbols       []string          `yaml:"mysterySymbols"`
	MapMysteryWeight     map[string]string `yaml:"mapMysteryWeight"`
}

type ReelSetMystery struct {
	*BasicComponent
	Config             *ReelSetMysteryConfig
	MapMysteryWeights  map[string]*sgc7game.ValWeights2
	MysterySymbolCodes []int
}

// Init -
func (reelSetMystery *ReelSetMystery) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReelSetMystery.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReelSetMysteryConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReelSetMystery.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	reelSetMystery.Config = cfg

	for k, v := range reelSetMystery.Config.MapMysteryWeight {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(pool.Config.GetPath(v), "val", "weight", pool.DefaultPaytables)
		if err != nil {
			goutils.Error("ReelSetMystery.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", v),
				zap.Error(err))

			return err
		}

		reelSetMystery.MapMysteryWeights[k] = vw2
	}

	for _, v := range reelSetMystery.Config.MysterySymbols {
		reelSetMystery.MysterySymbolCodes = append(reelSetMystery.MysterySymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	reelSetMystery.onInit(&cfg.BasicComponentConfig)

	return nil
}

func (reelSetMystery *ReelSetMystery) hasMystery(gs *sgc7game.GameScene) bool {
	for _, v := range reelSetMystery.MysterySymbolCodes {
		if gs.HasSymbol(v) {
			return true
		}
	}

	return false
}

// playgame
func (reelSetMystery *ReelSetMystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[reelSetMystery.Name].(*MysteryData)

	gs := reelSetMystery.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")
	if !reelSetMystery.hasMystery(gs) {
		reelSetMystery.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, &cd.BasicComponentData)
	} else {
		vw2, isok := reelSetMystery.MapMysteryWeights[gameProp.GetTagStr(TagCurReels)]
		if !isok {
			goutils.Error("ReelSetMystery.OnPlayGame:MapMysteryWeights",
				zap.String("TagCurReels", gameProp.GetTagStr(TagCurReels)),
				zap.Error(ErrIvalidTagCurReels))

			return ErrIvalidTagCurReels
		}

		if reelSetMystery.Config.MysteryRNG != "" {
			rng := gameProp.GetTagInt(reelSetMystery.Config.MysteryRNG)
			cs := vw2.Vals[rng]

			curmcode := cs.Int()
			cd.CurMysteryCode = curmcode

			// gameProp.SetVal(GamePropCurMystery, curmcode)

			sc2 := gs.Clone()
			for _, v := range reelSetMystery.MysterySymbolCodes {
				sc2.ReplaceSymbol(v, curmcode)
			}

			reelSetMystery.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)
		} else {
			curm, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ReelSetMystery.OnPlayGame:RandVal",
					zap.Error(err))

				return err
			}

			curmcode := curm.Int()
			cd.CurMysteryCode = curmcode

			// gameProp.SetVal(GamePropCurMystery, curmcode)

			sc2 := gs.Clone()
			for _, v := range reelSetMystery.MysterySymbolCodes {
				sc2.ReplaceSymbol(v, curmcode)
			}

			reelSetMystery.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)
		}
	}

	reelSetMystery.onStepEnd(gameProp, curpr, gp, "")

	gp.AddComponentData(reelSetMystery.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (reelSetMystery *ReelSetMystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[reelSetMystery.Name].(*ReelSetMysteryData)

	if len(cd.UsedScenes) > 0 {
		fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(cd.CurMysteryCode))
		asciigame.OutputScene("after symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (reelSetMystery *ReelSetMystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (reelSetMystery *ReelSetMystery) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.ReelSetMysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("ReelSetMystery.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return reelSetMystery.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (reelSetMystery *ReelSetMystery) NewComponentData() IComponentData {
	return &ReelSetMysteryData{}
}

// EachUsedResults -
func (reelSetMystery *ReelSetMystery) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.ReelSetMysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("ReelSetMystery.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewReelSetMystery(name string) IComponent {
	mystery := &ReelSetMystery{
		BasicComponent:    NewBasicComponent(name),
		MapMysteryWeights: make(map[string]*sgc7game.ValWeights2),
	}

	return mystery
}
