package lowcode

import (
	"github.com/zhs007/slotsgamecore7/stats2"
)

type Stats2 struct {
	MapStats map[string]*stats2.Stats
}

func (s2 *Stats2) onBet(bet int64) {
	for _, v := range s2.MapStats {
		v.OnBet(bet)
	}
}

func (s2 *Stats2) onStats(componentName string, ic IComponent, icd IComponentData) {
	sd2 := s2.MapStats[componentName]
	if sd2 != nil {

	}
}

func NewStats2(components *ComponentList) *Stats2 {
	s2 := &Stats2{
		MapStats: make(map[string]*stats2.Stats),
	}

	for key, ic := range components.MapComponents {
		sd2 := ic.NewStats2()
		if sd2 != nil {
			s2.MapStats[key] = sd2
		}
	}

	return s2
}
