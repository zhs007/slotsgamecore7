package mathtoolset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../unittestdata/genmathcfg.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	t.Logf("Test_LoadConfig OK")
}
