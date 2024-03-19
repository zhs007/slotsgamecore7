package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Feature struct {
	Parent      string
	RootTrigger *StatsRootTrigger // 只有respin和foreach才需要这个
	Trigger     *StatsTrigger     // 普通的trigger，如果在respin或foreach下面，则需要配合它们才能得到正确的统计
	Wins        *StatsWins        // wins
}

func (f2 *Feature) onStatsGame(bet int) {
	if f2.Trigger != nil {
		f2.Trigger.TotalTimes++
	}

	if f2.RootTrigger != nil {
		f2.RootTrigger.TotalTimes++
	}
}

func (f2 *Feature) procCacheStatsWins(win int64) {
	if f2.Wins != nil {
		f2.Wins.TotalWin += win
	}
}

func (f2 *Feature) procCacheStatsTrigger() {
	f2.Trigger.TriggerTimes++
}

func (f2 *Feature) procCacheStatsRootTrigger() {
	if f2.RootTrigger != nil {
		if f2.RootTrigger.TriggerTimes == 0 {
			f2.RootTrigger.TriggerTimes++
		}

		f2.RootTrigger.RunTimes++
	}
}

func (f2 *Feature) procCacheStatsRootTriggerWins(win int64) {
	if f2.RootTrigger != nil {
		f2.RootTrigger.TotalWins += win
	}
}

// func (f2 *Feature) OnWins(win int64) {
// 	if f2.Wins != nil {
// 		f2.Wins.TotalWin += win
// 	}
// }

// func (f2 *Feature) OnBet(bet int64) {
// 	if f2.Wins != nil {
// 		f2.Wins.TotalBet += bet
// 	}

// 	if f2.Trigger != nil {
// 		f2.Trigger.TotalTimes++
// 	}
// }

// func (f2 *Feature) OnStep() {
// 	if f2.StepTrigger != nil {
// 		f2.StepTrigger.TotalTimes++
// 	}
// }

// func (f2 *Feature) OnStepTrigger(isTrigger bool) {
// 	if isTrigger {
// 		if f2.StepTrigger != nil {
// 			f2.StepTrigger.TriggerTimes++
// 		}
// 	}
// }

// func (f2 *Feature) OnTrigger(td *triggerData) {
// 	// if isTrigger {
// 	if f2.Trigger != nil {
// 		f2.Trigger.TriggerTimes++
// 	}
// 	// }
// }

func (f2 *Feature) Clone() *Feature {
	target := &Feature{}

	if f2.Trigger != nil {
		target.Trigger = f2.Trigger.Clone()
	}

	if f2.RootTrigger != nil {
		target.RootTrigger = f2.RootTrigger.Clone()
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

	if f2.RootTrigger != nil && src.RootTrigger != nil {
		f2.RootTrigger.Merge(src.RootTrigger)
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

	if f2.RootTrigger != nil {
		sn := fmt.Sprintf("%v - root trigger", sheet)
		f.NewSheet(sn)

		f2.RootTrigger.SaveSheet(f, sn)
	}

	if f2.Wins != nil {
		sn := fmt.Sprintf("%v - wins", sheet)
		f.NewSheet(sn)

		f2.Wins.SaveSheet(f, sn)
	}
}

func NewFeature(parent string, opts Options) *Feature {
	f2 := &Feature{
		Parent:  parent,
		Trigger: &StatsTrigger{},
	}

	if opts.Has(OptWins) {
		f2.Wins = &StatsWins{}
	}

	if opts.Has(OptRootTrigger) {
		f2.RootTrigger = &StatsRootTrigger{}
	}

	return f2
}
