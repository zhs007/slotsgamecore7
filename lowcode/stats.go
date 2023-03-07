package lowcode

import (
	"github.com/zhs007/slotsgamecore7/mathtoolset"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
)

func NewStats(parent *sgc7stats.Feature, name string, onAnalyze sgc7stats.FuncAnalyzeFeature, width int, symbols []mathtoolset.SymbolType) *sgc7stats.Feature {
	var feature *sgc7stats.Feature

	if parent != nil {
		feature = sgc7stats.NewFeature(name, sgc7stats.FeatureBasic, onAnalyze, parent)
	} else {
		feature = sgc7stats.NewFeature(name, sgc7stats.FeatureBasic, onAnalyze, nil)
	}

	feature.Reels = sgc7stats.NewReels(width, symbols)
	feature.Symbols = sgc7stats.NewSymbolsRTP(width, symbols)
	feature.AllWins = sgc7stats.NewWins()
	feature.CurWins = sgc7stats.NewWins()

	return feature
}
