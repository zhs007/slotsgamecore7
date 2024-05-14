package mathtoolset2

import (
	"io"
	"log/slog"

	"github.com/zhs007/goutils"
)

// func genStackReelStrict(rs2 *ReelStats2, stack []int, excludeSymbol []string) ([]string, error) {
// 	ssds := []*stackSymbolData{}

// 	for _, s := range excludeSymbol {
// 		for i := 0; i < rs2.MapSymbols[s]; i++ {
// 			ssds = append(ssds, &stackSymbolData{
// 				symbol: s,
// 				num:    1,
// 			})
// 		}
// 	}

// 	for s, n := range rs2.MapSymbols {
// 		if n <= 0 {
// 			continue
// 		}

// 		if goutils.IndexOfStringSlice(excludeSymbol, s, 0) < 0 {
// 			nssds, isok := genStackSymbol(s, n, stack)
// 			if !isok {
// 				goutils.Error("genStackReelStrict:genStackReel",
// 					goutils.Err(ErrGenStackReel))

// 				return nil, ErrGenStackReel
// 			}

// 			ssds = append(ssds, nssds...)
// 		}
// 	}

// 	reel, isok := genReelSSD(ssds, nil, "")
// 	if isok {
// 		return reel, nil
// 	}

// 	goutils.Error("GenStackReels:genReelSSD",
// 		goutils.Err(ErrGenStackReel))

// 	return nil, ErrGenStackReel
// }

func genStackReelsStrict(reader io.Reader, stack []int, excludeSymbol []string) ([][]string, error) {
	rss2, err := LoadReelsStats2(reader)
	if err != nil {
		goutils.Error("genStackReelsStrick:LoadReelsStats2",
			goutils.Err(err))

		return nil, err
	}

	rd := [][]string{}

	for i, r := range rss2.Reels {
		reel, err := genStackReel(r, stack, excludeSymbol)
		if err != nil {
			goutils.Error("genStackReelsStrick:genStackReel",
				slog.Int("reelIndex", i),
				goutils.Err(err))

			return nil, err
		}

		rd = append(rd, reel)
	}

	return rd, nil
}
