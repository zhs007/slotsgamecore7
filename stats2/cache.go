package stats2

type Cache struct {
	MapStats map[string]*Feature
	Bet      int
	TotalWin int64
}

func (s2 *Cache) GetFeature(name string) *Feature {
	return s2.MapStats[name]
}

func (s2 *Cache) HasFeature(name string) bool {
	_, isok := s2.MapStats[name]

	return isok
}

func (s2 *Cache) AddFeature(name string, feature *Feature) {
	s2.MapStats[name] = feature
}

func (s2 *Cache) ProcStatsWins(name string, win int64, isRealWin bool) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsWins(win)

		if isRealWin && f2.Parent != "" {
			p2, isok := s2.MapStats[f2.Parent]
			if isok {
				p2.procCacheStatsRootTriggerWins(win)
			}
		}
	}

	if isRealWin {
		s2.TotalWin += win
	}
}

func (s2 *Cache) ProcStatsTrigger(name string) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsTrigger()
	}
}

func (s2 *Cache) ProcStatsRootTrigger(name string, wins int64, isEnding bool) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsRootTrigger(wins, isEnding)
	}
}

func NewCache(bet int) *Cache {
	s2 := &Cache{
		MapStats: make(map[string]*Feature),
		Bet:      bet,
	}

	return s2
}
