package gamecollection

import (
	"context"
	"net"

	goutils "github.com/zhs007/goutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - DreamTech Service
type Serv struct {
	sgc7pb.UnimplementedGameLogicCollectionServer
	lis      net.Listener
	grpcServ *grpc.Server
}

// NewServ -
func NewServ(bindaddr string, version string, useOpenTelemetry bool) (*Serv, error) {
	lis, err := net.Listen("tcp", bindaddr)
	if err != nil {
		goutils.Error("NewServ.Listen",
			zap.Error(err))

		return nil, err
	}

	var grpcServ *grpc.Server

	if useOpenTelemetry {
		grpcServ = grpc.NewServer(
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

	sgc7pb.RegisterGameLogicCollectionServer(grpcServ, serv)

	goutils.Info("NewServ OK.",
		zap.String("addr", bindaddr),
		zap.String("ver", version),
		zap.String("corever", sgc7ver.Version))

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
func (serv *Serv) InitGame(ctx context.Context, req *sgc7pb.RequestInitGame) (*sgc7pb.ReplyInitGame, error) {
	return nil, nil
}
