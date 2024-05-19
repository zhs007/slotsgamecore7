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

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

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

		genSymbolValsWithSymbol.AddOtherScene(gameProp, curpr, nos, cd)

		nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := genSymbolValsWithSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (genSymbolValsWithPos *GenSymbolValsWithSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after GenSymbolValsWithSymbol", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// // OnStats
// func (genSymbolValsWithPos *GenSymbolValsWithSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

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
	IsUseSource string   `json:"isUseSource"`
	IsAlwaysGen string   `json:"isAlwaysGen"`
}

func (jcfg *jsonGenSymbolValsWithSymbol) build() *GenSymbolValsWithSymbolConfig {
	cfg := &GenSymbolValsWithSymbolConfig{
		StrType:     jcfg.Type,
		Symbols:     jcfg.Symbols,
		Weight:      jcfg.Weight,
		DefaultVal:  jcfg.DefaultVal,
		IsUseSource: jcfg.IsUseSource == "true",
		IsAlwaysGen: jcfg.IsAlwaysGen == "true",
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
