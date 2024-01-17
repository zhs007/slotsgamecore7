package gamecollection

import (
	"context"
	"net"

	goutils "github.com/zhs007/goutils"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Serv - Service
type Serv struct {
	sgc7pb.UnimplementedGameLogicCollectionServer
	lis      net.Listener
	grpcServ *grpc.Server
	mgrGame  *GameMgr
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
		mgrGame:  NewGameMgr(),
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
	goutils.Debug("Serv.InitGame",
		goutils.JSON("req", req))

	err := serv.mgrGame.InitGame(req.GameCode, []byte(req.Config))
	if err != nil {
		goutils.Error("Serv.InitGame:InitGame",
			zap.Error(err))

		return &sgc7pb.ReplyInitGame{
			IsOK: false,
			Err:  err.Error(),
		}, nil
	}

	return &sgc7pb.ReplyInitGame{
		IsOK: true,
	}, nil
}

// GetGameConfig - get game config
func (serv *Serv) GetGameConfig(ctx context.Context, req *sgc7pb.RequestGameConfig) (*sgc7pb.ReplyGameConfig, error) {
	goutils.Debug("Serv.GetGameConfig",
		goutils.JSON("req", req))

	cfg, err := serv.mgrGame.GetGameConfig(req.GameCode)
	if err != nil {
		goutils.Error("Serv.GetGameConfig:GetGameConfig",
			zap.Error(err))

		return &sgc7pb.ReplyGameConfig{
			IsOK: false,
			Err:  err.Error(),
		}, nil
	}

	return &sgc7pb.ReplyGameConfig{
		IsOK:       true,
		GameConfig: sgc7pbutils.BuildPBGameConfig(cfg),
	}, nil
}

// InitializeGamePlayer - initialize a player
func (serv *Serv) InitializeGamePlayer(ctx context.Context, req *sgc7pb.RequestInitializeGamePlayer) (*sgc7pb.ReplyInitializeGamePlayer, error) {
	goutils.Debug("Serv.InitializeGamePlayer",
		goutils.JSON("req", req))

	ps, err := serv.mgrGame.InitializeGamePlayer(req.GameCode)
	if err != nil {
		goutils.Error("Serv.InitializeGamePlayer:InitializeGamePlayer",
			zap.Error(err))

		return &sgc7pb.ReplyInitializeGamePlayer{
			IsOK: false,
			Err:  err.Error(),
		}, nil
	}

	return &sgc7pb.ReplyInitializeGamePlayer{
		IsOK:        true,
		PlayerState: ps,
	}, nil
}

// PlayGame - play game
func (serv *Serv) PlayGame(req *sgc7pb.RequestPlayGame, stream sgc7pb.GameLogicCollection_PlayGameServer) error {
	goutils.Debug("Serv.PlayGame",
		goutils.JSON("req", req))

	res, err := serv.mgrGame.PlayGame(req.GameCode, req.Play)
	if err != nil {
		goutils.Error("Serv.PlayGame:PlayGame",
			zap.Error(err))

		stream.Send(&sgc7pb.ReplyPlayGame{
			IsOK: false,
			Err:  err.Error(),
		})

		return nil
	}

	return stream.Send(&sgc7pb.ReplyPlayGame{
		IsOK: true,
		Play: res,
	})
}

// PlayGame2 - play game
func (serv *Serv) PlayGame2(ctx context.Context, req *sgc7pb.RequestPlayGame) (*sgc7pb.ReplyPlayGame, error) {
	goutils.Debug("Serv.PlayGame2",
		goutils.JSON("req", req))

	res, err := serv.mgrGame.PlayGame(req.GameCode, req.Play)
	if err != nil {
		goutils.Error("Serv.PlayGame:PlayGame2",
			zap.Error(err))

		return &sgc7pb.ReplyPlayGame{
			IsOK: false,
			Err:  err.Error(),
		}, nil
	}

	return &sgc7pb.ReplyPlayGame{
		IsOK: true,
		Play: res,
	}, nil
}
