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

	rss2, err := LoadReelsStats("../unittestdata/reelsstats.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, rss2)

	assert.Equal(t, len(rss2.Reels), 5)
	assert.Equal(t, len(rss2.Reels[0].MapSymbols), 13)
	assert.Equal(t, len(rss2.Reels[1].MapSymbols), 13)
	assert.Equal(t, len(rss2.Reels[2].MapSymbols), 13)
	assert.Equal(t, len(rss2.Reels[3].MapSymbols), 13)
	assert.Equal(t, len(rss2.Reels[4].MapSymbols), 13)

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

func Test_ReelStats_GetCanAddSymbols(t *testing.T) {
	rs := NewReelStats()

	rs.AddSymbol(0, 1)
	rs.AddSymbol(1, 1)
	rs.AddSymbol(2, 1)
	rs.AddSymbol(3, 1)
	rs.AddSymbol(4, 1)
	rs.AddSymbol(5, 1)
	rs.AddSymbol(6, 1)
	rs.AddSymbol(7, 1)
	rs.AddSymbol(8, 1)
	rs.AddSymbol(9, 1)

	lst := rs.GetCanAddSymbols([]SymbolType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, len(lst), 1)

	t.Logf("Test_ReelStats_GetCanAddSymbols OK")
}

func Test_ReelStats_GetCanAddSymbols2(t *testing.T) {
	rs := NewReelStats()

	rs.AddSymbol(0, 1)
	rs.AddSymbol(1, 2)
	rs.AddSymbol(2, 3)
	rs.AddSymbol(3, 4)
	rs.AddSymbol(4, 5)
	rs.AddSymbol(5, 6)
	rs.AddSymbol(6, 7)
	rs.AddSymbol(7, 8)
	rs.AddSymbol(8, 8)
	rs.AddSymbol(9, 8)

	lst := rs.GetCanAddSymbols([]SymbolType{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	assert.Equal(t, len(lst), 8)

	t.Logf("Test_ReelStats_GetCanAddSymbols2 OK")
}
