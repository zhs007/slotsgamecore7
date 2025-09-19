package mathtoolset

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type SymbolsWinsFileMode int

const (
	SWFModeRTP        SymbolsWinsFileMode = 1
	SWFModeWins       SymbolsWinsFileMode = 2
	SWFModeWinsNum    SymbolsWinsFileMode = 3
	SWFModeWinsNumPer SymbolsWinsFileMode = 4
)

type SymbolWinsStats struct {
	Symbol  SymbolType
	WinsNum []int64
	Wins    []int64
	Total   int64
}

func (sws *SymbolWinsStats) IsFine() bool {
	for i := 0; i+1 < len(sws.WinsNum); i++ {
		if sws.WinsNum[i] > sws.WinsNum[i+1] {
			return false
		}
	}

	return true
}

func (sws *SymbolWinsStats) Clone() *SymbolWinsStats {
	nsws := &SymbolWinsStats{
		Symbol:  sws.Symbol,
		WinsNum: make([]int64, len(sws.WinsNum)),
		Wins:    make([]int64, len(sws.Wins)),
		Total:   sws.Total,
	}

	copy(nsws.Wins, sws.Wins)
	copy(nsws.WinsNum, sws.WinsNum)

	return nsws
}

func (sws *SymbolWinsStats) Merge(sws1 *SymbolWinsStats) {
	for i, v := range sws1.Wins {
		sws.Wins[i] += v
	}

	for i, v := range sws1.WinsNum {
		sws.WinsNum[i] += v
	}

	sws.Total += sws1.Total
}

func (sws *SymbolWinsStats) MergeWithMulti(sws1 *SymbolWinsStats, multi int) {
	for i, v := range sws1.Wins {
		sws.Wins[i] += v
	}

	for i, v := range sws1.WinsNum {
		sws.WinsNum[i] += v
	}

	sws.Total += sws1.Total * int64(multi)
}

func newSymbolWinsStats(symbol SymbolType, num int) *SymbolWinsStats {
	return &SymbolWinsStats{
		Symbol:  symbol,
		WinsNum: make([]int64, num),
		Wins:    make([]int64, num),
	}
}

type SymbolsWinsStats struct {
	MapSymbols map[SymbolType]*SymbolWinsStats
	Symbols    []SymbolType
	Num        int
	TotalBet   int64
	TotalWins  int64
}

func (ssws *SymbolsWinsStats) Clone() *SymbolsWinsStats {
	nssws := &SymbolsWinsStats{
		MapSymbols: make(map[SymbolType]*SymbolWinsStats),
		Symbols:    make([]SymbolType, len(ssws.Symbols)),
		Num:        ssws.Num,
		TotalBet:   ssws.TotalBet,
		TotalWins:  ssws.TotalWins,
	}

	copy(nssws.Symbols, ssws.Symbols)
	for k, v := range ssws.MapSymbols {
		nssws.MapSymbols[k] = v.Clone()
	}

	return nssws
}

func (ssws *SymbolsWinsStats) CountRTP() float64 {
	return float64(ssws.TotalWins) / float64(ssws.TotalBet)
}

func (ssws *SymbolsWinsStats) IsFine() bool {
	for _, v := range ssws.MapSymbols {
		if !v.IsFine() {
			return false
		}
	}

	return true
}

func (ssws *SymbolsWinsStats) Merge(ssws1 *SymbolsWinsStats) {
	for s, v := range ssws1.MapSymbols {
		sws := ssws.GetSymbolWinsStats(s)

		sws.Merge(v)
	}
}

func (ssws *SymbolsWinsStats) MergeWithMulti(ssws1 *SymbolsWinsStats, multi int) {
	for s, v := range ssws1.MapSymbols {
		sws := ssws.GetSymbolWinsStats(s)

		sws.Merge(v)
	}
}

