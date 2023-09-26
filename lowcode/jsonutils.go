package lowcode

import (
	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

func isIntValWeights(lst []*weightData) bool {
	for _, v := range lst {
		_, err := goutils.String2Int64(v.Val)
		if err != nil {
			return false
		}
	}

	return true
}

func parseValWeights(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseValWeights:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseValWeights:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	if isIntValWeights(lst) {
		vals := []sgc7game.IVal{}
		weights := []int{}

		for _, v := range lst {
			i64, err := goutils.String2Int64(v.Val)
			if err != nil {
				goutils.Error("parseValWeights:String2Int64",
					zap.Error(err))

				return nil, err
			}

			vals = append(vals, sgc7game.NewIntValEx[int](int(i64)))
			weights = append(weights, v.Weight)
		}

		return sgc7game.NewValWeights2(vals, weights)
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		vals = append(vals, sgc7game.NewStrValEx(v.Val))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}
