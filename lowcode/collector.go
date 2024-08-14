package lowcode

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const CollectorTypeName = "collector"

type CollectorData struct {
	BasicComponentData
	Val          int // 当前总值, Current total value
	NewCollector int // 这一个step收集到的, The values collected in this step
}

// OnNewGame -
func (collectorData *CollectorData) OnNewGame(gameProp *GameProperty, component IComponent) {
	collectorData.BasicComponentData.OnNewGame(gameProp, component)

	collectorData.Val = 0
	collectorData.NewCollector = 0
}

// onNewStep -
func (collectorData *CollectorData) onNewStep() {
	collectorData.NewCollector = 0
}

// // ChgConfigIntVal -
// func (collectorData *CollectorData) ChgConfigIntVal(key string, off int) {
// 	if key == CCVValueNum {
// 		// collectorData.AddRespinTimes(off)
// 	} else {
// 		collectorData.BasicComponentData.ChgConfigIntVal(key, off)
// 	}
// }

// SetConfigIntVal -
func (collectorData *CollectorData) SetConfigIntVal(key string, val int) {
	if key == CCVValueNum {
		collectorData.Val = val
	} else {
		collectorData.BasicComponentData.ChgConfigIntVal(key, val)
	}
}

// GetOutput -
func (collectorData *CollectorData) GetOutput() int {
	return collectorData.Val
}

// GetVal -
func (collectorData *CollectorData) GetVal(key string) (int, bool) {
	if key == CVValue {
		return collectorData.Val, true
	}

	return 0, false
}

// // ChgConfigIntVal -
// func (collectorData *CollectorData) ChgConfigIntVal(key string, off int) {
// 	if key == CCVValueNum {
// 		collectorData.Val += off
// 	}
// }

// Clone
func (collectorData *CollectorData) Clone() IComponentData {
	target := &CollectorData{
		BasicComponentData: collectorData.CloneBasicComponentData(),
		Val:                collectorData.Val,
		NewCollector:       collectorData.NewCollector,
	}

	return target
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
	BasicComponentConfig `yaml:",inline" json:",inline"`
	// Symbol               string           `yaml:"symbol" json:"symbol"`
	MaxVal           int              `yaml:"maxVal" json:"maxVal"`
	PerLevelAwards   []*Award         `yaml:"perLevelAwards" json:"perLevelAwards"`
	MapSPLevelAwards map[int][]*Award `yaml:"mapSPLevelAwards" json:"mapSPLevelAwards"`
	IsCycle          bool             `yaml:"isCycle" json:"isCycle"`
}

// SetLinkComponent
func (cfg *CollectorConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type Collector struct {
	*BasicComponent `json:"-"`
	Config          *CollectorConfig `json:"config"`
	// SymbolCode      int              `json:"-"`
}

// Init -
func (collector *Collector) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Collector.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &CollectorConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Collector.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return collector.InitEx(cfg, pool)
}

// InitEx -
func (collector *Collector) InitEx(cfg any, pool *GamePropertyPool) error {
	collector.Config = cfg.(*CollectorConfig)
	collector.Config.ComponentType = CollectorTypeName

	// collector.SymbolCode = pool.DefaultPaytables.MapSymbols[collector.Config.Symbol]

	if collector.Config.PerLevelAwards != nil {
		for _, v := range collector.Config.PerLevelAwards {
			v.Init()
		}
	}

	if collector.Config.MapSPLevelAwards != nil {
		for _, lst := range collector.Config.MapSPLevelAwards {
			for _, v := range lst {
				v.Init()
			}
		}
	}

	collector.onInit(&collector.Config.BasicComponentConfig)

	return nil
}

// // OnNewGame - 因为 BasicComponent 考虑到效率，没有执行ComponentData的OnNewGame，所以这里需要特殊处理
// func (collector *Collector) OnNewGame(gameProp *GameProperty) error {
// 	cd := gameProp.MapComponentData[collector.Name]

// 	cd.OnNewGame()

// 	return nil
// }

// add -
func (collector *Collector) add(plugin sgc7plugin.IPlugin, num int, cd *CollectorData, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, noProcLevelUp bool) error {
	if num <= 0 {
		return nil
	}

	// if cd == nil {
	// 	cd = gameProp.MapComponentData[collector.Name].(*CollectorData)
	// }

	if collector.Config.MaxVal > 0 && !collector.Config.IsCycle && cd.Val == collector.Config.MaxVal {
		return nil
	}

	cd.NewCollector += num
	oldval := cd.Val
	cd.Val += num
	if collector.Config.MaxVal > 0 {
		if cd.Val >= collector.Config.MaxVal {
			if collector.Config.IsCycle {
				cd.Val -= collector.Config.MaxVal
			} else {
				cd.Val = collector.Config.MaxVal
			}
		}
	}

	if num > 0 && !noProcLevelUp && oldval != cd.Val {
		for i := 1; i <= num; i++ {
			cl := oldval + i
			if cl >= collector.Config.MaxVal {
				collector.onLevelUp(plugin, gameProp, curpr, gp, -1, false)

				break
			} else {
				collector.onLevelUp(plugin, gameProp, curpr, gp, cl, false)
			}
		}
	}

	return nil
}

