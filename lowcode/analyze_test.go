package lowcode

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_ParseStepParentChildren(t *testing.T) {
	// SetJsonMode()

	game, err := NewGame2("../unittestdata/testgame.json", func() sgc7plugin.IPlugin {
		return sgc7plugin.NewFastPlugin()
	}, NewBasicRNG, NewEmptyFeatureLevel)
	assert.NoError(t, err)
	assert.NotNil(t, game)

	bet := 20
	components := game.Pool.mapComponents[bet]
	// gameProp := game.Pool.newGameProp(bet)

	node, err := ParseStepParentChildren(components, game.Pool.Config.MapBetConfigs[bet].Start)
	assert.NoError(t, err)
	assert.NotNil(t, node)
	assert.True(t, node.CountComponentNum() == 24)
	assert.True(t, node.CountDeep() == 2)
	assert.True(t, node.CountParentNum() == 2)

	t.Logf("Test_ParseStepParentChildren OK")
}
