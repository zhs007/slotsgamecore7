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

const ChgSymbolValsTypeName = "chgSymbolVals"

type ChgSymbolValsType int

const (
	CSVTypeInc ChgSymbolValsType = 0 // ++
	CSVTypeDec ChgSymbolValsType = 1 // --
)

func parseChgSymbolValsType(strType string) ChgSymbolValsType {
	if strType == "dec" {
		return CSVTypeDec
	}

	return CSVTypeInc
}

type ChgSymbolValsSourceType int

const (
	CSVSTypePositionCollection ChgSymbolValsSourceType = 0 // positionCollection
	CSVSTypeWinResult          ChgSymbolValsSourceType = 1 // winResult
)

func parseChgSymbolValsSourceType(strType string) ChgSymbolValsSourceType {
	if strType == "positioncollection" {
		return CSVSTypePositionCollection
	}

	return CSVSTypeWinResult
}

// ChgSymbolValsConfig - configuration for ChgSymbolVals
type ChgSymbolValsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	StrType              string                  `yaml:"type" json:"type"`
	Type                 ChgSymbolValsType       `yaml:"-" json:"-"`
	StrSourceType        string                  `yaml:"sourceType" json:"sourceType"`
	SourceType           ChgSymbolValsSourceType `yaml:"-" json:"-"`
	PositionCollection   string                  `yaml:"positionCollection" json:"positionCollection"`
	WinResultComponents  []string                `yaml:"winResultComponents" json:"winResultComponents"`
	MaxNumber            int                     `yaml:"maxNumber" json:"maxNumber"`
}

// SetLinkComponent
func (cfg *ChgSymbolValsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ChgSymbolVals struct {
	*BasicComponent `json:"-"`
	Config          *ChgSymbolValsConfig `json:"config"`
}

// Init -
func (chgSymbolVals *ChgSymbolVals) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ChgSymbolVals.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ChgSymbolValsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ChgSymbolVals.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return chgSymbolVals.InitEx(cfg, pool)
}

// InitEx -
func (chgSymbolVals *ChgSymbolVals) InitEx(cfg any, pool *GamePropertyPool) error {
	chgSymbolVals.Config = cfg.(*ChgSymbolValsConfig)
	chgSymbolVals.Config.ComponentType = ChgSymbolValsTypeName

	chgSymbolVals.Config.Type = parseChgSymbolValsType(chgSymbolVals.Config.StrType)
	chgSymbolVals.Config.SourceType = parseChgSymbolValsSourceType(chgSymbolVals.Config.StrSourceType)

	chgSymbolVals.onInit(&chgSymbolVals.Config.BasicComponentConfig)

	return nil
}

func (chgSymbolVals *ChgSymbolVals) rebuildPos(pos []int, plugin sgc7plugin.IPlugin) ([]int, error) {
	if chgSymbolVals.Config.MaxNumber <= 0 {
		return pos, nil
	}

	if len(pos)/2 <= chgSymbolVals.Config.MaxNumber {
		return pos, nil
	}

	npos := []int{}

	for i := 0; i < chgSymbolVals.Config.MaxNumber; i++ {
		cr, err := plugin.Random(context.Background(), len(pos)/2)
		if err != nil {
			goutils.Error("ChgSymbolVals.rebuildPos:Random",
				goutils.Err(err))

			return nil, err
		}

		npos = append(npos, pos[cr*2], pos[cr*2+1])

		pos = append(pos[:cr*2], pos[(cr+1)*2:]...)
	}

	return npos, nil
}

// playgame
func (chgSymbolVals *ChgSymbolVals) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// symbolVal2.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	os := chgSymbolVals.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	if os != nil {
		nos := os

		if chgSymbolVals.Config.SourceType == CSVSTypePositionCollection {
			pc, isok := gameProp.Components.MapComponents[chgSymbolVals.Config.PositionCollection]
			if isok {
				pccd := gameProp.GetComponentData(pc)
				pos := pccd.GetPos()
				if len(pos) > 0 {
					npos, err := chgSymbolVals.rebuildPos(pos, plugin)
					if err != nil {
						goutils.Error("ChgSymbolVals.OnPlayGame:rebuildPos",
							goutils.Err(err))

						return "", nil
					}

					if chgSymbolVals.Config.Type == CSVTypeInc {
						nos = os.CloneEx(gameProp.PoolScene)

						for i := 0; i < len(npos)/2; i++ {
							nos.Arr[npos[i*2]][npos[i*2+1]]++
						}
					} else if chgSymbolVals.Config.Type == CSVTypeDec {
						for i := 0; i < len(npos)/2; i++ {
							if nos.Arr[npos[i*2]][npos[i*2+1]] > 0 {
								if nos == os {
									nos = os.CloneEx(gameProp.PoolScene)
								}

								nos.Arr[npos[i*2]][npos[i*2+1]]--
							}
						}
					}
				}
			}
		} else {
			for _, cn := range chgSymbolVals.Config.WinResultComponents {
				ccd := gameProp.GetComponentDataWithName(cn)
				// ccd := gameProp.MapComponentData[cn]
				lst := ccd.GetResults()
				for _, ri := range lst {
					pos := curpr.Results[ri].Pos
					if len(pos) > 0 {
						if chgSymbolVals.Config.Type == CSVTypeInc {
							if nos == os {
								nos = os.CloneEx(gameProp.PoolScene)
							}

							for i := 0; i < len(pos)/2; i++ {
								nos.Arr[pos[i*2]][pos[i*2+1]]++
							}
						} else if chgSymbolVals.Config.Type == CSVTypeDec {
							for i := 0; i < len(pos)/2; i++ {
								if nos.Arr[pos[i*2]][pos[i*2+1]] > 0 {
									if nos == os {
										nos = os.CloneEx(gameProp.PoolScene)
									}

									nos.Arr[pos[i*2]][pos[i*2+1]]--
								}
							}
						}
					}
				}
			}
		}

		if nos == os {
			nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		chgSymbolVals.AddOtherScene(gameProp, curpr, nos, cd)

		nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

		return nc, nil
	}

	nc := chgSymbolVals.onStepEnd(gameProp, curpr, gp, "")

	return nc, ErrComponentDoNothing
}

// OnAsciiGame - outpur to asciigame
func (chgSymbolVals *ChgSymbolVals) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {

	cd := icd.(*BasicComponentData)

	if len(cd.UsedOtherScenes) > 0 {
		asciigame.OutputOtherScene("after ChgSymbolVals", pr.OtherScenes[cd.UsedOtherScenes[0]])
	}

	return nil
}

// // OnStats
// func (chgSymbolVals *ChgSymbolVals) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewChgSymbolVals(name string) IComponent {
	return &ChgSymbolVals{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

type jsonChgSymbolVals struct {
	Type                string   `json:"type"`
	SourceType          string   `json:"sourceType"`
	PositionCollection  string   `json:"positionCollection"`
	WinResultComponents []string `json:"winResultComponents"`
	MaxNumber           int      `json:"maxNumber"`
}

func (jcfg *jsonChgSymbolVals) build() *ChgSymbolValsConfig {
	cfg := &ChgSymbolValsConfig{
		StrType:             jcfg.Type,
		StrSourceType:       strings.ToLower(jcfg.SourceType),
		PositionCollection:  jcfg.PositionCollection,
		WinResultComponents: jcfg.WinResultComponents,
		MaxNumber:           jcfg.MaxNumber,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseChgSymbolVals(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseChgSymbolVals:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseChgSymbolVals:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonChgSymbolVals{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseChgSymbolVals:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ChgSymbolValsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