// onLevelUp -
func (collector *Collector) onLevelUp(plugin sgc7plugin.IPlugin, gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, newLevel int, noProcSPLevel bool) error {
	if collector.Config.PerLevelAwards != nil {
		for _, v := range collector.Config.PerLevelAwards {
			gameProp.procAward(plugin, v, curpr, gp, false)
		}
	}

	if noProcSPLevel {
		return nil
	}

	sp, isok := collector.Config.MapSPLevelAwards[newLevel]
	if isok {
		for _, v := range sp {
			gameProp.procAward(plugin, v, curpr, gp, false)
		}
	} else {
		if newLevel == -1 {
			sp1, isok1 := collector.Config.MapSPLevelAwards[collector.Config.MaxVal]
			if isok1 {
				for _, v := range sp1 {
					gameProp.procAward(plugin, v, curpr, gp, false)
				}
			}
		} else if newLevel == collector.Config.MaxVal {
			sp1, isok1 := collector.Config.MapSPLevelAwards[-1]
			if isok1 {
				for _, v := range sp1 {
					gameProp.procAward(plugin, v, curpr, gp, false)
				}
			}
		}
	}

	return nil
}

// playgame
func (collector *Collector) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	ccd := cd.(*CollectorData)
	ccd.onNewStep()

	off, isok := ccd.GetConfigIntVal(CCVValueNum)
	if isok {
		err := collector.add(plugin, off, ccd, gameProp, curpr, gp, false)
		if err != nil {
			goutils.Error("Collector.OnPlayGame:add:off",
				goutils.Err(err))

			return "", err
		}

		ccd.ClearConfigIntVal(CCVValueNum)
	}

	// gs := collector.GetTargetScene3(gameProp, curpr, prs, 0)

	// nn := gs.CountSymbolEx(func(cursymbol int, x, y int) bool {
	// 	return cursymbol == collector.SymbolCode
	// })

	// collector.add(plugin, nn, ccd, gameProp, curpr, gp, false)

	nc := collector.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (collector *Collector) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {

	ccd := cd.(*CollectorData)

	if ccd.NewCollector <= 0 {
		fmt.Printf("%v dose not collect new value, the collector value is %v\n", collector.Name, ccd.Val)
	} else {
		fmt.Printf("%v collect %v. the collector value is %v\n", collector.Name, ccd.NewCollector, ccd.Val)
	}

	return nil
}

// // OnStats
// func (collector *Collector) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	if feature != nil && len(lst) > 0 {
// 		if feature.RespinEndingStatus != nil {
// 			pbcd, lastpr := findLastPBComponentDataEx(lst, feature.RespinEndingName, collector.Name)

// 			if pbcd != nil {
// 				collector.OnStatsWithPB(feature, pbcd, lastpr)
// 			}
// 		}

// 		if feature.RespinStartStatus != nil {
// 			pbcd, lastpr := findFirstPBComponentDataEx(lst, feature.RespinStartName, collector.Name)

// 			if pbcd != nil {
// 				collector.OnStatsWithPB(feature, pbcd, lastpr)
// 			}
// 		}

// 		if feature.RespinStartStatusEx != nil {
// 			pbs, prs := findAllPBComponentDataEx(lst, feature.RespinStartNameEx, collector.Name)

// 			if len(pbs) > 0 {
// 				for i, v := range pbs {
// 					collector.OnStatsWithPB(feature, v, prs[i])
// 				}
// 			}
// 		}
// 	}

// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (collector *Collector) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	pbcd, isok := pbComponentData.(*sgc7pb.CollectorData)
// 	if !isok {
// 		goutils.Error("Collector.OnStatsWithPB",
// 			goutils.Err(ErrIvalidProto))

// 		return 0, ErrIvalidProto
// 	}

// 	if feature.RespinEndingStatus != nil {
// 		feature.RespinEndingStatus.AddStatus(int(pbcd.Val))
// 	}

// 	if feature.RespinStartStatus != nil {
// 		feature.RespinStartStatus.AddStatus(int(pbcd.Val - pbcd.NewCollector))
// 	}

// 	if feature.RespinStartStatusEx != nil {
// 		feature.RespinStartStatusEx.AddStatus(int(pbcd.Val - pbcd.NewCollector))
// 	}

// 	return 0, nil
// }

// NewComponentData -
func (collector *Collector) NewComponentData() IComponentData {
	return &CollectorData{}
}

// EachUsedResults -
func (collector *Collector) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewCollector(name string) IComponent {
	collector := &Collector{
		BasicComponent: NewBasicComponent(name, 1),
	}

	return collector
}

// "configuration": {},
type jsonCollector struct {
	MaxVal  int  `json:"maxVal"`
	IsCycle bool `json:"isCycle"`
}

func (jcfg *jsonCollector) build() *CollectorConfig {
	cfg := &CollectorConfig{
		MaxVal:  jcfg.MaxVal,
		IsCycle: jcfg.IsCycle,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseCollector(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseCollector:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseCollector:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonCollector{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseCollector:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, mapawards, err := parseCollectorControllers(ctrls)
		if err != nil {
			goutils.Error("parseScatterTrigger:parseCollectorControllers",
				goutils.Err(err))

			return "", err
		}

		if len(awards) > 0 {
			cfgd.PerLevelAwards = awards
		}

		if len(mapawards) > 0 {
			cfgd.MapSPLevelAwards = mapawards
		}
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: CollectorTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
