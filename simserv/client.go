package simserv

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
)

// Client - client
type Client struct {
	ServURL string
}

// NewClient - new Client, servurl is like http://127.0.0.1:7891/game/
func NewClient(servurl string) *Client {
	return &Client{
		ServURL: servurl,
	}
}

// GetConfig - get configuration
func (client *Client) GetConfig() (*sgc7game.Config, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	sc, buff, err := sgc7http.HTTPGet(
		goutils.AppendString(client.ServURL, "config"),
		nil)
	if err != nil {
		goutils.Error("gatiserv.Client.GetConfig:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.GetConfig:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	cfg := &sgc7game.Config{}
	err = json.Unmarshal(buff, cfg)
	if err != nil {
		goutils.Error("gatiserv.Client.GetConfig:JSON",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return cfg, nil
}

// Initialize - initialize a player
func (client *Client) Initialize() (*sgc7pb.PlayerState, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	sc, buff, err := sgc7http.HTTPGet(
		goutils.AppendString(client.ServURL, "initialize"),
		nil)
	if err != nil {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	ps := &sgc7pb.PlayerState{}
	err = json.Unmarshal(buff, ps)
	if err != nil {
		goutils.Error("gatiserv.Client.Initialize:Unmarshal",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return ps, nil
}

// PlayEx - play with string parameter
func (client *Client) PlayEx(param *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	buff, err := json.Marshal(param)
	if err != nil {
		goutils.Error("gatiserv.Client.PlayEx:Marshal",
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	sc, buff, err := sgc7http.HTTPPostEx(
		goutils.AppendString(client.ServURL, "play"),
		nil,
		buff,
	)
	if err != nil {
		goutils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	rp := &sgc7pb.ReplyPlay{}
	err = json.Unmarshal(buff, rp)
	if err != nil {
		goutils.Error("gatiserv.Client.PlayEx:Unmarshal",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return rp, nil
}

// Initialize2 - initialize a player
func (client *Client) Initialize2() (string, error) {
	sc, buff, err := sgc7http.HTTPGet(
		goutils.AppendString(client.ServURL, "initialize"),
		nil)
	if err != nil {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return "", err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return "", ErrNonStatusOK
	}

	return string(buff), nil
}

// PlayEx2 - play with string parameter
func (client *Client) PlayEx2(param string) (string, error) {
	sc, buff, err := sgc7http.HTTPPostEx(
		goutils.AppendString(client.ServURL, "play"),
		nil,
		[]byte(param),
	)
	if err != nil {
		goutils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return "", err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.PlayEx:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return "", ErrNonStatusOK
	}

	return string(buff), nil
}
