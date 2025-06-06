package stats2

import "github.com/zhs007/goutils"

type Cache struct {
	MapStats  map[string]*Feature
	Bet       int
	TotalWin  int64
	RespinArr []string
	rngs      []int
}

func (s2 *Cache) check() {
	for _, v := range s2.MapStats {
		v.check()
	}
}

func (s2 *Cache) OnStepEnd(respinArr []string) {
	newArr := []string{}
	for _, v := range s2.RespinArr {
		f, isok := s2.MapStats[v]
		if isok {
			if !(!f.RootTrigger.IsStarted && goutils.IndexOfStringSlice(respinArr, v, 0) < 0) {
				newArr = append(newArr, v)
			}
		}
	}

	s2.RespinArr = newArr
}

func (s2 *Cache) GetFeature(name string) *Feature {
	return s2.MapStats[name]
}

func (s2 *Cache) HasFeature(name string) bool {
	_, isok := s2.MapStats[name]

	return isok
}

func (s2 *Cache) AddFeature(name string, feature *Feature, isRespin bool) {
	s2.MapStats[name] = feature

	if isRespin {
		if goutils.IndexOfStringSlice(s2.RespinArr, name, 0) < 0 {
			s2.RespinArr = append(s2.RespinArr, name)
		}
	}
}

func (s2 *Cache) ProcStatsOnEnding(win int64, rngs []int) {
	s2.TotalWin = win
	s2.rngs = rngs
}

func (s2 *Cache) ProcStatsWins(name string, win int64) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsWins(win)
	}
}

func (s2 *Cache) ProcStatsTrigger(name string) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsTrigger()
	}
}

func (s2 *Cache) ProcStatsIntVal(name string, val int) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsIntVal(val)
	}
}

func (s2 *Cache) ProcStatsIntVal2(name string, val int) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsIntVal2(val)
	}
}

func (s2 *Cache) ProcStatsStrVal(name string, val string) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsStrVal(val)
	}
}

func (s2 *Cache) ProcStatsRespinTrigger(name string, wins int64, isEnding bool) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsRespinTrigger(wins, isEnding)
	}
}

func (s2 *Cache) ProcStatsForeachTrigger(name string, runtimes int, wins int64) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsForeachTrigger(runtimes, wins)
	}
}

func NewCache(bet int) *Cache {
	s2 := &Cache{
		MapStats: make(map[string]*Feature),
		Bet:      bet,
	}

	return s2
}
