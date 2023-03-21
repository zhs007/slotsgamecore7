package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BookOf(t *testing.T) {
	bookof := NewBookOf("bg-bookof")
	assert.NotNil(t, bookof)

	t.Logf("Test_BookOf OK")
}
