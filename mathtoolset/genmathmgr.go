package mathtoolset

import (
	"fmt"

	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

type GenMathMgr struct {
	Paytables     *sgc7game.PayTables
	MapPaytables  map[string]*sgc7game.PayTables
	MapReelsStats map[string]*ReelsStats
	MapReelsData  map[string]*sgc7game.ReelsData
	RTP           float32
	RSS           *ReelsStats
	RetStats      []*SymbolsWinsStats
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

func (mgr *GenMathMgr) LoadReelsData2(paytablesfn string, fn string) (*sgc7game.ReelsData, error) {
	mgr.LoadPaytables(paytablesfn)

	rd, isok := mgr.MapReelsData[fn]
	if !isok {
		paytables1, err := sgc7game.LoadReelsFromExcel2(fn, mgr.Paytables)
		if err != nil {
			goutils.Error("GenMathMgr.LoadReelsData2:LoadReelsFromExcel2",
				zap.String("fn", fn),
				zap.Error(err))

			return nil, err
		}

		mgr.MapReelsData[fn] = paytables1
		rd = paytables1
	}

	return rd, nil
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

func (mgr *GenMathMgr) Save() error {
	for i, v := range mgr.RetStats {
		v.SaveExcel(fmt.Sprintf("ssws-%v.xlsx", i), []SymbolsWinsFileMode{SWFModeRTP, SWFModeWins, SWFModeWinsNum})
	}

	return nil
}

func NewGamMathMgr() *GenMathMgr {
	return &GenMathMgr{
		MapPaytables:  make(map[string]*sgc7game.PayTables),
		MapReelsStats: make(map[string]*ReelsStats),
	}
}