func (ssws *SymbolsWinsStats) GetSymbolWinsStats(symbol SymbolType) *SymbolWinsStats {
	sws, isok := ssws.MapSymbols[symbol]
	if isok {
		return sws
	}

	ssws.MapSymbols[symbol] = newSymbolWinsStats(symbol, ssws.Num)
	ssws.Symbols = append(ssws.Symbols, symbol)

	return ssws.MapSymbols[symbol]
}

func (ssws *SymbolsWinsStats) onBuildEnd() {
	ssws.Symbols = nil
	for s, v := range ssws.MapSymbols {
		v.Total = ssws.TotalBet
		ssws.Symbols = append(ssws.Symbols, s)
	}

	sort.Slice(ssws.Symbols, func(i, j int) bool {
		return ssws.Symbols[i] < ssws.Symbols[j]
	})
}

func (ssws *SymbolsWinsStats) SaveExcelSheet(f *excelize.File, fm SymbolsWinsFileMode) error {
	sheet := "rtp"

	if fm == SWFModeWins {
		sheet = "wins"
	} else if fm == SWFModeWinsNum {
		sheet = "wins number"
	} else if fm == SWFModeWinsNumPer {
		sheet = "wins number percent"
	}

	f.NewSheet(sheet)

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "symbol")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "total")

	si := 2

	for i := 0; i < ssws.Num; i++ {
		f.SetCellStr(sheet, goutils.Pos2Cell(i+si, 0), fmt.Sprintf("X%v", i+1))
	}

	f.SetCellStr(sheet, goutils.Pos2Cell(si+ssws.Num, 0), "sum")

	y := 1
	trtp := 0.0
	twinnum := int64(0)
	twinnumper := 0.0
	twins := int64(0)

	for _, s := range ssws.Symbols {
		sws := ssws.GetSymbolWinsStats(s)

	f.SetCellValue(sheet, goutils.Pos2Cell(0, y), int(s))
		f.SetCellValue(sheet, goutils.Pos2Cell(1, y), sws.Total)

		rtp := 0.0
		winnum := int64(0)
		winnumper := 0.0
		wins := int64(0)

		for i := 0; i < ssws.Num; i++ {
			if fm == SWFModeRTP {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), float64(sws.Wins[i])*100.0/float64(sws.Total))

				rtp += float64(sws.Wins[i]) * 100.0 / float64(sws.Total)
			} else if fm == SWFModeWinsNum {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), sws.WinsNum[i])

				winnum += sws.WinsNum[i]
			} else if fm == SWFModeWinsNumPer {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), float64(sws.WinsNum[i])*100.0/float64(sws.Total))

				winnumper += float64(sws.WinsNum[i]) * 100.0 / float64(sws.Total)
			} else {
				f.SetCellValue(sheet, goutils.Pos2Cell(i+si, y), sws.Wins[i])

				wins += sws.Wins[i]
			}
		}

		if fm == SWFModeRTP {
			f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), rtp)
			trtp += rtp
		} else if fm == SWFModeWinsNum {
			f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), winnum)
			twinnum += winnum
		} else if fm == SWFModeWinsNumPer {
			f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), winnum)
			twinnumper += winnumper
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), wins)
			twins += wins
		}

		y++
	}

	if fm == SWFModeRTP {
		f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), trtp)
	} else if fm == SWFModeWinsNum {
		f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), twinnum)
	} else if fm == SWFModeWinsNumPer {
		f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), twinnumper)
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(si+ssws.Num, y), twins)
	}

	return nil
}

func (ssws *SymbolsWinsStats) SaveExcel(fn string, fms []SymbolsWinsFileMode) error {
	f := excelize.NewFile()

	for _, fm := range fms {
		err := ssws.SaveExcelSheet(f, fm)
		if err != nil {
			goutils.Error("SymbolsWinsStats.SaveExcel",
				goutils.Err(err))

			return err
		}
	}

	f.DeleteSheet(f.GetSheetName(0))

	return f.SaveAs(fn)
}

