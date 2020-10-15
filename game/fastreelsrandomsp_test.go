package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_FastReelsRandomSP(t *testing.T) {
	bp := sgc7plugin.NewBasicPlugin()

	reels, err := LoadReels5JSON("../unittestdata/reels.json")
	assert.NoError(t, err)

	frr := NewFastReelsRandomSP(reels, func(r *ReelsData, x int, y int) []int {
		if r.Reels[x][y] == 9 {
			// X
			// ?
			// ?
			arr := []int{y}

			// ?
			// X
			// ?
			if y >= 1 {
				arr = append(arr, y-1)
			}

			// ?
			// ?
			// X
			if y >= 2 {
				arr = append(arr, y-2)
			}

			return arr
		}

		return nil
	})

	for i := 0; i < 100000; i++ {
		arr, err := frr.Random(bp)
		assert.NoError(t, err)

		gs, err := NewGameSceneWithReels(reels, 5, 3, arr)
		assert.NoError(t, err)

		// 9 is on every wheel
		for x := range arr {
			assert.Equal(t, gs.Arr[x][0] == 9 || gs.Arr[x][1] == 9 || gs.Arr[x][2] == 9, true)
		}
	}

	t.Logf("Test_FastReelsRandomSP OK")
}
