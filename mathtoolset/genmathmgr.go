package mathtoolset

import (
	"fmt"
	"log/slog"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type GenMathMgr struct {
	Paytables     *sgc7game.PayTables
	MapPaytables  map[string]*sgc7game.PayTables
	MapReelsStats map[string]*ReelsStats
	MapReelsData  map[string]*sgc7game.ReelsData
	RTP           float32
	RSS           *ReelsStats
	RetStats      []*SymbolsWinsStats
	Rets          []float64
	MapRets       map[string]float64
	Config        *Config
}

func (mgr *GenMathMgr) LoadPaytables(fn string) error {
	paytables, isok := mgr.MapPaytables[fn]
	if !isok {
		paytables1, err := sgc7game.LoadPaytablesFromExcel(fn)
		if err != nil {
			goutils.Error("GenMathMgr.LoadPaytables:LoadPaytablesFromExcel",
				slog.String("fn", fn),
				goutils.Err(err))

			return err
		}

		mgr.MapPaytables[fn] = paytables1
		paytables = paytables1
	}

	mgr.Paytables = paytables

	return nil
}

func (mgr *GenMathMgr) LoadReelsData(paytablesfn string, fn string, isStrReel bool) (*sgc7game.ReelsData, error) {
	mgr.LoadPaytables(paytablesfn)

	rd, isok := mgr.MapReelsData[fn]
	if !isok {
		if isStrReel {
			rd1, err := sgc7game.LoadReelsFromExcel2(fn, mgr.Paytables)
			if err != nil {
				goutils.Error("GenMathMgr.LoadReelsData:LoadReelsFromExcel2",
					slog.String("fn", fn),
					goutils.Err(err))

				return nil, err
			}

			rd = rd1
		} else {
			rd1, err := sgc7game.LoadReelsFromExcel(fn)
			if err != nil {
				goutils.Error("GenMathMgr.LoadReelsData:LoadReelsFromExcel",
					slog.String("fn", fn),
					goutils.Err(err))

				return nil, err
			}

			rd = rd1
		}

		mgr.MapReelsData[fn] = rd
	}

	return rd, nil
}

func (mgr *GenMathMgr) LoadReelsState(fn string) error {
	rss, isok := mgr.MapReelsStats[fn]
	if !isok {
		rss1, err := LoadReelsStats(fn)
		if err != nil {
			goutils.Error("GenMathMgr.LoadReelsState:LoadReelsStats",
				slog.String("fn", fn),
				goutils.Err(err))

			return err
		}

		mgr.MapReelsStats[fn] = rss1
		rss = rss1
	}

	mgr.RSS = rss

	return nil
}

func (mgr *GenMathMgr) Save() error {
	mgr.saveResults("genmath.xlsx")

	for i, v := range mgr.RetStats {
		v.SaveExcel(fmt.Sprintf("ssws-%v.xlsx", i), []SymbolsWinsFileMode{SWFModeRTP, SWFModeWins, SWFModeWinsNum})
	}

	return nil
}

func (mgr *GenMathMgr) saveResults(fn string) error {
	if mgr.Config == nil {
		return nil
	}

	f := excelize.NewFile()

	sheet := f.GetSheetList()[0]

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "name")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "retsult")

	si := 1

	for _, v := range mgr.Config.Codes {
		f.SetCellStr(sheet, goutils.Pos2Cell(0, si), fmt.Sprintf("%v", v.Name))
		f.SetCellStr(sheet, goutils.Pos2Cell(1, si), fmt.Sprintf("%v", mgr.MapRets[v.Name]))

		si++
	}

	return f.SaveAs(fn)
}

func (mgr *GenMathMgr) RunCode(i int) error {
	script, err := NewScriptCore(mgr)
	if err != nil {
		goutils.Error("GenMathMgr.RunCode:NewScriptCore",
			goutils.Err(err))

		return err
	}

	err = script.Compile(mgr.Config.Codes[i].Code)
	if err != nil {
		goutils.Error("GenMathMgr.RunCode:Compile",
			goutils.Err(err))

		return err
	}

	out, err := script.Eval(mgr)
	if err != nil {
		goutils.Error("GenMathMgr.RunCode:Eval",
			goutils.Err(err))

		return err
	}

	mgr.Rets = append(mgr.Rets, out.Value().(float64))

	fmt.Printf("RunCode %v is %v\n", i, out.Value().(float64))

	return nil
}

func (mgr *GenMathMgr) RunCodeEx(name string) (ref.Val, error) {
	for _, v := range mgr.Config.Codes {
		if v.Name == name {
			script, err := NewScriptCore(mgr)
			if err != nil {
				goutils.Error("GenMathMgr.RunCode:NewScriptCore",
					goutils.Err(err))

				return nil, err
			}

			err = script.Compile(v.Code)
			if err != nil {
				goutils.Error("GenMathMgr.RunCode:Compile",
					goutils.Err(err))

				return nil, err
			}

			out, err := script.Eval(mgr)
			if err != nil {
				goutils.Error("GenMathMgr.RunCode:Eval",
					goutils.Err(err))

				return nil, err
			}

			mgr.Rets = append(mgr.Rets, out.Value().(float64))
			mgr.MapRets[name] = out.Value().(float64)

			fmt.Printf("RunCode %v is %v\n", name, out.Value().(float64))

			return out, nil
		}
	}

	return types.Double(0), nil
}

func NewGamMathMgr(cfg *Config) *GenMathMgr {
	return &GenMathMgr{
		MapPaytables:  make(map[string]*sgc7game.PayTables),
		MapReelsStats: make(map[string]*ReelsStats),
		MapReelsData:  make(map[string]*sgc7game.ReelsData),
		MapRets:       make(map[string]float64),
		Config:        cfg,
	}
}
