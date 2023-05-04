package lowcode

import (
	"sync/atomic"
	"time"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"github.com/zhs007/slotsgamecore7/mathtoolset"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
)

type StatsConfig struct {
	Name      string         `yaml:"name"`
	Component string         `yaml:"component"`
	Status    []string       `yaml:"status"` // component -> status
	Children  []*StatsConfig `yaml:"children"`
}

func NewStatsFeature(parent *sgc7stats.Feature, name string, onAnalyze sgc7stats.FuncAnalyzeFeature, width int, symbols []mathtoolset.SymbolType, isStatus bool) *sgc7stats.Feature {
	var feature *sgc7stats.Feature

	if parent != nil {
		feature = sgc7stats.NewFeature(name, sgc7stats.FeatureBasic, onAnalyze, parent)
	} else {
		feature = sgc7stats.NewFeature(name, sgc7stats.FeatureBasic, onAnalyze, nil)
	}

	if isStatus {
		feature.Status = sgc7stats.NewStatus()
	} else {
		feature.Reels = sgc7stats.NewReels(width, symbols)
		feature.Symbols = sgc7stats.NewSymbolsRTP(width, symbols)
		feature.AllWins = sgc7stats.NewWins()
		feature.CurWins = sgc7stats.NewWins()
	}

	return feature
}

type StatsParam struct {
	Stake   *sgc7game.Stake
	Results []*sgc7game.PlayResult
}

type Stats struct {
	Root      *sgc7stats.Feature
	chanStats chan *StatsParam
	lastNum   int32
	TotalNum  int64
}

func (stats *Stats) StartWorker() {
	for {
		param := <-stats.chanStats

		stats.Root.OnResults(param.Stake, param.Results)

		atomic.AddInt32(&stats.lastNum, -1)
	}
}

func (stats *Stats) Push(stake *sgc7game.Stake, results []*sgc7game.PlayResult) {
	param := &StatsParam{
		Stake:   stake,
		Results: results,
	}

	atomic.AddInt32(&stats.lastNum, 1)
	atomic.AddInt64(&stats.TotalNum, 1)

	stats.chanStats <- param
}

func (stats *Stats) Wait() {
	for {
		v := atomic.LoadInt32(&stats.lastNum)
		if v > 0 {
			time.Sleep(time.Second)
		} else {
			break
		}
	}

}

func NewStats(root *sgc7stats.Feature) *Stats {
	stats := &Stats{
		Root:      root,
		chanStats: make(chan *StatsParam, 1024),
	}

	return stats
}
