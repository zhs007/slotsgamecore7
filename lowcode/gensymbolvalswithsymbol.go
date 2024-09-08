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
	"github.com/zhs007/slotsgamecore7/stats2"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v2"
)

const GenSymbolValsWithSymbolTypeName = "genSymbolValsWithSymbol"

type GenSymbolValsWithSymbolType int

const (
	GSVWSTypeNormal   GenSymbolValsWithSymbolType = 0
	GSVWSTypeNonClear GenSymbolValsWithSymbolType = 1
)

func parseGenSymbolValsWithSymbolType(strType string) GenSymbolValsWithSymbolType {
	if strType == "non-clear" {
		return GSVWSTypeNonClear
	}

	return GSVWSTypeNormal
}

type GenSymbolValsWithSymbolData struct {
	BasicComponentData
	GenVals []int
}

// OnNewGame -
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) OnNewGame(gameProp *GameProperty, component IComponent) {
	genSymbolValsWithSymbolData.BasicComponentData.OnNewGame(gameProp, component)
}

// OnNewStep -
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) onNewStep() {
	genSymbolValsWithSymbolData.UsedOtherScenes = nil
	genSymbolValsWithSymbolData.GenVals = nil
}

// Clone
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) Clone() IComponentData {
	target := &GenSymbolValsWithSymbolData{
		BasicComponentData: genSymbolValsWithSymbolData.CloneBasicComponentData(),
	}

	target.GenVals = make([]int, len(genSymbolValsWithSymbolData.GenVals))
	copy(target.GenVals, genSymbolValsWithSymbolData.GenVals)

	return target
}

// BuildPBComponentData
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) BuildPBComponentData() proto.Message {
	return genSymbolValsWithSymbolData.BasicComponentData.BuildPBComponentData()
}

// addVal
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) addVal(val int) {
	genSymbolValsWithSymbolData.GenVals = append(genSymbolValsWithSymbolData.GenVals, val)
}

// GetValEx -
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) GetVal(key string) (int, bool) {
	return genSymbolValsWithSymbolData.GetValEx(key, GCVTypeNormal)
}

// GetValEx -
func (genSymbolValsWithSymbolData *GenSymbolValsWithSymbolData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVSymbolVal && len(genSymbolValsWithSymbolData.GenVals) > 0 {
		if getType == GCVTypeMin {
			return slices.Min(genSymbolValsWithSymbolData.GenVals), true
		} else if getType == GCVTypeMax {
			return slices.Max(genSymbolValsWithSymbolData.GenVals), true
		} else {
			return genSymbolValsWithSymbolData.GenVals[len(genSymbolValsWithSymbolData.GenVals)-1], true
		}
	}

	return 0, false
}

// GenSymbolValsWithSymbolConfig - configuration for GenSymbolValsWithSymbol
type GenSymbolValsWithSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                      `yaml:"type" json:"type"`
	Type                 GenSymbolValsWithSymbolType `yaml:"-" json:"-"`
	Symbols              []string                    `yaml:"symbols" json:"symbols"`
	SymbolCodes          []int                       `yaml:"-" json:"-"`
	Weight               string                      `yaml:"weight" json:"weight"`
	WeightVW2            *sgc7game.ValWeights2       `yaml:"-" json:"-"`
	DefaultVal           int                         `yaml:"defaultVal" json:"defaultVal"`
	IsUseSource          bool                        `yaml:"isUseSource" json:"isUseSource"`
	IsAlwaysGen          bool                        `yaml:"isAlwaysGen" json:"isAlwaysGen"`
}

// SetLinkComponent
func (cfg *GenSymbolValsWithSymbolConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolValsWithSymbol struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolValsWithSymbolConfig `json:"config"`
}

