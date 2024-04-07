package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PowInt(t *testing.T) {

	t0 := PowInt(2, 0)
	assert.Equal(t, t0, 1)

	t1 := PowInt(2, -1)
	assert.Equal(t, t1, 1)

	t2 := PowInt(2, 1)
	assert.Equal(t, t2, 2)

	t3 := PowInt(2, 2)
	assert.Equal(t, t3, 4)

	t4 := PowInt(2, 3)
	assert.Equal(t, t4, 8)

	t5 := PowInt(2, 4)
	assert.Equal(t, t5, 16)

	t.Logf("Test_PowInt OK")
}
