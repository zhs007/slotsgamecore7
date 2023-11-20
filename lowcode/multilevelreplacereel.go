package lowcode

import (
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

const MultiLevelReplaceReelTypeName = "multiLevelReplaceReel"

type MultiLevelReplaceReelData struct {
	BasicComponentData
	CurLevel int
}

// OnNewGame -
func (multiLevelReplaceReelData *MultiLevelReplaceReelData) OnNewGame() {
	multiLevelReplaceReelData.BasicComponentData.OnNewGame()

	multiLevelReplaceReelData.CurLevel = 0
}

// OnNewStep -
func (multiLevelReplaceReelData *MultiLevelReplaceReelData) OnNewStep() {
	multiLevelReplaceReelData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (multiLevelReplaceReelData *MultiLevelReplaceReelData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MultiLevelReplaceReelData{
		BasicComponentData: multiLevelReplaceReelData.BuildPBBasicComponentData(),
		CurLevel:           int32(multiLevelReplaceReelData.CurLevel),
	}

	return pbcd
}

// MultiLevelReplaceReelLevelConfig - configuration for MultiLevelReplaceReelData's Level
type MultiLevelReplaceReelLevelConfig struct {
	Reels           map[int][]string `yaml:"reels" json:"reels"` // x - [0, width)
	SymbolCodeReels map[int][]int    `yaml:"-" json:"-"`
	Collector       string           `yaml:"collector" json:"collector"`
	CollectorVal    int              `yaml:"collectorVal" json:"collectorVal"`
}

// MultiLevelReplaceReelDataConfig - configuration for MultiLevelReplaceReelData
type MultiLevelReplaceReelDataConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Levels               []*MultiLevelReplaceReelLevelConfig `yaml:"levels" json:"levels"`
}

type MultiLevelReplaceReel struct {
	*BasicComponent `json:"-"`
	Config          *MultiLevelReplaceReelDataConfig `json:"config"`
}

// Init -
func (multiLevelReplaceReel *MultiLevelReplaceReel) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiLevelReplaceReel.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiLevelReplaceReelDataConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiLevelReplaceReel.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return multiLevelReplaceReel.InitEx(cfg, pool)
}

// InitEx -
func (multiLevelReplaceReel *MultiLevelReplaceReel) InitEx(cfg any, pool *GamePropertyPool) error {
	multiLevelReplaceReel.Config = cfg.(*MultiLevelReplaceReelDataConfig)
	multiLevelReplaceReel.Config.ComponentType = MultiLevelReplaceReelTypeName

	for _, v := range multiLevelReplaceReel.Config.Levels {
		if v.Reels != nil {
			v.SymbolCodeReels = make(map[int][]int)

			for ri, symbols := range v.Reels {
				scs := []int{}
				for _, s := range symbols {
					scs = append(scs, pool.DefaultPaytables.MapSymbols[s])
				}

				v.SymbolCodeReels[ri] = scs
			}
		}
	}

	multiLevelReplaceReel.onInit(&multiLevelReplaceReel.Config.BasicComponentConfig)

	return nil
}

// // OnNewGame -
// func (multiLevelReplaceReel *MultiLevelReplaceReel) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[multiLevelReplaceReel.Name]

// 	cd.OnNewGame()

// 	return nil
// }

// OnNewStep -
func (multiLevelReplaceReel *MultiLevelReplaceReel) OnNewStep(gameProp *GameProperty) error {
	multiLevelReplaceReel.BasicComponent.OnNewStep(gameProp)

	cd := gameProp.MapComponentData[multiLevelReplaceReel.Name].(*MultiLevelReplaceReelData)

	for i := cd.CurLevel + 1; i < len(multiLevelReplaceReel.Config.Levels); i++ {
		v := multiLevelReplaceReel.Config.Levels[i]

		collectorData, isok := gameProp.MapComponentData[v.Collector].(*CollectorData)
		if isok {
			if collectorData.Val >= v.CollectorVal {
				cd.CurLevel = i
			} else {
				break
			}
		}
	}

	return nil
}

// playgame
func (multiLevelReplaceReel *MultiLevelReplaceReel) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	multiLevelReplaceReel.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[multiLevelReplaceReel.Name].(*MultiLevelReplaceReelData)

	if multiLevelReplaceReel.Config.Levels[cd.CurLevel].SymbolCodeReels != nil {
		gs := multiLevelReplaceReel.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")

		// sc := gs.Clone()
		sc := gs.CloneEx(gameProp.PoolScene)

		for x, reel := range multiLevelReplaceReel.Config.Levels[cd.CurLevel].SymbolCodeReels {
			copy(sc.Arr[x], reel)
			// for y, s := range reel {
			// 	sc.Arr[x][y] = s
			// }
		}

		multiLevelReplaceReel.AddScene(gameProp, curpr, sc, &cd.BasicComponentData)
	} else {
		multiLevelReplaceReel.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")

		multiLevelReplaceReel.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, &cd.BasicComponentData)
	}

	multiLevelReplaceReel.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(multiLevelReplaceReel.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelReplaceReel *MultiLevelReplaceReel) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[multiLevelReplaceReel.Name].(*MultiLevelReplaceReelData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after replaceReel symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (multiLevelReplaceReel *MultiLevelReplaceReel) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (multiLevelReplaceReel *MultiLevelReplaceReel) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.MultiLevelReplaceReelData)
	if !isok {
		goutils.Error("MultiLevelReplaceReel.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
	}

	return multiLevelReplaceReel.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (multiLevelReplaceReel *MultiLevelReplaceReel) NewComponentData() IComponentData {
	return &MultiLevelReplaceReelData{}
}

// EachUsedResults -
func (multiLevelReplaceReel *MultiLevelReplaceReel) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.MultiLevelReplaceReelData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("multiLevelReplaceReel.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewMultiLevelReplaceReel(name string) IComponent {
	multiLevelReplaceReel := &MultiLevelReplaceReel{
		BasicComponent: NewBasicComponent(name),
	}

	return multiLevelReplaceReel
}
