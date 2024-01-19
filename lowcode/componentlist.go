package lowcode

import "github.com/zhs007/slotsgamecore7/stats2"

type ComponentList struct {
	Components    []IComponent
	MapComponents map[string]IComponent
	Stats2        *stats2.Stats
}

func (lst *ComponentList) AddComponent(name string, component IComponent) {
	lst.Components = append(lst.Components, component)

	lst.MapComponents[name] = component
}

func (lst *ComponentList) onInit() {
	if gAllowStats2 {
		lst.Stats2 = NewStats2(lst)

		lst.Stats2.Start()
	}
}

func NewComponentList() *ComponentList {
	lst := &ComponentList{
		MapComponents: make(map[string]IComponent),
	}

	return lst
}
