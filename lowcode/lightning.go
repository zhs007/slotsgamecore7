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

const (
	LightningTypeVal = "val"
	LightningTypeMul = "mul"

	LightningFeatureCollector = "collector"
)

type LightningData struct {
	BasicComponentData
	Collector    int
	Val          int
	Mul          int
	NewConnector int
}

// OnNewGame -
func (lightningData *LightningData) OnNewGame() {
	lightningData.BasicComponentData.OnNewGame()
}

// OnNewGame -
func (lightningData *LightningData) OnNewStep() {
	lightningData.BasicComponentData.OnNewStep()

	lightningData.Collector = 0
	lightningData.Val = 0
	lightningData.Mul = 0
	lightningData.NewConnector = 0
}

// BuildPBComponentData
func (lightningData *LightningData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.LightningData{
		BasicComponentData: lightningData.BuildPBBasicComponentData(),
		Collector:          int32(lightningData.Collector),
		Val:                int32(lightningData.Val),
		Mul:                int32(lightningData.Mul),
		NewConnector:       int32(lightningData.NewConnector),
	}

	return pbcd
}

// LightningTriggerFeatureConfig - configuration for lightning trigger feature
type LightningTriggerFeatureConfig struct {
	Symbol  string `yaml:"symbol"`  // like NEW_COLLECTOR
	Feature string `yaml:"feature"` // like collector
}

// LightningSymbolValConfig - configuration for symbol value
type LightningSymbolValConfig struct {
	Symbol string `yaml:"symbol"`
	Weight string `yaml:"weight"`
	Type   string `yaml:"type"` // like val or mul
}

// LightningConfig - configuration for Lightning
type LightningConfig struct {
	BasicComponentConfig  `yaml:",inline"`
	Symbol                string                           `yaml:"symbol"`
	Weight                string                           `yaml:"weight"`
	SymbolVals            []*LightningSymbolValConfig      `yaml:"symbolVals"`
	SymbolTriggerFeatures []*LightningTriggerFeatureConfig `yaml:"symbolTriggerFeatures"`
	EndingFirstComponent  string                           `yaml:"endingFirstComponent"`
}

// LightningSymbolData - symbol data for Lightning
type LightningSymbolData struct {
	SymbolCode int
	Weight     *sgc7game.ValWeights2
	Config     *LightningSymbolValConfig
}

type Lightning struct {
	*BasicComponent
	Config                   *LightningConfig
	SymbolCode               int
	Weight                   *sgc7game.ValWeights2
	MapSymbols               map[int]*LightningSymbolData
	MapSymbolTriggerFeatures map[int]*LightningTriggerFeatureConfig
	ValSymbolCode            int
	MulSymbolCode            int
	CollectorSymbolCode      int
}

// Init -
func (lightning *Lightning) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("Lightning.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &LightningConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("Lightning.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	lightning.Config = cfg

	lightning.SymbolCode = pool.DefaultPaytables.MapSymbols[cfg.Symbol]

	if lightning.Config.Weight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(pool.Config.GetPath(lightning.Config.Weight, lightning.Config.UseFileMapping), "val", "weight", pool.DefaultPaytables)
		if err != nil {
			goutils.Error("Lightning.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("Weight", lightning.Config.Weight),
				zap.Error(err))

			return err
		}

		lightning.Weight = vw2
	}

	for _, v := range lightning.Config.SymbolVals {
		symbolCode := pool.DefaultPaytables.MapSymbols[v.Symbol]

		sd := &LightningSymbolData{
			SymbolCode: symbolCode,
			Config:     v,
		}

		vw2, err := sgc7game.LoadValWeights2FromExcel(pool.Config.GetPath(v.Weight, lightning.Config.UseFileMapping), "val", "weight", sgc7game.NewIntVal[int])
		if err != nil {
			goutils.Error("Lightning.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("Weight", v.Weight),
				zap.Error(err))

			return err
		}

		sd.Weight = vw2

		lightning.MapSymbols[symbolCode] = sd

		if v.Type == LightningTypeVal {
			lightning.ValSymbolCode = symbolCode
		} else if v.Type == LightningTypeMul {
			lightning.MulSymbolCode = symbolCode
		}
	}

	for _, v := range cfg.SymbolTriggerFeatures {
		symbolCode := pool.DefaultPaytables.MapSymbols[v.Symbol]

		lightning.MapSymbolTriggerFeatures[symbolCode] = v

		if v.Feature == LightningFeatureCollector {
			lightning.CollectorSymbolCode = symbolCode
		}
	}

	lightning.onInit(&cfg.BasicComponentConfig)

	return nil
}

