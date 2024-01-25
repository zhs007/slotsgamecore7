package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
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

type ReplaceSymbolGroup struct {
	*BasicComponent `json:"-"`
	Config          *ReplaceSymbolGroupConfig `json:"config"`
}

// Init -
func (replaceSymbolGroup *ReplaceSymbolGroup) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReplaceSymbolGroup.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReplaceSymbolGroupConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReplaceSymbolGroup.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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
			zap.Int("src", len(replaceSymbolGroup.Config.SrcSymbolCodes)),
			zap.Int("target", len(replaceSymbolGroup.Config.TargetSymbolCodes)),
			zap.Error(ErrIvalidComponentConfig))

		return ErrIvalidComponentConfig
	}

	replaceSymbolGroup.onInit(&replaceSymbolGroup.Config.BasicComponentConfig)

	return nil
}

// playgame
func (replaceSymbolGroup *ReplaceSymbolGroup) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	replaceSymbolGroup.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	if len(replaceSymbolGroup.Config.SrcSymbolCodes) == 0 {
		replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	cd := gameProp.MapComponentData[replaceSymbolGroup.Name].(*BasicComponentData)

	gs := replaceSymbolGroup.GetTargetScene3(gameProp, curpr, cd, replaceSymbolGroup.Name, "", 0)
	ngs := gs

	for x, arr := range gs.Arr {
		for y, srcSymbol := range arr {
			si := goutils.IndexOfIntSlice(replaceSymbolGroup.Config.SrcSymbolCodes, srcSymbol, 0)
			if si >= 0 {
				if ngs == gs {
					ngs = gs.Clone()
				}

				ngs.Arr[x][y] = replaceSymbolGroup.Config.TargetSymbolCodes[si]
			}
		}
	}

	if ngs == gs {
		replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	replaceSymbolGroup.AddScene(gameProp, curpr, ngs, cd)

	replaceSymbolGroup.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (replaceSymbolGroup *ReplaceSymbolGroup) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[replaceSymbolGroup.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after replaceSymbolGroup", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (replaceSymbolGroup *ReplaceSymbolGroup) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewReplaceSymbolGroup(name string) IComponent {
	return &ReplaceSymbolGroup{
		BasicComponent: NewBasicComponent(name, 1),
	}
}