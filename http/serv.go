package sgc7http

import (
	"net"

	"github.com/bytedance/sonic"
	"github.com/valyala/fasthttp"
	goutils "github.com/zhs007/goutils"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// APIHandle - handle
type APIHandle func(ctx *fasthttp.RequestCtx, serv *Serv)

// Serv -
type Serv struct {
	bindAddr    string
	mapAPI      map[string]APIHandle
	isDebugMode bool
	listener    net.Listener
}

// NewServ - new a serv
func NewServ(bindAddr string, isDebugMode bool) *Serv {
	s := &Serv{
		bindAddr:    bindAddr,
		mapAPI:      make(map[string]APIHandle),
		isDebugMode: isDebugMode,
	}

	return s
}

// RegHandle - register a handle
func (s *Serv) RegHandle(name string, handle APIHandle) {
	s.mapAPI[name] = handle
}

// HandleFastHTTP -
func (s *Serv) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	if s.isDebugMode {
		s.outputDebugInfo(ctx)
	}

	h, isok := s.mapAPI[string(ctx.Path())]
	if isok && h != nil {
		h(ctx, s)
	} else {
		s.SetHTTPStatus(ctx, fasthttp.StatusNotFound)
	}
}

// Stop - stop a server
func (s *Serv) Stop() error {
	if s.listener != nil {
		s.listener.Close()

		s.listener = nil
	}

	return nil
}

// Start - start a server
func (s *Serv) Start() error {
	if s.listener != nil {
		s.Stop()
	}

	ln, err := net.Listen("tcp4", s.bindAddr)
	if err != nil {
		goutils.Error("gatiserv.Serv.Start:Listen",
			zap.Error(err))

		return err
	}

	s.listener = ln

	return fasthttp.Serve(ln, s.HandleFastHTTP)
}

// SetResponse - set a response
func (s *Serv) SetResponse(ctx *fasthttp.RequestCtx, jsonObj any) {
	if jsonObj == nil {
		ctx.SetContentType("application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte(""))

		return
	}

	b, err := sonic.Marshal(jsonObj)
	if err != nil {
		goutils.Warn("gatiserv.Serv.SetResponse",
			zap.Error(err))

		s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetContentType("application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(b)

	goutils.Debug("gatiserv.Serv.SetResponse",
		zap.String("RequestURI", string(ctx.RequestURI())),
		zap.String("body", string(b)))
}

// SetResponse - set a response
func (s *Serv) SetPBResponse(ctx *fasthttp.RequestCtx, msg proto.Message) {
	if msg == nil {
		ctx.SetContentType("application/json;charset=UTF-8")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody([]byte(""))

		return
	}

	// m := protojson.MarshalOptions{
	// 	Resolver:
	//   }
	result, err := protojson.Marshal(msg)
	if err != nil {
		goutils.Warn("gatiserv.Serv.SetResponse",
			zap.Error(err))

		s.SetHTTPStatus(ctx, fasthttp.StatusInternalServerError)

		return
	}

	ctx.SetContentType("application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(result))

	goutils.Debug("gatiserv.Serv.SetResponse",
		zap.String("RequestURI", string(ctx.RequestURI())),
		zap.String("body", string(result)))
}

// SetStringResponse - set a response with string
func (s *Serv) SetStringResponse(ctx *fasthttp.RequestCtx, str string) {
	ctx.SetContentType("application/json;charset=UTF-8")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte(str))

	goutils.Debug("gatiserv.Serv.SetStringResponse",
		zap.String("RequestURI", string(ctx.RequestURI())),
		zap.String("body", str))
}

// SetHTTPStatus - set a response with status
func (s *Serv) SetHTTPStatus(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetStatusCode(statusCode)

	goutils.Debug("gatiserv.Serv.SetHTTPStatus",
		zap.String("RequestURI", string(ctx.RequestURI())),
		zap.Int("statusCode", statusCode))
}

func (s *Serv) outputDebugInfo(ctx *fasthttp.RequestCtx) {
	goutils.Debug("Request infomation",
		zap.String("Method", string(ctx.Method())),
		zap.String("RequestURI", string(ctx.RequestURI())),
		zap.String("Path", string(ctx.Path())),
		zap.String("Host", string(ctx.Host())),
		zap.String("UserAgent", string(ctx.UserAgent())),
		zap.String("RemoteIP", ctx.RemoteIP().String()),
		zap.Uint64("ConnRequestNum", ctx.ConnRequestNum()),
		zap.Time("ConnTime", ctx.ConnTime()),
		zap.Time("Time", ctx.Time()),
	)

	if ctx.QueryArgs() != nil {
		goutils.Debug("Request infomation QueryArgs",
			zap.String("QueryArgs", ctx.QueryArgs().String()),
		)
	}

	if ctx.PostArgs() != nil {
		goutils.Debug("Request infomation PostArgs",
			zap.String("PostArgs", ctx.PostArgs().String()),
		)
	}

	if ctx.PostBody() != nil {
		goutils.Debug("Request infomation PostBody",
			zap.String("PostBody", string(ctx.PostBody())),
		)
	}
}

// ParseBody - parse body
func (s *Serv) ParseBody(ctx *fasthttp.RequestCtx, params any) error {
	err := sonic.Unmarshal(ctx.PostBody(), params)
	if err != nil {
		return err
	}

	return nil
}
