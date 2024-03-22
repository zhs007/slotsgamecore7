package mathtoolset2

import (
	"context"
	"log/slog"
	"net"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// Serv - Service
type Serv struct {
	sgc7pb.UnimplementedMathToolsetServer
	lis      net.Listener
	grpcServ *grpc.Server
}

// NewServ -
func NewServ(bindaddr string, version string, useOpenTelemetry bool) (*Serv, error) {
	lowcode.SetJsonMode()

	lis, err := net.Listen("tcp", bindaddr)
	if err != nil {
		goutils.Error("NewServ.Listen",
			goutils.Err(err))

		return nil, err
	}

	var grpcServ *grpc.Server

	if useOpenTelemetry {
		grpcServ = grpc.NewServer(
			grpc.MaxRecvMsgSize(1024*1024*10),
			grpc.MaxSendMsgSize(1024*1024*10),
			grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
			grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
		)
	} else {
		grpcServ = grpc.NewServer()
	}

	serv := &Serv{
		lis:      lis,
		grpcServ: grpcServ,
	}

	sgc7pb.RegisterMathToolsetServer(grpcServ, serv)

	goutils.Info("NewServ OK.",
		slog.String("addr", bindaddr),
		slog.String("ver", version),
		slog.String("corever", sgc7ver.Version))

	return serv, nil
}

// Start - start a service
func (serv *Serv) Start(ctx context.Context) error {
	return serv.grpcServ.Serve(serv.lis)
}

// Stop - stop service
func (serv *Serv) Stop() {
	serv.lis.Close()
}

// initGame - initial game
func (serv *Serv) RunScript(ctx context.Context, req *sgc7pb.RunScript) (*sgc7pb.ReplyRunScript, error) {
	goutils.Debug("Serv.RunScript",
		slog.Any("req", req))

	sc, err := NewScriptCore(req.MapFiles)
	if err != nil {
		goutils.Error("Serv.RunScript:NewScriptCore",
			goutils.Err(err))

		return nil, err
	}

	err = sc.Run(req.Script)
	if err != nil {
		goutils.Error("Serv.RunScript:Run",
			goutils.Err(err))

		return nil, err
	}

	reply := &sgc7pb.ReplyRunScript{}

	for _, v := range sc.ErrInRun {
		reply.ScriptErrs = append(reply.ScriptErrs, v.Error())
	}

	str, err := sc.MapOutputFiles.ToJson()
	if err != nil {
		goutils.Error("Serv.RunScript:ToJson",
			goutils.Err(err))

		return nil, err
	}

	reply.MapFiles = str

	return reply, nil
}
