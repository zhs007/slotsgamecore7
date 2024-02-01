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

const MultiLevelReelsTypeName = "multiLevelReels"

type MultiLevelReelsData struct {
	BasicComponentData
	CurLevel int
}

// OnNewGame -
func (multiLevelReelsData *MultiLevelReelsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	multiLevelReelsData.BasicComponentData.OnNewGame(gameProp, component)

	multiLevelReelsData.CurLevel = 0
}

// OnNewStep -
func (multiLevelReelsData *MultiLevelReelsData) OnNewStep(gameProp *GameProperty, component IComponent) {
	multiLevelReelsData.BasicComponentData.OnNewStep(gameProp, component)
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
	Reel           string `yaml:"reel" json:"reel"`
	ReelSetsWeight string `yaml:"reelSetWeight" json:"reelSetWeight"`
	Collector      string `yaml:"collector" json:"collector"`
	CollectorVal   int    `yaml:"collectorVal" json:"collectorVal"`
}

// MultiLevelReelsConfig - configuration for MultiLevelReels
type MultiLevelReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Levels               []*MultiLevelReelsLevelConfig `yaml:"levels" json:"levels"`
}

type MultiLevelReels struct {
	*BasicComponent     `json:"-"`
	Config              *MultiLevelReelsConfig  `json:"config"`
	LevelReelSetWeights []*sgc7game.ValWeights2 `json:"-"`
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

	return multiLevelReels.InitEx(cfg, pool)
}

// InitEx -
func (multiLevelReels *MultiLevelReels) InitEx(cfg any, pool *GamePropertyPool) error {
	multiLevelReels.Config = cfg.(*MultiLevelReelsConfig)
	multiLevelReels.Config.ComponentType = MultiLevelReelsTypeName

	for _, v := range multiLevelReels.Config.Levels {
		if v.ReelSetsWeight != "" {
			vw2, err := pool.LoadStrWeights(v.ReelSetsWeight, multiLevelReels.Config.UseFileMapping)
			if err != nil {
				goutils.Error("MultiLevelReels.Init:LoadSymbolWeights",
					zap.String("Weight", v.ReelSetsWeight),
					zap.Error(err))

				return err
			}

			multiLevelReels.LevelReelSetWeights = append(multiLevelReels.LevelReelSetWeights, vw2)
		}
	}

	if len(multiLevelReels.LevelReelSetWeights) > 0 && len(multiLevelReels.LevelReelSetWeights) != len(multiLevelReels.Config.Levels) {
		goutils.Error("MultiLevelReels.Init:check levels",
			zap.Int("reelsetLength", len(multiLevelReels.LevelReelSetWeights)),
			zap.Int("levelLength", len(multiLevelReels.Config.Levels)),
			zap.Error(ErrIvalidMultiLevelReelsConfig))

		return ErrIvalidMultiLevelReelsConfig
	}

	multiLevelReels.onInit(&multiLevelReels.Config.BasicComponentConfig)

	return nil
}

// // OnNewGame -
// func (multiLevelReels *MultiLevelReels) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[multiLevelReels.Name]

// 	cd.OnNewGame()

// 	return nil
// }

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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) error {

	multiLevelReels.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	mlrcd := cd.(*MultiLevelReelsData)

	reelname := ""
	if multiLevelReels.LevelReelSetWeights != nil {
		val, err := multiLevelReels.LevelReelSetWeights[mlrcd.CurLevel].RandVal(plugin)
		if err != nil {
			goutils.Error("MultiLevelReels.OnPlayGame:ReelSetWeights.RandVal",
				zap.Int("curLevel", mlrcd.CurLevel),
				zap.Error(err))

			return err
		}

		curreels := val.String()
		gameProp.TagStr(TagCurReels, curreels)

		rd, isok := gameProp.Pool.Config.MapReels[curreels]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", mlrcd.CurLevel),
				zap.String("reelset", val.String()),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.CurReels = rd
		reelname = curreels
	} else {
		rd, isok := gameProp.Pool.Config.MapReels[multiLevelReels.Config.Levels[mlrcd.CurLevel].Reel]
		if !isok {
			goutils.Error("MultiLevelReels.OnPlayGame:MapReels",
				zap.Int("curLevel", mlrcd.CurLevel),
				zap.String("reelset", multiLevelReels.Config.Levels[mlrcd.CurLevel].Reel),
				zap.Error(ErrInvalidReels))

			return ErrInvalidReels
		}

		gameProp.TagStr(TagCurReels, multiLevelReels.Config.Levels[mlrcd.CurLevel].Reel)

		gameProp.CurReels = rd
		reelname = multiLevelReels.Config.Levels[mlrcd.CurLevel].Reel
	}

	sc := gameProp.PoolScene.New(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), false)
	sc.ReelName = reelname
	// sc, err := sgc7game.NewGameScene(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight))
	// if err != nil {
	// 	goutils.Error("MultiLevelReels.OnPlayGame:NewGameScene",
	// 		zap.Error(err))

	// 	return err
	// }

	sc.RandReelsWithReelData(gameProp.CurReels, plugin)

	multiLevelReels.AddScene(gameProp, curpr, sc, &mlrcd.BasicComponentData)

	multiLevelReels.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(multiLevelReels.Name, cd)
	// multiLevelReels.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelReels *MultiLevelReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	mlrcd := cd.(*MultiLevelReelsData)

	if len(mlrcd.UsedScenes) > 0 {
		asciigame.OutputScene("initial symbols", pr.Scenes[mlrcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (multiLevelReels *MultiLevelReels) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (multiLevelReels *MultiLevelReels) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	pbcd, isok := pbComponentData.(*sgc7pb.MultiLevelReelsData)
	if !isok {
		goutils.Error("MultiLevelReels.OnStatsWithPB",
			zap.Error(ErrIvalidProto))

		return 0, ErrIvalidProto
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
		BasicComponent: NewBasicComponent(name, 0),
	}

	return multiLevelReels
}
