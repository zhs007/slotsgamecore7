package lowcode

import (
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

func getConfigInCell(cell *ast.Node) (*ast.Node, string, error) {
	componentValues := cell.Get("componentValues")
	if componentValues == nil {
		goutils.Error("getConfigInCell:componentValues",
			zap.Error(ErrNoComponentValues))

		return nil, "", ErrNoComponentValues
	}

	cfg := componentValues.Get("configuration")
	if cfg == nil {
		goutils.Error("getConfigInCell:configuration",
			zap.Error(ErrInvalidJsonNode))

		return nil, "", ErrInvalidJsonNode
	}

	label := componentValues.Get("label")
	if cfg == nil {
		goutils.Error("getConfigInCell:label",
			zap.Error(ErrInvalidJsonNode))

		return nil, "", ErrInvalidJsonNode
	}

	str, err := label.String()
	if cfg == nil {
		goutils.Error("getConfigInCell:label.String",
			zap.Error(err))

		return nil, "", err
	}

	return cfg, str, nil
}
