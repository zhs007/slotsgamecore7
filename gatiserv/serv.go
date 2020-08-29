package gatiserv

import (
	"github.com/valyala/fasthttp"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
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

	return s
}
