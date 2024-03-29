package mathtoolset

import (
	"math"

	"github.com/xuri/excelize/v2"
	"github.com/zhs007/goutils"
)

type MultiLevelRTPNode struct {
	SpinNum     int
	EndingLevel int
	RTP         float64
	Percent     float64
	TotalRTP    float64
}

type MultiLevelRTPData struct {
	Nodes []*MultiLevelRTPNode
}

func (rtpdata *MultiLevelRTPData) add(spinnum int, endinglevel int, rtp float64, per float64) {
	// for _, v := range rtpdata.Nodes {
	// 	if v.SpinNum == spinnum && v.EndingLevel == endinglevel {
	// 		v.RTP += rtp
	// 		v.Percent += per

	// 		return
	// 	}
	// }

	rtpdata.Nodes = append(rtpdata.Nodes, &MultiLevelRTPNode{
		SpinNum:     spinnum,
		EndingLevel: endinglevel,
		RTP:         rtp,
		Percent:     per,
		TotalRTP:    rtp * per,
	})
}

func (rtpdata *MultiLevelRTPData) addEx(spinnum int, endinglevel int, rtp float64, per float64) {
	for _, v := range rtpdata.Nodes {
		if v.SpinNum == spinnum && v.EndingLevel == endinglevel {
			v.TotalRTP += rtp * per
			v.Percent += per

			return
		}
	}

	rtpdata.Nodes = append(rtpdata.Nodes, &MultiLevelRTPNode{
		SpinNum:     spinnum,
		EndingLevel: endinglevel,
		TotalRTP:    rtp * per,
		Percent:     per,
	})
}

func (rtpdata *MultiLevelRTPData) calcMulLevelRTP2(prelevel int, levelRTPs []float64, levelUpProbs []map[int]float64, spinNum int, levelUpAddSpinNum []int,
	totalSpinNum int, totalRTP float64, curPer float64) float64 {

	// 如果最后一次spin了
	if spinNum == 0 {
		return 0
	}

	if prelevel >= len(levelRTPs) {
		prelevel = len(levelRTPs) - 1
	}

	// 	已经到最高级了，这里不考虑最高级依然可以增加次数，所以直接返回即可
	if prelevel == len(levelRTPs)-1 {
		rtpdata.add(totalSpinNum+spinNum, prelevel, totalRTP+levelRTPs[prelevel]*float64(spinNum), curPer)

		return levelRTPs[prelevel] * float64(spinNum)
	}

	currtp := levelRTPs[prelevel]

	mapProbs := levelUpProbs[prelevel]

	for k, v := range mapProbs {
		if k == 0 {
			// 如果最后一次
			if spinNum == 1 {
				rtpdata.add(totalSpinNum+spinNum, prelevel, totalRTP+levelRTPs[prelevel], curPer*v)
			} else {
				currtp += rtpdata.calcMulLevelRTP2(prelevel, levelRTPs, levelUpProbs, spinNum-1, levelUpAddSpinNum, totalSpinNum+1, totalRTP+levelRTPs[prelevel], curPer*v) * v
			}
		} else {
			addnum := 0
			for i := 1; i <= k; i++ {
				cl := prelevel + i

				if cl < len(levelRTPs) {
					addnum += levelUpAddSpinNum[cl]
				}
			}

			if spinNum-1+addnum > 0 {
				// 考虑升级的情况
				currtp += rtpdata.calcMulLevelRTP2(prelevel+k, levelRTPs, levelUpProbs, spinNum-1+addnum, levelUpAddSpinNum, totalSpinNum+1, totalRTP+levelRTPs[prelevel], curPer*v) * v
			} else if spinNum == 1 {
				rtpdata.add(totalSpinNum+spinNum, prelevel+k, totalRTP+levelRTPs[prelevel], curPer*v)
			}
		}
	}

	return currtp
}

func (rtpdata *MultiLevelRTPData) CalcMulLevelRTP2(levelRTPs []float64, levelUpProbs []map[int]float64, spinNum int, levelUpAddSpinNum []int) float64 {
	if spinNum <= 0 {
		return 0
	}

	if spinNum == 1 {
		rtpdata.add(1, 0, levelRTPs[0], 1)

		return levelRTPs[0]
	}

	levelUpProbs1 := formatProbsSlice(levelUpProbs)

	return rtpdata.calcMulLevelRTP2(0, levelRTPs, levelUpProbs1, spinNum, levelUpAddSpinNum, 0, 0, 1)
}

func (rtpdata *MultiLevelRTPData) Format() *MultiLevelRTPData {
	rtpdata1 := NewMultiLevelRTPData()

	for _, n := range rtpdata.Nodes {
		rtpdata1.addEx(n.SpinNum, n.EndingLevel, n.RTP, n.Percent)
	}

	return rtpdata1
}

