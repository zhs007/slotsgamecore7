package lowcode

import (
	"log/slog"
	"os"
	"slices"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

// GenSymbolValsInReelsTypeName - component name
const GenSymbolValsInReelsTypeName = "genSymbolValsInReels"

type GenSymbolValsInReelsSrcSymbolValsType int

const (
	GSVIRSSVTypeNone  GenSymbolValsInReelsSrcSymbolValsType = 0
	GSVIRSSVTypeClone GenSymbolValsInReelsSrcSymbolValsType = 1
)

func parseGenSymbolValsInReelsSrcSymbolValsType(str string) GenSymbolValsInReelsSrcSymbolValsType {
	switch str {
	case "clone":
		return GSVIRSSVTypeClone
	}

	return GSVIRSSVTypeNone
}

type GenSymbolValsInReelsCoreType int

const (
	GSVIRCTypeNumber GenSymbolValsInReelsCoreType = 0
	GSVIRCTypeWeight GenSymbolValsInReelsCoreType = 1
)

func parseGenSymbolValsInReelsCoreType(str string) GenSymbolValsInReelsCoreType {
	switch str {
	case "weight":
		return GSVIRCTypeWeight
	}

	return GSVIRCTypeNumber
}

// GenSymbolValsInReelsConfig - minimal configuration for the placeholder
type GenSymbolValsInReelsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrSrcSymbolValsType string                                `yaml:"srcSymbolValsType" json:"srcSymbolValsType"`
	SrcSymbolValsType    GenSymbolValsInReelsSrcSymbolValsType `yaml:"-" json:"-"`
	DefaultVal           int                                   `yaml:"defaultVal" json:"defaultVal"`
	IsForceRefresh       bool                                  `yaml:"isForceRefresh" json:"isForceRefresh"`
	StrGenType           string                                `yaml:"genType" json:"genType"`
	GenType              GenSymbolValsInReelsCoreType          `yaml:"-" json:"-"`
	IsAlwaysGen          bool                                  `yaml:"isAlwaysGen" json:"isAlwaysGen"`
	MapSrcSymbols        map[int][]string                      `yaml:"mapSrcSymbols" json:"mapSrcSymbols"`
	MapSrcSymbolCodes    map[int][]int                         `yaml:"-" json:"-"`
	MapNumber            map[int]int                           `yaml:"mapNumber" json:"mapNumber"`
	MapNumberWeight      map[int]string                        `yaml:"mapNumberWeight" json:"mapNumberWeight"`
	MapNumberWeightVW    map[int]*sgc7game.ValWeights2         `yaml:"-" json:"-"`
	SpGrid               string                                `yaml:"spGrid" json:"spGrid"`
	MapControllers       map[string][]*Award                   `yaml:"mapControllers" json:"mapControllers"`
	JumpToComponent      string                                `yaml:"jumpToComponent" json:"jumpToComponent"`
}

func (cfg *GenSymbolValsInReelsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type GenSymbolValsInReels struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolValsInReelsConfig `json:"config"`
}

// Init - read config from file
func (gsv *GenSymbolValsInReels) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolValsInReels.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolValsInReelsConfig{}
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolValsInReels.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return gsv.InitEx(cfg, pool)
}

