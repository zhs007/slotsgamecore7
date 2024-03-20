package lowcode

import (
	"github.com/zhs007/slotsgamecore7/stats2"
)

func NewStats2(components *ComponentList) *stats2.Stats {
	lst := components.statsNodeData.GetComponents()
	s2 := stats2.NewStats(lst)

	for _, key := range lst {
		ic, isok := components.MapComponents[key]
		if isok {
			p := components.statsNodeData.GetParent(key)
			sd2 := ic.NewStats2(p)
			if sd2 != nil {
				s2.AddFeature(key, sd2)
			}
		}
	}

	return s2
}
