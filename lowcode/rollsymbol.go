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
	"gopkg.in/yaml.v2"
)

// RollSymbolTypeName is the component type name used to identify
// roll symbol components in configuration and during initialization.
const RollSymbolTypeName = "rollSymbol"

// RollSymbolData holds the runtime data for a RollSymbol component.
// It embeds BasicComponentData to inherit common component state and
// includes the list of generated SymbolCodes for the current play.
type RollSymbolData struct {
	BasicComponentData
	// SymbolCodes is the sequence of symbol integer codes produced
	// by the component during OnPlayGame.
	SymbolCodes []int
}

// OnNewGame resets or initializes runtime data when a new game starts.
// It forwards to the embedded BasicComponentData.OnNewGame implementation.
func (rollSymbolData *RollSymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	rollSymbolData.BasicComponentData.OnNewGame(gameProp, component)
}

// Clone creates a deep copy of RollSymbolData. It is used when the
// component system needs a fresh copy of component data (for example
// between parallel plays or for snapshotting state).
func (rollSymbolData *RollSymbolData) Clone() IComponentData {
	target := &RollSymbolData{
		BasicComponentData: rollSymbolData.CloneBasicComponentData(),
	}

	target.SymbolCodes = make([]int, len(rollSymbolData.SymbolCodes))
	copy(target.SymbolCodes, rollSymbolData.SymbolCodes)

	return target
}

// BuildPBComponentData converts runtime data into the protobuf
// representation used for serialization or inter-service communication.
func (rollSymbolData *RollSymbolData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RollSymbolData{
		BasicComponentData: rollSymbolData.BuildPBBasicComponentData(),
	}

	for _, v := range rollSymbolData.SymbolCodes {
		pbcd.SymbolCodes = append(pbcd.SymbolCodes, int32(v))
	}

	return pbcd
}

// GetValEx returns extended integer values for this component data
// identified by key. RollSymbol currently does not expose any
// extended values, so it returns false.
func (rollSymbolData *RollSymbolData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	return 0, false
}

// RollSymbolConfig defines the static configuration for a RollSymbol
// component. Fields are populated from YAML/JSON configuration and
// consumed during InitEx and OnPlayGame.
type RollSymbolConfig struct {
	// BasicComponentConfig contains common config fields such as Name
	// and link information used by BasicComponent logic.
	BasicComponentConfig   `yaml:",inline" json:",inline"`

	// SymbolNum is the default number of symbols to roll/generate.
	SymbolNum              int                   `yaml:"symbolNum" json:"symbolNum"`

	// SymbolNumComponent optionally points to another component whose
	// output overrides SymbolNum at runtime.
	SymbolNumComponent     string                `json:"symbolNumComponent"`

	// Weight identifies the symbol weight resource used to randomly
	// choose symbols (referencing a ValWeights2 configuration).
	Weight                 string                `yaml:"weight" json:"weight"`
	// WeightVW is the loaded weight data (populated in InitEx).
	WeightVW               *sgc7game.ValWeights2 `json:"-"`

	// SrcSymbolCollection allows limiting weights to symbols present
	// in a given collection at runtime.
	SrcSymbolCollection    string                `yaml:"srcSymbolCollection" json:"srcSymbolCollection"`

	// IgnoreSymbolCollection lists symbols to exclude from the weight
	// selection.
	IgnoreSymbolCollection string                `yaml:"ignoreSymbolCollection" json:"ignoreSymbolCollection"`

	// TargetSymbolCollection, if set, receives the generated symbols
	// so other components can use them.
	TargetSymbolCollection string                `yaml:"targetSymbolCollection" json:"targetSymbolCollection"`
}

