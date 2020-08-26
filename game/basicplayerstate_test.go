package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicPlayerState(t *testing.T) {
	ps := NewBasicPlayerState("BG")
	assert.NotNil(t, ps, "Test_BasicPlayerState NewBasicPlayerState")

	assert.Equal(t, ps.Public.CurGameMod, "BG", "Test_BasicPlayerState BasicPlayerPublicState CurGameMod")

	ps.SetPublic(&BasicPlayerPublicState{CurGameMod: "FG"})
	assert.Equal(t, ps.Public.CurGameMod, "FG", "Test_BasicPlayerState SetPublic CurGameMod")

	ps.SetPrivate(&BasicPlayerPrivateState{})

	var ips IPlayerState
	ips = ps
	assert.NotNil(t, ips, "Test_BasicPlayerState IPlayerState")

	t.Logf("Test_BasicPlayerState OK")
}
