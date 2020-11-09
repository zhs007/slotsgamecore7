package sgc7plugin

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicPlugin(t *testing.T) {
	bp := NewBasicPlugin()

	var lstr []int

	for i := 0; i < 1000; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err, "Test_BasicPlugin Random")
		assert.True(t, func() bool {
			return r >= 0 && r < 100
		}(), "Test_BasicPlugin Random range")

		lstr = append(lstr, r)
	}

	lst := bp.GetUsedRngs()
	assert.NotNil(t, lst, "Test_BasicPlugin GetUsedRngs")
	assert.Equal(t, len(lst), 1000, "Test_BasicPlugin GetUsedRngs len")

	for i := 0; i < 1000; i++ {
		assert.Equal(t, lst[i].Value, lstr[i], "Test_BasicPlugin GetUsedRngs value")
	}

	bp.ClearUsedRngs()

	lst1 := bp.GetUsedRngs()
	assert.Nil(t, lst1, "Test_BasicPlugin GetUsedRngs")

	lstcache := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	bp.SetCache(lstcache)

	for i := 0; i < 10; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.Equal(t, r, lstcache[i], "Test_BasicPlugin Random Cache value")
		assert.Equal(t, len(bp.Cache), 9-i, "Test_BasicPlugin Random ClearCache")
	}

	bp.SetCache(lstcache)

	for i := 0; i < 5; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.Equal(t, r, lstcache[i], "Test_BasicPlugin Random Cache value")
		assert.Equal(t, len(bp.Cache), 9-i, "Test_BasicPlugin Random ClearCache")
	}

	bp.ClearCache()
	assert.Equal(t, len(bp.Cache), 0, "Test_BasicPlugin Random ClearCache")

	var ip IPlugin
	ip = bp
	assert.NotNil(t, ip, "Test_BasicPlugin IPlugin")

	t.Logf("Test_BasicPlugin OK")
}

func Test_BasicPlugin2(t *testing.T) {
	bp := NewBasicPlugin()

	var lstr []int

	for i := 0; i < 10; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.True(t, func() bool {
			return r >= 0 && r < 100
		}())

		lstr = append(lstr, r)
	}

	bp.TagUsedRngs()
	// assert.Equal(t, tag0, 10)

	for i := 0; i < 10; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.True(t, func() bool {
			return r >= 0 && r < 100
		}())

		lstr = append(lstr, r)
	}

	usedrngs := bp.GetUsedRngs()
	assert.Equal(t, len(usedrngs), 20)

	for i := 0; i < len(usedrngs); i++ {
		assert.Equal(t, usedrngs[i].Value, lstr[i])
	}

	err := bp.RollbackUsedRngs()
	assert.NoError(t, err)

	usedrngs = bp.GetUsedRngs()
	assert.Equal(t, len(usedrngs), 10)

	for i := 0; i < len(usedrngs); i++ {
		assert.Equal(t, usedrngs[i].Value, lstr[i])
	}

	bp.ClearUsedRngs()
	err = bp.RollbackUsedRngs()
	assert.Error(t, err)

	t.Logf("Test_BasicPlugin2 OK")
}
