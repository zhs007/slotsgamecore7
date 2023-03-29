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

type CollectorData struct {
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
}

// OnNewGame -
func (collectorData *CollectorData) OnNewGame() {
	collectorData.Val = 0
}

// OnNewGame -
func (collectorData *CollectorData) OnNewStep() {
	collectorData.NewCollector = 0
}

// BuildPBComponentData
func (collectorData *CollectorData) BuildPBComponentData() proto.Message {
	return &sgc7pb.CollectorData{
		Val:          int32(collectorData.Val),
		NewCollector: int32(collectorData.NewCollector),
	}
}

// CollectorConfig - configuration for Collector
type CollectorConfig struct {
	BasicComponentConfig `yaml:",inline"`
	MaxVal               int `yaml:"maxVal"`
}

type Collector struct {
	*BasicComponent
	Config *CollectorConfig
}

// Init -
func (collector *Collector) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("BasicReels.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &CollectorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("BasicReels.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	collector.Config = cfg

	collector.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame - 因为 BasicComponent 考虑到效率，没有执行ComponentData的OnNewGame，所以这里需要特殊处理
func (collector *Collector) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[collector.Name]

	cd.OnNewGame()

	return nil
}

// playgame
func (collector *Collector) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	gameProp.SetStrVal(GamePropNextComponent, collector.Config.DefaultNextComponent)

	collector.onStepEnd(gameProp, curpr, gp)

	gp.AddComponentData(collector.Name, gameProp.MapComponentData[collector.Name])

	return nil
}

// OnAsciiGame - outpur to asciigame
func (collector *Collector) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[collector.Name].(*CollectorData)

	if cd.NewCollector <= 0 {
		fmt.Printf("%v dose not collect new value, the collector value is %v", collector.Name, cd.Val)
	} else {
		fmt.Printf("%v collect %v. the collector value is %v", collector.Name, cd.NewCollector, cd.Val)
	}

	return nil
}

// OnStats
func (collector *Collector) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (collector *Collector) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.CollectorData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("BasicComponent.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return 0, nil
}

// NewComponentData -
func (collector *Collector) NewComponentData() IComponentData {
	return &CollectorData{}
}

// EachUsedResults -
func (collector *Collector) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewCollector(name string) IComponent {
	collector := &Collector{
		BasicComponent: NewBasicComponent(name),
	}

	return collector
}
