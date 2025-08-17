package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestPayTables - new PayTables for test
func newTestPayTablesForJules(t *testing.T) *PayTables {
	pt := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	pt.MapPay[10] = []int{0, 0, 10, 20, 50}
	pt.MapSymbols["SC"] = 10

	return pt
}

// isScatterTest - is a symbol a scatter
func isScatterTestForJules(scatter int, cursymbol int) bool {
	return cursymbol == scatter
}

func Test_CalcScatter_Jules(t *testing.T) {
	pt := newTestPayTablesForJules(t)
	scatter := 10
	bet := 1
	coins := 1

	t.Run("no scatter", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
			{1, 2, 3},
			{4, 5, 6},
		})
		assert.NoError(t, err)

		result := CalcScatter(scene, pt, scatter, bet, coins, isScatterTestForJules)
		assert.Nil(t, result)
	})

	t.Run("3 scatters", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{4, 5, 6},
		})
		assert.NoError(t, err)

		result := CalcScatter(scene, pt, scatter, bet, coins, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.Symbol)
		assert.Equal(t, int(RTScatter), int(result.Type))
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 10, result.CoinWin)
		assert.Equal(t, 10, result.CashWin)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2}, result.Pos)
	})

	t.Run("6 scatters", func(t *testing.T) {
		// CalcScatter nums is capped at len(scene.Arr) which is 5
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 10, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter(scene, pt, scatter, bet, coins, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.Symbol)
		assert.Equal(t, int(RTScatter), int(result.Type))
		assert.Equal(t, 50, result.Mul)
		assert.Equal(t, 50, result.CoinWin)
		assert.Equal(t, 50, result.CashWin)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2, 3, 1, 4, 0, 4, 2}, result.Pos)
	})
}

func Test_CalcScatter2_Jules(t *testing.T) {
	pt := newTestPayTablesForJules(t)
	scatter := 10
	bet := 1
	coins := 1

	pt.MapPay[10] = []int{0, 0, 10, 20, 50, 100}

	t.Run("6 scatters", func(t *testing.T) {
		// CalcScatter2 nums is capped at len(pt.MapPay[scatter]) which is 6
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 10, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter2(scene, pt, scatter, bet, coins, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 10, result.Symbol)
		assert.Equal(t, int(RTScatter), int(result.Type))
		assert.Equal(t, 100, result.Mul)
		assert.Equal(t, 100, result.CoinWin)
		assert.Equal(t, 100, result.CashWin)
		assert.Equal(t, 6, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2, 3, 1, 4, 0, 4, 2}, result.Pos)
	})
}

func Test_CalcScatter3_Jules(t *testing.T) {
	pt := newTestPayTablesForJules(t)
	scatter := 10
	bet := 1
	coins := 1

	t.Run("only one on reel", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter3(scene, pt, scatter, bet, coins, isScatterTestForJules, true)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2, 4, 0}, result.Pos)
	})

	t.Run("all on reel", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter3(scene, pt, scatter, bet, coins, isScatterTestForJules, false)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2, 4, 0, 4, 2}, result.Pos)
	})
}

func Test_CalcScatter4_Jules(t *testing.T) {
	pt := newTestPayTablesForJules(t)
	scatter := 10
	bet := 1

	t.Run("only one on reel", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter4(scene, pt, scatter, bet, isScatterTestForJules, true)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 2, 2, 4, 0}, result.Pos)
		assert.Equal(t, 20, result.CoinWin)
		assert.Equal(t, 20, result.CashWin)
	})

	t.Run("all on reel", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter4(scene, pt, scatter, bet, isScatterTestForJules, false)
		assert.NotNil(t, result)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, 50, result.CoinWin)
		assert.Equal(t, 50, result.CashWin)
	})
}

func Test_CalcScatter5_Jules(t *testing.T) {
	pt := newTestPayTablesForJules(t)
	scatter := 10
	bet := 1

	t.Run("height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter5(scene, pt, scatter, bet, isScatterTestForJules, false, 2, false)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, []int{0, 0, 1, 1, 4, 0}, result.Pos)
	})

	t.Run("reversal height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatter5(scene, pt, scatter, bet, isScatterTestForJules, false, 2, true)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, []int{1, 1, 2, 2, 4, 2}, result.Pos)
	})
}

func Test_CalcScatterEx_Jules(t *testing.T) {
	scatter := 10
	nums := 3

	t.Run("3 scatters", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{4, 5, 6},
		})
		assert.NoError(t, err)

		result := CalcScatterEx(scene, scatter, nums, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
	})
}

func Test_CalcScatterEx2_Jules(t *testing.T) {
	scatter := 10
	nums := 2

	t.Run("height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatterEx2(scene, scatter, nums, isScatterTestForJules, 2, false)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
	})

	t.Run("reversal height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatterEx2(scene, scatter, nums, isScatterTestForJules, 2, true)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
	})
}

