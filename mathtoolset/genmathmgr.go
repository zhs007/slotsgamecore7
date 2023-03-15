package mathtoolset

import (
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type GenMathMgr struct {
	Paytables     *sgc7game.PayTables
	MapPaytables  map[string]*sgc7game.PayTables
	MapReelsStats map[string]*ReelsStats
	RTP           float32
	RSS           *ReelsStats
}

func (mgr *GenMathMgr) LoadPaytables(fn string) error {
	paytables, isok := mgr.MapPaytables[fn]
	if !isok {
		paytables1, err := sgc7game.LoadPaytablesFromExcel(fn)
		if err != nil {
			goutils.Error("GenMathMgr.LoadPaytables:LoadPaytablesFromExcel",
				zap.String("fn", fn),
				zap.Error(err))

			return err
		}

		mgr.MapPaytables[fn] = paytables1
		paytables = paytables1
	}

	mgr.Paytables = paytables

	return nil
}

func (mgr *GenMathMgr) LoadReelsState(fn string) error {
	rss, isok := mgr.MapReelsStats[fn]
	if !isok {
		rss1, err := LoadReelsStats(fn)
		if err != nil {
			goutils.Error("GenMathMgr.LoadReelsState:LoadReelsStats",
				zap.String("fn", fn),
				zap.Error(err))

			return err
		}

		mgr.MapReelsStats[fn] = rss1
		rss = rss1
	}

	mgr.RSS = rss

	return nil
}

func NewGamMathMgr() *GenMathMgr {
	return &GenMathMgr{
		MapPaytables:  make(map[string]*sgc7game.PayTables),
		MapReelsStats: make(map[string]*ReelsStats),
	}
}
