package gati

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
)

func Test_PluginGATI(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const URL = "http://127.0.0.1:50000/numbers"
	res := []int{123, 123, 123}

	resbuff, err := sonic.Marshal(res)
	assert.NoError(t, err)

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s?size=%d", URL, 3),
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.Header.Get("X-Game-ID"), "1019", "they should be equal")

			return httpmock.NewStringResponder(200, string(resbuff))(req)
		})

	bp := NewPluginGATI(&Config{
		GameID:  "1019",
		RNGURL:  URL,
		RngNums: 3,
	})

	var lstr []int

	for i := 0; i < 1000; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err, "Test_PluginGATI Random")
		assert.True(t, func() bool {
			return r >= 0 && r < 100
		}(), "Test_PluginGATI Random range")

		assert.Equal(t, r, 23, "Test_PluginGATI Random range")

		lstr = append(lstr, r)
	}

	lst := bp.GetUsedRngs()
	assert.NotNil(t, lst, "Test_PluginGATI GetUsedRngs")
	assert.Equal(t, len(lst), 1000, "Test_PluginGATI GetUsedRngs len")

	for i := 0; i < 1000; i++ {
		assert.Equal(t, lst[i].Value, lstr[i], "Test_PluginGATI GetUsedRngs value")
	}

	bp.ClearUsedRngs()

	lst1 := bp.GetUsedRngs()
	assert.Nil(t, lst1, "Test_PluginGATI GetUsedRngs")

	lstcache := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	bp.SetCache(lstcache)

	for i := 0; i < 10; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.Equal(t, r, lstcache[i], "Test_PluginGATI Random Cache value")
		assert.Equal(t, len(bp.Rngs), 9-i, "Test_PluginGATI Random ClearCache")
	}

	bp.SetCache(lstcache)

	for i := 0; i < 5; i++ {
		r, err := bp.Random(context.Background(), 100)
		assert.NoError(t, err)
		assert.Equal(t, r, lstcache[i], "Test_PluginGATI Random Cache value")
		assert.Equal(t, len(bp.Rngs), 9-i, "Test_PluginGATI Random ClearCache")
	}

	bp.ClearCache()
	assert.Equal(t, len(bp.Rngs), 0, "Test_PluginGATI Random ClearCache")

	var ip sgc7plugin.IPlugin
	ip = bp
	assert.NotNil(t, ip, "Test_PluginGATI IPlugin")

	t.Logf("Test_PluginGATI OK")
}
