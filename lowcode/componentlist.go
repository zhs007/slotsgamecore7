package lowcode

import (
	"log/slog"

	"github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/stats2"
)

type ComponentList struct {
	Components    []IComponent          `yaml:"-" json:"-"`
	MapComponents map[string]IComponent `yaml:"mapComponents" json:"mapComponents"`
	Stats2        *stats2.Stats         `yaml:"-" json:"-"`
	statsNodeData *SPCNode              `yaml:"-" json:"-"`
	RngLib        *RngLib               `yaml:"-" json:"-"`
}

func (lst *ComponentList) AddComponent(name string, component IComponent) {
	lst.Components = append(lst.Components, component)

	lst.MapComponents[name] = component
}

func (lst *ComponentList) onInit(start string) error {
	if gAllowStats2 {
		node, err := ParseStepParentChildren(lst, start)
		if err != nil {
			goutils.Error("ComponentList.onInit:ParseStepParentChildren",
				slog.String("start", start),
				goutils.Err(err))

			return err
		}

		lst.statsNodeData = node

		lst.Stats2 = NewStats2(lst)

		if gRngLibConfig != "" {
			rnglib, err := LoadRngLib(gRngLibConfig)
			if err != nil {
				goutils.Error("ComponentList.onInit:LoadRngLib",
					slog.String("gRngLibConfig", gRngLibConfig),
					goutils.Err(err))

				// return err
			} else {
				lst.RngLib = rnglib
			}
		}

		lst.Stats2.Start()
	}

	for _, v := range lst.MapComponents {
		v.OnGameInited(lst)
	}

	return nil
}

// GetAllLinkComponents - get all link components
func (lst *ComponentList) GetAllLinkComponents(componentName string) []string {
	ic, isok := lst.MapComponents[componentName]
	if !isok {
		goutils.Error("ComponentList.GetAllLinkComponents",
			slog.String("name", componentName),
			goutils.Err(ErrInvalidComponentName))

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
