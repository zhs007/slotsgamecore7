package gatiserv

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
)

//--------------------------------------------------------------------------------------
// testPlayerState

type testPlayerState struct {
}

// SetPublic - set player public state
func (ps *testPlayerState) SetPublic(pub any) error {
	return nil
}

// SetPrivate - set player private state
func (ps *testPlayerState) SetPrivate(pri any) error {
	return nil
}

// // SetPublicString - set player public state
// func (ps *testPlayerState) SetPublicString(pub string) error {
// 	return nil
// }

// // SetPrivateString - set player private state
// func (ps *testPlayerState) SetPrivateString(pri string) error {
// 	return nil
// }

// GetPublic - get player public state
func (ps *testPlayerState) GetPublic() any {
	return sgc7game.BasicPlayerPublicState{
		CurGameMod: "BG",
	}
}

// GetPrivate - get player private state
func (ps *testPlayerState) GetPrivate() any {
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
func (sv *testService) Initialize() *PlayerState {
	if sv.initmode == 0 {
		return nil
	}

	ips := &testPlayerState{}

	return &PlayerState{
		Public:  ips.GetPublic(),
		Private: ips.GetPrivate(),
	}
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
	return []*ComponentChecksum{{ID: 1, Checksum: "1234567"}}, nil
}

// Version - version
func (sv *testService) Version() *VersionInfo {
	return &VersionInfo{}
}

// // NewBoostData - new a BoostData
// func (sv *testService) NewBoostData() any {
// 	return nil
// }

// // NewBoostDataList - new a list for BoostData
// func (sv *testService) NewBoostDataList() []any {
// 	return nil
// }

// // NewPlayerBoostData - new a PlayerBoostData
// func (sv *testService) NewPlayerBoostData() any {
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
	goutils.InitLogger("", "", "debug", true, "")

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
			goutils.Error("Test_Serv Start error",
				zap.Error(err))
			// t.Fatalf("Test_Serv Start error %v",
			// 	err)
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
	assert.Equal(t, string(buff), "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}", "they should be equal")

	service.initmode = 1

	sc, buff, err = sgc7http.HTTPGet("http://127.0.0.1:7891/v2/games/1019/initialize", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"playerStatePublic\":{\"curgamemod\":\"BG\",\"nextm\":0},\"playerStatePrivate\":{}}", "they should be equal")

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

	assert.Equal(t, sc, 405, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/validate", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 405, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	validateParams := &ValidateParams{}
	sc, buff, err = sgc7http.HTTPPost("http://127.0.0.1:7891/v2/games/1019/validate", nil, validateParams)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 405, "they should be equal")
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

	playParams = &PlayParams{
		PlayerState: &PlayerState{
			Public: &sgc7game.BasicPlayerPublicState{
				CurGameMod: "BG",
			},
			Private: &sgc7game.BasicPlayerPrivateState{},
		},
	}
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

	clientps := &PlayerState{
		Public:  &sgc7game.BasicPlayerPublicState{},
		Private: &sgc7game.BasicPlayerPrivateState{},
	}
	err = client.Initialize(clientps)
	assert.Nil(t, err)
	assert.NotNil(t, clientps)

	retChecksum, err := client.Checksum([]*CriticalComponent{
		{ID: 1, Name: "test", Location: "test/test"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, retChecksum)
	assert.Equal(t, len(retChecksum), 1)
	assert.Equal(t, retChecksum[0].ID, 1)
	assert.Equal(t, retChecksum[0].Checksum, "1234567")

	serv.Stop()

	t.Logf("Test_Serv OK")
}
