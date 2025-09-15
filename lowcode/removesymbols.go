// Package lowcode contains configurable, pluggable "components" that modify
// game scenes and play results at runtime. Components are defined by YAML/JSON
// configurations and operate on GameScene objects to implement mechanics such as
// removing symbols, adding symbols, collecting positions, etc.
package lowcode

import (
	"log/slog"
	"os"
	"slices"

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

// RemoveSymbolsTypeName is the component type name used in configuration files
// to identify the removeSymbols component.
const RemoveSymbolsTypeName = "removeSymbols"

// RemoveSymbolsType enumerates different removal behaviors supported by the
// removeSymbols component.
type RemoveSymbolsType int

const (
	RSTypeBasic       RemoveSymbolsType = 0
	RSTypeAdjacentPay RemoveSymbolsType = 1
)

// parseRemoveSymbolsType parses a string representation of the component
// subtype and returns the corresponding RemoveSymbolsType. Unknown types
// default to RSTypeBasic.
func parseRemoveSymbolsType(strType string) RemoveSymbolsType {
	if strType == "adjacentPay" {
		return RSTypeAdjacentPay
	}

	return RSTypeBasic
}

// RemoveSymbolsData holds runtime data for a removeSymbols component.
// It embeds BasicComponentData and tracks how many symbols were removed in
// the current step and the average height of removed symbols (AvgHeight).
// AvgHeight is stored as an int fixed-point value where 100 represents 1.0.
type RemoveSymbolsData struct {
	BasicComponentData
	RemovedNum int
	AvgHeight  int // average removed symbol height; fixed-point: 100==1.0
}

// GetValEx returns the integer value associated with the provided key for
// this component data. It supports "CVAvgHeight" to expose the calculated
// average remove height.
func (removeSymbolsData *RemoveSymbolsData) GetValEx(key string, getType GetComponentValType) (int, bool) {
	if key == CVAvgHeight {
		return removeSymbolsData.AvgHeight, true
	}

	return 0, false
}

// OnNewGame is called when a new game starts and allows the component data
// to initialize or reset any per-game state.
func (removeSymbolsData *RemoveSymbolsData) OnNewGame(gameProp *GameProperty, component IComponent) {
	removeSymbolsData.BasicComponentData.OnNewGame(gameProp, component)
}

// onNewStep resets per-step transient state. It should be called at the
// beginning of each component execution step.
func (removeSymbolsData *RemoveSymbolsData) onNewStep() {
	removeSymbolsData.RemovedNum = 0
	removeSymbolsData.UsedScenes = nil
	removeSymbolsData.UsedOtherScenes = nil

	// Always initialize AvgHeight to avoid carrying old values between steps
	removeSymbolsData.AvgHeight = 0
}

// Clone creates a shallow copy of RemoveSymbolsData suitable for storing as
// independent component data on cloned branches.
func (removeSymbolsData *RemoveSymbolsData) Clone() IComponentData {
	target := &RemoveSymbolsData{
		BasicComponentData: removeSymbolsData.CloneBasicComponentData(),
		RemovedNum:         removeSymbolsData.RemovedNum,
	}

	return target
}

// BuildPBComponentData builds the protobuf representation of this component
// data for serialization or transport.
func (removeSymbolsData *RemoveSymbolsData) BuildPBComponentData() proto.Message {
	pbcd := &sgc7pb.RemoveSymbolsData{
		BasicComponentData: removeSymbolsData.BuildPBBasicComponentData(),
		RemovedNum:         int32(removeSymbolsData.RemovedNum),
	}

	return pbcd
}

// RemoveSymbolsConfig is the YAML/JSON configuration for a removeSymbols
// component. It defines targets, ignored symbols, output wiring and
// additional behavior flags.
type RemoveSymbolsConfig struct {
	BasicComponentConfig     `yaml:",inline" json:",inline"`
	StrType                  string            `yaml:"type" json:"type"`
	Type                     RemoveSymbolsType `yaml:"-" json:"-"`
	AddedSymbol              string            `yaml:"addedSymbol" json:"addedSymbol"`
	AddedSymbolCode          int               `yaml:"-" json:"-"`
	JumpToComponent          string            `yaml:"jumpToComponent" json:"jumpToComponent"`                   // jump to
	TargetComponents         []string          `yaml:"targetComponents" json:"targetComponents"`                 // 这些组件的中奖会需要参与remove
	SourcePositionCollection []string          `yaml:"sourcePositionCollection" json:"sourcePositionCollection"` // 源位置集合
	IgnoreSymbols            []string          `yaml:"ignoreSymbols" json:"ignoreSymbols"`                       // 忽略的symbol
	IgnoreSymbolCodes        []int             `yaml:"-" json:"-"`                                               // 忽略的symbol
	IsNeedProcSymbolVals     bool              `yaml:"isNeedProcSymbolVals" json:"isNeedProcSymbolVals"`         // 是否需要同时处理symbolVals
	EmptySymbolVal           int               `yaml:"emptySymbolVal" json:"emptySymbolVal"`                     // 空的symbolVal是什么
	OutputToComponent        string            `yaml:"outputToComponent" json:"outputToComponent"`               // outputToComponent
	Awards                   []*Award          `yaml:"awards" json:"awards"`                                     // 新的奖励系统
}

// SetLinkComponent sets a named downstream component link. Supported link
// values are "next" (default next component) and "jump" (jump target).
func (cfg *RemoveSymbolsConfig) SetLinkComponent(link string, componentName string) {
	switch link {
	case "next":
		cfg.DefaultNextComponent = componentName
	case "jump":
		cfg.JumpToComponent = componentName
	}
}

// RemoveSymbols implements the removeSymbols component behavior. It holds a
// reference to its configuration and a BasicComponent instance for common
// component behavior.
type RemoveSymbols struct {
	*BasicComponent `json:"-"`
	Config          *RemoveSymbolsConfig `json:"config"`
}

// Init reads a YAML file from disk and initializes the component configuration
// using InitEx. It is a convenience wrapper used when configurations are
// stored in files.
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

// InitEx initializes the component from a parsed configuration object. It
// performs symbol code resolution using the provided GamePropertyPool and
// prepares awards and other derived config fields.
func (removeSymbols *RemoveSymbols) InitEx(cfg any, pool *GamePropertyPool) error {
	rcfg, ok := cfg.(*RemoveSymbolsConfig)
	if !ok || rcfg == nil {
		goutils.Error("RemoveSymbols.InitEx:InvalidConfig",
			goutils.Err(ErrInvalidComponentConfig))

		return ErrInvalidComponentConfig
	}

	removeSymbols.Config = rcfg
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

// canRemove checks whether the symbol at position (x,y) in the provided
// GameScene is eligible for removal according to this component's
// configuration. The function is defensive and will return false if the
// coordinates are out-of-range.
func (removeSymbols *RemoveSymbols) canRemove(x, y int, gs *sgc7game.GameScene) bool {
	// Guard against out-of-range indices to avoid panic. Caller generally should
	// ensure x,y are valid, but be defensive here.
	if x < 0 || y < 0 || x >= gs.Width || y >= gs.Height {
		return false
	}

	curs := gs.Arr[x][y]

	// If curs < 0 we treat it as removable (matches previous behavior).
	if curs < 0 {
		return true
	}

	if len(removeSymbols.Config.IgnoreSymbolCodes) > 0 {
		return goutils.IndexOfIntSlice(removeSymbols.Config.IgnoreSymbolCodes, curs, 0) < 0
	}

	return true
}

// onBasic implements the default remove behavior. It removes positions
// collected from either configured target components or explicit source
// position collections. If an "other scene" (os) is provided, it will be
// updated accordingly (e.g. empty value mapping).
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

		if len(removeSymbols.Config.TargetComponents) > 0 {
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
		} else if len(removeSymbols.Config.SourcePositionCollection) > 0 {
			for _, cn := range removeSymbols.Config.SourcePositionCollection {

				curpos := gameProp.GetComponentPos(cn)
				for pi := 0; pi < len(curpos)/2; pi++ {
					x := curpos[pi*2]
					y := curpos[pi*2+1]
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

		if ngs == gs {
			return ErrComponentDoNothing
		}

		removeSymbols.AddOtherScene(gameProp, curpr, nos, &rscd.BasicComponentData)
	} else {
		if len(removeSymbols.Config.TargetComponents) > 0 {
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
		} else if len(removeSymbols.Config.SourcePositionCollection) > 0 {
			for _, cn := range removeSymbols.Config.SourcePositionCollection {
				curpos := gameProp.GetComponentPos(cn)

				for pi := 0; pi < len(curpos)/2; pi++ {
					x := curpos[pi*2]
					y := curpos[pi*2+1]
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

		if ngs == gs {
			return ErrComponentDoNothing
		}
	}

	if !gIsReleaseMode {
		if rscd.RemovedNum > 0 {
			rscd.AvgHeight = totalHeight * 100 / rscd.RemovedNum
		} else {
			rscd.AvgHeight = 0
		}
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &rscd.BasicComponentData)

	return nil
}

// onAdjacentPay implements an adjacency-aware removal used for "adjacentPay"
// style components. It retains the middle symbol of certain adjacent groups
// (e.g. for 3-in-a-row or 5-in-a-row) and replaces them with an added symbol
// after processing.
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
		if rscd.RemovedNum > 0 {
			rscd.AvgHeight = totalHeight * 100 / rscd.RemovedNum
		} else {
			rscd.AvgHeight = 0
		}
	}

	removeSymbols.AddScene(gameProp, curpr, ngs, &rscd.BasicComponentData)

	return nil
}

// OnPlayGame is the component entry point executed for a play step. It applies
// the configured removal behavior (basic or adjacentPay) to the target
// GameScene(s), processes controllers/awards and returns the next component
// to execute (or an ErrComponentDoNothing error when no change occurred).
func (removeSymbols *RemoveSymbols) OnPlayGame(gameProp *GameProperty, curpr *sgc7game.PlayResult, gp *GameParams, plugin sgc7plugin.IPlugin,
	cmd string, param string, ps sgc7game.IPlayerState, stake *sgc7game.Stake, prs []*sgc7game.PlayResult, cd IComponentData) (string, error) {
	rscd, ok := cd.(*RemoveSymbolsData)
	if !ok || rscd == nil {
		goutils.Error("RemoveSymbols.OnPlayGame:InvalidComponentData",
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}
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
			goutils.Err(ErrInvalidComponentConfig))

		return "", ErrInvalidComponentConfig
	}

	removeSymbols.ProcControllers(gameProp, plugin, curpr, gp, -1, "")

	nc := removeSymbols.onStepEnd(gameProp, curpr, gp, removeSymbols.Config.JumpToComponent)

	return nc, nil
}

// ProcControllers executes configured award controllers after removal is
// performed. It's called by OnPlayGame after the scene changes are applied.
func (removeSymbols *RemoveSymbols) ProcControllers(gameProp *GameProperty, plugin sgc7plugin.IPlugin, curpr *sgc7game.PlayResult, gp *GameParams, val int, strVal string) {
	if len(removeSymbols.Config.Awards) > 0 {
		gameProp.procAwards(plugin, removeSymbols.Config.Awards, curpr, gp)
	}
}

// OnAsciiGame outputs a textual representation of the primary scene affected
// by this component to aid debugging and development. It expects a
// *RemoveSymbolsData as component data and will safely no-op if UsedScenes is
// empty.
func (removeSymbols *RemoveSymbols) OnAsciiGame(gameProp *GameProperty, pr *sgc7game.PlayResult, lst []*sgc7game.PlayResult, mapSymbolColor *asciigame.SymbolColorMap, cd IComponentData) error {
	bcd := cd.(*RemoveSymbolsData)

	if len(bcd.UsedScenes) > 0 {
		asciigame.OutputScene("after removedSymbols", pr.Scenes[bcd.UsedScenes[0]], mapSymbolColor)
	}

	return nil
}

// NewComponentData returns a new, zeroed RemoveSymbolsData instance. The
// runtime will call this to allocate per-play component data.
func (removeSymbols *RemoveSymbols) NewComponentData() IComponentData {
	return &RemoveSymbolsData{}
}

// EachUsedResults is required by the component interface but not used by
// removeSymbols; it is intentionally a no-op here.
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

// NewRemoveSymbols constructs a RemoveSymbols component with the given name.
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
// "outputToComp": "bg-pos-rmoved",
// "sourcePositionCollection": [
//
//	"bg-scatterzone"
//
// ]
// jsonRemoveSymbols is the lightweight JSON/YAML mapping used when parsing an
// inline removeSymbols configuration from AST nodes. It is converted to the
// full RemoveSymbolsConfig via build().
type jsonRemoveSymbols struct {
	Type                     string   `json:"type"`                     // type
	AddedSymbol              string   `json:"addedSymbol"`              // addedSymbol
	TargetComponents         []string `json:"targetComponents"`         // target components whose results cause removal
	IgnoreSymbols            []string `json:"ignoreSymbols"`            // symbols to ignore when removing
	IsNeedProcSymbolVals     bool     `json:"isNeedProcSymbolVals"`     // whether to process other scene symbol values
	EmptySymbolVal           int      `json:"emptySymbolVal"`           // value to set in other scene for emptied slots
	OutputToComponent        string   `json:"outputToComp"`             // optional component to output removed positions to
	SourcePositionCollection []string `json:"sourcePositionCollection"` // alternative explicit list of positions to remove
}

// build converts the jsonRemoveSymbols helper into a full
// RemoveSymbolsConfig instance; slices are cloned to avoid sharing memory
// between parsed AST structures and the resulting config.
func (jcfg *jsonRemoveSymbols) build() *RemoveSymbolsConfig {
	cfg := &RemoveSymbolsConfig{
		StrType:                  jcfg.Type,
		TargetComponents:         slices.Clone(jcfg.TargetComponents),
		IgnoreSymbols:            slices.Clone(jcfg.IgnoreSymbols),
		IsNeedProcSymbolVals:     jcfg.IsNeedProcSymbolVals,
		EmptySymbolVal:           jcfg.EmptySymbolVal,
		AddedSymbol:              jcfg.AddedSymbol,
		OutputToComponent:        jcfg.OutputToComponent,
		SourcePositionCollection: slices.Clone(jcfg.SourcePositionCollection),
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
