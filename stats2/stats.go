package stats2

import (
	"time"

	"github.com/xuri/excelize/v2"
)

type Stats struct {
	MapStats       map[string]*Feature
	chanBet        chan int
	chanCache      chan *Cache
	BetTimes       int64
	BetEndingTimes int64
	Components     []string
}

func (s2 *Stats) PushBet(bet int) {
	s2.chanBet <- bet
}

func (s2 *Stats) PushCache(cache *Cache) {
	s2.chanCache <- cache
}

func (s2 *Stats) onCache(cache *Cache) {
	for k, v := range cache.MapStats {
		v.onStatsGame(cache)

		s, isok := s2.MapStats[k]
		if isok {
			s.Merge(v)
		}
	}
}

func (s2 *Stats) SaveExcel(fn string) error {
	f := excelize.NewFile()

	for _, cn := range s2.Components {
		f2, isok := s2.MapStats[cn]
		if isok {
			f2.SaveSheet(f, cn)
		}
	}

	return f.SaveAs(fn)
}

func (s2 *Stats) Start() {
	go func() {
		for {
			<-s2.chanBet

			s2.BetTimes++
		}
	}()

	go func() {
		for {
			cache := <-s2.chanCache

			s2.onCache(cache)

			s2.BetEndingTimes++
		}
	}()
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

func NewStats(components []string) *Stats {
	s2 := &Stats{
		MapStats:   make(map[string]*Feature),
		chanBet:    make(chan int, 1024),
		chanCache:  make(chan *Cache, 1024),
		Components: components,
	}

	return s2
}
