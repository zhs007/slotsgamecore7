package gatiserv

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

type testService struct {
	cfg *Config
}

// Config - get configuration
func (sv *testService) Config() interface{} {
	return sv.cfg
}

func Test_Serv(t *testing.T) {
	cfg := &Config{
		GameID:      "1019",
		BindAddr:    "127.0.0.1:7891",
		IsDebugMode: true,
	}
	serv := NewServ(&testService{cfg}, cfg)

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

	rr := &Config{}
	err = json.Unmarshal(buff, rr)
	if err != nil {
		t.Fatalf("Test_Serv Unmarshal error %v",
			err)
	}

	assert.Equal(t, rr.GameID, cfg.GameID, "they should be equal")
	assert.Equal(t, rr.BindAddr, cfg.BindAddr, "they should be equal")
	assert.Equal(t, rr.IsDebugMode, cfg.IsDebugMode, "they should be equal")

	serv.Stop()

	t.Logf("Test_Serv OK")
}
