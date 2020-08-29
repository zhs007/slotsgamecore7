package gatiserv

import (
	"github.com/valyala/fasthttp"
	sgc7http "github.com/zhs007/slotsgamecore7/http"
)

// APIHandle - handle
type APIHandle func(ctx *fasthttp.RequestCtx, serv *Serv)

// Serv -
type Serv struct {
	sgc7http.Serv
}

// NewServ - new a serv
func NewServ(bindAddr string, isDebugMode bool) *Serv {
	s := &Serv{
		sgc7http.NewServ(bindAddr, isDebugMode),
	}

	return s
}