func newSymbolsWinsStatsWithPaytables(paytables *sgc7game.PayTables, symbols []SymbolType) *SymbolsWinsStats {
	num := 0
	for _, arr := range paytables.MapPay {
		if len(arr) > num {
			num = len(arr)
		}
	}

	ssws := NewSymbolsWinsStats(num)

	for s := range paytables.MapPay {
		if HasSymbol(symbols, SymbolType(s)) {
			ssws.GetSymbolWinsStats(SymbolType(s))
		}
	}

	return ssws
}

// lst is like [S, W, S+W, No S+W, All]
func CalcSymbolWins(rss *ReelsStats, wilds []SymbolType, symbol SymbolType, symbol2 SymbolType, lst []InReelSymbolType) (int64, error) {
	curwins := int64(1)

	for i, t := range lst {
		cn := rss.GetNum(i, symbol, symbol2, wilds, t)
		if cn < 0 {
			goutils.Error("CalcSymbolWins:GetNum",
				slog.Int("InReelSymbolType", int(t)),
				goutils.Err(ErrInvalidInReelSymbolType))

			return 0, ErrInvalidInReelSymbolType
		}

		if cn == 0 {
			return 0, nil
		}

		curwins *= int64(cn)
	}

	return curwins, nil
}

// calcNotWildWins - 这个接口只能用于处理非wild的赢得，symbol必须不是wild，且只处理 S 开头的情况
func calcNotWildWins(rss *ReelsStats, wilds []SymbolType, symbol SymbolType, num int) int64 {
	curwins := int64(rss.GetNum(0, symbol, -1, wilds, IRSTypeSymbol))

	for i := 1; i < num; i++ {
		cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeSymbolAndWild)

		if cn <= 0 {
			return 0
		}

		curwins *= int64(cn)
	}

	if num == len(rss.Reels) {
		return curwins
	}

	curwins *= int64(rss.GetNum(num, symbol, -1, wilds, IRSTypeNoSymbolAndNoWild))

	for i := num + 1; i < len(rss.Reels); i++ {
		cn := rss.GetNum(i, symbol, -1, wilds, IRSTypeAll)

		if cn <= 0 {
			return 0
		}

		curwins *= int64(cn)
	}

	return curwins
}

