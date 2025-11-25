package lowcode

import (
	"context"
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

const GenPositionCollectionTypeName = "genPositionCollection"

type GenPositionCollectionSourceType int

const (
	GPCSTypeAll                GenPositionCollectionSourceType = 0
	GPCSTypeMask               GenPositionCollectionSourceType = 1
	GPCSTypeRowMask            GenPositionCollectionSourceType = 2
	GPCSTypeAllMask            GenPositionCollectionSourceType = 3
	GPCSTypePositionCollection GenPositionCollectionSourceType = 4
)

func parseGenPositionCollectionSourceType(str string) GenPositionCollectionSourceType {
	switch str {
	case "mask":
		return GPCSTypeMask
	case "positioncollection":
		return GPCSTypePositionCollection
	case "rowmask":
		return GPCSTypeRowMask
	case "allmask":
		return GPCSTypeAllMask
	}

	return GPCSTypeAll
}

type GenPositionCollectionCoreType int

const (
	GPCCTypeNumber        GenPositionCollectionCoreType = 0
	GPCCTypeNumberWeight  GenPositionCollectionCoreType = 1
	GPCCTypeEachPosRandom GenPositionCollectionCoreType = 2
)

func parseGenPositionCollectionCoreType(str string) GenPositionCollectionCoreType {
	switch str {
	case "numberweight":
		return GPCCTypeNumberWeight
	case "eachposrandom":
		return GPCCTypeEachPosRandom
	}

	return GPCCTypeNumber
}

// GenPositionCollectionConfig - configuration for GenPositionCollection
type GenPositionCollectionConfig struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	StrSrcType               string                          `yaml:"srcType" json:"srcType"`
	SrcType                  GenPositionCollectionSourceType `yaml:"-" json:"-"`
	StrCoreType              string                          `yaml:"coreType" json:"coreType"`
	CoreType                 GenPositionCollectionCoreType   `yaml:"-" json:"-"`
	Number                   int                             `yaml:"number" json:"number"`
	NumberWeight             string                          `yaml:"numberWeight" json:"numberWeight"`
	NumberWeightVW2          *sgc7game.ValWeights2           `yaml:"-" json:"-"`
	OutputPositionCollection string                          `yaml:"outputPositionCollection" json:"outputPositionCollection"`
	RowMask                  string                          `json:"rowMask" json:"rowMask"`
	Mask                     string                          `json:"mask" json:"mask"`
	SrcPositionCollection    string                          `json:"srcPositionCollection" json:"srcPositionCollection"`
	MapControllers           map[string][]*Award             `yaml:"controllers" json:"controllers"`
	JumpToComponent          string                          `yaml:"jumpToComponent" json:"jumpToComponent"`
}

// SetLinkComponent
func (cfg *GenPositionCollectionConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

type GenPositionCollection struct {
	*BasicComponent `json:"-"`
	Config          *GenPositionCollectionConfig `json:"config"`
}

// Init -
func (genPositionCollection *GenPositionCollection) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenPositionCollection.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenPositionCollectionConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenPositionCollection.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genPositionCollection.InitEx(cfg, pool)
}

// InitEx -
func (genPositionCollection *GenPositionCollection) InitEx(cfg any, pool *GamePropertyPool) error {
	genPositionCollection.Config = cfg.(*GenPositionCollectionConfig)
	genPositionCollection.Config.ComponentType = GenPositionCollectionTypeName

	genPositionCollection.Config.SrcType = parseGenPositionCollectionSourceType(genPositionCollection.Config.StrSrcType)
	genPositionCollection.Config.CoreType = parseGenPositionCollectionCoreType(genPositionCollection.Config.StrCoreType)

	if genPositionCollection.Config.NumberWeight != "" {
		vw2, err := pool.LoadIntWeights(genPositionCollection.Config.NumberWeight, true)
		if err != nil {
			goutils.Error("GenPositionCollection.InitEx:ParseValWeights2",
				slog.String("NumberWeight", genPositionCollection.Config.NumberWeight),
				goutils.Err(err))

			return err
		}

		genPositionCollection.Config.NumberWeightVW2 = vw2
	}

	for _, arr := range genPositionCollection.Config.MapControllers {
		for _, aw := range arr {
			aw.Init()
		}
	}

	genPositionCollection.onInit(&genPositionCollection.Config.BasicComponentConfig)

	return nil
}

