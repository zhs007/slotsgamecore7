package sgc7game

import (
	jsoniter "github.com/json-iterator/go"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// PlayResult - result for play
type PlayResult struct {
	CurGameMod       string       `json:"curgamemod"`
	CurGameModParams interface{}  `json:"curgamemodparams"`
	NextGameMod      string       `json:"nextgamemod"`
	Scenes           []*GameScene `json:"scenes"`
	Results          []*Result    `json:"results"`
	NextCmds         []string     `json:"-"`
	NextCmdParams    []string     `json:"-"`
	CoinWin          int          `json:"-"`
	CashWin          int64        `json:"-"`
	IsFinish         bool         `json:"-"`
	IsWait           bool         `json:"-"`
}

// PlayResult2JSON - PlayResult => json
func PlayResult2JSON(pr *PlayResult) ([]byte, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	b, err := json.Marshal(pr)
	if err != nil {
		sgc7utils.Warn("sgc7game.PlayResult2JSON",
			zap.Error(err))

		return nil, err
	}

	return b, nil
}

// JSON2PlayResult - json => PlayResult
func JSON2PlayResult(buf []byte, pr *PlayResult) (*PlayResult, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	// pr := &PlayResult{}
	err := json.Unmarshal(buf, &pr)
	if err != nil {
		sgc7utils.Warn("sgc7game.JSON2PlayResult",
			zap.Error(err))
		return nil, err
	}

	return pr, nil
}
