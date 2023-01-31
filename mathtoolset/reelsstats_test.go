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

	rss, err := BuildReelsStats(reels, nil)
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	assert.Equal(t, len(rss.Reels), 5)
	assert.Equal(t, len(rss.Reels[0].MapSymbols), 11)
	assert.Equal(t, len(rss.Reels[1].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[2].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[3].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[4].MapSymbols), 12)

	err = rss.SaveExcel("../unittestdata/reelsstats.xlsx")
	assert.NoError(t, err)

	t.Logf("Test_BuildReelsStats OK")
}

func Test_BuildReelsStatsMapping(t *testing.T) {

	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	ms := NewSymbolsMapping()
	ms.Add(12, 1)

	rss, err := BuildReelsStats(reels, ms)
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	assert.Equal(t, len(rss.Reels), 5)
	assert.Equal(t, len(rss.Reels[0].MapSymbols), 10)
	assert.Equal(t, len(rss.Reels[1].MapSymbols), 12)
	assert.Equal(t, len(rss.Reels[2].MapSymbols), 12)
	assert.Equal(t, len(rss.Reels[3].MapSymbols), 12)
	assert.Equal(t, len(rss.Reels[4].MapSymbols), 11)

	err = rss.SaveExcel("../unittestdata/reelsstatsmapping.xlsx")
	assert.NoError(t, err)

	t.Logf("Test_BuildReelsStatsMapping OK")
}
