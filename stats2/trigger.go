package stats2

import (
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsTrigger struct {
	TriggerTimes int64
	TotalTimes   int64
}

func (trigger *StatsTrigger) Clone() *StatsTrigger {
	return &StatsTrigger{
		TotalTimes:   trigger.TotalTimes,
		TriggerTimes: trigger.TriggerTimes,
	}
}

func (trigger *StatsTrigger) Merge(src *StatsTrigger) {
	trigger.TotalTimes += src.TotalTimes
	trigger.TriggerTimes += src.TriggerTimes
}

func (trigger *StatsTrigger) SaveSheet(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "spin times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "trigger times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "percent")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), trigger.TotalTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), trigger.TriggerTimes)
	if trigger.TotalTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(trigger.TriggerTimes)/float64(trigger.TotalTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}
}
