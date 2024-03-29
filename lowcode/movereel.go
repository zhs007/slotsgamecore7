package lowcode

import (
	"log/slog"
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"gopkg.in/yaml.v2"
)

const MoveReelTypeName = "moveReel"

// MoveReelConfig - configuration for MoveReel
type MoveReelConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	MoveReelIndex        []int `yaml:"moveReelIndex" json:"moveReelIndex"`           // 每个轴的移动幅度，-1是上移
	EmptyOtherSceneVal   int   `yaml:"emptyOtherSceneVal" json:"emptyOtherSceneVal"` // 如果要移动otherscene时，这个是移出去以后的默认值
}

type MoveReel struct {
	*BasicComponent `json:"-"`
	Config          *MoveReelConfig `json:"config"`
}

// Init -
func (moveReel *MoveReel) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MoveReel.Init:ReadFile",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &MoveReelConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MoveReel.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	return moveReel.InitEx(cfg, pool)
}

// InitEx -
func (moveReel *MoveReel) InitEx(cfg any, pool *GamePropertyPool) error {
	moveReel.Config = cfg.(*MoveReelConfig)
	moveReel.Config.ComponentType = MoveReelTypeName

	moveReel.onInit(&moveReel.Config.BasicComponentConfig)

	return nil
}

// playgame
func (moveReel *MoveReel) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// moveReel.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	gs := moveReel.GetTargetScene3(gameProp, curpr, prs, 0)

	sc2 := gs.CloneEx(gameProp.PoolScene)

	for x, v := range moveReel.Config.MoveReelIndex {
		sc2.ResetReelIndex2(gameProp.Pool.Config.MapReels[sc2.ReelName], x, sc2.Indexes[x]+v)
	}

	moveReel.AddScene(gameProp, curpr, sc2, bcd)

	os := moveReel.GetTargetOtherScene3(gameProp, curpr, prs, 0)

	if os != nil {
		os2 := os.CloneEx(gameProp.PoolScene)

		for x, v := range moveReel.Config.MoveReelIndex {
			if v == 0 {
				continue
			}

			if v < 0 {
				v = -v

				for y := 0; y < os2.Height; y++ {
					if y < v {
						os2.Arr[x][y] = moveReel.Config.EmptyOtherSceneVal
					} else if y-v < os2.Height {
						os2.Arr[x][y] = os.Arr[x][y-v]
					}
				}
			} else {
				for y := 0; y < os2.Height; y++ {
					if y+v < os2.Height {
						os2.Arr[x][y] = os.Arr[x][y+v]
					} else {
						os2.Arr[x][y] = moveReel.Config.EmptyOtherSceneVal
					}
				}
			}
		}

		moveReel.AddOtherScene(gameProp, curpr, os2, bcd)
	}

	nc := moveReel.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (moveReel *MoveReel) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after moveReel", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (moveReel *MoveReel) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

func NewMoveReel(name string) IComponent {
	return &MoveReel{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
