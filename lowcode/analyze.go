package lowcode

import (
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

type SPCSPNode struct {
	Component string
	Children  *SPCNode
}

type SPCNode struct {
	Parent           *SPCNode
	NormalComponents []string
	Children         []*SPCSPNode
}

func isParentComponentInSPC(component IComponent) bool {
	if component.IsRespin() {
		return true
	}

	return false
}

func ParseStepParentChildren(lst *ComponentList, start string) (*SPCNode, error) {
	if lst == nil || len(lst.MapComponents) <= 0 {
		return nil, nil
	}

	node := &SPCNode{}

	cn := start
	for {
		ic, isok := lst.MapComponents[cn]
		if !isok {
			goutils.Error("ParseStepParentChildren:MapComponents",
				zap.String("name", cn),
				zap.Error(ErrInvalidComponentName))

			return nil, ErrInvalidComponentName
		}

		if isParentComponentInSPC(ic) {
			child, err := ParseStepParentChildren(lst, cn)
			if err != nil {
				goutils.Error("ParseStepParentChildren:ParseStepParentChildren",
					zap.String("name", cn),
					zap.Error(err))

				return nil, err
			}

			cp := &SPCSPNode{
				Component: cn,
				Children:  child,
			}

			node.Children = append(node.Children, cp)
		} else {
			node.NormalComponents = append(node.NormalComponents, cn)
		}
	}

	if len(node.NormalComponents) == 0 {
		return nil, nil
	}

	return node, nil
}
