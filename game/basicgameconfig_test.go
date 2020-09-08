package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_BasicGameConfig(t *testing.T) {
	cfg := BasicGameConfig{}

	err := LoadGameConfig("../unittestdata/gamecfg.yaml", &cfg)
	assert.NoError(t, err)

	game := NewBasicGame(func() sgc7plugin.IPlugin {
		return sgc7plugin.NewBasicPlugin()
	})

	err = cfg.Init5(game)
	assert.NoError(t, err)

	t.Logf("Test_BasicGameConfig OK")
}
