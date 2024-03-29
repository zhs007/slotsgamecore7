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

func (node *SPCNode) IsIn(componentName string) bool {
	if node.IsInNormal(componentName) {
		return true
	}

	if node.IsInChildren(componentName) {
		return true
	}

	return false
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

func parseNextComponents(lst *ComponentList, start string, historys []string) (*SPCNode, []string, error) {
	if start == "" {
		return nil, historys, nil
	}

	node := &SPCNode{}
	cn := start

	ic, isok := lst.MapComponents[cn]
	if !isok {
		goutils.Error("parseNextComponents:MapComponents",
			slog.String("name", cn),
			goutils.Err(ErrInvalidComponentName))

		return nil, nil, ErrInvalidComponentName
	}

	if isParentComponentInSPC(ic) {
		children := ic.GetChildLinkComponents()
		if len(children) == 1 {
			childNode, nh, err := parseNextComponents(lst, children[0], historys)
			if err != nil {
				goutils.Error("parseNextComponents:parseNextComponents",
					slog.String("name", children[0]),
					goutils.Err(err))

				return nil, nil, err
			}

			historys = nh

			if childNode != nil {
				childNode.Root = cn

				historys = append(historys, children[0])
				node.AddChild(childNode)
			}
		} else if len(children) > 1 {
			goutils.Error("parseNextComponents",
				slog.String("name", cn),
				slog.Any("children", children),
				goutils.Err(ErrInvalidComponentChildren))

			return nil, nil, ErrInvalidComponentChildren
		}
	} else {
		historys = append(historys, cn)

		node.AddNormal(cn)
	}

	nextComponents := ic.GetNextLinkComponents()
	for _, curcomponent := range nextComponents {
		if curcomponent != "" && goutils.IndexOfStringSlice(historys, curcomponent, 0) < 0 {
			nextNode, nh, err := parseNextComponents(lst, curcomponent, historys)
			if err != nil {
				goutils.Error("parseNextComponents:nextComponents:parseNextComponents",
					slog.String("name", curcomponent),
					goutils.Err(err))

				return nil, nil, err
			}

			historys = nh

			for _, child := range nextNode.Children {
				historys = append(historys, child.Root)
				node.AddChild(child)
			}

			for _, nc := range nextNode.NormalComponents {
				historys = append(historys, nc)
				node.AddNormal(nc)
			}
		}
	}

	return node, historys, nil
}

func ParseStepParentChildren(lst *ComponentList, start string) (*SPCNode, error) {
	if lst == nil || len(lst.MapComponents) <= 0 {
		return nil, nil
	}

	node, _, err := parseNextComponents(lst, start, []string{})
	if err != nil {
		goutils.Error("ParseStepParentChildren:parseNextComponents",
			slog.String("name", start),
			goutils.Err(err))

		return nil, err
	}

	node.Format()

	return node, nil
}
