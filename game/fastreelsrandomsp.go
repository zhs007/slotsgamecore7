package sgc7game

import (
	"context"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// FuncOnSelectReelIndex - onSelectReelIndex
type FuncOnSelectReelIndex func(reels *ReelsData, x int, y int) []int

// FastReelsRandomSP - fast random for a special scene, it's like scatter x 5
type FastReelsRandomSP struct {
	Reels    *ReelsData
	ArrIndex [][]int
}

// NewFastReelsRandomSP - new a FastReelsRandomSP
func NewFastReelsRandomSP(reels *ReelsData, onSelectReelIndex FuncOnSelectReelIndex) *FastReelsRandomSP {
	frr := &FastReelsRandomSP{
		Reels: reels,
	}

	for x, l := range reels.Reels {
		arr := []int{}

		for y := range l {
			carr := onSelectReelIndex(reels, x, y)
			if carr != nil {
				for _, v := range carr {
					arr = sgc7utils.InsUniqueIntSlice(arr, v)
				}
			}
		}

		frr.ArrIndex = append(frr.ArrIndex, arr)
	}

	return frr
}

// Random - random
func (frr *FastReelsRandomSP) Random(plugin sgc7plugin.IPlugin) ([]int, error) {
	arr := []int{}

	for _, l := range frr.ArrIndex {
		y, err := plugin.Random(context.Background(), len(l))
		if err != nil {
			return nil, err
		}

		arr = append(arr, l[y])
	}

	return arr, nil
}
