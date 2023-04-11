package mathtoolset

// countScatterNumTimes -
func countScatterNumTimes(rss *ReelsStats, symbol SymbolType, num int, ci int, height int) int64 {
	if ci == len(rss.Reels)-1 {
		if num > 1 {
			return 0
		}

		if num > 0 {
			cn := rss.GetScatterNum(ci, symbol, IRSTypeSymbol, height)

			return int64(cn)
		}

		cnn := rss.GetScatterNum(ci, symbol, IRSTypeNoSymbol, height)

		return int64(cnn)
	}

	tn := int64(0)

	if num > 0 {
		cn := rss.GetScatterNum(ci, symbol, IRSTypeSymbol, height)
		if cn > 0 {
			nn := countScatterNumTimes(rss, symbol, num-1, ci+1, height)

			tn += int64(cn) * nn
		}
	}

	cnn := rss.GetScatterNum(ci, symbol, IRSTypeNoSymbol, height)
	if cnn > 0 {
		nn := countScatterNumTimes(rss, symbol, num, ci+1, height)

		tn += int64(cnn) * nn
	}

	return tn
}

func CalcScatterProbability(rss *ReelsStats, symbol SymbolType, num int, height int) float64 {

	totalnum := int64(1)

	for _, rs := range rss.Reels {
		totalnum *= int64(rs.TotalSymbolNum)
	}

	symbolnum := countScatterNumTimes(rss, symbol, num, 0, height)

	return float64(symbolnum) / float64(totalnum)
}
