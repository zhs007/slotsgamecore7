package lowcode

type ComponentList struct {
	Components    []IComponent
	MapComponents map[string]IComponent
	Stats2        *Stats2
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
