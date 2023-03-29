package sgc7rtp

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7stats "github.com/zhs007/slotsgamecore7/stats"
	"go.uber.org/zap"
)

type FuncOnRTPResults func(lst []*sgc7game.PlayResult, gameData any)

type RTPReturnData struct {
	Return     float64
	TotalTimes int64
	Times      int64
	Total      float64
}

// RTP -
type RTP struct {
	WinNums             int64
	BetNums             int64
	TotalBet            int64
	TotalWins           int64
	Root                *RTPNode
	MapHR               map[string]*HitRateNode
	MapFeature          map[string]*FeatureNode
	Variance            float64
	StdDev              float64
	Returns             []float64
	ReturnWeights       []float64
	MaxReturn           int64
	MaxReturnNums       int64
	MapPlayerPool       map[string]*PlayerPoolData
	MapHitFrequencyData map[string]*HitFrequencyData
	MapReturn           map[string]*RTPReturnDataList
	MapStats            map[string]*RTPStats
	MaxCoincidingWin    float64
	Stats2              *sgc7stats.Feature
	FuncRTPResults      FuncOnRTPResults
}

// NewRTP - new RTP
func NewRTP() *RTP {
	return &RTP{
		Root:                NewRTPRoot(),
		MapHR:               make(map[string]*HitRateNode),
		MapFeature:          make(map[string]*FeatureNode),
		MapPlayerPool:       make(map[string]*PlayerPoolData),
		MapHitFrequencyData: make(map[string]*HitFrequencyData),
		MapReturn:           make(map[string]*RTPReturnDataList),
		MapStats:            make(map[string]*RTPStats),
	}
}

// Clone - clone
func (rtp *RTP) Clone() *RTP {
	nrtp := &RTP{
		BetNums:             rtp.BetNums,
		TotalBet:            rtp.TotalBet,
		TotalWins:           rtp.TotalWins,
		Root:                rtp.Root.Clone(),
		MapHR:               make(map[string]*HitRateNode),
		MapFeature:          make(map[string]*FeatureNode),
		MapPlayerPool:       make(map[string]*PlayerPoolData),
		MapHitFrequencyData: make(map[string]*HitFrequencyData),
		MapReturn:           make(map[string]*RTPReturnDataList),
		MapStats:            make(map[string]*RTPStats),
		MaxCoincidingWin:    rtp.MaxCoincidingWin,
		FuncRTPResults:      rtp.FuncRTPResults,
	}

	for k, v := range rtp.MapHR {
		nrtp.MapHR[k] = v.Clone()
	}

	for k, v := range rtp.MapFeature {
		nrtp.MapFeature[k] = v.Clone()
	}

	for k, v := range rtp.MapPlayerPool {
		nrtp.MapPlayerPool[k] = v.Clone()
	}

	for k, v := range rtp.MapHitFrequencyData {
		nrtp.MapHitFrequencyData[k] = v.Clone()
	}

	for k, v := range rtp.MapReturn {
		nrtp.MapReturn[k] = v.Clone()
	}

	for k, v := range rtp.MapStats {
		nrtp.MapStats[k] = v.Clone()
	}

	if rtp.Stats2 != nil {
		nrtp.Stats2 = rtp.Stats2.CloneIncludeChildren()
	}

	return nrtp
}

// Add - add
func (rtp *RTP) Add(rtp1 *RTP) {
	rtp.WinNums += rtp1.WinNums
	rtp.BetNums += rtp1.BetNums
	rtp.TotalBet += rtp1.TotalBet
	rtp.TotalWins += rtp1.TotalWins
	rtp.Returns = append(rtp.Returns, rtp1.Returns...)
	rtp.ReturnWeights = append(rtp.ReturnWeights, rtp1.ReturnWeights...)

	if rtp1.MaxReturn > rtp.MaxReturn {
		rtp.MaxReturn = rtp1.MaxReturn
		rtp.MaxReturnNums = rtp1.MaxReturnNums
	} else if rtp1.MaxReturn == rtp.MaxReturn {
		rtp.MaxReturnNums += rtp1.MaxReturnNums
	}

	if rtp1.MaxCoincidingWin > rtp.MaxCoincidingWin {
		rtp.MaxCoincidingWin = rtp1.MaxCoincidingWin
	}

	rtp.Root.Add(rtp1.Root)

	for k, v := range rtp.MapHR {
		v.Add(rtp1.MapHR[k])
	}

	for k, v := range rtp.MapFeature {
		v.Add(rtp1.MapFeature[k])
	}

	for k, v := range rtp.MapPlayerPool {
		v.Add(rtp1.MapPlayerPool[k])
	}

	for k, v := range rtp.MapHitFrequencyData {
		v.Add(rtp1.MapHitFrequencyData[k])
	}

	for k, v := range rtp.MapReturn {
		v.Merge(rtp1.MapReturn[k])
	}

	for k, v := range rtp.MapStats {
		v.Merge(rtp1.MapStats[k])
	}

	if rtp.Stats2 != nil && rtp1.Stats2 != nil {
		rtp.Stats2.Merge(rtp1.Stats2)
	}
}

// NewStats -
func (rtp *RTP) NewStats(tagname string) {
	rtp.MapStats[tagname] = &RTPStats{
		TagName: tagname,
	}
}

