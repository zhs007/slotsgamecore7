package lowcode

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func SavePaytable(f *excelize.File, sheet string, pt *sgc7game.PayTables) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SavePaytable.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	num := 0
	maxCode := -1
	for code, arr := range pt.MapPay {
		num = len(arr)

		if code > maxCode {
			maxCode = code
		}
	}

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "Code")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "Symbol")

	for i := 1; i <= num; i++ {
		f.SetCellStr(sheet, goutils.Pos2Cell(1+i, 0), fmt.Sprintf("X%v", i))
	}

	for code := 0; code <= maxCode; code++ {
		arr, isok := pt.MapPay[code]
		if isok {
			f.SetCellInt(sheet, goutils.Pos2Cell(0, code+1), code)
			f.SetCellStr(sheet, goutils.Pos2Cell(1, code+1), pt.GetStringFromInt(code))

			for i := 0; i < num; i++ {
				f.SetCellInt(sheet, goutils.Pos2Cell(2+i, code+1), arr[i])
			}
		}
	}

	return nil
}

func SaveReels(f *excelize.File, sheet string, pt *sgc7game.PayTables, reels *sgc7game.ReelsData) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SaveReels.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	num := len(reels.Reels)
	maxline := 0

	for _, arr := range reels.Reels {
		if maxline < len(arr) {
			maxline = len(arr)
		}
	}

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "line")

	for i := 1; i <= num; i++ {
		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), fmt.Sprintf("R%v", i))
	}

	for y := 0; y < maxline; y++ {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y+1), y)

		for i := 0; i < num; i++ {
			if y < len(reels.Reels[i]) {
				f.SetCellStr(sheet, goutils.Pos2Cell(i+1, y+1), pt.GetStringFromInt(reels.Reels[i][y]))
			}
		}
	}

	return nil
}

func SaveLineData(f *excelize.File, sheet string, ld *sgc7game.LineData) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SaveLineData.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	maxline := len(ld.Lines)
	num := len(ld.Lines[0])

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "line")

	for i := 1; i <= num; i++ {
		f.SetCellStr(sheet, goutils.Pos2Cell(i, 0), fmt.Sprintf("R%v", i))
	}

	for y := 0; y < maxline; y++ {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y+1), y+1)

		for i := 0; i < num; i++ {
			f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y+1), ld.Lines[y][i])
		}
	}

	return nil
}

func SaveSymbolWeights(f *excelize.File, sheet string, pt *sgc7game.PayTables, vw *sgc7game.ValWeights2) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SaveSymbolWeights.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	maxCode := -1
	for code := range pt.MapPay {
		if code > maxCode {
			maxCode = code
		}
	}

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "val")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "weight")

	y := 1
	for c := 0; c <= maxCode; c++ {
		cv := sgc7game.NewIntValEx(c)
		w := vw.GetWeight(cv)
		if w > 0 {
			f.SetCellStr(sheet, goutils.Pos2Cell(0, y), pt.GetStringFromInt(c))
			f.SetCellInt(sheet, goutils.Pos2Cell(1, y), w)
		}
	}

	return nil
}

func SaveIntWeights(f *excelize.File, sheet string, vw *sgc7game.ValWeights2) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SaveIntWeights.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	vals := []int{}

	for _, v := range vw.Vals {
		vals = append(vals, v.Int())
	}

	sort.Slice(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "val")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "weight")

	for i := 0; i < len(vals); i++ {
		cv := sgc7game.NewIntValEx(vals[i])
		w := vw.GetWeight(cv)
		if w > 0 {
			f.SetCellInt(sheet, goutils.Pos2Cell(0, i+1), vals[i])
			f.SetCellInt(sheet, goutils.Pos2Cell(1, i+1), w)
		}
	}

	return nil
}

func SaveStrWeights(f *excelize.File, sheet string, vw *sgc7game.ValWeights2) error {
	_, err := f.NewSheet(sheet)
	if err != nil {
		goutils.Error("SaveStrWeights.NewSheet",
			slog.String("sheet", sheet),
			goutils.Err(err))

		return err
	}

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "val")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "weight")

	for i := 0; i < len(vw.Vals); i++ {
		if vw.Weights[i] > 0 {
			f.SetCellStr(sheet, goutils.Pos2Cell(0, i+1), vw.Vals[i].String())
			f.SetCellInt(sheet, goutils.Pos2Cell(1, i+1), vw.Weights[i])
		}
	}

	return nil
}
