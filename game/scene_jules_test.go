package sgc7game

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

var ErrTest = errors.New("test error")

func Test_NewGameScene_jules(t *testing.T) {
	gs, err := NewGameScene(5, 3)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 5)
	assert.Equal(t, gs.Height, 3)

	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			assert.Equal(t, gs.Arr[x][y], -1)
		}
	}

	t.Logf("Test_NewGameScene_jules OK")
}

func Test_NewGameScene2_jules(t *testing.T) {
	gs, err := NewGameScene2(5, 3, 7)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 5)
	assert.Equal(t, gs.Height, 3)

	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			assert.Equal(t, gs.Arr[x][y], 7)
		}
	}

	gs, err = NewGameScene2(5, 3, 0)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 5)
	assert.Equal(t, gs.Height, 3)

	for x := 0; x < 5; x++ {
		for y := 0; y < 3; y++ {
			assert.Equal(t, gs.Arr[x][y], 0)
		}
	}

	t.Logf("Test_NewGameScene2_jules OK")
}

func Test_NewGameSceneEx_jules(t *testing.T) {
	heights := []int{3, 4, 5}
	gs, err := NewGameSceneEx(heights)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 3)
	assert.Equal(t, gs.Height, 5)
	assert.Equal(t, len(gs.HeightEx), 3)

	for i, h := range heights {
		assert.Equal(t, len(gs.Arr[i]), h)
		for y := 0; y < h; y++ {
			assert.Equal(t, gs.Arr[i][y], -1)
		}
	}

	t.Logf("Test_NewGameSceneEx_jules OK")
}

func Test_NewGameSceneWithArr_jules(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 6}
	gs, err := NewGameSceneWithArr(2, 3, arr)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 2)
	assert.Equal(t, gs.Height, 3)
	assert.Equal(t, gs.Arr[0][0], 1)
	assert.Equal(t, gs.Arr[0][1], 2)
	assert.Equal(t, gs.Arr[0][2], 3)
	assert.Equal(t, gs.Arr[1][0], 4)
	assert.Equal(t, gs.Arr[1][1], 5)
	assert.Equal(t, gs.Arr[1][2], 6)

	_, err = NewGameSceneWithArr(2, 3, []int{1, 2, 3})
	assert.Error(t, err)

	t.Logf("Test_NewGameSceneWithArr_jules OK")
}

func Test_NewGameSceneWithArr2_jules(t *testing.T) {
	arr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	gs, err := NewGameSceneWithArr2(arr)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 2)
	assert.Equal(t, gs.Height, 3)
	assert.Equal(t, gs.Arr[0][0], 1)
	assert.Equal(t, gs.Arr[0][1], 2)
	assert.Equal(t, gs.Arr[0][2], 3)
	assert.Equal(t, gs.Arr[1][0], 4)
	assert.Equal(t, gs.Arr[1][1], 5)
	assert.Equal(t, gs.Arr[1][2], 6)

	arr2 := [][]int{
		{1, 2, 3},
		{4, 5},
	}
	_, err = NewGameSceneWithArr2(arr2)
	assert.Error(t, err)

	t.Logf("Test_NewGameSceneWithArr2_jules OK")
}

func Test_NewGameSceneWithArr2Ex_jules(t *testing.T) {
	arr := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	gs, err := NewGameSceneWithArr2Ex(arr)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 2)
	assert.Equal(t, gs.Height, 3)

	arr2 := [][]int{
		{1, 2, 3},
		{4, 5},
	}
	gs, err = NewGameSceneWithArr2Ex(arr2)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 2)
	assert.Equal(t, gs.Height, 3)
	assert.Equal(t, len(gs.HeightEx), 2)
	assert.Equal(t, gs.HeightEx[0], 3)
	assert.Equal(t, gs.HeightEx[1], 2)

	t.Logf("Test_NewGameSceneWithArr2Ex_jules OK")
}
func Test_NewGameSceneWithReels_jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}
	arr := []int{0, 0, 0}
	gs, err := NewGameSceneWithReels(reels, 3, 3, arr)
	assert.NoError(t, err)
	assert.NotNil(t, gs)
	assert.Equal(t, gs.Width, 3)
	assert.Equal(t, gs.Height, 3)

	assert.Equal(t, gs.Arr[0][0], 1)
	assert.Equal(t, gs.Arr[0][1], 2)
	assert.Equal(t, gs.Arr[0][2], 3)
	assert.Equal(t, gs.Arr[1][0], 4)
	assert.Equal(t, gs.Arr[1][1], 5)
	assert.Equal(t, gs.Arr[1][2], 6)
	assert.Equal(t, gs.Arr[2][0], 7)
	assert.Equal(t, gs.Arr[2][1], 8)
	assert.Equal(t, gs.Arr[2][2], 9)
}

