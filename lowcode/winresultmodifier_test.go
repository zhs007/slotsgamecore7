package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseWinResultModifierType(t *testing.T) {
	cases := map[string]WinResultModifierType{
		"existsymbol":       WRMTypeExistSymbol,
		"addsymbolmulti":    WRMTypeAddSymbolMulti,
		"mulsymbolmulti":    WRMTypeMulSymbolMulti,
		"symbolmultionways": WRMTypeSymbolMultiOnWays,
		"divide":            WRMTypeDivide,
		"multiply":          WRMTypeMultiply,
		"unknown":           WRMTypeExistSymbol, // default
	}

	for k, v := range cases {
		got := parseWinResultModifierType(k)
		assert.Equal(t, v, got, "parseWinResultModifierType(%s)", k)
	}
}

func TestWinResultModifierTypeMethods(t *testing.T) {
	// isNeedMultiply: all except Divide
	assert.True(t, WRMTypeExistSymbol.isNeedMultiply())
	assert.True(t, WRMTypeAddSymbolMulti.isNeedMultiply())
	assert.True(t, WRMTypeMulSymbolMulti.isNeedMultiply())
	assert.True(t, WRMTypeSymbolMultiOnWays.isNeedMultiply())
	assert.False(t, WRMTypeDivide.isNeedMultiply())
	assert.True(t, WRMTypeMultiply.isNeedMultiply())

	// isNeedGameScene only for some types
	assert.True(t, WRMTypeExistSymbol.isNeedGameScene())
	assert.True(t, WRMTypeAddSymbolMulti.isNeedGameScene())
	assert.True(t, WRMTypeMulSymbolMulti.isNeedGameScene())
	assert.True(t, WRMTypeSymbolMultiOnWays.isNeedGameScene())
	assert.False(t, WRMTypeDivide.isNeedGameScene())
	assert.False(t, WRMTypeMultiply.isNeedGameScene())
}

func TestGetWinMulti(t *testing.T) {
	wrm := &WinResultModifier{
		Config: &WinResultModifierConfig{
			WinMulti: 5,
		},
	}

	// basicCD with no CCVWinMulti -> use Config.WinMulti
	basicCD := &BasicComponentData{}
	got := wrm.GetWinMulti(basicCD)
	assert.Equal(t, 5, got)

	// basicCD with CCVWinMulti set to positive value
	basicCD.MapConfigIntVals = map[string]int{CCVWinMulti: 3}
	got = wrm.GetWinMulti(basicCD)
	assert.Equal(t, 3, got)

	// basicCD with CCVWinMulti set to 0 or negative -> returns 1
	basicCD.MapConfigIntVals[CCVWinMulti] = 0
	got = wrm.GetWinMulti(basicCD)
	assert.Equal(t, 1, got)

	basicCD.MapConfigIntVals[CCVWinMulti] = -10
	got = wrm.GetWinMulti(basicCD)
	assert.Equal(t, 1, got)
}
