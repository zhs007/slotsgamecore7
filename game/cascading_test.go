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
