package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadReelsStats2(t *testing.T) {

	reels, err := LoadReelsStats2("../unittestdata/reelsstats2.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	t.Logf("Test_LoadReelsStats2 OK")
}
