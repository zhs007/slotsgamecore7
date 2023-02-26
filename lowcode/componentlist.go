package lowcode

type ComponentList struct {
	Components []IComponent
}

func (lst *ComponentList) AddComponent(component IComponent) {
	lst.Components = append(lst.Components, component)
}

func NewComponentList() *ComponentList {
	lst := &ComponentList{}

	return lst
}
