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

type MultiLevelMysteryData struct {
	BasicComponentData
	CurLevel       int
	CurMysteryCode int
}

// OnNewGame -
func (multiLevelMysteryData *MultiLevelMysteryData) OnNewGame() {
	multiLevelMysteryData.BasicComponentData.OnNewGame()

	multiLevelMysteryData.CurLevel = 0
}

// OnNewStep -
func (multiLevelMysteryData *MultiLevelMysteryData) OnNewStep() {
	multiLevelMysteryData.BasicComponentData.OnNewStep()

	multiLevelMysteryData.CurMysteryCode = -1
}

// BuildPBComponentData
func (multiLevelMysteryData *MultiLevelMysteryData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.MultiLevelMysteryData{
		BasicComponentData: multiLevelMysteryData.BuildPBBasicComponentData(),
		CurLevel:           int32(multiLevelMysteryData.CurLevel),
		CurMysteryCode:     int32(multiLevelMysteryData.CurMysteryCode),
	}

	return pbcd
}

// MultiLevelMysteryLevelConfig - configuration for MultiLevelMystery's Level
type MultiLevelMysteryLevelConfig struct {
	MysteryWeight string `yaml:"mysteryWeight"`
	Collector     string `yaml:"collector"`
	CollectorVal  int    `yaml:"collectorVal"`
}

// MultiLevelMysteryConfig - configuration for MultiLevelMystery
type MultiLevelMysteryConfig struct {
	BasicComponentConfig   `yaml:",inline"`
	Mystery                string                          `yaml:"mystery"`
	Mysterys               []string                        `yaml:"mysterys"`
	Levels                 []*MultiLevelMysteryLevelConfig `yaml:"levels"`
	MysteryTriggerFeatures []*MysteryTriggerFeatureConfig  `yaml:"mysteryTriggerFeatures"`
}

type MultiLevelMystery struct {
	*BasicComponent
	Config                   *MultiLevelMysteryConfig
	MapMysteryTriggerFeature map[int]*MysteryTriggerFeatureConfig
	LevelMysteryWeights      []*sgc7game.ValWeights2
	MysterySymbols           []int
}

// maskOtherScene -
func (multiLevelMystery *MultiLevelMystery) maskOtherScene(gs *sgc7game.GameScene, symbolCode int) *sgc7game.GameScene {
	cgs := gs.Clone()

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
func (multiLevelMystery *MultiLevelMystery) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MultiLevelMystery.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MultiLevelMysteryConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MultiLevelMystery.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	multiLevelMystery.Config = cfg

	for _, v := range cfg.Levels {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(pool.Config.GetPath(v.MysteryWeight), "val", "weight", pool.DefaultPaytables)
		if err != nil {
			goutils.Error("MultiLevelMystery.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("MysteryWeight", v.MysteryWeight),
				zap.Error(err))

			return err
		}

		multiLevelMystery.LevelMysteryWeights = append(multiLevelMystery.LevelMysteryWeights, vw2)
	}

	if len(multiLevelMystery.Config.Mysterys) > 0 {
		for _, v := range multiLevelMystery.Config.Mysterys {
			multiLevelMystery.MysterySymbols = append(multiLevelMystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[v])
		}
	} else {
		multiLevelMystery.MysterySymbols = append(multiLevelMystery.MysterySymbols, pool.DefaultPaytables.MapSymbols[multiLevelMystery.Config.Mystery])
	}

	for _, v := range cfg.MysteryTriggerFeatures {
		symbolCode := pool.DefaultPaytables.MapSymbols[v.Symbol]

		multiLevelMystery.MapMysteryTriggerFeature[symbolCode] = v
	}

	multiLevelMystery.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (multiLevelMystery *MultiLevelMystery) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[multiLevelMystery.Name]

	cd.OnNewGame()

	return nil
}

// OnNewStep -
func (multiLevelMystery *MultiLevelMystery) OnNewStep(gameProp *GameProperty) error {
	multiLevelMystery.BasicComponent.OnNewStep(gameProp)

	cd := gameProp.MapComponentData[multiLevelMystery.Name].(*MultiLevelMysteryData)

	for i := cd.CurLevel + 1; i < len(multiLevelMystery.Config.Levels); i++ {
		v := multiLevelMystery.Config.Levels[i]

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
func (multiLevelMystery *MultiLevelMystery) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[multiLevelMystery.Name].(*MultiLevelMysteryData)

	gs := multiLevelMystery.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")

	if gs.HasSymbols(multiLevelMystery.MysterySymbols) {
		curm, err := multiLevelMystery.LevelMysteryWeights[cd.CurLevel].RandVal(plugin)
		if err != nil {
			goutils.Error("MultiLevelMystery.OnPlayGame:RandVal",
				zap.Error(err))

			return err
		}

		curmcode := curm.Int()
		cd.CurMysteryCode = curmcode

		// gameProp.SetVal(GamePropCurMystery, curm.Int())

		sc2 := gs.Clone()
		for _, v := range multiLevelMystery.MysterySymbols {
			sc2.ReplaceSymbol(v, curm.Int())
		}

		multiLevelMystery.AddScene(gameProp, curpr, sc2, &cd.BasicComponentData)

		v, isok := multiLevelMystery.MapMysteryTriggerFeature[curmcode]
		if isok {
			if v.RespinFirstComponent != "" {
				os := multiLevelMystery.maskOtherScene(sc2, curmcode)

				gameProp.Respin(curpr, gp, v.RespinFirstComponent, sc2, os)

				return nil
			}
		}
	} else {
		multiLevelMystery.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, &cd.BasicComponentData)
	}

	multiLevelMystery.onStepEnd(gameProp, curpr, gp, "")

	// gp.AddComponentData(multiLevelMystery.Name, cd)
	// multiLevelMystery.BuildPBComponent(gp)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (multiLevelMystery *MultiLevelMystery) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[multiLevelMystery.Name].(*MultiLevelMysteryData)

	if len(cd.UsedScenes) > 0 {
		fmt.Printf("mystery is %v\n", gameProp.GetStrVal(cd.CurMysteryCode))
		asciigame.OutputScene("after symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (multiLevelMystery *MultiLevelMystery) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (multiLevelMystery *MultiLevelMystery) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.MultiLevelMysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("MultiLevelMystery.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return multiLevelMystery.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (multiLevelMystery *MultiLevelMystery) NewComponentData() IComponentData {
	return &MultiLevelMysteryData{}
}

// EachUsedResults -
func (multiLevelMystery *MultiLevelMystery) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.MultiLevelMysteryData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("MultiLevelMystery.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewMultiLevelMystery(name string) IComponent {
	multiLevelMystery := &MultiLevelMystery{
		BasicComponent:           NewBasicComponent(name),
		MapMysteryTriggerFeature: make(map[int]*MysteryTriggerFeatureConfig),
	}

	return multiLevelMystery
}
