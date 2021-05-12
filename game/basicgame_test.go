package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_BasicGame(t *testing.T) {
	game := NewBasicGame(func() sgc7plugin.IPlugin {
		return sgc7plugin.NewBasicPlugin()
	})

	assert.NotNil(t, game.GetConfig(), "Test_BasicGame GetConfig")

	p0 := game.NewPlugin()
	assert.NotNil(t, p0, "Test_BasicGame NewPlugin")

	game.FreePlugin(p0)

	ps := game.Initialize()
	assert.NotNil(t, ps, "Test_BasicGame Initialize")
	bps, isok := ps.(*BasicPlayerState)
	assert.Equal(t, isok, true, "Test_BasicGame BasicPlayerState")
	assert.Equal(t, bps.Public.CurGameMod, "bg", "Test_BasicGame BasicPlayerPublicState CurGameMod")

	bg := NewBasicGameMod("bg", 5, 3)
	err := game.AddGameMod(bg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod bg")

	fg := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(fg)
	assert.Nil(t, err, "Test_BasicGame AddGameMod fg")

	fg1 := NewBasicGameMod("fg", 6, 4)
	err = game.AddGameMod(fg1)
	assert.Equal(t, err, ErrDuplicateGameMod, "Test_BasicGame AddGameMod fg1")

	var igame IGame = game
	// igame = game
	assert.NotNil(t, igame, "Test_BasicGame IGame")

	t.Logf("Test_BasicGame OK")
}
