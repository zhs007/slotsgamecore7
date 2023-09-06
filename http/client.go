package sgc7http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
)

// HTTPGet - get
func HTTPGet(url string, header map[string]string) (int, []byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1, nil, err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, nil
}

// HTTPPost - post
func HTTPPost(url string, header map[string]string, bodyObj any) (int, []byte, error) {
	var req *http.Request
	var err error
	client := &http.Client{}

	if bodyObj != nil {
		bb, err := sonic.Marshal(bodyObj)
		if err != nil {
			goutils.Warn("sgc7http.HTTPPost:Marshal",
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

	for k, v := range header {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, body, nil
}

// HTTPPostEx - post
func HTTPPostEx(url string, header map[string]string, body []byte) (int, []byte, error) {
	var req *http.Request
	var err error
	client := &http.Client{}

	if body != nil {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			return -1, nil, err
		}
	} else {
		req, err = http.NewRequest("POST", url, nil)
		if err != nil {
			return -1, nil, err
		}
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return -1, nil, err
	}

	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, err
	}

	return resp.StatusCode, rbody, nil
}
