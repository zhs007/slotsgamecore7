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

const ReelSetMysteryTypeName = "reelSetMystery"

type ReelSetMysteryData struct {
	BasicComponentData
	CurMysteryCode int
}

// OnNewGame -
func (reelSetMysteryData *ReelSetMysteryData) OnNewGame(gameProp *GameProperty, component IComponent) {
	reelSetMysteryData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (reelSetMysteryData *ReelSetMysteryData) OnNewStep(gameProp *GameProperty, component IComponent) {
	reelSetMysteryData.BasicComponentData.OnNewStep(gameProp, component)

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
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MysteryRNG           string            `yaml:"mysteryRNG" json:"mysteryRNG"` // 强制用已经使用的随机数结果做 ReelSetMystery
	MysterySymbols       []string          `yaml:"mysterySymbols" json:"mysterySymbols"`
	MapMysteryWeight     map[string]string `yaml:"mapMysteryWeight" json:"mapMysteryWeight"`
}

type ReelSetMystery struct {
	*BasicComponent    `json:"-"`
	Config             *ReelSetMysteryConfig            `json:"config"`
	MapMysteryWeights  map[string]*sgc7game.ValWeights2 `json:"-"`
	MysterySymbolCodes []int                            `json:"-"`
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

	return reelSetMystery.InitEx(cfg, pool)
}

// InitEx -
func (reelSetMystery *ReelSetMystery) InitEx(cfg any, pool *GamePropertyPool) error {
	reelSetMystery.Config = cfg.(*ReelSetMysteryConfig)
	reelSetMystery.Config.ComponentType = ReelSetMysteryTypeName

	for k, v := range reelSetMystery.Config.MapMysteryWeight {
		vw2, err := pool.LoadSymbolWeights(v, "val", "weight", pool.DefaultPaytables, reelSetMystery.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ReelSetMystery.Init:LoadSymbolWeights",
				zap.String("Weight", v),
				zap.Error(err))

			return err
		}

		reelSetMystery.MapMysteryWeights[k] = vw2
	}

	for _, v := range reelSetMystery.Config.MysterySymbols {
		reelSetMystery.MysterySymbolCodes = append(reelSetMystery.MysterySymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	reelSetMystery.onInit(&reelSetMystery.Config.BasicComponentConfig)

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	reelSetMystery.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*ReelSetMysteryData)

	gs := reelSetMystery.GetTargetScene3(gameProp, curpr, prs, &bcd.BasicComponentData, reelSetMystery.Name, "", 0)
	if !reelSetMystery.hasMystery(gs) {
		reelSetMystery.ReTagScene(gameProp, curpr, bcd.TargetSceneIndex, &bcd.BasicComponentData)
	} else {
		vw2, isok := reelSetMystery.MapMysteryWeights[gameProp.GetTagStr(TagCurReels)]
		if !isok {
			goutils.Error("ReelSetMystery.OnPlayGame:MapMysteryWeights",
				zap.String("TagCurReels", gameProp.GetTagStr(TagCurReels)),
				zap.Error(ErrIvalidTagCurReels))

			return "", ErrIvalidTagCurReels
		}

		if reelSetMystery.Config.MysteryRNG != "" {
			rng := gameProp.GetTagInt(reelSetMystery.Config.MysteryRNG)
			cs := vw2.Vals[rng]

			curmcode := cs.Int()
			bcd.CurMysteryCode = curmcode

			// gameProp.SetVal(GamePropCurMystery, curmcode)

			// sc2 := gs.Clone()
			sc2 := gs.CloneEx(gameProp.PoolScene)
			for _, v := range reelSetMystery.MysterySymbolCodes {
				sc2.ReplaceSymbol(v, curmcode)
			}

			reelSetMystery.AddScene(gameProp, curpr, sc2, &bcd.BasicComponentData)
		} else {
			curm, err := vw2.RandVal(plugin)
			if err != nil {
				goutils.Error("ReelSetMystery.OnPlayGame:RandVal",
					zap.Error(err))

				return "", err
			}

			curmcode := curm.Int()
			bcd.CurMysteryCode = curmcode

			// gameProp.SetVal(GamePropCurMystery, curmcode)

			// sc2 := gs.Clone()
			sc2 := gs.CloneEx(gameProp.PoolScene)
			for _, v := range reelSetMystery.MysterySymbolCodes {
				sc2.ReplaceSymbol(v, curmcode)
			}

			reelSetMystery.AddScene(gameProp, curpr, sc2, &bcd.BasicComponentData)
		}
	}

	nc := reelSetMystery.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(reelSetMystery.Name, cd)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (reelSetMystery *ReelSetMystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*ReelSetMysteryData)

	if len(bcd.UsedScenes) > 0 {
		fmt.Printf("mystery is %v\n", gameProp.CurPaytables.GetStringFromInt(bcd.CurMysteryCode))
		asciigame.OutputScene("after symbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (reelSetMystery *ReelSetMystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (reelSetMystery *ReelSetMystery) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.ReelSetMysteryData)
	if !isok {
		goutils.Error("ReelSetMystery.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
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
		BasicComponent:    NewBasicComponent(name, 1),
		MapMysteryWeights: make(map[string]*sgc7game.ValWeights2),
	}

	return mystery
}
