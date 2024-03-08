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

func parseReelSetWeights(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseReelSetWeights:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseReelSetWeights:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		vals = append(vals, sgc7game.NewStrValEx(v.Val))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseSymbolWeights(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolWeights:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseSymbolWeights:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for i, v := range lst {
		sc, isok := paytables.MapSymbols[v.Val]
		if !isok {
			goutils.Error("parseSymbolWeights:MapSymbols",
				zap.Int("i", i),
				zap.String("symbol", v.Val),
				zap.Error(ErrIvalidSymbol))

			return nil, ErrIvalidSymbol
		}

		vals = append(vals, sgc7game.NewIntValEx(sc))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseStrWeights(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseStrWeights:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseStrWeights:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		vals = append(vals, sgc7game.NewStrValEx(v.Val))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseIntValMappingFile(n *ast.Node) (*sgc7game.ValMapping2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValMappingFile:MarshalJSON",
			zap.Error(err))

		return nil, err
	}

	lst := []*mappingData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValMappingFile:Unmarshal",
			zap.Error(err))

		return nil, err
	}

	invals := []int{}
	outvals := []sgc7game.IVal{}

	for _, v := range lst {
		inv, err := goutils.String2Int64(v.In)
		if err != nil {
			goutils.Error("parseIntValMappingFile:String2Int64",
				zap.String("in", v.In),
				zap.Error(err))

			return nil, err
		}

		outv, err := goutils.String2Int64(v.Out)
		if err != nil {
			goutils.Error("parseIntValMappingFile:String2Int64",
				zap.String("out", v.Out),
				zap.Error(err))

			return nil, err
		}

		outvals = append(outvals, sgc7game.NewIntValEx(outv))
		invals = append(invals, int(inv))
	}

	return sgc7game.NewValMapping2(invals, outvals)
}
