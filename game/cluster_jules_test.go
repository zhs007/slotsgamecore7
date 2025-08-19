package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock functions for cluster tests
func isValidSymbol(s int) bool {
	return s >= 0 && !isWild(s)
}

func isWild(s int) bool {
	return s == 1 // Assuming 1 is the wild symbol
}

func isSameSymbol(s1, s2 int) bool {
	if isWild(s1) {
		return true
	}

	return s1 == s2
}

func getSymbol(s int) int {
	return s
}

func Test_CalcClusterResult_Jules(t *testing.T) {
	scene := &GameScene{
		Width:  3,
		Height: 3,
		Arr: [][]int{
			{2, 2, 3},
			{2, 1, 3},
			{3, 3, 3},
		},
	}
	pt := &PayTables{
		MapPay: map[int][]int{
			2: {0, 0, 10, 20}, // 3 symbols = 10, 4 symbols = 20
			3: {0, 0, 0, 0, 30}, // 5 symbols = 30
		},
	}
	bet := 1

	results, err := CalcClusterResult(scene, pt, bet, isValidSymbol, isWild, isSameSymbol, getSymbol)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 2, len(results))

	// Cluster of 2s (3 symbols) + 1 wild = 4 symbols
	assert.Equal(t, 2, results[0].Symbol)
	assert.Equal(t, 4, results[0].SymbolNums)
	assert.Equal(t, 20, results[0].CoinWin)

	// Cluster of 3s (5 symbols)
	assert.Equal(t, 3, results[1].Symbol)
	assert.Equal(t, 5, results[1].SymbolNums)
	assert.Equal(t, 30, results[1].CoinWin)
}

func Test_CalcClusterSymbol_Jules(t *testing.T) {
	scene := &GameScene{
		Width:  3,
		Height: 3,
		Arr: [][]int{
			{2, 2, 3},
			{2, 1, 3},
			{3, 3, 3},
		},
	}

	pos := calcClusterSymbol(scene, 0, 0, 2, []int{}, isSameSymbol)
	assert.Equal(t, 8, len(pos)) // 4 positions * 2 coordinates
}