// OnStats -
func (rtp *RTP) OnStats(tagname string, val int64) {
	stats, isok := rtp.MapStats[tagname]
	if isok {
		stats.OnRTPStats(val)
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

	for _, v := range rtp.MapHitFrequencyData {
		v.OnHitFrequencyBet(v, bet)
	}
}

// OnResult -
func (rtp *RTP) OnResult(stake *sgc7game.Stake, pr *sgc7game.PlayResult, gameData any) {
	rtp.Root.OnResult(pr, gameData)

	for _, v := range rtp.MapHR {
		v.FuncOnResult(rtp, v, pr)
	}

	for _, v := range rtp.MapHitFrequencyData {
		v.OnHitFrequencyResult(v, pr)
	}

	curwin := float64(pr.CashWin) / float64(stake.CashBet)
	if curwin > rtp.MaxCoincidingWin {
		rtp.MaxCoincidingWin = curwin
	}
}

// OnResults -
func (rtp *RTP) OnResults(lst []*sgc7game.PlayResult, gameData any) {
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
		v.FuncOnResults(v, lst, gameData)
	}

	for _, v := range rtp.MapReturn {
		v.onResults(v, lst)
	}

	if rtp.FuncRTPResults != nil {
		rtp.FuncRTPResults(lst, gameData)
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

// AddReturnNode -
func (rtp *RTP) AddReturnNode(tag string, valRange []float64, funcOnResults FuncRDLOnResults) {
	rtp.MapReturn[tag] = NewRTPReturnDataList(tag, valRange, funcOnResults)
}

// Save2CSV -
func (rtp *RTP) Save2CSV(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		goutils.Error("sgc7rtp.RTP.Save2CSV",
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
		strhead = goutils.AppendString(strhead, ",X", strconv.Itoa(v))
	}
	strhead = goutils.AppendString(strhead, ",totalwin\n")

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

	if len(rtp.MapPlayerPool) > 0 {
		f.WriteString("\n\n\n")

		f.WriteString("name,playernums,values\n")

		keys := []string{}
		for k := range rtp.MapPlayerPool {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, v := range keys {
			str := rtp.MapPlayerPool[v].GenString()
			f.WriteString(str)
		}
	}

	if len(rtp.MapHitFrequencyData) > 0 {
		f.WriteString("\n\n\n")

		f.WriteString("name,total,trigger times\n")

		keys := []string{}
		for k := range rtp.MapHitFrequencyData {
			keys = append(keys, k)
		}

		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})

		for _, v := range keys {
			str := rtp.MapHitFrequencyData[v].GenString()
			f.WriteString(str)
		}
	}

	if len(rtp.MapStats) > 0 {
		f.WriteString("\n\n\n")

		f.WriteString("name,total,min,max,avg,times\n")

		for _, v := range rtp.MapStats {
			str := v.GenString()
			f.WriteString(str)
		}
	}

	f.WriteString("\n\n\n")
	f.WriteString("totalnums,winnums,Hit Frequency,Variance,StdDev,MaxReturn,MaxReturnNums,MaxCoincidingWin\n")
	str = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n",
		rtp.BetNums, rtp.WinNums, float64(rtp.WinNums)/float64(rtp.BetNums), rtp.Variance, rtp.StdDev, rtp.MaxReturn, rtp.MaxReturnNums, rtp.MaxCoincidingWin)
	f.WriteString(str)

	f.Sync()

	return nil
}

// SaveReturns2CSV -
func (rtp *RTP) SaveReturns2CSV(fn string) error {
	results := []*RTPReturnData{}
	totaltimes := int64(0)
	for i, v := range rtp.Returns {
		results = addResults(results, v, int64(rtp.ReturnWeights[i]))

		totaltimes += int64(rtp.ReturnWeights[i])
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Return < results[j].Return
	})

	f, err := os.Create(fn)
	if err != nil {
		goutils.Error("sgc7rtp.RTP.SaveReturns2CSV",
			zap.Error(err))

		return err
	}
	defer f.Close()

	f.WriteString("returns,totaltimes,times,per\n")
	for _, v := range results {
		str := fmt.Sprintf("%v,%v,%v,%v\n",
			v.Return, totaltimes, v.Times, float64(v.Times)/float64(totaltimes))
		f.WriteString(str)
	}

	f.Sync()

	return nil
}

// AddReturns -
func (rtp *RTP) AddReturns(ret float64) {
	if rtp.MaxReturn < int64(ret) {
		rtp.MaxReturn = int64(ret)
		rtp.MaxReturnNums = 1
	} else if rtp.MaxReturn == int64(ret) {
		rtp.MaxReturnNums++
	}

	for i, v := range rtp.Returns {
		if goutils.IsFloatEquals(v, ret) {
			rtp.ReturnWeights[i]++

			return
		}
	}

	rtp.Returns = append(rtp.Returns, ret)
	rtp.ReturnWeights = append(rtp.ReturnWeights, 1)
}

// AddPlayerPoolData -
func (rtp *RTP) AddPlayerPoolData(tag string, funcOnPlayer FuncOnPlayer) {
	rtp.MapPlayerPool[tag] = NewPlayerPoolData(tag, funcOnPlayer)
}

// OnPlayerPoolData -
func (rtp *RTP) OnPlayerPoolData(ps sgc7game.IPlayerState) {
	for _, v := range rtp.MapPlayerPool {
		v.OnPlayer(v, ps)
	}
}

// AddHitFrequencyData -
func (rtp *RTP) AddHitFrequencyData(tag string, onHitFrequencyBet FuncOnHitFrequencyBet, onHitFrequencyResult FuncOnHitFrequencyResult) {
	rtp.MapHitFrequencyData[tag] = NewHitFrequencyData(tag, onHitFrequencyBet, onHitFrequencyResult)
}

// SaveAllReturns -
func (rtp *RTP) SaveAllReturns(dir string, fnprefix string) {
	for k, v := range rtp.MapReturn {
		fn := path.Join(dir, fmt.Sprintf("fnprefix_%v.csv", k))
		v.SaveReturns2CSV(fn)
	}
}
