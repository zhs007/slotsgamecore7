package lowcode

import (
	"time"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/slotsgamecore7/stats2"
)

type stepTriggerData struct {
	Name      string
	IsTrigger bool
}

type Stats2 struct {
	MapStats        map[string]*stats2.Stats
	chanBet         chan int64
	chanStepTrigger chan *stepTriggerData
	chanBetEnding   chan int
	chanStep        chan int
	BetTimes        int64
	BetEndingTimes  int64
}

func (s2 *Stats2) onStatsBet(bet int64) {
	for _, v := range s2.MapStats {
		v.OnBet(bet)
	}
}

func (s2 *Stats2) pushBet(bet int64) {
	s2.chanBet <- bet
}

func (s2 *Stats2) pushStep() {
	s2.chanStep <- 0
}

func (s2 *Stats2) pushBetEnding() {
	s2.chanBetEnding <- 0
}

func (s2 *Stats2) onStatsStep() {
	for _, v := range s2.MapStats {
		v.OnStep()
	}
}

// func (s2 *Stats2) onStep() {
// 	for _, v := range s2.MapStats {
// 		v.OnStep()
// 	}
// }

// func (s2 *Stats2) onStats(componentName string, ic IComponent, icd IComponentData) {
// 	sd2 := s2.MapStats[componentName]
// 	if sd2 != nil {
// 		ic.OnStats2(icd, sd2)
// 	}
// }

func (s2 *Stats2) pushStepStats(componentName string, isTrigger bool) {
	s2.chanStepTrigger <- &stepTriggerData{
		Name:      componentName,
		IsTrigger: isTrigger,
	}
}

func (s2 *Stats2) onStepStats(ic IComponent, icd IComponentData) {
	ic.OnStats2(icd, s2)
}

func (s2 *Stats2) onStatsStepTrigger(std *stepTriggerData) {
	s2.MapStats[std.Name].OnStepTrigger(std.IsTrigger)
}

func (s2 *Stats2) SaveExcel(fn string) error {
	f := excelize.NewFile()

	for k, v := range s2.MapStats {
		v.SaveSheet(f, k)
	}

	return f.SaveAs(fn)
}

func (s2 *Stats2) Start() {
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
}

func (s2 *Stats2) WaitEnding() {
	for {
		if s2.BetTimes == s2.BetEndingTimes {
			return
		}

		time.Sleep(time.Second)
	}
}

func NewStats2(components *ComponentList) *Stats2 {
	s2 := &Stats2{
		MapStats:        make(map[string]*stats2.Stats),
		chanBet:         make(chan int64),
		chanStepTrigger: make(chan *stepTriggerData),
		chanBetEnding:   make(chan int),
		chanStep:        make(chan int),
	}

	for key, ic := range components.MapComponents {
		sd2 := ic.NewStats2()
		if sd2 != nil {
			s2.MapStats[key] = sd2
		}
	}

	return s2
}
