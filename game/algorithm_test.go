package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcScatter(t *testing.T) {
	// Load test paytable
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test case 1: Basic scatter calculation with 3 scatters
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 9, 2},
		{3, 9, 4},
		{5, 9, 6},
	})
	assert.NoError(t, err)

	result1 := CalcScatter(scene1, pt, 9, 1, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.NotNil(t, result1)
	assert.Equal(t, 9, result1.Symbol)
	assert.Equal(t, 1, int(result1.Type))  // Added explicit conversion to int
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 6, len(result1.Pos))

	// Test case 2: Less than 3 scatter symbols (should return nil)
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 9, 3},
		{4, 9, 6},
		{7, 8, 1},
	})
	assert.NoError(t, err)

	result2 := CalcScatter(scene2, pt, 9, 1, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.Nil(t, result2)

	// Test case 3: More than 3 scatters
	scene3, err := NewGameSceneWithArr2([][]int{
		{9, 9, 9},
		{9, 9, 9},
		{9, 9, 9},
	})
	assert.NoError(t, err)

	result3 := CalcScatter(scene3, pt, 9, 2, 3, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.NotNil(t, result3)
	assert.Equal(t, 3, result3.SymbolNums) // Should be capped at number of reels
	assert.True(t, result3.CoinWin > 0)
	assert.True(t, result3.CashWin > 0)

	// Test case 4: Complex scatter pattern with exactly 3 scatters
	scene4 := &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 11, 9},
			{7, 1, 7},
			{6, 11, 11},
			{1, 9, 0},
		},
	}

	result4 := CalcScatter(scene4, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Equal(t, result4.Symbol, 11)
	assert.Equal(t, result4.Mul, 5)
	assert.Equal(t, result4.CoinWin, 50)
	assert.Equal(t, result4.CashWin, 100)
	assert.Equal(t, len(result4.Pos), 6)

	// Test case 5: Complex scatter pattern with 5 scatters
	scene5 := &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 11, 9},
			{11, 1, 7},
			{6, 11, 11},
			{1, 11, 0},
		},
	}

	result5 := CalcScatter(scene5, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Equal(t, result5.Symbol, 11)
	assert.Equal(t, result5.Mul, 100)
	assert.Equal(t, result5.CoinWin, 1000)
	assert.Equal(t, result5.CashWin, 2000)
	assert.Equal(t, len(result5.Pos), 10)

	t.Logf("Test_CalcScatter OK")
}

func Test_CalcLine(t *testing.T) {
	// Load test paytable
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test case 1: Basic line win with no wilds
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 6, 7},
		{8, 9, 1},
		{2, 3, 4},
	})
	assert.NoError(t, err)

	line := []int{0, 0, 0, 1, 2} // Line pattern
	result1 := CalcLine(scene1, pt, line, 1,
		func(cursymbol int) bool { return cursymbol >= 0 },    // isValidSymbol
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start },     // isSameSymbol
		func(cursymbol int) int { return cursymbol },          // getSymbol
	)

	assert.NotNil(t, result1)
	assert.Equal(t, 1, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)

	// Test case 2: Line win with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{0, 4, 5},
		{1, 6, 7},
		{0, 9, 1},
		{2, 3, 4},
	})
	assert.NoError(t, err)

	result2 := CalcLine(scene2, pt, line, 1,
		func(cursymbol int) bool { return cursymbol >= 0 },    // isValidSymbol
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start || cur == 0 }, // isSameSymbol with wild
		func(cursymbol int) int { return cursymbol },          // getSymbol
	)

	assert.NotNil(t, result2)
	assert.Equal(t, 1, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)

	// Test case 3: No line win
	scene3, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{2, 4, 5},
		{3, 6, 7},
		{4, 9, 1},
		{5, 3, 4},
	})
	assert.NoError(t, err)

	result3 := CalcLine(scene3, pt, line, 1,
		func(cursymbol int) bool { return cursymbol >= 0 },    // isValidSymbol
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start },     // isSameSymbol
		func(cursymbol int) int { return cursymbol },          // getSymbol
	)

	assert.Nil(t, result3)

	t.Logf("Test_CalcLine OK")
}

func Test_CalcLineRL(t *testing.T) {
	// Load test paytable
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test case 1: Basic right to left win
	scene1, err := NewGameSceneWithArr2([][]int{
		{2, 2, 3},
		{1, 4, 5},
		{1, 6, 7},
		{1, 9, 1},
		{1, 3, 4},
	})
	assert.NoError(t, err)

	line := []int{0, 0, 0, 1, 2} // Line pattern
	result1 := CalcLineRL(scene1, pt, line, 1,
		func(cursymbol int) bool { return cursymbol >= 0 },    // isValidSymbol
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start },     // isSameSymbol
		func(cursymbol int) int { return cursymbol },          // getSymbol
	)

	assert.NotNil(t, result1)
	assert.Equal(t, 1, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)

	// Test case 2: Right to left with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{2, 2, 3},
		{0, 4, 5},
		{1, 6, 7},
		{1, 9, 1},
		{1, 3, 4},
	})
	assert.NoError(t, err)

	result2 := CalcLineRL(scene2, pt, line, 1,
		func(cursymbol int) bool { return cursymbol >= 0 },    // isValidSymbol
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start || cur == 0 }, // isSameSymbol with wild
		func(cursymbol int) int { return cursymbol },          // getSymbol
	)

	assert.NotNil(t, result2)
	assert.Equal(t, 1, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)

	t.Logf("Test_CalcLineRL OK")
}