// InitEx - placeholder init with parsed config
func (gsv *GenSymbolValsInReels) InitEx(cfg any, pool *GamePropertyPool) error {
	gsv.Config = cfg.(*GenSymbolValsInReelsConfig)
	gsv.Config.ComponentType = GenSymbolValsInReelsTypeName

	gsv.Config.SrcSymbolValsType = parseGenSymbolValsInReelsSrcSymbolValsType(gsv.Config.StrSrcSymbolValsType)
	gsv.Config.GenType = parseGenSymbolValsInReelsCoreType(gsv.Config.StrGenType)

	if gsv.Config.MapSrcSymbols != nil {
		gsv.Config.MapSrcSymbolCodes = make(map[int][]int, len(gsv.Config.MapSrcSymbols))

		for k, arr := range gsv.Config.MapSrcSymbols {
			for _, v := range arr {
				code, isok := pool.DefaultPaytables.MapSymbols[v]
				if isok {
					gsv.Config.MapSrcSymbolCodes[k] = append(gsv.Config.MapSrcSymbolCodes[k], code)
				}
			}
		}
	}

	if len(gsv.Config.MapNumberWeight) > 0 {
		gsv.Config.MapNumberWeightVW = make(map[int]*sgc7game.ValWeights2, len(gsv.Config.MapNumberWeight))

		for k, v := range gsv.Config.MapNumberWeight {
			vw2, err := pool.LoadIntWeights(v, gsv.Config.UseFileMapping)
			if err != nil {
				goutils.Error("GenSymbolValsInReels.InitEx:MapNumberWeight:LoadIntWeights",
					slog.String("Weight", v),
					goutils.Err(err))
				return err
			}

			gsv.Config.MapNumberWeightVW[k] = vw2
		}
	}

	for _, arr := range gsv.Config.MapControllers {
		for _, v := range arr {
			v.Init()
		}
	}

	gsv.onInit(&gsv.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (gsv *GenSymbolValsInReels) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	arr := gsv.Config.MapControllers[strVal]
	if len(arr) > 0 {
		gameProp.procAwards(plugin, arr, curpr, gp)
	}
}

// getSrcOtherScene
func (gsv *GenSymbolValsInReels) getSrcOtherScene(gameProp *GameProperty, curpr *sgc7game.PlayResult,
	prs []*sgc7game.PlayResult) (*sgc7game.GameScene, error) {

	if gsv.Config.SrcSymbolValsType == GSVIRSSVTypeNone {
		return nil, nil
	}

	os := gsv.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	return os, nil
}

// procNumber
func (gsv *GenSymbolValsInReels) procNumber(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, plugin sgc7plugin.IPlugin,
	curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, gp *GameParams, bcd *BasicComponentData, stackSPGrid *SPGridStack) (*sgc7game.GameScene, error) {

	if os == nil {
		os = gsv.newOutput(gameProp, curpr, prs, bcd, stackSPGrid)
	} else {
		os = os.CloneEx(gameProp.PoolScene)
	}

	isTrigger := false

	if len(gsv.Config.MapSrcSymbolCodes) > 0 {
		for x, arr := range gsv.Config.MapSrcSymbolCodes {
			num := gsv.Config.MapNumber[x]

			for y := 0; y < gs.Height; y++ {
				sym := gs.Arr[x][y]
				if slices.Contains(arr, sym) {
					os.Arr[x][y] = num

					isTrigger = true
				}
			}
		}
	} else {
		for x, num := range gsv.Config.MapNumber {
			for y := 0; y < gs.Height; y++ {
				os.Arr[x][y] = num

				isTrigger = true
			}
		}
	}

	if isTrigger {
		gsv.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")

		gsv.setOutput(gameProp, curpr, prs, bcd, os)

		return os, nil
	}

	return nil, nil
}

// procWeight
func (gsv *GenSymbolValsInReels) procWeight(gameProp *GameProperty, gs *sgc7game.GameScene, os *sgc7game.GameScene, plugin sgc7plugin.IPlugin,
	curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult,
	gp *GameParams, bcd *BasicComponentData, stackSPGrid *SPGridStack) (*sgc7game.GameScene, error) {

	if os == nil {
		os = gsv.newOutput(gameProp, curpr, prs, bcd, stackSPGrid)
	} else {
		os = os.CloneEx(gameProp.PoolScene)
	}

	isTrigger := false

	if len(gsv.Config.MapSrcSymbolCodes) > 0 {
		for x, arr := range gsv.Config.MapSrcSymbolCodes {
			vw := gsv.Config.MapNumberWeightVW[x]
			if vw != nil {
				for y := 0; y < gs.Height; y++ {
					sym := gs.Arr[x][y]
					if slices.Contains(arr, sym) {
						cr, err := vw.RandVal(plugin)
						if err != nil {
							goutils.Error("GenSymbolValsInReels.procWeight:RandVal",
								slog.Int("x", x),
								slog.Int("y", y),
								goutils.Err(err))

							return nil, err
						}

						os.Arr[x][y] = cr.Int()

						isTrigger = true
					}
				}
			}
		}
	} else {
		for x, vw := range gsv.Config.MapNumberWeightVW {
			if vw != nil {
				for y := 0; y < gs.Height; y++ {

					cr, err := vw.RandVal(plugin)
					if err != nil {
						goutils.Error("GenSymbolValsInReels.procWeight:RandVal",
							slog.Int("x", x),
							slog.Int("y", y),
							goutils.Err(err))

						return nil, err
					}

					os.Arr[x][y] = cr.Int()

					isTrigger = true
				}
			}
		}
	}

	if isTrigger {
		gsv.ProcControllers(gameProp, plugin, curpr, gp, 0, "<trigger>")

		gsv.setOutput(gameProp, curpr, prs, bcd, os)

		return os, nil
	}

	return nil, nil
}

func (gsv *GenSymbolValsInReels) getOutput(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult,
	_ *BasicComponentData) (*sgc7game.GameScene, *SPGridStack, error) {
	if gsv.Config.SpGrid != "" {
		stackSPGrid, isok := gameProp.MapSPGridStack[gsv.Config.SpGrid]
		if !isok {
			goutils.Error("GenSymbolValsInReels.getOutput:MapSPGridStack",
				slog.String("SPGrid", gsv.Config.SpGrid),
				goutils.Err(ErrInvalidComponentConfig))

			return nil, nil, ErrInvalidComponentConfig
		}

		spgrid := stackSPGrid.Stack.GetTopSPGridEx(gsv.Config.SpGrid, curpr, prs)

		return spgrid, stackSPGrid, nil
	}

	os, err := gsv.getSrcOtherScene(gameProp, curpr, prs)
	if err != nil {
		goutils.Error("GenSymbolValsInReels.getOutput:getSrcOtherScene",
			goutils.Err(err))

		return nil, nil, err
	}

	return os, nil, nil
}

func (gsv *GenSymbolValsInReels) newOutput(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult, bcd *BasicComponentData,
	stackSPGrid *SPGridStack) *sgc7game.GameScene {

	if stackSPGrid != nil {
		spgrid := gameProp.PoolScene.New2(stackSPGrid.Width, stackSPGrid.Height, gsv.Config.DefaultVal)

		return spgrid
	}

	os := gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight),
		gsv.Config.DefaultVal)

	return os
}

