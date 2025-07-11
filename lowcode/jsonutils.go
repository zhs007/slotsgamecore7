package lowcode

import (
	"log/slog"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/ast"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
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
			goutils.Err(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseValWeights:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	if isIntValWeights(lst) {
		vals := []sgc7game.IVal{}
		weights := []int{}

		for _, v := range lst {
			i64, err := goutils.String2Int64(v.Val)
			if err != nil {
				goutils.Error("parseValWeights:String2Int64",
					goutils.Err(err))

				return nil, err
			}

			vals = append(vals, sgc7game.NewIntValEx(int(i64)))
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
			goutils.Err(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseReelSetWeights:Unmarshal",
			goutils.Err(err))

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

func parseReelSetWeights2(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseReelSetWeights2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseReelSetWeights2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		vals = append(vals, sgc7game.NewStrValEx(strings.TrimSpace(v["val"])))
		i64, err := goutils.String2Int64(strings.TrimSpace(v["weight"]))
		if err != nil {
			goutils.Error("parseReelSetWeights2:String2Int64",
				slog.String("weight", v["weight"]),
				goutils.Err(err))

			return nil, err
		}

		weights = append(weights, int(i64))
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseSymbolWeights(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolWeights:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseSymbolWeights:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for i, v := range lst {
		sc, isok := paytables.MapSymbols[v.Val]
		if !isok {
			goutils.Error("parseSymbolWeights:MapSymbols",
				slog.Int("i", i),
				slog.String("symbol", v.Val),
				goutils.Err(ErrInvalidSymbol))

			return nil, ErrInvalidSymbol
		}

		vals = append(vals, sgc7game.NewIntValEx(sc))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseSymbolWeights2(n *ast.Node, paytables *sgc7game.PayTables) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseSymbolWeights2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseSymbolWeights2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for i, v := range lst {
		sc, isok := paytables.MapSymbols[strings.TrimSpace(v["val"])]
		if !isok {
			goutils.Error("parseSymbolWeights2:MapSymbols",
				slog.Int("i", i),
				slog.String("symbol", v["val"]),
				goutils.Err(ErrInvalidSymbol))

			return nil, ErrInvalidSymbol
		}

		vals = append(vals, sgc7game.NewIntValEx(sc))

		i64, err := goutils.String2Int64(strings.TrimSpace(v["weight"]))
		if err != nil {
			goutils.Error("parseSymbolWeights2:String2Int64",
				slog.String("weight", v["weight"]),
				goutils.Err(err))

			return nil, err
		}

		weights = append(weights, int(i64))
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseStrWeights2(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseStrWeights2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseStrWeights2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		vals = append(vals, sgc7game.NewStrValEx(v["val"]))

		i64, err := goutils.String2Int64(strings.TrimSpace(v["weight"]))
		if err != nil {
			goutils.Error("parseStrWeights2:String2Int64",
				slog.String("weight", v["weight"]),
				goutils.Err(err))

			return nil, err
		}

		weights = append(weights, int(i64))
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseStrWeights(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseStrWeights:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseStrWeights:Unmarshal",
			goutils.Err(err))

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

func parseIntValWeights2(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValWeights2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValWeights2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		i64, err := goutils.String2Int64(strings.TrimSpace(v["val"]))
		if err != nil {
			goutils.Error("parseIntValWeights2:String2Int64",
				slog.String("val", v["val"]),
				goutils.Err(err))

			return nil, err
		}

		vals = append(vals, sgc7game.NewIntValEx(i64))

		i64, err = goutils.String2Int64(strings.TrimSpace(v["weight"]))
		if err != nil {
			goutils.Error("parseIntValWeights2:String2Int64",
				slog.String("weight", v["weight"]),
				goutils.Err(err))

			return nil, err
		}

		weights = append(weights, int(i64))
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseIntValWeights(n *ast.Node) (*sgc7game.ValWeights2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValWeights:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*weightData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValWeights:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	vals := []sgc7game.IVal{}
	weights := []int{}

	for _, v := range lst {
		iv, err := goutils.String2Int64(v.Val)
		if err != nil {
			goutils.Error("parseIntValWeights:String2Int64",
				slog.String("val", v.Val),
				goutils.Err(err))

			return nil, err
		}

		vals = append(vals, sgc7game.NewIntValEx(int(iv)))
		weights = append(weights, v.Weight)
	}

	return sgc7game.NewValWeights2(vals, weights)
}

func parseStringValMappingFile2(n *ast.Node) (*sgc7game.ValMapping2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseStringValMappingFile2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseStringValMappingFile2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	invals := []int{}
	outvals := []sgc7game.IVal{}

	for _, v := range lst {
		inv, err := goutils.String2Int64(strings.TrimSpace(v["in"]))
		if err != nil {
			goutils.Error("parseStringValMappingFile2:String2Int64",
				slog.String("in", v["in"]),
				goutils.Err(err))

			return nil, err
		}

		outvals = append(outvals, sgc7game.NewStrValEx(strings.TrimSpace(v["out"])))
		invals = append(invals, int(inv))
	}

	return sgc7game.NewValMapping2(invals, outvals)
}

func parseIntValMappingFile2(n *ast.Node) (*sgc7game.ValMapping2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValMappingFile2:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []map[string]string{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValMappingFile2:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	invals := []int{}
	outvals := []sgc7game.IVal{}

	for _, v := range lst {
		inv, err := goutils.String2Int64(strings.TrimSpace(v["in"]))
		if err != nil {
			goutils.Error("parseIntValMappingFile2:String2Int64",
				slog.String("in", v["in"]),
				goutils.Err(err))

			return nil, err
		}

		outv, err := goutils.String2Int64(strings.TrimSpace(v["out"]))
		if err != nil {
			goutils.Error("parseIntValMappingFile2:String2Int64",
				slog.String("out", v["out"]),
				goutils.Err(err))

			return nil, err
		}

		outvals = append(outvals, sgc7game.NewIntValEx(outv))
		invals = append(invals, int(inv))
	}

	return sgc7game.NewValMapping2(invals, outvals)
}

func parseIntValMappingFile(n *ast.Node) (*sgc7game.ValMapping2, error) {
	buf, err := n.MarshalJSON()
	if err != nil {
		goutils.Error("parseIntValMappingFile:MarshalJSON",
			goutils.Err(err))

		return nil, err
	}

	lst := []*mappingData{}

	err = sonic.Unmarshal(buf, &lst)
	if err != nil {
		goutils.Error("parseIntValMappingFile:Unmarshal",
			goutils.Err(err))

		return nil, err
	}

	invals := []int{}
	outvals := []sgc7game.IVal{}

	for _, v := range lst {
		inv, err := goutils.String2Int64(v.In)
		if err != nil {
			goutils.Error("parseIntValMappingFile:String2Int64",
				slog.String("in", v.In),
				goutils.Err(err))

			return nil, err
		}

		outv, err := goutils.String2Int64(v.Out)
		if err != nil {
			goutils.Error("parseIntValMappingFile:String2Int64",
				slog.String("out", v.Out),
				goutils.Err(err))

			return nil, err
		}

		outvals = append(outvals, sgc7game.NewIntValEx(outv))
		invals = append(invals, int(inv))
	}

	return sgc7game.NewValMapping2(invals, outvals)
}
