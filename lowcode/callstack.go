package lowcode

// 关于 onNewGame 和 onNewStep
// 1. callstack不会跨game的保留数据，所以每次新游戏，callstack都会是一个空的
// 2. callstack只会保留主调用堆栈的数据，且当componentData第一次被获取时，执行OnNewStep

// FuncOnEachHistoryComponent - if return false then break
type FuncOnEachHistoryComponent func(tag string, gameProp *GameProperty, ic IComponent, cd IComponentData) bool

type callStackNode struct {
	Name             string
	MapComponentData map[string]IComponentData
	mapHistory       map[string]IComponentData
}

func (csn *callStackNode) GetComponentData(gameProp *GameProperty, ic IComponent) IComponentData {
	name := ic.GetName()
	cd, isok := csn.MapComponentData[name]
	if !isok {
		cd = ic.NewComponentData()

		csn.MapComponentData[name] = cd

		csn.mapHistory[name] = cd
	} else {
		_, isok := csn.mapHistory[name]
		if !isok {
			cd.OnNewStep(gameProp, ic)

			csn.mapHistory[name] = cd
		}
	}

	return cd
}

func (csn *callStackNode) OnNewStep() {
	csn.mapHistory = make(map[string]IComponentData)
}

func newCallStackNode(name string) *callStackNode {
	return &callStackNode{
		Name:             name,
		mapHistory:       make(map[string]IComponentData),
		MapComponentData: make(map[string]IComponentData),
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
}

func (cs *CallStack) GetCurComponentData(gameProp *GameProperty, ic IComponent) IComponentData {
	return cs.nodes[len(cs.nodes)-1].GetComponentData(gameProp, ic)
}

func (cs *CallStack) GetComponentNum() int {
	return len(cs.historyNodes)
}

func (cs *CallStack) OnNewGame() {
	cs.nodes = cs.nodes[:0]

	cs.nodes = append(cs.nodes, newCallStackNode(""))
}

func (cs *CallStack) OnNewStep() {
	cs.nodes = cs.nodes[0:1]

	cs.nodes[0].OnNewStep()
}

func (cs *CallStack) Each(gameProp *GameProperty, onEach FuncOnEachHistoryComponent) {
	for _, node := range cs.historyNodes {
		onEach(node.tag, gameProp, node.component, node.cd)
	}
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

func NewCallStack(name string) *CallStack {
	cs := &CallStack{}

	return cs
}
