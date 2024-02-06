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

const ReplaceReelTypeName = "replaceReel"

// ReplaceReelConfig - configuration for ReplaceReel
type ReplaceReelConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MapReels             map[int]string `yaml:"mapReels" json:"mapReels"`
	MapReelsCode         map[int]int    `yaml:"-" json:"-"`
}

type ReplaceReel struct {
	*BasicComponent `json:"-"`
	Config          *ReplaceReelConfig `json:"config"`
}

// Init -
func (replaceReel *ReplaceReel) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReplaceReel.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReplaceReelConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReplaceReel.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return replaceReel.InitEx(cfg, pool)
}

// InitEx -
func (replaceReel *ReplaceReel) InitEx(cfg any, pool *GamePropertyPool) error {
	replaceReel.Config = cfg.(*ReplaceReelConfig)
	replaceReel.Config.ComponentType = ReplaceReelTypeName

	replaceReel.Config.MapReelsCode = make(map[int]int)

	for k, v := range replaceReel.Config.MapReels {
		sc, isok := pool.DefaultPaytables.MapSymbols[v]
		if !isok {
			goutils.Error("ReplaceReel.InitEx:MapReels",
				zap.String("symbol", v),
				zap.Error(ErrInvalidSymbol))

			return ErrInvalidSymbol
		}

		replaceReel.Config.MapReelsCode[k] = sc
	}

	replaceReel.onInit(&replaceReel.Config.BasicComponentConfig)

	return nil
}

// playgame
func (replaceReel *ReplaceReel) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, icd IComponentData) (string, error) {

	// replaceReel.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := icd.(*BasicComponentData)

	gs := replaceReel.GetTargetScene3(gameProp, curpr, prs, cd, replaceReel.Name, "", 0)

	sc2 := gs.CloneEx(gameProp.PoolScene)

	for x, target := range replaceReel.Config.MapReelsCode {
		arr := sc2.Arr[x]
		for y := range arr {
			sc2.Arr[x][y] = target
		}
	}

	replaceReel.AddScene(gameProp, curpr, sc2, cd)

	nc := replaceReel.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (replaceReel *ReplaceReel) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, icd IComponentData) error {
	cd := icd.(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after replaceReel", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (replaceReel *ReplaceReel) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewReplaceReel(name string) IComponent {
	return &ReplaceReel{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
