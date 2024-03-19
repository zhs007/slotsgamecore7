package stats2

import (
	"time"

	"github.com/xuri/excelize/v2"
)

// type triggerData struct {
// 	Parent string
// 	Name   string
// }

type Stats struct {
	MapStats       map[string]*Feature
	chanBet        chan int
	chanCache      chan *Cache
	BetTimes       int64
	BetEndingTimes int64
}

// func (s2 *Stats) onStatsBet(bet int64) {
// 	for _, v := range s2.MapStats {
// 		v.OnBet(bet)
// 	}
// }

func (s2 *Stats) PushBet(bet int) {
	s2.chanBet <- bet
}

// func (s2 *Stats) PushCache() {
// 	s2.chanStep <- 0
// }

// func (s2 *Stats) PushBetEnding() {
// 	s2.chanBetEnding <- 0
// }

// func (s2 *Stats) onStatsStep() {
// 	for _, v := range s2.MapStats {
// 		v.OnStep()
// 	}
// }

// func (s2 *Stats) PushStepTrigger(componentName string, isTrigger bool) {
// 	s2.chanStepTrigger <- &triggerData{
// 		Name:      componentName,
// 		IsTrigger: isTrigger,
// 	}
// }

func (s2 *Stats) PushCache(cache *Cache) {
	s2.chanCache <- cache
}

// func (s2 *Stats) onStatsStepTrigger(std *triggerData) {
// 	s2.MapStats[std.Name].OnStepTrigger(std.IsTrigger)
// }

// func (s2 *Stats) onStatsTrigger(std *triggerData) {
// 	s := s2.MapStats[std.Name]
// 	if s != nil {
// 		s.OnTrigger(std)
// 	}
// }

func (s2 *Stats) onCache(cache *Cache) {
	for k, v := range cache.MapStats {
		v.onStatsGame(cache.Bet)

		s, isok := s2.MapStats[k]
		if isok {
			s.Merge(v)
		}
	}
}

func (s2 *Stats) SaveExcel(fn string) error {
	f := excelize.NewFile()

	for k, v := range s2.MapStats {
		v.SaveSheet(f, k)
	}

	return f.SaveAs(fn)
}

func (s2 *Stats) Start() {
	go func() {
		for {
			<-s2.chanBet

			// s2.onStatsBet(bet)
			s2.BetTimes++
		}
	}()

	// go func() {
	// 	for {
	// 		std := <-s2.chanStepTrigger

	// 		s2.onStatsStepTrigger(std)
	// 	}
	// }()

	go func() {
		for {
			cache := <-s2.chanCache

			s2.onCache(cache)

			s2.BetEndingTimes++
		}
	}()

	// go func() {
	// 	for {
	// 		<-s2.chanStep

	// 		s2.onStatsStep()
	// 	}
	// }()

	// go func() {
	// 	for {
	// 		td := <-s2.chanTrigger

	// 		s2.onStatsTrigger(td)
	// 	}
	// }()
}

func (s2 *Stats) WaitEnding() {
	for {
		if s2.BetTimes == s2.BetEndingTimes {
			return
		}

		time.Sleep(time.Second)
	}
}

func (s2 *Stats) AddFeature(name string, feature *Feature) {
	s2.MapStats[name] = feature
}

func NewStats() *Stats {
	s2 := &Stats{
		MapStats:  make(map[string]*Feature),
		chanBet:   make(chan int, 1024),
		chanCache: make(chan *Cache, 1024),
	}

	return s2
}