// SetLinkComponent
// SetLinkComponent sets linked component names (currently only "next"
// is supported and is stored in DefaultNextComponent).
func (cfg *RollSymbolConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

// RollSymbol is the runtime component implementation. It composes
// BasicComponent for common behavior and holds a typed Config pointer.
type RollSymbol struct {
	*BasicComponent `json:"-"`
	Config          *RollSymbolConfig `json:"config"`
}

// Init loads the component configuration from a YAML file and calls
// InitEx with the parsed configuration. The filename should be the
// path to a YAML file describing RollSymbolConfig.
func (rollSymbol *RollSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RollSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RollSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RollSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return rollSymbol.InitEx(cfg, pool)
}

// InitEx initializes the RollSymbol from an already-parsed config
// object (typically produced by parsing YAML/JSON). It validates the
// cfg type, loads weight resources, and performs any additional
// configuration initialization.
func (rollSymbol *RollSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	rcfg, ok := cfg.(*RollSymbolConfig)
	if !ok {
		goutils.Error("RollSymbol.InitEx:invalid cfg type",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	rollSymbol.Config = rcfg
	rollSymbol.Config.ComponentType = RollSymbolTypeName

	// Load the symbol weights referenced by the configuration. A
	// missing Weight is considered invalid for this component.
	if rollSymbol.Config.Weight != "" {
		vw2, err := pool.LoadSymbolWeights(rollSymbol.Config.Weight, "val", "weight", pool.DefaultPaytables, rollSymbol.Config.UseFileMapping)
		if err != nil {
			goutils.Error("RollSymbol.Init:LoadStrWeights",
				slog.String("Weight", rollSymbol.Config.Weight),
				goutils.Err(err))

			return err
		}

		rollSymbol.Config.WeightVW = vw2
	} else {
		goutils.Error("RollSymbol.InitEx:Weight",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	// Call common BasicComponent initialization logic (links, names, etc.).
	rollSymbol.onInit(&rollSymbol.Config.BasicComponentConfig)

	return nil
}

func (rollSymbol *RollSymbol) getValWeight(gameProp *GameProperty) *sgc7game.ValWeights2 {
	// If no collection filtering is configured, return the loaded weight
	// directly to avoid unnecessary cloning.
	if rollSymbol.Config.SrcSymbolCollection == "" && rollSymbol.Config.IgnoreSymbolCollection == "" {
		return rollSymbol.Config.WeightVW
	}

	var vw *sgc7game.ValWeights2

	// If a source symbol collection is specified, clone the weight set
	// but limit values only to symbols present in that collection.
	if rollSymbol.Config.SrcSymbolCollection != "" {
		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.SrcSymbolCollection)

		vw = rollSymbol.Config.WeightVW.CloneWithIntArray(symbols)

		if vw == nil {
			return nil
		}
	}

	// If cloning from source collection wasn't done, clone the full set
	// so that subsequent modifications don't alter the original.
	if vw == nil {
		vw = rollSymbol.Config.WeightVW.Clone()
	}

	// Remove ignored symbols from the working ValWeights2 clone.
	if rollSymbol.Config.IgnoreSymbolCollection != "" {
		symbols := gameProp.GetComponentSymbols(rollSymbol.Config.IgnoreSymbolCollection)

		if len(symbols) > 0 {
			vw = vw.CloneWithoutIntArray(symbols)
		}

		if vw == nil {
			return nil
		}
	}

	if len(vw.Vals) == 0 {
		return nil
	}

	return vw
}

func (rollSymbol *RollSymbol) getSymbolNum(gameProp *GameProperty, basicCD *BasicComponentData) int {
	v, isok := basicCD.GetConfigIntVal(CCVSymbolNum)
	if isok {
		return v
	}

	if rollSymbol.Config.SymbolNumComponent != "" {
		cd := gameProp.GetComponentDataWithName(rollSymbol.Config.SymbolNumComponent)
		if cd != nil {
			return cd.GetOutput()
		}
	}

	// Default to the configured SymbolNum value.
	return rollSymbol.Config.SymbolNum
}

// OnPlayGame is the core execution method invoked when the component
// should produce symbol rolls for a single play. It populates
// RollSymbolData.SymbolCodes and optionally registers generated
// symbols into a target collection for other components to consume.
func (rollSymbol *RollSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// Ensure component data is of expected type before usage.
	rsd, ok := icd.(*RollSymbolData)
	if !ok {
		goutils.Error("RollSymbol.OnPlayGame:invalid component data",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	rsd.SymbolCodes = nil

	// Resolve how many symbols to generate for this play. The number
	// can be overridden at runtime by another component via
	// SymbolNumComponent or by a dynamic component-config value.
	sn := rollSymbol.getSymbolNum(gameProp, &rsd.BasicComponentData)

	for i := 0; i < sn; i++ {
		// Each roll may depend on the current state of symbol
		// collections, so obtain a working ValWeights2 instance which
		// may be a clone filtered by source/ignore collections.
		vw := rollSymbol.getValWeight(gameProp)
		if vw == nil {
			// No available weights -> stop producing symbols.
			break
		}

		// Randomly choose a value according to the weight set.
		cr, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("RollSymbol.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		sc := cr.Int()

		// Record generated symbol code in component data.
		rsd.SymbolCodes = append(rsd.SymbolCodes, sc)

		// Optionally export the generated symbol to a named
		// collection for other components to access.
		if rollSymbol.Config.TargetSymbolCollection != "" {
			gameProp.AddComponentSymbol(rollSymbol.Config.TargetSymbolCollection, sc)
		}
	}

	if len(rsd.SymbolCodes) == 0 {
		nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	nc := rollSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame outputs the component runtime result to the ASCII game
// printer (useful for debugging / CLI visualization). It prints the
// list of generated symbol names using the default paytable mapping.
func (rollSymbol *RollSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	rsd, ok := icd.(*RollSymbolData)
	if !ok {
		goutils.Error("RollSymbol.OnAsciiGame:invalid component data",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	fmt.Printf("rollSymbol %v, got ", rollSymbol.GetName())

	for _, v := range rsd.SymbolCodes {
		// Translate symbol code into a string using the default
		// paytables for human-friendly display.
		fmt.Printf("%v ", gameProp.Pool.DefaultPaytables.GetStringFromInt(v))
	}

	fmt.Print("\n")

	return nil
}

// NewComponentData -
func (rollSymbol *RollSymbol) NewComponentData() IComponentData {
	return &RollSymbolData{}
}

func NewRollSymbol(name string) IComponent {
	return &RollSymbol{
		BasicComponent: NewBasicComponent(name, 0),
	}
}

//	"configuration": {
//		"weight": "fgbookofsymbol",
//		"symbolNum": 3,
//	    "symbolNumComponent": "bg-symnum",
//		"ignoreSymbolCollection": "fg-syms",
//		"targetSymbolCollection": "fg-syms"
//	},
type jsonRollSymbol struct {
	Weight                 string `json:"weight"`
	SymbolNum              int    `json:"symbolNum"`
	SymbolNumComponent     string `json:"symbolNumComponent"`
	SrcSymbolCollection    string `json:"srcSymbolCollection"`
	IgnoreSymbolCollection string `json:"ignoreSymbolCollection"`
	TargetSymbolCollection string `json:"targetSymbolCollection"`
}

func (jcfg *jsonRollSymbol) build() *RollSymbolConfig {
	cfg := &RollSymbolConfig{
		Weight:                 jcfg.Weight,
		SymbolNum:              jcfg.SymbolNum,
		SymbolNumComponent:     jcfg.SymbolNumComponent,
		SrcSymbolCollection:    jcfg.SrcSymbolCollection,
		IgnoreSymbolCollection: jcfg.IgnoreSymbolCollection,
		TargetSymbolCollection: jcfg.TargetSymbolCollection,
	}

	return cfg
}

func parseRollSymbol(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRollSymbol:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRollSymbol:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRollSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRollSymbol:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RollSymbolTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
