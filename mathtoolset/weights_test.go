package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_forEachArrWithLength(t *testing.T) {
	num0 := 0
	forEachArrWithLength(nil, []int{1, 2, 3}, 1, func([]int) {
		num0++
	})
	assert.Equal(t, num0, 3)

	num1 := 0
	forEachArrWithLength(nil, []int{1, 2, 3}, 2, func([]int) {
		num1++
	})
	assert.Equal(t, num1, 3)

	num2 := 0
	forEachArrWithLength(nil, []int{1, 2, 3}, 3, func([]int) {
		num2++
	})
	assert.Equal(t, num2, 1)

	num3 := 0
	forEachArrWithLength(nil, []int{1, 2, 3, 4}, 1, func([]int) {
		num3++
	})
	assert.Equal(t, num3, 4)

	num4 := 0
	forEachArrWithLength(nil, []int{1, 2, 3, 4}, 2, func([]int) {
		num4++
	})
	assert.Equal(t, num4, 6)

	num5 := 0
	forEachArrWithLength(nil, []int{1, 2, 3, 4}, 3, func([]int) {
		num5++
	})
	assert.Equal(t, num5, 4)

	num6 := 0
	forEachArrWithLength(nil, []int{1, 2, 3, 4}, 4, func([]int) {
		num6++
	})
	assert.Equal(t, num6, 1)

	num7 := 0
	forEachArrWithLength(nil, []int{1, 2, 3, 4}, 5, func([]int) {
		num7++
	})
	assert.Equal(t, num7, 0)

	t.Logf("Test_forEachArrWithLength OK")
}
