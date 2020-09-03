package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicGameMod(t *testing.T) {
	mod := NewBasicGameMod("bg", 5, 3)

	assert.Equal(t, mod.GetName(), "bg", "Test_BasicGameMod GetName")

	// assert.NotNil(t, mod.GetGameScene(), "Test_BasicGameMod GetGameScene")

	assert.Equal(t, mod.Width, 5, "Test_BasicGameMod GetGameScene Width")
	assert.Equal(t, mod.Height, 3, "Test_BasicGameMod GetGameScene Height")

	var gamemod IGameMod
	gamemod = mod
	assert.NotNil(t, gamemod, "Test_BasicGameMod IGameMod")

	t.Logf("Test_BasicGameMod OK")
}
