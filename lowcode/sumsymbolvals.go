package lowcode

import (
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
	"gopkg.in/yaml.v2"
)

const SumSymbolValsTypeName = "sumSymbolVals"

type SumSymbolValsType int

const (
	SSVTypeNone       SumSymbolValsType = 0 // none
	SSVTypeEqu        SumSymbolValsType = 1 // ==
	SSVTypeGreaterEqu SumSymbolValsType = 2 // >=
	SSVTypeLessEqu    SumSymbolValsType = 3 // <=
	SSVTypeGreater    SumSymbolValsType = 4 // >
	SSVTypeLess       SumSymbolValsType = 5 // <
	SSVTypeInAreaLR   SumSymbolValsType = 6 // In [min, max]
	SSVTypeInAreaR    SumSymbolValsType = 7 // In (min, max]
	SSVTypeInAreaL    SumSymbolValsType = 8 // In [min, max)
	SSVTypeInArea     SumSymbolValsType = 9 // In (min, max)
)

func parseSumSymbolValsType(strType string) SumSymbolValsType {
	if strType == "==" {
		return SSVTypeEqu
	} else if strType == ">=" {
		return SSVTypeGreaterEqu
	} else if strType == "<=" {
		return SSVTypeLessEqu
	} else if strType == ">" {
		return SSVTypeGreater
	} else if strType == "<" {
		return SSVTypeLess
	} else if strType == "In [min, max]" {
		return SSVTypeInAreaLR
	} else if strType == "In (min, max]" {
		return SSVTypeInAreaR
	} else if strType == "In [min, max)" {
		return SSVTypeInAreaL
	} else if strType == "In (min, max)" {
		return SSVTypeInArea
	}

	return SSVTypeNone
}

type SumSymbolValsData struct {
	BasicComponentData
	Number int
}

// OnNewGame -
func (sumSymbolValsData *SumSymbolValsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	sumSymbolValsData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone
func (sumSymbolValsData *SumSymbolValsData) Clone() IComponentData {
	target := &SumSymbolValsData{
		BasicComponentData: sumSymbolValsData.CloneBasicComponentData(),
		Number:             sumSymbolValsData.Number,
	}

	return target
}

// BuildPBComponentData
func (sumSymbolValsData *SumSymbolValsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SumSymbolValsData{
		BasicComponentData: sumSymbolValsData.BuildPBBasicComponentData(),
		Number:             int32(sumSymbolValsData.Number),
	}

	return pbcd
}

// GetValEx -
func (sumSymbolValsData *SumSymbolValsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVNumber || key == CVOutputInt {
		return sumSymbolValsData.Number, true
	}

	return 0, false
}

// // SetConfigIntVal - CCVValueNum的set和chg逻辑不太一样，等于的时候不会触发任何的 controllers
// func (sumSymbolValsData *SumSymbolValsData) SetConfigIntVal(key string, val int) {
// 	if key == CCVForceValNow {
// 		sumSymbolValsData.Number = val
// 	} else {
// 		sumSymbolValsData.BasicComponentData.SetConfigIntVal(key, val)
// 	}
// }

// // ChgConfigIntVal -
// func (sumSymbolValsData *SumSymbolValsData) ChgConfigIntVal(key string, off int) int {
// 	if key == CCVForceValNow {
// 		sumSymbolValsData.Number += off

// 		return sumSymbolValsData.Number
// 	}

// 	return sumSymbolValsData.BasicComponentData.ChgConfigIntVal(key, off)
// }

// GetOutput -
func (rollNumberData *SumSymbolValsData) GetOutput() int {
	return rollNumberData.Number
}

// SumSymbolValsConfig - configuration for SumSymbolVals
type SumSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string            `yaml:"type" json:"type"`
	Type                 SumSymbolValsType `yaml:"-" json:"-"`
	Value                int               `yaml:"value" json:"value"`
	Min                  int               `yaml:"min" json:"min"`
	Max                  int               `yaml:"max" json:"max"`
	SourceComponent      string            `yaml:"sourceComponent" json:"sourceComponent"`
	Awards               []*Award          `yaml:"awards" json:"awards"` // 新的奖励系统
}

// SetLinkComponent
func (cfg *SumSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type SumSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *SumSymbolValsConfig `json:"config"`
}

