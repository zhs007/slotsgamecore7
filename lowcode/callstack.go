package lowcode

import (
	"fmt"
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// 关于 onNewGame 和 onNewStep
// 1. callstack不会跨game的保留数据，所以每次新游戏，callstack都会是一个空的
// 2. callstack只会保留主调用堆栈的数据，且当componentData第一次被获取时，执行OnNewStep

// 关于 callstack 层级
// 1. 至少有一个 global，所有component都会有一个global对象
// 2. 可能会有 foreach 层，一个循环体，同时只会有一层（i = 0; i < 10;）这种，同时也只会有一层，这个 foreach 层，只会有当前子组件对象（只可能是连线对象）
//  这些对象在一开始构造，其它的component对象在这一层是找不到的
// 3. respin 不能出现在 foreach 层

// FuncOnEachHistoryComponent -
type FuncOnEachHistoryComponent func(tag string, gameProp *GameProperty, ic IComponent, cd IComponentData) error

type callStackNode struct {
	CoreComponent        IComponent
	Name                 string
	MapComponentData     map[string]IComponentData
	mapHistory           map[string]IComponentData
	SymbolCode           int
	CurIndex             int
	isNoAutoNew          bool
	cacheSceneIndex      int
	cacheOtherSceneIndex int
}

func (csn *callStackNode) IsSame(component IComponent, symbolCode int, i int) bool {
	return csn.CoreComponent == component && csn.SymbolCode == symbolCode && csn.CurIndex == i
}

func (csn *callStackNode) IsInCallStack(componentName string) bool {
	_, isok := csn.mapHistory[componentName]

	return isok
}

func (csn *callStackNode) addComponentData(ic IComponent) {
	cd := ic.NewComponentData()
	csn.MapComponentData[ic.GetName()] = cd
}

func (csn *callStackNode) onCallEnd(ic IComponent, cd IComponentData) {
	csn.mapHistory[ic.GetName()] = cd
}

func (csn *callStackNode) genTag(ic IComponent) string {
	if csn.Name == "" {
		return ic.GetName()
	}

	return fmt.Sprintf("%v/%v", csn.Name, ic.GetName())
}

func (csn *callStackNode) getComponentData(gameProp *GameProperty, ic IComponent, _ *CallStack) IComponentData {
	name := ic.GetName()
	cd, isok := csn.MapComponentData[name]
	if !isok {
		// 以前没有考虑 step 的情况，所以在 global 时需要 new 一个
		if csn.isNoAutoNew {
			return nil
		}

		// ncd := cs.getComponentDataFromPreSteps(gameProp, ic)
		// if ncd != nil {
		// 	csn.MapComponentData[name] = ncd

		// 	csn.mapHistory[name] = ncd

		// 	return ncd
		// }

		cd = ic.NewComponentData()

		cd.OnNewGame(gameProp, ic)

		csn.MapComponentData[name] = cd

		csn.mapHistory[name] = cd
	}

	return cd
}

func (csn *callStackNode) OnNewStep() {
	csn.mapHistory = make(map[string]IComponentData)
}

func newGlobalCallStackNode() *callStackNode {
	return &callStackNode{
		mapHistory:       make(map[string]IComponentData),
		MapComponentData: make(map[string]IComponentData),
	}
}

func newEachSymbolCallStackNode(component IComponent, symbolCode int, i int, cacheSceneIndex int, cacheOtherSceneIndex int, pt *sgc7game.PayTables) *callStackNode {
	return &callStackNode{
		CoreComponent:        component,
		Name:                 fmt.Sprintf("%v:%v>%v", component.GetName(), i, pt.GetStringFromInt(symbolCode)),
		mapHistory:           make(map[string]IComponentData),
		MapComponentData:     make(map[string]IComponentData),
		SymbolCode:           symbolCode,
		CurIndex:             i,
		isNoAutoNew:          true,
		cacheSceneIndex:      cacheSceneIndex,
		cacheOtherSceneIndex: cacheOtherSceneIndex,
	}
}

type callStackHistoryNode struct {
	tag       string
	component IComponent
	cd        IComponentData
}

type CallStack struct {
	nodes        []*callStackNode
	historyNodes []*callStackHistoryNode
	// stepsCallStack []*CallStack
	// isNew          bool
}

// GetGlobalComponentData -
func (cs *CallStack) GetGlobalComponentData(gameProp *GameProperty, ic IComponent) IComponentData {
	return cs.nodes[0].getComponentData(gameProp, ic, cs)
}

func (cs *CallStack) GetCurComponentData(gameProp *GameProperty, ic IComponent) IComponentData {
	return cs.nodes[len(cs.nodes)-1].getComponentData(gameProp, ic, cs)
}

// func (cs *CallStack) getComponentDataFromPreSteps(gameProp *GameProperty, ic IComponent) IComponentData {
// 	if len(cs.stepsCallStack) == 0 {
// 		return nil
// 	}

// 	return cs.stepsCallStack[len(cs.stepsCallStack)-1].GetComponentData(gameProp, ic)

// }

func (cs *CallStack) GetComponentData(gameProp *GameProperty, ic IComponent) IComponentData {
	if len(cs.nodes) == 1 {
		return cs.nodes[0].getComponentData(gameProp, ic, cs)
	}

	cd := cs.nodes[len(cs.nodes)-1].getComponentData(gameProp, ic, cs)
	if cd == nil {
		return cs.nodes[0].getComponentData(gameProp, ic, cs)
	}

	return cd
}

func (cs *CallStack) GetComponentNum() int {
	return len(cs.historyNodes)
}

func (cs *CallStack) OnNewGame() {
	cs.nodes = cs.nodes[:0]

	cs.nodes = append(cs.nodes, newGlobalCallStackNode())

	cs.historyNodes = nil

	// cs.isNew = true
}

func (cs *CallStack) OnNewStep() *CallStack {
	// if
	// newCS := NewCallStack()
	// newCS.OnNewGame()

	// if len(cs.stepsCallStack) > 1 {
	// 	newCS.stepsCallStack = make([]*CallStack, len(cs.stepsCallStack), len(cs.stepsCallStack)+1)
	// 	copy(newCS.stepsCallStack, cs.stepsCallStack)
	// 	newCS.stepsCallStack = append(newCS.stepsCallStack, cs)
	// } else {
	// 	newCS.stepsCallStack = []*CallStack{cs}
	// }

	cs.nodes = cs.nodes[0:1]

	cs.nodes[0].OnNewStep()

	cs.historyNodes = nil

	return nil
}

func (cs *CallStack) Each(gameProp *GameProperty, onEach FuncOnEachHistoryComponent) error {
	for i, node := range cs.historyNodes {
		err := onEach(node.tag, gameProp, node.component, node.cd)
		if err != nil {
			goutils.Error("CallStack.Each",
				slog.Int("i", i),
				slog.String("tag", node.tag),
				goutils.Err(err))

			return err
		}
	}

	return nil
}

func (cs *CallStack) genTag(component IComponent) string {
	if len(cs.nodes) == 1 {
		return component.GetName()
	}

	tag := ""

	for i, node := range cs.nodes {
		if i != 0 {
			tag += "/"
		}

		tag += node.Name
	}

	tag += "/"
	tag += component.GetName()

	return tag
}

func (cs *CallStack) ComponentDone(gameProp *GameProperty, component IComponent, cd IComponentData) {
	cs.historyNodes = append(cs.historyNodes, &callStackHistoryNode{
		tag:       cs.genTag(component),
		component: component,
		cd:        cd,
	})
}

func (cs *CallStack) IsInCurCallStack(componentName string) bool {
	return cs.nodes[len(cs.nodes)-1].IsInCallStack(componentName)
}

func (cs *CallStack) GetCurCallStackSymbol() int {
	if len(cs.nodes) == 1 {
		return -1
	}

	return cs.nodes[len(cs.nodes)-1].SymbolCode
}

func (cs *CallStack) OnCallEnd(ic IComponent, cd IComponentData) string {
	cs.nodes[len(cs.nodes)-1].onCallEnd(ic, cd)

	tag := cs.nodes[len(cs.nodes)-1].genTag(ic)
	cs.historyNodes = append(cs.historyNodes, &callStackHistoryNode{
		tag:       tag,
		component: ic,
		cd:        cd,
	})

	return tag
}

func (cs *CallStack) StartEachSymbols(gameProp *GameProperty, component IComponent, children []string, symbolCode int, i int) error {
	node := newEachSymbolCallStackNode(component, symbolCode, i, len(gameProp.SceneStack.Scenes), len(gameProp.OtherSceneStack.Scenes), gameProp.Pool.DefaultPaytables)

	components := gameProp.Components

	for _, v := range children {
		ic, isok := components.MapComponents[v]
		if !isok {
			goutils.Error("CallStack.StartEachSymbols:children",
				goutils.Err(ErrInvalidComponentName))

			return ErrInvalidComponentName
		}

		node.addComponentData(ic)
	}

	cs.nodes = append(cs.nodes, node)

	return nil
}

func (cs *CallStack) onEachSymbolsEnd(gameProp *GameProperty, component IComponent, symbolCode int, i int) error {
	if len(cs.nodes) > 1 && cs.nodes[len(cs.nodes)-1].IsSame(component, symbolCode, i) {
		gameProp.SceneStack.TruncateTo(cs.nodes[len(cs.nodes)-1].cacheSceneIndex)
		gameProp.OtherSceneStack.TruncateTo(cs.nodes[len(cs.nodes)-1].cacheOtherSceneIndex)

		cs.nodes = cs.nodes[:len(cs.nodes)-1]

		return nil
	}

	goutils.Error("CallStack.onEachSymbolsEnd",
		slog.String("component", component.GetName()),
		slog.Int("symbolCode", symbolCode),
		slog.Int("i", i),
		goutils.Err(ErrInvalidCallStackNode))

	return ErrInvalidCallStackNode
}

func NewCallStack() *CallStack {
	cs := &CallStack{}

	return cs
}
