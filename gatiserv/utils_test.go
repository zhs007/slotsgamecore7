package gatiserv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

//--------------------------------------------------------------------------------------
// errPlayerState

type errPlayerState struct {
}

// SetPublic - set player public state
func (eps *errPlayerState) SetPublic(pub interface{}) error {
	return nil

}

// SetPrivate - set player private state
func (eps *errPlayerState) SetPrivate(pri interface{}) error {
	return nil
}

// SetPublicString - set player public state
func (eps *errPlayerState) SetPublicString(pub string) error {
	return nil
}

// SetPrivateString - set player private state
func (eps *errPlayerState) SetPrivateString(pri string) error {
	return nil
}

// GetPublic - get player public state
func (eps *errPlayerState) GetPublic() interface{} {
	return nil
}

// GetPrivate - get player private state
func (eps *errPlayerState) GetPrivate() interface{} {
	return nil
}

func Test_BuildIPlayerState(t *testing.T) {
	ips := &sgc7game.BasicPlayerState{}

	err := BuildIPlayerState(ips, PlayerState{
		Public:  "{\"curgamemod\":\"BG\"}",
		Private: "{}",
	})
	assert.Nil(t, err)

	assert.Equal(t, ips.Public.CurGameMod, "BG")

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

func Test_BuildPlayerStateString(t *testing.T) {
	str, err := BuildPlayerStateString(nil)
	assert.Nil(t, err)
	assert.Equal(t, str, "{\"playerStatePublic\":\"{}\",\"playerStatePrivate\":\"{}\"}")

	ps, err := ParsePlayerState(str)
	assert.Nil(t, err)
	assert.NotNil(t, ps)

	ips := sgc7game.NewBasicPlayerState("BG")
	str, err = BuildPlayerStateString(ips)
	assert.Nil(t, err)
	assert.Equal(t, str, "{\"playerStatePublic\":\"{\\\"curgamemod\\\":\\\"BG\\\"}\",\"playerStatePrivate\":\"{}\"}")

	ps, err = ParsePlayerState(str)
	assert.Nil(t, err)
	assert.NotNil(t, ps)

	ps, err = ParsePlayerState("")
	assert.NotNil(t, err)
	assert.Nil(t, ps)

	// eps := &errPlayerState{}
	// str, err = BuildPlayerStateString(eps)
	// assert.Nil(t, err)
	// assert.Equal(t, str, "{\"playerStatePublic\":\"{\\\"curgamemod\\\":\\\"BG\\\"}\",\"playerStatePrivate\":\"{}\"}")

	t.Logf("Test_BuildPlayerStateString OK")
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

func Test_ParsePlayParams(t *testing.T) {
	pp, err := ParsePlayParams("{\"stakeValue\":{\"coinBet\":0.0500,\"cashBet\":1.00,\"currency\":\"EUR\"},\"playerState\":{\"playerStatePublic\":\"{\\\"curgamemod\\\":\\\"BG\\\"}\",\"playerStatePrivate\":\"{}\"},\"clientParams\":null,\"cheat\":\"\",\"command\":\"\",\"freespinsActive\":false,\"jackpotStakeValue\":null,\"jackpotValues\":null}")
	assert.Nil(t, err)
	assert.NotNil(t, pp)

	assert.Equal(t, pp.Stake.CoinBet, 0.05)
	assert.Equal(t, pp.Stake.CashBet, 1.0)
	assert.Equal(t, pp.Stake.Currency, "EUR")

	bps := sgc7game.NewBasicPlayerStateEx(pp.PlayerState.Public, pp.PlayerState.Private)
	assert.Equal(t, bps.Public.CurGameMod, "BG")

	assert.Equal(t, pp.Params, "")

	assert.Equal(t, pp.Cheat, "")

	assert.Equal(t, pp.Cmd, "")

	assert.Equal(t, pp.FreespinsActive, false)

	assert.Equal(t, pp.JackpotStakeValue, float64(0))

	assert.Nil(t, pp.JackpotValues)

	pp, err = ParsePlayParams("")
	assert.NotNil(t, err)
	assert.Nil(t, pp)

	t.Logf("Test_ParsePlayParams OK")
}