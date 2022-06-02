package relaxutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Config(t *testing.T) {
	cfg := &Config{
		General: "General game information",
		SD:      1.56,
	}

	err := SaveConfig("../unittestdata/relaxconfig.xml", cfg)
	assert.NoError(t, err)

	t.Logf("Test_CalcScatter OK")
}