// lst is like [S, W, S+W, All, All]
// 如果是 www 开头，且 w 作为 symbol符号，这里会需要减去比3w大的情况
func calcSymbolWinsFromList(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, ci int, num int, lst []InReelSymbolType, wildPayoutSymbol SymbolType, wildNum int) (int64, error) {

	if ci == num {
		// 第2种情况
		if wildPayoutSymbol != symbol && wildNum > 0 && ci == wildNum {
			if IsFirstWild(lst, ci) {
				return 0, nil
			}
		}

		if ci == len(rss.Reels) {
			return CalcSymbolWins(rss, wilds, symbol, -1, lst)
		}

		// 如果 wild 赔付就是 sumbol，那么需要处理3w以后的可能性
		if wildPayoutSymbol == symbol {
			// 如果要计算 3a，如果是w开头，且前面至少有1个a，那么第4个只要不是a和w就好了
			if IsFirstWild(lst, num) {
				// 如果是 3w，用加法来算
				curwin := int64(0)
				ps := paytables.MapPay[int(symbol)][num-1]

				for _, s := range symbols {
					if s == symbol {
						continue
					}

					if HasSymbol(wilds, s) {
						continue
					}

					parr := paytables.MapPay[int(s)]
					cn := -1
					for j := num; j < len(parr); j++ {
						if parr[j] > ps {
							cn = j

							break
						}
					}

					// 如果4b大于3w，那么5b也一定大于3w，所以就彻底排除接下来如果是符号b的情况
					if cn == num {
						continue
					}

					lst[ci] = IRSTypeSymbol2
					for j := ci + 1; j < len(parr); j++ {
						lst[j] = IRSTypeAll
					}

					cw, err := CalcSymbolWins(rss, wilds, symbol, s, lst)
					if err != nil {
						goutils.Error("calcSymbolWinsFromList:CalcSymbolWins",
							goutils.Err(err))

						return 0, err
					}

					curwin += cw

					// 如果 4b、5b 都小于 3w，那么需要加入接下来是b的情况
					if cn < 0 {
						continue
					} else {
						// 剩下就是4b小于3w，且5b大于3w，那么我们需要加入wwwbx，减去wwwbb和wwwbw
						// 如果是6个，则加入 wwwbxx，减去 wwwbbx 和 wwwbwx
						lst[ci+1] = IRSTypeSymbol2AndWild

						for j := ci + 2; j < len(parr); j++ {
							lst[j] = IRSTypeAll
						}

						cw0, err := CalcSymbolWins(rss, wilds, symbol, s, lst)
						if err != nil {
							goutils.Error("calcSymbolWinsFromList:CalcSymbolWins",
								goutils.Err(err))

							return 0, err
						}

						curwin -= cw0
					}
				}

				return curwin, nil
			}

			lst[ci] = IRSTypeNoSymbolAndNoWild

			return CalcSymbolWins(rss, wilds, symbol, -1, lst)
		}

		lst[ci] = IRSTypeNoSymbolAndNoWild
		for j := ci + 1; j < len(rss.Reels); j++ {
			lst[j] = IRSTypeAll
		}

		return CalcSymbolWins(rss, wilds, symbol, -1, lst)
	}

	// 第2种情况
	if wildPayoutSymbol != symbol && wildNum > 0 && ci == wildNum {
		if IsFirstWild(lst, ci) {
			return 0, nil
		}
	}

	curwins := int64(0)

	for t := IRSTypeSymbol; t <= IRSTypeWild; t++ {
		lst[ci] = t

		cw, err := calcSymbolWinsFromList(paytables, rss, symbols, wilds, symbol, ci+1, num, lst, wildPayoutSymbol, wildNum)
		if err != nil {
			goutils.Error("calcSymbolWinsFromList:calcSymbolWinsFromList",
				slog.Int("ci", ci),
				slog.Int("InReelSymbolType", int(t)),
				goutils.Err(err))

			return 0, err
		}

		curwins += cw
	}

	return curwins, nil
}

// calcSymbolFirstWildWins - 这个接口只能用于处理非wild赢得，symbol必须不是wild，且以wild开头
func calcSymbolFirstWildWins(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, num int, wildPayoutSymbol SymbolType, wildNum int) (int64, error) {

	lst := NewInReelSymbolTypeArr(len(rss.Reels))
	lst[0] = IRSTypeWild

	return calcSymbolWinsFromList(paytables, rss, symbols, wilds, symbol, 1, num, lst, wildPayoutSymbol, wildNum)
}

// CalcSymbolWinsInReelsWithLine -
// symbol 不可能是 wild
func CalcSymbolWinsInReelsWithLine(paytables *sgc7game.PayTables, rss *ReelsStats, symbols []SymbolType,
	wilds []SymbolType, symbol SymbolType, num int, wildPayoutSymbol SymbolType) (int64, error) {

	// 分2种情况，分别是s开头和w开头
	// s开头时，wild是确定的，所以计算起来非常简单
	curwins := calcNotWildWins(rss, wilds, symbol, num)

	// 如果是w开头，还会分为2种情况
	// 1. 当前s就是最大赔付，所以3个w就是3个s，但这时需要考虑4个b比3个s大的情况，这样 wwwbx 就是 bbbbx，而不是 sssbx
	// 2. 当前s不是最大赔付，这时需要考虑2个w比3个s大的情况，也就是 wwsxx，应该算作 aasxx
	cw, err := calcSymbolFirstWildWins(paytables, rss, symbols, wilds, symbol, num, wildPayoutSymbol,
		analyzeWildNum(paytables, symbol, num, wildPayoutSymbol))
	if err != nil {
		goutils.Error("CalcSymbolWinsInReelsWithLine:calcSymbolFirstWildWins",
			goutils.Err(err))

		return 0, err
	}

	curwins += cw

	return curwins, nil
}

