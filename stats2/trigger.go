package stats2

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsTrigger struct {
	TriggerTimes int64
}

func (trigger *StatsTrigger) Clone() *StatsTrigger {
	return &StatsTrigger{
		TriggerTimes: trigger.TriggerTimes,
	}
}

func (trigger *StatsTrigger) Merge(src *StatsTrigger) {
	trigger.TriggerTimes += src.TriggerTimes
}

func (trigger *StatsTrigger) SaveSheet(f *excelize.File, sheet string, parent string, s2 *Stats) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "root run times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "trigger times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "percent")

	totaltimes := s2.GetRunTimes(parent)

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), totaltimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), trigger.TriggerTimes)
	if totaltimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(trigger.TriggerTimes)/float64(totaltimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}
}
