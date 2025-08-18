package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_CalcLine2_Jules - Test for CalcLine2 with more cases
func Test_CalcLine2_Jules(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 10, 100, 1000, 10000}, // Symbol 'A', 2 symbols pay 10
			2: {0, 0, 50, 500, 5000},     // Symbol 'B'
			9: {0, 0, 200, 2000, 20000},  // Wild, pays more than A
		},
	}

	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }

	// Test case 1: No win
	scene1, err := NewGameSceneWithArr2([][]int{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 0},
	})
	assert.NoError(t, err)
	result1 := CalcLine2(scene1, pt, []int{0, 1, 2}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result1, "Test_CalcLine2_Jules: No win expected")

	// Test case 2: Invalid start symbol
	scene2, err := NewGameSceneWithArr2([][]int{
		{0}, {1}, {1},
	})
	assert.NoError(t, err)
	result2 := CalcLine2(scene2, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result2, "Test_CalcLine2_Jules: Invalid start symbol should result in no win")

	// Test case 3: All wilds line
	scene3, err := NewGameSceneWithArr2([][]int{
		{9}, {9}, {9},
	})
	assert.NoError(t, err)
	result3 := CalcLine2(scene3, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result3, "Test_CalcLine2_Jules: All wilds should be a win")
	assert.Equal(t, 9, result3.Symbol, "Test_CalcLine2_Jules: Symbol should be wild")
	assert.Equal(t, 3, result3.SymbolNums, "Test_CalcLine2_Jules: Should be 3 symbols")
	assert.Equal(t, 200, result3.CoinWin, "Test_CalcLine2_Jules: incorrect coin win for all wilds")

	// Test case 4: Wilds at the start, followed by a regular symbol
	scene4, err := NewGameSceneWithArr2([][]int{
		{9}, {9}, {1},
	})
	assert.NoError(t, err)
	result4 := CalcLine2(scene4, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result4, "Test_CalcLine2_Jules: Wilds followed by symbol should win")
	assert.Equal(t, 1, result4.Symbol, "Test_CalcLine2_Jules: Symbol should be 1")
	assert.Equal(t, 3, result4.SymbolNums, "Test_CalcLine2_Jules: Symbol count should be 3")
	assert.Equal(t, 100, result4.CoinWin, "Test_CalcLine2_Jules: Coin win for 3x symbol 1")

	// Test case 5: Wild payout is better
	scene5, err := NewGameSceneWithArr2([][]int{
		{9}, {9}, {2},
	})
	assert.NoError(t, err)
	pt.MapPay[9] = []int{0, 20, 200, 2000, 20000} // make wild pay for 2
	pt.MapPay[2][2] = 10                          // 3 of symbol 2 pays 10
	result5 := CalcLine2(scene5, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result5, "Test_CalcLine2_Jules: Wilds followed by symbol should win")
	assert.Equal(t, 9, result5.Symbol, "Test_CalcLine2_Jules: Symbol should be 9 as it pays more")
	assert.Equal(t, 2, result5.SymbolNums, "Test_CalcLine2_Jules: Symbol count for wild should be 2")
	assert.Equal(t, 20, result5.CoinWin, "Test_CalcLine2_Jules: Coin win for 2x wild")
	pt.MapPay[2][2] = 50 // reset paytable
	pt.MapPay[9] = []int{0, 0, 200, 2000, 20000}

	// Test case 6: Broken line
	scene6, err := NewGameSceneWithArr2([][]int{
		{1}, {1}, {2}, {1}, {1},
	})
	assert.NoError(t, err)
	result6 := CalcLine2(scene6, pt, []int{0, 0, 0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result6, "Test_CalcLine2_Jules: Broken line should have a win of 2 symbols")
	assert.Equal(t, 1, result6.Symbol, "Test_CalcLine2_Jules: Symbol should be 1")
	assert.Equal(t, 2, result6.SymbolNums, "Test_CalcLine2_Jules: Symbol count should be 2")
	assert.Equal(t, 10, result6.CoinWin, "Test_CalcLine2_Jules: Payout for 2 symbols is 10")
}