func Test_CalcFullLine(t *testing.T) {
	// Load test paytable
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test case 1: Basic full line calculation
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 6, 7},
		{1, 9, 1},
		{1, 3, 4},
	})
	assert.NoError(t, err)

	results1 := CalcFullLine(scene1, pt, 1,
		func(cursymbol int, scene *GameScene, x, y int) bool { return cursymbol >= 0 }, // isValidSymbolEx
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start },     // isSameSymbol
	)

	assert.NotNil(t, results1)
	assert.Greater(t, len(results1), 0)
	assert.Equal(t, 1, results1[0].Symbol)
	assert.Equal(t, 5, results1[0].SymbolNums)

	// Test case 2: Full line with wilds
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{0, 4, 5},
		{1, 6, 7},
		{0, 9, 1},
		{1, 3, 4},
	})
	assert.NoError(t, err)

	results2 := CalcFullLine(scene2, pt, 1,
		func(cursymbol int, scene *GameScene, x, y int) bool { return cursymbol >= 0 }, // isValidSymbolEx
		func(cursymbol int) bool { return cursymbol == 0 },    // isWild
		func(cur, start int) bool { return cur == start || cur == 0 }, // isSameSymbol with wild
	)

	assert.NotNil(t, results2)
	assert.Greater(t, len(results2), 0)

	t.Logf("Test_CalcFullLine OK")
}

func Test_CountScatterInArea(t *testing.T) {
	// Test case 1: Basic area scatter count
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 9, 2},
		{9, 9, 4},
		{5, 9, 6},
	})
	assert.NoError(t, err)

	result1 := CountScatterInArea(scene1, 9, 3,
		func(x, y int) bool { return x < 2 && y < 2 }, // 2x2 area in top-left
		func(scatter, cursymbol int) bool { return cursymbol == scatter },
	)

	assert.NotNil(t, result1)
	assert.Equal(t, 9, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)

	// Test case 2: No scatters in area
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 9},
		{3, 4, 9},
		{5, 6, 9},
	})
	assert.NoError(t, err)

	result2 := CountScatterInArea(scene2, 9, 2,
		func(x, y int) bool { return x < 2 && y < 2 }, // 2x2 area in top-left
		func(scatter, cursymbol int) bool { return cursymbol == scatter },
	)

	assert.Nil(t, result2)

	t.Logf("Test_CountScatterInArea OK")
}

func Test_CalcReelScatterEx(t *testing.T) {
	// Test case 1: Basic reel scatter with exactly required number
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 9, 2},
		{3, 9, 4},
		{5, 9, 6},
	})
	assert.NoError(t, err)

	result1 := CalcReelScatterEx(scene1, 9, 3, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.NotNil(t, result1)
	assert.Equal(t, 9, result1.Symbol)
	assert.Equal(t, int(RTScatterEx), int(result1.Type))  // Added explicit conversion to int
	assert.Equal(t, 3, result1.SymbolNums)

	// Test case 2: More than required scatters
	scene2, err := NewGameSceneWithArr2([][]int{
		{9, 9, 9},
		{9, 9, 9},
		{9, 9, 9},
	})
	assert.NoError(t, err)

	result2 := CalcReelScatterEx(scene2, 9, 2, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.NotNil(t, result2)
	assert.Equal(t, 3, result2.SymbolNums)

	// Test case 3: Less than required scatters
	scene3, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{9, 4, 5},
		{6, 7, 8},
	})
	assert.NoError(t, err)

	result3 := CalcReelScatterEx(scene3, 9, 2, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	})

	assert.Nil(t, result3)

	t.Logf("Test_CalcReelScatterEx OK")
}

func Test_CalcScatter5(t *testing.T) {
	// Load test paytable
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test case 1: Basic scatter with height limit from top
	scene1, err := NewGameSceneWithArr2([][]int{
		{9, 1, 2},
		{9, 4, 5},
		{9, 7, 8},
	})
	assert.NoError(t, err)

	result1 := CalcScatter5(scene1, pt, 9, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	}, true, 2, false)

	assert.NotNil(t, result1)
	assert.Equal(t, 9, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)

	// Test case 2: Basic scatter with height limit from bottom
	scene2, err := NewGameSceneWithArr2([][]int{
		{1, 2, 9},
		{4, 5, 9},
		{7, 8, 9},
	})
	assert.NoError(t, err)

	result2 := CalcScatter5(scene2, pt, 9, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	}, true, 2, true)

	assert.NotNil(t, result2)
	assert.Equal(t, 9, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)

	// Test case 3: Multiple scatters per reel with height limit
	scene3, err := NewGameSceneWithArr2([][]int{
		{9, 9, 9},
		{9, 9, 9},
		{9, 9, 9},
	})
	assert.NoError(t, err)

	result3 := CalcScatter5(scene3, pt, 9, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	}, true, 2, false)

	assert.NotNil(t, result3)
	assert.Equal(t, 9, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums) // Should be capped at number of reels

	// Test case 4: Invalid height parameter
	scene4, err := NewGameSceneWithArr2([][]int{
		{9, 1, 2},
		{9, 4, 5},
		{9, 7, 8},
	})
	assert.NoError(t, err)

	result4 := CalcScatter5(scene4, pt, 9, 1, func(scatter, cursymbol int) bool {
		return cursymbol == scatter
	}, true, 0, false)

	assert.NotNil(t, result4)
	assert.Equal(t, 9, result4.Symbol)
	assert.Equal(t, 3, result4.SymbolNums)

	t.Logf("Test_CalcScatter5 OK")
}
