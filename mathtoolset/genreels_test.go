package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenReels(t *testing.T) {
	rss, err := LoadReelsStats("../unittestdata/genreelssrc.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	assert.Equal(t, len(rss.Reels), 5)
	assert.Equal(t, len(rss.Reels[0].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[1].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[2].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[3].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[4].MapSymbols), 13)

	rd, err := GenReels(rss, 2, 100)
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	rd.SaveExcel("../unittestdata/genreels.xlsx")

	t.Logf("Test_GenReels OK")
}

func Test_GenReelsMainSymbolsDistance(t *testing.T) {
	rss, err := LoadReelsStats("../unittestdata/genreelssrc.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	assert.Equal(t, len(rss.Reels), 5)
	assert.Equal(t, len(rss.Reels[0].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[1].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[2].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[3].MapSymbols), 13)
	assert.Equal(t, len(rss.Reels[4].MapSymbols), 13)

	rd, err := GenReelsMainSymbolsDistance(rss, []SymbolType{0, 1, 2, 3}, 2, 100)
	assert.NoError(t, err)
	assert.NotNil(t, rd)

	rd.SaveExcel("../unittestdata/genreelsmainsymbolsdistance.xlsx")

	t.Logf("Test_GenReels OK")
}
