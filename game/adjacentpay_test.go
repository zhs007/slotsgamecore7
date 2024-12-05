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

	scene2, err := NewGameSceneWithArr2([][]int{
		{0, -1, 8, 3, 3},
		{4, 5, 8, 4, 5},
		{6, 6, 0, 7, 7},
		{3, 5, 9, 2, 0},
		{-1, -1, 9, -1, -1},
	})
	assert.NoError(t, err)

	result2, err := CalcAdjacentPay(scene2, pt, 10, func(cursymbol int) bool {
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
	assert.Equal(t, len(result2), 4)

	t.Logf("Test_CalcAdjacentPay OK")
}

func Test_CalcAdjacentPayEdgeCases(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test empty scene
	scene, err := NewGameSceneWithArr2([][]int{
		{-1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1},
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
	assert.Equal(t, len(result), 0)

	// Test all wild symbols
	scene2, err := NewGameSceneWithArr2([][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	})
	assert.NoError(t, err)

	result2, err := CalcAdjacentPay(scene2, pt, 10, func(cursymbol int) bool {
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
	assert.Greater(t, len(result2), 0)

	// Test vertical adjacent pays
	scene3, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3, 4, 5},
		{1, 7, 3, 4, 5},
		{1, 8, 3, 4, 5},
		{0, 9, 3, 4, 5},
		{1, 0, 3, 4, 5},
	})
	assert.NoError(t, err)

	result3, err := CalcAdjacentPay(scene3, pt, 10, func(cursymbol int) bool {
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
	assert.Greater(t, len(result3), 0)

	// Test mixed horizontal and vertical pays
	scene4, err := NewGameSceneWithArr2([][]int{
		{2, 2, 2, 4, 5},
		{2, 7, 3, 4, 5},
		{2, 8, 3, 4, 5},
		{0, 9, 3, 4, 5},
		{2, 0, 3, 4, 5},
	})
	assert.NoError(t, err)

	result4, err := CalcAdjacentPay(scene4, pt, 10, func(cursymbol int) bool {
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
	assert.Greater(t, len(result4), 0)

	// Test scene with all wild symbols
	sceneAllWild, err := NewGameSceneWithArr2([][]int{
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0},
	})
	assert.NoError(t, err)

	resultAllWild, err := CalcAdjacentPay(sceneAllWild, pt, 10, func(cursymbol int) bool {
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
	assert.NotNil(t, resultAllWild)
	assert.Equal(t, 10, len(resultAllWild))

	// Test scene with alternating symbols and wilds
	sceneAlternate, err := NewGameSceneWithArr2([][]int{
		{1, 0, 1, 0, 1},
		{0, 2, 0, 2, 0},
		{3, 0, 3, 0, 3},
		{0, 4, 0, 4, 0},
		{5, 0, 5, 0, 5},
	})
	assert.NoError(t, err)

	resultAlternate, err := CalcAdjacentPay(sceneAlternate, pt, 10, func(cursymbol int) bool {
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
	assert.NotNil(t, resultAlternate)

	// Test scene with single valid line
	sceneSingleLine, err := NewGameSceneWithArr2([][]int{
		{-1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1},
		{1, 1, 1, 1, 1},
		{-1, -1, -1, -1, -1},
		{-1, -1, -1, -1, -1},
	})
	assert.NoError(t, err)

	resultSingleLine, err := CalcAdjacentPay(sceneSingleLine, pt, 10, func(cursymbol int) bool {
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
	assert.Equal(t, 1, len(resultSingleLine))

	// Test scene with wild at start positions
	sceneWildStart, err := NewGameSceneWithArr2([][]int{
		{0, 1, 1, 1, 1},
		{0, 2, 2, 2, 2},
		{0, 3, 3, 3, 3},
		{0, 4, 4, 4, 4},
		{0, 5, 5, 5, 5},
	})
	assert.NoError(t, err)

	resultWildStart, err := CalcAdjacentPay(sceneWildStart, pt, 10, func(cursymbol int) bool {
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
	assert.NotNil(t, resultWildStart)

	// Test scene with wild at end positions
	sceneWildEnd, err := NewGameSceneWithArr2([][]int{
		{1, 1, 1, 1, 0},
		{2, 2, 2, 2, 0},
		{3, 3, 3, 3, 0},
		{4, 4, 4, 4, 0},
		{5, 5, 5, 5, 0},
	})
	assert.NoError(t, err)

	resultWildEnd, err := CalcAdjacentPay(sceneWildEnd, pt, 10, func(cursymbol int) bool {
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
	assert.NotNil(t, resultWildEnd)

	t.Logf("Test_CalcAdjacentPayEdgeCases OK")
}
