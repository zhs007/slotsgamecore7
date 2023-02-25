package lowcode

type ComponentMgr struct {
	MapComponent map[string]FuncNewComponent
}

func NewComponentMgr() *ComponentMgr {
	return &ComponentMgr{
		MapComponent: make(map[string]FuncNewComponent),
	}
}
