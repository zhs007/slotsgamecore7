package sgc7rtp

import (
	"strconv"

	goutils "github.com/zhs007/goutils"
)

type RTPStats struct {
	TagName  string
	MaxVal   int64
	MinVal   int64
	TotalVal int64
	Times    int64
	AvgVal   int64
}

func (rtpstats *RTPStats) OnRTPStats(val int64) {
	if rtpstats.Times == 0 {
		rtpstats.MaxVal = val
		rtpstats.MinVal = val

		rtpstats.TotalVal += val
		rtpstats.Times++

		return
	}

	if val > rtpstats.MaxVal {
		rtpstats.MaxVal = val
	}

	if val < rtpstats.MinVal {
		rtpstats.MinVal = val
	}

	rtpstats.TotalVal += val
	rtpstats.Times++
}

func (rtpstats *RTPStats) OnEnd() {
	if rtpstats.Times > 0 {
		rtpstats.AvgVal = rtpstats.TotalVal / rtpstats.Times
	}
}

// Clone - clone
func (rtpstats *RTPStats) Clone() *RTPStats {
	stats1 := &RTPStats{
		TagName:  rtpstats.TagName,
		MaxVal:   rtpstats.MaxVal,
		MinVal:   rtpstats.MinVal,
		TotalVal: rtpstats.TotalVal,
		Times:    rtpstats.Times,
		AvgVal:   rtpstats.AvgVal,
	}

	return stats1
}

// Add - add
func (rtpstats *RTPStats) Merge(stats1 *RTPStats) {
	if rtpstats.TagName == stats1.TagName {
		rtpstats.TotalVal += stats1.TotalVal
		rtpstats.Times += stats1.Times

		rtpstats.AvgVal = rtpstats.TotalVal / rtpstats.Times

		if stats1.MaxVal > rtpstats.MaxVal {
			rtpstats.MaxVal = stats1.MaxVal
		}

		if stats1.MinVal < rtpstats.MinVal {
			rtpstats.MinVal = stats1.MinVal
		}
	}
}

func (rtpstats *RTPStats) GenString() string {
	return goutils.AppendString(rtpstats.TagName, ",",
		strconv.FormatInt(rtpstats.TotalVal, 10), ",",
		strconv.FormatInt(rtpstats.MinVal, 10), ",",
		strconv.FormatInt(rtpstats.MaxVal, 10), ",",
		strconv.FormatInt(rtpstats.AvgVal, 10), ",",
		strconv.FormatInt(rtpstats.Times, 10), "\n")
}
