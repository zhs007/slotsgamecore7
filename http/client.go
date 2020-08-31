package sgc7http

import (
	"bytes"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// HTTPGet - get
func HTTPGet(url string, header map[string]string) (int, []byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, nil, err
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
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

// HTTPPost - post
func HTTPPost(url string, header map[string]string, bodyObj interface{}) (int, []byte, error) {
	var req *http.Request
	var err error
	client := &http.Client{}

	if bodyObj != nil {
		json := jsoniter.ConfigCompatibleWithStandardLibrary

		bb, err := json.Marshal(bodyObj)
		if err != nil {
			sgc7utils.Warn("sgc7http.HTTPPost:Marshal",
				zap.Error(err))

			return -1, nil, err
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(bb))
		if err != nil {
			return -1, nil, err
		}
	} else {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return -1, nil, err
		}
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}

	req.Header.Set("Content-Type", "application/json")

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
