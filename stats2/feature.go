package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Feature struct {
	Parent      string            `json:"parent"`
	RootTrigger *StatsRootTrigger `json:"rootTrigger"` // 只有respin和foreach才需要这个
	Trigger     *StatsTrigger     `json:"trigger"`     // 普通的trigger，如果在respin或foreach下面，则需要配合它们才能得到正确的统计
	Wins        *StatsWins        `json:"wins"`        // wins
}

// func (f2 *Feature) onStatsGame(cache *Cache) {
// 	if f2.Parent != "" {
// 		p := cache.GetFeature(f2.Parent)

// 		if p != nil && p.RootTrigger != nil {
// 			if f2.Trigger != nil {
// 				f2.Trigger.TotalTimes += p.RootTrigger.RunTimes
// 			}

// 			if f2.RootTrigger != nil {
// 				f2.RootTrigger.TotalTimes += p.RootTrigger.RunTimes
// 			}

// 			return
// 		}
// 	}

// 	if f2.Trigger != nil {
// 		f2.Trigger.TotalTimes++
// 	}

// 	if f2.RootTrigger != nil {
// 		f2.RootTrigger.TotalTimes++
// 	}
// }

func (f2 *Feature) procCacheStatsWins(win int64) {
	if f2.Wins != nil {
		f2.Wins.AddWin(win)
	}
}

func (f2 *Feature) procCacheStatsTrigger() {
	f2.Trigger.TriggerTimes++
}

func (f2 *Feature) procCacheStatsRootTrigger(wins int64, isEnding bool) {
	if f2.RootTrigger != nil {
		if !f2.RootTrigger.IsStarted {
			f2.RootTrigger.TriggerTimes++
			f2.RootTrigger.IsStarted = true
		}

		f2.RootTrigger.RunTimes++

		if isEnding {
			f2.RootTrigger.IsStarted = false
			f2.RootTrigger.TotalWins = wins
		}
	}
}

func (f2 *Feature) procCacheStatsForeachTrigger(runtimes int, win int64) {
	if f2.RootTrigger != nil {
		f2.RootTrigger.TriggerTimes++
		f2.RootTrigger.RunTimes += int64(runtimes)
		f2.RootTrigger.TotalWins += win
	}
}

// func (f2 *Feature) Clone() *Feature {
// 	target := &Feature{}

// 	if f2.Trigger != nil {
// 		target.Trigger = f2.Trigger.Clone()
// 	}

// 	if f2.RootTrigger != nil {
// 		target.RootTrigger = f2.RootTrigger.Clone()
// 	}

// 	if f2.Wins != nil {
// 		target.Wins = f2.Wins.Clone()
// 	}

// 	return target
// }

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

func (f2 *Feature) SaveSheet(f *excelize.File, sheet string, s2 *Stats) {
	if f2.Trigger != nil {
		sn := fmt.Sprintf("%v - trigger", sheet)
		f.NewSheet(sn)

		f2.Trigger.SaveSheet(f, sn, f2.Parent, s2)
	}

	if f2.RootTrigger != nil {
		sn := fmt.Sprintf("%v - root trigger", sheet)
		f.NewSheet(sn)

		f2.RootTrigger.SaveSheet(f, sn, f2.Parent, s2)
	}

	if f2.Wins != nil {
		sn := fmt.Sprintf("%v - wins", sheet)
		f.NewSheet(sn)

		f2.Wins.SaveSheet(f, sn, s2.TotalBet)
	}
}

func NewFeature(parent string, opts Options) *Feature {
	f2 := &Feature{
		Parent: parent,
	}

	if opts.Has(OptWins) {
		f2.Wins = NewStatsWins()
	}

	if opts.Has(OptRootTrigger) {
		f2.RootTrigger = NewStatsRootTrigger()
	} else {
		f2.Trigger = NewStatsTrigger()
	}

	return f2
}
