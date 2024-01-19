package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Feature struct {
	Trigger     *StatsTrigger
	StepTrigger *StatsTrigger
	Wins        *StatsWins
}

func (f2 *Feature) OnWins(win int64) {
	if f2.Wins != nil {
		f2.Wins.TotalWin += win
	}
}

func (f2 *Feature) OnBet(bet int64) {
	if f2.Wins != nil {
		f2.Wins.TotalBet += bet
	}

	if f2.Trigger != nil {
		f2.Trigger.TotalTimes++
	}
}

func (f2 *Feature) OnStep() {
	if f2.StepTrigger != nil {
		f2.StepTrigger.TotalTimes++
	}
}

func (f2 *Feature) OnStepTrigger(isTrigger bool) {
	if isTrigger {
		if f2.StepTrigger != nil {
			f2.StepTrigger.TriggerTimes++
		}
	}
}

func (f2 *Feature) OnTrigger(isTrigger bool) {
	if isTrigger {
		if f2.Trigger != nil {
			f2.Trigger.TriggerTimes++
		}
	}
}

func (f2 *Feature) Clone() *Feature {
	target := &Feature{}

	if f2.Trigger != nil {
		target.Trigger = f2.Trigger.Clone()
	}

	if f2.StepTrigger != nil {
		target.StepTrigger = f2.StepTrigger.Clone()
	}

	if f2.Wins != nil {
		target.Wins = f2.Wins.Clone()
	}

	return target
}

func (f2 *Feature) Merge(src *Feature) {
	if f2.Trigger != nil && src.Trigger != nil {
		f2.Trigger.Merge(src.Trigger)
	}

	if f2.StepTrigger != nil && src.StepTrigger != nil {
		f2.StepTrigger.Merge(src.StepTrigger)
	}

	if f2.Wins != nil && src.Wins != nil {
		f2.Wins.Merge(src.Wins)
	}
}

func (f2 *Feature) SaveSheet(f *excelize.File, sheet string) {
	if f2.Trigger != nil {
		sn := fmt.Sprintf("%v - trigger", sheet)
		f.NewSheet(sn)

		f2.Trigger.SaveSheet(f, sn)
	}

	if f2.StepTrigger != nil {
		sn := fmt.Sprintf("%v - stepTrigger", sheet)
		f.NewSheet(sn)

		f2.StepTrigger.SaveSheet(f, sn)
	}

	if f2.Wins != nil {
		sn := fmt.Sprintf("%v - wins", sheet)
		f.NewSheet(sn)

		f2.Wins.SaveSheet(f, sn)
	}
}

func NewFeature(opts Options) *Feature {
	f2 := &Feature{
		Trigger: &StatsTrigger{},
	}

	if opts.Has(OptWins) {
		f2.Wins = &StatsWins{}
	}

	if opts.Has(OptStepTrigger) {
		f2.StepTrigger = &StatsTrigger{}
	}

	return f2
}
