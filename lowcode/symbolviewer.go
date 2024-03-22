package lowcode

import (
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type SymbolViewerData struct {
	Code   int
	Symbol string
	Output string
	Color  string
}

type SymbolsViewer struct {
	MapSymbols map[int]*SymbolViewerData
}

func NewSymbolViewerFromPaytables(paytables *sgc7game.PayTables) *SymbolsViewer {
	symViewer := &SymbolsViewer{
		MapSymbols: make(map[int]*SymbolViewerData),
	}

	for key, sym := range paytables.MapSymbols {
		svd := &SymbolViewerData{
			Code:   sym,
			Symbol: key,
			Output: key,
			Color:  "low",
		}

		symViewer.MapSymbols[sym] = svd
	}

	return symViewer
}

func LoadSymbolsViewer(fn string) (*SymbolsViewer, error) {
	symbols := []int{}
	symbolstr := []string{}
	outputs := []string{}
	colors := []string{}

	err := sgc7game.LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		data = strings.TrimSpace(data)
		if header == "code" {
			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadSymbolsViewer:LoadExcel:String2Int64",
					slog.String("header", header),
					slog.String("data", data),
					goutils.Err(err))

				return err
			}

			symbols = append(symbols, int(v))
		} else if header == "symbol" {
			symbolstr = append(symbolstr, data)
		} else if header == "output" {
			outputs = append(outputs, data)
		} else if header == "color" {
			colors = append(colors, data)
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadSymbolsViewer:LoadExcel",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	sv := &SymbolsViewer{
		MapSymbols: make(map[int]*SymbolViewerData),
	}

	for i, v := range symbols {
		svd := &SymbolViewerData{
			Code:   v,
			Symbol: symbolstr[i],
			Output: outputs[i],
			Color:  colors[i],
		}

		sv.MapSymbols[v] = svd
	}

	return sv, nil
}
