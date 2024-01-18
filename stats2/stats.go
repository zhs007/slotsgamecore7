package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Stats struct {
	StepTrigger *StatsTrigger
	Wins        *StatsWins
}

func (s2 *Stats) OnWins(win int64) {
	if s2.Wins != nil {
		s2.Wins.TotalWin += win
	}
}

func (s2 *Stats) OnBet(bet int64) {
	if s2.Wins != nil {
		s2.Wins.TotalBet += bet
	}
}

func (s2 *Stats) OnStep() {
	if s2.StepTrigger != nil {
		s2.StepTrigger.TotalTimes++
	}
}

func (s2 *Stats) OnStepTrigger(isTrigger bool) {
	if s2.StepTrigger != nil && isTrigger {
		s2.StepTrigger.TriggerTimes++
	}
}

func (s2 *Stats) Clone() *Stats {
	target := &Stats{}

	if s2.StepTrigger != nil {
		target.StepTrigger = s2.StepTrigger.Clone()
	}

	if s2.Wins != nil {
		target.Wins = s2.Wins.Clone()
	}

	return target
}

func (s2 *Stats) Merge(src *Stats) {
	if s2.StepTrigger != nil && src.StepTrigger != nil {
		s2.StepTrigger.Merge(src.StepTrigger)
	}

	if s2.Wins != nil && src.Wins != nil {
		s2.Wins.Merge(src.Wins)
	}
}

func (s2 *Stats) SaveSheet(f *excelize.File, sheet string) {
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

	if opts.Has(OptStepTrigger) {
		s2.StepTrigger = &StatsTrigger{}
	}

	return s2
}
