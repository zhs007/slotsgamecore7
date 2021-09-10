package sgc7rtp

import (
	"strconv"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

// FuncOnHitFrequencyBet - onHitFrequencyBet(*HitFrequencyData, int64)
type FuncOnHitFrequencyBet func(hfd *HitFrequencyData, bet int64)

// FuncOnHitFrequencyResult - onHitFrequencyResult(*HitFrequencyData, *sgc7game.PlayResult)
type FuncOnHitFrequencyResult func(hfd *HitFrequencyData, pr *sgc7game.PlayResult)

// HitFrequencyData -
type HitFrequencyData struct {
	TagName              string
	Total                int64
	TriggerTimes         int64
	OnHitFrequencyBet    FuncOnHitFrequencyBet
	OnHitFrequencyResult FuncOnHitFrequencyResult
}

// NewHitFrequencyData - new HitFrequencyData
func NewHitFrequencyData(tag string, onHitFrequencyBet FuncOnHitFrequencyBet, onHitFrequencyResult FuncOnHitFrequencyResult) *HitFrequencyData {
	return &HitFrequencyData{
		TagName:              tag,
		OnHitFrequencyBet:    onHitFrequencyBet,
		OnHitFrequencyResult: onHitFrequencyResult,
	}
}

// Clone - clone
func (hfd *HitFrequencyData) Clone() *HitFrequencyData {
	hfd1 := &HitFrequencyData{
		TagName:              hfd.TagName,
		Total:                hfd.Total,
		TriggerTimes:         hfd.TriggerTimes,
		OnHitFrequencyBet:    hfd.OnHitFrequencyBet,
		OnHitFrequencyResult: hfd.OnHitFrequencyResult,
	}

	return hfd1
}

// Add - add
func (hfd *HitFrequencyData) Add(hfd1 *HitFrequencyData) {
	if hfd.TagName == hfd1.TagName {
		hfd.Total += hfd1.Total
		hfd.TriggerTimes += hfd1.TriggerTimes
	}
}

func (hfd *HitFrequencyData) GenString() string {
	return goutils.AppendString(hfd.TagName, ",", strconv.FormatInt(hfd.Total, 10), ",", strconv.FormatInt(hfd.TriggerTimes, 10), "\n")
}
