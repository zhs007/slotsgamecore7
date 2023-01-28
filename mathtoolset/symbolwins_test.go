package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_AnalyzeReelsWithLine(t *testing.T) {
	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	paytables, err := sgc7game.LoadPaytablesFromExcel("../unittestdata/paytables.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, paytables)

	ssws, err := AnalyzeReelsWithLine(paytables, reels, []SymbolType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	assert.NoError(t, err)
	assert.NotNil(t, ssws)

	err = ssws.SaveExcel("../unittestdata/symbolswinsstats.xlsx")
	assert.NoError(t, err)

	t.Logf("Test_AnalyzeReelsWithLine OK")
}
