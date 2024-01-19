package stats2

import (
	"time"

	"github.com/xuri/excelize/v2"
)

type triggerData struct {
	Name      string
	IsTrigger bool
}

type Stats struct {
	MapStats        map[string]*Feature
	chanBet         chan int64
	chanStepTrigger chan *triggerData
	chanBetEnding   chan int
	chanStep        chan int
	chanTrigger     chan *triggerData
	BetTimes        int64
	BetEndingTimes  int64
}

func (s2 *Stats) onStatsBet(bet int64) {
	for _, v := range s2.MapStats {
		v.OnBet(bet)
	}
}

func (s2 *Stats) PushBet(bet int64) {
	s2.chanBet <- bet
}

func (s2 *Stats) PushStep() {
	s2.chanStep <- 0
}

func (s2 *Stats) PushBetEnding() {
	s2.chanBetEnding <- 0
}

func (s2 *Stats) onStatsStep() {
	for _, v := range s2.MapStats {
		v.OnStep()
	}
}

func (s2 *Stats) PushStepTrigger(componentName string, isTrigger bool) {
	s2.chanStepTrigger <- &triggerData{
		Name:      componentName,
		IsTrigger: isTrigger,
	}
}

func (s2 *Stats) PushTrigger(componentName string, isTrigger bool) {
	s2.chanTrigger <- &triggerData{
		Name:      componentName,
		IsTrigger: isTrigger,
	}
}

func (s2 *Stats) onStatsStepTrigger(std *triggerData) {
	s2.MapStats[std.Name].OnStepTrigger(std.IsTrigger)
}

func (s2 *Stats) onStatsTrigger(std *triggerData) {
	s := s2.MapStats[std.Name]
	if s != nil {
		s.OnTrigger(std.IsTrigger)
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
			bet := <-s2.chanBet

			s2.onStatsBet(bet)
			s2.BetTimes++
		}
	}()

	go func() {
		for {
			std := <-s2.chanStepTrigger

			s2.onStatsStepTrigger(std)
		}
	}()

	go func() {
		for {
			<-s2.chanBetEnding

			s2.BetEndingTimes++
		}
	}()

	go func() {
		for {
			<-s2.chanStep

			s2.onStatsStep()
		}
	}()

	go func() {
		for {
			td := <-s2.chanTrigger

			s2.onStatsTrigger(td)
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

func NewStats() *Stats {
	s2 := &Stats{
		MapStats:        make(map[string]*Feature),
		chanBet:         make(chan int64),
		chanStepTrigger: make(chan *triggerData),
		chanBetEnding:   make(chan int),
		chanStep:        make(chan int),
		chanTrigger:     make(chan *triggerData),
	}

	return s2
}
