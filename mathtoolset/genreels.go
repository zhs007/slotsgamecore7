package mathtoolset

import (
	"context"

	"github.com/zhs007/goutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"go.uber.org/zap"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func genReel(plugin sgc7plugin.IPlugin, rs *ReelStats, minoff int) ([]int, error) {
	nrs := rs.Clone()
	reel := []int{}
	excsym := []SymbolType{}

	for nrs.TotalSymbolNum > 0 {
		syms, err := nrs.BuildSymbolsWithWeights(excsym)
		if err != nil || syms.MaxWeight <= 0 {
			goutils.Error("genReel:BuildSymbolsWithWeights",
				goutils.JSON("excludeSymbols", excsym),
				zap.Int("lastnum", nrs.TotalSymbolNum),
				zap.Error(ErrNoValidSymbols))

			return nil, ErrNoValidSymbols
		}

		// ci, err := plugin.Random(context.Background(), len(syms))
		s, err := syms.RandVal(plugin)
		if err != nil {
			goutils.Error("GenReels:Random",
				zap.Error(err))

			return nil, err
		}

		// s := syms[ci]

		nrs.RemoveSymbol(SymbolType(s), 1)

		reel = append(reel, s)

		if minoff > 0 {
			if len(excsym) >= minoff {
				excsym = excsym[1:]
			}

			excsym = append(excsym, SymbolType(s))
		}
	}

	return reel, nil
}

// 纯随机，重复符号需要不重复
func GenReels(rss *ReelsStats, minoff int, trytimes int) (*sgc7game.ReelsData, error) {
	reels := sgc7game.NewReelsData(len(rss.Reels))

	plugin := sgc7plugin.NewBasicPlugin()

	for i, rs := range rss.Reels {
		for j := 0; j < trytimes; j++ {
			reel, err := genReel(plugin, rs, minoff)
			if err != nil {
				goutils.Error("GenReels:genReel",
					zap.Error(err))
			} else {
				reels.SetReel(i, reel)

				break
			}
		}
	}

	return reels, nil
}

// 随机，主要符号需要尽量分散开
func genReelsMainSymbolsDistance(plugin sgc7plugin.IPlugin, rs *ReelStats,
	mainSymbols []SymbolType, minoff int) ([]int, error) {

	// 应该用圆的切割算法
	// 1个点没办法将圆切成2段，2个点才能切成2段
	msn := rs.CountSymbolsNum(mainSymbols)
	if msn <= 1 {
		return genReel(plugin, rs, minoff)
	}

	nrs := rs.Clone()
	reel := []int{}
	excsym := []SymbolType{}
	var cacheexcsym []SymbolType

	// 这里需要注意的是，譬如35，切6份，余5
	// 如果把5留到最后，就会出现间隔11，所以需要把5分摊到中间去

	msoff := rs.TotalSymbolNum / msn
	lastoff := rs.TotalSymbolNum % msn
	curi := 0
	curoff := msoff

	for nrs.TotalSymbolNum > 0 {
		var syms *sgc7game.ValWeights
		var err error

		if curi == 0 {
			curoff = msoff

			if lastoff > 0 && msn > 0 {
				cr, err := plugin.Random(context.Background(), msn)
				if err != nil {
					goutils.Error("genReelsMainSymbolsDistance:Random",
						goutils.JSON("msn", msn),
						zap.Int("lastoff", lastoff),
						zap.Error(err))

					return nil, err
				}

				if cr < lastoff {
					lastoff--

					curoff++
				}

				msn--
			}

			syms, err = nrs.BuildSymbolsWithWeightsEx(mainSymbols)
			if err != nil || syms.MaxWeight <= 0 {
				syms, err = nrs.BuildSymbolsWithWeights2(excsym, mainSymbols)
				if err != nil || syms.MaxWeight <= 0 {
					// goutils.Error("genReelsMainSymbolsDistance:BuildSymbols",
					// 	goutils.JSON("excludeSymbols", excsym),
					// 	zap.Int("lastnum", nrs.TotalSymbolNum),
					// 	zap.Error(ErrNoValidSymbols))

					return nil, ErrNoValidSymbols
				}
			}
		} else {
			syms, err = nrs.BuildSymbolsWithWeights2(excsym, mainSymbols)
			if err != nil || syms.MaxWeight <= 0 {
				// goutils.Error("genReelsMainSymbolsDistance:BuildSymbols",
				// 	goutils.JSON("excludeSymbols", excsym),
				// 	zap.Int("lastnum", nrs.TotalSymbolNum),
				// 	zap.Error(ErrNoValidSymbols))

				return nil, ErrNoValidSymbols
			}
		}

		if curi == curoff-1 {
			curi = 0
		} else {
			curi++
		}

		s, err := syms.RandVal(plugin)
		if err != nil {
			goutils.Error("genReelsMainSymbolsDistance:Random",
				zap.Error(err))

			return nil, err
		}

		nrs.RemoveSymbol(SymbolType(s), 1)

		reel = append(reel, s)

		if cacheexcsym != nil {
			excsym = cacheexcsym

			cacheexcsym = nil
		}

		if len(excsym) >= minoff {
			excsym = excsym[1:]
		}

		excsym = append(excsym, SymbolType(s))

		if nrs.TotalSymbolNum >= minoff {
			cacheexcsym = CloneSymbols(excsym)

			for tt := 0; tt < minoff-nrs.TotalSymbolNum-1; tt++ {
				excsym = append(excsym, SymbolType(reel[tt]))
			}
		}
	}

	return reel, nil
}

// 随机，主要符号需要尽量分散开
func GenReelsMainSymbolsDistance(rss *ReelsStats, mainSymbols []SymbolType, minoff int, trytimes int) (*sgc7game.ReelsData, error) {
	reels := sgc7game.NewReelsData(len(rss.Reels))

	plugin := sgc7plugin.NewBasicPlugin()

	for i, rs := range rss.Reels {
		isok := false
		for j := 0; j < trytimes; j++ {
			reel, err := genReelsMainSymbolsDistance(plugin, rs, mainSymbols, minoff)
			if err != nil {
				// goutils.Error("GenReelsMainSymbolsDistance:genReelsMainSymbolsDistance",
				// 	zap.Error(err))
			} else {
				isok = true

				reels.SetReel(i, reel)

				break
			}
		}

		if !isok {
			goutils.Error("GenReelsMainSymbolsDistance:genReelsMainSymbolsDistance",
				zap.Error(ErrNoValidSymbols))

			return nil, ErrNoValidSymbols
		}
	}

	return reels, nil
}