func analyzeWildNum(paytables *sgc7game.PayTables, symbol SymbolType, num int, wild SymbolType) int {
	if symbol == wild {
		return 0
	}

	sp := paytables.MapPay[int(symbol)][num-1]
	warr := paytables.MapPay[int(wild)]

	for i := 0; i < len(warr); i++ {
		if warr[i] >= sp {
			return i + 1
		}
	}

	return 0
}

func getMaxPayoutSymbol(paytables *sgc7game.PayTables, symbols []SymbolType, num int) SymbolType {
	maxs := SymbolType(-1)
	maxpayout := -1

	for _, s := range symbols {
		if maxs < 0 {
			maxs = s
			maxpayout = paytables.MapPay[int(s)][num-1]
		} else if maxpayout < paytables.MapPay[int(s)][num-1] {
			maxs = s
			maxpayout = paytables.MapPay[int(s)][num-1]
		}
	}

	return maxs
}

func AnalyzeReelsWithLine(paytables *sgc7game.PayTables, reels *sgc7game.ReelsData,
	symbols []SymbolType, wilds []SymbolType, mapSymbols *SymbolsMapping, betMul int, lineNum int) (*SymbolsWinsStats, error) {

	rss, err := BuildReelsStats(reels, mapSymbols)
	if err != nil {
		goutils.Error("AnalyzeReelsWithLine:BuildReelsStats",
			goutils.Err(err))

		return nil, err
	}

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, arr := range reels.Reels {
		ssws.TotalBet *= int64(len(arr))
	}

	ssws.TotalBet *= int64(betMul)

	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcSymbolWinsInReelsWithLine(paytables, rss, symbols, wilds, s, i+1, getMaxPayoutSymbol(paytables, symbols, i+1))
					if err != nil {
						goutils.Error("AnalyzeReelsWithLine:CalcSymbolWinsInReelsWithLine",
							goutils.Err(err))

						return nil, err
					}

					sws.WinsNum[i] = cw
					sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(lineNum)

					ssws.TotalWins += int64(arrPay[i]) * sws.WinsNum[i] * int64(lineNum)
				}
			}
		}
	}

	ssws.onBuildEnd()

	return ssws, nil
}

func AnalyzeReelsWithLineEx(paytables *sgc7game.PayTables, rss *ReelsStats,
	symbols []SymbolType, wilds []SymbolType, betMul int, lineNum int) (*SymbolsWinsStats, error) {

	ssws := newSymbolsWinsStatsWithPaytables(paytables, symbols)

	ssws.TotalBet = 1
	for _, rs := range rss.Reels {
		ssws.TotalBet *= int64(rs.TotalSymbolNum)
	}

	ssws.TotalBet *= int64(betMul)

	ssws.TotalWins = 0

	for _, s := range symbols {
		sws := ssws.GetSymbolWinsStats(s)

		arrPay, isok := paytables.MapPay[int(s)]
		if isok {
			for i := 0; i < len(arrPay); i++ {
				if arrPay[i] > 0 {
					cw, err := CalcSymbolWinsInReelsWithLine(paytables, rss, symbols, wilds, s, i+1, getMaxPayoutSymbol(paytables, symbols, i+1))
					if err != nil {
						goutils.Error("AnalyzeReelsWithLine:CalcSymbolWinsInReelsWithLine",
							goutils.Err(err))

						return nil, err
					}

					sws.WinsNum[i] = cw
					sws.Wins[i] = int64(arrPay[i]) * sws.WinsNum[i] * int64(lineNum)

					ssws.TotalWins += int64(arrPay[i]) * sws.WinsNum[i] * int64(lineNum)
				}
			}
		}
	}

	ssws.onBuildEnd()

	return ssws, nil
}

func NewSymbolsWinsStats(num int) *SymbolsWinsStats {
	return &SymbolsWinsStats{
		MapSymbols: make(map[SymbolType]*SymbolWinsStats),
		Num:        num,
	}
}
