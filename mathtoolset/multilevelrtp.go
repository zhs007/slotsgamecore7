package mathtoolset

import (
	"math"

	"github.com/zhs007/goutils"
	"go.uber.org/zap"
)

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
				zap.Error(ErrCannotBeConverged))

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
				zap.Error(ErrCannotBeConverged))

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
			addnum += levelUpAddSpinNum[prelevel+i]
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
