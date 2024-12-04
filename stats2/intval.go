package stats2

import (
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsIntVal struct {
	TotalUsedTimes int64         `json:"totalUsedTimes"`
	MapUsedTimes   map[int]int64 `json:"mapUsedTimes"`
}

func (intVal *StatsIntVal) UseVal(val int) {
	intVal.TotalUsedTimes++
	intVal.MapUsedTimes[val]++
}

func (intVal *StatsIntVal) Merge(src *StatsIntVal) {
	intVal.TotalUsedTimes += src.TotalUsedTimes

	for k, v := range src.MapUsedTimes {
		intVal.MapUsedTimes[k] += v
	}
}

func (intVal *StatsIntVal) SaveSheet(f *excelize.File, sheet string, s2 *Stats) {
	intVal.saveSheet(f, sheet, 0, 0, s2)
}

func (intVal *StatsIntVal) saveSheet(f *excelize.File, sheet string, sx, sy int, _ *Stats) {
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+0, sy+0), "totalTimes")

	f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+0), intVal.TotalUsedTimes)

	lstvals := []int{}
	for k := range intVal.MapUsedTimes {
		lstvals = append(lstvals, k)
	}

	sort.Slice(lstvals, func(i, j int) bool {
		return lstvals[i] < lstvals[j]
	})

	f.SetCellValue(sheet, goutils.Pos2Cell(sx+3, sy+2), "val")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+4, sy+2), "times")
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+2), "percent")

	y := 3
	for _, k := range lstvals {
		v := intVal.MapUsedTimes[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+3, sy+y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+4, sy+y), v)

		if intVal.TotalUsedTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), float64(v)/float64(intVal.TotalUsedTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), 0)
		}

		y++
	}
}

func NewStatsIntVal() *StatsIntVal {
	return &StatsIntVal{
		MapUsedTimes: make(map[int]int64),
	}
}
