package sgc7utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsFloatEquals(t *testing.T) {
	assert.Equal(t, IsFloatEquals(0.123456789, 0.123456780), true)
	assert.Equal(t, IsFloatEquals(0.1234567, 0.12345678), false)

	t.Logf("Test_IsFloatEquals OK")
}
