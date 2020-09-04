package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicGameConfig(t *testing.T) {
	cfg := BasicGameConfig{}

	err := LoadGameConfig("../unittestdata/gamecfg.yaml", &cfg)
	assert.NoError(t, err)

	game := NewBasicGame()

	err = cfg.Init5(game)
	assert.NoError(t, err)

	t.Logf("Test_BasicGameConfig OK")
}
