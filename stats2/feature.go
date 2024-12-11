package stats2

import (
	"fmt"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type Feature struct {
	Parent      string            `json:"parent"`
	RootTrigger *StatsRootTrigger `json:"rootTrigger"` // 只有respin和foreach才需要这个
	Trigger     *StatsTrigger     `json:"trigger"`     // 普通的trigger，如果在respin或foreach下面，则需要配合它们才能得到正确的统计
	Wins        *StatsWins        `json:"wins"`        // wins
	IntVal      *StatsIntVal      `json:"intVal"`      // intVal
	StrVal      *StatsStrVal      `json:"strVal"`      // strVal
}

func (f2 *Feature) check() {
	if f2.RootTrigger != nil {
		if f2.RootTrigger.CurWins > 0 {
			goutils.Error("Feature.check:f2.RootTrigger.CurWins")
		}
	}
}

func (f2 *Feature) procCacheStatsIntVal(val int) {
	if f2.IntVal != nil {
		f2.IntVal.UseVal(val)
	}
}

func (f2 *Feature) procCacheStatsStrVal(val string) {
	if f2.StrVal != nil {
		f2.StrVal.UseVal(val)
	}
}

func (f2 *Feature) procCacheStatsWins(win int64) {
	if f2.Wins != nil {
		f2.Wins.AddWin(win)
	}
}

func (f2 *Feature) procCacheStatsTrigger() {
	f2.Trigger.TriggerTimes++
}

func (f2 *Feature) procCacheStatsRespinTrigger(wins int64, isEnding bool) {
	if f2.RootTrigger != nil {
		if !f2.RootTrigger.IsStarted {
			f2.RootTrigger.TriggerTimes++
			f2.RootTrigger.IsStarted = true
		}

		f2.RootTrigger.RunTimes++

		f2.RootTrigger.TotalWins += wins
		f2.RootTrigger.CurWins += wins

		if isEnding {
			f2.RootTrigger.IsStarted = false

			f2.RootTrigger.Wins.AddWin(f2.RootTrigger.CurWins)
			f2.RootTrigger.CurWins = 0
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

	if f2.IntVal != nil && src.IntVal != nil {
		f2.IntVal.Merge(src.IntVal)
	}

	if f2.StrVal != nil && src.StrVal != nil {
		f2.StrVal.Merge(src.StrVal)
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

		f2.Wins.SaveSheet(f, sn, s2)
	}

	if f2.IntVal != nil {
		sn := fmt.Sprintf("%v - intVal", sheet)
		f.NewSheet(sn)

		f2.IntVal.SaveSheet(f, sn, s2)
	}

	if f2.StrVal != nil {
		sn := fmt.Sprintf("%v - strVal", sheet)
		f.NewSheet(sn)

		f2.StrVal.SaveSheet(f, sn, s2)
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

	if opts.Has(OptIntVal) {
		f2.IntVal = NewStatsIntVal()
	}

	if opts.Has(OptStrVal) {
		f2.StrVal = NewStatsStrVal()
	}

	return f2
}
