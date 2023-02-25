package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../data/game001/rtp96.yaml")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	t.Logf("Test_LoadConfig OK")
}
