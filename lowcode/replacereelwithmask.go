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

const ReplaceReelWithMaskTypeName = "replaceReelWithMask"

// ReplaceReelWithMaskConfig - configuration for ReplaceReelWithMask
type ReplaceReelWithMaskConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Symbol               string `yaml:"symbol" json:"symbol"`
	SymbolCode           int    `yaml:"-" json:"-"`
	Mask                 string `yaml:"mask" json:"mask"`
}

// SetLinkComponent
func (cfg *ReplaceReelWithMaskConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	}
}

type ReplaceReelWithMask struct {
	*BasicComponent `json:"-"`
	Config          *ReplaceReelWithMaskConfig `json:"config"`
}

// Init -
func (replaceReelWithMask *ReplaceReelWithMask) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReplaceReelWithMask.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReplaceReelWithMaskConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReplaceReelWithMask.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return replaceReelWithMask.InitEx(cfg, pool)
}

// InitEx -
func (replaceReelWithMask *ReplaceReelWithMask) InitEx(cfg any, pool *GamePropertyPool) error {
	replaceReelWithMask.Config = cfg.(*ReplaceReelWithMaskConfig)
	replaceReelWithMask.Config.ComponentType = ReplaceReelWithMaskTypeName

	sc, isok := pool.DefaultPaytables.MapSymbols[replaceReelWithMask.Config.Symbol]
	if !isok {
		goutils.Error("ReplaceReelWithMask.InitEx:Symbol",
			zap.String("symbol", replaceReelWithMask.Config.Symbol),
			zap.Error(ErrInvalidSymbol))

		return ErrInvalidSymbol
	}

	replaceReelWithMask.Config.SymbolCode = sc

	replaceReelWithMask.onInit(&replaceReelWithMask.Config.BasicComponentConfig)

	return nil
}

// playgame
func (replaceReelWithMask *ReplaceReelWithMask) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) error {

	replaceReelWithMask.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := replaceReelWithMask.GetTargetScene3(gameProp, curpr, prs, cd, replaceReelWithMask.Name, "", 0)
	ngs := gs

	maskVal, err := gameProp.GetMask(replaceReelWithMask.Config.Mask)
	if err != nil {
		goutils.Error("ReplaceReelWithMask.OnPlayGame:GetMask",
			zap.Error(err))

		return err
	}

	for x, v := range maskVal {
		if v {
			if ngs == gs {
				ngs = gs.CloneEx(gameProp.PoolScene)
			}

			arr := ngs.Arr[x]
			for y := range arr {
				ngs.Arr[x][y] = replaceReelWithMask.Config.SymbolCode
			}
		}
	}

	if ngs == gs {
		replaceReelWithMask.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	replaceReelWithMask.AddScene(gameProp, curpr, ngs, cd)

	replaceReelWithMask.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (replaceReel *ReplaceReelWithMask) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after replaceReelWithMask", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (replaceReel *ReplaceReelWithMask) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewReplaceReelWithMask(name string) IComponent {
	return &ReplaceReelWithMask{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"targetSymbols": "J",
//		"srcMask": "fg-bookof"
//	},
type jsonReplaceReelWithMask struct {
	TargetSymbols string `json:"targetSymbols"`
	SrcMask       string `json:"srcMask"`
}

func (jcfg *jsonReplaceReelWithMask) build() *ReplaceReelWithMaskConfig {
	cfg := &ReplaceReelWithMaskConfig{
		Symbol: jcfg.TargetSymbols,
		Mask:   jcfg.SrcMask,
	}

	cfg.UseSceneV3 = true

	return cfg
}

func parseReplaceReelWithMask(gamecfg *Config, cell *ast.Node) (string, error) {
	cfg, label, _, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseReplaceReelWithMask:getConfigInCell",
			zap.Error(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseReplaceReelWithMask:MarshalJSON",
			zap.Error(err))

		return "", err
	}

	data := &jsonReplaceReelWithMask{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseReplaceReelWithMask:Unmarshal",
			zap.Error(err))

		return "", err
	}

	cfgd := data.build()

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: ReplaceReelWithMaskTypeName,
	}

	gamecfg.GameMods[0].Components = append(gamecfg.GameMods[0].Components, ccfg)

	return label, nil
}
