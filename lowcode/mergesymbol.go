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

const MergeSymbolTypeName = "mergeSymbol"

// MergeSymbolConfig - configuration for MergeSymbol
type MergeSymbolConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	SrcScene             []string `yaml:"srcScene" json:"srcScene"`                     // 2个scene，mask false表示用0，true表示用1
	TargetMask           string   `yaml:"targetMask" json:"targetMask"`                 // mask
	EmptyOtherSceneVal   int      `yaml:"emptyOtherSceneVal" json:"emptyOtherSceneVal"` // 如果要合并otherscene时，某一个otherscene不存在时，就用这个作默认值
}

type MergeSymbol struct {
	*BasicComponent `json:"-"`
	Config          *MergeSymbolConfig `json:"config"`
}

// Init -
func (mergeSymbol *MergeSymbol) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("MergeSymbol.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &MergeSymbolConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("MergeSymbol.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return mergeSymbol.InitEx(cfg, pool)
}

// InitEx -
func (mergeSymbol *MergeSymbol) InitEx(cfg any, pool *GamePropertyPool) error {
	mergeSymbol.Config = cfg.(*MergeSymbolConfig)
	mergeSymbol.Config.ComponentType = MoveReelTypeName

	mergeSymbol.onInit(&mergeSymbol.Config.BasicComponentConfig)

	return nil
}

// playgame
func (mergeSymbol *MergeSymbol) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// mergeSymbol.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*BasicComponentData)

	gs1 := mergeSymbol.GetTargetScene3(gameProp, curpr, prs, bcd, mergeSymbol.Name, mergeSymbol.Config.SrcScene[0], 0)
	gs2 := mergeSymbol.GetTargetScene3(gameProp, curpr, prs, bcd, mergeSymbol.Name+":1", mergeSymbol.Config.SrcScene[1], 1)

	sc2 := gs1.CloneEx(gameProp.PoolScene)

	mask, err := gameProp.Pool.GetMask(mergeSymbol.Config.TargetMask, gameProp)
	if err != nil {
		goutils.Error("MergeSymbol.OnPlayGame:GetMask",
			zap.Error(err))

		return "", err
	}

	for x, arr := range gs2.Arr {
		if mask[x] {
			copy(sc2.Arr[x], arr)
			// for y, s := range arr {
			// 	sc2.Arr[x][y] = s
			// }

			sc2.Indexes[x] = gs2.Indexes[x]
		}
	}

	mergeSymbol.AddScene(gameProp, curpr, sc2, bcd)

	os1 := mergeSymbol.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	os2 := mergeSymbol.GetTargetOtherScene3(gameProp, curpr, prs, 1)
	if os1 != nil || os2 != nil {
		var os3 *sgc7game.GameScene

		mask, err := gameProp.Pool.GetMask(mergeSymbol.Config.TargetMask, gameProp)
		if err != nil {
			goutils.Error("MergeSymbol.OnPlayGame:GetMask",
				zap.Error(err))

			return "", err
		}

		if os1 == nil {
			os3 = os2.CloneEx(gameProp.PoolScene)

			for x, arr := range os3.Arr {
				if !mask[x] {
					for y := range arr {
						os3.Arr[x][y] = mergeSymbol.Config.EmptyOtherSceneVal
					}
				}
			}
		} else if os2 == nil {
			os3 = os1.CloneEx(gameProp.PoolScene)

			for x, arr := range os3.Arr {
				if mask[x] {
					for y := range arr {
						os3.Arr[x][y] = mergeSymbol.Config.EmptyOtherSceneVal
					}
				}
			}
		} else {
			os3 = os1.CloneEx(gameProp.PoolScene)

			for x, arr := range os2.Arr {
				if mask[x] {
					copy(os3.Arr[x], arr)
				}
			}
		}

		mergeSymbol.AddOtherScene(gameProp, curpr, os3, bcd)
	}

	nc := mergeSymbol.onStepEnd(gameProp, curpr, gp, "")

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (mergeSymbol *MergeSymbol) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*BasicComponentData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after mergeSymbol", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (mergeSymbol *MergeSymbol) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewMergeSymbol(name string) IComponent {
	return &MergeSymbol{
		BasicComponent: NewBasicComponent(name, 2),
	}
}
