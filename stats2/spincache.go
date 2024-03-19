package stats2

// type SpinCache struct {
// 	mapTrigger map[string]bool
// }

// func (cache *SpinCache) OnStepTrigger(name string) {
// 	cache.mapTrigger[name] = true
// }

// func (cache *SpinCache) OnBetEnding(s2 *Stats) {
// 	for k := range cache.mapTrigger {
// 		s2.PushTrigger(k, true)
// 	}
// }

// func (cache *SpinCache) Clear() {
// 	maps.Clear(cache.mapTrigger)
// }

// func NewSpinCache() *SpinCache {
// 	return &SpinCache{
// 		mapTrigger: make(map[string]bool),
// 	}
// }
