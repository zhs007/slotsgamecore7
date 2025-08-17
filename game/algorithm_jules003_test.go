package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CheckLineRL_Jules(t *testing.T) {
	// Symbol 0 is wild
	// Symbol -1 is invalid
	isWild := func(cursymbol int) bool {
		return cursymbol == 0
	}

	isValidSymbol := func(cursymbol int) bool {
		return cursymbol >= 0
	}

	isSameSymbol := func(cursymbol int, startsymbol int) bool {
		return cursymbol == startsymbol || cursymbol == 0
	}

	getSymbol := func(cursymbol int) int {
		return cursymbol
	}

	type testCase struct {
		name          string
		scene         *GameScene
		line          []int
		minnum        int
		expectedWin   bool
		expectedSym   int
		expectedNums  int
		expectedWilds int
		expectedPos   []int
	}

	testCases := []testCase{
		{
			name: "Simple win, no wilds",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 1, 1},
					{8, 1, 2},
					{9, 1, 3},
				},
			},
			line:          []int{0, 0, 1, 1, 1},
			minnum:        3,
			expectedWin:   true,
			expectedSym:   1,
			expectedNums:  3,
			expectedWilds: 0,
			expectedPos:   []int{4, 1, 3, 1, 2, 1},
		},
		{
			name: "No win, not enough symbols",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 1},
					{9, 1, 2},
					{9, 3, 3},
				},
			},
			line:        []int{0, 0, 1, 1, 1},
			minnum:      3,
			expectedWin: false,
		},
		{
			name: "Win with wilds",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 0, 1}, // wild
					{8, 1, 2},
					{9, 1, 3},
				},
			},
			line:          []int{0, 0, 1, 1, 1},
			minnum:        3,
			expectedWin:   true,
			expectedSym:   1,
			expectedNums:  3,
			expectedWilds: 1,
			expectedPos:   []int{4, 1, 3, 1, 2, 1},
		},
		{
			name: "Win starts with wild (from right)",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 1, 1},
					{8, 1, 2},
					{9, 0, 3}, // wild
				},
			},
			line:          []int{0, 0, 1, 1, 1},
			minnum:        3,
			expectedWin:   true,
			expectedSym:   1,
			expectedNums:  3,
			expectedWilds: 1,
			expectedPos:   []int{4, 1, 3, 1, 2, 1},
		},
		{
			name: "Win all wilds",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{-1, -1, -1}, // Use invalid symbols to ensure the line breaks after wilds
					{7, 0, 1},    // Wild
					{8, 0, 2},    // Wild
					{9, 0, 3},    // Wild
				},
			},
			line:          []int{0, 0, 1, 1, 1},
			minnum:        3,
			expectedWin:   true,
			expectedSym:   0, // wild symbol
			expectedNums:  3,
			expectedWilds: 3,
			expectedPos:   []int{4, 1, 3, 1, 2, 1},
		},
		{
			name: "Special case: Wilds at start, non-wild determines symbol",
			// C B W W W -> reading from right: W W W B C
			scene: &GameScene{
				Arr: [][]int{
					{5, 2, 3}, // C
					{4, 2, 6}, // B
					{7, 0, 1}, // W
					{8, 0, 2}, // W
					{9, 0, 3}, // W
				},
			},
			line:          []int{0, 1, 1, 1, 1},
			minnum:        4,
			expectedWin:   true,
			expectedSym:   2, // Should be symbol B
			expectedNums:  4,
			expectedWilds: 3,
			expectedPos:   []int{4, 1, 3, 1, 2, 1, 1, 1},
		},
		{
			name: "Invalid start symbol",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 1, 1},
					{8, 1, 2},
					{9, -1, 3}, // invalid
				},
			},
			line:        []int{0, 0, 1, 1, 1},
			minnum:      3,
			expectedWin: false,
		},
		{
			name: "Line broken by non-matching symbol",
			scene: &GameScene{
				Arr: [][]int{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 1}, // non-matching
					{9, 1, 2},
					{9, 1, 3},
				},
			},
			line:        []int{0, 0, 1, 1, 1},
			minnum:      3,
			expectedWin: false, // Only 2 symbols match from right
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CheckLineRL(tc.scene, tc.line, tc.minnum, isValidSymbol, isWild, isSameSymbol, getSymbol)

			if !tc.expectedWin {
				assert.Nil(t, result, "Expected no win, but got one")
				return
			}

			assert.NotNil(t, result, "Expected a win, but got none")
			assert.Equal(t, tc.expectedSym, result.Symbol, "Winning symbol mismatch")
			assert.Equal(t, tc.expectedNums, result.SymbolNums, "Symbol count mismatch")
			assert.Equal(t, tc.expectedWilds, result.Wilds, "Wild count mismatch")
			assert.Equal(t, tc.expectedPos, result.Pos, "Positions mismatch")
		})
	}
}
