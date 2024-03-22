package lowcode

import (
	"log/slog"

	"github.com/zhs007/goutils"
)

type SPCNode struct {
	Parent           *SPCNode
	NormalComponents []string
	Children         []*SPCNode
	Root             string
}

func (node *SPCNode) GetComponents() []string {
	lst := []string{}

	if node.Root != "" {
		lst = append(lst, node.Root)
	}

	lst = append(lst, node.NormalComponents...)

	for _, v := range node.Children {
		curlst := v.GetComponents()

		lst = append(lst, curlst...)
	}

	return lst
}

func (node *SPCNode) GetParent(component string) string {
	if node.Root == component {
		return node.Parent.Root
	}

	if goutils.IndexOfStringSlice(node.NormalComponents, component, 0) >= 0 {
		return node.Root
	}

	for _, v := range node.Children {
		p := v.GetParent(component)
		if p != "" {
			return p
		}
	}

	return ""
}

func (node *SPCNode) CountComponentNum() int {
	num := 0

	num += len(node.NormalComponents)
	num += len(node.Children)

	for _, c := range node.Children {
		num += c.CountComponentNum()
	}

	return num
}

func (node *SPCNode) CountParentNum() int {
	num := 0

	num += len(node.Children)

	for _, c := range node.Children {
		num += c.CountParentNum()
	}

	return num
}

func (node *SPCNode) CountDeep() int {
	if len(node.Children) > 0 {
		num := 1

		for _, c := range node.Children {
			num += c.CountDeep()
		}

		return num
	}

	return 0
}

func (node *SPCNode) AddNormal(componentName string) {
	if !node.IsInNormal(componentName) {
		node.NormalComponents = append(node.NormalComponents, componentName)
	}
}

func (node *SPCNode) AddChild(child *SPCNode) {
	if !node.IsChildren(child.Root) {
		child.Parent = node

		node.Children = append(node.Children, child)
	}
}

func (node *SPCNode) Format() {
	lst := []string{}
	for _, v := range node.NormalComponents {
		if !node.IsInChildren(v) {
			lst = append(lst, v)
		}
	}

	node.NormalComponents = lst

	for _, c := range node.Children {
		c.Format()
	}
}

func (node *SPCNode) IsInNormal(componentName string) bool {
	return goutils.IndexOfStringSlice(node.NormalComponents, componentName, 0) >= 0
}

func (node *SPCNode) IsChildren(componentName string) bool {
	for _, v := range node.Children {
		if v.Root == componentName {
			return true
		}
	}

	return false
}

func (node *SPCNode) IsInChildren(componentName string) bool {
	for _, v := range node.Children {
		if v.Root == componentName {
			return true
		}

		for _, n := range v.NormalComponents {
			for n == componentName {
				return true
			}
		}

		if v.IsInChildren(componentName) {
			return true
		}
	}

	return false
}

func isParentComponentInSPC(component IComponent) bool {
	if component.IsRespin() || component.IsForeach() {
		return true
	}

	return false
}

func parseNextComponents(lst *ComponentList, start string) (*SPCNode, error) {
	node := &SPCNode{}
	cn := start

	ic, isok := lst.MapComponents[cn]
	if !isok {
		goutils.Error("parseNextComponents:MapComponents",
			slog.String("name", cn),
			goutils.Err(ErrInvalidComponentName))

		return nil, ErrInvalidComponentName
	}

	if isParentComponentInSPC(ic) {
		children := ic.GetChildLinkComponents()
		if len(children) == 1 {
			childNode, err := parseNextComponents(lst, children[0])
			if err != nil {
				goutils.Error("parseNextComponents:parseNextComponents",
					slog.String("name", children[0]),
					goutils.Err(err))

				return nil, err
			}

			childNode.Root = cn

			node.AddChild(childNode)
		} else if len(children) > 1 {
			goutils.Error("parseNextComponents",
				slog.String("name", cn),
				slog.Any("children", children),
				goutils.Err(ErrInvalidComponentChildren))

			return nil, ErrInvalidComponentChildren
		}
	} else {
		node.AddNormal(cn)
	}

	nextComponents := ic.GetNextLinkComponents()
	for _, curcomponent := range nextComponents {
		if curcomponent != "" {
			nextNode, err := parseNextComponents(lst, curcomponent)
			if err != nil {
				goutils.Error("parseNextComponents:nextComponents:parseNextComponents",
					slog.String("name", curcomponent),
					goutils.Err(err))

				return nil, err
			}

			for _, child := range nextNode.Children {
				node.AddChild(child)
			}

			for _, nc := range nextNode.NormalComponents {
				node.AddNormal(nc)
			}
		}
	}

	return node, nil
}

func ParseStepParentChildren(lst *ComponentList, start string) (*SPCNode, error) {
	if lst == nil || len(lst.MapComponents) <= 0 {
		return nil, nil
	}

	node, err := parseNextComponents(lst, start)
	if err != nil {
		goutils.Error("ParseStepParentChildren:parseNextComponents",
			slog.String("name", start),
			goutils.Err(err))

		return nil, err
	}

	node.Format()

	return node, nil
}