func Test_ReplaceSymbol_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 1},
		{2, 1, 2},
	})

	gs.ReplaceSymbol(1, 3)

	assert.Equal(t, gs.Arr[0][0], 3)
	assert.Equal(t, gs.Arr[0][1], 2)
	assert.Equal(t, gs.Arr[0][2], 3)
	assert.Equal(t, gs.Arr[1][0], 2)
	assert.Equal(t, gs.Arr[1][1], 3)
	assert.Equal(t, gs.Arr[1][2], 2)
}

func Test_Clear_jules(t *testing.T) {
	gs, _ := NewGameScene(3, 4)
	gs.Clear(5)

	for x := 0; x < 3; x++ {
		for y := 0; y < 4; y++ {
			assert.Equal(t, gs.Arr[x][y], 5)
		}
	}
}

func Test_Fill_jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}
	gs, _ := NewGameScene(2, 3)
	arr := []int{1, 2}
	gs.Fill(reels, arr)

	assert.Equal(t, gs.Arr[0][0], 2)
	assert.Equal(t, gs.Arr[0][1], 3)
	assert.Equal(t, gs.Arr[0][2], 1)
	assert.Equal(t, gs.Arr[1][0], 6)
	assert.Equal(t, gs.Arr[1][1], 4)
	assert.Equal(t, gs.Arr[1][2], 5)
}

func Test_ResetReelIndex2_jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
		},
	}
	gs, _ := NewGameScene(1, 3)

	err := gs.ResetReelIndex2(reels, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 3, gs.Arr[0][0])

	err = gs.ResetReelIndex2(reels, 0, 3)
	assert.NoError(t, err)
	assert.Equal(t, 1, gs.Arr[0][0])

	err = gs.ResetReelIndex2(reels, -1, 1)
	assert.Error(t, err)

	err = gs.ResetReelIndex2(reels, 1, 1)
	assert.Error(t, err)

	err = gs.ResetReelIndex2(reels, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, gs.Arr[0][0])
}

type mockPluginForRand struct {
	sgc7plugin.IPlugin
	results []int
	cursor  int
}

func (p *mockPluginForRand) Random(ctx context.Context, r int) (int, error) {
	if p.cursor < len(p.results) {
		v := p.results[p.cursor]
		p.cursor++
		return v, nil
	}
	return -1, ErrTest
}

func Test_RandReelsWithReelData_jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15},
		},
	}
	gs, _ := NewGameScene(3, 3)
	plugin := &mockPluginForRand{
		results: []int{0, 1, 2},
	}

	err := gs.RandReelsWithReelData(reels, plugin)
	assert.NoError(t, err)

	assert.Equal(t, 1, gs.Arr[0][0])
	assert.Equal(t, 2, gs.Arr[0][1])
	assert.Equal(t, 3, gs.Arr[0][2])

	assert.Equal(t, 7, gs.Arr[1][0])
	assert.Equal(t, 8, gs.Arr[1][1])
	assert.Equal(t, 9, gs.Arr[1][2])

	assert.Equal(t, 13, gs.Arr[2][0])
	assert.Equal(t, 14, gs.Arr[2][1])
	assert.Equal(t, 15, gs.Arr[2][2])
}

func Test_ForEach_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2},
		{3, 4},
	})
	count := 0
	gs.ForEach(func(x, y, val int) {
		assert.Equal(t, gs.Arr[x][y], val)
		count++
	})
	assert.Equal(t, 4, count)
}

func Test_ForEachAround_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	count := 0
	gs.ForEachAround(1, 1, func(x, y, val int) {
		count++
	})
	assert.Equal(t, 8, count)

	gsWithEx, _ := NewGameSceneEx([]int{2, 3, 2})
	count = 0
	gsWithEx.ForEachAround(1, 1, func(x, y, val int) {
		count++
	})
	assert.Equal(t, 6, count)
}

