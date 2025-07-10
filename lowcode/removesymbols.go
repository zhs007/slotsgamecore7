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

type RemoveSymbolsType int

const (
	RSTypeBasic       RemoveSymbolsType = 0
	RSTypeAdjacentPay RemoveSymbolsType = 1
)

func parseRemoveSymbolsType(strType string) RemoveSymbolsType {
	if strType == "adjacentPay" {
		return RSTypeAdjacentPay
	}

	return RSTypeBasic
}

type RemoveSymbolsData struct {
	BasicComponentData
	RemovedNum int
	AvgHeight  int // 平均移除图标的高度，用int表示浮点数，因此100表示1
}

// GetValEx -
func (removeSymbolsData *RemoveSymbolsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
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
	StrType              string            `yaml:"type" json:"type"`
	Type                 RemoveSymbolsType `yaml:"-" json:"-"`
	AddedSymbol          string            `yaml:"addedSymbol" json:"addedSymbol"`
	AddedSymbolCode      int               `yaml:"-" json:"-"`
	JumpToComponent      string            `yaml:"jumpToComponent" json:"jumpToComponent"`           // jump to
	TargetComponents     []string          `yaml:"targetComponents" json:"targetComponents"`         // 这些组件的中奖会需要参与remove
	IgnoreSymbols        []string          `yaml:"ignoreSymbols" json:"ignoreSymbols"`               // 忽略的symbol
	IgnoreSymbolCodes    []int             `yaml:"-" json:"-"`                                       // 忽略的symbol
	IsNeedProcSymbolVals bool              `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int               `yaml:"emptySymbolVal" json:"emptySymbolVal"`             // 空的symbolVal是什么
	OutputToComponent    string            `yaml:"outputToComponent" json:"outputToComponent"`       // outputToComponent
	Awards               []*Award          `yaml:"awards" json:"awards"`                             // 新的奖励系统
}

// SetLinkComponent
func (cfg *RemoveSymbolsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
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

	removeSymbols.Config.AddedSymbolCode = pool.DefaultPaytables.MapSymbols[removeSymbols.Config.AddedSymbol]
	removeSymbols.Config.Type = parseRemoveSymbolsType(removeSymbols.Config.StrType)

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
		return true
	}

	if len(removeSymbols.Config.IgnoreSymbolCodes) > 0 {
		return goutils.IndexOfIntSlice(removeSymbols.Config.IgnoreSymbolCodes, curs, 0) < 0
	}

	return true
}

// onBasic
func (removeSymbols *RemoveSymbols) onBasic(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, rscd *RemoveSymbolsData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) error {
	ngs := gs
	totalHeight := 0

	var outputCD IComponentData
	if removeSymbols.Config.OutputToComponent != "" {
		outputCD = gameProp.GetCurComponentDataWithName(removeSymbols.Config.OutputToComponent)
	}

	if os != nil {
		nos := os

		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn)
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

							if outputCD != nil {
								outputCD.AddPos(x, y)
							}

							rscd.RemovedNum++
						}
					}
				}
			}
		}

		if ngs == gs {
			return ErrComponentDoNothing
		}

		removeSymbols.AddOtherScene(gameProp, curpr, nos, &rscd.BasicComponentData)
	} else {
		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn)
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

							if outputCD != nil {
								outputCD.AddPos(x, y)
							}

							rscd.RemovedNum++
						}
					}

				}
			}
		}

		if ngs == gs {
			return ErrComponentDoNothing
		}
	}

	if !gIsReleaseMode {
		rscd.AvgHeight = totalHeight * 100 / rscd.RemovedNum
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &rscd.BasicComponentData)

	return nil
}

// onAdjacentPay
func (removeSymbols *RemoveSymbols) onAdjacentPay(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, rscd *RemoveSymbolsData,
	gs *sgc7game.GameScene, os *sgc7game.GameScene) error {
	ngs := gs
	totalHeight := 0
	npos := []int{}

	var outputCD IComponentData
	if removeSymbols.Config.OutputToComponent != "" {
		outputCD = gameProp.GetCurComponentDataWithName(removeSymbols.Config.OutputToComponent)
	}

	if os != nil {
		nos := os

		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn)
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

							if outputCD != nil {
								outputCD.AddPos(x, y)
							}

							isNeedRMOtherScene := true
							if len(curpr.Results[ri].Pos)/2 == 3 && pi == 1 {
								isNeedRMOtherScene = false
							} else if len(curpr.Results[ri].Pos)/2 == 5 && pi == 2 {
								isNeedRMOtherScene = false
							}

							if isNeedRMOtherScene {
								ngs.Arr[x][y] = -1
								nos.Arr[x][y] = removeSymbols.Config.EmptySymbolVal
							} else {
								npos = append(npos, x, y)
							}

							rscd.RemovedNum++
						}
					}
				}
			}
		}

		for pi := 0; pi < len(npos)/2; pi++ {
			ngs.Arr[npos[pi*2]][npos[pi*2+1]] = removeSymbols.Config.AddedSymbolCode
		}

		if ngs == gs {
			return ErrComponentDoNothing
		}

		removeSymbols.AddOtherScene(gameProp, curpr, nos, &rscd.BasicComponentData)
	} else {
		for _, cn := range removeSymbols.Config.TargetComponents {
			// 如果前面没有执行过，就可能没有清理数据，所以这里需要跳过
			if goutils.IndexOfStringSlice(gp.HistoryComponents, cn, 0) < 0 {
				continue
			}

			ccd := gameProp.GetCurComponentDataWithName(cn)
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

							if outputCD != nil {
								outputCD.AddPos(x, y)
							}

							isNeedRMOtherScene := true
							if len(curpr.Results[ri].Pos)/2 == 3 && pi == 1 {
								isNeedRMOtherScene = false
							} else if len(curpr.Results[ri].Pos)/2 == 5 && pi == 2 {
								isNeedRMOtherScene = false
							}

							if isNeedRMOtherScene {
								ngs.Arr[x][y] = -1
							} else {
								npos = append(npos, x, y)
							}

							rscd.RemovedNum++
						}
					}
				}
			}
		}

		for pi := 0; pi < len(npos)/2; pi++ {
			ngs.Arr[npos[pi*2]][npos[pi*2+1]] = removeSymbols.Config.AddedSymbolCode
		}

		if ngs == gs {
			return ErrComponentDoNothing
		}
	}

	if !gIsReleaseMode {
		rscd.AvgHeight = totalHeight * 100 / rscd.RemovedNum
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &rscd.BasicComponentData)

	return nil
}

// playgame
func (removeSymbols *RemoveSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {

	rscd := cd.(*RemoveSymbolsData)
	rscd.onNewStep()

	gs := removeSymbols.GetTargetScene3(gameProp, curpr, prs, 0)
	rscd.RemovedNum = 0

	var os *sgc7game.GameScene
	if removeSymbols.Config.IsNeedProcSymbolVals {
		os = removeSymbols.GetTargetOtherScene3(gameProp, curpr, prs, 0)
	}

	switch removeSymbols.Config.Type {
	case RSTypeBasic:
		err := removeSymbols.onBasic(gameProp, curpr, gp, rscd, gs, os)
		if err == ErrComponentDoNothing {
			nc := removeSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		} else if err != nil {
			goutils.Error("RemoveSymbols.OnPlayGame:onBasic",
				goutils.Err(err))

			return "", err
		}
	case RSTypeAdjacentPay:
		err := removeSymbols.onAdjacentPay(gameProp, curpr, gp, rscd, gs, os)
		if err == ErrComponentDoNothing {
			nc := removeSymbols.onStepEnd(gameProp, curpr, gp, "")

			return nc, ErrComponentDoNothing
		} else if err != nil {
			goutils.Error("RemoveSymbols.OnPlayGame:onAdjacentPay",
				goutils.Err(err))

			return "", err
		}
	default:
		goutils.Error("RemoveSymbols.OnPlayGame:InvalidType",
			slog.Int("type", int(removeSymbols.Config.Type)),
			goutils.Err(ErrIvalidComponentConfig))

		return "", ErrIvalidComponentConfig
	}

	removeSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := removeSymbols.onStepEnd(gameProp, curpr, gp, removeSymbols.Config.JumpToComponent)

	return nc, nil
}

// OnProcControllers -
func (removeSymbols *RemoveSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(removeSymbols.Config.Awards) > 0 {
		gameProp.procAwards(plugin, removeSymbols.Config.Awards, curpr, gp)
	}
}

// OnAsciiGame - outpur to asciigame
func (removeSymbols *RemoveSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*RemoveSymbolsData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after removedSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

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

// "isNeedProcSymbolVals": false,
// "emptySymbolVal": -1,
// "targetComponents": [
//
//	"bg-pay"
//
// ],
// "type": "adjacentPay",
// "addedSymbol": "WL"
// "outputToComp": "bg-pos-rmoved"
type jsonRemoveSymbols struct {
	Type                 string   `json:"type"`                 // type
	AddedSymbol          string   `json:"addedSymbol"`          // addedSymbol
	TargetComponents     []string `json:"targetComponents"`     // 这些组件的中奖会需要参与remove
	IgnoreSymbols        []string `json:"ignoreSymbols"`        // 忽略的symbol
	IsNeedProcSymbolVals bool     `json:"isNeedProcSymbolVals"` // 是否需要同时处理symbolVals
	EmptySymbolVal       int      `json:"emptySymbolVal"`       // 空的symbolVal是什么
	OutputToComponent    string   `json:"outputToComp"`         // outputToComp
}

func (jcfg *jsonRemoveSymbols) build() *RemoveSymbolsConfig {
	cfg := &RemoveSymbolsConfig{
		StrType:              jcfg.Type,
		TargetComponents:     jcfg.TargetComponents,
		IgnoreSymbols:        jcfg.IgnoreSymbols,
		IsNeedProcSymbolVals: jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:       jcfg.EmptySymbolVal,
		AddedSymbol:          jcfg.AddedSymbol,
		OutputToComponent:    jcfg.OutputToComponent,
	}

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