func Test_CalcReelScatterEx_Jules(t *testing.T) {
	scatter := 10
	nums := 3

	t.Run("4 scatters", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 3},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcReelScatterEx(scene, scatter, nums, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
	})
}

func Test_CalcReelScatterEx2_Jules(t *testing.T) {
	scatter := 10
	nums := 2

	t.Run("height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 10},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcReelScatterEx2(scene, scatter, nums, isScatterTestForJules, 2, false)
		assert.NotNil(t, result)
		assert.Equal(t, 3, result.SymbolNums)
	})

	t.Run("reversal height 2", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 2, 10},
			{4, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcReelScatterEx2(scene, scatter, nums, isScatterTestForJules, 2, true)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
	})
}

func Test_CountScatterInArea_Jules(t *testing.T) {
	scatter := 10
	nums := 2

	isInArea := func(x, y int) bool {
		return x < 2 && y < 2
	}

	t.Run("2 scatters in area", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 10, 10},
			{10, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CountScatterInArea(scene, scatter, nums, isInArea, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
	})
}

func Test_CalcScatterOnReels_Jules(t *testing.T) {
	scatter := 10
	nums := 3

	t.Run("4 scatters on reels", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{10, 10, 10},
			{10, 10, 6},
			{7, 8, 10},
			{1, 2, 3},
			{10, 5, 10},
		})
		assert.NoError(t, err)

		result := CalcScatterOnReels(scene, scatter, nums, isScatterTestForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 4, result.SymbolNums)
	})
}

func newTestPayTablesForLineTests(t *testing.T) *PayTables {
	pt := &PayTables{
		MapPay:     make(map[int][]int),
		MapSymbols: make(map[string]int),
	}

	pt.MapPay[0] = []int{0, 10, 20, 50, 100}
	pt.MapSymbols["WILD"] = 0

	pt.MapPay[1] = []int{0, 5, 10, 20, 40}
	pt.MapSymbols["A"] = 1

	pt.MapPay[2] = []int{0, 3, 6, 12, 24}
	pt.MapSymbols["B"] = 2

	return pt
}

func isValidSymbolForJules(cursymbol int) bool {
	return cursymbol >= 0
}

func isWildForJules(cursymbol int) bool {
	return cursymbol == 0
}

func isSameSymbolForJules(cursymbol int, startsymbol int) bool {
	if isWildForJules(startsymbol) {
		return isWildForJules(cursymbol)
	}

	return cursymbol == startsymbol || isWildForJules(cursymbol)
}

func getSymbolForJules(cursymbol int) int {
	return cursymbol
}

func calcOtherMulForJules(scene *GameScene, result *Result) int {
	if result.Wilds > 0 {
		return 2
	}

	return 1
}

func calcOtherMulExForJules(scene *GameScene, symbol int, pos []int) int {
	hasWild := false
	for i := 0; i < len(pos)/2; i++ {
		if scene.Arr[pos[i*2]][pos[i*2+1]] == 0 {
			hasWild = true
			break
		}
	}

	if hasWild {
		return 2
	}

	return 1
}

func Test_CalcLine_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1
	line := []int{0, 0, 0, 0, 0} // A straight line on the first row

	t.Run("3 of a kind", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLine(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 0, result.Wilds)
		assert.Equal(t, []int{0, 0, 1, 0, 2, 0}, result.Pos)
	})

	t.Run("3 of a kind with wild", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{0, 4, 5}, // Wild
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLine(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
		assert.Equal(t, []int{0, 0, 1, 0, 2, 0}, result.Pos)
	})

	t.Run("starts with wild", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3}, // Wild
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLine(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol) // Should be symbol A
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul) // Pay for 3xA
		assert.Equal(t, 1, result.Wilds)
		assert.Equal(t, []int{0, 0, 1, 0, 2, 0}, result.Pos)
	})

	t.Run("starts with wild, wild pays more", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3}, // Wild
			{0, 4, 5},
			{2, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLine(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Symbol) // Should be symbol 2
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, 24, result.Mul) // Pay for 5xSymbol2
		assert.Equal(t, 2, result.Wilds)
	})

	t.Run("no win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{2, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{1, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLine(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.Nil(t, result)
	})
}

func Test_CalcLineEx_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1
	line := []int{0, 0, 0, 0, 0}

	t.Run("3 of a kind with wild and othermul", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{0, 4, 5}, // Wild
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineEx(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, calcOtherMulForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
		assert.Equal(t, 2, result.OtherMul)
		assert.Equal(t, 20, result.CoinWin) // 10 * 2
		assert.Equal(t, 20, result.CashWin) // 10 * 2 * 1
	})
}

