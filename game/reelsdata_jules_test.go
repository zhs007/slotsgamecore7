package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidRI5_Jules(t *testing.T) {
	assert.False(t, isValidRI5(nil), "Test with nil")
	assert.False(t, isValidRI5([]reelsInfo5{}), "Test with empty")
	assert.False(t, isValidRI5([]reelsInfo5{{Line: 0}}), "Test with 0")
	assert.True(t, isValidRI5([]reelsInfo5{{Line: 1}}), "Test with 1")
	assert.True(t, isValidRI5([]reelsInfo5{{Line: 0}, {Line: 1}}), "Test with 0 and 1")
}

func Test_LoadReels3JSON_Jules(t *testing.T) {
	rd, err := LoadReels3JSON("testdata/reels3.json")
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	assert.Equal(t, 3, len(rd.Reels))
	assert.Equal(t, []int{1, 4}, rd.Reels[0])
	assert.Equal(t, []int{2, 5}, rd.Reels[1])
	assert.Equal(t, []int{3, 6}, rd.Reels[2])

	rd, err = LoadReels3JSON("testdata/invalidfile.json")
	assert.Error(t, err)
	assert.Nil(t, rd)
}

func Test_LoadReels5JSON_Jules(t *testing.T) {
	rd, err := LoadReels5JSON("testdata/reels5.json")
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	assert.Equal(t, 5, len(rd.Reels))
	assert.Equal(t, []int{1, 6}, rd.Reels[0])
	assert.Equal(t, []int{2, 7}, rd.Reels[1])
	assert.Equal(t, []int{3, 8}, rd.Reels[2])
	assert.Equal(t, []int{4, 9}, rd.Reels[3])
	assert.Equal(t, []int{5, 10}, rd.Reels[4])

	rd, err = LoadReels5JSON("testdata/invalidfile.json")
	assert.Error(t, err)
	assert.Nil(t, rd)
}

func Test_ReelsData_SetReel_Jules(t *testing.T) {
	rd := NewReelsData(3)
	rd.SetReel(1, []int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, rd.Reels[1])
}

func Test_ReelsData_DropDownIntoGameScene_Jules(t *testing.T) {
	rd := &ReelsData{
		Reels: [][]int{
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15},
		},
	}
	scene := &GameScene{
		Arr: [][]int{
			{-1, -1, 1},
			{2, -1, -1},
			{-1, 3, -1},
		},
	}
	indexes := []int{2, 3, 4}

	newIndexes, err := rd.DropDownIntoGameScene(scene, indexes)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, newIndexes)

	expectedScene := &GameScene{
		Arr: [][]int{
			{2, 1, 1},
			{2, 8, 7},
			{14, 3, 13},
		},
	}
	assert.Equal(t, expectedScene.Arr, scene.Arr)
}

func Test_ReelsData_DropDownIntoGameScene2_Jules(t *testing.T) {
	rd := &ReelsData{
		Reels: [][]int{
			{1, 2, 3, 4, 5},
			{6, 7, 8, 9, 10},
			{11, 12, 13, 14, 15},
		},
	}
	scene := &GameScene{
		Arr: [][]int{
			{-1, -1, 1},
			{2, -1, -1},
			{-1, 3, -1},
		},
	}
	indexes := []int{2, 3, 4}

	newIndexes, err := rd.DropDownIntoGameScene2(scene, indexes)
	assert.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, newIndexes)

	expectedScene := &GameScene{
		Arr: [][]int{
			{1, 2, 1},
			{2, 7, 8},
			{13, 3, 14},
		},
	}
	assert.Equal(t, expectedScene.Arr, scene.Arr)
}

func Test_ReelsData_BuildReelsPosData_Jules(t *testing.T) {
	rd := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}
	rpd, err := rd.BuildReelsPosData(func(rd *ReelsData, x, y int) bool {
		return rd.Reels[x][y]%2 == 0
	})
	assert.NoError(t, err)
	assert.NotNil(t, rpd)
	assert.Equal(t, 2, len(rpd.ReelsPos))
	assert.Equal(t, 1, len(rpd.ReelsPos[0]))
	assert.Equal(t, 2, len(rpd.ReelsPos[1]))
	assert.Equal(t, []int{1}, rpd.ReelsPos[0])
	assert.Equal(t, []int{0, 2}, rpd.ReelsPos[1])
}

func Test_NewReelsData_Jules(t *testing.T) {
	rd := NewReelsData(5)
	assert.NotNil(t, rd)
	assert.Equal(t, 5, len(rd.Reels))
}
