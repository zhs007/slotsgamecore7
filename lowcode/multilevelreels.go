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

type MultiLevelReelsData struct {
	BasicComponentData
	CurLevel int
}

// OnNewGame -
func (multiLevelReelsData *MultiLevelReelsData) OnNewGame() {
	multiLevelReelsData.BasicComponentData.OnNewGame()

	multiLevelReelsData.CurLevel = 0
}

// OnNewGame -
func (multiLevelReelsData *MultiLevelReelsData) OnNewStep() {
	multiLevelReelsData.BasicComponentData.OnNewStep()
}

// BuildPBComponentData
func (multiLevelReelsData *MultiLevelReelsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MultiLevelReelsData{
		BasicComponentData: multiLevelReelsData.BuildPBBasicComponentData(),
		CurLevel:           int32(multiLevelReelsData.CurLevel),
	}

	return pbcd
}

// MultiLevelReelsLevelConfig - configuration for MultiLevelReels's Level
type MultiLevelReelsLevelConfig struct {
	Reel           string `yaml:"reel"`
	ReelSetsWeight string `yaml:"reelSetWeight"`
	Collector      string `yaml:"collector"`
	CollectorVal   int    `yaml:"collectorVal"`
}

// MultiLevelReelsConfig - configuration for MultiLevelReels
type MultiLevelReelsConfig struct {
	BasicComponentConfig `yaml:",inline"`
	Levels               []*MultiLevelReelsLevelConfig `yaml:"levels"`
}

type MultiLevelReels struct {
	*BasicComponent
	Config              *MultiLevelReelsConfig
	LevelReelSetWeights []*sgc7game.ValWeights2
}

// Init -
func (multiLevelReels *MultiLevelReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiLevelReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiLevelReelsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiLevelReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiLevelReels.Config = cfg

	for _, v := range cfg.Levels {
		if v.ReelSetsWeight != "" {
			vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(v.ReelSetsWeight), "val", "weight", sgc7game.NewStrVal)
			if err != nil {
				goutils.Error("MultiLevelReels.Init:LoadValWeights2FromExcel",
					zap.String("ReelSetsWeight", v.ReelSetsWeight),
					zap.Error(err))

				return err
			}

			multiLevelReels.LevelReelSetWeights = append(multiLevelReels.LevelReelSetWeights, vw2)
		}
	}

	if len(multiLevelReels.LevelReelSetWeights) > 0 && len(multiLevelReels.LevelReelSetWeights) != len(cfg.Levels) {
		goutils.Error("MultiLevelReels.Init:check levels",
			zap.Int("reelsetLength", len(multiLevelReels.LevelReelSetWeights)),
			zap.Int("levelLength", len(cfg.Levels)),
			zap.Error(ErrIvalidMultiLevelReelsConfig))

		return ErrIvalidMultiLevelReelsConfig
	}

	multiLevelReels.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (multiLevelReels *MultiLevelReels) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[multiLevelReels.Name]

	cd.OnNewGame()

	return nil
}

// OnNewStep -
func (multiLevelReels *MultiLevelReels) OnNewStep(gameProp *GameProperty) error {
	multiLevelReels.BasicComponent.OnNewStep(gameProp)

	cd := gameProp.MapComponentData[multiLevelReels.Name].(*MultiLevelReelsData)

	for i := cd.CurLevel + 1; i < len(multiLevelReels.Config.Levels); i++ {
		v := multiLevelReels.Config.Levels[i]

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
func (multiLevelReels *MultiLevelReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[multiLevelReels.Name].(*MultiLevelReelsData)

	if multiLevelReels.LevelReelSetWeights != nil {
		val, err := multiLevelReels.LevelReelSetWeights[cd.CurLevel].RandVal(plugin)
		if err != nil {
			goutils.Error("MultiLevelReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Int("curLevel", cd.CurLevel),
				zap.Error(err))

			return err
		}

		curreels := val.String()
		gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", cd.CurLevel),
				zap.String("reelset", val.String()),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
	} else {
		rd, isok := gameProp.Pool.Config.MapReels[multiLevelReels.Config.Levels[cd.CurLevel].Reel]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", cd.CurLevel),
				zap.String("reelset", multiLevelReels.Config.Levels[cd.CurLevel].Reel),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.TagStr(TagCurReels, multiLevelReels.Config.Levels[cd.CurLevel].Reel)

		gameProp.CurReels = rd
	}

	sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	if err != nil {
		goutils.Error("MultiLevelReels.OnPlayGame:NewGameScene",
			zap.Error(err))

		return err
	}

	sc.RandReelsWithReelData(gameProp.CurReels, plugin)

	multiLevelReels.AddScene(gameProp, curpr, sc, &cd.BasicComponentData)

	multiLevelReels.onStepEnd(gameProp, curpr, gp, "")

	gp.AddComponentData(multiLevelReels.Name, cd)
	// multiLevelReels.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelReels *MultiLevelReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[multiLevelReels.Name].(*MultiLevelReelsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (multiLevelReels *MultiLevelReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (multiLevelReels *MultiLevelReels) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.MultiLevelReelsData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("MultiLevelReels.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return multiLevelReels.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (multiLevelReels *MultiLevelReels) NewComponentData() IComponentData {
	return &MultiLevelReelsData{}
}

// EachUsedResults -
func (multiLevelReels *MultiLevelReels) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.MultiLevelReelsData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("MultiLevelReels.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewMultiLevelReels(name string) IComponent {
	multiLevelReels := &MultiLevelReels{
		BasicComponent: NewBasicComponent(name),
	}

	return multiLevelReels
}
