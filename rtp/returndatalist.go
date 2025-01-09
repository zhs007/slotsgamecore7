package sgc7rtp

import (
	"fmt"
	"os"
	"sort"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"gonum.org/v1/gonum/stat"
)

// FuncRDLOnResults - onResult(*RTPReturnDataList, []*sgc7game.PlayResult)
type FuncRDLOnResults func(rdlst *RTPReturnDataList, lst []*sgc7game.PlayResult) bool

type RTPReturnDataList struct {
	Tag           string
	Returns       []int64
	ReturnWeights []int64
	MaxReturn     int64
	MaxReturnNums int64
	MaxReturnRNGs []int
	ValRange      []float64
	TotalReturns  []float64
	onResults     FuncRDLOnResults
}

func NewRTPReturnDataList(tag string, valRange []float64, onResults FuncRDLOnResults) *RTPReturnDataList {
	return &RTPReturnDataList{
		Tag:       tag,
		ValRange:  valRange,
		onResults: onResults,
	}
}

// AddReturns -
func (rdlst *RTPReturnDataList) AddReturns(fret float64, rngs []int) {
	iret := int64(fret * 100)

	if rdlst.MaxReturn < iret {
		rdlst.MaxReturn = iret
		rdlst.MaxReturnNums = 1
		rdlst.MaxReturnRNGs = rngs
	} else if rdlst.MaxReturn == iret {
		rdlst.MaxReturnNums++
	}

	for i, v := range rdlst.Returns {
		if v == iret {
			rdlst.ReturnWeights[i]++
			rdlst.TotalReturns[i] += fret

			return
		}
	}

	rdlst.Returns = append(rdlst.Returns, iret)
	rdlst.ReturnWeights = append(rdlst.ReturnWeights, 1)
	rdlst.TotalReturns = append(rdlst.TotalReturns, fret)
}

// AddReturnsEx -
func (rdlst *RTPReturnDataList) addReturnsEx(ret int64, times int64, rngs []int) {
	if rdlst.MaxReturn < ret {
		rdlst.MaxReturn = ret
		rdlst.MaxReturnNums = times

		if rngs != nil {
			rdlst.MaxReturnRNGs = rngs
		}
	} else if rdlst.MaxReturn == ret {
		rdlst.MaxReturnNums++
	}

	fret := float64(ret) / 100 * float64(times)

	for i, v := range rdlst.Returns {
		if v == ret {
			rdlst.ReturnWeights[i] += times
			rdlst.TotalReturns[i] += fret

			return
		}
	}

	rdlst.Returns = append(rdlst.Returns, ret)
	rdlst.ReturnWeights = append(rdlst.ReturnWeights, times)
	rdlst.TotalReturns = append(rdlst.TotalReturns, fret)
}

// Merge -
func (rdlst *RTPReturnDataList) Merge(lst *RTPReturnDataList) {
	if rdlst.MaxReturn < lst.MaxReturn {
		rdlst.MaxReturnRNGs = lst.MaxReturnRNGs
	}

	for i, v := range lst.Returns {
		rdlst.addReturnsEx(v, lst.ReturnWeights[i], nil)
	}
}

// Clone -
func (rdlst *RTPReturnDataList) Clone() *RTPReturnDataList {
	rdl := &RTPReturnDataList{
		MaxReturn:     rdlst.MaxReturn,
		MaxReturnNums: rdlst.MaxReturnNums,
		MaxReturnRNGs: rdlst.MaxReturnRNGs,
		ValRange:      rdlst.ValRange,
		onResults:     rdlst.onResults,
	}

	copy(rdl.Returns, rdlst.Returns)
	copy(rdl.ReturnWeights, rdlst.ReturnWeights)
	copy(rdl.TotalReturns, rdlst.TotalReturns)

	return rdl
}

// SaveReturns2CSV -
func (rdlst *RTPReturnDataList) SaveReturns2CSV(fn string) error {
	if len(rdlst.Returns) != len(rdlst.ReturnWeights) ||
		len(rdlst.Returns) != len(rdlst.TotalReturns) {
		// goutils.Error("RTPReturnDataList:SaveReturns2CSV",
		// 	slog.Int("Returns len", len(rdlst.Returns)),
		// 	slog.Int("ReturnWeights len", len(rdlst.ReturnWeights)),
		// 	slog.Int("TotalReturns len", len(rdlst.TotalReturns)),
		// 	goutils.Err(ErrInvalidReturnLen))

		return ErrInvalidReturnLen
	}

	results := []*RTPReturnData{}
	totaltimes := int64(0)
	for i, v := range rdlst.Returns {
		results = addResults3(results, v, rdlst.ReturnWeights[i], rdlst.TotalReturns[i])

		totaltimes += rdlst.ReturnWeights[i]
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Return < results[j].Return
	})

	f, err := os.Create(fn)
	if err != nil {
		// goutils.Error("sgc7rtp.RTPReturnDataList.SaveReturns2CSV",
		// 	goutils.Err(err))

		return err
	}
	defer f.Close()

	f.WriteString("Standard Deviation\n")
	f.WriteString(fmt.Sprintf("%v\n\n\n", rdlst.calcSD()))

	f.WriteString("returns,totaltimes,times,per,total\n")
	for _, v := range results {
		str := fmt.Sprintf("%v,%v,%v,%v,%v\n",
			v.Return, totaltimes, v.Times, float64(v.Times)/float64(totaltimes), v.Total)
		f.WriteString(str)
	}

	arr2 := rdlst.procValRange()
	f.WriteString("\n\n\n")
	f.WriteString("returns,totaltimes,times,per,total\n")
	for _, v := range arr2 {
		str := fmt.Sprintf("%v,%v,%v,%v,%v\n",
			v.Return, totaltimes, v.Times, float64(v.Times)/float64(totaltimes), v.Total)
		f.WriteString(str)
	}

	f.Sync()

	return nil
}

func (rdlst *RTPReturnDataList) calcSD() float64 {
	lstRets := []float64{}
	lstWeights := []float64{}

	for i, v := range rdlst.Returns {
		lstRets = append(lstRets, float64(v)/100)
		lstWeights = append(lstWeights, float64(rdlst.ReturnWeights[i]))
	}

	return stat.StdDev(lstRets, lstWeights)
}

func (rdlst *RTPReturnDataList) procValRange() []*RTPReturnData {
	results := []*RTPReturnData{}
	for i, v := range rdlst.Returns {
		vv := rdlst.countValRange(v)
		results = addResults3(results, vv, rdlst.ReturnWeights[i], rdlst.TotalReturns[i])
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Return < results[j].Return
	})

	return results
}

func (rdlst *RTPReturnDataList) countValRange(val int64) int64 {
	for _, v := range rdlst.ValRange {
		if float64(val)/100 < v {
			return int64(v * 100)
		}
	}

	return int64(rdlst.ValRange[len(rdlst.ValRange)-1] * 100)
}