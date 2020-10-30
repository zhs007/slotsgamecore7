package dtserv

import (
	"context"
	"net"

	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7utils "github.com/zhs007/slotsgamecore7/utils"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - DreamTech Service
type Serv struct {
	lis      net.Listener
	grpcServ *grpc.Server
	service  IService
}

// NewServ -
func NewServ(service IService, bindaddr string, version string) (*Serv, error) {
	lis, err := net.Listen("tcp", bindaddr)
	if err != nil {
		sgc7utils.Error("NewServ.Listen",
			zap.Error(err))

		return nil, err
	}

	grpcServ := grpc.NewServer()

	serv := &Serv{
		lis:      lis,
		grpcServ: grpcServ,
		service:  service,
	}

	sgc7pb.RegisterDTGameServiceServer(grpcServ, serv)

	sgc7utils.Info("NewServ OK.",
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

	return
}

// GetConfig - get config
func (serv *Serv) GetConfig(ctx context.Context, req *sgc7pb.RequestConfig) (*sgc7pb.GameConfig, error) {
	cfg := serv.service.GetConfig()

	return BuildGameConfig(cfg), nil
}

// Initialize - initialize a player
func (serv *Serv) Initialize(ctx context.Context, req *sgc7pb.RequestInitialize) (*sgc7pb.PlayerState, error) {
	ps := serv.service.Initialize()

	return serv.service.BuildPlayerStatePB(ps), nil
}

// Play - play game
func (serv *Serv) Play(ctx context.Context, req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	return nil, nil
}
