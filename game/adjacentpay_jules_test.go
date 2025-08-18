package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalcAdjacentPay_Jules_WildWin(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// In paytables.json, 5 of symbol 8 pays 800, while 3 of wild (symbol 0) pays 50.
	// The code should choose the win with symbol 8.
	sceneWildWin, err := NewGameSceneWithArr2([][]int{
		{0, 0, 0, 8, 8},
	})
	assert.NoError(t, err)

	resultWildWin, err := CalcAdjacentPay(sceneWildWin, pt, 1, func(cursymbol int) bool {
		return cursymbol >= 0
	}, func(cursymbol int) bool {
		return cursymbol == 0
	}, func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol || cursymbol == 0
	}, func(cursymbol int) int {
		return cursymbol
	})
	assert.NoError(t, err)
	assert.NotNil(t, resultWildWin)
	assert.Equal(t, 1, len(resultWildWin))
	assert.Equal(t, 8, resultWildWin[0].Symbol) // win with symbol 8 is better

	t.Logf("Test_CalcAdjacentPay_Jules_WildWin OK")
}

func Test_calcAdjacentPayWithX_Jules_Wilds(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// A single wild does not pay, so this should return nil
	scene, err := NewGameSceneWithArr2([][]int{
		{0},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithX(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.Nil(t, result)
}

func Test_calcAdjacentPayWithY_WildTakeOver(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// 3 wilds pay more than 5 of symbol 1
	pt.MapPay[0][2] = 1000
	pt.MapPay[1][4] = 10

	scene, err := NewGameSceneWithArr2([][]int{
		{0, 0, 0, 1, 1},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithY(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Symbol)
	assert.Equal(t, 3, result.SymbolNums)
}

func Test_calcAdjacentPayWithY_Jules_Wilds(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// 5 wilds pay
	scene, err := NewGameSceneWithArr2([][]int{
		{0, 0, 0, 0, 0},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithY(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.NotNil(t, result)
	assert.Equal(t, 5, result.SymbolNums)
}

func Test_calcAdjacentPay_NoWin(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	scene, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	})
	assert.NoError(t, err)

	result, err := CalcAdjacentPay(scene, pt, 1, func(s int) bool { return s >= 0 }, func(s int) bool { return s == 0 }, func(cs, ss int) bool { return cs == ss }, func(s int) int { return s })
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func Test_calcAdjacentPay_InvalidSymbol(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	scene, err := NewGameSceneWithArr2([][]int{
		{-1, -1, -1},
		{-1, -1, -1},
		{-1, -1, -1},
	})
	assert.NoError(t, err)

	result, err := CalcAdjacentPay(scene, pt, 1, func(s int) bool { return s >= 0 }, func(s int) bool { return s == 0 }, func(cs, ss int) bool { return cs == ss }, func(s int) int { return s })
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func Test_calcAdjacentPayWithX_WildTakeOver(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// 3 wilds pay more than 5 of symbol 1
	pt.MapPay[0][2] = 1000
	pt.MapPay[1][4] = 10

	scene, err := NewGameSceneWithArr2([][]int{
		{0},
		{0},
		{0},
		{1},
		{1},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithX(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Symbol)
	assert.Equal(t, 3, result.SymbolNums)
}

func Test_CalcAdjacentPay_Isprocx(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	scene, err := NewGameSceneWithArr2([][]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 6, 7},
	})
	assert.NoError(t, err)

	// It should only find the vertical win, as the symbols are marked as processed.
	result, err := CalcAdjacentPay(scene, pt, 1, func(s int) bool { return s >= 0 }, func(s int) bool { return s == 0 }, func(cs, ss int) bool { return cs == ss }, func(s int) int { return s })
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result))
}

func Test_calcAdjacentPayWithX_NoWildPay(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	pt.MapPay[0][0] = 0 // 1 wild doesn't pay

	scene, err := NewGameSceneWithArr2([][]int{
		{0, 1, 2},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithX(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.Nil(t, result)
}

func Test_calcAdjacentPayWithX_NoRegPay(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	pt.MapPay[1][1] = 0 // 2 of symbol 1 don't pay
	pt.MapPay[0][0] = 0 // 1 wild doesn't pay

	scene, err := NewGameSceneWithArr2([][]int{
		{0, 1, 1},
	})
	assert.NoError(t, err)

	result := calcAdjacentPayWithX(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.Nil(t, result)
}

func Test_calcAdjacentPayWithY_WildLogic(t *testing.T) {
	pt, err := LoadPayTables5JSON("../unittestdata/paytables.json")
	assert.NoError(t, err)

	// wild > regular
	pt.MapPay[0][0] = 100
	pt.MapPay[1][1] = 10
	scene, err := NewGameSceneWithArr2([][]int{{0, 1, 1}})
	assert.NoError(t, err)
	result := calcAdjacentPayWithY(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.Symbol)

	// regular > wild
	pt.MapPay[0][0] = 10
	pt.MapPay[1][1] = 100
	result = calcAdjacentPayWithY(scene, 0, 0, 0, pt, 1, func(cs, ss int) bool { return cs == ss || cs == 0 }, func(cs int) bool { return cs == 0 })
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.Symbol)
}
