package stats2

import (
	"log/slog"
	"time"

	"github.com/bytedance/sonic"
	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type Stats struct {
	MapStats       map[string]*Feature `json:"mapStats"`
	chanBet        chan int            `json:"-"`
	chanCache      chan *Cache         `json:"-"`
	TotalBet       int64               `json:"totalBet"`
	BetTimes       int64               `json:"betTimes"`
	BetEndingTimes int64               `json:"-"`
	Components     []string            `json:"components"`
}

func (s2 *Stats) PushBet(bet int) {
	s2.chanBet <- bet
}

func (s2 *Stats) PushCache(cache *Cache) {
	s2.chanCache <- cache
}

func (s2 *Stats) onCache(cache *Cache) {
	for k, v := range cache.MapStats {
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
			f2.SaveSheet(f, cn, s2)
		}
	}

	return f.SaveAs(fn)
}

func (s2 *Stats) ExportExcel() ([]byte, error) {
	f := excelize.NewFile()

	for _, cn := range s2.Components {
		f2, isok := s2.MapStats[cn]
		if isok {
			f2.SaveSheet(f, cn, s2)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		goutils.Error("Stats.ExportExcel:WriteToBuffer",
			goutils.Err(err))

		return nil, err
	}

	return buf.Bytes(), nil
}

func (s2 *Stats) Start() {
	go func() {
		for {
			bet := <-s2.chanBet

			s2.BetTimes++
			s2.TotalBet += int64(bet)
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

func (s2 *Stats) GetRunTimes(name string) int64 {
	if name == "" {
		return s2.BetTimes
	}

	f2, isok := s2.MapStats[name]
	if isok {
		return f2.RootTrigger.RunTimes
	}

	return 0
}

func (s2 *Stats) Merge(src *Stats) {
	for k, v := range src.MapStats {
		cv, isok := s2.MapStats[k]
		if isok {
			cv.Merge(v)
		}
	}
}

func (s2 *Stats) ToJson() string {
	str, _ := sonic.MarshalString(s2)

	return str
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

func LoadStats(str string) (*Stats, error) {
	s2 := &Stats{}

	err := sonic.UnmarshalString(str, s2)
	if err != nil {
		goutils.Error("LoadStats:UnmarshalString",
			slog.String("str", str),
			goutils.Err(err))

		return nil, err
	}

	return s2, nil
}
