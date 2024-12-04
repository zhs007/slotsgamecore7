package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DropDownSymbols(t *testing.T) {
	gs, err := NewGameSceneWithArr2([][]int{
		{8, 10, -1},
		{11, -1, 7},
		{-1, 4, 6},
		{-1, 8, -1},
		{5, -1, -1},
	})
	assert.NoError(t, err)
	assert.NotNil(t, gs)

	err = DropDownSymbols(gs)
	assert.NoError(t, err)

	assert.Equal(t, gs.Arr[0][0], -1)
	assert.Equal(t, gs.Arr[0][1], 8)
	assert.Equal(t, gs.Arr[0][2], 10)

	assert.Equal(t, gs.Arr[1][0], -1)
	assert.Equal(t, gs.Arr[1][1], 11)
	assert.Equal(t, gs.Arr[1][2], 7)

	assert.Equal(t, gs.Arr[2][0], -1)
	assert.Equal(t, gs.Arr[2][1], 4)
	assert.Equal(t, gs.Arr[2][2], 6)

	assert.Equal(t, gs.Arr[3][0], -1)
	assert.Equal(t, gs.Arr[3][1], -1)
	assert.Equal(t, gs.Arr[3][2], 8)

	assert.Equal(t, gs.Arr[4][0], -1)
	assert.Equal(t, gs.Arr[4][1], -1)
	assert.Equal(t, gs.Arr[4][2], 5)

	t.Logf("Test_DropDownSymbols OK")
}

func Test_DropDownSymbols_2(t *testing.T) {
	gs, err := NewGameSceneWithArr2([][]int{
		{6, 3, 7, 8, 6, 6, 4},
		{8, 6, 8, 8, 4, 6, 7},
		{8, 1, 6, 3, 8, 7, 5},
		{4, 7, 8, 4, 8, 4, 7},
		{8, 8, 7, -1, -1, -1, 8},
		{7, 5, 5, 4, 8, -1, 6},
		{8, 7, 8, 3, 5, -1, 4},
	})
	assert.NoError(t, err)
	assert.NotNil(t, gs)

	err = DropDownSymbols(gs)
	assert.NoError(t, err)

	assert.Equal(t, gs.Arr[4][0], -1)
	assert.Equal(t, gs.Arr[4][1], -1)
	assert.Equal(t, gs.Arr[4][2], -1)
	assert.Equal(t, gs.Arr[4][3], 8)
	assert.Equal(t, gs.Arr[4][4], 8)
	assert.Equal(t, gs.Arr[4][5], 7)
	assert.Equal(t, gs.Arr[4][6], 8)

	assert.Equal(t, gs.Arr[5][0], -1)
	assert.Equal(t, gs.Arr[5][1], 7)
	assert.Equal(t, gs.Arr[5][2], 5)
	assert.Equal(t, gs.Arr[5][3], 5)
	assert.Equal(t, gs.Arr[5][4], 4)
	assert.Equal(t, gs.Arr[5][5], 8)
	assert.Equal(t, gs.Arr[5][6], 6)

	assert.Equal(t, gs.Arr[6][0], -1)
	assert.Equal(t, gs.Arr[6][1], 8)
	assert.Equal(t, gs.Arr[6][2], 7)
	assert.Equal(t, gs.Arr[6][3], 8)
	assert.Equal(t, gs.Arr[6][4], 3)
	assert.Equal(t, gs.Arr[6][5], 5)
	assert.Equal(t, gs.Arr[6][6], 4)

	t.Logf("Test_DropDownSymbols_2 OK")
}

func Test_DropDownSymbols2(t *testing.T) {
	gs, err := NewGameSceneWithArr2([][]int{
		{8, 10, -1},
		{11, -1, 7},
		{-1, 4, 6},
		{-1, 8, -1},
		{5, -1, -1},
	})
	assert.NoError(t, err)
	assert.NotNil(t, gs)

	err = DropDownSymbols2(gs)
	assert.NoError(t, err)

	assert.Equal(t, gs.Arr[0][0], 8)
	assert.Equal(t, gs.Arr[0][1], 10)
	assert.Equal(t, gs.Arr[0][2], -1)

	assert.Equal(t, gs.Arr[1][0], 11)
	assert.Equal(t, gs.Arr[1][1], 7)
	assert.Equal(t, gs.Arr[1][2], -1)

	assert.Equal(t, gs.Arr[2][0], 4)
	assert.Equal(t, gs.Arr[2][1], 6)
	assert.Equal(t, gs.Arr[2][2], -1)

	assert.Equal(t, gs.Arr[3][0], 8)
	assert.Equal(t, gs.Arr[3][1], -1)
	assert.Equal(t, gs.Arr[3][2], -1)

	assert.Equal(t, gs.Arr[4][0], 5)
	assert.Equal(t, gs.Arr[4][1], -1)
	assert.Equal(t, gs.Arr[4][2], -1)

	t.Logf("Test_DropDownSymbols2 OK")
}

