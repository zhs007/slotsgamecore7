package gati

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// GetRngs - get rngs
func GetRngs(rngURL string, gameID int, nums int) ([]int, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	url := sgc7utils.AppendString(rngURL, "?size=", strconv.Itoa(nums))

	mapHeader := make(map[string]string)
	mapHeader["X-Game-ID"] = strconv.Itoa(gameID)
	code, body, err := sgc7http.HTTPGet(url, mapHeader)
	if err != nil {
		sgc7utils.Error("gati.GetRngs:HTTPGet",
			zap.Error(err),
			zap.String("url", url))

		return nil, err
	}

	if code != 200 {
		sgc7utils.Error("gati.GetRngs:HTTPGet",
			zap.Int("code", code),
			zap.String("url", url))

		return nil, err
	}

	lst := []int{}
	err = json.Unmarshal(body, &lst)
	if err != nil {
		return nil, err
	}

	return lst, nil
}
