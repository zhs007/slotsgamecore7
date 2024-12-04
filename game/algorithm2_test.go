package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcLine2(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},  // normal symbol
			2: {0, 0, 50, 500, 5000},     // normal symbol
			9: {0, 0, 20, 200, 2000},     // wild symbol
		},
	}

	// Test case 1: Basic line win with no wilds
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	})
	assert.NoError(t, err)

	result1 := CalcLine2(scene1, pt, []int{1, 1, 1}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 })

	assert.NotNil(t, result1)
	assert.Equal(t, 1, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 100, result1.CashWin)
	assert.Equal(t, 0, result1.Wilds)

	// Test case 2: Line win with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{9, 1, 1},
		{9, 1, 1},
		{1, 1, 1},
	})
	assert.NoError(t, err)

	result2 := CalcLine2(scene2, pt, []int{0, 0, 0}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 })

	assert.NotNil(t, result2)
	assert.Equal(t, 1, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)
	assert.Equal(t, 100, result2.CashWin)
	assert.Equal(t, 2, result2.Wilds)

	// Test case 3: Line win with multipliers
	scene3, err := NewGameSceneWithArr2([][]int{
		{1, 1, 1},
		{1, 1, 1},
		{1, 1, 1},
	})
	assert.NoError(t, err)

	result3 := CalcLine2(scene3, pt, []int{1, 1, 1}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 2 }) // 2x multiplier for each position

	assert.NotNil(t, result3)
	assert.Equal(t, 1, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums)
	assert.Equal(t, 800, result3.CashWin) // 100 * 2 * 2 * 2 = 800
	assert.Equal(t, 0, result3.Wilds)
}

func Test_CalcLineRL2(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},  // normal symbol
			2: {0, 0, 50, 500, 5000},     // normal symbol
			9: {0, 0, 20, 200, 2000},     // wild symbol
		},
	}

	// Test case 1: Basic right-to-left win
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 1, 2},
		{1, 1, 2},
		{1, 1, 2},
	})
	assert.NoError(t, err)

	result1 := CalcLineRL2(scene1, pt, []int{2, 2, 2}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 })

	assert.NotNil(t, result1)
	assert.Equal(t, 2, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 50, result1.CashWin)
	assert.Equal(t, 0, result1.Wilds)

	// Test case 2: Right-to-left win with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{2, 1, 1},
		{9, 1, 1},  // Only one wild in the middle
		{2, 1, 1},
	})
	assert.NoError(t, err)

	result2 := CalcLineRL2(scene2, pt, []int{0, 0, 0}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 })

	assert.NotNil(t, result2)
	assert.Equal(t, 2, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)
	assert.Equal(t, 50, result2.CashWin)
	assert.Equal(t, 1, result2.Wilds)  // Only one wild in the winning pattern

	// Additional test case for wild counting
	scene2b, err := NewGameSceneWithArr2([][]int{
		{2, 1, 1},
		{9, 1, 1},  // Only one wild in the middle
		{2, 1, 1},
	})
	assert.NoError(t, err)

	result2b := CalcLineRL2(scene2b, pt, []int{0, 0, 0}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 })

	assert.NotNil(t, result2b)
	assert.Equal(t, 2, result2b.Symbol)
	assert.Equal(t, 3, result2b.SymbolNums)
	assert.Equal(t, 50, result2b.CashWin)
	assert.Equal(t, 1, result2b.Wilds)  // Only one wild in the winning pattern

	// Test case 3: Right-to-left win with multipliers
	scene3, err := NewGameSceneWithArr2([][]int{
		{2, 1, 2},
		{2, 2, 2},
		{1, 1, 2},
	})
	assert.NoError(t, err)

	result3 := CalcLineRL2(scene3, pt, []int{2, 2, 2}, 1,
		func(cursymbol int) bool { return cursymbol >= 1 && cursymbol <= 9 },
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 2 }) // 2x multiplier for each position

	assert.NotNil(t, result3)
	assert.Equal(t, 2, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums)
	assert.Equal(t, 400, result3.CashWin) // 50 * 2 * 2 * 2 = 400
	assert.Equal(t, 0, result3.Wilds)
}

func Test_CountSymbolOnLine(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},  // normal symbol
			2: {0, 0, 50, 500, 5000},     // normal symbol
			9: {0, 0, 20, 200, 2000},     // wild symbol
		},
	}

	// Test case 1: Basic symbol count
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 1, 1},
		{1, 2, 2},
		{1, 2, 2},
	})
	assert.NoError(t, err)

	result1 := CountSymbolOnLine(scene1, pt, []int{0, 0, 0}, 1, 1,
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 },
		func(src int, target int) int { return src * target })

	assert.NotNil(t, result1)
	assert.Equal(t, 1, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 100, result1.CashWin)
	assert.Equal(t, 0, result1.Wilds)

	// Test case 2: Symbol count with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 9, 1},
		{1, 9, 2},
		{1, 2, 2},
	})
	assert.NoError(t, err)

	result2 := CountSymbolOnLine(scene2, pt, []int{0, 1, 0}, 1, 1,
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 1 },
		func(src int, target int) int { return src * target })

	assert.NotNil(t, result2)
	assert.Equal(t, 1, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)
	assert.Equal(t, 100, result2.CashWin)
	assert.Equal(t, 1, result2.Wilds)  // Only counting one wild in the winning pattern

	// Test case 3: Symbol count with multipliers
	scene3, err := NewGameSceneWithArr2([][]int{
		{1, 1, 1},
		{2, 1, 2},
		{2, 2, 1},
	})
	assert.NoError(t, err)

	result3 := CountSymbolOnLine(scene3, pt, []int{0, 1, 2}, 1, 1,
		func(cursymbol int) bool { return cursymbol == 9 },
		func(cursymbol int, startsymbol int) bool { return cursymbol == startsymbol || cursymbol == 9 },
		func(cursymbol int) int { return cursymbol },
		func(x, y int) int { return 2 }, // 2x multiplier for each position
		func(src int, target int) int { return src * target })

	assert.NotNil(t, result3)
	assert.Equal(t, 1, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums)
	assert.Equal(t, 800, result3.CashWin) // 100 * 2 * 2 * 2 = 800
	assert.Equal(t, 0, result3.Wilds)
}
