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
		syms := nrs.BuildSymbols(excsym)
		if len(syms) <= 0 {
			goutils.Error("genReel:BuildSymbols",
				goutils.JSON("excludeSymbols", excsym),
				zap.Int("lastnum", nrs.TotalSymbolNum),
				zap.Error(ErrNoValidSymbols))

			return nil, ErrNoValidSymbols
		}

		ci, err := plugin.Random(context.Background(), len(syms))
		if err != nil {
			goutils.Error("GenReels:Random",
				zap.Error(err))

			return nil, err
		}

		s := syms[ci]

		nrs.RemoveSymbol(s, 1)

		reel = append(reel, int(s))

		if len(excsym) >= minoff {
			excsym = excsym[1:]
		}

		excsym = append(excsym, s)
	}

	return reel, nil
}

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
