package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcAdjacentPay(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	scene, err := NewGameSceneWithArr2([][]int{
		{1, 2, 2, 3, 3},
		{4, 5, 4, 4, 5},
		{-1, -1, -1, 6, 7},
		{3, 5, 0, 2, 0},
		{-1, -1, -1, -1, -1},
	})
	assert.NoError(t, err)

	result, err := CalcAdjacentPay(scene, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		if cursymbol == startsymbol {
			return true
		}

		return cursymbol == 0
	}, func(cursymbol int) int {
		return cursymbol
	})
	assert.NoError(t, err)
	assert.Equal(t, len(result), 1)

	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 2, 8, 3, 3},
		{4, 5, 8, 4, 5},
		{6, 6, 0, 7, 7},
		{3, 5, 9, 2, 0},
		{-1, -1, 9, -1, -1},
	})
	assert.NoError(t, err)

	result1, err := CalcAdjacentPay(scene1, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		if cursymbol == startsymbol {
			return true
		}

		return cursymbol == 0
	}, func(cursymbol int) int {
		return cursymbol
	})
	assert.NoError(t, err)
	assert.Equal(t, len(result1), 4)

	t.Logf("Test_CalcAdjacentPay OK")
}
