package lowcode

import (
    "fmt"
    "path/filepath"
    "testing"

    "github.com/xuri/excelize/v2"
    sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func createXLSX(t *testing.T, headers []string, rows [][]string) string {
    t.Helper()

    f := excelize.NewFile()
    sheet := f.GetSheetName(0)

    // write headers to row 1
    for i, h := range headers {
        cell := fmt.Sprintf("%c%d", 'A'+i, 1)
        if err := f.SetCellValue(sheet, cell, h); err != nil {
            t.Fatalf("SetCellValue header failed: %v", err)
        }
    }

    // write data rows starting at row 2
    for r, row := range rows {
        for c, v := range row {
            cell := fmt.Sprintf("%c%d", 'A'+c, 2+r)
            if err := f.SetCellValue(sheet, cell, v); err != nil {
                t.Fatalf("SetCellValue row failed: %v", err)
            }
        }
    }

    dir := t.TempDir()
    fn := filepath.Join(dir, "test.xlsx")
    if err := f.SaveAs(fn); err != nil {
        t.Fatalf("SaveAs failed: %v", err)
    }

    return fn
}

func TestLoadSymbolsViewer_HappyPath(t *testing.T) {
    headers := []string{"Code", "Symbol", "Output", "Color"}
    rows := [][]string{
        {"1", "A", "a_out", "red"},
        {"2", "B", "", ""},
    }

    fn := createXLSX(t, headers, rows)
    sv, err := LoadSymbolsViewer(fn)
    if err != nil {
        t.Fatalf("LoadSymbolsViewer failed: %v", err)
    }

    if sv == nil {
        t.Fatalf("expected non-nil SymbolsViewer")
    }

    if len(sv.MapSymbols) != 2 {
        t.Fatalf("expected 2 symbols, got %d", len(sv.MapSymbols))
    }

    s1, ok := sv.MapSymbols[1]
    if !ok {
        t.Fatalf("missing code 1")
    }
    if s1.Symbol != "A" || s1.Output != "a_out" || s1.Color != "red" {
        t.Fatalf("unexpected data for code 1: %+v", s1)
    }

    s2, ok := sv.MapSymbols[2]
    if !ok {
        t.Fatalf("missing code 2")
    }
    if s2.Symbol != "B" {
        t.Fatalf("unexpected symbol for code 2: %s", s2.Symbol)
    }
    if s2.Color != DefaultSymbolColor {
        t.Fatalf("expected default color for code 2, got %s", s2.Color)
    }
}

func TestLoadSymbolsViewer_MissingCodeSkipped(t *testing.T) {
    headers := []string{"Code", "Symbol", "Output", "Color"}
    rows := [][]string{
        {"", "NoCode", "out", "blue"},
        {"5", "E", "eout", "green"},
    }

    fn := createXLSX(t, headers, rows)
    sv, err := LoadSymbolsViewer(fn)
    if err != nil {
        t.Fatalf("LoadSymbolsViewer failed: %v", err)
    }

    if _, ok := sv.MapSymbols[0]; ok {
        t.Fatalf("unexpected entry for code 0")
    }
    if _, ok := sv.MapSymbols[5]; !ok {
        t.Fatalf("expected entry for code 5")
    }
}

func TestLoadSymbolsViewer_InvalidCodeError(t *testing.T) {
    headers := []string{"Code", "Symbol", "Output", "Color"}
    rows := [][]string{
        {"abc", "Bad", "out", ""},
    }

    fn := createXLSX(t, headers, rows)
    _, err := LoadSymbolsViewer(fn)
    if err == nil {
        t.Fatalf("expected error for invalid code, got nil")
    }
}

func TestLoadSymbolsViewer_DuplicateCodeOverwrite(t *testing.T) {
    headers := []string{"Code", "Symbol", "Output", "Color"}
    rows := [][]string{
        {"3", "X", "xo", "red"},
        {"3", "Y", "yo", "green"},
    }

    fn := createXLSX(t, headers, rows)
    sv, err := LoadSymbolsViewer(fn)
    if err != nil {
        t.Fatalf("LoadSymbolsViewer failed: %v", err)
    }

    s, ok := sv.MapSymbols[3]
    if !ok {
        t.Fatalf("expected code 3 present")
    }
    if s.Symbol != "Y" || s.Output != "yo" || s.Color != "green" {
        t.Fatalf("duplicate code not overwritten by last row: %+v", s)
    }
}

func TestNewSymbolViewerFromPaytables(t *testing.T) {
    pt := &sgc7game.PayTables{
        MapPay:     make(map[int][]int),
        MapSymbols: map[string]int{"symA": 10},
    }

    sv := NewSymbolViewerFromPaytables(pt)
    if sv == nil {
        t.Fatalf("NewSymbolViewerFromPaytables returned nil")
    }
    v, ok := sv.MapSymbols[10]
    if !ok {
        t.Fatalf("expected map contains code 10")
    }
    if v.Symbol != "symA" || v.Output != "symA" || v.Color != DefaultSymbolColor {
        t.Fatalf("unexpected symbolviewer data: %+v", v)
    }
}