func (gsv *GenSymbolValsInReels) setOutput(gameProp *GameProperty, curpr *sgc7game.PlayResult, prs []*sgc7game.PlayResult,
	bcd *BasicComponentData, os *sgc7game.GameScene) {

	if gsv.Config.SpGrid != "" {
		gsv.AddSPGrid(gsv.Config.SpGrid, gameProp, curpr, os, bcd)
	} else {
		gsv.AddOtherScene(gameProp, curpr, os, bcd)
	}
}

// OnPlayGame - placeholder implementation; do nothing
func (gsv *GenSymbolValsInReels) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams,
	plugin sgc7plugin.IPlugin, cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult,
	icd IComponentData) (string, error) {

	// placeholder: nothing implemented yet
	bcd := icd.(*BasicComponentData)
	bcd.UsedOtherScenes = nil

	os, stackSPGrid, err := gsv.getOutput(gameProp, curpr, prs, bcd)
	if err != nil {
		goutils.Error("GenSymbolValsInReels.OnPlayGame:getOutput",
			goutils.Err(err))

		return "", err
	}

	gs := gsv.GetTargetScene3(gameProp, curpr, prs, 0)

	switch gsv.Config.GenType {
	case GSVIRCTypeNumber:
		nos, err := gsv.procNumber(gameProp, gs, os, plugin, curpr, prs, gp, bcd, stackSPGrid)
		if err != nil {
			goutils.Error("GenSymbolValsInReels.OnPlayGame:procNumber",
				goutils.Err(err))

			return "", err
		}

		if nos == nil {
			nc := gsv.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	case GSVIRCTypeWeight:
		nos, err := gsv.procWeight(gameProp, gs, os, plugin, curpr, prs, gp, bcd, stackSPGrid)
		if err != nil {
			goutils.Error("GenSymbolValsInReels.OnPlayGame:procWeight",
				goutils.Err(err))

			return "", err
		}

		if nos == nil {
			nc := gsv.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	default:
		goutils.Error("GenSymbolValsInReels.OnPlayGame:invalid GenType",
			slog.Int("GenType", int(gsv.Config.GenType)),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	nc := gsv.onStepEnd(gameProp, curpr, gp, gsv.Config.JumpToComponent)

	return nc, nil
}

// OnAsciiGame - placeholder - no output for ascii game
func (gsv *GenSymbolValsInReels) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

// NewGenSymbolValsInReels - factory
func NewGenSymbolValsInReels(name string) IComponent {
	return &GenSymbolValsInReels{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "srcSymbolValsType": "none",
// "defaultVal": 0,
// "isForceRefresh": false,
// "genType": "number",
// "isAlwaysGen": false,
// "mapSrcSymbols": [
//
//	[
//	    1,
//	    [
//	        "WL"
//	    ]
//	],
//	[
//	    2,
//	    [
//	        "WL"
//	    ]
//	],
//	[
//	    3,
//	    [
//	        "WL"
//	    ]
//	],
//	[
//	    4,
//	    [
//	        "WL"
//	    ]
//	],
//	[
//	    5,
//	    [
//	        "WL"
//	    ]
//	]
//
// ],
// "mapNumber": [
//
//	[
//	    1,
//	    1
//	],
//	[
//	    2,
//	    2
//	],
//	[
//	    3,
//	    3
//	],
//	[
//	    4,
//	    4
//	],
//	[
//	    5,
//	    5
//	]
//
// ]
// "spGrid": "fg-wl-spgrid"
type jsonGenSymbolValsInReels struct {
	SrcSymbolValsType string          `json:"srcSymbolValsType"`
	DefaultVal        int             `json:"defaultVal"`
	IsForceRefresh    bool            `json:"isForceRefresh"`
	GenType           string          `json:"genType"`
	IsAlwaysGen       bool            `json:"isAlwaysGen"`
	MapSrcSymbols     [][]interface{} `json:"mapSrcSymbols"`
	MapNumber         [][]interface{} `json:"mapNumber"`
	MapNumberWeight   [][]interface{} `json:"mapNumberWeight"`
	SpGrid            string          `json:"spGrid"`
}

func (jcfg *jsonGenSymbolValsInReels) build() *GenSymbolValsInReelsConfig {
	cfg := &GenSymbolValsInReelsConfig{
		StrSrcSymbolValsType: jcfg.SrcSymbolValsType,
		DefaultVal:           jcfg.DefaultVal,
		IsForceRefresh:       jcfg.IsForceRefresh,
		StrGenType:           jcfg.GenType,
		IsAlwaysGen:          jcfg.IsAlwaysGen,
		SpGrid:               jcfg.SpGrid,
	}

	if len(jcfg.MapNumber) > 0 {
		cfg.MapNumber = make(map[int]int, len(jcfg.MapNumber))

		for _, arr := range jcfg.MapNumber {
			if len(arr) != 2 {
				goutils.Error("jsonGenSymbolValsInReels.build:MapNumber:arr")

				return nil
			}

			k := arr[0].(float64)
			v := arr[1].(float64)
			cfg.MapNumber[int(k)-1] = int(v) // [1,w] => [0,w)
		}
	}

	if len(jcfg.MapSrcSymbols) > 0 {
		cfg.MapSrcSymbols = make(map[int][]string, len(jcfg.MapSrcSymbols))

		for _, arr := range jcfg.MapSrcSymbols {
			if len(arr) != 2 {
				goutils.Error("jsonGenSymbolValsInReels.build:MapSrcSymbols:arr")

				return nil
			}

			k := arr[0].(float64)
			arr1 := arr[1].([]interface{})
			for _, v := range arr1 {
				cfg.MapSrcSymbols[int(k)-1] = append(cfg.MapSrcSymbols[int(k)-1], v.(string))
			}
		}
	}

	if len(jcfg.MapNumberWeight) > 0 {
		cfg.MapNumberWeight = make(map[int]string, len(jcfg.MapNumberWeight))

		for _, arr := range jcfg.MapNumberWeight {
			if len(arr) != 2 {
				goutils.Error("jsonGenSymbolValsInReels.build:MapNumberWeight:arr")

				return nil
			}

			k := arr[0].(float64)
			v := arr[1].(string)
			cfg.MapNumberWeight[int(k)-1] = v // [1,w] => [0,w)
		}
	}

	return cfg
}

func parseGenSymbolValsInReels(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolValsInReels:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolValsInReels:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolValsInReels{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolValsInReels:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		mapControllers, err := parseAllAndStrMapControllers2(ctrls)
		if err != nil {
			goutils.Error("parseChgSymbols2:parseAllAndStrMapControllers2",
				goutils.Err(err))

			return "", err
		}

		cfgd.MapControllers = mapControllers
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolValsInReelsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
