package gatiserv

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	"go.uber.org/zap"
)

// APIHandle - handle
type APIHandle func(ctx *fasthttp.RequestCtx, serv *Serv)

// Serv -
type Serv struct {
	bindAddr string
	mapAPI   map[string]APIHandle
}

// NewServ - new a serv
func NewServ(bindAddr string) *Serv {
	s := &Serv{
		bindAddr: bindAddr,
		mapAPI:   make(map[string]APIHandle),
	}

	return s
}

// RegHandle - register a handle
func (s *Serv) RegHandle(name string, handle APIHandle) {
	s.mapAPI[name] = handle
}

// HandleFastHTTP -
func (s *Serv) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	h, isok := s.mapAPI[string(ctx.Path())]
	if isok && h != nil {
		h(ctx, s)
	} else {
		s.SetHTTPStatus(ctx, fasthttp.StatusNotFound)
	}
}

// Start - start a server
func (s *Serv) Start() {
	fasthttp.ListenAndServe(s.bindAddr, s.HandleFastHTTP)
}

// SetResponse - set a response
func (s *Serv) SetResponse(ctx *fasthttp.RequestCtx, jsonObj interface{}) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(jsonObj)
	if err != nil {
		sgc7utils.Warn("gatiserv.Serv.SetResponse",
			zap.Error(err))

		s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetContentType("application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(b)
}

// SetHTTPStatus - set a response with status
func (s *Serv) SetHTTPStatus(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetStatusCode(statusCode)
}
