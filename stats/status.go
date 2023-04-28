package stats

import (
	"sort"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type statusdata struct {
	status int
	times  int64
}

type Status struct {
	MapStatus  map[int]int64
	TotalTimes int64
}

func (status *Status) Clone() *Status {
	nw := &Status{
		MapStatus:  make(map[int]int64),
		TotalTimes: status.TotalTimes,
	}

	for k, v := range status.MapStatus {
		nw.MapStatus[k] = v
	}

	return nw
}

func (status *Status) Merge(src *Status) {
	if src == nil {
		return
	}

	status.TotalTimes += src.TotalTimes

	for k, v := range src.MapStatus {
		status.MapStatus[k] += v
	}
}

func (status *Status) genData() []*statusdata {
	lst := []*statusdata{}

	for k, v := range status.MapStatus {
		lst = append(lst, &statusdata{
			status: k,
			times:  v,
		})
	}

	return lst
}

func (status *Status) AddStatus(win int) {
	status.TotalTimes++

	status.MapStatus[win]++
}

func (status *Status) SaveSheet(f *excelize.File, sheet string) error {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "status")
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), "times")

	lst := status.genData()

	sort.Slice(lst, func(i, j int) bool {
		return lst[i].status < lst[j].status
	})

	y := 1
	for _, v := range lst {
		f.SetCellValue(sheet, goutils.Pos2Cell(0, y), v.status)
		f.SetCellValue(sheet, goutils.Pos2Cell(1, y), v.times)

		y++
	}

	return nil
}

func NewStatus() *Status {
	return &Status{
		MapStatus: make(map[int]int64),
	}
}
