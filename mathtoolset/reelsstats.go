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

func (rs *ReelStats) canAdd(i int, symbols []SymbolType) bool {
	s := symbols[i]
	n := rs.MapSymbols[s].Num

	for ri := i + 1; ri < len(symbols); ri++ {
		cs := symbols[ri]
		cn := rs.MapSymbols[cs].Num

		if cn <= n {
			return false
		}
	}

	return true
}

func (rs *ReelStats) GetCanAddSymbols(symbols []SymbolType) []SymbolType {
	lst := []SymbolType{}

	for ri := len(symbols) - 1; ri >= 0; ri-- {
		if rs.canAdd(ri, symbols) {
			lst = append(lst, symbols[ri])
		}
	}

	return lst
}

func (rs *ReelStats) BuildSymbols(excludeSymbols []SymbolType) []SymbolType {
	symbols := []SymbolType{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && !HasSymbol(excludeSymbols, s) {
			symbols = append(symbols, s)
		}
	}

	return symbols
}

func (rs *ReelStats) BuildSymbolsWithWeights(excludeSymbols []SymbolType) (*sgc7game.ValWeights, error) {
	vals := []int{}
	weights := []int{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && !HasSymbol(excludeSymbols, s) {
			vals = append(vals, int(s))
			weights = append(weights, v.Num)
		}
	}

	return sgc7game.NewValWeights(vals, weights)
}

func (rs *ReelStats) BuildSymbolsEx(symbols []SymbolType) []SymbolType {
	newarr := []SymbolType{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && HasSymbol(symbols, s) {
			newarr = append(newarr, s)
		}
	}

	return newarr
}

func (rs *ReelStats) BuildSymbolsWithWeightsEx(symbols []SymbolType) (*sgc7game.ValWeights, error) {
	vals := []int{}
	weights := []int{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && HasSymbol(symbols, s) {
			vals = append(vals, int(s))
			weights = append(weights, v.Num)
		}
	}

	return sgc7game.NewValWeights(vals, weights)
}

func (rs *ReelStats) BuildSymbols2(excludeSymbols1 []SymbolType, excludeSymbols2 []SymbolType) []SymbolType {
	symbols := []SymbolType{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && !HasSymbol(excludeSymbols1, s) && !HasSymbol(excludeSymbols2, s) {
			symbols = append(symbols, s)
		}
	}

	return symbols
}

func (rs *ReelStats) BuildSymbolsWithWeights2(excludeSymbols1 []SymbolType, excludeSymbols2 []SymbolType) (*sgc7game.ValWeights, error) {
	vals := []int{}
	weights := []int{}

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && !HasSymbol(excludeSymbols1, s) && !HasSymbol(excludeSymbols2, s) {
			vals = append(vals, int(s))
			weights = append(weights, v.Num)
		}
	}

	return sgc7game.NewValWeights(vals, weights)
}

func (rs *ReelStats) CountSymbolsNum(symbols []SymbolType) int {
	num := 0

	for s, v := range rs.MapSymbols {
		if v.Num > 0 && HasSymbol(symbols, s) {
			num += v.Num
		}
	}

	return num
}

func (rs *ReelStats) GetSymbolWithIndex(index int) SymbolType {
	i := 0
	for s := range rs.MapSymbols {
		if i == index {
			return s
		}

		i++
	}

	return -1
}

func (rs *ReelStats) Clone() *ReelStats {
	nrs := &ReelStats{
		MapSymbols: make(map[SymbolType]*SymbolStats),
	}

	for k, v := range rs.MapSymbols {
		nrs.AddSymbol(k, v.Num)
	}

	return nrs
}

func (rs *ReelStats) GetSymbolStats(symbol SymbolType) *SymbolStats {
	v, isok := rs.MapSymbols[symbol]
	if isok {
		return v
	}

	rs.MapSymbols[symbol] = newSymbolStats(symbol, 0)

	return rs.MapSymbols[symbol]
}

func (rs *ReelStats) AddSymbol(symbol SymbolType, num int) {
	ss := rs.GetSymbolStats(symbol)

	ss.Num += num
	rs.TotalSymbolNum += num
}

func (rs *ReelStats) RemoveSymbol(symbol SymbolType, num int) int {
	ss := rs.GetSymbolStats(symbol)

	if ss.Num > num {
		ss.Num -= num
		rs.TotalSymbolNum -= num

		return num
	}

	curnum := ss.Num

	rs.TotalSymbolNum -= ss.Num
	ss.Num = 0

	delete(rs.MapSymbols, symbol)

	return curnum
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

func (rss *ReelsStats) ClearEmptySymbols() {
	rss.Symbols = nil

	for _, r := range rss.Reels {
		for s, v := range r.MapSymbols {
			if v.Num > 0 {
				if !rss.HasSymbols(s) {
					rss.Symbols = append(rss.Symbols, s)
				}
			}
		}
	}
}

func (rss *ReelsStats) CloneWithMapping(mapSymbols *SymbolsMapping) *ReelsStats {
	nrss := &ReelsStats{}

	for _, rs := range rss.Reels {
		nrs := NewReelStats()

		for s, ss := range rs.MapSymbols {
			if mapSymbols != nil && mapSymbols.Has(s) {
				nrs.AddSymbol(mapSymbols.MapSymbols[s], ss.Num)
			} else {
				nrs.AddSymbol(s, ss.Num)
			}
		}

		nrss.Reels = append(nrss.Reels, nrs)
	}

	nrss.buildSortedSymbols()

	return nrss
}

func (rss *ReelsStats) Clone() *ReelsStats {
	nrss := &ReelsStats{}

	for _, rs := range rss.Reels {
		nrs := NewReelStats()

		for s, ss := range rs.MapSymbols {
			nrs.AddSymbol(s, ss.Num)
		}

		nrss.Reels = append(nrss.Reels, nrs)
	}

	nrss.buildSortedSymbols()

	return nrss
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
	defer f.Close()

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
			rs.AddSymbol(s, arr[i])
		}

		rss.Reels = append(rss.Reels, rs)
	}

	rss.buildSortedSymbols()

	return rss, nil
}

func BuildBasicReelsStats(reelnum int, paytables *sgc7game.PayTables, excludeSyms []SymbolType) (*ReelsStats, error) {
	rss := &ReelsStats{}

	for i := 0; i < reelnum; i++ {
		rs := NewReelStats()

		for s := range paytables.MapPay {
			if !HasSymbol(excludeSyms, SymbolType(s)) {
				rs.AddSymbol(SymbolType(s), 1)
			}
		}

		rss.Reels = append(rss.Reels, rs)
	}

	rss.buildSortedSymbols()

	return rss, nil
}

func BuildBasicReelsStatsEx(reelnum int, syms []SymbolType) (*ReelsStats, error) {
	rss := &ReelsStats{}

	for i := 0; i < reelnum; i++ {
		rs := NewReelStats()

		for _, s := range syms {
			rs.AddSymbol(SymbolType(s), 1)
		}

		rss.Reels = append(rss.Reels, rs)
	}

	rss.buildSortedSymbols()

	return rss, nil
}
