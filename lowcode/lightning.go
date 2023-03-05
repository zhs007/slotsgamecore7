package lowcode

import (
	"fmt"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const (
	LightningTypeVal = "val"
	LightningTypeMul = "mul"

	LightningFeatureCollector = "collector"
)

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
	Collector                int
	Val                      int
	Mul                      int
	NewConnector             int
}

// Init -
func (lightning *Lightning) Init(fn string, gameProp *GameProperty) error {
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

	lightning.SymbolCode = gameProp.CurPaytables.MapSymbols[cfg.Symbol]

	if lightning.Config.Weight != "" {
		vw2, err := sgc7game.LoadValWeights2FromExcelWithSymbols(lightning.Config.Weight, "val", "weight", gameProp.CurPaytables)
		if err != nil {
			goutils.Error("Lightning.Init:LoadValWeights2FromExcelWithSymbols",
				zap.String("Weight", lightning.Config.Weight),
				zap.Error(err))

			return err
		}

		lightning.Weight = vw2
	}

	for _, v := range lightning.Config.SymbolVals {
		symbolCode := gameProp.CurPaytables.MapSymbols[v.Symbol]

		sd := &LightningSymbolData{
			SymbolCode: symbolCode,
			Config:     v,
		}

		vw2, err := sgc7game.LoadValWeights2FromExcel(v.Weight, "val", "weight", sgc7game.NewIntVal[int])
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
		symbolCode := gameProp.CurPaytables.MapSymbols[v.Symbol]

		lightning.MapSymbolTriggerFeatures[symbolCode] = v

		if v.Feature == LightningFeatureCollector {
			lightning.CollectorSymbolCode = symbolCode
		}
	}

	lightning.BasicComponent.onInit(&cfg.BasicComponentConfig)

	return nil
}

// OnNewGame -
func (lightning *Lightning) OnNewGame(gameProp *GameProperty) error {
	return nil
}

// OnNewStep -
func (lightning *Lightning) OnNewStep(gameProp *GameProperty) error {

	lightning.BasicComponent.OnNewStep()

	lightning.Collector = 0
	lightning.Val = 0
	lightning.Mul = 0
	lightning.NewConnector = 0

	return nil
}

// playgame
func (lightning *Lightning) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

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

	gs := pregp.LastScene.Clone()
	os := pregp.LastOtherScene.Clone()

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

	lightning.AddScene(gameProp, curpr, gs)
	lightning.AddOtherScene(gameProp, curpr, os)

	lightning.Collector = lastCollector
	lightning.Val = val
	lightning.Mul = mul

	if isCollector {
		lightning.NewConnector = lastCollector + val*mul

		os.Arr[collectorX][collectorY] = lightning.NewConnector

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

		lightning.AddResult(curpr, ret)

		if lightning.Config.EndingFirstComponent != "" {
			gameProp.Respin(curpr, gp, lightning.Config.EndingFirstComponent, gs, os)
		} else {
			gameProp.SetStrVal(GamePropNextComponent, lightning.Config.DefaultNextComponent)
		}
	}

	return nil
}

// OnAsciiGame - outpur to asciigame
func (lightning *Lightning) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	if len(lightning.UsedScenes) > 0 {
		// fmt.Printf("mystery is %v\n", gameProp.GetStrVal(GamePropCurMystery))
		asciigame.OutputScene("respin symbols", pr.Scenes[lightning.UsedScenes[0]], mapSymbolColor)

		fmt.Printf("ligntning val-%v mul-%v lastCollector-%v\n", lightning.Val, lightning.Mul, lightning.Collector)

		if lightning.NewConnector > 0 {
			fmt.Printf("new collect is %v\n", lightning.NewConnector)
		}

		asciigame.OutputResults(fmt.Sprintf("%v wins", lightning.Name), pr, func(i int, ret *sgc7game.Result) bool {
			return goutils.IndexOfIntSlice(lightning.UsedResults, i, 0) >= 0
		}, mapSymbolColor)
	}

	return nil
}

// OnStats
func (lightning *Lightning) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewLightning(name string) IComponent {
	return &Lightning{
		BasicComponent:           NewBasicComponent(name),
		MapSymbols:               make(map[int]*LightningSymbolData),
		MapSymbolTriggerFeatures: make(map[int]*LightningTriggerFeatureConfig),
	}
}
