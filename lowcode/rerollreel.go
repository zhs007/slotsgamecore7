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

const ReRollReelTypeName = "reRollReel"

// ReRollReelConfig - configuration for ReRollReel
type ReRollReelConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
}

type ReRollReel struct {
	*BasicComponent `json:"-"`
	Config          *ReRollReelConfig `json:"config"`
}

// Init -
func (reRollReel *ReRollReel) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("ReRollReel.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &ReRollReelConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("ReRollReel.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return reRollReel.InitEx(cfg, pool)
}

// InitEx -
func (reRollReel *ReRollReel) InitEx(cfg any, pool *GamePropertyPool) error {
	reRollReel.Config = cfg.(*ReRollReelConfig)
	reRollReel.Config.ComponentType = MoveReelTypeName

	reRollReel.onInit(&reRollReel.Config.BasicComponentConfig)

	return nil
}

// playgame
func (reRollReel *ReRollReel) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	reRollReel.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[reRollReel.Name].(*BasicComponentData)

	gs := reRollReel.GetTargetScene2(gameProp, curpr, cd, reRollReel.Name, "")

	sc2 := gs.CloneEx(gameProp.PoolScene)

	sc2.RandReelsWithReelData(gameProp.Pool.Config.MapReels[sc2.ReelName], plugin)

	reRollReel.AddScene(gameProp, curpr, sc2, cd)

	reRollReel.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (reRollReel *ReRollReel) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[reRollReel.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after reRollReel", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (reRollReel *ReRollReel) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewReRollReel(name string) IComponent {
	return &ReRollReel{
		BasicComponent: NewBasicComponent(name),
	}
}
