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

const ReplaceSymbolGroupTypeName = "replaceSymbolGroup"

// ReplaceSymbolGroupConfig - configuration for ReplaceSymbolGroup
type ReplaceSymbolGroupConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SrcSymbols           []string `yaml:"srcSymbols" json:"srcSymbols"`
	SrcSymbolCodes       []int    `yaml:"-" json:"-"`
	TargetSymbols        []string `yaml:"targetSymbols" json:"targetSymbols"`
	TargetSymbolCodes    []int    `yaml:"-" json:"-"`
	Mask                 string   `yaml:"mask" json:"mask"`
}

// SetLinkComponent
func (cfg *ReplaceSymbolGroupConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ReplaceSymbolGroup struct {
	*BasicComponent `json:"-"`
	Config          *ReplaceSymbolGroupConfig `json:"config"`
}

// Init -
func (replaceSymbolGroup *ReplaceSymbolGroup) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReplaceSymbolGroup.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &ReplaceSymbolGroupConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReplaceSymbolGroup.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return replaceSymbolGroup.InitEx(cfg, pool)
}

// InitEx -
func (replaceSymbolGroup *ReplaceSymbolGroup) InitEx(cfg any, pool *GamePropertyPool) error {
	replaceSymbolGroup.Config = cfg.(*ReplaceSymbolGroupConfig)
	replaceSymbolGroup.Config.ComponentType = ReplaceSymbolGroupTypeName

	for _, v := range replaceSymbolGroup.Config.SrcSymbols {
		replaceSymbolGroup.Config.SrcSymbolCodes = append(replaceSymbolGroup.Config.SrcSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	for _, v := range replaceSymbolGroup.Config.TargetSymbols {
		replaceSymbolGroup.Config.TargetSymbolCodes = append(replaceSymbolGroup.Config.TargetSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	if len(replaceSymbolGroup.Config.SrcSymbolCodes) > len(replaceSymbolGroup.Config.TargetSymbolCodes) {
		goutils.Error("ReplaceSymbolGroup.InitEx:invalid symbols",
			slog.Int("src", len(replaceSymbolGroup.Config.SrcSymbolCodes)),
			slog.Int("target", len(replaceSymbolGroup.Config.TargetSymbolCodes)),
			goutils.Err(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	replaceSymbolGroup.onInit(&replaceSymbolGroup.Config.BasicComponentConfig)

	return nil
}

// playgame
func (replaceSymbolGroup *ReplaceSymbolGroup) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceSymbolGroup.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	if len(replaceSymbolGroup.Config.SrcSymbolCodes) == 0 {
		nc := replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	cd := icd.(*BasicComponentData)

	gs := replaceSymbolGroup.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	for x, arr := range gs.Arr {
		for y, srcSymbol := range arr {
			si := goutils.IndexOfIntSlice(replaceSymbolGroup.Config.SrcSymbolCodes, srcSymbol, 0)
			if si >= 0 {
				if ngs == gs {
					ngs = gs.CloneEx(gameProp.PoolScene)
				}

				ngs.Arr[x][y] = replaceSymbolGroup.Config.TargetSymbolCodes[si]
			}
		}
	}

	if ngs == gs {
		nc := replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

		return nc, ErrComponentDoNothing
	}

	replaceSymbolGroup.AddScene(gameProp, curpr, ngs, cd)

	nc := replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (replaceSymbolGroup *ReplaceSymbolGroup) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after replaceSymbolGroup", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (replaceSymbolGroup *ReplaceSymbolGroup) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewReplaceSymbolGroup(name string) IComponent {
	return &ReplaceSymbolGroup{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"srcSymbol targetSymbol": [
//			{
//				"srcSymbols": "RH",
//				"targetSymbol": "GH"
//			},
//			{
//				"srcSymbols": "RM",
//				"targetSymbol": "GM"
//			},
//			{
//				"srcSymbols": "RL",
//				"targetSymbol": "GL"
//			}
//		]
//	},
type jsonReplaceSymbolGroupNode struct {
	SrcSymbols   string `json:"srcSymbols"`   // src
	TargetSymbol string `json:"targetSymbol"` // target
}

type jsonReplaceSymbolGroup struct {
	Symbols []*jsonReplaceSymbolGroupNode `json:"srcSymbol targetSymbol"`
}

func (jcfg *jsonReplaceSymbolGroup) build() *ReplaceSymbolGroupConfig {
	cfg := &ReplaceSymbolGroupConfig{}

	for _, s := range jcfg.Symbols {
		cfg.SrcSymbols = append(cfg.SrcSymbols, s.SrcSymbols)
		cfg.TargetSymbols = append(cfg.TargetSymbols, s.TargetSymbol)
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseReplaceSymbolGroup(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseReplaceSymbolGroup:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseReplaceSymbolGroup:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonReplaceSymbolGroup{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseReplaceSymbolGroup:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ReplaceSymbolGroupTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
