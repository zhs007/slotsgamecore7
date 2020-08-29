package gatiserv

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// BasicURL - basic url
const BasicURL = "/v2/games/"

// Serv -
type Serv struct {
	*sgc7http.Serv
	Service IService
	Cfg     *Config
}

// NewServ - new a serv
func NewServ(service IService, cfg *Config) *Serv {
	s := &Serv{
		sgc7http.NewServ(cfg.BindAddr, cfg.IsDebugMode),
		service,
		cfg,
	}

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/config"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			ret := s.Service.Config()
			if ret == nil {
				s.SetStringResponse(ctx, "{}")
			} else {
				s.SetResponse(ctx, ret)
			}
		})

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/initialize"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			ps := s.Service.Initialize()
			str, err := s.BuildPlayerStateString(ps)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.initialize:BuildPlayerStateString",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetStringResponse(ctx, str)
		})

	return s
}

// BuildPlayerStateString - sgc7game.IPlayerState => string
func (serv *Serv) BuildPlayerStateString(ps sgc7game.IPlayerState) (string, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	if ps == nil {
		return "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}", nil
	}

	psb, err := json.Marshal(ps.GetPublic())
	if err != nil {
		sgc7utils.Warn("gatiserv.Serv.BuildPlayerStateString:Marshal GetPublic",
			zap.Error(err))

		return "", err
	}

	psp, err := json.Marshal(ps.GetPrivate())
	if err != nil {
		sgc7utils.Warn("gatiserv.Serv.BuildPlayerStateString:Marshal GetPrivate",
			zap.Error(err))

		return "", err
	}

	return sgc7utils.AppendString(
		"{\"playerStatePublic\":", string(psb), ",\"playerStatePrivate\":", string(psp), "}"), nil
}
