package stats

import (
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

type FuncAnalyzeFeature func(*Feature, *sgc7game.Stake, []*sgc7game.PlayResult) (bool, int64, int64)

type Feature struct {
	Name         string
	PlayTimes    int64
	TotalBets    int64
	TotalWins    int64
	TriggerTimes int64
	OnAnalyze    FuncAnalyzeFeature
	Obj          interface{}
}

func (feature *Feature) OnResults(stake *sgc7game.Stake, lst []*sgc7game.PlayResult) {
	feature.PlayTimes++

	istrigger, bet, wins := feature.OnAnalyze(feature, stake, lst)
	if istrigger {
		feature.TriggerTimes++

		feature.TotalWins += wins
	}

	feature.TotalBets += bet
}

func NewFeature(name string, onanalyze FuncAnalyzeFeature) *Feature {
	return &Feature{
		Name:      name,
		OnAnalyze: onanalyze,
	}
}
