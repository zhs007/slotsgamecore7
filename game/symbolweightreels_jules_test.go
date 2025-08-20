package sgc7game

import (
	"testing"

	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	"github.com/stretchr/testify/assert"
)

func Test_LoadSymbolWeightReels5JSON_Jules(t *testing.T) {
	swr, err := LoadSymbolWeightReels5JSON("./testdata/symbolweightreels5.json")
	assert.NoError(t, err)
	assert.NotNil(t, swr)

	assert.Equal(t, swr.Width, 5)
	assert.Equal(t, len(swr.Sets), 1)
	assert.Equal(t, len(swr.Sets[0].Arr), 2)

	assert.Equal(t, len(swr.Sets[0].Arr[0].Reels), 5)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].MaxWeights, 30)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[0], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[1], 2)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[0], 10)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[1], 20)

	assert.Equal(t, len(swr.Sets[0].Arr[1].Reels), 5)
	assert.Equal(t, swr.Sets[0].Arr[1].Reels[0].MaxWeights, 30)
	assert.Equal(t, swr.Sets[0].Arr[1].Reels[0].Symbols[0], 3)
	assert.Equal(t, swr.Sets[0].Arr[1].Reels[0].Weights[0], 30)

	// test invalid file
	_, err = LoadSymbolWeightReels5JSON("./testdata/invalidfile.json")
	assert.Error(t, err)

	// test invalid json
	_, err = LoadSymbolWeightReels5JSON("./testdata/reels3_invalid.json")
	assert.Error(t, err)

	// test empty json
	swr2, err := LoadSymbolWeightReels5JSON("./testdata/empty.json")
	assert.NoError(t, err)
	assert.Nil(t, swr2)
}

func Test_SymbolWeightReels_RandomScene_Jules(t *testing.T) {
	swr, err := LoadSymbolWeightReels5JSON("./testdata/symbolweightreels5.json")
	assert.NoError(t, err)
	assert.NotNil(t, swr)

	plugin := sgc7plugin.NewMockPlugin()
	// weights are [10, 20], so 0-9 -> index 0, 10-29 -> index 1
	plugin.Cache = []int{0, 10, 0, 10, 0, 10, 0, 10, 0, 10, 0, 10, 0, 10, 0}

	gs := &GameScene{
		Width:  5,
		Height: 3,
		Arr:    make([][]int, 5),
	}
	for i := 0; i < 5; i++ {
		gs.Arr[i] = make([]int, 3)
	}

	err = swr.RandomScene(gs, plugin, 0, 0, true)
	assert.NoError(t, err)

	assert.Equal(t, 1, gs.Arr[0][0])
	assert.Equal(t, 2, gs.Arr[0][1])
	assert.Equal(t, 1, gs.Arr[0][2])
	assert.Equal(t, 2, gs.Arr[1][0])
	assert.Equal(t, 1, gs.Arr[1][1])

	// test nocheck = false
	gs.Arr[0][0] = -1
	gs.Arr[0][1] = 100
	gs.Arr[0][2] = -1
	plugin.Cache = []int{0, 10}
	err = swr.RandomScene(gs, plugin, 0, 0, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, gs.Arr[0][0])
	assert.Equal(t, 100, gs.Arr[0][1])
	assert.Equal(t, 2, gs.Arr[0][2])

	// test invalid width
	gs.Width = 4
	err = swr.RandomScene(gs, plugin, 0, 0, true)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightWidthReels)
	gs.Width = 5

	// test invalid settype1
	err = swr.RandomScene(gs, plugin, 1, 0, true)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType1)

	err = swr.RandomScene(gs, plugin, -1, 0, true)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType1)

	// test invalid settype2
	err = swr.RandomScene(gs, plugin, 0, 2, true)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType2)

	err = swr.RandomScene(gs, plugin, 0, -1, true)
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType2)

	// test insData5 error
	swr2 := &SymbolWeightReels{}
	err = swr2.insData5(symbolWeightReels{SetType1: 0})
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType1)

	err = swr2.insData5(symbolWeightReels{SetType1: 1, SetType2: 0})
	assert.Error(t, err)
	assert.Equal(t, err, ErrInvalidSymbolWeightReelsSetType2)
}