// Init -
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolValsWithSymbol.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &GenSymbolValsWithSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolValsWithSymbol.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return genSymbolValsWithSymbol.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolValsWithSymbol.Config = cfg.(*GenSymbolValsWithSymbolConfig)
	genSymbolValsWithSymbol.Config.ComponentType = GenSymbolValsWithSymbolTypeName

	genSymbolValsWithSymbol.Config.Type = parseGenSymbolValsWithSymbolType(genSymbolValsWithSymbol.Config.StrType)

	for _, s := range genSymbolValsWithSymbol.Config.Symbols {
		sc, isok := pool.DefaultPaytables.MapSymbols[s]
		if !isok {
			goutils.Error("GenSymbolValsWithSymbol.InitEx:Symbol",
				slog.String("symbol", s),
				goutils.Err(ErrIvalidSymbol))
		}

		genSymbolValsWithSymbol.Config.SymbolCodes = append(genSymbolValsWithSymbol.Config.SymbolCodes, sc)
	}

	if genSymbolValsWithSymbol.Config.Weight != "" {
		vw2, err := pool.LoadIntWeights(genSymbolValsWithSymbol.Config.Weight, genSymbolValsWithSymbol.Config.UseFileMapping)
		if err != nil {
			goutils.Error("ChgSymbols.InitEx:LoadIntWeights",
				slog.String("Weight", genSymbolValsWithSymbol.Config.Weight),
				goutils.Err(err))

			return err
		}

		genSymbolValsWithSymbol.Config.WeightVW2 = vw2
	}

	genSymbolValsWithSymbol.onInit(&genSymbolValsWithSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	cd := icd.(*GenSymbolValsWithSymbolData)

	cd.onNewStep()

	gs := genSymbolValsWithSymbol.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs != nil {
		var os *sgc7game.GameScene

		if genSymbolValsWithSymbol.Config.IsUseSource {
			os = genSymbolValsWithSymbol.GetTargetOtherScene3(gameProp, curpr, prs, 0)
		}

		nos := os

		if genSymbolValsWithSymbol.Config.Type == GSVWSTypeNormal {
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(genSymbolValsWithSymbol.Config.SymbolCodes, s, 0) >= 0 {
						if nos == nil {
							curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
							if err != nil {
								goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
									goutils.Err(err))

								return "", err
							}

							nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)

							cd.addVal(curv.Int())
							nos.Arr[x][y] = curv.Int()
						} else if nos.Arr[x][y] == genSymbolValsWithSymbol.Config.DefaultVal {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
							if err != nil {
								goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
									goutils.Err(err))

								return "", err
							}

							cd.addVal(curv.Int())
							nos.Arr[x][y] = curv.Int()
						}
					} else {
						if nos != nil && nos.Arr[x][y] != genSymbolValsWithSymbol.Config.DefaultVal {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							nos.Arr[x][y] = genSymbolValsWithSymbol.Config.DefaultVal
						}
					}
				}
			}
		} else if genSymbolValsWithSymbol.Config.Type == GSVWSTypeNonClear {
			for x, arr := range gs.Arr {
				for y, s := range arr {
					if goutils.IndexOfIntSlice(genSymbolValsWithSymbol.Config.SymbolCodes, s, 0) >= 0 {
						if nos == nil {
							curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
							if err != nil {
								goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
									goutils.Err(err))

								return "", err
							}

							nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)

							cd.addVal(curv.Int())
							nos.Arr[x][y] = curv.Int()
						} else if nos.Arr[x][y] == genSymbolValsWithSymbol.Config.DefaultVal {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							curv, err := genSymbolValsWithSymbol.Config.WeightVW2.RandVal(plugin)
							if err != nil {
								goutils.Error("GenSymbolValsWithSymbol.OnPlayGame:RandVal",
									goutils.Err(err))

								return "", err
							}

							cd.addVal(curv.Int())
							nos.Arr[x][y] = curv.Int()
						}
					}
				}
			}
		}

		if nos == nil && genSymbolValsWithSymbol.Config.IsAlwaysGen {
			nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), genSymbolValsWithSymbol.Config.DefaultVal)
		}

		if nos == os {
			nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		genSymbolValsWithSymbol.AddOtherScene(gameProp, curpr, nos, &cd.BasicComponentData)

		nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after GenSymbolValsWithSymbol", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// NewComponentData -
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) NewComponentData() IComponentData {
	return &GenSymbolValsWithSymbolData{}
}

// OnStats2
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) OnStats2(icd IComponentData, s2 *stats2.Cache, gameProp *GameProperty, gp *GameParams, pr *sgc7game.PlayResult) {
	genSymbolValsWithSymbol.BasicComponent.OnStats2(icd, s2, gameProp, gp, pr)

	cd := icd.(*GenSymbolValsWithSymbolData)

	for _, v := range cd.GenVals {
		s2.ProcStatsIntVal(genSymbolValsWithSymbol.GetName(), v)
	}
}

// NewStats2 -
func (genSymbolValsWithSymbol *GenSymbolValsWithSymbol) NewStats2(parent string) *stats2.Feature {
	return stats2.NewFeature(parent, []stats2.Option{stats2.OptIntVal})
}

func NewGenSymbolValsWithSymbol(name string) IComponent {
	return &GenSymbolValsWithSymbol{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonGenSymbolValsWithSymbol struct {
	Type        string   `json:"type"`
	Symbols     []string `json:"symbols"`
	Weight      string   `json:"weight"`
	DefaultVal  int      `json:"defaultVal"`
	IsUseSource bool     `json:"isUseSource"`
	IsAlwaysGen bool     `json:"isAlwaysGen"`
}

func (jcfg *jsonGenSymbolValsWithSymbol) build() *GenSymbolValsWithSymbolConfig {
	cfg := &GenSymbolValsWithSymbolConfig{
		StrType:     jcfg.Type,
		Symbols:     jcfg.Symbols,
		Weight:      jcfg.Weight,
		DefaultVal:  jcfg.DefaultVal,
		IsUseSource: jcfg.IsUseSource,
		IsAlwaysGen: jcfg.IsAlwaysGen,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseGenSymbolValsWithSymbol(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolValsWithSymbol:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolValsWithSymbol:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonGenSymbolValsWithSymbol{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolValsWithSymbol:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolValsWithSymbolTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
