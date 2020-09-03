package sgc7plugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicPlugin(t *testing.T) {
	bp := NewBasicPlugin()

	var lstr []int

	for i := 0; i < 1000; i++ {
		r, err := bp.Random(100)
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

	var ip IPlugin
	ip = bp
	assert.NotNil(t, ip, "Test_BasicPlugin IPlugin")

	t.Logf("Test_BasicPlugin OK")
}
