package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_FastReelsRandomSP_New(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{0, 1, 2, 3, 4, 5, 6},
			{10, 11, 12, 13, 14, 15},
			{20, 21, 22, 23, 24},
		},
	}

	// This function selects even numbers from the reel
	onSelectEven := func(reels *ReelsData, x int, y int) []int {
		if reels.Reels[x][y]%2 == 0 {
			return []int{y}
		}
		return nil
	}

	frr := NewFastReelsRandomSP(reels, onSelectEven)
	assert.NotNil(t, frr, "Test_FastReelsRandomSP_New: NewFastReelsRandomSP should not return nil")
	assert.Equal(t, frr.Reels, reels, "Test_FastReelsRandomSP_New: Reels should be the same")

	expectedArrIndex := [][]int{
		{0, 2, 4, 6}, // Even indices for reel 0: 0, 2, 4, 6
		{1, 3, 5},    // Even indices for reel 1: 11, 13, 15 -> indices 1, 3, 5. Wait, no, values are 10,12,14. indices are 0,2,4
		{0, 2, 4},    // Even indices for reel 2: 20, 22, 24 -> indices 0, 2, 4
	}

	// Let's correct expectedArrIndex for reel 1
	expectedArrIndex[1] = []int{0, 2, 4} // Values 10, 12, 14 are at indices 0, 2, 4

	assert.Equal(t, len(expectedArrIndex), len(frr.ArrIndex), "Test_FastReelsRandomSP_New: ArrIndex length mismatch")
	for i, arr := range frr.ArrIndex {
		assert.Equal(t, expectedArrIndex[i], arr, "Test_FastReelsRandomSP_New: ArrIndex content mismatch for reel %d", i)
	}
}

func Test_FastReelsRandomSP_Random(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{0, 1, 2, 3, 4, 5, 6},
			{10, 11, 12, 13, 14, 15},
			{20, 21, 22, 23, 24},
		},
	}

	// This function selects all numbers
	onSelectAll := func(reels *ReelsData, x int, y int) []int {
		return []int{y}
	}

	frr := NewFastReelsRandomSP(reels, onSelectAll)
	assert.NotNil(t, frr, "Test_FastReelsRandomSP_Random: NewFastReelsRandomSP should not return nil")

	plugin := sgc7plugin.NewMockPlugin()
	// We want to select index 2 from reel 0 (len 7), index 1 from reel 1 (len 6), and index 4 from reel 2 (len 5)
	// Random will do cache[i] % len(arr)
	// Reel 0: len is 7. To get 2, we can use 2, 9, 16... Let's use 9. 9 % 7 = 2
	// Reel 1: len is 6. To get 1, we can use 1, 7, 13... Let's use 7. 7 % 6 = 1
	// Reel 2: len is 5. To get 4, we can use 4, 9, 14... Let's use 4. 4 % 5 = 4
	plugin.Cache = []int{9, 7, 4}

	randomIndices, err := frr.Random(plugin)
	assert.NoError(t, err, "Test_FastReelsRandomSP_Random: Random should not return an error")

	// frr.ArrIndex will have [[0,1,2,3,4,5,6], [0,1,2,3,4,5], [0,1,2,3,4]]
	// Reel 0: frr.ArrIndex[0] has len 7. Random(7) with cache 9 gives 9%7=2. So we get frr.ArrIndex[0][2] which is 2.
	// Reel 1: frr.ArrIndex[1] has len 6. Random(6) with cache 7 gives 7%6=1. So we get frr.ArrIndex[1][1] which is 1.
	// Reel 2: frr.ArrIndex[2] has len 5. Random(5) with cache 4 gives 4%5=4. So we get frr.ArrIndex[2][4] which is 4.
	expectedIndices := []int{2, 1, 4}

	assert.Equal(t, expectedIndices, randomIndices, "Test_FastReelsRandomSP_Random: The returned indices are not what we expected")
}