// Test_CalcLineRL2_Jules - Test for CalcLineRL2 with more cases
func Test_CalcLineRL2_Jules(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000}, // Symbol 'A'
			2: {0, 0, 50, 500, 5000},    // Symbol 'B'
			9: {0, 0, 200, 2000, 20000}, // Wild
		},
	}

	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }

	// Test case 1: Basic RL win
	scene1, err := NewGameSceneWithArr2([][]int{
		{1}, {1}, {1},
	})
	assert.NoError(t, err)
	result1 := CalcLineRL2(scene1, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result1)
	assert.Equal(t, 1, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 100, result1.CoinWin)

	// Test case 2: No win
	scene2, err := NewGameSceneWithArr2([][]int{
		{1}, {2}, {1},
	})
	assert.NoError(t, err)
	result2 := CalcLineRL2(scene2, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result2, "Test_CalcLineRL2_Jules: No win expected")

	// Test case 3: All wilds
	scene3, err := NewGameSceneWithArr2([][]int{
		{9}, {9}, {9},
	})
	assert.NoError(t, err)
	result3 := CalcLineRL2(scene3, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result3)
	assert.Equal(t, 9, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums)
	assert.Equal(t, 200, result3.CoinWin)
}

// Test_CountSymbolOnLine_Jules - Test for CountSymbolOnLine with more cases
func Test_CountSymbolOnLine_Jules(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},
			9: {0, 0, 200, 2000, 20000},
		},
	}
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }
	calcMulti := func(a, b int) int { return a * b }

	// Test case 1: Symbol not present
	scene1, err := NewGameSceneWithArr2([][]int{{2}, {3}, {4}})
	assert.NoError(t, err)
	result1 := CountSymbolOnLine(scene1, pt, []int{0, 0, 0}, 1, 1, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result1)

	// Test case 2: Symbol present but not at start
	scene2, err := NewGameSceneWithArr2([][]int{{2}, {1}, {1}})
	assert.NoError(t, err)
	result2 := CountSymbolOnLine(scene2, pt, []int{0, 0, 0}, 1, 1, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result2, "Symbol must be at the start of the line")

	// Test case 3: Starts with wild, counts as symbol
	scene3, err := NewGameSceneWithArr2([][]int{{9}, {1}, {1}})
	assert.NoError(t, err)
	result3 := CountSymbolOnLine(scene3, pt, []int{0, 0, 0}, 1, 1, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result3)
	assert.Equal(t, 1, result3.Symbol)
	assert.Equal(t, 3, result3.SymbolNums)
	assert.Equal(t, 100, result3.CoinWin)
}

// Test_CalcLine3_Jules - Test for CalcLine3 (untested function)
func Test_CalcLine3_Jules(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},
			9: {0, 0, 200, 2000, 20000},
		},
	}
	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }

	// Test with multiplicative multiplier
	getMulti := func(x, y int) int { return 2 }
	calcMulti := func(a, b int) int { return a * b }

	scene, err := NewGameSceneWithArr2([][]int{{1}, {1}, {1}})
	assert.NoError(t, err)
	result := CalcLine3(scene, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Symbol)
	assert.Equal(t, 3, result.SymbolNums)
	assert.Equal(t, 800, result.CashWin)

	// Test with additive multiplier
	getMulti = func(x, y int) int { return 2 }
	calcMulti = func(a, b int) int { return a + b }
	result2 := CalcLine3(scene, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result2)
	assert.Equal(t, 1, result2.Symbol)
	assert.Equal(t, 3, result2.SymbolNums)
	assert.Equal(t, 6, result2.OtherMul)
	assert.Equal(t, 600, result2.CashWin)
}

// Test_CalcLineRL3_Jules - Test for CalcLineRL3 (untested function)
func Test_CalcLineRL3_Jules(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000},
			9: {0, 0, 200, 2000, 20000},
		},
	}
	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 2 }
	calcMulti := func(a, b int) int { return a * b }

	scene, err := NewGameSceneWithArr2([][]int{{9}, {9}, {2}})
	assert.NoError(t, err)
	// RL win: starts with 2, then sees 9, 9. 3 symbols.
	// Pay for three 2s is 50.
	// Multiplier is 2*2*2=8.
	// Win is 50 * 8 = 400
	pt.MapPay[2] = []int{0, 0, 50, 500, 5000}
	result := CalcLineRL3(scene, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result)
	assert.Equal(t, 2, result.Symbol)
	assert.Equal(t, 3, result.SymbolNums)
	assert.Equal(t, 400, result.CashWin)
}

