package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type StatsWins struct {
	TotalWin int64
	TotalBet int64
}

func (wins *StatsWins) Clone() *StatsWins {
	return &StatsWins{
		TotalWin: wins.TotalWin,
		TotalBet: wins.TotalBet,
	}
}

func (wins *StatsWins) Merge(src *StatsWins) {
	wins.TotalBet += src.TotalBet
	wins.TotalWin += src.TotalWin
}

func (wins *StatsWins) SaveSheet(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "win")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "bet")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "rtp")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), wins.TotalWin)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), wins.TotalBet)
	if wins.TotalBet > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(wins.TotalWin)/float64(wins.TotalBet))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}
}

type Stats struct {
	TriggerTimes int64
	TotalTimes   int64
	Wins         *StatsWins
}

func (s2 *Stats) OnWins(win int64) {
	if s2.Wins != nil {
		s2.Wins.TotalWin += win
	}
}

func (s2 *Stats) OnTrigger(isTrigger bool) {
	if isTrigger {
		s2.TriggerTimes++
	}
}

func (s2 *Stats) OnBet(bet int64) {
	s2.TotalTimes++

	if s2.Wins != nil {
		s2.Wins.TotalBet += bet
	}
}

func (s2 *Stats) Clone() *Stats {
	target := &Stats{
		TotalTimes:   s2.TotalTimes,
		TriggerTimes: s2.TriggerTimes,
	}

	if s2.Wins != nil {
		target.Wins = s2.Wins.Clone()
	}

	return target
}

func (s2 *Stats) Merge(src *Stats) {
	s2.TotalTimes += src.TotalTimes
	s2.TriggerTimes += src.TriggerTimes

	if s2.Wins != nil && src.Wins != nil {
		s2.Wins.Merge(src.Wins)
	}
}

func (s2 *Stats) SaveSheet(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 0), "spin times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 1), "trigger times")
	f.SetCellValue(sheet, goutils.Pos2Cell(0, 2), "percent")

	f.SetCellValue(sheet, goutils.Pos2Cell(1, 0), s2.TotalTimes)
	f.SetCellValue(sheet, goutils.Pos2Cell(1, 1), s2.TriggerTimes)
	if s2.TotalTimes > 0 {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), float64(s2.TriggerTimes)/float64(s2.TotalTimes))
	} else {
		f.SetCellValue(sheet, goutils.Pos2Cell(1, 2), 0)
	}

	if s2.Wins != nil {
		sn := fmt.Sprintf("%v - wins", sheet)
		f.NewSheet(sn)

		s2.Wins.SaveSheet(f, sn)
	}
}

func NewStats(opts Options) *Stats {
	s2 := &Stats{}

	if opts.Has(OptWins) {
		s2.Wins = &StatsWins{}
	}

	return s2
}