func Test_RemoveSymbolWithResult(t *testing.T) {
	// Test case 1: Basic symbol removal
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 1, 2},
		{1, 1, 3},
		{2, 3, 4},
	})
	assert.NoError(t, err)
	assert.NotNil(t, scene1)

	result1 := &PlayResult{
		Results: []*Result{
			{
				Symbol: 1,
				Pos:    []int{0, 0, 0, 1, 1, 0, 1, 1}, // 2x2 cluster of symbol 1
			},
		},
	}

	err = RemoveSymbolWithResult(scene1, result1)
	assert.NoError(t, err)

	// Check if symbols were removed (set to -1)
	assert.Equal(t, -1, scene1.Arr[0][0])
	assert.Equal(t, -1, scene1.Arr[0][1])
	assert.Equal(t, -1, scene1.Arr[1][0])
	assert.Equal(t, -1, scene1.Arr[1][1])
	assert.Equal(t, 2, scene1.Arr[0][2]) // Unchanged
	assert.Equal(t, 3, scene1.Arr[1][2]) // Unchanged
	assert.Equal(t, 2, scene1.Arr[2][0]) // Unchanged
	assert.Equal(t, 3, scene1.Arr[2][1]) // Unchanged
	assert.Equal(t, 4, scene1.Arr[2][2]) // Unchanged

	// Test case 2: Multiple winning combinations
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 2},
		{3, 2, 2},
		{1, 2, 1},
	})
	assert.NoError(t, err)
	assert.NotNil(t, scene2)

	result2 := &PlayResult{
		Results: []*Result{
			{
				Symbol: 2,
				Pos:    []int{0, 1, 0, 2, 1, 1, 1, 2}, // Symbol 2 cluster
			},
			{
				Symbol: 1,
				Pos:    []int{0, 0, 2, 0, 2, 2}, // Symbol 1 positions
			},
		},
	}

	err = RemoveSymbolWithResult(scene2, result2)
	assert.NoError(t, err)

	// Check if all winning symbols were removed
	assert.Equal(t, -1, scene2.Arr[0][0])
	assert.Equal(t, -1, scene2.Arr[0][1])
	assert.Equal(t, -1, scene2.Arr[0][2])
	assert.Equal(t, 3, scene2.Arr[1][0])   // Unchanged
	assert.Equal(t, -1, scene2.Arr[1][1])
	assert.Equal(t, -1, scene2.Arr[1][2])
	assert.Equal(t, -1, scene2.Arr[2][0])
	assert.Equal(t, 2, scene2.Arr[2][1])   // Unchanged
	assert.Equal(t, -1, scene2.Arr[2][2])

	t.Logf("Test_RemoveSymbolWithResult OK")
}

func Test_RemoveSymbolWithResult2(t *testing.T) {
	// Test case 1: Selective symbol removal with custom conditions
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 1, 2},
		{1, 1, 3},
		{2, 3, 4},
	})
	assert.NoError(t, err)
	assert.NotNil(t, scene1)

	result1 := &PlayResult{
		Results: []*Result{
			{
				Symbol: 1,
				Pos:    []int{0, 0, 0, 1, 1, 0, 1, 1}, // 2x2 cluster of symbol 1
			},
		},
	}

	// Only remove symbols in even positions (x+y is even)
	err = RemoveSymbolWithResult2(scene1, result1, 
		func(r *Result) bool { return true }, // Accept all results
		func(x, y int) bool { return (x+y)%2 == 0 }, // Only remove on even positions
	)
	assert.NoError(t, err)

	// Check selective removal
	assert.Equal(t, -1, scene1.Arr[0][0])  // Even position (0,0)
	assert.Equal(t, 1, scene1.Arr[0][1])   // Odd position (0,1)
	assert.Equal(t, 1, scene1.Arr[1][0])   // Odd position (1,0)
	assert.Equal(t, -1, scene1.Arr[1][1])  // Even position (1,1)

	// Test case 2: Selective result removal
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 2},
		{3, 2, 2},
		{1, 2, 1},
	})
	assert.NoError(t, err)
	assert.NotNil(t, scene2)

	result2 := &PlayResult{
		Results: []*Result{
			{
				Symbol: 2,
				Pos:    []int{0, 1, 0, 2, 1, 1, 1, 2}, // Symbol 2 cluster
			},
			{
				Symbol: 1,
				Pos:    []int{0, 0, 2, 0, 2, 2}, // Symbol 1 positions
			},
		},
	}

	// Only remove results with symbol 2
	err = RemoveSymbolWithResult2(scene2, result2,
		func(r *Result) bool { return r.Symbol == 2 }, // Only process symbol 2
		func(x, y int) bool { return true }, // Remove all positions
	)
	assert.NoError(t, err)

	// Check if only symbol 2 clusters were removed
	assert.Equal(t, 1, scene2.Arr[0][0])   // Symbol 1 unchanged
	assert.Equal(t, -1, scene2.Arr[0][1])  // Symbol 2 removed
	assert.Equal(t, -1, scene2.Arr[0][2])  // Symbol 2 removed
	assert.Equal(t, 3, scene2.Arr[1][0])   // Unchanged
	assert.Equal(t, -1, scene2.Arr[1][1])  // Symbol 2 removed
	assert.Equal(t, -1, scene2.Arr[1][2])  // Symbol 2 removed
	assert.Equal(t, 1, scene2.Arr[2][0])   // Symbol 1 unchanged
	assert.Equal(t, 2, scene2.Arr[2][1])   // Unchanged
	assert.Equal(t, 1, scene2.Arr[2][2])   // Symbol 1 unchanged

	t.Logf("Test_RemoveSymbolWithResult2 OK")
}
