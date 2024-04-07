package gamecollection

import (
	"context"
	"log/slog"
	"net"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/lowcode"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
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
func NewServ(bindaddr string, version string, useOpenTelemetry bool, funcNewRNG lowcode.FuncNewRNG) (*Serv, error) {
	// lowcode.SetJsonMode()

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
		mgrGame:  NewGameMgr(funcNewRNG),
	}

	sgc7pb.RegisterGameLogicCollectionServer(grpcServ, serv)

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
func (serv *Serv) InitGame(ctx context.Context, req *sgc7pb.RequestInitGame) (*sgc7pb.ReplyInitGame, error) {
	goutils.Debug("Serv.InitGame",
		slog.Any("req", req))

	err := serv.mgrGame.InitGame(req.GameCode, []byte(req.Config))
	if err != nil {
		goutils.Error("Serv.InitGame:InitGame",
			goutils.Err(err))

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
		slog.Any("req", req))

	cfg, err := serv.mgrGame.GetGameConfig(req.GameCode)
	if err != nil {
		goutils.Error("Serv.GetGameConfig:GetGameConfig",
			goutils.Err(err))

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
		slog.Any("req", req))

	ps, err := serv.mgrGame.InitializeGamePlayer(req.GameCode)
	if err != nil {
		goutils.Error("Serv.InitializeGamePlayer:InitializeGamePlayer",
			goutils.Err(err))

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
		slog.Any("req", req))

	res, err := serv.mgrGame.PlayGame(req.GameCode, req.Play)
	if err != nil {
		goutils.Error("Serv.PlayGame:PlayGame",
			goutils.Err(err))

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
		slog.Any("req", req))

	res, err := serv.mgrGame.PlayGame(req.GameCode, req.Play)
	if err != nil {
		goutils.Error("Serv.PlayGame:PlayGame2",
			goutils.Err(err))

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
