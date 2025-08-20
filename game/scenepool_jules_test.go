package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GameScenePoolEx_jules(t *testing.T) {
	pool := NewGameScenePoolEx()
	assert.NotNil(t, pool)

	assert.Equal(t, len(pool.MapPools), 0)

	s1 := pool.New(5, 3)
	assert.NotNil(t, s1)
	assert.Equal(t, s1.Width, 5)
	assert.Equal(t, s1.Height, 3)

	assert.Equal(t, len(pool.MapPools), 1)
	assert.Equal(t, len(pool.MapPools[5]), 1)
	assert.NotNil(t, pool.MapPools[5][3])
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 0)
	assert.Equal(t, len(pool.MapPools[5][3].used), 1)

	s2 := pool.New(5, 3)
	assert.NotNil(t, s2)
	assert.Equal(t, s2.Width, 5)
	assert.Equal(t, s2.Height, 3)

	assert.Equal(t, len(pool.MapPools), 1)
	assert.Equal(t, len(pool.MapPools[5]), 1)
	assert.NotNil(t, pool.MapPools[5][3])
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 0)
	assert.Equal(t, len(pool.MapPools[5][3].used), 2)

	pool.Put(s1)
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 1)
	assert.Equal(t, len(pool.MapPools[5][3].used), 2) // put will not modify used

	s3 := pool.New(5, 3)
	assert.NotNil(t, s3)
	assert.Equal(t, s3.Width, 5)
	assert.Equal(t, s3.Height, 3)
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 0)
	assert.Equal(t, len(pool.MapPools[5][3].used), 3)

	pool.Reset()
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 3)
	assert.Nil(t, pool.MapPools[5][3].used)

	s4 := pool.New2(5, 3, 1)
	assert.NotNil(t, s4)
	assert.Equal(t, s4.Width, 5)
	assert.Equal(t, s4.Height, 3)
	assert.Equal(t, len(pool.MapPools[5][3].Pools), 2)
	assert.Equal(t, len(pool.MapPools[5][3].used), 1)
	for x := 0; x < s4.Width; x++ {
		for y := 0; y < s4.Height; y++ {
			assert.Equal(t, 1, s4.Arr[x][y])
		}
	}

	s5 := pool.New2(5, 4, 1)
	assert.NotNil(t, s5)
	assert.Equal(t, s5.Width, 5)
	assert.Equal(t, s5.Height, 4)
	assert.Equal(t, len(pool.MapPools[5]), 2)
	assert.Equal(t, len(pool.MapPools[5][4].Pools), 0)
	assert.Equal(t, len(pool.MapPools[5][4].used), 1)
	for x := 0; x < s5.Width; x++ {
		for y := 0; y < s5.Height; y++ {
			assert.Equal(t, 1, s5.Arr[x][y])
		}
	}

	pool.Put(s5)

	s6 := pool.New2(5, 4, 2)
	assert.NotNil(t, s6)
	assert.Equal(t, s6.Width, 5)
	assert.Equal(t, s6.Height, 4)
	assert.Equal(t, len(pool.MapPools[5]), 2)
	assert.Equal(t, len(pool.MapPools[5][4].Pools), 0)
	assert.Equal(t, len(pool.MapPools[5][4].used), 2)
	for x := 0; x < s6.Width; x++ {
		for y := 0; y < s6.Height; y++ {
			assert.Equal(t, 2, s6.Arr[x][y])
		}
	}

	// Test that New does not clear the scene
	s6.Arr[0][0] = 99
	pool.Put(s6)
	s7 := pool.New(5, 4)
	assert.Equal(t, 99, s7.Arr[0][0])

	// test newGameScenePool
	gp := newGameScenePool()
	assert.NotNil(t, gp)

	gs1, _ := NewGameScene(1, 1)
	gp.put(gs1)
	assert.Equal(t, 1, len(gp.Pools))

	gs2 := gp.get()
	assert.Equal(t, gs1, gs2)
	assert.Equal(t, 0, len(gp.Pools))

	gs2_1 := gp.get()
	assert.Nil(t, gs2_1)

	gp.inUsed(gs2)
	assert.Equal(t, 1, len(gp.used))

	gp.reset()
	assert.Equal(t, 1, len(gp.Pools))
	assert.Nil(t, gp.used)

	// test new and new2
	pool2 := NewGameScenePoolEx()
	s8 := pool2.new(1, 1)
	assert.NotNil(t, s8)
	s9 := pool2.new2(1, 1, 7)
	assert.NotNil(t, s9)
	assert.Equal(t, 7, s9.Arr[0][0])
}
