package gatiserv

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
)

func httpGet(url string) (int, []byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, nil
}

//--------------------------------------------------------------------------------------
// testPlayerState

type testPlayerState struct {
}

// SetPublic - set player public state
func (ps *testPlayerState) SetPublic(pub interface{}) error {
	return nil
}

// SetPrivate - set player private state
func (ps *testPlayerState) SetPrivate(pri interface{}) error {
	return nil
}

// GetPublic - get player public state
func (ps *testPlayerState) GetPublic() interface{} {
	return sgc7game.BasicPlayerPublicState{
		CurGameMod: "BG",
	}
}

// GetPrivate - get player private state
func (ps *testPlayerState) GetPrivate() interface{} {
	return sgc7game.BasicPlayerPrivateState{}
}

//--------------------------------------------------------------------------------------
// testService

type testService struct {
	cfg      *sgc7game.Config
	initmode int
}

// Config - get configuration
func (sv *testService) Config() *sgc7game.Config {
	return sv.cfg
}

// Initialize - initialize a player
func (sv *testService) Initialize() sgc7game.IPlayerState {
	if sv.initmode == 0 {
		return nil
	}

	return &testPlayerState{}
}

func Test_Serv(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	cfg := &Config{
		GameID:      "1019",
		BindAddr:    "127.0.0.1:7891",
		IsDebugMode: true,
	}

	service := &testService{
		&sgc7game.Config{
			Width:  5,
			Height: 3,
		},
		0,
	}

	serv := NewServ(service, cfg)

	go func() {
		err := serv.Start()
		if err != nil {
			t.Fatalf("Test_Serv Start error %v",
				err)
		}
	}()

	time.Sleep(time.Second * 3)

	sc, buff, err := httpGet("http://127.0.0.1:7891/v2/games/1019/config")
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	rr := &sgc7game.Config{}
	err = json.Unmarshal(buff, rr)
	if err != nil {
		t.Fatalf("Test_Serv Unmarshal error %v",
			err)
	}

	assert.Equal(t, rr.Width, 5, "they should be equal")
	assert.Equal(t, rr.Height, 3, "they should be equal")

	sc, buff, err = httpGet("http://127.0.0.1:7891/v2/games/1019/initialize")
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}", "they should be equal")

	service.initmode = 1

	sc, buff, err = httpGet("http://127.0.0.1:7891/v2/games/1019/initialize")
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"playerStatePublic\":{\"CurGameMod\":\"BG\"},\"playerStatePrivate\":{}}", "they should be equal")

	serv.Stop()

	t.Logf("Test_Serv OK")
}
