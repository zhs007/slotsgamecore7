package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_RandWithWeights(t *testing.T) {
	bp := sgc7plugin.NewBasicPlugin()

	c, err := RandWithWeights(bp, 0, nil)
	assert.EqualError(t, err, ErrInvalidWeights.Error())
	assert.Equal(t, c, -1)

	c, err = RandWithWeights(bp, 100, nil)
	assert.EqualError(t, err, ErrInvalidWeights.Error())
	assert.Equal(t, c, -1)

	c, err = RandWithWeights(bp, 0, []int{1, 2, 3})
	assert.EqualError(t, err, ErrInvalidWeights.Error())
	assert.Equal(t, c, -1)

	bp.SetCache([]int{2})
	c, err = RandWithWeights(bp, 6, []int{1, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, c, 1)

	bp.SetCache([]int{7})
	c, err = RandWithWeights(bp, 8, []int{1, 2, 3})
	assert.EqualError(t, err, ErrInvalidWeights.Error())
	assert.Equal(t, c, -1)

	bp.SetCache([]int{2})
	c, err = RandWithWeights(bp, 8, []int{1, 0, 2, 3})
	assert.NoError(t, err)
	assert.Equal(t, c, 2)

	t.Logf("Test_RandWithWeights OK")
}