// Init -
func (sumSymbolVals *SumSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("SumSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &SumSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("SumSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return sumSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (sumSymbolVals *SumSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	sumSymbolVals.Config = cfg.(*SumSymbolValsConfig)
	sumSymbolVals.Config.ComponentType = CheckSymbolValsTypeName

	sumSymbolVals.Config.Type = parseSumSymbolValsType(sumSymbolVals.Config.StrType)

	for _, v := range sumSymbolVals.Config.Awards {
		v.Init()
	}

	sumSymbolVals.onInit(&sumSymbolVals.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (sumSymbolVals *SumSymbolVals) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(sumSymbolVals.Config.Awards) > 0 {
		gameProp.procAwards(plugin, sumSymbolVals.Config.Awards, curpr, gp)
	}
}

func (sumSymbolVals *SumSymbolVals) checkVal(v int) bool {
	if sumSymbolVals.Config.Type == SSVTypeEqu {
		return v == sumSymbolVals.Config.Value
	} else if sumSymbolVals.Config.Type == SSVTypeGreaterEqu {
		return v >= sumSymbolVals.Config.Value
	} else if sumSymbolVals.Config.Type == SSVTypeLessEqu {
		return v <= sumSymbolVals.Config.Value
	} else if sumSymbolVals.Config.Type == SSVTypeGreater {
		return v > sumSymbolVals.Config.Value
	} else if sumSymbolVals.Config.Type == SSVTypeLess {
		return v < sumSymbolVals.Config.Value
	} else if sumSymbolVals.Config.Type == SSVTypeInAreaLR {
		return v >= sumSymbolVals.Config.Min && v <= sumSymbolVals.Config.Max
	} else if sumSymbolVals.Config.Type == SSVTypeInAreaR {
		return v > sumSymbolVals.Config.Min && v <= sumSymbolVals.Config.Max
	} else if sumSymbolVals.Config.Type == SSVTypeInAreaL {
		return v >= sumSymbolVals.Config.Min && v < sumSymbolVals.Config.Max
	} else if sumSymbolVals.Config.Type == SSVTypeInArea {
		return v > sumSymbolVals.Config.Min && v < sumSymbolVals.Config.Max
	}

	return false
}

func (sumSymbolVals *SumSymbolVals) sum(gameProp *GameProperty, os *sgc7game.GameScene) int {
	sumVal := 0

	pc, isok := gameProp.Components.MapComponents[sumSymbolVals.Config.SourceComponent]
	if isok {
		pccd := gameProp.GetComponentData(pc)
		pos := pccd.GetPos()

		if len(pos) > 0 {
			for i := 0; i < len(pos)/2; i++ {
				x := pos[i*2]
				y := pos[i*2+1]
				v := os.Arr[x][y]

				if sumSymbolVals.checkVal(v) {
					sumVal += v
				}
			}
		} else {
			for _, arr := range os.Arr {
				for _, v := range arr {
					if sumSymbolVals.checkVal(v) {
						sumVal += v
					}
				}
			}
		}
	} else {
		for _, arr := range os.Arr {
			for _, v := range arr {
				if sumSymbolVals.checkVal(v) {
					sumVal += v
				}
			}
		}
	}

	return sumVal
}

// playgame
func (sumSymbolVals *SumSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*SumSymbolValsData)
	cd.Number = 0

	os := sumSymbolVals.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	if os != nil {
		val := sumSymbolVals.sum(gameProp, os)
		cd.Number = val

		sumSymbolVals.ProcControllers(gameProp, plugin, curpr, gp, -1, "")
	}

	nc := sumSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (sumSymbolVals *SumSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewComponentData -
func (sumSymbolVals *SumSymbolVals) NewComponentData() IComponentData {
	return &SumSymbolValsData{}
}

func NewSumSymbolVals(name string) IComponent {
	return &SumSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "type": ">",
// "value": 0,
// "outputToComponent": "bg-pos-rmoved",
// "sourceComponent": "bg-pos-rmoved"
type jsonSumSymbolVals struct {
	Type            string `json:"type"`
	Value           int    `json:"value"`
	Min             int    `json:"min"`
	Max             int    `json:"max"`
	SourceComponent string `json:"sourceComponent"`
}

func (jcfg *jsonSumSymbolVals) build() *SumSymbolValsConfig {
	cfg := &SumSymbolValsConfig{
		StrType: jcfg.Type,
		// OutputToComponent: jcfg.OutputToComponent,
		Value:           jcfg.Value,
		Min:             jcfg.Min,
		Max:             jcfg.Max,
		SourceComponent: jcfg.SourceComponent,
	}

	return cfg
}

func parseSumSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseSumSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseSumSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonSumSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseSumSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseSumSymbolVals:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: SumSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
