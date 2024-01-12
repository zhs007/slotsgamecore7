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

const FixSymbolsTypeName = "fixSymbols"

const (
	FixSymbolsTypeUnknow    int = 0 // unknow
	FixSymbolsTypeMergeDown int = 1 // merge & down
)

func parseFixSymbolsType(str string) int {
	if str == "mergedown" {
		return FixSymbolsTypeMergeDown
	}

	return FixSymbolsTypeUnknow
}

// FixSymbolsConfig - configuration for FixSymbols feature
type FixSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	Type                 string   `yaml:"type" json:"type"`
	Symbols              []string `yaml:"symbols" json:"symbols"`
}

type FixSymbols struct {
	*BasicComponent `json:"-"`
	Config          *FixSymbolsConfig `json:"config"`
	SymbolCodes     []int             `json:"-"`
	Type            int               `json:"-"`
}

// checkMergeDown -
func (fixSymbols *FixSymbols) isNeedMergeDownWithYArr(yarr []int) bool {
	if len(yarr) >= 2 {
		sy := yarr[0]
		for i := 1; i < len(yarr); i++ {
			if yarr[i] > sy+1 {
				return true
			}

			sy = yarr[i]
		}
	}

	return false
}

// checkMergeDown -
func (fixSymbols *FixSymbols) isNeedMergeDown(gs *sgc7game.GameScene) ([]int, [][]int) {
	xarr := make([]int, 0, len(gs.Arr))
	yarrs := make([][]int, 0, len(gs.Arr))

	for x, arr := range gs.Arr {
		yarr := make([]int, 0, len(arr))

		for y, s := range arr {
			if goutils.IndexOfIntSlice(fixSymbols.SymbolCodes, s, 0) >= 0 {
				yarr = append(yarr, y)
			}
		}

		if fixSymbols.isNeedMergeDownWithYArr(yarr) {
			xarr = append(xarr, x)
			yarrs = append(yarrs, yarr)
		}
	}

	return xarr, yarrs
}

// Init -
func (fixSymbols *FixSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("FixSymbols.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &FixSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("FixSymbols.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return fixSymbols.InitEx(cfg, pool)
}

// InitEx -
func (fixSymbols *FixSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	fixSymbols.Config = cfg.(*FixSymbolsConfig)
	fixSymbols.Config.ComponentType = FixSymbolsTypeName

	for _, v := range fixSymbols.Config.Symbols {
		fixSymbols.SymbolCodes = append(fixSymbols.SymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	fixSymbols.Type = parseFixSymbolsType(fixSymbols.Config.Type)

	fixSymbols.onInit(&fixSymbols.Config.BasicComponentConfig)

	return nil
}

// playgame
func (fixSymbols *FixSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	fixSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[fixSymbols.Name].(*BasicComponentData)

	gs := fixSymbols.GetTargetScene3(gameProp, curpr, cd, fixSymbols.Name, "", 0)

	needReTag := true
	if fixSymbols.Type == FixSymbolsTypeMergeDown {
		xarr, _ := fixSymbols.isNeedMergeDown(gs)
		if len(xarr) > 0 {
			ngs := gs.CloneEx(gameProp.PoolScene)
			// ngs := gs.Clone()

			// 3可以是个特例
			if len(gs.Arr[0]) == 3 {
				for _, x := range xarr {
					ngs.Arr[x][1] = ngs.Arr[x][0]
					ngs.Arr[x][0] = gs.Arr[x][1]
				}
			}

			fixSymbols.AddScene(gameProp, curpr, ngs, cd)

			needReTag = false
		}
	}

	if needReTag {
		fixSymbols.ReTagScene(gameProp, curpr, cd.TargetSceneIndex, cd)
	}

	fixSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (fixSymbols *FixSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {

	cd := gameProp.MapComponentData[fixSymbols.Name].(*BasicComponentData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("The value of the symbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (fixSymbols *FixSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

func NewFixSymbols(name string) IComponent {
	return &FixSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}
