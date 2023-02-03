package mathtoolset

import (
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

		if len(excsym) >= minoff {
			excsym = excsym[1:]
		}

		excsym = append(excsym, SymbolType(s))
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
func genReelsMainSymbolsDistance(plugin sgc7plugin.IPlugin, rs *ReelStats, mainSymbols []SymbolType, minoff int) ([]int, error) {
	nrs := rs.Clone()
	reel := []int{}
	excsym := []SymbolType{}

	msn := rs.CountSymbolsNum(mainSymbols)
	msoff := rs.TotalSymbolNum / msn
	curi := 0

	for nrs.TotalSymbolNum > 0 {
		var syms *sgc7game.ValWeights
		var err error

		if curi == msoff-1 {
			curi = 0

			syms, err = nrs.BuildSymbolsWithWeightsEx(mainSymbols)
			if err != nil || syms.MaxWeight <= 0 {
				syms, err = nrs.BuildSymbolsWithWeights2(excsym, mainSymbols)
				if err != nil || syms.MaxWeight <= 0 {
					goutils.Error("genReelsMainSymbolsDistance:BuildSymbols",
						goutils.JSON("excludeSymbols", excsym),
						zap.Int("lastnum", nrs.TotalSymbolNum),
						zap.Error(ErrNoValidSymbols))

					return nil, ErrNoValidSymbols
				}
			}
		} else {
			curi++

			syms, err = nrs.BuildSymbolsWithWeights2(excsym, mainSymbols)
			if err != nil || syms.MaxWeight <= 0 {
				goutils.Error("genReelsMainSymbolsDistance:BuildSymbols",
					goutils.JSON("excludeSymbols", excsym),
					zap.Int("lastnum", nrs.TotalSymbolNum),
					zap.Error(ErrNoValidSymbols))

				return nil, ErrNoValidSymbols
			}
		}

		s, err := syms.RandVal(plugin)
		if err != nil {
			goutils.Error("genReelsMainSymbolsDistance:Random",
				zap.Error(err))

			return nil, err
		}

		nrs.RemoveSymbol(SymbolType(s), 1)

		reel = append(reel, s)

		if len(excsym) >= minoff {
			excsym = excsym[1:]
		}

		excsym = append(excsym, SymbolType(s))
	}

	return reel, nil
}

// 随机，主要符号需要尽量分散开
func GenReelsMainSymbolsDistance(rss *ReelsStats, mainSymbols []SymbolType, minoff int, trytimes int) (*sgc7game.ReelsData, error) {
	reels := sgc7game.NewReelsData(len(rss.Reels))

	plugin := sgc7plugin.NewBasicPlugin()

	for i, rs := range rss.Reels {
		for j := 0; j < trytimes; j++ {
			reel, err := genReelsMainSymbolsDistance(plugin, rs, mainSymbols, minoff)
			if err != nil {
				goutils.Error("GenReelsMainSymbolsDistance:genReelsMainSymbolsDistance",
					zap.Error(err))
			} else {
				reels.SetReel(i, reel)

				break
			}
		}
	}

	return reels, nil
}