func (rtpdata *MultiLevelRTPData) SaveResults(fn string) error {
	f := excelize.NewFile()

	sheet := f.GetSheetList()[0]

	f.SetCellStr(sheet, goutils.Pos2Cell(0, 0), "spinNum")
	f.SetCellStr(sheet, goutils.Pos2Cell(1, 0), "EndingLevel")
	f.SetCellStr(sheet, goutils.Pos2Cell(2, 0), "Percent")
	f.SetCellStr(sheet, goutils.Pos2Cell(3, 0), "RTP")
	f.SetCellStr(sheet, goutils.Pos2Cell(4, 0), "TotalRTP")

	si := 1

	for _, v := range rtpdata.Nodes {
		f.SetCellInt(sheet, goutils.Pos2Cell(0, si), v.SpinNum)
		f.SetCellInt(sheet, goutils.Pos2Cell(1, si), v.EndingLevel)
		f.SetCellFloat(sheet, goutils.Pos2Cell(2, si), v.Percent, 5, 64)
		f.SetCellFloat(sheet, goutils.Pos2Cell(3, si), v.RTP, 5, 64)
		f.SetCellFloat(sheet, goutils.Pos2Cell(4, si), v.TotalRTP, 5, 64)

		si++
	}

	return f.SaveAs(fn)
}

func NewMultiLevelRTPData() *MultiLevelRTPData {
	return &MultiLevelRTPData{}
}

func calcMulLevelRTP(prelevel int, levelRTPs []float64, levelUpProbs []float64, spinNum int, levelUpAddSpinNum []int) float64 {
	// 如果最后一次spin了
	if spinNum == 0 {
		return 0
	}

	// 	已经到最高级了
	if prelevel == len(levelRTPs)-1 {
		// 如果不能增加次数，则直接返回即可
		if levelUpAddSpinNum[prelevel] == 0 {
			return levelRTPs[prelevel] * float64(spinNum)
		}

		x := float64(levelUpAddSpinNum[prelevel]) * levelUpProbs[prelevel]
		if x >= 1 {
			goutils.Error("calcMulLevelRTP",
				goutils.Err(ErrCannotBeConverged))

			return math.NaN()
		}

		// 否则需要返回无穷级数求和
		return levelRTPs[prelevel] * float64(spinNum) / (1 - x)
	}

	currtp := float64(0)

	// 先便利本次不能升级的情况
	currtp += calcMulLevelRTP(prelevel, levelRTPs, levelUpProbs, spinNum-1, levelUpAddSpinNum) * (1 - levelUpProbs[prelevel])

	// 再考虑升级的情况
	currtp += calcMulLevelRTP(prelevel+1, levelRTPs, levelUpProbs, spinNum-1+levelUpAddSpinNum[prelevel], levelUpAddSpinNum) * levelUpProbs[prelevel]

	return currtp
}

// 计算可升级的rtp，这里只考虑一次升级
func CalcMulLevelRTP(levelRTPs []float64, levelUpProbs []float64, spinNum int, levelUpAddSpinNum []int) float64 {
	if spinNum <= 0 {
		return 0
	}

	if spinNum == 1 {
		return levelRTPs[0]
	}

	return levelRTPs[0] + calcMulLevelRTP(0, levelRTPs, levelUpProbs, spinNum-1, levelUpAddSpinNum)
}

func calcMulLevelRTP2(prelevel int, levelRTPs []float64, levelUpProbs []map[int]float64, spinNum int, levelUpAddSpinNum []int) float64 {
	if prelevel >= len(levelRTPs) {
		prelevel = len(levelRTPs) - 1
	}

	// 如果最后一次spin了
	if spinNum == 0 {
		return 0
	}

	// 	已经到最高级了
	if prelevel == len(levelRTPs)-1 {
		// 如果不能增加次数，则直接返回即可
		if levelUpAddSpinNum[prelevel] == 0 {
			return levelRTPs[prelevel] * float64(spinNum)
		}

		mapProbs := levelUpProbs[prelevel]
		x := float64(0)
		for k, v := range mapProbs {
			x += float64(k) * v
		}

		if x >= 1 {
			goutils.Error("calcMulLevelRTP2",
				goutils.Err(ErrCannotBeConverged))

			return math.NaN()
		}

		// 否则需要返回无穷级数求和
		return levelRTPs[prelevel] * float64(spinNum) / (1 - x)
	}

	currtp := float64(0)

	mapProbs := levelUpProbs[prelevel]

	noupprob := float64(1)

	for k, v := range mapProbs {
		addnum := 0
		for i := 0; i < k; i++ {
			cl := prelevel + i
			if cl >= len(levelRTPs) {
				cl = len(levelRTPs) - 1
			}

			addnum += levelUpAddSpinNum[cl]
		}

		// 考虑升级的情况
		currtp += calcMulLevelRTP2(prelevel+k, levelRTPs, levelUpProbs, spinNum-1+addnum, levelUpAddSpinNum) * v

		noupprob -= v
	}

	// 本次不能升级的情况
	currtp += calcMulLevelRTP2(prelevel, levelRTPs, levelUpProbs, spinNum-1, levelUpAddSpinNum) * noupprob

	return currtp
}

// 计算可升级的rtp，这里只考虑一次升级
func CalcMulLevelRTP2(levelRTPs []float64, levelUpProbs []map[int]float64, spinNum int, levelUpAddSpinNum []int) float64 {
	if spinNum <= 0 {
		return 0
	}

	if spinNum == 1 {
		return levelRTPs[0]
	}

	return levelRTPs[0] + calcMulLevelRTP2(0, levelRTPs, levelUpProbs, spinNum-1, levelUpAddSpinNum)
}

func formatProbsSlice(arrProbs []map[int]float64) []map[int]float64 {
	newarr := []map[int]float64{}
	for _, m := range arrProbs {
		newm := make(map[int]float64)

		totalv := float64(0)

		for _, v := range m {
			totalv += v
		}

		for k, v := range m {
			newm[k] = v / totalv
		}

		newarr = append(newarr, newm)
	}

	return newarr
}
