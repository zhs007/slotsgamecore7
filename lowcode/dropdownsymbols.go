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
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const DropDownSymbolsTypeName = "dropDownSymbols"

// DropDownSymbolsConfig - configuration for DropDownSymbols
type DropDownSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	HoldSymbols          []string `yaml:"holdSymbols" json:"holdSymbols"`                   // 不需要下落的symbol
	HoldSymbolCodes      []int    `yaml:"-" json:"-"`                                       // 不需要下落的symbol
	IsNeedProcSymbolVals bool     `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
}

// SetLinkComponent
func (cfg *DropDownSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type DropDownSymbols struct {
	*BasicComponent `json:"-"`
	Config          *DropDownSymbolsConfig `json:"config"`
}

// Init -
func (dropDownSymbols *DropDownSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("DropDownSymbols.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &DropDownSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("DropDownSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return dropDownSymbols.InitEx(cfg, pool)
}

// InitEx -
func (dropDownSymbols *DropDownSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	dropDownSymbols.Config = cfg.(*DropDownSymbolsConfig)
	dropDownSymbols.Config.ComponentType = DropDownSymbolsTypeName

	for _, v := range dropDownSymbols.Config.HoldSymbols {
		dropDownSymbols.Config.HoldSymbolCodes = append(dropDownSymbols.Config.HoldSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	dropDownSymbols.onInit(&dropDownSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (dropDownSymbols *DropDownSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// dropDownSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	gs := dropDownSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	if gs == nil {
		goutils.Error("DropDownSymbols.OnPlayGame",
			goutils.Err(ErrInvalidScene))

		return "", ErrInvalidScene
	}

	if !gs.HasSymbol(-1) {
		nc := dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	ngs := gs.CloneEx(gameProp.PoolScene)

	var os *sgc7game.GameScene
	if dropDownSymbols.Config.IsNeedProcSymbolVals {
		os = dropDownSymbols.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	if os != nil {
		nos := os.CloneEx(gameProp.PoolScene)

		for x, arr := range ngs.Arr {
			for y := len(arr) - 1; y >= 0; {
				if arr[y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						if arr[y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
							arr[y] = arr[y1]
							arr[y1] = -1

							nos.Arr[x][y] = nos.Arr[x][y1]
							nos.Arr[x][y1] = dropDownSymbols.Config.EmptySymbolVal

							hass = true
							y--
							break
						}
					}

					if !hass {
						break
					}
				} else {
					y--
				}
			}
		}

		dropDownSymbols.AddOtherScene(gameProp, curpr, nos, bcd)
	} else {
		for x, arr := range ngs.Arr {
			for y := len(arr) - 1; y >= 0; {
				if arr[y] == -1 {
					hass := false
					for y1 := y - 1; y1 >= 0; y1-- {
						if arr[y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
							arr[y] = arr[y1]
							arr[y1] = -1

							hass = true
							y--
							break
						}
					}

					if !hass {
						break
					}
				} else {
					y--
				}
			}
		}
	}

	dropDownSymbols.AddScene(gameProp, curpr, ngs, bcd)

	nc := dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (dropDownSymbols *DropDownSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after dropDownSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (dropDownSymbols *DropDownSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (dropDownSymbols *DropDownSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	return 0, nil
// }

// EachUsedResults -
func (dropDownSymbols *DropDownSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewDropDownSymbols(name string) IComponent {
	return &DropDownSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

// "configuration": {},
type jsonDropDownSymbols struct {
	HoldSymbols          []string `json:"holdSymbols"`                                      // 不需要下落的symbol
	IsNeedProcSymbolVals bool     `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
}

func (jcfg *jsonDropDownSymbols) build() *DropDownSymbolsConfig {
	cfg := &DropDownSymbolsConfig{
		HoldSymbols:          jcfg.HoldSymbols,
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseDropDownSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseDropDownSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseDropDownSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonDropDownSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseDropDownSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: DropDownSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
