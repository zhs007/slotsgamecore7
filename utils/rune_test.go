package sgc7utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Rune2Int(t *testing.T) {
	str := "0123456789"
	out := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	for i, v := range str {
		co := Rune2Int(v)
		assert.Equal(t, co, out[i])
	}

	t.Logf("Test_Rune2Int OK")
}
