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
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"gopkg.in/yaml.v2"
)

const RemoveSymbolsTypeName = "removeSymbols"

type RemoveSymbolsData struct {
	BasicComponentData
	RemovedNum int
	AvgHeight  int // 平均移除图标的高度，用int表示浮点数，因此100表示1
}

// GetVal -
func (removeSymbolsData *RemoveSymbolsData) GetVal(key string) (int, bool) {
	if key == CVAvgHeight {
		return removeSymbolsData.AvgHeight, true
	}

	return 0, false
}

// OnNewGame -
func (removeSymbolsData *RemoveSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	removeSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep -
func (removeSymbolsData *RemoveSymbolsData) onNewStep() {
	// removeSymbolsData.BasicComponentData.OnNewStep(gameProp, component)

	removeSymbolsData.RemovedNum = 0
	removeSymbolsData.UsedScenes = nil
	removeSymbolsData.UsedOtherScenes = nil

	if gIsReleaseMode {
		removeSymbolsData.AvgHeight = 0
	}
}

// Clone
func (removeSymbolsData *RemoveSymbolsData) Clone() IComponentData {
	target := &RemoveSymbolsData{
		BasicComponentData: removeSymbolsData.CloneBasicComponentData(),
		RemovedNum:         removeSymbolsData.RemovedNum,
	}

	return target
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
	JumpToComponent      string   `yaml:"jumpToComponent" json:"jumpToComponent"`           // jump to
	TargetComponents     []string `yaml:"targetComponents" json:"targetComponents"`         // 这些组件的中奖会需要参与remove
	IgnoreSymbols        []string `yaml:"ignoreSymbols" json:"ignoreSymbols"`               // 忽略的symbol
	IgnoreSymbolCodes    []int    `yaml:"-" json:"-"`                                       // 忽略的symbol
	IsNeedProcSymbolVals bool     `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	Awards               []*Award `yaml:"awards" json:"awards"`                             // 新的奖励系统
}

// SetLinkComponent
func (cfg *RemoveSymbolsConfig) SetLinkComponent(link string, componentName string) {
	if link == "next" {
		cfg.DefaultNextComponent = componentName
	} else if link == "jump" {
		cfg.JumpToComponent = componentName
	}
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
			slog.String("fn", fn),
			goutils.Err(err))

		return err
	}

	cfg := &RemoveSymbolsConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		goutils.Error("RemoveSymbols.Init:Unmarshal",
			slog.String("fn", fn),
			goutils.Err(err))

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

	for _, v := range removeSymbols.Config.Awards {
		v.Init()
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
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	// removeSymbols.onPlayGame(gameProp, curpr, gp, plugin, cmd, param, ps, stake, prs)

	bcd := cd.(*RemoveSymbolsData)
	bcd.onNewStep()

	gs := removeSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	ngs := gs

	bcd.RemovedNum = 0

	totalHeight := 0

	var os *sgc7game.GameScene
	if removeSymbols.Config.IsNeedProcSymbolVals {
		os = removeSymbols.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	if os != nil {
		nos := os

		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn) //gameProp.MapComponentData[cn]
			if ccd != nil {
				lst := ccd.GetResults()
				for _, ri := range lst {
					for pi := 0; pi < len(curpr.Results[ri].Pos)/2; pi++ {
						x := curpr.Results[ri].Pos[pi*2]
						y := curpr.Results[ri].Pos[pi*2+1]
						if removeSymbols.canRemove(x, y, ngs) {
							if ngs == gs {
								ngs = gs.CloneEx(gameProp.PoolScene)
								nos = os.CloneEx(gameProp.PoolScene)
							}

							if !gIsReleaseMode {
								totalHeight += y
							}

							ngs.Arr[x][y] = -1
							nos.Arr[x][y] = removeSymbols.Config.EmptySymbolVal

							bcd.RemovedNum++
						}
					}
				}
			}
		}

		if ngs == gs {
			nc := removeSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}

		removeSymbols.AddOtherScene(gameProp, curpr, nos, &bcd.BasicComponentData)
	} else {
		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn) //gameProp.MapComponentData[cn]
			if ccd != nil {
				lst := ccd.GetResults()
				for _, ri := range lst {
					for pi := 0; pi < len(curpr.Results[ri].Pos)/2; pi++ {
						x := curpr.Results[ri].Pos[pi*2]
						y := curpr.Results[ri].Pos[pi*2+1]
						if removeSymbols.canRemove(x, y, ngs) {
							if ngs == gs {
								ngs = gs.CloneEx(gameProp.PoolScene)
							}

							if !gIsReleaseMode {
								totalHeight += y
							}

							ngs.Arr[x][y] = -1

							bcd.RemovedNum++
						}
					}
				}
			}
		}

		if ngs == gs {
			nc := removeSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		}
	}

	if !gIsReleaseMode {
		bcd.AvgHeight = totalHeight * 100 / bcd.RemovedNum
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &bcd.BasicComponentData)

	if len(removeSymbols.Config.Awards) > 0 {
		gameProp.procAwards(plugin, removeSymbols.Config.Awards, curpr, gp)
	}

	nc := removeSymbols.onStepEnd(gameProp, curpr, gp, removeSymbols.Config.JumpToComponent)

	return nc, nil
}

// OnAsciiGame - outpur to asciigame
func (removeSymbols *RemoveSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*RemoveSymbolsData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after removedSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// // OnStats
// func (removeSymbols *RemoveSymbols) OnStats(feature *sgc7stats.Feature, stake *sgc7game.Stake, lst []*sgc7game.PlayResult) (bool, int64, int64) {
// 	return false, 0, 0
// }

// // OnStatsWithPB -
// func (removeSymbols *RemoveSymbols) OnStatsWithPB(feature *sgc7stats.Feature, pbComponentData proto.Message, pr *sgc7game.PlayResult) (int64, error) {
// 	return 0, nil
// }

// NewComponentData -
func (removeSymbols *RemoveSymbols) NewComponentData() IComponentData {
	return &RemoveSymbolsData{}
}

// EachUsedResults -
func (removeSymbols *RemoveSymbols) EachUsedResults(pr *sgc7game.PlayResult, pbComponentData *anypb.Any, oneach FuncOnEachUsedResult) {
}

// GetAllLinkComponents - get all link components
func (removeSymbols *RemoveSymbols) GetAllLinkComponents() []string {
	return []string{removeSymbols.Config.DefaultNextComponent, removeSymbols.Config.JumpToComponent}
}

// GetNextLinkComponents - get next link components
func (removeSymbols *RemoveSymbols) GetNextLinkComponents() []string {
	return []string{removeSymbols.Config.DefaultNextComponent, removeSymbols.Config.JumpToComponent}
}

func NewRemoveSymbols(name string) IComponent {
	return &RemoveSymbols{
		BasicComponent: NewBasicComponent(name, 1),
	}
}

//	"configuration": {
//		"targetComponents": [
//			"bg-payblue"
//		]
//	},
type jsonRemoveSymbols struct {
	TargetComponents     []string `json:"targetComponents"`                                 // 这些组件的中奖会需要参与remove
	IgnoreSymbols        []string `json:"ignoreSymbols"`                                    // 忽略的symbol
	IsNeedProcSymbolVals bool     `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
}

func (jcfg *jsonRemoveSymbols) build() *RemoveSymbolsConfig {
	cfg := &RemoveSymbolsConfig{
		TargetComponents:     jcfg.TargetComponents,
		IgnoreSymbols:        jcfg.IgnoreSymbols,
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
	}

	// cfg.UseSceneV3 = true

	return cfg
}

func parseRemoveSymbols(gamecfg *BetConfig, cell *ast.Node) (string, error) {
	cfg, label, ctrls, err := getConfigInCell(cell)
	if err != nil {
		goutils.Error("parseRemoveSymbols:getConfigInCell",
			goutils.Err(err))

		return "", err
	}

	buf, err := cfg.MarshalJSON()
	if err != nil {
		goutils.Error("parseRemoveSymbols:MarshalJSON",
			goutils.Err(err))

		return "", err
	}

	data := &jsonRemoveSymbols{}

	err = sonic.Unmarshal(buf, data)
	if err != nil {
		goutils.Error("parseRemoveSymbols:Unmarshal",
			goutils.Err(err))

		return "", err
	}

	cfgd := data.build()

	if ctrls != nil {
		awards, err := parseControllers(ctrls)
		if err != nil {
			goutils.Error("parseRemoveSymbols:parseControllers",
				goutils.Err(err))

			return "", err
		}

		cfgd.Awards = awards
	}

	gamecfg.mapConfig[label] = cfgd
	gamecfg.mapBasicConfig[label] = &cfgd.BasicComponentConfig

	ccfg := &ComponentConfig{
		Name: label,
		Type: RemoveSymbolsTypeName,
	}

	gamecfg.Components = append(gamecfg.Components, ccfg)

	return label, nil
}
