package simserv

import (
	"github.com/valyala/fasthttp"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.uber.org/zap"
)

// BasicURL - basic url
const BasicURL = "/game"

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

	s.RegHandle(goutils.AppendString(BasicURL, "/config"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret := s.Service.Config()
			if ret == nil {
				s.SetStringResponse(ctx, "{}")
			} else {
				s.SetResponse(ctx, ret)
			}
		})

	s.RegHandle(goutils.AppendString(BasicURL, "/initialize"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ps := s.Service.Initialize()
			if ps == nil {
				s.SetStringResponse(ctx, "{\"playerStatePublic\":{},\"playerStatePrivate\":{}}")

				return
			}

			s.SetResponse(ctx, ps)
		})

	s.RegHandle(goutils.AppendString(BasicURL, "/play"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &sgc7pb.RequestPlay{}
			err := s.ParseBody(ctx, params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:ParseBody",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.Service.Play(params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:Play",
					zap.Error(err))

				if err == sgc7game.ErrInvalidStake {
					s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

					return
				}

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetResponse(ctx, ret)
		})

	return s
}
