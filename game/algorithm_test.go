package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcScatter(t *testing.T) {
	scene := &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 11, 9},
			{7, 1, 7},
			{6, 11, 11},
			{1, 9, 0},
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
			{1, 0, 1},
			{9, 11, 9},
			{11, 1, 7},
			{6, 11, 11},
			{1, 11, 0},
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
			{11, 0, 11},
			{9, 11, 9},
			{11, 1, 7},
			{6, 11, 11},
			{1, 11, 0},
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
			{1, 0, 1},
			{9, 1, 9},
			{7, 1, 7},
			{6, 1, 6},
			{1, 9, 0},
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
			return cs == s || s == 0
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
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{9, 0, 11},
			{9, 0, 11},
			{0, 0, 11},
			{9, 0, 11},
			{1, 1, 11},
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
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{9, 0, 11},
			{9, 0, 11},
			{0, 0, 11},
			{9, 0, 11},
			{1, 2, 11},
		},
	}

	// 0,0,0,0,2 => 0x4 | 2x5 => 0x4
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 0, "they should be equal")
	assert.Equal(t, result.Mul, 500, "they should be equal")
	assert.Equal(t, result.CoinWin, 500, "they should be equal")
	assert.Equal(t, result.CashWin, 500, "they should be equal")
	assert.Equal(t, len(result.Pos), 8, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{9, 0, 11},
			{9, 0, 11},
			{0, 0, 11},
			{9, 0, 11},
			{1, 2, 11},
		},
	}

	// 0,0,0,0,3 => 0x4 | 3x5 => 0x4
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 0, "they should be equal")
	assert.Equal(t, result.Mul, 500, "they should be equal")
	assert.Equal(t, result.CoinWin, 500, "they should be equal")
	assert.Equal(t, result.CashWin, 500, "they should be equal")
	assert.Equal(t, len(result.Pos), 8, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{9, 0, 11},
			{9, 0, 11},
			{0, 0, 11},
			{9, 0, 11},
			{1, 0, 11},
		},
	}

	// 0,0,0,0,0 => 0x5
	result = CalcLine(scene, pt, ld.Lines[0], 1,
		func(cs int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 0, "they should be equal")
	assert.Equal(t, result.Mul, 2000, "they should be equal")
	assert.Equal(t, result.CoinWin, 2000, "they should be equal")
	assert.Equal(t, result.CashWin, 2000, "they should be equal")
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
			return cs == s || s == 0
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
			return cs == s || s == 0
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
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 9, "they should be equal")
	assert.Equal(t, result.Mul, 15, "they should be equal")
	assert.Equal(t, result.CoinWin, 15, "they should be equal")
	assert.Equal(t, result.CashWin, 15, "they should be equal")
	assert.Equal(t, len(result.Pos), 8, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 1, 9},
			{7, 1, 7},
			{6, 0, 6},
			{1, 1, 0},
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
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 1, "they should be equal")
	assert.Equal(t, result.Mul, 1000, "they should be equal")
	assert.Equal(t, result.CoinWin, 1000, "they should be equal")
	assert.Equal(t, result.CashWin, 1000, "they should be equal")
	assert.Equal(t, len(result.Pos), 10, "they should be equal")

	scene = &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 1, 9},
			{7, 2, 7},
			{6, 2, 6},
			{1, 0, 0},
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
			return cs == s || s == 0
		})

	assert.Nil(t, result, "it should be nil")

	scene = &GameScene{
		Arr: [][]int{
			{1, 0, 1},
			{9, 1, 9},
			{7, 1, 7},
			{6, 2, 6},
			{1, 0, 0},
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
			return cs == s || s == 0
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
			{8, 10, 1},
			{11, 10, 7},
			{0, 4, 6},
			{7, 8, 0},
			{1, 9, 5},
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
			return cs == s || s == 0
		})

	assert.Equal(t, result.Symbol, 10, "they should be equal")
	assert.Equal(t, result.Mul, 5, "they should be equal")
	assert.Equal(t, result.CoinWin, 5, "they should be equal")
	assert.Equal(t, result.CashWin, 5, "they should be equal")
	assert.Equal(t, len(result.Pos), 6, "they should be equal")

	t.Logf("Test_CalcLine2 OK")
}

func Test_CalcFullLine(t *testing.T) {
	scene, err := NewGameSceneWithArr2([][]int{
		{8, 10, 1},
		{11, 10, 7},
		{0, 4, 6},
		{7, 8, 0},
		{1, 9, 5},
	})
	assert.NoError(t, err)

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// 0,1,1,1,9 => 1x4
	results := CalcFullLine(scene, pt, 1,
		func(cs int, scene *GameScene, x, y int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, len(results), 1)
	assert.Equal(t, len(results[0].Pos), 8)
	assert.Equal(t, results[0].Symbol, 10)
	assert.Equal(t, results[0].Mul, 15)

	t.Logf("Test_CalcFullLine OK")
}

func Test_CalcFullLine2(t *testing.T) {
	scene, err := NewGameSceneWithArr2([][]int{
		{8, 10, 7},
		{11, 10, 7},
		{0, 4, 6},
		{7, 8, 0},
		{1, 9, 5},
	})
	assert.NoError(t, err)

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	results := CalcFullLine(scene, pt, 1,
		func(cs int, scene *GameScene, x, y int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, len(results), 3)

	assert.Equal(t, len(results[0].Pos), 8)
	assert.Equal(t, results[0].Symbol, 10)
	assert.Equal(t, results[0].Mul, 15)

	assert.Equal(t, len(results[1].Pos), 8)
	assert.Equal(t, results[1].Symbol, 7)

	assert.Equal(t, len(results[2].Pos), 8)
	assert.Equal(t, results[2].Symbol, 7)

	t.Logf("Test_CalcFullLine2 OK")
}

func Test_CalcFullLineEx(t *testing.T) {
	scene, err := NewGameSceneWithArr2([][]int{
		{8, 10, 7},
		{11, 10, 7},
		{0, 4, 6},
		{7, 8, 0},
		{1, 9, 5},
	})
	assert.NoError(t, err)

	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// 0,1,1,1,9 => 1x4
	results := CalcFullLineEx(scene, pt, 1,
		func(cs int, scene *GameScene, x, y int) bool {
			return cs != 11
		},
		func(cs int) bool {
			return cs == 0
		},
		func(s int, cs int) bool {
			return cs == s || s == 0
		})

	assert.Equal(t, len(results), 2)

	assert.Equal(t, len(results[0].Pos), 8)
	assert.Equal(t, results[0].Symbol, 10)
	assert.Equal(t, results[0].Mul, 15)

	assert.Equal(t, len(results[1].Pos), 10)
	assert.Equal(t, results[1].Symbol, 7)
	assert.Equal(t, results[1].Mul, 30)
	assert.Equal(t, results[1].CoinWin, 30*2)

	t.Logf("Test_CalcFullLineEx OK")
}
