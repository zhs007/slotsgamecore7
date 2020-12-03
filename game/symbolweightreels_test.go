package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadSymbolWeightReels5JSON(t *testing.T) {
	swr, err := LoadSymbolWeightReels5JSON("../unittestdata/symbolweightreels.json")
	assert.NoError(t, err)
	assert.NotNil(t, swr)

	assert.Equal(t, len(swr.Sets), 2)
	assert.Equal(t, len(swr.Sets[0].Arr), 4)
	assert.Equal(t, len(swr.Sets[1].Arr), 4)

	assert.Equal(t, len(swr.Sets[0].Arr[0].Reels), 5)
	assert.Equal(t, len(swr.Sets[0].Arr[0].Reels[0].Symbols), 11)
	assert.Equal(t, len(swr.Sets[0].Arr[0].Reels[0].Weights), 11)

	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[0], 0)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[1], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[2], 2)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[3], 3)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[4], 4)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[5], 5)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[6], 6)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[7], 7)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[8], 8)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[9], 9)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Symbols[10], 10)

	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[0], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[1], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[2], 6)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[3], 3)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[4], 3)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[5], 8)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[6], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[7], 9)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[8], 10)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[9], 1)
	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].Weights[10], 1)

	assert.Equal(t, swr.Sets[0].Arr[0].Reels[0].MaxWeights, 44)

	assert.Equal(t, len(swr.Sets[0].Arr[3].Reels), 5)
	assert.Equal(t, len(swr.Sets[0].Arr[3].Reels[0].Symbols), 11)
	assert.Equal(t, len(swr.Sets[0].Arr[3].Reels[0].Weights), 11)

	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[0], 0)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[1], 1)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[2], 2)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[3], 3)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[4], 4)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[5], 5)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[6], 6)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[7], 7)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[8], 8)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[9], 9)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Symbols[10], 10)

	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[0], 5)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[1], 1)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[2], 1)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[3], 6)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[4], 12)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[5], 3)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[6], 15)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[7], 7)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[8], 2)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[9], 10)
	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].Weights[10], 1)

	assert.Equal(t, swr.Sets[0].Arr[3].Reels[2].MaxWeights, 63)

	t.Logf("Test_LoadSymbolWeightReels5JSON OK")
}
