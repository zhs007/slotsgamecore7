package sgc7rtp

import (
	"os"
	"sort"
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// RTP -
type RTP struct {
	BetNums  int64
	TotalBet int64
	Root     *RTPNode
	MapHR    map[string]*HitRateNode
}

// NewRTP - new RTP
func NewRTP() *RTP {
	return &RTP{
		Root:  NewRTPRoot(),
		MapHR: make(map[string]*HitRateNode),
	}
}

// CalcRTP -
func (rtp *RTP) CalcRTP() {
	rtp.Root.CalcRTP(rtp.TotalBet)
}

// Bet -
func (rtp *RTP) Bet(bet int64) {
	rtp.BetNums++
	rtp.TotalBet += bet

	for _, v := range rtp.MapHR {
		v.BetNums++
	}
}

// OnResult -
func (rtp *RTP) OnResult(pr *sgc7game.PlayResult) {
	rtp.Root.OnResult(pr)

	for _, v := range rtp.MapHR {
		v.funcOnResult(v, pr)
	}
}

// AddHitRateNode -
func (rtp *RTP) AddHitRateNode(tag string, funcOnResult FuncHROnResult) {
	rtp.MapHR[tag] = NewSpecialHitRate(tag, funcOnResult)
}

// Save2CSV -
func (rtp *RTP) Save2CSV(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		sgc7utils.Error("sgc7rtp.RTP.Save2CSV",
			zap.Error(err))

		return err
	}
	defer f.Close()

	gms := rtp.Root.GetGameMods(nil)
	sn := rtp.Root.GetSymbolNums(nil)
	symbols := rtp.Root.GetSymbols(nil)

	sort.Slice(sn, func(i, j int) bool {
		return sn[i] < sn[j]
	})

	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i] < symbols[j]
	})

	sort.Slice(gms, func(i, j int) bool {
		return gms[i] < gms[j]
	})

	strhead := "gamemod,tag,symbol,totalbet"
	for _, v := range sn {
		strhead = sgc7utils.AppendString(strhead, ",X", strconv.Itoa(v))
	}
	strhead = sgc7utils.AppendString(strhead, ",totalwin\n")

	f.WriteString(strhead)

	for _, gm := range gms {
		tags := rtp.Root.GetTags(nil, gm)
		if len(tags) == 0 {
			tags = append(tags, "")
		} else {
			sort.Slice(tags, func(i, j int) bool {
				return tags[i] < tags[j]
			})
		}

		for _, tag := range tags {
			for _, symbol := range symbols {
				str := rtp.Root.GenSymbolString(gm, tag, symbol, sn, rtp.TotalBet)
				f.WriteString(str)
			}

			str := rtp.Root.GenTagString(gm, tag, sn, rtp.TotalBet)
			f.WriteString(str)
		}

		str := rtp.Root.GenGameModString(gm, sn, rtp.TotalBet)
		f.WriteString(str)
	}

	str := rtp.Root.GenRootString(sn, rtp.TotalBet)
	f.WriteString(str)

	if len(rtp.MapHR) > 0 {
		f.WriteString("\n\n\n")

		f.WriteString("name,hitrate,average\n")

		keys := []string{}
		for k := range rtp.MapHR {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, v := range keys {
			str := rtp.MapHR[v].GenString()
			f.WriteString(str)
		}
	}

	f.Sync()

	return nil
}
