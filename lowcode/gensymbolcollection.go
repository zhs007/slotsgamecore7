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
	"gopkg.in/yaml.v2"
)

const GenSymbolCollectionTypeName = "genSymbolCollection"

type GenSymbolCollectionCoreType int

const (
	GSCCoreTypeNumber       GenSymbolCollectionCoreType = 0
	GSCCoreTypeNumberWeight GenSymbolCollectionCoreType = 1
)

func parseGenSymbolCollectionCoreType(strType string) GenSymbolCollectionCoreType {
	if strType == "numberweight" {
		return GSCCoreTypeNumberWeight
	}

	return GSCCoreTypeNumber
}

// GenSymbolCollectionConfig - configuration for GenSymbolCollection
type GenSymbolCollectionConfig struct {
	BasicComponentConfig   `yaml:",inline" json:",inline"`
	StrCoreType            string                      `yaml:"coreType" json:"coreType"`
	CoreType               GenSymbolCollectionCoreType `yaml:"-" json:"-"`
	Number                 int                         `yaml:"number" json:"number"`
	NumberWeight           string                      `yaml:"numberWeight" json:"numberWeight"`
	NumberWeightVM         *sgc7game.ValWeights2       `yaml:"-" json:"-"`
	SymbolWeight           string                      `yaml:"symbolWeight" json:"symbolWeight"`
	SymbolWeightVM         *sgc7game.ValWeights2       `yaml:"-" json:"-"`
	OutputSymbolCollection string                      `yaml:"outputSymbolCollection" json:"outputSymbolCollection"`
}

// SetLinkComponent
func (cfg *GenSymbolCollectionConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolCollection struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolCollectionConfig `json:"config"`
}

// Init -
func (genSymbolCollection *GenSymbolCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolCollection.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolCollection.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genSymbolCollection.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolCollection *GenSymbolCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolCollection.Config = cfg.(*GenSymbolCollectionConfig)
	genSymbolCollection.Config.ComponentType = GenSymbolCollectionTypeName

	genSymbolCollection.Config.CoreType = parseGenSymbolCollectionCoreType(genSymbolCollection.Config.StrCoreType)

	if genSymbolCollection.Config.NumberWeight != "" {
		vw2, err := pool.LoadIntWeights(genSymbolCollection.Config.NumberWeight, true)
		if err != nil {
			goutils.Error("GenSymbolCollection.InitEx:LoadIntWeights",
				slog.String("NumberWeight", genSymbolCollection.Config.NumberWeight),
				goutils.Err(err))

			return err
		}

		genSymbolCollection.Config.NumberWeightVM = vw2
	} else if genSymbolCollection.Config.CoreType == GSCCoreTypeNumberWeight {
		goutils.Error("GenSymbolCollection.InitEx:NumberWeight",
			goutils.Err(ErrInvalidWeightVal))

		return ErrInvalidWeightVal
	}

	if genSymbolCollection.Config.SymbolWeight != "" {
		vw2, err := pool.LoadIntWeights(genSymbolCollection.Config.SymbolWeight, true)
		if err != nil {
			goutils.Error("GenSymbolCollection.InitEx:LoadIntWeights",
				slog.String("SymbolWeight", genSymbolCollection.Config.SymbolWeight),
				goutils.Err(err))

			return err
		}

		genSymbolCollection.Config.SymbolWeightVM = vw2
	} else {
		goutils.Error("GenSymbolCollection.InitEx:SymbolWeight",
			goutils.Err(ErrInvalidWeightVal))

		return ErrInvalidWeightVal
	}

	genSymbolCollection.onInit(&genSymbolCollection.Config.BasicComponentConfig)

	return nil
}

// getNumberWeight -
func (genSymbolCollection *GenSymbolCollection) getNumberWeight(gameProp *GameProperty, cd *BasicComponentData) *sgc7game.ValWeights2 {
	str := cd.GetConfigVal(CCVNumberWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return genSymbolCollection.Config.NumberWeightVM
}

// getSymbolWeight -
func (genSymbolCollection *GenSymbolCollection) getSymbolWeight(gameProp *GameProperty, cd *BasicComponentData) *sgc7game.ValWeights2 {
	str := cd.GetConfigVal(CCVSymbolWeight)
	if str != "" {
		vw2, _ := gameProp.Pool.LoadIntWeights(str, true)

		return vw2
	}

	return genSymbolCollection.Config.SymbolWeightVM
}

// playgame
func (genSymbolCollection *GenSymbolCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	bcd := icd.(*BasicComponentData)

	maxNumber := genSymbolCollection.Config.Number

	if genSymbolCollection.Config.CoreType == GSCCoreTypeNumberWeight {
		numberVM := genSymbolCollection.getNumberWeight(gameProp, bcd)
		cr, err := numberVM.RandVal(plugin)
		if err != nil {
			goutils.Error("GenSymbolCollection.OnPlayGame:RandVal:NumberWeightVM",
				goutils.Err(err))

			return "", err
		}

		maxNumber = cr.Int()
	}

	if genSymbolCollection.Config.OutputSymbolCollection == "" {
		goutils.Error("GenSymbolCollection.procPos:OutputSymbolCollection is empty",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	isc := gameProp.GetComponentDataWithName(genSymbolCollection.Config.OutputSymbolCollection)
	if isc == nil {
		goutils.Error("GenSymbolCollection.procPos:GetComponentDataWithName",
			slog.String("OutputSymbolCollection", genSymbolCollection.Config.OutputSymbolCollection),
			goutils.Err(ErrNoComponent))

		return "", ErrNoComponent
	}

	symVW := genSymbolCollection.getSymbolWeight(gameProp, bcd)

	for range maxNumber {
		cr, err := symVW.RandVal(plugin)
		if err != nil {
			goutils.Error("GenSymbolCollection.OnPlayGame:RandVal:SymbolWeightVM",
				goutils.Err(err))

			return "", err
		}

		isc.AddSymbolCodes([]int{cr.Int()})
	}

	nc := genSymbolCollection.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genSymbolCollection *GenSymbolCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewGenSymbolCollection(name string) IComponent {
	return &GenSymbolCollection{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "outputSymbolCollection": "bg-spsyms",
// "coreType": "numberWeight",
// "numberWeight": "bgspnumweight",
// "symbolWeight": "bgspsymweight"
type jsonGenSymbolCollection struct {
	CoreType               string `json:"coreType"`
	Number                 int    `json:"number"`
	NumberWeight           string `json:"numberWeight"`
	SymbolWeight           string `json:"symbolWeight"`
	OutputSymbolCollection string `json:"outputSymbolCollection"`
}

func (jcfg *jsonGenSymbolCollection) build() *GenSymbolCollectionConfig {
	cfg := &GenSymbolCollectionConfig{
		StrCoreType:            strings.ToLower(jcfg.CoreType),
		Number:                 jcfg.Number,
		NumberWeight:           jcfg.NumberWeight,
		SymbolWeight:           jcfg.SymbolWeight,
		OutputSymbolCollection: jcfg.OutputSymbolCollection,
	}

	return cfg
}

func parseGenSymbolCollection(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolCollection:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolCollection:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolCollection{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolCollection:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolCollectionTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
