package gatiserv

import (
	"github.com/valyala/fasthttp"
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

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/initialize"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ps := s.Service.Initialize()
			str, err := BuildPlayerStateString(ps)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.initialize:BuildPlayerStateString",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetStringResponse(ctx, str)
		})

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/validate"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &ValidateParams{}
			err := s.ParseBody(ctx, params)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.validate:ParseBody",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret := s.Service.Validate(params)
			s.SetResponse(ctx, ret)
		})

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/play"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &PlayParams{}
			err := s.ParseBody(ctx, params)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.play:ParseBody",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.Service.Play(params)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.play:Play",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetResponse(ctx, ret)
		})

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/checksum"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := []*CriticalComponent{}
			err := s.ParseBody(ctx, params)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.checksum:ParseBody",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.Service.Checksum(params)
			if err != nil {
				sgc7utils.Warn("gatiserv.Serv.checksum:Checksum",
					zap.Error(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetResponse(ctx, ret)
		})

	s.RegHandle(sgc7utils.AppendString(BasicURL, cfg.GameID, "/version"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret := s.Service.Version()

			s.SetResponse(ctx, ret)
		})

	return s
}