// playgame
func (lightning *Lightning) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	lightning.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[lightning.Name].(*LightningData)

	if len(prs) <= 0 {
		goutils.Error("Lightning:prs",
			zap.Error(ErrIvalidPlayResultLength))

		return ErrIvalidPlayResultLength
	}

	preps := prs[len(prs)-1]
	pregp, isok := preps.CurGameModParams.(*GameParams)
	if !isok {
		goutils.Error("Lightning:preps",
			zap.Error(ErrIvalidCurGameModParams))

		return ErrIvalidCurGameModParams
	}

	// gs := pregp.LastScene.Clone()
	// os := pregp.LastOtherScene.Clone()
	gs := pregp.LastScene.CloneEx(gameProp.PoolScene)
	os := pregp.LastOtherScene.CloneEx(gameProp.PoolScene)

	isCollector := false
	collectorX := 0
	collectorY := 0
	val := 0
	mul := 1
	lastCollector := 0
	arrpos := []int{}
	for x, arr := range os.Arr {
		for y, v := range arr {
			if v >= 0 {
				arrpos = append(arrpos, x, y)

				if gs.Arr[x][y] == lightning.CollectorSymbolCode {
					lastCollector = os.Arr[x][y]
				}

				cs, err := lightning.Weight.RandVal(plugin)
				if err != nil {
					goutils.Error("Lightning.OnPlayGame:RandVal",
						zap.Error(err))

					return err
				}

				symbolCode := cs.Int()

				gs.Arr[x][y] = symbolCode
				sw, isok := lightning.MapSymbols[symbolCode]
				if isok {
					cv, err := sw.Weight.RandVal(plugin)
					if err != nil {
						goutils.Error("Lightning.OnPlayGame:RandVal",
							zap.Int("symbol", symbolCode),
							zap.Error(err))

						return err
					}

					os.Arr[x][y] = cv.Int()

					if sw.Config.Type == LightningTypeVal {
						val += os.Arr[x][y]
					} else if sw.Config.Type == LightningTypeMul {
						mul *= os.Arr[x][y]
					}
				} else {
					os.Arr[x][y] = 0
				}

				_, isok = lightning.MapSymbolTriggerFeatures[symbolCode]
				if isok {
					collectorX = x
					collectorY = y

					isCollector = true
				}
			}
		}
	}

	lightning.AddScene(gameProp, curpr, gs, &cd.BasicComponentData)
	lightning.AddOtherScene(gameProp, curpr, os, &cd.BasicComponentData)

	cd.Collector = lastCollector
	cd.Val = val
	cd.Mul = mul

	if isCollector {
		cd.NewConnector = lastCollector + val*mul

		os.Arr[collectorX][collectorY] = cd.NewConnector

		gameProp.Respin(curpr, gp, lightning.Name, gs, os)
	} else {
		ret := &sgc7game.Result{
			Symbol:     lightning.SymbolCode,
			Type:       100,
			Mul:        1,
			CoinWin:    lastCollector + val*mul,
			CashWin:    (lastCollector + val*mul) * int(stake.CoinBet),
			Pos:        arrpos,
			Wilds:      0,
			SymbolNums: len(arrpos) / 2,
		}

		lightning.AddResult(curpr, ret, &cd.BasicComponentData)

		if lightning.Config.EndingFirstComponent != "" {
			gameProp.Respin(curpr, gp, lightning.Config.EndingFirstComponent, gs, os)
		} else {
			gameProp.SetStrVal(GamePropNextComponent, lightning.Config.DefaultNextComponent)
		}
	}

	// gp.AddComponentData(lightning.Name, cd)

	return nil
}

// OnAsciiGame - outpur to asciigame
func (lightning *Lightning) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[lightning.Name].(*LightningData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("respin symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)

		fmt.Printf("ligntning val-%v mul-%v lastCollector-%v\n", cd.Val, cd.Mul, cd.Collector)

		if cd.NewConnector > 0 {
			fmt.Printf("new collect is %v\n", cd.NewConnector)
		}

		asciigame.OutputResults(fmt.Sprintf("%v wins", lightning.Name), pr, func(i int, ret *sgc7game.Result) bool {
			return goutils.IndexOfIntSlice(cd.UsedResults, i, 0) >= 0
		}, mapSymbolColor)
	}

	return nil
}

// OnStats
func (lightning *Lightning) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (lightning *Lightning) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData *anypb.Any, pr *sgc7game.PlayResult) (int64, error) {
	pbcd := &sgc7pb.LightningData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Lightning.OnStatsWithPB:UnmarshalTo",
			zap.Error(err))

		return 0, err
	}

	return lightning.OnStatsWithPBBasicComponentData(feature, pbcd.BasicComponentData, pr), nil
}

// NewComponentData -
func (lightning *Lightning) NewComponentData() IComponentData {
	return &LightningData{}
}

// EachUsedResults -
func (lightning *Lightning) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.LightningData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("Lightning.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewLightning(name string) IComponent {
	return &Lightning{
		BasicComponent:           NewBasicComponent(name),
		MapSymbols:               make(map[int]*LightningSymbolData),
		MapSymbolTriggerFeatures: make(map[int]*LightningTriggerFeatureConfig),
	}
}
