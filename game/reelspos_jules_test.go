package sgc7game

import (
	"context"
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_NewReelsPosData_Jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}
	rpd := NewReelsPosData(reels)
	assert.NotNil(t, rpd)
	assert.Equal(t, 2, len(rpd.ReelsPos))
	assert.Equal(t, 0, len(rpd.ReelsPos[0]))
	assert.Equal(t, 0, len(rpd.ReelsPos[1]))
}

func Test_ReelsPosData_AddPos_Jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}
	rpd := NewReelsPosData(reels)
	rpd.AddPos(0, 1)
	rpd.AddPos(1, 0)
	rpd.AddPos(1, 2)
	assert.Equal(t, []int{1}, rpd.ReelsPos[0])
	assert.Equal(t, []int{0, 2}, rpd.ReelsPos[1])

	// Test out of bounds
	rpd.AddPos(2, 0)
	assert.Equal(t, 2, len(rpd.ReelsPos))
}

func Test_ReelsPosData_RandReel_Jules(t *testing.T) {
	reels := &ReelsData{
		Reels: [][]int{
			{1, 2, 3},
			{4, 5, 6},
		},
	}
	rpd := NewReelsPosData(reels)
	rpd.AddPos(0, 1)
	rpd.AddPos(1, 0)
	rpd.AddPos(1, 2)

	plugin := sgc7plugin.NewMockPlugin()

	// Test case 1
	plugin.SetCache([]int{0})
	pos, err := rpd.RandReel(context.Background(), plugin, 0)
	assert.NoError(t, err)
	assert.Equal(t, 1, pos)

	// Test case 2
	plugin.SetCache([]int{1})
	pos, err = rpd.RandReel(context.Background(), plugin, 1)
	assert.NoError(t, err)
	assert.Equal(t, 2, pos)

	// Test invalid x
	pos, err = rpd.RandReel(context.Background(), plugin, 2)
	assert.Error(t, err)
	assert.Equal(t, -1, pos)

	pos, err = rpd.RandReel(context.Background(), plugin, -1)
	assert.Error(t, err)
	assert.Equal(t, -1, pos)
}
