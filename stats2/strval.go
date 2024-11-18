package stats2

import (
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsStrVal struct {
	TotalUsedTimes int64            `json:"totalUsedTimes"`
	MapUsedTimes   map[string]int64 `json:"mapUsedTimes"`
}

func (strVal *StatsStrVal) UseVal(val string) {
	strVal.TotalUsedTimes++
	strVal.MapUsedTimes[val]++
}

func (strVal *StatsStrVal) Merge(src *StatsStrVal) {
	strVal.TotalUsedTimes += src.TotalUsedTimes

	for k, v := range src.MapUsedTimes {
		strVal.MapUsedTimes[k] += v
	}
}

func (strVal *StatsStrVal) SaveSheet(f *excelize.File, sheet string, s2 *Stats) {
	strVal.saveSheet(f, sheet, 0, 0, s2)
}

func (strVal *StatsStrVal) saveSheet(f *excelize.File, sheet string, sx, sy int, s2 *Stats) {
	f.SetCellValue(sheet, goutils.Pos2Cell(sx+0, sy+0), "totalTimes")

	f.SetCellValue(sheet, goutils.Pos2Cell(sx+1, sy+0), strVal.TotalUsedTimes)

	lstvals := []string{}
	for k := range strVal.MapUsedTimes {
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
		v := strVal.MapUsedTimes[k]
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+3, sy+y), k)
		f.SetCellValue(sheet, goutils.Pos2Cell(sx+4, sy+y), v)

		if strVal.TotalUsedTimes > 0 {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), float64(v)/float64(strVal.TotalUsedTimes))
		} else {
			f.SetCellValue(sheet, goutils.Pos2Cell(sx+5, sy+y), 0)
		}

		y++
	}
}

func NewStatsStrVal() *StatsStrVal {
	return &StatsStrVal{
		MapUsedTimes: make(map[string]int64),
	}
}
