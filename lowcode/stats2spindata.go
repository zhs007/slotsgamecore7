package lowcode

type stats2SpinData struct {
	mapTrigger map[string]bool
}

func (ssd *stats2SpinData) onStepStats(ic IComponent, icd IComponentData) {
	ssd.mapTrigger[ic.GetName()] = true
}

func (ssd *stats2SpinData) onBetEnding(s2 *Stats2) {
	for k := range ssd.mapTrigger {
		s2.pushTrigger(k, true)
	}
}

func newStats2SpinData() *stats2SpinData {
	return &stats2SpinData{
		mapTrigger: make(map[string]bool),
	}
}
