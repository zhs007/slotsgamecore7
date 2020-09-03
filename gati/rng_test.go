package gati

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
)

// func Test_GetRngs(t *testing.T) {

// 	lst, err := GetRngs("http://127.0.0.1:50000/numbers", 936207324, 100)
// 	if err != nil {
// 		t.Fatalf("Test_GetRngs GetRngs error %v",
// 			err)
// 	}

// 	if len(lst) != 100 {
// 		t.Fatalf("Test_GetRngs GetRngs lst err")
// 	}

// 	t.Logf("Test_GetRngs OK")
// }

func Test_GetRngsMock(t *testing.T) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	const URL = "http://127.0.0.1:50000/numbers"
	res := []int{123, 1, 2}

	resbuff, err := json.Marshal(res)

	// httpmock.RegisterResponder("GET",
	// 	fmt.Sprintf("%s?size=%d", URL, 3),
	// 	httpmock.NewStringResponder(200, string(resbuff)))

	httpmock.RegisterResponder("GET",
		fmt.Sprintf("%s?size=%d", URL, 3),
		func(req *http.Request) (*http.Response, error) {
			assert.Equal(t, req.Header.Get("X-Game-ID"), "936207324", "they should be equal")

			return httpmock.NewStringResponder(200, string(resbuff))(req)
		})

	lst, err := GetRngs(URL, "936207324", 3)
	if err != nil {
		t.Fatalf("Test_GetRngs GetRngs error %v",
			err)
	}

	if len(lst) != 3 {
		t.Fatalf("Test_GetRngs GetRngs lst err")
	}

	for i, v := range lst {
		if v != res[i] {
			t.Fatalf("Test_GetRngs lst[%d] %d != %d",
				i, v, res[i])
		}
	}

	t.Logf("Test_GetRngsMock OK")
}
