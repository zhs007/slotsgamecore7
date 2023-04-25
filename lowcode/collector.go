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
	BasicComponentData
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
}

// OnNewGame -
func (collectorData *CollectorData) OnNewGame() {
	collectorData.BasicComponentData.OnNewGame()

	collectorData.Val = 0
}

// OnNewGame -
func (collectorData *CollectorData) OnNewStep() {
	collectorData.BasicComponentData.OnNewStep()

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
	Symbol               string                 `yaml:"symbol"`
	MaxVal               int                    `yaml:"maxVal"`
	PerLevelAwards       []*AwardConfig         `yaml:"perLevelAwards"`
	MapSPLevelAwards     map[int][]*AwardConfig `yaml:"mapSPLevelAwards"`
}

type Collector struct {
	*BasicComponent
	Config           *CollectorConfig
	SymbolCode       int
	PerLevelAwards   []*Award
	MapSPLevelAwards map[int][]*Award
}

// Init -
func (collector *Collector) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Collector.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &CollectorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Collector.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	collector.Config = cfg

	collector.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	if cfg.PerLevelAwards != nil {
		for _, v := range cfg.PerLevelAwards {
			collector.PerLevelAwards = append(collector.PerLevelAwards, NewArard(v))
		}
	}

	if cfg.MapSPLevelAwards != nil {
		collector.MapSPLevelAwards = make(map[int][]*Award)

		for k, lst := range cfg.MapSPLevelAwards {
			awards := []*Award{}

			for _, v := range lst {
				awards = append(awards, NewArard(v))
			}

			collector.MapSPLevelAwards[k] = awards
		}
	}

	collector.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame - 因为 BasicComponent 考虑到效率，没有执行ComponentData的OnNewGame，所以这里需要特殊处理
func (collector *Collector) OnNewGame(gameProp *GameProperty) error {
	cd := gameProp.MapComponentData[collector.Name]

	cd.OnNewGame()

	return nil
}

// onLevelUp -
func (collector *Collector) onLevelUp(gameProp *GameProperty, curpr *sgc7game.PlayResult, newLevel int, noProcSPLevel bool) error {
	if collector.PerLevelAwards != nil {
		for _, v := range collector.PerLevelAwards {
			gameProp.procAward(v, curpr)
		}
	}

	if noProcSPLevel {
		return nil
	}

	sp, isok := collector.MapSPLevelAwards[newLevel]
	if isok {
		for _, v := range sp {
			gameProp.procAward(v, curpr)
		}
	}

	return nil
}

// playgame
func (collector *Collector) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	cd := gameProp.MapComponentData[collector.Name].(*CollectorData)

	gs := collector.GetTargetScene(gameProp, curpr, &cd.BasicComponentData, "")

	nn := gs.CountSymbolEx(func(cursymbol int, x, y int) bool {
		return cursymbol == collector.SymbolCode
	})

	cd.NewCollector = nn
	cd.Val += nn
	if collector.Config.MaxVal > 0 {
		if cd.Val > collector.Config.MaxVal {
			cd.Val = collector.Config.MaxVal
		}
	}

	if nn > 0 {
		for i := 1; i <= nn; i++ {
			cl := cd.Val + i
			if cl > collector.Config.MaxVal {
				collector.onLevelUp(gameProp, curpr, collector.Config.MaxVal, false)
			} else {
				collector.onLevelUp(gameProp, curpr, cl, true)
			}
		}
	}

	// gameProp.SetStrVal(GamePropNextComponent, collector.Config.DefaultNextComponent)

	collector.onStepEnd(gameProp, curpr, gp, "")

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
		goutils.Error("Collector.OnStatsWithPB:UnmarshalTo",
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
