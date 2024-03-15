package lowcode

import (
	"os"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const GenSymbolValsWithPosTypeName = "genSymbolValsWithPos"

type GenSymbolValsWithPosType int

const (
	GSVWPTypeAdd               GenSymbolValsWithPosType = 0
	GSVWPTypeMask              GenSymbolValsWithPosType = 1
	GSVWPTypeAddWithIntMapping GenSymbolValsWithPosType = 2
)

func parseGenSymbolValsWithPosType(strType string) GenSymbolValsWithPosType {
	if strType == "mask" {
		return GSVWPTypeMask
	} else if strType == "addWithIntMapping" {
		return GSVWPTypeAddWithIntMapping
	}

	return GSVWPTypeAdd
}

// GenSymbolValsWithPosConfig - configuration for GenSymbolValsWithPos
type GenSymbolValsWithPosConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	TargetComponents     []string                 `yaml:"targetComponents" json:"targetComponents"`
	StrType              string                   `yaml:"genType" json:"genType"`
	Type                 GenSymbolValsWithPosType `yaml:"-" json:"-"`
	ValMapping           string                   `yaml:"valMapping" json:"valMapping"`
	ValMappingVM         *sgc7game.ValMapping2    `yaml:"-" json:"-"`
	IsUseSource          bool                     `yaml:"isUseSource" json:"isUseSource"`
	IsAlwaysGen          bool                     `yaml:"isAlwaysGen" json:"isAlwaysGen"`
}

// SetLinkComponent
func (cfg *GenSymbolValsWithPosConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type GenSymbolValsWithPos struct {
	*BasicComponent `json:"-"`
	Config          *GenSymbolValsWithPosConfig `json:"config"`
}

// Init -
func (genSymbolValsWithPos *GenSymbolValsWithPos) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("GenSymbolValsWithPos.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &GenSymbolValsWithPosConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("GenSymbolValsWithPos.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return genSymbolValsWithPos.InitEx(cfg, pool)
}

// InitEx -
func (genSymbolValsWithPos *GenSymbolValsWithPos) InitEx(cfg any, pool *GamePropertyPool) error {
	genSymbolValsWithPos.Config = cfg.(*GenSymbolValsWithPosConfig)
	genSymbolValsWithPos.Config.ComponentType = GenSymbolValsWithPosTypeName

	genSymbolValsWithPos.Config.Type = parseGenSymbolValsWithPosType(genSymbolValsWithPos.Config.StrType)

	if genSymbolValsWithPos.Config.ValMapping != "" {
		vm2 := pool.LoadIntMapping(genSymbolValsWithPos.Config.ValMapping)
		if vm2 == nil {
			goutils.Error("GenSymbolValsWithPos.Init:LoadIntMapping",
				zap.String("ValMapping", genSymbolValsWithPos.Config.ValMapping),
				zap.Error(ErrInvalidIntValMappingFile))

			return ErrInvalidIntValMappingFile
		}

		genSymbolValsWithPos.Config.ValMappingVM = vm2
	}

	genSymbolValsWithPos.onInit(&genSymbolValsWithPos.Config.BasicComponentConfig)

	return nil
}

// playgame
func (genSymbolValsWithPos *GenSymbolValsWithPos) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	var os *sgc7game.GameScene

	if genSymbolValsWithPos.Config.IsUseSource {
		os = genSymbolValsWithPos.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	nos := os

	if genSymbolValsWithPos.Config.Type == GSVWPTypeAdd {
		for _, cn := range genSymbolValsWithPos.Config.TargetComponents {
			ccd := gameProp.GetCurComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				for pi := 0; pi < len(curpr.Results[ri].Pos)/2; pi++ {
					x := curpr.Results[ri].Pos[pi*2]
					y := curpr.Results[ri].Pos[pi*2+1]

					if nos == os {
						if os != nil {
							nos = os.CloneEx(gameProp.PoolScene)
						} else {
							nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), 0)
						}
					}

					nos.Arr[x][y]++
				}
			}
		}
	} else if genSymbolValsWithPos.Config.Type == GSVWPTypeMask {
		for _, cn := range genSymbolValsWithPos.Config.TargetComponents {
			ccd := gameProp.GetCurComponentDataWithName(cn)
			lst := ccd.GetResults()
			for _, ri := range lst {
				for pi := 0; pi < len(curpr.Results[ri].Pos)/2; pi++ {
					x := curpr.Results[ri].Pos[pi*2]
					y := curpr.Results[ri].Pos[pi*2+1]

					if os != nil && os.Arr[x][y] > 0 {
						continue
					}

					if nos == os {
						if os != nil {
							nos = os.CloneEx(gameProp.PoolScene)
						} else {
							nos = gameProp.PoolScene.New2(gameProp.GetVal(GamePropWidth), gameProp.GetVal(GamePropHeight), 0)
						}
					}

					nos.Arr[x][y]++
				}
			}
		}
	}

	if nos == os {
		nc := genSymbolValsWithPos.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	genSymbolValsWithPos.AddOtherScene(gameProp, curpr, nos, cd)

	nc := genSymbolValsWithPos.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (genSymbolValsWithPos *GenSymbolValsWithPos) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after GenSymbolValsWithPos", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// OnStats
func (genSymbolValsWithPos *GenSymbolValsWithPos) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewGenSymbolValsWithPos(name string) IComponent {
	return &GenSymbolValsWithPos{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "isAlwaysGen": "true",
// "isUseSource": "true",
// "targetComponents": [
//
//	"bg-pay",
//	"bg-firstcheckwins"
//
// ],
// "genType": "add"
type jsonGenSymbolValsWithPos struct {
	TargetComponents []string `json:"targetComponents"`
	StrType          string   `json:"genType"`
	ValMapping       string   `json:"valMapping"`
	IsUseSource      string   `json:"isUseSource"`
	IsAlwaysGen      string   `json:"isAlwaysGen"`
}

func (jcfg *jsonGenSymbolValsWithPos) build() *GenSymbolValsWithPosConfig {
	cfg := &GenSymbolValsWithPosConfig{
		StrType:          jcfg.StrType,
		TargetComponents: jcfg.TargetComponents,
		ValMapping:       jcfg.ValMapping,
		IsUseSource:      jcfg.IsUseSource == "true",
		IsAlwaysGen:      jcfg.IsAlwaysGen == "true",
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseGenSymbolValsWithPos(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseGenSymbolValsWithPos:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseGenSymbolValsWithPos:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonGenSymbolValsWithPos{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseGenSymbolValsWithPos:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: GenSymbolValsWithPosTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
