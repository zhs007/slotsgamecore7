package mathtoolset

func buildWaysSymbolStatsKey(symbol SymbolType, numInWindow int) int {
	return int(symbol)*100 + numInWindow
}

type WaysSymbolStats struct {
	Symbol      SymbolType
	NumInWindow int
	Num         int
}

func newWaysSymbolStats(s SymbolType, numInWindow int, num int) *WaysSymbolStats {
	return &WaysSymbolStats{
		Symbol:      s,
		NumInWindow: numInWindow,
		Num:         num,
	}
}

type WaysReelStats struct {
	MapSymbols     map[int]*WaysSymbolStats
	TotalSymbolNum int
}
