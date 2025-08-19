package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RemoveSymbolWithResult_Jules(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}
	result := &PlayResult{
		Results: []*Result{
			{
				Pos: []int{0, 0, 0, 1}, // (0,0) and (0,1)
			},
		},
	}

	RemoveSymbolWithResult(scene, result)

	assert.Equal(t, -1, scene.Arr[0][0])
	assert.Equal(t, -1, scene.Arr[0][1])
	assert.Equal(t, 3, scene.Arr[0][2])
}

func Test_RemoveSymbolWithResult2_Jules(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}
	result := &PlayResult{
		Results: []*Result{
			{
				Symbol: 1,
				Pos:    []int{0, 0, 0, 1},
			},
			{
				Symbol: 2,
				Pos:    []int{1, 0, 1, 1},
			},
		},
	}

	// Only remove results with symbol 1
	canRemoveResult := func(r *Result) bool {
		return r.Symbol == 1
	}
	// Always remove symbols
	canRemoveSymbol := func(x, y int) bool {
		return true
	}

	RemoveSymbolWithResult2(scene, result, canRemoveResult, canRemoveSymbol)

	assert.Equal(t, -1, scene.Arr[0][0])
	assert.Equal(t, -1, scene.Arr[0][1])
	assert.Equal(t, 4, scene.Arr[1][0]) // Not removed
	assert.Equal(t, 5, scene.Arr[1][1]) // Not removed
}

func Test_DropDownSymbols_Jules(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			{1, -1, 3},
			{-1, 5, -1},
			{7, -1, 9},
		},
	}

	DropDownSymbols(scene)

	expected := [][]int{
		{-1, 1, 3},
		{-1, -1, 5},
		{-1, 7, 9},
	}

	assert.Equal(t, expected, scene.Arr)
}

func Test_DropDownSymbols2_Jules(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			{-1, 2, -1},
			{4, -1, 6},
			{-1, 8, -1},
		},
	}

	DropDownSymbols2(scene)

	expected := [][]int{
		{2, -1, -1},
		{4, 6, -1},
		{8, -1, -1},
	}

	assert.Equal(t, expected, scene.Arr)
}