func Test_CalcLineRL_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1
	line := []int{0, 0, 0, 0, 0}

	t.Run("3 of a kind RL", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{2, 2, 3},
			{2, 4, 5},
			{1, 6, 7},
			{1, 9, 1},
			{1, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineRL(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
	})

	t.Run("3 of a kind RL with wild", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{2, 2, 3},
			{2, 4, 5},
			{1, 6, 7},
			{0, 9, 1}, // Wild
			{1, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineRL(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
	})

	t.Run("starts with wild RL", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{2, 2, 3},
			{2, 4, 5},
			{1, 6, 7},
			{1, 9, 1},
			{0, 3, 4}, // Wild
		})
		assert.NoError(t, err)

		result := CalcLineRL(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
	})

	t.Run("starts with wild, wild pays more RL", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{2, 2, 3},
			{2, 4, 5},
			{2, 6, 7},
			{0, 9, 1}, // Wild
			{0, 3, 4}, // Wild
		})
		assert.NoError(t, err)

		result := CalcLineRL(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Symbol)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, 24, result.Mul)
		assert.Equal(t, 2, result.Wilds)
	})
}

func Test_CalcLineRLEx_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1
	line := []int{0, 0, 0, 0, 0}

	t.Run("3 of a kind RL with wild and othermul", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{2, 2, 3},
			{2, 4, 5},
			{1, 6, 7},
			{0, 9, 1}, // Wild
			{1, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineRLEx(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, calcOtherMulForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
		assert.Equal(t, 2, result.OtherMul)
		assert.Equal(t, 20, result.CoinWin)
		assert.Equal(t, 20, result.CashWin)
	})
}

func Test_CalcLineOtherMul_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1
	line := []int{0, 0, 0, 0, 0}

	t.Run("3 of a kind with wild and othermul", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{0, 4, 5}, // Wild
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineOtherMul(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, calcOtherMulExForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 1, result.Symbol)
		assert.Equal(t, 3, result.SymbolNums)
		assert.Equal(t, 10, result.Mul)
		assert.Equal(t, 1, result.Wilds)
		assert.Equal(t, 2, result.OtherMul)
		assert.Equal(t, 20, result.CoinWin)
		assert.Equal(t, 20, result.CashWin)
	})

	t.Run("starts with wild, wild pays more with othermul", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3}, // Wild
			{0, 4, 5},
			{0, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		result := CalcLineOtherMul(scene, pt, line, bet, isValidSymbolForJules, isWildForJules, isSameSymbolForJules, calcOtherMulExForJules, getSymbolForJules)
		assert.NotNil(t, result)
		assert.Equal(t, 2, result.Symbol)
		assert.Equal(t, 5, result.SymbolNums)
		assert.Equal(t, 24, result.Mul)
		assert.Equal(t, 3, result.Wilds)
		assert.Equal(t, 2, result.OtherMul)
		assert.Equal(t, 48, result.CoinWin)
	})
}

func isValidSymbolExForJules(pt *PayTables, cursymbol int, scene *GameScene, x, y int) bool {
	if cursymbol < 0 {
		return false
	}
	_, ok := pt.MapPay[cursymbol]
	return ok
}

func getMultiForJules(x, y int) int {
	return 1
}

func Test_CalcFullLineEx_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1

	t.Run("simple ways win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CalcFullLineEx(scene, pt, bet, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}

func Test_CalcFullLineExWithMulti_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1

	t.Run("ways win with multi", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CalcFullLineExWithMulti(scene, pt, bet, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules, getMultiForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}

func Test_CheckWays_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	minnum := 2

	t.Run("ways win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CheckWays(scene, minnum, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}

func Test_CalcFullLineEx2_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1

	t.Run("simple ways win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CalcFullLineEx2(scene, pt, bet, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})

	t.Run("starts with wild", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{0, 2, 3}, // Wild
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CalcFullLineEx2(scene, pt, bet, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}

func Test_CalcFullLine_Jules(t *testing.T) {
	pt := newTestPayTablesForLineTests(t)
	bet := 1

	t.Run("simple ways win", func(t *testing.T) {
		scene, err := NewGameSceneWithArr2([][]int{
			{1, 2, 3},
			{1, 4, 5},
			{1, 6, 7},
			{2, 9, 1},
			{2, 3, 4},
		})
		assert.NoError(t, err)

		results := CalcFullLine(scene, pt, bet, func(cursymbol int, scene *GameScene, x, y int) bool {
			return isValidSymbolExForJules(pt, cursymbol, scene, x, y)
		}, isWildForJules, isSameSymbolForJules)
		assert.NotNil(t, results)
		assert.Equal(t, 1, len(results))
	})
}
