package mathtoolset2

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadReelsStats2(t *testing.T) {
	file, err := os.Open("../unittestdata/reelsstats2.xlsx")
	assert.NoError(t, err)

	reels, err := LoadReelsStats2(file)
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	t.Logf("Test_LoadReelsStats2 OK")
}
