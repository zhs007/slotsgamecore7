package sgc7game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasicPlayerState(t *testing.T) {
	ps := NewBasicPlayerState("BG")
	assert.NotNil(t, ps, "Test_BasicPlayerState NewBasicPlayerState")

	assert.Equal(t, ps.Public.CurGameMod, "BG", "Test_BasicPlayerState BasicPlayerPublicState CurGameMod")

	ps.SetPublic(BasicPlayerPublicState{CurGameMod: "FG"})
	assert.Equal(t, ps.Public.CurGameMod, "FG", "Test_BasicPlayerState SetPublic CurGameMod")

	ps.SetPrivate(BasicPlayerPrivateState{})

	// ipspub := ps.GetPublic()
	// bppub, isok := ipspub.(BasicPlayerPublicState)
	// assert.Equal(t, isok, true, "Test_BasicPlayerState BasicPlayerPublicState")
	// assert.Equal(t, bppub.CurGameMod, "FG", "Test_BasicPlayerState BasicPlayerPublicState CurGameMod")

	// ipspri := ps.GetPrivate()
	// bppri, isok := ipspri.(BasicPlayerPrivateState)
	// assert.Equal(t, isok, true, "Test_BasicPlayerState BasicPlayerPrivateState")
	// assert.NotNil(t, bppri, "Test_BasicPlayerState BasicPlayerPrivateState")

	var ips IPlayerState
	ips = ps
	assert.NotNil(t, ips, "Test_BasicPlayerState IPlayerState")

	ipspub := ips.GetPublic()
	bppub, isok := ipspub.(BasicPlayerPublicState)
	assert.Equal(t, isok, true, "Test_BasicPlayerState BasicPlayerPublicState")
	assert.Equal(t, bppub.CurGameMod, "FG", "Test_BasicPlayerState BasicPlayerPublicState CurGameMod")

	ipspri := ips.GetPrivate()
	bppri, isok := ipspri.(BasicPlayerPrivateState)
	assert.Equal(t, isok, true, "Test_BasicPlayerState BasicPlayerPrivateState")
	assert.NotNil(t, bppri, "Test_BasicPlayerState BasicPlayerPrivateState")

	t.Logf("Test_BasicPlayerState OK")
}
