package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicGame(t *testing.T) {
	game := NewBasicGame()

	assert.NotNil(t, game.GetConfig(), "Test_BasicGame GetConfig")

	assert.Nil(t, game.GetPlugin(), "Test_BasicGame GetPlugin")

	ps := game.Initialize()
	assert.NotNil(t, ps, "Test_BasicGame Initialize")
	bps, isok := ps.(*BasicPlayerState)
	assert.Equal(t, isok, true, "Test_BasicGame BasicPlayerState")
	assert.Equal(t, bps.Public.CurGameMod, "BG", "Test_BasicGame BasicPlayerPublicState CurGameMod")

	bg := NewBasicGameMod("bg", 5, 3)
	err := game.AddGameMod(&bg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod bg")

	fg := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(&fg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod fg")

	fg1 := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(&fg1)
	assert.Equal(t, err, ErrDuplicateGameMod, "Test_BasicGame AddGameMod fg1")

	var igame IGame
	igame = &game
	assert.NotNil(t, igame, "Test_BasicGame IGame")

	t.Logf("Test_BasicGame OK")
}
