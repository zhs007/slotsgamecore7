package mathtoolset2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ParseExRule(t *testing.T) {
	rule, err := ParseExRule("OFF_3,SC,WL")
	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.Equal(t, "OFF", rule.Code)
	assert.Equal(t, 1, len(rule.Params))
	assert.Equal(t, 3, rule.Params[0])
	assert.Equal(t, 2, len(rule.Symbols))
	assert.Equal(t, "SC", rule.Symbols[0])
	assert.Equal(t, "WL", rule.Symbols[1])

	t.Logf("Test_ParseExRule OK")
}
