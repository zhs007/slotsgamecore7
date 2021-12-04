package sgc7rtp

import (
	"fmt"
	"os"
	"sort"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	"go.uber.org/zap"
)

// FuncRDLOnResults - onResult(*RTPReturnDataList, []*sgc7game.PlayResult)
type FuncRDLOnResults func(rdlst *RTPReturnDataList, lst []*sgc7game.PlayResult) bool

type RTPReturnDataList struct {
	Tag           string
	Returns       []int64
	ReturnWeights []int64
	MaxReturn     int64
	MaxReturnNums int64
	onResults     FuncRDLOnResults
}

func NewRTPReturnDataList(tag string, onResults FuncRDLOnResults) *RTPReturnDataList {
	return &RTPReturnDataList{
		Tag:       tag,
		onResults: onResults,
	}
}

// AddReturns -
func (rdlst *RTPReturnDataList) AddReturns(fret float64) {
	iret := int64(fret * 100)

	if rdlst.MaxReturn < iret {
		rdlst.MaxReturn = iret
		rdlst.MaxReturnNums = 1
	} else if rdlst.MaxReturn == iret {
		rdlst.MaxReturnNums++
	}

	for i, v := range rdlst.Returns {
		if v == iret {
			rdlst.ReturnWeights[i]++

			return
		}
	}

	rdlst.Returns = append(rdlst.Returns, iret)
	rdlst.ReturnWeights = append(rdlst.ReturnWeights, 1)
}

// AddReturnsEx -
func (rdlst *RTPReturnDataList) addReturnsEx(ret int64, times int64) {
	if rdlst.MaxReturn < ret {
		rdlst.MaxReturn = ret
		rdlst.MaxReturnNums = times
	} else if rdlst.MaxReturn == ret {
		rdlst.MaxReturnNums++
	}

	for i, v := range rdlst.Returns {
		if v == ret {
			rdlst.ReturnWeights[i] += times

			return
		}
	}

	rdlst.Returns = append(rdlst.Returns, ret)
	rdlst.ReturnWeights = append(rdlst.ReturnWeights, times)
}

// Merge -
func (rdlst *RTPReturnDataList) Merge(lst *RTPReturnDataList) {
	for i, v := range lst.Returns {
		rdlst.addReturnsEx(v, lst.ReturnWeights[i])
	}
}

// Clone -
func (rdlst *RTPReturnDataList) Clone() *RTPReturnDataList {
	return &RTPReturnDataList{
		Returns:       rdlst.Returns[0:],
		ReturnWeights: rdlst.ReturnWeights[0:],
		MaxReturn:     rdlst.MaxReturn,
		MaxReturnNums: rdlst.MaxReturnNums,
		onResults:     rdlst.onResults,
	}
}

// SaveReturns2CSV -
func (rdlst *RTPReturnDataList) SaveReturns2CSV(fn string) error {
	results := []*RTPReturnData{}
	totaltimes := int64(0)
	for i, v := range rdlst.Returns {
		results = addResults2(results, v, rdlst.ReturnWeights[i])

		totaltimes += rdlst.ReturnWeights[i]
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Return < results[j].Return
	})

	f, err := os.Create(fn)
	if err != nil {
		goutils.Error("sgc7rtp.RTPReturnDataList.SaveReturns2CSV",
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
