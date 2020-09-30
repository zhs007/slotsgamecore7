package gatiserv

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

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

// SetPublicString - set player public state
func (ps *testPlayerState) SetPublicString(pub string) error {
	return nil
}

// SetPrivateString - set player private state
func (ps *testPlayerState) SetPrivateString(pri string) error {
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
	cfg        *sgc7game.Config
	initmode   int
	GameConfig *GATIGameConfig
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

// Validate - validate game
func (sv *testService) Validate(params *ValidateParams) []ValidationError {
	return nil
}

// Play - play game
func (sv *testService) Play(params *PlayParams) (*PlayResult, error) {
	return nil, nil
}

// Checksum - checksum
func (sv *testService) Checksum(lst []*CriticalComponent) ([]*ComponentChecksum, error) {
	return nil, nil
}

// Version - version
func (sv *testService) Version() *VersionInfo {
	return &VersionInfo{}
}

// // NewBoostData - new a BoostData
// func (sv *testService) NewBoostData() interface{} {
// 	return nil
// }

// // NewBoostDataList - new a list for BoostData
// func (sv *testService) NewBoostDataList() []interface{} {
// 	return nil
// }

// // NewPlayerBoostData - new a PlayerBoostData
// func (sv *testService) NewPlayerBoostData() interface{} {
// 	return nil
// }

// OnPlayBoostData - after call Play
func (sv *testService) OnPlayBoostData(params *PlayParams, result *PlayResult) error {
	return nil
}

// GetGameConfig - get GATIGameConfig
func (sv *testService) GetGameConfig() *GATIGameConfig {
	return nil
}

// Evaluate -
func (sv *testService) Evaluate(params *EvaluateParams, id string) (*EvaluateResult, error) {
	return nil, nil
}

func Test_Serv(t *testing.T) {
	sgc7utils.InitLogger("", "", "debug", true, "")

	json := jsoniter.ConfigCompatibleWithStandardLibrary

	cfg := &Config{
		GameID:      "1019",
		BindAddr:    "127.0.0.1:7891",
		IsDebugMode: true,
	}

	gc, err := LoadGATIGameConfig("../unittestdata/game_configuration.json")
	assert.NoError(t, err)

	service := &testService{
		&sgc7game.Config{
			Width:  5,
			Height: 3,
		},
		0,
		gc,
	}

	serv := NewServ(service, cfg)
	client := NewClient("http://127.0.0.1:7891/v2/games/1019/")

	go func() {
		err := serv.Start()
		if err != nil {
			t.Fatalf("Test_Serv Start error %v",
				err)
		}
	}()

	time.Sleep(time.Second * 3)

	sc, buff, err := sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/config", nil)
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

	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/config", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/initialize", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"playerStatePublic\":\"{}\",\"playerStatePrivate\":\"{}\"}", "they should be equal")

	service.initmode = 1

	sc, buff, err = sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/initialize", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"playerStatePublic\":\"{\\\"curgamemod\\\":\\\"BG\\\"}\",\"playerStatePrivate\":\"{}\"}", "they should be equal")

	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/initialize", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/validate", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/validate", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	validateParams := &ValidateParams{}
	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/validate", nil, validateParams)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/play", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/play", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	playParams := &PlayParams{}
	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/play", nil, playParams)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	clientcfg, err := client.GetConfig()
	assert.Nil(t, err)
	assert.NotNil(t, clientcfg)

	clientps, err := client.Initialize()
	assert.Nil(t, err)
	assert.NotNil(t, clientps)

	serv.Stop()

	t.Logf("Test_Serv OK")
}
