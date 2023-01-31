package mathtoolset

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type SymbolStats struct {
	Symbol SymbolType
	Num    int
}

func newSymbolStats(s SymbolType, num int) *SymbolStats {
	return &SymbolStats{
		Symbol: s,
		Num:    num,
	}
}

type ReelStats struct {
	MapSymbols     map[SymbolType]*SymbolStats
	TotalSymbolNum int
}

func (rs *ReelStats) GetSymbolStats(symbol SymbolType) *SymbolStats {
	v, isok := rs.MapSymbols[symbol]
	if isok {
		return v
	}

	rs.MapSymbols[symbol] = newSymbolStats(symbol, 0)

	return rs.MapSymbols[symbol]
}

func NewReelStats() *ReelStats {
	return &ReelStats{
		MapSymbols: make(map[SymbolType]*SymbolStats),
	}
}

func BuildReelStats(reel []int, mapSymbols *SymbolsMapping) (*ReelStats, error) {
	if len(reel) == 0 {
		goutils.Error("ReelStats.AnalyzeReel",
			zap.Error(ErrInvalidReel))

		return nil, ErrInvalidReel
	}

	rs := NewReelStats()

	for _, s := range reel {
		st := SymbolType(s)
		if mapSymbols != nil && mapSymbols.Has(st) {
			st = mapSymbols.MapSymbols[st]
		}

		ss := rs.GetSymbolStats(st)
		ss.Num++
	}

	rs.TotalSymbolNum = len(reel)

	return rs, nil
}

type ReelsStats struct {
	Reels   []*ReelStats
	Symbols []SymbolType
}

func (rss *ReelsStats) Add(symbol SymbolType) {

}

func (rss *ReelsStats) HasSymbols(symbol SymbolType) bool {
	for _, v := range rss.Symbols {
		if v == symbol {
			return true
		}
	}

	return false
}

func (rss *ReelsStats) buildSortedSymbols() {
	rss.Symbols = nil

	for _, r := range rss.Reels {
		for s := range r.MapSymbols {
			if !rss.HasSymbols(s) {
				rss.Symbols = append(rss.Symbols, s)
			}
		}
	}

	sort.Slice(rss.Symbols, func(i, j int) bool {
		return rss.Symbols[i] < rss.Symbols[j]
	})
}

func (rss *ReelsStats) GetNum(reelindex int, symbol SymbolType, symbol2 SymbolType,
	wilds []SymbolType, irstype InReelSymbolType) int {

	if irstype == IRSTypeAll {
		return rss.Reels[reelindex].TotalSymbolNum
	}

	if irstype == IRSTypeSymbol2 || irstype == IRSTypeSymbol2AndWild {
		if symbol2 < 0 {
			return -1
		}

		s2s := rss.Reels[reelindex].GetSymbolStats(symbol2)

		if irstype == IRSTypeSymbol2 {
			return s2s.Num
		}

		wildnum := 0
		for _, w := range wilds {
			if w == symbol {
				continue
			}

			ws := rss.Reels[reelindex].GetSymbolStats(w)
			if ws.Num > 0 {
				wildnum += ws.Num
			}
		}

		return s2s.Num + wildnum
	}

	wildnum := 0

	if irstype == IRSTypeWild ||
		irstype == IRSTypeNoSymbolAndNoWild ||
		irstype == IRSTypeSymbolAndWild ||
		irstype == IRSTypeNoWild {

		for _, w := range wilds {
			if w == symbol {
				continue
			}

			ws := rss.Reels[reelindex].GetSymbolStats(w)
			if ws.Num > 0 {
				wildnum += ws.Num
			}
		}

		if irstype == IRSTypeWild {
			return wildnum
		}

		if irstype == IRSTypeNoWild {
			return rss.Reels[reelindex].TotalSymbolNum - wildnum
		}
	}

	ss := rss.Reels[reelindex].GetSymbolStats(symbol)

	if irstype == IRSTypeSymbol {
		return ss.Num
	}

	if irstype == IRSTypeSymbolAndWild {
		return ss.Num + wildnum
	}

	if irstype == IRSTypeNoSymbolAndNoWild {
		return rss.Reels[reelindex].TotalSymbolNum - ss.Num - wildnum
	}

	if irstype == IRSTypeNoSymbol {
		return rss.Reels[reelindex].TotalSymbolNum - ss.Num
	}

	return -1
}

func (rss *ReelsStats) GetScatterNum(reelindex int, symbol SymbolType, irstype InReelSymbolType, height int) int {
	ss := rss.Reels[reelindex].GetSymbolStats(symbol)

	if irstype == IRSTypeSymbol {
		return ss.Num * height
	}

	if irstype == IRSTypeNoSymbol {
		return rss.Reels[reelindex].TotalSymbolNum - ss.Num*height
	}

	return -1
}

func (rss *ReelsStats) GetSymbolNum(reelindex int, symbol SymbolType, wilds []SymbolType) int {
	ss := rss.Reels[reelindex].GetSymbolStats(symbol)

	wildnum := 0
	for _, w := range wilds {
		if w == symbol {
			continue
		}

		ws := rss.Reels[reelindex].GetSymbolStats(w)
		if ws.Num > 0 {
			wildnum += ws.Num
		}
	}

	return ss.Num + wildnum
}

func (rss *ReelsStats) GetSymbolNumNoWild(reelindex int, symbol SymbolType, wilds []SymbolType) int {
	ss := rss.Reels[reelindex].GetSymbolStats(symbol)

	if HasSymbol(wilds, symbol) {
		wildnum := 0
		for _, w := range wilds {
			if w == symbol {
				continue
			}

			ws := rss.Reels[reelindex].GetSymbolStats(w)
			if ws.Num > 0 {
				wildnum += ws.Num
			}
		}

		return ss.Num + wildnum
	}

	return ss.Num
}

