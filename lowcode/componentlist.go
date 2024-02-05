package lowcode

import (
	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/stats2"
	"go.uber.org/zap"
)

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

	for _, v := range lst.MapComponents {
		v.OnGameInited(lst)
	}
}

// GetAllLinkComponents - get all link components
func (lst *ComponentList) GetAllLinkComponents(componentName string) []string {
	ic, isok := lst.MapComponents[componentName]
	if !isok {
		goutils.Error("ComponentList.GetAllLinkComponents",
			zap.String("name", componentName),
			zap.Error(ErrInvalidComponentName))

		return nil
	}

	allcomponents := []string{componentName}

	curlst := ic.GetAllLinkComponents()
	for _, v := range curlst {
		if v == "" {
			continue
		}

		allcomponents = InsStringSliceNonRep(allcomponents, v)
		// allcomponents = append(allcomponents, v)
		children := lst.GetAllLinkComponents(v)
		allcomponents = InsSliceNonRep(allcomponents, children)
		// allcomponents = append(allcomponents, children...)
	}

	return allcomponents
}

func NewComponentList() *ComponentList {
	lst := &ComponentList{
		MapComponents: make(map[string]IComponent),
	}

	return lst
}
