package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_BuildWaysReelsStats(t *testing.T) {

	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	wrss := BuildWaysReelsStats(reels, 3)
	assert.NotNil(t, wrss)

	assert.Equal(t, len(wrss.Reels), 5)
	assert.Equal(t, wrss.Reels[0].GetNumWithSymbolNumInWindow(1, 1), 9)
	assert.Equal(t, wrss.Reels[0].GetNumWithSymbolNumInWindow(12, 1), 2)
	assert.Equal(t, wrss.Reels[0].GetNumWithSymbolNumInWindow(12, 2), 2)
	assert.Equal(t, wrss.Reels[0].GetNumWithSymbolNumInWindow(12, 3), 3)

	wrss.SaveExcel("../unittestdata/wrss.xlsx")

	t.Logf("Test_BuildWaysReelsStats OK")
}