func (rss *ReelsStats) GetReelLengthNoSymbol(reelindex int, symbol SymbolType, wilds []SymbolType) int {
	ss := rss.Reels[reelindex].GetSymbolStats(symbol)

	if HasSymbol(wilds, symbol) {
		wildnum := 0
		for _, w := range wilds {
			if w == symbol {
				continue
			}

			ws := rss.Reels[reelindex].GetSymbolStats(w)
			if ws.Num > 0 {
				wildnum += ws.Num
			}
		}

		return rss.Reels[reelindex].TotalSymbolNum - (ss.Num + wildnum)
	}

	return rss.Reels[reelindex].TotalSymbolNum - ss.Num
}

func (rss *ReelsStats) GetReelLength(reelindex int) int {
	return rss.Reels[reelindex].TotalSymbolNum
}

func (rss *ReelsStats) SaveExcel(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetName(0)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "symbol")
	for i := range rss.Reels {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+1, 0), fmt.Sprintf("reel%v", i+1))
	}

	y := 1

	for _, s := range rss.Symbols {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, y), int(s))

		for i, reel := range rss.Reels {
			statsSymbol, isok := reel.MapSymbols[s]
			if isok {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), statsSymbol.Num)
			} else {
				f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), 0)
			}
		}

		y++
	}

	for i, rs := range rss.Reels {
		f.SetCellInt(sheet, goutils.Pos2Cell(i+1, y), rs.TotalSymbolNum)
	}

	return f.SaveAs(fn)
}

func BuildReelsStats(reels *sgc7game.ReelsData, mapSymbols *SymbolsMapping) (*ReelsStats, error) {
	rss := &ReelsStats{
		Reels: make([]*ReelStats, len(reels.Reels)),
	}

	for i, r := range reels.Reels {
		rs, err := BuildReelStats(r, mapSymbols)
		if err != nil {
			goutils.Error("BuildReelsStats:BuildReelStats",
				zap.Error(err))

			return nil, err
		}

		rss.Reels[i] = rs
	}

	rss.buildSortedSymbols()

	return rss, nil
}

func getReelNum(mapcolname map[int]string) (int, error) {
	lst := []int{}

	for _, v := range mapcolname {
		if strings.Index(v, "reel") == 0 {
			arr := strings.Split(v, "reel")
			if len(arr) != 2 {
				goutils.Error("getReelNum",
					zap.String("colname", v),
					zap.Error(ErrInvalidReelsStatsExcelColname))

				return -1, ErrInvalidReelsStatsExcelColname
			}

			i64, err := goutils.String2Int64(arr[1])
			if err != nil {
				goutils.Error("getReelNum:String2Int64",
					zap.String("colname", v),
					zap.Error(err))

				return -1, err
			}

			lst = append(lst, int(i64))
		}
	}

	sort.Slice(lst, func(i, j int) bool {
		return lst[i] < lst[j]
	})

	if len(lst) != lst[len(lst)-1] {
		goutils.Error("getReelNum",
			goutils.JSON("mapcolname", mapcolname),
			zap.Error(ErrInvalidReelsStatsExcelColname))

		return -1, ErrInvalidReelsStatsExcelColname
	}

	return len(lst), nil
}

func LoadReelsStats(fn string) (*ReelsStats, error) {
	f, err := excelize.OpenFile(fn)
	if err != nil {
		goutils.Error("LoadReelsStats:OpenFile",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	lstname := f.GetSheetList()
	if len(lstname) <= 0 {
		goutils.Error("LoadReelsStats:GetSheetList",
			goutils.JSON("SheetList", lstname),
			zap.String("fn", fn),
			zap.Error(ErrInvalidReelsStatsExcelFile))

		return nil, ErrInvalidReelsStatsExcelFile
	}

	rows, err := f.GetRows(lstname[0])
	if err != nil {
		goutils.Error("LoadReelsStats:GetRows",
			zap.String("fn", fn),
			zap.Error(err))

		return nil, err
	}

	mapcolname := make(map[int]string)
	reelnum := 0
	rss := &ReelsStats{}
	mapSymbols := make(map[SymbolType][]int)

	for y, row := range rows {
		if y == 0 {
			for x, colCell := range row {
				mapcolname[x] = strings.ToLower(colCell)
			}

			reelnum, err = getReelNum(mapcolname)
			if err != nil {
				goutils.Error("LoadReelsStats:getReelNum",
					zap.String("fn", fn),
					zap.Error(err))

				return nil, err
			}
		} else {

			symbol := -1
			vals := []int{}

			for x, colCell := range row {
				colname, isok := mapcolname[x]
				if isok {
					if colname == "symbol" {
						i64, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadReelsStats:String2Int64",
								zap.String("val", colCell),
								zap.Error(err))

							break
							// return nil, err
						}

						symbol = int(i64)
					} else if strings.Index(colname, "reel") == 0 {
						v, err := goutils.String2Int64(colCell)
						if err != nil {
							goutils.Error("LoadReelsStats:String2Int64",
								zap.String("val", colCell),
								zap.Error(err))

							return nil, err
						}

						vals = append(vals, int(v))
					}
				}
			}

			if symbol >= 0 {
				mapSymbols[SymbolType(symbol)] = vals
			}
		}
	}

	for i := 0; i < reelnum; i++ {
		rs := NewReelStats()

		for s, arr := range mapSymbols {
			ss := rs.GetSymbolStats(s)
			ss.Num = arr[i]
		}

		rss.Reels = append(rss.Reels, rs)
	}

	rss.buildSortedSymbols()

	return rss, nil
}
