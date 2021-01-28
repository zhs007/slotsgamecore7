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
