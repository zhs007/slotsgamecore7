package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewPlayResult_Jules(t *testing.T) {
	pr := NewPlayResult("bg", 0, -1, "NORMAL")
	assert.NotNil(t, pr)
	assert.Equal(t, "bg", pr.CurGameMod)
	assert.Equal(t, 0, pr.CurIndex)
	assert.Equal(t, -1, pr.ParentIndex)
	assert.Equal(t, "NORMAL", pr.ModType)

	t.Logf("Test_NewPlayResult_Jules OK")
}

func Test_GetPlayResultCurIndex_Jules(t *testing.T) {
	prs := []*PlayResult{}
	assert.Equal(t, 0, GetPlayResultCurIndex(prs))

	prs = append(prs, &PlayResult{CurIndex: 0})
	assert.Equal(t, 1, GetPlayResultCurIndex(prs))

	prs = append(prs, &PlayResult{CurIndex: 1})
	assert.Equal(t, 2, GetPlayResultCurIndex(prs))

	t.Logf("Test_GetPlayResultCurIndex_Jules OK")
}

func Test_PlayResult2JSON_Jules(t *testing.T) {
	pr := &PlayResult{
		CurGameMod: "bg",
		CoinWin:    100,
	}

	json, err := PlayResult2JSON(pr)
	assert.NoError(t, err)

	pr2, err := JSON2PlayResult(json, &PlayResult{})
	assert.NoError(t, err)

	assert.Equal(t, pr.CurGameMod, pr2.CurGameMod)
	// Note: CoinWin is not serialized
	assert.Equal(t, 0, pr2.CoinWin)

	t.Logf("Test_PlayResult2JSON_Jules OK")
}

func Test_CountEndingSymbols_Jules(t *testing.T) {
	pr := &PlayResult{}
	assert.Nil(t, pr.CountEndingSymbols([]int{0, 1}))

	gs := &GameScene{
		Arr: [][]int{
			{0, 1, 2},
			{1, 2, 0},
			{2, 0, 1},
		},
	}
	pr.Scenes = append(pr.Scenes, gs)

	counts := pr.CountEndingSymbols([]int{0, 1})
	assert.Equal(t, 2, len(counts))
	assert.Equal(t, 3, counts[0])
	assert.Equal(t, 3, counts[1])

	t.Logf("Test_CountEndingSymbols_Jules OK")
}