// OnProcControllers -
func (genPositionCollection *GenPositionCollection) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	ctrls, isok := genPositionCollection.Config.MapControllers[strVal]
	if isok {
		gameProp.procAwards(plugin, ctrls, curpr, gp)
	}
}

// getSrcPos
func (genPositionCollection *GenPositionCollection) getSrcPos(gameProp *GameProperty, plugin sgc7plugin.IPlugin, gs *sgc7game.GameScene) (*PosData, error) {

	pos := gameProp.posPool.Get()

	switch genPositionCollection.Config.SrcType {
	case GPCSTypePositionCollection:
		curpos := gameProp.GetComponentPos(genPositionCollection.Config.SrcPositionCollection)
		if len(curpos) > 0 {
			for i := range len(curpos) / 2 {
				x := curpos[i*2]
				y := curpos[i*2+1]

				if !pos.Has(x, y) {
					pos.Add(x, y)
				}
			}
		}
	case GPCSTypeAll:
		for x := 0; x < gameProp.GetVal(GamePropWidth); x++ {
			for y := 0; y < gameProp.GetVal(GamePropHeight); y++ {
				pos.Add(x, y)
			}
		}
	case GPCSTypeRowMask:
		imaskd := gameProp.GetComponentDataWithName(genPositionCollection.Config.RowMask)
		if imaskd != nil {
			arr := imaskd.GetMask()
			if len(arr) != gs.Height {
				goutils.Error("GenPositionCollection.getSrcPos:RowMask:len(arr)!=gs.Height",
					goutils.Err(ErrInvalidComponentConfig))

				return nil, ErrInvalidComponentConfig
			}

			for y := 0; y < gs.Height; y++ {
				if arr[y] {
					for x := 0; x < gs.Width; x++ {
						pos.Add(x, y)
					}
				}
			}
		}
	case GPCSTypeMask:
		imaskd := gameProp.GetComponentDataWithName(genPositionCollection.Config.Mask)
		if imaskd != nil {
			arr := imaskd.GetMask()
			if len(arr) != gs.Width {
				goutils.Error("GenPositionCollection.getSrcPos:Mask:len(arr)!=gs.Width",
					slog.String("componentName", genPositionCollection.GetName()),
					goutils.Err(ErrInvalidComponentConfig))

				return nil, ErrInvalidComponentConfig
			}

			for x := 0; x < gs.Width; x++ {
				if arr[x] {
					for y := 0; y < gs.Height; y++ {
						pos.Add(x, y)
					}
				}
			}
		}
	case GPCSTypeAllMask:
		imaskd1 := gameProp.GetComponentDataWithName(genPositionCollection.Config.Mask)
		imaskd2 := gameProp.GetComponentDataWithName(genPositionCollection.Config.RowMask)
		if imaskd1 != nil && imaskd2 != nil {
			arr1 := imaskd1.GetMask()
			if len(arr1) != gs.Width {
				goutils.Error("GenPositionCollection.getSrcPos:CS2STypeAllMask:len(arr1)!=gs.Width",
					goutils.Err(ErrInvalidComponentConfig))

				return nil, ErrInvalidComponentConfig
			}

			arr2 := imaskd2.GetMask()
			if len(arr2) != gs.Height {
				goutils.Error("GenPositionCollection.getSrcPos:CS2STypeAllMask:len(arr2)!=gs.Height",
					goutils.Err(ErrInvalidComponentConfig))

				return nil, ErrInvalidComponentConfig
			}

			for x := 0; x < gs.Width; x++ {
				if arr1[x] {
					for y := 0; y < gs.Height; y++ {
						if arr2[y] {
							pos.Add(x, y)
						}
					}
				}
			}
		}
	default:
		goutils.Error("GenPositionCollection.getSrcPos:ErrUnsupportedSourceType",
			slog.String("srcType", genPositionCollection.Config.StrSrcType),
			goutils.Err(ErrInvalidComponentConfig))

		return nil, ErrInvalidComponentConfig
	}

	return pos, nil
}

