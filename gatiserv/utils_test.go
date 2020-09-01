package gatiserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_BuildIPlayerState(t *testing.T) {
	ips := &sgc7game.BasicPlayerState{}

	err := BuildIPlayerState(ips, PlayerState{
		Public:  "{\"curgamemod\":\"BG\"}",
		Private: "{}",
	})
	assert.Nil(t, err)

	pbs, isok := ips.GetPublic().(sgc7game.BasicPlayerPublicState)
	assert.Equal(t, isok, true)
	assert.Equal(t, pbs.CurGameMod, "BG")

	pps, isok := ips.GetPrivate().(sgc7game.BasicPlayerPrivateState)
	assert.Equal(t, isok, true)
	assert.NotNil(t, pps)

	err = BuildIPlayerState(ips, PlayerState{
		Public:  "",
		Private: "",
	})
	assert.NotNil(t, err)

	err = BuildIPlayerState(ips, PlayerState{
		Public:  "{}",
		Private: "",
	})
	assert.NotNil(t, err)

	err = BuildIPlayerState(ips, PlayerState{
		Public:  "{}",
		Private: "{}",
	})
	assert.Nil(t, err)

	pbs, isok = ips.GetPublic().(sgc7game.BasicPlayerPublicState)
	assert.Equal(t, isok, true)
	assert.NotNil(t, pbs)

	pps, isok = ips.GetPrivate().(sgc7game.BasicPlayerPrivateState)
	assert.Equal(t, isok, true)
	assert.NotNil(t, pps)

	t.Logf("Test_BuildIPlayerState OK")
}

func Test_BuildStake(t *testing.T) {
	bs := BuildStake(Stake{
		0.01,
		1,
		"EUR",
	})

	assert.Equal(t, bs.CoinBet, int64(1))
	assert.Equal(t, bs.CashBet, int64(100))
	assert.Equal(t, bs.Currency, "EUR")

	bs = BuildStake(Stake{
		0.004,
		100.125,
		"USD",
	})

	assert.Equal(t, bs.CoinBet, int64(0))
	assert.Equal(t, bs.CashBet, int64(10012))
	assert.Equal(t, bs.Currency, "USD")

	t.Logf("Test_BuildStake OK")
}
