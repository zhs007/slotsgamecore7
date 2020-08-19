package gati

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
)

// GetRngs - get rngs
func GetRngs(rngURL string, gameID int, nums int) ([]int, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", sgc7utils.AppendString(rngURL, "?size=", strconv.Itoa(nums)), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Game-ID", strconv.Itoa(gameID))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	lst := []int{}
	err = json.Unmarshal(body, &lst)
	if err != nil {
		return nil, err
	}

	return lst, nil
}
