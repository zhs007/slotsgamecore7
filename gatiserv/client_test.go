package gatiserv

import (
	"net/http"
	"testing"

	"github.com/bytedance/sonic"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func Test_ClientGetConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const URL = "http://127.0.0.1:7891/v2/games/1019/"
	const configURL = URL + "config"

	client := NewClient(URL)

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(404, "")(req)
		})

	cfg, err := client.GetConfig()
	assert.Equal(t, err, ErrNonStatusOK, "Test_ClientGetConfig GetConfig")
	assert.Nil(t, cfg, "Test_ClientGetConfig GetConfig")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return nil, ErrNonStatusOK
		})

	cfg, err = client.GetConfig()
	assert.NotNil(t, err, "Test_ClientGetConfig GetConfig")
	assert.Nil(t, cfg, "Test_ClientGetConfig GetConfig")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(200, "")(req)
		})

	cfg, err = client.GetConfig()
	assert.NotNil(t, err, "Test_ClientGetConfig GetConfig")
	assert.Nil(t, cfg, "Test_ClientGetConfig GetConfig")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			rcfg := &sgc7game.Config{}
			resbuff, err := sonic.Marshal(rcfg)
			assert.Nil(t, err, "Test_ClientGetConfig Marshal")

			return httpmock.NewStringResponder(200, string(resbuff))(req)
		})

	cfg, err = client.GetConfig()
	assert.Nil(t, err, "Test_ClientGetConfig GetConfig")
	assert.NotNil(t, cfg, "Test_ClientGetConfig GetConfig")

	t.Logf("Test_ClientGetConfig OK")
}

func Test_ClientInitialize(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const URL = "http://127.0.0.1:7891/v2/games/1019/"
	const configURL = URL + "initialize"

	client := NewClient(URL)

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(404, "")(req)
		})

	ps := &PlayerState{}

	err := client.Initialize(ps)
	assert.Equal(t, err, ErrNonStatusOK, "Test_ClientInitialize Initialize")
	// assert.Nil(t, ps, "Test_ClientInitialize Initialize")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return nil, ErrNonStatusOK
		})

	err = client.Initialize(ps)
	assert.NotNil(t, err, "Test_ClientInitialize Initialize")
	// assert.Nil(t, ps, "Test_ClientInitialize Initialize")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(200, "")(req)
		})

	err = client.Initialize(ps)
	assert.NotNil(t, err, "Test_ClientInitialize Initialize")
	// assert.Nil(t, ps, "Test_ClientInitialize Initialize")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			ps := &PlayerState{}
			resbuff, err := sonic.Marshal(ps)
			assert.Nil(t, err, "Test_ClientInitialize Marshal")

			return httpmock.NewStringResponder(200, string(resbuff))(req)
		})

	err = client.Initialize(ps)
	assert.Nil(t, err, "Test_ClientInitialize Initialize")
	assert.NotNil(t, ps, "Test_ClientInitialize Initialize")

	httpmock.RegisterResponder("GET",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(200, "{\"playerStatePublic\":{\"curgamemod\":\"BG\"},\"playerStatePrivate\":{}}")(req)
		})

	ps = &PlayerState{
		Public:  &sgc7game.BasicPlayerPublicState{},
		Private: &sgc7game.BasicPlayerPrivateState{},
	}
	err = client.Initialize(ps)
	assert.Nil(t, err, "Test_ClientInitialize Initialize")
	assert.NotNil(t, ps, "Test_ClientInitialize Initialize")

	bps, isok := ps.Public.(*sgc7game.BasicPlayerPublicState)
	assert.Equal(t, isok, true)
	assert.Equal(t, bps.CurGameMod, "BG")

	t.Logf("Test_ClientInitialize OK")
}

func Test_ClientPlayEx(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const URL = "http://127.0.0.1:7891/v2/games/1019/"
	const configURL = URL + "play"

	client := NewClient(URL)

	httpmock.RegisterResponder("POST",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(404, "")(req)
		})

	pr, err := client.PlayEx("")
	assert.Equal(t, err, ErrNonStatusOK, "Test_ClientPlayEx PlayEx")
	assert.Nil(t, pr, "Test_ClientPlayEx PlayEx")

	httpmock.RegisterResponder("POST",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return nil, ErrNonStatusOK
		})

	pr, err = client.PlayEx("")
	assert.NotNil(t, err, "Test_ClientPlayEx PlayEx")
	assert.Nil(t, pr, "Test_ClientPlayEx PlayEx")

	httpmock.RegisterResponder("POST",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponder(200, "")(req)
		})

	pr, err = client.PlayEx("")
	assert.NotNil(t, err, "Test_ClientPlayEx PlayEx")
	assert.Nil(t, pr, "Test_ClientPlayEx PlayEx")

	httpmock.RegisterResponder("POST",
		configURL,
		func(req *http.Request) (*http.Response, error) {
			pr := &PlayResult{}
			resbuff, err := sonic.Marshal(pr)
			assert.Nil(t, err, "Test_ClientPlayEx Marshal")

			return httpmock.NewStringResponder(200, string(resbuff))(req)
		})

	pr, err = client.PlayEx("")
	assert.Nil(t, err, "Test_ClientPlayEx PlayEx")
	assert.NotNil(t, pr, "Test_ClientPlayEx PlayEx")

	t.Logf("Test_ClientPlayEx OK")
}
