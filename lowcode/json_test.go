package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_NewGame2(t *testing.T) {
	SetJsonMode()

	game, err := NewGame2("../data/game002.json", func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	})
	assert.NoError(t, err)
	assert.NotNil(t, game)

	t.Logf("Test_NewGame2 OK")
}
