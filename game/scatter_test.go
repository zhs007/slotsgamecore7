package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcScatterComprehensive(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test 1: Basic scatter calculation
	scene1 := &GameScene{
		Arr: [][]int{
			{1, 11, 1},
			{11, 2, 9},
			{7, 11, 7},
		},
	}

	result1 := CalcScatter(scene1, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result1)
	assert.Equal(t, result1.Symbol, 11)
	assert.Equal(t, len(result1.Pos), 6)

	// Test 2: Maximum scatter symbols (more than reel count)
	scene2 := &GameScene{
		Arr: [][]int{
			{11, 11, 11},
			{11, 11, 11},
			{11, 11, 11},
		},
	}

	result2 := CalcScatter(scene2, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result2)
	assert.Equal(t, result2.Symbol, 11)
	assert.Equal(t, result2.SymbolNums, len(scene2.Arr)) // Should be capped at reel count

	// Test 3: No scatter symbols
	scene3 := &GameScene{
		Arr: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	result3 := CalcScatter(scene3, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Nil(t, result3)

	// Test 8: Different bet and coin values
	scene8 := &GameScene{
		Arr: [][]int{
			{11, 2, 11},
			{4, 11, 6},
			{11, 8, 11},
		},
	}

	result8 := CalcScatter(scene8, pt, 11, 5, 20, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result8)
	assert.Equal(t, result8.Symbol, 11)
	assert.True(t, result8.CashWin > 0)
	assert.True(t, result8.CoinWin > 0)

	// Test 9: Irregular grid
	scene9 := &GameScene{
		Arr: [][]int{
			{11, 2},
			{4, 11, 6, 11},
			{11},
		},
	}

	result9 := CalcScatter(scene9, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result9)
	assert.Equal(t, result9.Symbol, 11)

	t.Logf("Test_CalcScatterComprehensive OK")
}

func Test_CalcScatter2Comprehensive(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test 1: Basic scatter2 calculation
	scene1 := &GameScene{
		Arr: [][]int{
			{1, 11, 1},
			{11, 2, 9},
			{7, 11, 7},
		},
	}

	result1 := CalcScatter2(scene1, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result1)
	assert.Equal(t, result1.Symbol, 11)

	// Test 2: Maximum scatter symbols (more than paytable length)
	scene2 := &GameScene{
		Arr: [][]int{
			{11, 11, 11},
			{11, 11, 11},
			{11, 11, 11},
		},
	}

	result2 := CalcScatter2(scene2, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.NotNil(t, result2)
	assert.Equal(t, result2.Symbol, 11)
	assert.Equal(t, result2.SymbolNums, len(pt.MapPay[11])) // Should be capped at paytable length

	// Test 3: No scatter symbols
	scene3 := &GameScene{
		Arr: [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		},
	}

	result3 := CalcScatter2(scene3, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Nil(t, result3)

	t.Logf("Test_CalcScatter2Comprehensive OK")
}

func Test_CalcScatter3Comprehensive(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test 1: Basic scatter3 calculation with isOnlyOneOnReel=true
	scene1 := &GameScene{
		Arr: [][]int{
			{1, 11, 1},
			{11, 2, 9},
			{7, 11, 7},
		},
	}

	result1 := CalcScatter3(scene1, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	}, true)

	assert.NotNil(t, result1)
	assert.Equal(t, result1.Symbol, 11)

	// Test 2: Multiple scatters on same reel with isOnlyOneOnReel=true
	scene2 := &GameScene{
		Arr: [][]int{
			{11, 11, 11},
			{11, 2, 11},
			{11, 11, 11},
		},
	}

	result2 := CalcScatter3(scene2, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	}, true)

	assert.NotNil(t, result2)
	assert.Equal(t, result2.Symbol, 11)
	assert.Equal(t, result2.SymbolNums, 3) // Should only count one per reel

	// Test 3: Same test with isOnlyOneOnReel=false
	result3 := CalcScatter3(scene2, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	}, false)

	assert.NotNil(t, result3)
	assert.Equal(t, result3.Symbol, 11)
	assert.True(t, result3.SymbolNums > 3) // Should count all scatters

	t.Logf("Test_CalcScatter3Comprehensive OK")
}

func Test_CalcScatter5Comprehensive(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// Test 1: Basic scatter5 calculation
	scene1 := &GameScene{
		Arr: [][]int{
			{1, 11, 1},
			{11, 2, 9},
			{7, 11, 7},
		},
	}

	result1 := CalcScatter5(scene1, pt, 11, 2, func(s int, cs int) bool {
		return cs == s
	}, true, 3, false)

	assert.NotNil(t, result1)
	assert.Equal(t, result1.Symbol, 11)

	// Test 2: Test with reversal height
	scene2 := &GameScene{
		Arr: [][]int{
			{11, 11, 11},
			{11, 2, 11},
			{11, 11, 11},
		},
	}

	result2 := CalcScatter5(scene2, pt, 11, 2, func(s int, cs int) bool {
		return cs == s
	}, true, 3, true)

	assert.NotNil(t, result2)
	assert.Equal(t, result2.Symbol, 11)

	// Test 3: Test with different height
	scene3 := &GameScene{
		Arr: [][]int{
			{11, 11, 11, 11},
			{11, 2, 11, 11},
			{11, 11, 11, 11},
		},
	}

	result3 := CalcScatter5(scene3, pt, 11, 2, func(s int, cs int) bool {
		return cs == s
	}, true, 4, false)

	assert.NotNil(t, result3)
	assert.Equal(t, result3.Symbol, 11)

	t.Logf("Test_CalcScatter5Comprehensive OK")
}
