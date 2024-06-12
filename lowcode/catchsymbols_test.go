package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_mergePosWithoutSelected(t *testing.T) {
	srcpos1 := []int{0, 0, 1, 1, 2, 2}
	targetpos1 := []int{3, 3, 4, 4, 5, 5}

	pos1 := mergePosWithoutSelected(srcpos1, targetpos1, 0)
	assert.Equal(t, len(pos1), 10)
	assert.Equal(t, pos1[5], 5)
	assert.Equal(t, pos1[6], 1)
	assert.Equal(t, pos1[7], 1)
	assert.Equal(t, pos1[8], 2)

	srcpos2 := []int{0, 0, 1, 1, 2, 2}
	targetpos2 := []int{3, 3, 4, 4, 5, 5}

	pos2 := mergePosWithoutSelected(srcpos2, targetpos2, 1)
	assert.Equal(t, len(pos2), 10)
	assert.Equal(t, pos2[5], 5)
	assert.Equal(t, pos2[6], 0)
	assert.Equal(t, pos2[7], 0)
	assert.Equal(t, pos2[8], 2)

	srcpos3 := []int{0, 0, 1, 1, 2, 2}
	targetpos3 := []int{3, 3, 4, 4, 5, 5}

	pos3 := mergePosWithoutSelected(srcpos3, targetpos3, 2)
	assert.Equal(t, len(pos3), 10)
	assert.Equal(t, pos3[5], 5)
	assert.Equal(t, pos3[6], 0)
	assert.Equal(t, pos3[7], 0)
	assert.Equal(t, pos3[8], 1)

	t.Logf("Test_mergePosWithoutSelected OK")
}

func Test_findNearest(t *testing.T) {
	targetpos := []int{2, 0, 3, 3, 5, 5, 6, 6}

	pos1 := findNearest(0, 0, targetpos)
	assert.Equal(t, pos1, 0)

	pos2 := findNearest(4, 4, targetpos)
	assert.Equal(t, pos2, 1)

	t.Logf("Test_findNearest OK")
}
