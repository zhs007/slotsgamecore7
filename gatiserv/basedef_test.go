package gatiserv

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

func Test_PlayParams(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	pp := &PlayParams{}
	err := json.Unmarshal([]byte("{\"cheat\":\"\",\"command\":\"\",\"freespinsActive\":false,\"playerState\":{},\"stakeValue\":{\"cashBet\":3.00,\"coinBet\":0.1000,\"currency\":\"EUR\"}}"), pp)
	assert.NoError(t, err)

	t.Logf("Test_PlayParams OK")
}
