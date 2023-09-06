package gatiserv

import (
	"testing"

	"github.com/bytedance/sonic"
	"github.com/stretchr/testify/assert"
)

func Test_PlayParams(t *testing.T) {
	pp := &PlayParams{}
	err := sonic.Unmarshal([]byte("{\"cheat\":\"\",\"command\":\"\",\"freespinsActive\":false,\"playerState\":{},\"stakeValue\":{\"cashBet\":3.00,\"coinBet\":0.1000,\"currency\":\"EUR\"}}"), pp)
	assert.NoError(t, err)

	t.Logf("Test_PlayParams OK")
}
