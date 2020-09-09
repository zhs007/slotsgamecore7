package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcScatter(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 11, 9},
			[]int{7, 1, 7},
			[]int{6, 11, 11},
			[]int{1, 9, 0},
		},
	}

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	if err != nil {
		t.Fatalf("Test_CalcScatter LoadPayTables5JSON error %v",
			err)
	}

	result := CalcScatter(scene, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Equal(t, result.Symbol, 11, "they should be equal")
	assert.Equal(t, result.Mul, 5, "they should be equal")
	assert.Equal(t, result.CoinWin, 50, "they should be equal")
	assert.Equal(t, result.CashWin, 100, "they should be equal")
	assert.Equal(t, len(result.Pos), 6, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 11, 9},
			[]int{11, 1, 7},
			[]int{6, 11, 11},
			[]int{1, 11, 0},
		},
	}

	result = CalcScatter(scene, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Equal(t, result.Symbol, 11, "they should be equal")
	assert.Equal(t, result.Mul, 100, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 2000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			[]int{11, 0, 11},
			[]int{9, 11, 9},
			[]int{11, 1, 7},
			[]int{6, 11, 11},
			[]int{1, 11, 0},
		},
	}

	result = CalcScatter(scene, pt, 11, 2, 10, func(s int, cs int) bool {
		return cs == s
	})

	assert.Equal(t, result.Symbol, 11, "they should be equal")
	assert.Equal(t, result.Mul, 100, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 2000, "they should be equal")
	assert.Equal(t, len(result.Pos), 14, "they should be equal")

	t.Logf("Test_CalcScatter OK")
}

func Test_CalcLine(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 1, 9},
			[]int{7, 1, 7},
			[]int{6, 1, 6},
			[]int{1, 9, 0},
		},
	}

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	if err != nil {
		t.Fatalf("Test_CalcLine LoadPayTables5JSON error %v",
			err)
	}

	ld, err := LoadLine5JSON("../unittestdata/linedata.json")
	if err != nil {
		t.Fatalf("Test_CalcLine LoadLine5JSON error %v",
			err)
	}

	// 0,1,1,1,9 => 1x4
	result := CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 200, "they should be equal")
	assert.Equal(t, result.CoinWin, 200, "they should be equal")
	assert.Equal(t, result.CashWin, 200, "they should be equal")
	assert.Equal(t, len(result.Pos), 8, "they should be equal")

	// 1,1,1,1,0 => 1x5
	result = CalcLine(scene, pt, ld.Lines[10], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			[]int{9, 0, 11},
			[]int{9, 0, 11},
			[]int{0, 0, 11},
			[]int{9, 0, 11},
			[]int{1, 1, 11},
		},
	}

	// 0,0,0,0,1 => 0x4 | 1x5 => 1x5
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	// 11,0,0,0,11 => nil
	result = CalcLine(scene, pt, ld.Lines[10], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Nil(t, result, "it should be nil")

	// 11,11,11,11,11 => nil
	result = CalcLine(scene, pt, ld.Lines[2], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Nil(t, result, "it should be nil")

	// 9,9,0,9,1 => 9x4
	result = CalcLine(scene, pt, ld.Lines[1], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 9, "they should be equal")
	assert.Equal(t, result.Mul, 15, "they should be equal")
	assert.Equal(t, result.CoinWin, 15, "they should be equal")
	assert.Equal(t, result.CashWin, 15, "they should be equal")
	assert.Equal(t, len(result.Pos), 8, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 1, 9},
			[]int{7, 1, 7},
			[]int{6, 0, 6},
			[]int{1, 1, 0},
		},
	}

	// 0,1,1,0,1 => 1x5
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 1, 9},
			[]int{7, 2, 7},
			[]int{6, 2, 6},
			[]int{1, 0, 0},
		},
	}

	// 0,1,2,2,0 => nil
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Nil(t, result, "it should be nil")

	scene = &GameScene{
		Arr: [][]int{
			[]int{1, 0, 1},
			[]int{9, 1, 9},
			[]int{7, 1, 7},
			[]int{6, 2, 6},
			[]int{1, 0, 0},
		},
	}

	// 0,1,1,2,0 => 1x3
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 50, "they should be equal")
	assert.Equal(t, result.CoinWin, 50, "they should be equal")
	assert.Equal(t, result.CashWin, 50, "they should be equal")
	assert.Equal(t, len(result.Pos), 6, "they should be equal")

	t.Logf("Test_CalcLine OK")
}

func Test_CalcLine2(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			[]int{8, 10, 1},
			[]int{11, 10, 7},
			[]int{0, 4, 6},
			[]int{7, 8, 0},
			[]int{1, 9, 5},
		},
	}

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	ld, err := LoadLine5JSON("../unittestdata/linedata.json")
	assert.NoError(t, err)

	// 0,1,1,1,9 => 1x4
	result := CalcLine(scene, pt, ld.Lines[15], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || cs == 0 || s == 0
		})

	assert.Equal(t, result.Symbol, 10, "they should be equal")
	assert.Equal(t, result.Mul, 5, "they should be equal")
	assert.Equal(t, result.CoinWin, 5, "they should be equal")
	assert.Equal(t, result.CashWin, 5, "they should be equal")
	assert.Equal(t, len(result.Pos), 6, "they should be equal")
}
