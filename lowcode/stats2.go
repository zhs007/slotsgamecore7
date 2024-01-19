package lowcode

import (
	"github.com/zhs007/slotsgamecore7/stats2"
)

// type triggerData struct {
// 	Name      string
// 	IsTrigger bool
// }

// type Stats2 struct {
// 	MapStats        map[string]*stats2.Feature
// 	chanBet         chan int64
// 	chanStepTrigger chan *triggerData
// 	chanBetEnding   chan int
// 	chanStep        chan int
// 	chanTrigger     chan *triggerData
// 	BetTimes        int64
// 	BetEndingTimes  int64
// }

// func (s2 *Stats2) onStatsBet(bet int64) {
// 	for _, v := range s2.MapStats {
// 		v.OnBet(bet)
// 	}
// }

// func (s2 *Stats2) pushBet(bet int64) {
// 	s2.chanBet <- bet
// }

// func (s2 *Stats2) pushStep() {
// 	s2.chanStep <- 0
// }

// func (s2 *Stats2) pushBetEnding() {
// 	s2.chanBetEnding <- 0
// }

// func (s2 *Stats2) onStatsStep() {
// 	for _, v := range s2.MapStats {
// 		v.OnStep()
// 	}
// }

// func (s2 *Stats2) pushStepTrigger(componentName string, isTrigger bool) {
// 	s2.chanStepTrigger <- &triggerData{
// 		Name:      componentName,
// 		IsTrigger: isTrigger,
// 	}
// }

// func (s2 *Stats2) pushTrigger(componentName string, isTrigger bool) {
// 	s2.chanTrigger <- &triggerData{
// 		Name:      componentName,
// 		IsTrigger: isTrigger,
// 	}
// }

// func (s2 *Stats2) onStepStats(ic IComponent, icd IComponentData) {
// 	ic.OnStats2(icd, s2)
// }

// func (s2 *Stats2) onStatsStepTrigger(std *triggerData) {
// 	s2.MapStats[std.Name].OnStepTrigger(std.IsTrigger)
// }

// func (s2 *Stats2) onStatsTrigger(std *triggerData) {
// 	s := s2.MapStats[std.Name]
// 	if s != nil {
// 		s.OnTrigger(std.IsTrigger)
// 	}
// }

// func (s2 *Stats2) SaveExcel(fn string) error {
// 	f := excelize.NewFile()

// 	for k, v := range s2.MapStats {
// 		v.SaveSheet(f, k)
// 	}

// 	return f.SaveAs(fn)
// }

// func (s2 *Stats2) Start() {
// 	go func() {
// 		for {
// 			bet := <-s2.chanBet

// 			s2.onStatsBet(bet)
// 			s2.BetTimes++
// 		}
// 	}()

// 	go func() {
// 		for {
// 			std := <-s2.chanStepTrigger

// 			s2.onStatsStepTrigger(std)
// 		}
// 	}()

// 	go func() {
// 		for {
// 			<-s2.chanBetEnding

// 			s2.BetEndingTimes++
// 		}
// 	}()

// 	go func() {
// 		for {
// 			<-s2.chanStep

// 			s2.onStatsStep()
// 		}
// 	}()

// 	go func() {
// 		for {
// 			td := <-s2.chanTrigger

// 			s2.onStatsTrigger(td)
// 		}
// 	}()
// }

// func (s2 *Stats2) WaitEnding() {
// 	for {
// 		if s2.BetTimes == s2.BetEndingTimes {
// 			return
// 		}

// 		time.Sleep(time.Second)
// 	}
// }

func NewStats2(components *ComponentList) *stats2.Stats {
	s2 := stats2.NewStats()

	for key, ic := range components.MapComponents {
		sd2 := ic.NewStats2()
		if sd2 != nil {
			s2.AddFeature(key, sd2)
		}
	}

	return s2
}
