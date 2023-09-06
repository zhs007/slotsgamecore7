package gatiserv

import (
	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"

	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
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
	err = sonic.Unmarshal(buff, cfg)
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
func (client *Client) Initialize(ps *PlayerState) error {

	sc, buff, err := sgc7http.HTTPGet(
		goutils.AppendString(client.ServURL, "initialize"),
		nil)
	if err != nil {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.Initialize:HTTPGet",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return ErrNonStatusOK
	}

	err = ParsePlayerState(string(buff), ps)
	if err != nil {
		goutils.Error("gatiserv.Client.Initialize:ParsePlayerState",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return err
	}

	return nil
}

// PlayEx - play with string parameter
func (client *Client) PlayEx(param string) (*PlayResult, error) {
	sc, buff, err := sgc7http.HTTPPostEx(
		goutils.AppendString(client.ServURL, "play"),
		nil,
		[]byte(param),
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

	pr, err := ParsePlayResult(string(buff))
	if err != nil {
		goutils.Error("gatiserv.Client.PlayEx:ParsePlayResult",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return pr, nil
}

// Checksum - checksum
func (client *Client) Checksum(arr []*CriticalComponent) ([]*ComponentChecksum, error) {
	buf, err := sonic.Marshal(arr)
	if err != nil {
		goutils.Error("gatiserv.Client.Checksum:Marshal",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buf)),
			zap.Error(err))

		return nil, err
	}

	sc, buff, err := sgc7http.HTTPPostEx(
		goutils.AppendString(client.ServURL, "checksum"),
		nil,
		buf,
	)
	if err != nil {
		goutils.Error("gatiserv.Client.Checksum:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Error(err))

		return nil, err
	}

	if sc != fasthttp.StatusOK {
		goutils.Error("gatiserv.Client.Checksum:HTTPPostEx",
			zap.String("ServURL", client.ServURL),
			zap.Int("status", sc))

		return nil, ErrNonStatusOK
	}

	ret := []*ComponentChecksum{}
	err = sonic.Unmarshal(buff, &ret)
	if err != nil {
		goutils.Error("gatiserv.Client.Checksum:Unmarshal",
			zap.String("ServURL", client.ServURL),
			zap.String("body", string(buff)),
			zap.Error(err))

		return nil, err
	}

	return ret, nil
}
