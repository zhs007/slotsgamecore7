package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_CheckLine_Jules002 - test CheckLine
func Test_CheckLine_Jules002(t *testing.T) {
	// isValidSymbol - is it a valid symbol?
	isValidSymbol := func(cursymbol int) bool {
		return cursymbol >= 0
	}

	// isWild - is it a wild symbol?
	isWild := func(cursymbol int) bool {
		return cursymbol == 0
	}

	// isSameSymbol - is it the same symbol? (wild is always the same)
	isSameSymbol := func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol || cursymbol == 0
	}

	// getSymbol - get symbol
	getSymbol := func(cursymbol int) int {
		return cursymbol
	}

	// Test case 1: simple line, meets minnum
	t.Run("simple line, meets minnum", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{8, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 1, 2}
		result := CheckLine(scene, line, 3, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 0, result.Wilds)
	})

	// Test case 2: simple line, does not meet minnum
	t.Run("simple line, does not meet minnum", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{8, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 1, 2}
		result := CheckLine(scene, line, 4, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.Nil(t, result)
	})

	// Test case 3: line with wilds, not starting with wild
	t.Run("line with wilds, not starting with wild", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{0, 4, 5},
			{1, 6, 7},
			{0, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 1, 2}
		result := CheckLine(scene, line, 3, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 1, result.Wilds)
	})

	// Test case 4: all wilds
	t.Run("all wilds", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3},
			{0, 4, 5},
			{0, 6, 7},
			{0, 9, 1},
			{0, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 5, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Symbol)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, 5, result.Wilds)
	})

	// Test case 5: starts with wild, should return symbol B (2)
	t.Run("starts with wild, returns non-wild symbol", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 1, 3}, // W
			{0, 4, 5}, // W
			{2, 6, 7}, // B
			{2, 9, 1}, // B
			{3, 3, 4}, // C
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 4, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Symbol)
		assert.Equal(t, 4, result.SymbolNums)
		assert.Equal(t, 2, result.Wilds)
	})

	// Test case 6: starts with wild, but not enough symbols
	t.Run("starts with wild, not enough symbols", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 1, 3}, // W
			{0, 4, 5}, // W
			{2, 6, 7}, // B
			{3, 9, 1}, // C
			{4, 3, 4}, // D
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 4, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.Nil(t, result)
	})

	// Test case 7: invalid symbol at start
	t.Run("invalid symbol at start", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{-1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{8, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 1, 2}
		result := CheckLine(scene, line, 3, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.Nil(t, result)
	})

	// Test case 8: invalid symbol in the middle, breaking the line
	t.Run("invalid symbol in middle", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{-1, 6, 7},
			{1, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 3, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.Nil(t, result)
	})

	// Test case 9: W W W B C, should be 4 B
	t.Run("W W W B C", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 1, 3}, // W
			{0, 4, 5}, // W
			{0, 6, 7}, // W
			{2, 9, 1}, // B
			{3, 3, 4}, // C
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 4, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Symbol)
		assert.Equal(t, 4, result.SymbolNums)
		assert.Equal(t, 3, result.Wilds)
	})

	// Test case 10: only wild win, W W W, followed by invalid symbol
	t.Run("only wild win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 1, 3}, // W
			{0, 4, 5}, // W
			{0, 6, 7}, // W
			{-1, 9, 1}, // Invalid
			{5, 3, 4}, // E
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 3, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 0, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 3, result.Wilds)
	})

	// Test case 11: Wilds after non-wild sequence
	t.Run("wilds after non-wilds", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{0, 6, 7},
			{0, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 4, isValidSymbol, isWild, isSameSymbol, getSymbol)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 4, result.SymbolNums)
		assert.Equal(t, 2, result.Wilds)
	})

	// Test case 12: No win because first symbol is not valid
	t.Run("no win, first symbol invalid", func(t *testing.T) {
		isValid := func(s int) bool { return s > 0 }
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3},
			{0, 4, 5},
			{0, 6, 7},
			{0, 9, 1},
			{0, 3, 4},
		})
		assert.NoError(t, err)
		line := []int{0, 0, 0, 0, 0}
		result := CheckLine(scene, line, 5, isValid, isWild, isSameSymbol, getSymbol)
		assert.Nil(t, result)
	})
}
