package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcClusterResult(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test 1: Basic cluster with same symbols
	scene1, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene1.InitWithArr2([][]int{
		{1, 1, 2},
		{1, 1, 3},
		{2, 3, 4},
	})
	assert.NoError(t, err)

	results1, err := CalcClusterResult(scene1, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.NotNil(t, results1)
	assert.Greater(t, len(results1), 0)
	assert.Equal(t, results1[0].Symbol, 1)
	assert.Equal(t, len(results1[0].Pos), 8) // 4 positions * 2 coordinates

	// Test 2: Cluster with wild symbols
	scene2, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene2.InitWithArr2([][]int{
		{1, 0, 1},
		{0, 1, 0},
		{1, 0, 1},
	})
	assert.NoError(t, err)

	results2, err := CalcClusterResult(scene2, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol || cursymbol == 0
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.NotNil(t, results2)
	assert.Greater(t, len(results2), 0)

	// Test 3: No valid clusters
	scene3, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene3.InitWithArr2([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	assert.NoError(t, err)

	results3, err := CalcClusterResult(scene3, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.Nil(t, results3)

	// Test 4: Large cluster across the grid
	scene4, err := NewGameScene(5, 5)
	assert.NoError(t, err)
	err = scene4.InitWithArr2([][]int{
		{1, 1, 1, 1, 1},
		{1, 2, 1, 2, 1},
		{1, 1, 1, 1, 1},
		{2, 1, 2, 1, 2},
		{1, 1, 1, 1, 1},
	})
	assert.NoError(t, err)

	results4, err := CalcClusterResult(scene4, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.NotNil(t, results4)
	assert.Greater(t, len(results4), 0)
	assert.Greater(t, len(results4[0].Pos), 10) // Should have a large cluster

	// Test 5: Multiple separate clusters
	scene5, err := NewGameScene(4, 4)
	assert.NoError(t, err)
	err = scene5.InitWithArr2([][]int{
		{1, 1, 2, 2},
		{1, 1, 2, 2},
		{3, 3, 4, 4},
		{3, 3, 4, 4},
	})
	assert.NoError(t, err)

	results5, err := CalcClusterResult(scene5, pt, 10, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.NotNil(t, results5)
	assert.Equal(t, 4, len(results5)) // Should have 4 separate clusters

	// Test 6: Custom symbol validation
	scene7, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene7.InitWithArr2([][]int{
		{-1, 1, 1},
		{1, -1, 1},
		{1, 1, -1},
	})
	assert.NoError(t, err)

	results7, err := CalcClusterResult(scene7, pt, 10, func(cursymbol int) bool {
		return cursymbol > 0 // Only positive numbers are valid
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	}, func(cursymbol int) int {
		return cursymbol
	})

	assert.NoError(t, err)
	assert.NotNil(t, results7)

	t.Logf("Test_CalcClusterResult OK")
}

func Test_CalcClusterSymbol(t *testing.T) {
	// Test 1: Basic cluster formation
	scene1, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene1.InitWithArr2([][]int{
		{1, 1, 2},
		{1, 1, 3},
		{2, 3, 4},
	})
	assert.NoError(t, err)

	pos1 := calcClusterSymbol(scene1, 0, 0, 1, []int{}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	})

	assert.Equal(t, 8, len(pos1)) // 4 positions * 2 coordinates
	assert.True(t, containsPosition(pos1, 0, 0))
	assert.True(t, containsPosition(pos1, 0, 1))
	assert.True(t, containsPosition(pos1, 1, 0))
	assert.True(t, containsPosition(pos1, 1, 1))

	// Test 2: Single symbol cluster
	scene2, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene2.InitWithArr2([][]int{
		{1, 2, 2},
		{2, 2, 2},
		{2, 2, 2},
	})
	assert.NoError(t, err)

	pos2 := calcClusterSymbol(scene2, 0, 0, 1, []int{}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	})

	assert.Equal(t, 2, len(pos2)) // 1 position * 2 coordinates

	// Test 3: Full grid cluster
	scene3, err := NewGameScene(3, 3)
	assert.NoError(t, err)
	err = scene3.InitWithArr2([][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	})
	assert.NoError(t, err)

	pos3 := calcClusterSymbol(scene3, 0, 0, 1, []int{}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol
	})

	assert.Equal(t, 18, len(pos3)) // 9 positions * 2 coordinates

	t.Logf("Test_CalcClusterSymbol OK")
}

// Helper function to check if a position exists in the position array
func containsPosition(positions []int, x, y int) bool {
	for i := 0; i < len(positions); i += 2 {
		if positions[i] == x && positions[i+1] == y {
			return true
		}
	}
	return false
}
