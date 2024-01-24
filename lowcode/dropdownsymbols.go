package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const DropDownSymbolsTypeName = "dropDownSymbols"

// DropDownSymbolsConfig - configuration for DropDownSymbols
type DropDownSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	HoldSymbols          []string `yaml:"holdSymbols" json:"holdSymbols"` // 不需要下落的symbol
	HoldSymbolCodes      []int    `yaml:"-" json:"-"`                     // 不需要下落的symbol
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
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &DropDownSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("DropDownSymbols.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

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

func (dropDownSymbols *DropDownSymbols) canDropDown(x, y int, gs *sgc7game.GameScene) bool {
	curs := gs.Arr[x][y]
	if curs < 0 {
		return false
	}

	if len(dropDownSymbols.Config.HoldSymbolCodes) > 0 {
		return goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, curs, 0) < 0
	}

	return true
}

func (dropDownSymbols *DropDownSymbols) isNeedDropDown(gs *sgc7game.GameScene) bool {
	for _, arr := range gs.Arr {
		for _, s := range arr {
			if s < 0 {
				return true
			}
		}
	}

	return false
}

// playgame
func (dropDownSymbols *DropDownSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	dropDownSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[dropDownSymbols.Name].(*BasicComponentData)

	gs := dropDownSymbols.GetTargetScene3(gameProp, curpr, cd, dropDownSymbols.Name, "", 0)
	ngs := gs

	for x, arr := range ngs.Arr {
		for y := len(arr) - 1; y >= 0; {
			if arr[y] == -1 {
				hass := false
				for y1 := y - 1; y1 >= 0; y1-- {
					if arr[y1] != -1 && goutils.IndexOfIntSlice(dropDownSymbols.Config.HoldSymbolCodes, ngs.Arr[x][y1], 0) < 0 {
						if ngs == gs {
							ngs = gs.Clone()

							arr = ngs.Arr[x]
						}

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

	if ngs == gs {
		dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	dropDownSymbols.AddScene(gameProp, curpr, ngs, cd)

	dropDownSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (dropDownSymbols *DropDownSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[dropDownSymbols.Name].(*RemoveSymbolsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after dropDownSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (dropDownSymbols *DropDownSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (dropDownSymbols *DropDownSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	return 0, nil
}

// EachUsedResults -
func (dropDownSymbols *DropDownSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

func NewDropDownSymbols(name string) IComponent {
	return &DropDownSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
