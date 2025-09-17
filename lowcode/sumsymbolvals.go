package lowcode

import (
	"log/slog"
	"os"
	"strings"

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

// SumSymbolValsTypeName is the component type name used to register
// and identify SumSymbolVals components in configuration and component mgr.
const SumSymbolValsTypeName = "sumSymbolVals"

// SumSymbolValsType represents the comparison type used when deciding
// whether a symbol value should be counted in the sum. It covers
// equality, inequality and range membership forms.
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

// parseSumSymbolValsType parses a textual representation of the comparison
// type and returns the corresponding SumSymbolValsType. The input is expected
// to be one of the supported literal forms, e.g. "==", ">=", "in [min, max]".
// If the input is not recognized, SSVTypeNone is returned.
func parseSumSymbolValsType(strType string) SumSymbolValsType {
	switch strType {
	case "==":
		return SSVTypeEqu
	case ">=":
		return SSVTypeGreaterEqu
	case "<=":
		return SSVTypeLessEqu
	case ">":
		return SSVTypeGreater
	case "<":
		return SSVTypeLess
	case "in [min, max]":
		return SSVTypeInAreaLR
	case "in (min, max]":
		return SSVTypeInAreaR
	case "in [min, max)":
		return SSVTypeInAreaL
	case "in (min, max)":
		return SSVTypeInArea
	}

	return SSVTypeNone
}

// SumSymbolValsData holds runtime data for a SumSymbolVals component.
// It embeds BasicComponentData and stores the computed Number (the sum).
type SumSymbolValsData struct {
	BasicComponentData
	Number int
}

// OnNewGame implements the IComponentData lifecycle hook. It forwards
// to BasicComponentData.OnNewGame to initialize embedded fields.
func (sumSymbolValsData *SumSymbolValsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	sumSymbolValsData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone creates a deep copy of SumSymbolValsData for use in new play
// result contexts. This is required by the component framework to keep
// component data isolated between plays.
func (sumSymbolValsData *SumSymbolValsData) Clone() IComponentData {
	target := &SumSymbolValsData{
		BasicComponentData: sumSymbolValsData.CloneBasicComponentData(),
		Number:             sumSymbolValsData.Number,
	}

	return target
}

// BuildPBComponentData converts the component data into its protobuf
// representation for serialization across the system.
func (sumSymbolValsData *SumSymbolValsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.SumSymbolValsData{
		BasicComponentData: sumSymbolValsData.BuildPBBasicComponentData(),
		Number:             int32(sumSymbolValsData.Number),
	}

	return pbcd
}

// GetValEx returns named integer outputs exposed by this component's data.
// Supported keys are CVNumber and CVOutputInt, both returning the computed
// Number (sum) value. The boolean indicates whether the key was handled.
func (sumSymbolValsData *SumSymbolValsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVNumber || key == CVOutputInt {
		return sumSymbolValsData.Number, true
	}

	return 0, false
}

// GetOutput returns the numeric output value for ascii or other renderers.
func (rollNumberData *SumSymbolValsData) GetOutput() int {
	return rollNumberData.Number
}

// SumSymbolValsConfig is the YAML/JSON configuration structure for
// a SumSymbolVals component. It includes the comparison type (StrType),
// numeric parameters (Value, Min, Max), the source component name to
// read positions from, and optional awards to trigger when evaluated.
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

// SetLinkComponent implements BasicComponentConfig's link wiring helper.
// Currently SumSymbolVals supports linking the "next" component which
// sets DefaultNextComponent used by the framework.
func (cfg *SumSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// SumSymbolVals is a component that sums symbol values from a scene
// according to a configured predicate. The result is stored in
// SumSymbolValsData.Number and can be used by downstream components.
type SumSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *SumSymbolValsConfig `json:"config"`
}

// Init loads configuration from a YAML file and delegates to InitEx.
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

// InitEx initializes the component from an already-unmarshaled config
// object. This method performs type assertion and sets up internal
// fields such as parsed Type and initializes any Awards.
func (sumSymbolVals *SumSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	cd, ok := cfg.(*SumSymbolValsConfig)
    	if !ok {
    		goutils.Error("SumSymbolVals.InitEx:invalid cfg type",
    			slog.Any("cfg", cfg))
    		return ErrInvalidComponentConfig
    	}
	sumSymbolVals.Config = cd
	sumSymbolVals.Config.ComponentType = SumSymbolValsTypeName

	sumSymbolVals.Config.Type = parseSumSymbolValsType(sumSymbolVals.Config.StrType)

	for _, v := range sumSymbolVals.Config.Awards {
		v.Init()
	}

	sumSymbolVals.onInit(&sumSymbolVals.Config.BasicComponentConfig)

	return nil
}

// ProcControllers triggers any configured awards/controllers after the
// component has evaluated. It delegates to gameProp.procAwards if awards
// are present in configuration.
func (sumSymbolVals *SumSymbolVals) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(sumSymbolVals.Config.Awards) > 0 {
		gameProp.procAwards(plugin, sumSymbolVals.Config.Awards, curpr, gp)
	}
}

// checkVal evaluates whether a single symbol value v satisfies the
// configured predicate (Type/Value/Min/Max). It returns true when v
// should be included in the sum.
func (sumSymbolVals *SumSymbolVals) checkVal(v int) bool {
	switch sumSymbolVals.Config.Type {
	case SSVTypeEqu:
		return v == sumSymbolVals.Config.Value
	case SSVTypeGreaterEqu:
		return v >= sumSymbolVals.Config.Value
	case SSVTypeLessEqu:
		return v <= sumSymbolVals.Config.Value
	case SSVTypeGreater:
		return v > sumSymbolVals.Config.Value
	case SSVTypeLess:
		return v < sumSymbolVals.Config.Value
	case SSVTypeInAreaLR:
		return v >= sumSymbolVals.Config.Min && v <= sumSymbolVals.Config.Max
	case SSVTypeInAreaR:
		return v > sumSymbolVals.Config.Min && v <= sumSymbolVals.Config.Max
	case SSVTypeInAreaL:
		return v >= sumSymbolVals.Config.Min && v < sumSymbolVals.Config.Max
	case SSVTypeInArea:
		return v > sumSymbolVals.Config.Min && v < sumSymbolVals.Config.Max
	}
	// unknown type
	goutils.Debug("SumSymbolVals.checkVal: unknown type",
		slog.Int("type", int(sumSymbolVals.Config.Type)))

	return false
}

func (sumSymbolVals *SumSymbolVals) sum(gameProp *GameProperty, os *sgc7game.GameScene) int {
	sumVal := 0

	// If a source component is configured, use its positions; otherwise
	// iterate the whole scene.
	pc, isok := gameProp.Components.MapComponents[sumSymbolVals.Config.SourceComponent]
	if isok {
		pccd := gameProp.GetComponentData(pc)
		pos := pccd.GetPos()

		// pos is expected to be an even-length slice of coordinates [x0,y0,x1,y1,...].
		// If pos is non-empty we iterate pairs; otherwise we fallback to full scan.
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

	cd, ok := icd.(*SumSymbolValsData)
	if !ok {
		goutils.Error("SumSymbolVals.OnPlayGame: invalid component data type",
			slog.Any("icd", icd))
		return "", ErrInvalidComponentData
	}
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
		StrType:         strings.ToLower(jcfg.Type),
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