func Test_CountSymbol_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 1, 2},
		{3, 1, 4},
	})
	assert.Equal(t, 3, gs.CountSymbol(1))
	assert.Equal(t, 0, gs.CountSymbol(5))
}

func Test_CountSymbols_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{4, 5, 6},
	})
	counts := gs.CountSymbols([]int{1, 3, 5, 7})
	assert.Equal(t, []int{1, 1, 1, 0}, counts)
}

func Test_CountSymbolEx_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{4, 5, 6},
	})
	count := gs.CountSymbolEx(func(s, x, y int) bool {
		return s%2 == 0
	})
	assert.Equal(t, 3, count)
}

func Test_HasSymbol_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
	})
	assert.True(t, gs.HasSymbol(2))
	assert.False(t, gs.HasSymbol(4))
}

func Test_HasSymbols_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
	})
	assert.True(t, gs.HasSymbols([]int{4, 5, 1}))
	assert.False(t, gs.HasSymbols([]int{4, 5, 6}))
}

func Test_IsValidPos_jules(t *testing.T) {
	gs, _ := NewGameScene(3, 4)
	assert.True(t, gs.IsValidPos(0, 0))
	assert.True(t, gs.IsValidPos(2, 3))
	assert.False(t, gs.IsValidPos(3, 3))
	assert.False(t, gs.IsValidPos(2, 4))
	assert.False(t, gs.IsValidPos(-1, 0))
	assert.False(t, gs.IsValidPos(0, -1))
}

func Test_ToString_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2},
		{3, 4},
	})
	assert.Equal(t, "[[1,2],[3,4]]", gs.ToString())

	gs.Arr = [][]int{
		{1, 2},
		{3, 4},
	}
	gs.Arr[0] = nil
	assert.Equal(t, "[null,[3,4]]", gs.ToString())
}

func Test_Clone_jules(t *testing.T) {
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2},
		{3, 4},
	})
	gs.Indexes = []int{1, 2}
	gs.HeightEx = []int{2, 2}
	gs.ReelName = "test"

	clone := gs.Clone()

	assert.Equal(t, gs.Width, clone.Width)
	assert.Equal(t, gs.Height, clone.Height)
	assert.Equal(t, gs.ReelName, clone.ReelName)
	assert.Equal(t, len(gs.Arr), len(clone.Arr))
	assert.Equal(t, len(gs.Indexes), len(clone.Indexes))
	assert.Equal(t, len(gs.HeightEx), len(clone.HeightEx))

	gs.Arr[0][0] = 5
	gs.Indexes[0] = 5
	gs.HeightEx[0] = 5

	assert.NotEqual(t, gs.Arr[0][0], clone.Arr[0][0])
	assert.NotEqual(t, gs.Indexes[0], clone.Indexes[0])
	assert.NotEqual(t, gs.HeightEx[0], clone.HeightEx[0])
}

func Test_CloneEx_jules(t *testing.T) {
	pool := NewGameScenePoolEx()
	gs, _ := NewGameSceneWithArr2([][]int{
		{1, 2},
		{3, 4},
	})
	gs.Indexes = []int{1, 2}
	gs.HeightEx = []int{2, 2}
	gs.ReelName = "test"

	clone := gs.CloneEx(pool)

	assert.Equal(t, gs.Width, clone.Width)
	assert.Equal(t, gs.Height, clone.Height)
	assert.Equal(t, gs.ReelName, clone.ReelName)
	assert.Equal(t, len(gs.Arr), len(clone.Arr))
	assert.Equal(t, len(gs.Indexes), len(clone.Indexes))
	assert.Equal(t, len(gs.HeightEx), len(clone.HeightEx))

	gs.Arr[0][0] = 5
	gs.Indexes[0] = 5
	gs.HeightEx[0] = 5

	assert.NotEqual(t, gs.Arr[0][0], clone.Arr[0][0])
	assert.NotEqual(t, gs.Indexes[0], clone.Indexes[0])
	assert.NotEqual(t, gs.HeightEx[0], clone.HeightEx[0])
}

func Test_isArrEx_jules(t *testing.T) {
	assert.False(t, isArrEx([][]int{
		{1, 2},
		{3, 4},
	}))

	assert.True(t, isArrEx([][]int{
		{1, 2},
		{3},
	}))
}
