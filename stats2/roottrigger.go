package stats2

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsRootTrigger struct {
	RunTimes     int64 `json:"runTimes"`
	TriggerTimes int64 `json:"triggerTimes"`
	TotalWins    int64 `json:"totalWins"`
	IsStarted    bool  `json:"-"`
}

func (trigger *StatsRootTrigger) Merge(src *StatsRootTrigger) {
	trigger.RunTimes += src.RunTimes
	trigger.TriggerTimes += src.TriggerTimes
	trigger.TotalWins += src.TotalWins
}

func (trigger *StatsRootTrigger) SaveSheet(f *excelize.File, sheet string, parent string, s2 *Stats) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "root run times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "trigger times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "percent")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 3), "total run times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 4), "avg run times for per trigger")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 5), "total wins")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 6), "avg wins for per trigger")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 7), "avg wins for per running")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 8), "rtp")

	totaltimes := s2.GetRunTimes(parent)

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), totaltimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), trigger.TriggerTimes)

	if totaltimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(trigger.TriggerTimes)/float64(totaltimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 3), trigger.RunTimes)

	if trigger.TriggerTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 4), float64(trigger.RunTimes)/float64(trigger.TriggerTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 4), 0)
	}

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 5), trigger.TotalWins)

	if trigger.TriggerTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 6), float64(trigger.TotalWins)/float64(trigger.TriggerTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 6), 0)
	}

	if trigger.RunTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 7), float64(trigger.TotalWins)/float64(trigger.RunTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 7), 0)
	}

	if s2.TotalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 8), float64(trigger.TotalWins)/float64(s2.TotalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 8), 0)
	}
}

func NewStatsRootTrigger() *StatsRootTrigger {
	return &StatsRootTrigger{}
}
