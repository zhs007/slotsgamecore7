package lowcode

import (
	"os"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/asciigame"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RemoveSymbolsTypeName = "removeSymbols"

type RemoveSymbolsData struct {
	BasicComponentData
	RemovedNum int
}

// OnNewGame -
func (removeSymbolsData *RemoveSymbolsData) OnNewGame() {
	removeSymbolsData.BasicComponentData.OnNewGame()
}

// OnNewStep -
func (removeSymbolsData *RemoveSymbolsData) OnNewStep() {
	removeSymbolsData.BasicComponentData.OnNewStep()

	removeSymbolsData.RemovedNum = 0
}

// BuildPBComponentData
func (removeSymbolsData *RemoveSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RemoveSymbolsData{
		BasicComponentData: removeSymbolsData.BuildPBBasicComponentData(),
		RemovedNum:         int32(removeSymbolsData.RemovedNum),
	}

	return pbcd
}

// RemoveSymbolsConfig - configuration for RemoveSymbols
type RemoveSymbolsConfig struct {
	BasicComponentConfig `yaml:",inline" json:",inline"`
	TargetComponents     []string `yaml:"targetComponents" json:"targetComponents"` // 这些组件的中奖会需要参与remove
	IgnoreSymbols        []string `yaml:"ignoreSymbols" json:"ignoreSymbols"`       // 忽略的symbol
	IgnoreSymbolCodes    []int    `yaml:"-" json:"-"`                               // 忽略的symbol
}

type RemoveSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RemoveSymbolsConfig `json:"config"`
}

// Init -
func (removeSymbols *RemoveSymbols) Init(fn string, pool *GamePropertyPool) error {
	data, err := os.ReadFile(fn)
	if err != nil {
		goutils.Error("RemoveSymbols.Init:ReadFile",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	cfg := &RemoveSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RemoveSymbols.Init:Unmarshal",
			zap.String("fn", fn),
			zap.Error(err))

		return err
	}

	return removeSymbols.InitEx(cfg, pool)
}

// InitEx -
func (removeSymbols *RemoveSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	removeSymbols.Config = cfg.(*RemoveSymbolsConfig)
	removeSymbols.Config.ComponentType = RemoveSymbolsTypeName

	for _, v := range removeSymbols.Config.IgnoreSymbols {
		removeSymbols.Config.IgnoreSymbolCodes = append(removeSymbols.Config.IgnoreSymbolCodes, pool.DefaultPaytables.MapSymbols[v])
	}

	removeSymbols.onInit(&removeSymbols.Config.BasicComponentConfig)

	return nil
}

func (removeSymbols *RemoveSymbols) canRemove(x, y int, gs *sgc7game.GameScene) bool {
	curs := gs.Arr[x][y]
	if curs < 0 {
		return false
	}

	if len(removeSymbols.Config.IgnoreSymbolCodes) > 0 {
		return goutils.IndexOfIntSlice(removeSymbols.Config.IgnoreSymbolCodes, curs, 0) < 0
	}

	return true
}

// playgame
func (removeSymbols *RemoveSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult) error {

	removeSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	cd := gameProp.MapComponentData[removeSymbols.Name].(*RemoveSymbolsData)

	gs := removeSymbols.GetTargetScene3(gameProp, curpr, &cd.BasicComponentData, removeSymbols.Name, "", 0)
	ngs := gs

	cd.RemovedNum = 0

	for _, cn := range removeSymbols.Config.TargetComponents {
		ccd := gameProp.MapComponentData[cn]
		lst := ccd.GetResults()
		for _, ri := range lst {
			for pi := 0; pi < len(curpr.Results[ri].Pos)/2; pi++ {
				x := curpr.Results[ri].Pos[pi*2]
				y := curpr.Results[ri].Pos[pi*2+1]
				if removeSymbols.canRemove(x, y, ngs) {
					if ngs == gs {
						ngs = gs.Clone()
					}

					ngs.Arr[x][y] = -1

					cd.RemovedNum++
				}
			}
		}
	}

	if ngs == gs {
		removeSymbols.onStepEnd(gameProp, curpr, gp, "")

		return ErrComponentDoNothing
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &cd.BasicComponentData)

	removeSymbols.onStepEnd(gameProp, curpr, gp, "")

	return nil
}

// OnAsciiGame - outpur to asciigame
func (removeSymbols *RemoveSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap) error {
	cd := gameProp.MapComponentData[removeSymbols.Name].(*RemoveSymbolsData)

	if len(cd.UsedScenes) > 0 {
		asciigame.OutputScene("after removedSymbols", pr.Scenes[cd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// OnStats
func (removeSymbols *RemoveSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
	return false, 0, 0
}

// OnStatsWithPB -
func (removeSymbols *RemoveSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
	return 0, nil
}

// NewComponentData -
func (removeSymbols *RemoveSymbols) NewComponentData() IComponentData {
	return &RemoveSymbolsData{}
}

// EachUsedResults -
func (removeSymbols *RemoveSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
	pbcd := &sgc7pb.RemoveSymbolsData{}

	err := pbComponentData.UnmarshalTo(pbcd)
	if err != nil {
		goutils.Error("RemoveSymbols.EachUsedResults:UnmarshalTo",
			zap.Error(err))

		return
	}

	for _, v := range pbcd.BasicComponentData.UsedResults {
		oneach(pr.Results[v])
	}
}

func NewRemoveSymbols(name string) IComponent {
	mystery := &RemoveSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}

	return mystery
}
