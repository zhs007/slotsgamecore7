package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_BuildReelsStats(t *testing.T) {

	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	rss, err := BuildReelsStats(reels)
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	assert.Equal(t, len(rss.Reels), 5)
	assert.Equal(t, len(rss.Reels[0].MapSymbols), 11)
	assert.Equal(t, len(rss.Reels[1].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[2].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[3].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[4].MapSymbols), 12)

	t.Logf("Test_BuildReelsStats OK")
}
