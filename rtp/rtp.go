package sgc7rtp

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// RTP -
type RTP struct {
	WinNums       int64
	BetNums       int64
	TotalBet      int64
	Root          *RTPNode
	MapHR         map[string]*HitRateNode
	MapFeature    map[string]*FeatureNode
	Variance      float64
	Returns       []float64
	ReturnWeights []float64
}

// NewRTP - new RTP
func NewRTP() *RTP {
	return &RTP{
		Root:       NewRTPRoot(),
		MapHR:      make(map[string]*HitRateNode),
		MapFeature: make(map[string]*FeatureNode),
	}
}

// Clone - clone
func (rtp *RTP) Clone() *RTP {
	nrtp := &RTP{
		BetNums:    rtp.BetNums,
		TotalBet:   rtp.TotalBet,
		Root:       rtp.Root.Clone(),
		MapHR:      make(map[string]*HitRateNode),
		MapFeature: make(map[string]*FeatureNode),
	}

	for k, v := range rtp.MapHR {
		nrtp.MapHR[k] = v.Clone()
	}

	for k, v := range rtp.MapFeature {
		nrtp.MapFeature[k] = v.Clone()
	}

	return nrtp
}

// Add - add
func (rtp *RTP) Add(rtp1 *RTP) {
	rtp.WinNums += rtp1.WinNums
	rtp.BetNums += rtp1.BetNums
	rtp.TotalBet += rtp1.TotalBet
	rtp.Returns = append(rtp.Returns, rtp1.Returns...)
	rtp.ReturnWeights = append(rtp.ReturnWeights, rtp1.ReturnWeights...)

	rtp.Root.Add(rtp1.Root)

	for k, v := range rtp.MapHR {
		v.Add(rtp1.MapHR[k])
	}

	for k, v := range rtp.MapFeature {
		v.Add(rtp1.MapFeature[k])
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

	for _, v := range rtp.MapFeature {
		v.BetNums++
	}
}

// OnResult -
func (rtp *RTP) OnResult(pr *sgc7game.PlayResult) {
	rtp.Root.OnResult(pr)

	for _, v := range rtp.MapHR {
		v.FuncOnResult(v, pr)
	}
}

// OnResults -
func (rtp *RTP) OnResults(lst []*sgc7game.PlayResult) {
	iswin := false

	for _, v := range lst {
		if v.CoinWin > 0 {
			iswin = true

			break
		}
	}

	if iswin {
		rtp.WinNums++
	}

	for _, v := range rtp.MapFeature {
		v.FuncOnResults(v, lst)
	}
}

// AddHitRateNode -
func (rtp *RTP) AddHitRateNode(tag string, funcOnResult FuncHROnResult) {
	rtp.MapHR[tag] = NewSpecialHitRate(tag, funcOnResult)
}

// AddFeature -
func (rtp *RTP) AddFeature(tag string, funcOnResults FuncFeatureOnResults) {
	rtp.MapFeature[tag] = NewFeatureNode(tag, funcOnResults)
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

		f.WriteString("name,betnums,triggernums,totalnums,hitrate,average\n")

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

	if len(rtp.MapFeature) > 0 {
		f.WriteString("\n\n\n")

		f.WriteString("name,betnums,triggernums\n")

		keys := []string{}
		for k := range rtp.MapFeature {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, v := range keys {
			str := rtp.MapFeature[v].GenString()
			f.WriteString(str)
		}
	}

	f.WriteString("\n\n\n")
	f.WriteString("totalnums,winnums,Hit Frequency,Variance\n")
	str = fmt.Sprintf("%v,%v,%v,%v\n", rtp.BetNums, rtp.WinNums, float64(rtp.WinNums)/float64(rtp.BetNums), rtp.Variance)
	f.WriteString(str)

	f.Sync()

	return nil
}

// AddReturns -
func (rtp *RTP) AddReturns(ret float64) {
	for i, v := range rtp.Returns {
		if sgc7utils.IsFloatEquals(v, ret) {
			rtp.ReturnWeights[i]++

			return
		}
	}

	rtp.Returns = append(rtp.Returns, ret)
	rtp.ReturnWeights = append(rtp.ReturnWeights, 1)
}
