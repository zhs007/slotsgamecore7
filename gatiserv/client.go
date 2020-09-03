package gatiserv

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"

	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// Client - client
type Client struct {
	ServURL string
}

// NewClient - new Client, servurl is like http://127.0.0.1:7891/v2/games/1019/
func NewClient(servurl string) *Client {
	return &Client{
		ServURL: servurl,
	}
}

// GetConfig - get configuration
func (client *Client) GetConfig() (*sgc7game.Config, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	sc, buff, err := sgc7http.HTTPGet(
		sgc7utils.AppendString(client.ServURL, "config"),
		nil)
	if err != nil {
		sgc7utils.Error("gatiserv.Client.GetConfig:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		sgc7utils.Error("gatiserv.Client.GetConfig:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	cfg := &sgc7game.Config{}
	err = json.Unmarshal(buff, cfg)
	if err != nil {
		sgc7utils.Error("gatiserv.Client.GetConfig:JSON",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return cfg, nil
}

// Initialize - initialize a player
func (client *Client) Initialize() (*PlayerState, error) {

	sc, buff, err := sgc7http.HTTPGet(
		sgc7utils.AppendString(client.ServURL, "initialize"),
		nil)
	if err != nil {
		sgc7utils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		sgc7utils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	ps, err := ParsePlayerState(string(buff))
	if err != nil {
		sgc7utils.Error("gatiserv.Client.Initialize:ParsePlayerState",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return ps, nil
}

// PlayEx - play with string parameter
func (client *Client) PlayEx(param string) (*PlayResult, error) {
	sc, buff, err := sgc7http.HTTPPostEx(
		sgc7utils.AppendString(client.ServURL, "play"),
		nil,
		[]byte(param),
	)
	if err != nil {
		sgc7utils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		sgc7utils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	pr, err := ParsePlayResult(string(buff))
	if err != nil {
		sgc7utils.Error("gatiserv.Client.PlayEx:ParsePlayResult",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return pr, nil
}
