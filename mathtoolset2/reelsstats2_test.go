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

func Test_SaveLoadReelsStats2_RoundTrip(t *testing.T) {
	tmpf, err := os.CreateTemp("", "reelsstats2-*.xlsx")
	assert.NoError(t, err)
	tmpname := tmpf.Name()
	tmpf.Close()
	defer os.Remove(tmpname)

	rss := NewReelsStats2(3)
	rss.Symbols = []string{"A", "B_2"}

	rss.Reels[0].MapSymbols["A"] = 2
	rss.Reels[1].MapSymbols["A"] = 3
	rss.Reels[2].MapSymbols["A"] = 1

	rss.Reels[0].MapSymbols["B_2"] = 5
	rss.Reels[1].MapSymbols["B_2"] = 0
	rss.Reels[2].MapSymbols["B_2"] = 2

	rss.Reels[0].TotalSymbolNum = 7
	rss.Reels[1].TotalSymbolNum = 3
	rss.Reels[2].TotalSymbolNum = 3

	err = rss.SaveExcel(tmpname)
	assert.NoError(t, err)

	file, err := os.Open(tmpname)
	assert.NoError(t, err)
	defer file.Close()

	// no debug printing

	loaded, err := LoadReelsStats2(file)
	assert.NoError(t, err)
	assert.NotNil(t, loaded)

	assert.Equal(t, len(rss.Reels), len(loaded.Reels))

	for i := range rss.Reels {
		for _, sym := range rss.Symbols {
			exp := rss.Reels[i].MapSymbols[sym]
			got := loaded.Reels[i].MapSymbols[sym]
			assert.Equal(t, exp, got, "reel %d symbol %s", i, sym)
		}
		assert.Equal(t, rss.Reels[i].TotalSymbolNum, loaded.Reels[i].TotalSymbolNum)
	}
}

func Test_BuildReelsStats2(t *testing.T) {
	reels := [][]string{
		{"A", "B", "A"},
		{"B", "C"},
		{"A", "C", "C"},
	}

	rss, err := BuildReelsStats2(reels)
	assert.NoError(t, err)
	assert.NotNil(t, rss)

	// totals
	assert.Equal(t, 3, rss.Reels[0].TotalSymbolNum)
	assert.Equal(t, 2, rss.Reels[1].TotalSymbolNum)
	assert.Equal(t, 3, rss.Reels[2].TotalSymbolNum)

	// counts
	assert.Equal(t, 2, rss.Reels[0].MapSymbols["A"])
	assert.Equal(t, 1, rss.Reels[0].MapSymbols["B"])

	assert.Equal(t, 1, rss.Reels[1].MapSymbols["B"])
	assert.Equal(t, 1, rss.Reels[1].MapSymbols["C"])

	assert.Equal(t, 1, rss.Reels[2].MapSymbols["A"])
	assert.Equal(t, 2, rss.Reels[2].MapSymbols["C"])

	// symbols set (sorted)
	assert.Equal(t, []string{"A", "B", "C"}, rss.Symbols)
}

func Test_ReelsStats2_FullFlow_FromReels(t *testing.T) {
	reels := [][]string{
		{"X", "Y", "X", "Z"},
		{"Y", "Y", "X"},
	}

	rss, err := BuildReelsStats2(reels)
	assert.NoError(t, err)

	tmpf, err := os.CreateTemp("", "reelsstats2-flow-*.xlsx")
	assert.NoError(t, err)
	tmpname := tmpf.Name()
	tmpf.Close()
	defer os.Remove(tmpname)

	err = rss.SaveExcel(tmpname)
	assert.NoError(t, err)

	file, err := os.Open(tmpname)
	assert.NoError(t, err)
	defer file.Close()

	loaded, err := LoadReelsStats2(file)
	assert.NoError(t, err)

	// compare totals and counts
	assert.Equal(t, len(rss.Reels), len(loaded.Reels))
	for i := range rss.Reels {
		assert.Equal(t, rss.Reels[i].TotalSymbolNum, loaded.Reels[i].TotalSymbolNum)
		for _, s := range rss.Symbols {
			assert.Equal(t, rss.Reels[i].MapSymbols[s], loaded.Reels[i].MapSymbols[s])
		}
	}
}
