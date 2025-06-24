package mathtoolset2

import (
	"io"
	"log/slog"
	"slices"

	"github.com/zhs007/goutils"
)

func genReelDeep(pool *SymbolsPool, rules []*ExRule, rd []string) ([]string, error) {
	lst := BuildCurSymbols(rd, rules, pool)
	if len(lst) <= 0 {
		return nil, ErrUnkonow
	}

	// 新增全局可行性剪枝，递归前判断所有规则
	for _, rule := range rules {
		if !rule.IsFeasible(pool, rd) {
			return nil, ErrUnkonow
		}
	}

retry:
	nrd := slices.Clone(rd)
	npool := pool.Clone()

	sd := randomSymbolData(lst)

	for range sd.Num {
		nrd = append(nrd, sd.Symbol)
	}

	lst = rmSymbolData(lst, sd)
	npool.Remove(sd.Symbol, sd.Num)

	if len(npool.Pool) <= 0 {
		return nrd, nil
	}

	ret, err := genReelDeep(npool, rules, nrd)
	if err != nil {
		if len(lst) <= 0 {
			return nil, ErrUnkonow
		}

		goto retry
	}

	return ret, nil
}

func genReel(rs2 *ReelStats2, rules []*ExRule) ([]string, error) {
	rd := []string{}

	pool, err := rs2.genSymbolsPool()
	if err != nil {
		goutils.Error("genStackReel:genSymbolsPool",
			goutils.Err(err))

		return nil, err
	}

	return genReelDeep(pool, rules, rd)
}

func GenReels(reader io.Reader, strExRule string) ([][]string, error) {
	rss2, err := LoadReelsStats2(reader)
	if err != nil {
		goutils.Error("GenReels:LoadReelsStats2",
			goutils.Err(err))

		return nil, err
	}

	rules, err := ParseExRules(strExRule)
	if err != nil {
		goutils.Error("GenReels:ParseExRules",
			slog.String("strExRule", strExRule),
			goutils.Err(err))

		return nil, err
	}

	rd := [][]string{}

	for i, r := range rss2.Reels {
		reel, err := genReel(r, rules)
		if err != nil {
			goutils.Error("GenReels:genReel",
				slog.Int("reelIndex", i),
				goutils.Err(err))

			return nil, err
		}

		rd = append(rd, reel)
	}

	return rd, nil
}
