package lowcode

import (
	"log/slog"
	"strings"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// DefaultSymbolColor is the default color used when a symbol color is not provided.
const DefaultSymbolColor = "low"

// SymbolViewerData holds presentation data for a symbol.
type SymbolViewerData struct {
	Code   int
	Symbol string
	Output string
	Color  string
}

// SymbolViewer is a helper for mapping symbol code -> presentation data.
//
// Note: an alias "SymbolsViewer" is provided below for backwards compatibility
// with code that references the old exported name.
type SymbolViewer struct {
	MapSymbols map[int]*SymbolViewerData
}

// SymbolsViewer is an alias for SymbolViewer to keep compatibility with
// existing callers that expect the older name.
type SymbolsViewer = SymbolViewer

// NewSymbolViewerFromPaytables builds a SymbolViewer from game PayTables.
// It uses the symbol string as the default Output and assigns a default color
// when none is provided.
func NewSymbolViewerFromPaytables(paytables *sgc7game.PayTables) *SymbolsViewer {
	viewer := &SymbolViewer{
		MapSymbols: make(map[int]*SymbolViewerData),
	}

	for symStr, code := range paytables.MapSymbols {
		svd := &SymbolViewerData{
			Code:   code,
			Symbol: symStr,
			Output: symStr,
			Color:  DefaultSymbolColor,
		}

		viewer.MapSymbols[code] = svd
	}

	return (*SymbolsViewer)(viewer)
}

// LoadSymbolsViewer loads a symbols viewer from an Excel file. Required column:
//   - code (int)
// Optional columns:
//   - symbol, output, color
// The header matching is case-insensitive and trimmed.
func LoadSymbolsViewer(fn string) (*SymbolsViewer, error) {
	type tmpRow struct {
		code   *int
		symbol string
		output string
		color  string
		row    int
	}

	rows := map[int]*tmpRow{}

	// transform: normalize cells (header matching is lowercased)
	err := sgc7game.LoadExcel(fn, "", func(x int, str string) string {
		return strings.ToLower(strings.TrimSpace(str))
	}, func(x int, y int, header string, data string) error {
		data = strings.TrimSpace(data)
		r := rows[y]
		if r == nil {
			r = &tmpRow{row: y}
			rows[y] = r
		}

		switch header {
		case "code":
			if data == "" {
				// empty code cell => treat as missing code for this row
				return nil
			}

			v, err := goutils.String2Int64(data)
			if err != nil {
				goutils.Error("LoadSymbolsViewer:LoadExcel:String2Int64",
					slog.String("header", header),
					slog.String("data", data),
					slog.Int("row", y),
					goutils.Err(err))

				return err
			}
			iv := int(v)
			r.code = &iv
		case "symbol":
			r.symbol = data
		case "output":
			r.output = data
		case "color":
			r.color = data
		}

		return nil
	})
	if err != nil {
		goutils.Error("LoadSymbolsViewer:LoadExcel",
			slog.String("fn", fn),
			goutils.Err(err))

		return nil, err
	}

	sv := &SymbolViewer{
		MapSymbols: make(map[int]*SymbolViewerData),
	}

	for _, r := range rows {
		if r.code == nil {
			// skip rows without a code, but warn with row info
			goutils.Warn("LoadSymbolsViewer: missing code, skipping row",
				slog.Int("row", r.row))
			continue
		}

		code := *r.code

		// if duplicate code, warn and overwrite with last-seen
		if _, ok := sv.MapSymbols[code]; ok {
			goutils.Warn("LoadSymbolsViewer: duplicate code, overwriting",
				slog.Int("code", code),
				slog.Int("row", r.row))
		}

		color := r.color
		if color == "" {
			color = DefaultSymbolColor
		}

		sv.MapSymbols[code] = &SymbolViewerData{
			Code:   code,
			Symbol: r.symbol,
			Output: r.output,
			Color:  color,
		}
	}

	return (*SymbolsViewer)(sv), nil
}
