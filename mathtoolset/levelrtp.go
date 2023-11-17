package mathtoolset

type LevelRTPCollectData struct {
	Collector int
	Weight    float64
}

type LevelRTPData struct {
	Collector            map[int]*LevelRTPCollectData
	TotalCollectorWeight float64
	RTP                  float64
}

func (levelRTPData *LevelRTPData) AddCollector(collector int, weight float64) {
	levelRTPData.Collector[collector] = &LevelRTPCollectData{
		Collector: collector,
		Weight:    weight,
	}

	levelRTPData.TotalCollectorWeight += weight
}

func NewLevelRTPData(rtp float64) *LevelRTPData {
	return &LevelRTPData{
		Collector:            make(map[int]*LevelRTPCollectData),
		TotalCollectorWeight: 0,
		RTP:                  rtp,
	}
}

type LevelRTP struct {
	LevelData map[int]*LevelRTPData
	MaxLevel  int
}

func (levelRTP *LevelRTP) AddLevelData(level int, data *LevelRTPData) {
	levelRTP.LevelData[level] = data
}

// 计算RTP，一般来说，这种用于FG的计算，num是免费次数，awardNum 是到达某个level奖励的免费次数
func (levelRTP *LevelRTP) CalcRTP(startLevel int, num int, awardNum map[int]int) float64 {
	if num <= 0 {
		return 0
	}

	ld := levelRTP.LevelData[startLevel]
	curRTP := levelRTP.LevelData[startLevel].RTP

	for c, cd := range ld.Collector {
		curln := num
		curcn := startLevel + c

		if curcn != startLevel {
			if awardNum[curcn] > 0 {
				curln += awardNum[curcn]
			}
		}

		if curln > 1 {
			curRTP += levelRTP.calcRTP(curcn, curln-1, awardNum) * cd.Weight / ld.TotalCollectorWeight
		}
	}

	return curRTP
}

func (levelRTP *LevelRTP) calcRTP(curLevel int, lastnum int, awardNum map[int]int) float64 {
	if lastnum <= 0 {
		return 0
	}

	ld := levelRTP.LevelData[curLevel]
	curRTP := levelRTP.LevelData[curLevel].RTP

	for c, cd := range ld.Collector {
		curln := lastnum
		curcn := curLevel + c

		if curcn != curLevel {
			if awardNum[curcn] > 0 {
				curln += awardNum[curcn]
			}
		}

		if curln > 1 {
			curRTP += levelRTP.calcRTP(curcn, curln-1, awardNum) * cd.Weight / ld.TotalCollectorWeight
		}
	}

	return curRTP
}

func NewLevelRTP() *LevelRTP {
	return &LevelRTP{
		LevelData: make(map[int]*LevelRTPData),
	}
}
