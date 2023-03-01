package lowcode

type ComponentList struct {
	Components    []IComponent
	MapComponents map[string]IComponent
}

func (lst *ComponentList) AddComponent(name string, component IComponent) {
	lst.Components = append(lst.Components, component)

	lst.MapComponents[name] = component
}

func NewComponentList() *ComponentList {
	lst := &ComponentList{
		MapComponents: make(map[string]IComponent),
	}

	return lst
}