// Test_CalcLineRL2_Jules_More - Add more tests for CalcLineRL2
func Test_CalcLineRL2_Jules_More(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 100, 1000, 10000}, // Symbol 'A'
			2: {0, 10, 50, 500, 5000},    // Symbol 'B'
			9: {0, 20, 200, 2000, 20000}, // Wild
		},
	}

	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }

	// Wilds at the start (right), regular symbol pays more
	scene1, err := NewGameSceneWithArr2([][]int{
		{2}, {9}, {9},
	})
	assert.NoError(t, err)
	result1 := CalcLineRL2(scene1, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result1)
	assert.Equal(t, 2, result1.Symbol)
	assert.Equal(t, 3, result1.SymbolNums)
	assert.Equal(t, 50, result1.CoinWin)

	// Wilds at the start (right), wild pays more
	scene2, err := NewGameSceneWithArr2([][]int{
		{1}, {9}, {9},
	})
	assert.NoError(t, err)
	// 3 of '1' would pay 100. 2 of '9' pays 20. So '1' should win.
	// Let's make wild pay more.
	pt.MapPay[1][2] = 10 // 3 of '1' pays 10
	result2 := CalcLineRL2(scene2, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result2)
	assert.Equal(t, 9, result2.Symbol)
	assert.Equal(t, 2, result2.SymbolNums)
	assert.Equal(t, 20, result2.CoinWin)
}

// Test_CountSymbolOnLine_Jules_More - Add more tests for CountSymbolOnLine
func Test_CountSymbolOnLine_Jules_More(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			2: {0, 10, 50, 500, 5000},    // Symbol 'B'
			9: {0, 20, 200, 2000, 20000}, // Wild
		},
	}
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }
	calcMulti := func(a, b int) int { return a * b }

	// Starts with wild, wild payout is better
	scene1, err := NewGameSceneWithArr2([][]int{{9}, {9}, {2}})
	assert.NoError(t, err)
	result1 := CountSymbolOnLine(scene1, pt, []int{0, 0, 0}, 1, 2, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result1)
	// wnums=2, s0=9, wmul = 20
	// nums=3, ws=2, mul = 50
	// B win is better.
	assert.Equal(t, 2, result1.Symbol)

	// Make wild win better
	pt.MapPay[2][2] = 15 // 3 of '2' pays 15, less than 2 wilds (20)
	result2 := CountSymbolOnLine(scene1, pt, []int{0, 0, 0}, 1, 2, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result2)
	assert.Equal(t, 9, result2.Symbol)
	assert.Equal(t, 20, result2.CoinWin)
}

// Test_CalcLine3_Jules_Full - full tests for CalcLine3
func Test_CalcLine3_Jules_Full(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 10, 100, 1000, 10000},
			2: {0, 0, 50, 500, 5000},
			9: {0, 20, 200, 2000, 20000},
		},
	}

	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }
	calcMulti := func(a, b int) int { return a * b }

	// No win
	scene1, _ := NewGameSceneWithArr2([][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 0}})
	result1 := CalcLine3(scene1, pt, []int{0, 1, 2}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result1)

	// Invalid start symbol
	scene2, _ := NewGameSceneWithArr2([][]int{{0}, {1}, {1}})
	result2 := CalcLine3(scene2, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result2)

	// All wilds
	scene3, _ := NewGameSceneWithArr2([][]int{{9}, {9}, {9}})
	result3 := CalcLine3(scene3, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result3)
	assert.Equal(t, 9, result3.Symbol)
	assert.Equal(t, 200, result3.CoinWin)

	// Wilds + regular, regular pays more
	scene4, _ := NewGameSceneWithArr2([][]int{{9}, {9}, {1}})
	pt.MapPay[1][2] = 100 // 3 of '1' pays 100
	pt.MapPay[9][1] = 20  // 2 of '9' pays 20
	result4 := CalcLine3(scene4, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result4)
	assert.Equal(t, 1, result4.Symbol)
	assert.Equal(t, 100, result4.CoinWin)

	// Wilds + regular, wild pays more
	pt.MapPay[1][2] = 10 // 3 of '1' pays 10
	result5 := CalcLine3(scene4, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result5)
	assert.Equal(t, 9, result5.Symbol)
	assert.Equal(t, 20, result5.CoinWin)

	// Broken line
	scene6, _ := NewGameSceneWithArr2([][]int{{1}, {1}, {2}, {1}, {1}})
	result6 := CalcLine3(scene6, pt, []int{0, 0, 0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result6)
	assert.Equal(t, 1, result6.Symbol)
	assert.Equal(t, 2, result6.SymbolNums)
	assert.Equal(t, 10, result6.CoinWin)
}

