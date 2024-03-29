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

	ms := NewSymbolsMapping()
	ms.Add(12, 1)

	ssws, err := AnalyzeReelsWithLine(paytables, reels, []SymbolType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []SymbolType{0}, ms, 10, 10)
	assert.NoError(t, err)
	assert.NotNil(t, ssws)

	err = ssws.SaveExcel("../unittestdata/symbolswinsstats.xlsx", []SymbolsWinsFileMode{SWFModeRTP, SWFModeWinsNum, SWFModeWinsNumPer})
	assert.NoError(t, err)

	t.Logf("Test_AnalyzeReelsWithLine OK")
}

func Test_AnalyzeReelsScatter(t *testing.T) {
	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	paytables, err := sgc7game.LoadPaytablesFromExcel("../unittestdata/paytables.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, paytables)

	ms := NewSymbolsMapping()
	ms.Add(12, 1)

	ssws, err := AnalyzeReelsScatter(paytables, reels, []SymbolType{11}, ms, 3)
	assert.NoError(t, err)
	assert.NotNil(t, ssws)

	err = ssws.SaveExcel("../unittestdata/scatterswinsstats.xlsx", []SymbolsWinsFileMode{SWFModeRTP, SWFModeWinsNum, SWFModeWinsNumPer})
	assert.NoError(t, err)

	t.Logf("Test_AnalyzeReelsScatter OK")
}

func Test_MergeSSWS(t *testing.T) {
	reels, err := sgc7game.LoadReelsFromExcel("../unittestdata/reels.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, reels)

	paytables, err := sgc7game.LoadPaytablesFromExcel("../unittestdata/paytables.xlsx")
	assert.NoError(t, err)
	assert.NotNil(t, paytables)

	ms := NewSymbolsMapping()
	ms.Add(12, 1)

	ssws, err := AnalyzeReelsWithLine(paytables, reels, []SymbolType{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []SymbolType{0}, ms, 10, 10)
	assert.NoError(t, err)
	assert.NotNil(t, ssws)

	ssws1, err := AnalyzeReelsScatter(paytables, reels, []SymbolType{11}, ms, 3)
	assert.NoError(t, err)
	assert.NotNil(t, ssws1)

	ssws.Merge(ssws1)

	err = ssws.SaveExcel("../unittestdata/mergewinsstats.xlsx", []SymbolsWinsFileMode{SWFModeRTP, SWFModeWinsNum, SWFModeWinsNumPer})
	assert.NoError(t, err)

	t.Logf("Test_AnalyzeReelsScatter OK")
}
