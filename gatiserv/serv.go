package gatiserv

import (
	"log/slog"

	"github.com/valyala/fasthttp"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
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

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/config"),
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

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/initialize"),
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
			// str, err := BuildPlayerStateString(ps)
			// if err != nil {
			// 	goutils.Warn("gatiserv.Serv.initialize:BuildPlayerStateString",
			// 		goutils.Err(err))

			// 	s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

			// 	return
			// }

			s.SetResponse(ctx, ps)
		})

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/validate"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			s.SetHTTPStatus(ctx, fasthttp.StatusMethodNotAllowed)

			// return

			// if !ctx.Request.Header.IsPost() {
			// 	s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

			// 	return
			// }

			// params := &ValidateParams{}
			// err := s.ParseBody(ctx, params)
			// if err != nil {
			// 	goutils.Warn("gatiserv.Serv.validate:ParseBody",
			// 		goutils.Err(err))

			// 	s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

			// 	return
			// }

			// ret := s.Service.Validate(params)
			// s.SetResponse(ctx, ret)
		})

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/play"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &PlayParams{
				PlayerState: s.Service.Initialize(),
			}
			err := s.ParseBody(ctx, params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:ParseBody",
					goutils.Err(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.Service.Play(params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:Play",
					goutils.Err(err))

				if err == sgc7game.ErrInvalidStake {
					s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

					return
				}

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			err = s.Service.OnPlayBoostData(params, ret)
			if err != nil {
				goutils.Warn("gatiserv.Serv.play:OnPlayBoostData",
					goutils.Err(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetResponse(ctx, ret)
		})

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/checksum"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := []*CriticalComponent{}
			err := s.ParseBody(ctx, &params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.checksum:ParseBody",
					goutils.Err(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := s.Service.Checksum(params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.checksum:Checksum",
					goutils.Err(err))

				s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			s.SetResponse(ctx, ret)
		})

	s.RegHandle(goutils.AppendString(BasicURL, cfg.GameID, "/version"),
		func(ctx *fasthttp.RequestCtx, serv *sgc7http.Serv) {
			if !ctx.Request.Header.IsGet() {
				s.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret := s.Service.Version()

			s.SetResponse(ctx, ret)
		})

	gc := s.Service.GetGameConfig()
	if gc != nil {
		for _, v := range gc.GameObjectives {
			s.RegMission(v.ObjectiveID)
		}
	}

	return s
}

// RegMission -
func (serv *Serv) RegMission(id string) {
	goutils.Info("gatiserv.Serv.RegHandle",
		slog.String("id", id))

	serv.RegHandle(goutils.AppendString(BasicURL, serv.Cfg.GameID, "/evaluate/"+id),
		func(ctx *fasthttp.RequestCtx, serv1 *sgc7http.Serv) {
			if !ctx.Request.Header.IsPost() {
				serv.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			params := &EvaluateParams{}
			err := serv.ParseBody(ctx, params)
			if err != nil {
				goutils.Warn("gatiserv.Serv.evaluate:ParseBody",
					slog.String("id", id),
					goutils.Err(err))

				serv.SetHTTPStatus(ctx, fasthttp.StatusBadRequest)

				return
			}

			ret, err := serv.Service.Evaluate(params, id)
			if err != nil {
				goutils.Warn("gatiserv.Serv.evaluate:Evaluate",
					slog.String("id", id),
					goutils.Err(err))

				serv.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

				return
			}

			serv.SetResponse(ctx, ret)
		})
}