// Test_CalcLineRL3_Jules_Full - full tests for CalcLineRL3
func Test_CalcLineRL3_Jules_Full(t *testing.T) {
	pt := &PayTables{
		MapPay: map[int][]int{
			1: {0, 10, 100, 1000, 10000},
			2: {0, 0, 50, 500, 5000},
			9: {0, 20, 200, 2000, 20000},
		},
	}

	isValidSymbol := func(s int) bool { return s >= 1 && s <= 9 }
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }
	calcMulti := func(a, b int) int { return a * b }

	// No win
	scene1, _ := NewGameSceneWithArr2([][]int{{1}, {2}, {1}})
	result1 := CalcLineRL3(scene1, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result1)

	// Invalid start symbol
	scene2, _ := NewGameSceneWithArr2([][]int{{1}, {1}, {0}})
	result2 := CalcLineRL3(scene2, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result2)

	// All wilds
	scene3, _ := NewGameSceneWithArr2([][]int{{9}, {9}, {9}})
	result3 := CalcLineRL3(scene3, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result3)
	assert.Equal(t, 9, result3.Symbol)
	assert.Equal(t, 200, result3.CoinWin)

	// Wilds + regular, wild pays more
	scene4, _ := NewGameSceneWithArr2([][]int{{1}, {9}, {9}})
	pt.MapPay[1][2] = 10 // 3 of '1' pays 10
	pt.MapPay[9][1] = 20  // 2 of '9' pays 20
	result4 := CalcLineRL3(scene4, pt, []int{0, 0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.NotNil(t, result4)
	assert.Equal(t, 9, result4.Symbol)
	assert.Equal(t, 20, result4.CoinWin)
}

// Test_Algorithm2_Coverage - for improving coverage
func Test_Algorithm2_Coverage(t *testing.T) {
	// Gap 2 & 3: All wilds no pay, and both wild/regular no pay
	pt1 := &PayTables{
		MapPay: map[int][]int{
			1: {0, 0, 10, 100, 1000}, // 2 of 1 pays 0
			9: {0, 0, 20, 200, 2000}, // 2 of 9 pays 0
		},
	}
	isValidSymbol := func(s int) bool { return s >= 0 && s <= 9 } // 0 is valid now
	isWild := func(s int) bool { return s == 9 }
	isSameSymbol := func(s1, s2 int) bool { return s1 == s2 || s1 == 9 }
	getSymbol := func(s int) int { return s }
	getMulti := func(x, y int) int { return 1 }
	calcMulti := func(a, b int) int { return a * b }

	// Gap 2: All wilds, no payout
	scene1, _ := NewGameSceneWithArr2([][]int{{9}, {9}})
	result1 := CalcLine2(scene1, pt1, []int{0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result1, "Test_Algorithm2_Coverage: All wilds with no payout should be nil")

	// Gap 3: Both wild and regular lines have no pay
	scene2, _ := NewGameSceneWithArr2([][]int{{9}, {1}})
	result2 := CalcLine2(scene2, pt1, []int{0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result2, "Test_Algorithm2_Coverage: Both wild and regular no pay should be nil")

	// Gap 1: Broken by invalid symbol
	pt2 := &PayTables{MapPay: map[int][]int{
		1: {0, 10, 100, 1000, 10000},
		9: {0, 0, 0, 0, 0},
	}}
	isValidSymbol2 := func(s int) bool { return s == 1 || s == 9 } // 0 is invalid
	scene3, _ := NewGameSceneWithArr2([][]int{{9}, {1}, {0}, {1}})
	result3 := CalcLine2(scene3, pt2, []int{0, 0, 0, 0}, 1, isValidSymbol2, isWild, isSameSymbol, getSymbol, getMulti)
	assert.NotNil(t, result3)
	assert.Equal(t, 2, result3.SymbolNums, "Test_Algorithm2_Coverage: Should break at invalid symbol")

	// Also test the RL versions with one of the gap cases
	result4 := CalcLineRL2(scene1, pt1, []int{0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti)
	assert.Nil(t, result4, "Test_Algorithm2_Coverage: RL All wilds with no payout should be nil")

	result5 := CalcLine3(scene1, pt1, []int{0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result5, "Test_Algorithm2_Coverage: All wilds with no payout should be nil (CalcLine3)")

	result6 := CalcLineRL3(scene1, pt1, []int{0, 0}, 1, isValidSymbol, isWild, isSameSymbol, getSymbol, getMulti, calcMulti)
	assert.Nil(t, result6, "Test_Algorithm2_Coverage: RL All wilds with no payout should be nil (CalcLineRL3)")
}
