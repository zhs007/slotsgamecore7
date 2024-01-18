package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Stats struct {
	Trigger     *StatsTrigger
	StepTrigger *StatsTrigger
	Wins        *StatsWins
}

func (s2 *Stats) OnWins(win int64) {
	if s2.Wins != nil {
		s2.Wins.TotalWin += win
	}
}

func (s2 *Stats) OnTrigger(isTrigger bool) {
	if isTrigger && s2.Trigger != nil {
		s2.Trigger.TriggerTimes++
	}
}

func (s2 *Stats) OnBet(bet int64) {
	s2.Trigger.TotalTimes++

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
	if isTrigger && s2.StepTrigger != nil {
		s2.StepTrigger.TriggerTimes++
	}
}

func (s2 *Stats) Clone() *Stats {
	target := &Stats{
		Trigger: s2.Trigger.Clone(),
	}

	if s2.Wins != nil {
		target.Wins = s2.Wins.Clone()
	}

	return target
}

func (s2 *Stats) Merge(src *Stats) {
	s2.Trigger.Merge(src.Trigger)

	if s2.Wins != nil && src.Wins != nil {
		s2.Wins.Merge(src.Wins)
	}
}

func (s2 *Stats) SaveSheet(f *excelize.File, sheet string) {
	s2.Trigger.SaveSheet(f, sheet)

	if s2.Wins != nil {
		sn := fmt.Sprintf("%v - wins", sheet)
		f.NewSheet(sn)

		s2.Wins.SaveSheet(f, sn)
	}
}

func NewStats(opts Options) *Stats {
	s2 := &Stats{
		Trigger: &StatsTrigger{},
	}

	if opts.Has(OptWins) {
		s2.Wins = &StatsWins{}
	}

	if opts.Has(OptStepTrigger) {
		s2.StepTrigger = &StatsTrigger{}
	}

	return s2
}
