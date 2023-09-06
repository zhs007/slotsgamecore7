package gati

import (
	"strconv"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	"go.uber.org/zap"
)

// GetRngs - get rngs
func GetRngs(rngURL string, gameID string, nums int) ([]int, error) {
	url := goutils.AppendString(rngURL, "?size=", strconv.Itoa(nums))

	mapHeader := make(map[string]string)
	mapHeader["X-Game-ID"] = gameID
	code, body, err := sgc7http.HTTPGet(url, mapHeader)
	if err != nil {
		goutils.Error("gati.GetRngs:HTTPGet",
			zap.Error(err),
			zap.String("url", url))

		return nil, err
	}

	if code != 200 {
		goutils.Error("gati.GetRngs:HTTPGet",
			zap.Int("code", code),
			zap.String("url", url))

		return nil, err
	}

	lst := []int{}
	err = sonic.Unmarshal(body, &lst)
	if err != nil {
		return nil, err
	}

	return lst, nil
}
