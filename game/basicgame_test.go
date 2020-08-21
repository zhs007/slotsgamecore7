package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicGame(t *testing.T) {
	game := NewBasicGame()

	assert.NotNil(t, game.GetConfig(), "Test_BasicGame GetConfig")

	assert.Nil(t, game.GetPlugin(), "Test_BasicGame GetPlugin")

	assert.Nil(t, game.Initialize(), "Test_BasicGame Initialize")

	bg := NewBasicGameMod("bg", 5, 3)
	err := game.AddGameMod(&bg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod bg")

	fg := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(&fg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod fg")

	fg1 := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(&fg1)
	assert.Equal(t, err, ErrDuplicateGameMod, "Test_BasicGame AddGameMod fg1")

	t.Logf("Test_BasicGame OK")
}
