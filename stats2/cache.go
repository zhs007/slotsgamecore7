package stats2

import "github.com/zhs007/goutils"

type Cache struct {
	MapStats  map[string]*Feature
	Bet       int
	TotalWin  int64
	RespinArr []string
}

func (s2 *Cache) RemoveRespin(name string) {
	for i, v := range s2.RespinArr {
		if v == name {
			s2.RespinArr = append(s2.RespinArr[0:i], s2.RespinArr[i+1:]...)

			return
		}
	}
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

func (s2 *Cache) ProcStatsOnEnding(win int64) {
	s2.TotalWin = win
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

func (s2 *Cache) ProcStatsRespinTrigger(name string, isRunning bool, wins int64, isEnding bool) {
	f2, isok := s2.MapStats[name]
	if isok {
		f2.procCacheStatsRespinTrigger(isRunning, wins, isEnding)
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
