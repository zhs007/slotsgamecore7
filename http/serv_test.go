package sgc7http

import (
	"testing"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_Serv(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	serv := NewServ("127.0.0.1:7890", true)

	type requestBody struct {
		Param1 int    `json:"param1"`
		Param2 string `json:"param2"`
	}

	type response struct {
		Result string `json:"result"`
	}

	serv.RegHandle("/index", func(ctx *fasthttp.RequestCtx, serv *Serv) {
		if !ctx.Request.Header.IsGet() {
			serv.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

			return
		}

		hasa := false
		hasb := false

		ctx.QueryArgs().VisitAll(func(k []byte, v []byte) {
			if string(k) == "a" {
				hasa = true

				assert.Equal(t, string(v), "123", "they should be equal")
			} else if string(k) == "b" {
				hasb = true

				assert.Equal(t, string(v), "hello", "they should be equal")
			}
		})

		assert.Equal(t, hasa, true, "they should be equal")
		assert.Equal(t, hasb, true, "they should be equal")

		r := &response{
			Result: "OK",
		}

		serv.SetResponse(ctx, r)
	})

	serv.RegHandle("/index2", func(ctx *fasthttp.RequestCtx, serv *Serv) {
		hasa := false
		hasb := false

		ctx.QueryArgs().VisitAll(func(k []byte, v []byte) {
			if string(k) == "a" {
				hasa = true

				assert.Equal(t, string(v), "456", "they should be equal")
			} else if string(k) == "b" {
				hasb = true

				assert.Equal(t, string(v), "world", "they should be equal")
			}
		})

		assert.Equal(t, hasa, true, "they should be equal")
		assert.Equal(t, hasb, true, "they should be equal")

		serv.SetStringResponse(ctx, "{\"result\":\"123\"}")
	})

	serv.RegHandle("/index3", func(ctx *fasthttp.RequestCtx, serv *Serv) {
		hasheaderkey := false
		ctx.Request.Header.VisitAll(func(k []byte, v []byte) {
			// t.Logf("Test_Serv " + string(k) + " " + string(v))

			if string(k) == "Myheaderkey" {
				assert.Equal(t, string(v), "abc", "they should be equal")

				hasheaderkey = true
			}
		})

		assert.Equal(t, hasheaderkey, true, "they should be equal")

		hasa := false
		hasb := false

		ctx.QueryArgs().VisitAll(func(k []byte, v []byte) {
			if string(k) == "a" {
				hasa = true

				assert.Equal(t, string(v), "123", "they should be equal")
			} else if string(k) == "b" {
				hasb = true

				assert.Equal(t, string(v), "hello", "they should be equal")
			}
		})

		assert.Equal(t, hasa, true, "they should be equal")
		assert.Equal(t, hasb, true, "they should be equal")

		if ctx.Request.Header.IsPost() {
			params := &requestBody{}

			err := serv.ParseBody(ctx, params)
			assert.Nil(t, err, "there is non error")

			assert.Equal(t, params.Param1, 123, "they should be equal")
			assert.Equal(t, params.Param2, "hello", "they should be equal")
		}

		serv.SetResponse(ctx, nil)
	})

	go func() {
		err := serv.Start()
		if err != nil {
			t.Fatalf("Test_Serv Start error %v",
				err)
		}
	}()

	time.Sleep(time.Second * 3)

	sc, buff, err := HTTPGet("http://127.0.0.1:7890/index?a=123&b=hello", nil)
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

	sc, buff, err = HTTPGet("http://127.0.0.1:7890/index2?a=456&b=world", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"result\":\"123\"}", "they should be equal")

	err = json.Unmarshal(buff, rr)
	if err != nil {
		t.Fatalf("Test_Serv Unmarshal error %v",
			err)
	}

	assert.Equal(t, rr.Result, "123", "they should be equal")

	header3 := make(map[string]string)
	header3["myheaderkey"] = "abc"
	sc, buff, err = HTTPGet("http://127.0.0.1:7890/index3?a=123&b=hello", header3)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "", "they should be equal")

	sc, buff, err = HTTPGet("http://127.0.0.1:7890/abc?a=123&b=hello", nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 404, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = HTTPPost("http://127.0.0.1:7890/index?a=123&b=hello", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 400, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")

	sc, buff, err = HTTPPost("http://127.0.0.1:7890/index2?a=456&b=world", nil, nil)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "{\"result\":\"123\"}", "they should be equal")

	err = json.Unmarshal(buff, rr)
	if err != nil {
		t.Fatalf("Test_Serv Unmarshal error %v",
			err)
	}

	assert.Equal(t, rr.Result, "123", "they should be equal")

	post3 := &requestBody{
		Param1: 123,
		Param2: "hello",
	}
	sc, buff, err = HTTPPost("http://127.0.0.1:7890/index3?a=123&b=hello", header3, post3)
	if err != nil {
		t.Fatalf("Test_Serv httpGet error %v",
			err)
	}

	assert.Equal(t, sc, 200, "they should be equal")
	assert.NotNil(t, buff, "there is a valid buffer")
	assert.Equal(t, string(buff), "", "they should be equal")

	serv.Stop()

	t.Logf("Test_Serv OK")
}