// playgame
func (genPositionCollection *GenPositionCollection) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// bcd := icd.(*BasicComponentData)

	gs := genPositionCollection.GetTargetScene3(gameProp, curpr, prs, 0)

	pos, err := genPositionCollection.getSrcPos(gameProp, plugin, gs)
	if err != nil {
		goutils.Error("GenPositionCollection.OnPlayGame:getSrcPos",
			goutils.Err(err))

		return "", err
	}

	pccd := gameProp.GetComponentDataWithName(genPositionCollection.Config.OutputPositionCollection)
	if pccd == nil {
		goutils.Error("GenPositionCollection.OnPlayGame:GetComponentDataWithName",
			slog.String("OutputPositionCollection", genPositionCollection.Config.OutputPositionCollection),
			goutils.Err(ErrNoComponent))

		return "", ErrNoComponent
	}

	isTrigger := false

	switch genPositionCollection.Config.CoreType {
	case GPCCTypeNumber:
		n := genPositionCollection.Config.Number
		if n > pos.Len() {
			n = pos.Len()
		}

		for i := 0; i < n; i++ {
			ri, err := plugin.Random(context.Background(), pos.Len())
			if err != nil {
				goutils.Error("GenPositionCollection.OnPlayGame:Random",
					goutils.Err(err))

				return "", err
			}

			x := pos.pos[ri*2]
			y := pos.pos[ri*2+1]

			pccd.AddPos(x, y)

			pos.Del(ri)

			isTrigger = true
		}
	case GPCCTypeNumberWeight:
		vw := genPositionCollection.Config.NumberWeightVW2
		if vw == nil {
			goutils.Error("GenPositionCollection.OnPlayGame:NumberWeightVW2==nil",
				slog.String("NumberWeight", genPositionCollection.Config.NumberWeight),
				goutils.Err(ErrInvalidComponentConfig))

			return "", ErrInvalidComponentConfig
		}

		cn, err := vw.RandVal(plugin)
		if err != nil {
			goutils.Error("GenPositionCollection.OnPlayGame:RandVal",
				goutils.Err(err))

			return "", err
		}

		n := cn.Int()
		if n > pos.Len() {
			n = pos.Len()
		}

		for i := 0; i < n; i++ {
			ri, err := plugin.Random(context.Background(), pos.Len())
			if err != nil {
				goutils.Error("GenPositionCollection.OnPlayGame:Random",
					goutils.Err(err))

				return "", err
			}

			x := pos.pos[ri*2]
			y := pos.pos[ri*2+1]

			pccd.AddPos(x, y)

			pos.Del(ri)

			isTrigger = true
		}
	default:
		goutils.Error("GenPositionCollection.OnPlayGame:ErrUnsupportedCoreType",
			slog.String("coreType", genPositionCollection.Config.StrCoreType),
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	if isTrigger {
		genPositionCollection.ProcControllers(gameProp, plugin, curpr, gp, -1, "<trigger>")

		nc := genPositionCollection.onStepEnd(gameProp, curpr, gp, genPositionCollection.Config.JumpToComponent)

		return nc, nil
	}

	nc := genPositionCollection.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genPositionCollection *GenPositionCollection) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	return nil
}

func NewGenPositionCollection(name string) IComponent {
	return &GenPositionCollection{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "srcType": "all",
// "coreType": "numberWeight",
// "number": 1,
// "numberWeight": "cg-rollmunum",
// "outputPositionCollection": "cg-pos-multi"
type jsonGenPositionCollection struct {
	SrcType                  string `json:"srcType"`
	CoreType                 string `json:"coreType"`
	Number                   int    `json:"number"`
	NumberWeight             string `json:"numberWeight"`
	OutputPositionCollection string `json:"outputPositionCollection"`
	RowMask                  string `json:"rowMask"`
	Mask                     string `json:"mask"`
	SrcPositionCollection    string `json:"srcPositionCollection"`
}

func (jcfg *jsonGenPositionCollection) build() *GenPositionCollectionConfig {
	cfg := &GenPositionCollectionConfig{
		StrSrcType:               strings.ToLower(jcfg.SrcType),
		StrCoreType:              strings.ToLower(jcfg.CoreType),
		Number:                   jcfg.Number,
		NumberWeight:             jcfg.NumberWeight,
		OutputPositionCollection: jcfg.OutputPositionCollection,
		Mask:                     jcfg.Mask,
		RowMask:                  jcfg.RowMask,
		SrcPositionCollection:    jcfg.SrcPositionCollection,
	}

	return cfg
}

func parseGenPositionCollection(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenPositionCollection:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenPositionCollection:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenPositionCollection{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenPositionCollection:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenPositionCollectionTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
