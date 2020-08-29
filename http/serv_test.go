package sgc7http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
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

func Test_Serv(t *testing.T) {
	serv := NewServ("127.0.0.1:7890", true)

	type response struct {
		Result string `json:"result"`
	}

	serv.RegHandle("/index", func(ctx *fasthttp.RequestCtx, serv *Serv) {
		serv.SetHTTPStatus(ctx, 400)

		r := &response{
			Result: "OK",
		}

		serv.SetResponse(ctx, r)
	})

	go func() {
		err := serv.Start()
		if err != nil {
			t.Fatalf("Test_Serv Start error %v",
				err)
		}
	}()

	time.Sleep(time.Second * 3)

	sc, buff, err := httpGet("http://127.0.0.1:7890/index?a=123&b=hello")
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	rr := &response{}
	err = json.Unmarshal(buff, rr)
	if err != nil {
		t.Fatalf("Test_Serv Unmarshal error %v",
			err)
	}

	assert.Equal(t, rr.Result, "OK", "they should be equal")

	sc, buff, err = httpGet("http://127.0.0.1:7890/abc?a=123&b=hello")
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 404, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	serv.Stop()

	t.Logf("Test_Serv OK")
}
