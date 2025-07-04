package grpcserv

import (
	"context"
	"log/slog"
	"net"

	"github.com/bytedance/sonic"
	goutils "github.com/zhs007/goutils"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	sgc7ver "github.com/zhs007/slotsgamecore7/ver"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// Serv - Game Logic Service
type Serv struct {
	sgc7pb.UnimplementedGameLogicServer
	lis      net.Listener
	grpcServ *grpc.Server
	service  IService
	game     sgc7game.IGame
}

// NewServ -
func NewServ(service IService, game sgc7game.IGame, bindaddr string, version string, useOpenTelemetry bool) (*Serv, error) {
	lis, err := net.Listen("tcp", bindaddr)
	if err != nil {
		goutils.Error("NewServ.Listen",
			goutils.Err(err))

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
		service:  service,
		game:     game,
	}

	sgc7pb.RegisterGameLogicServer(grpcServ, serv)

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

// GetConfig - get config
func (serv *Serv) GetConfig(ctx context.Context, req *sgc7pb.RequestConfig) (*sgc7pb.GameConfig, error) {
	goutils.Debug("Serv.GetConfig",
		slog.Any("req", req))

	cfg := serv.game.GetConfig()

	res := sgc7pbutils.BuildPBGameConfig(cfg)

	return res, nil
}

// Initialize - initialize a player
func (serv *Serv) Initialize(ctx context.Context, req *sgc7pb.RequestInitialize) (*sgc7pb.PlayerState, error) {
	goutils.Debug("Serv.Initialize",
		slog.Any("req", req))

	ps := serv.game.Initialize()
	res, err := serv.service.BuildPBPlayerState(ps)
	if err != nil {
		goutils.Error("Serv.Initialize:BuildPBPlayerState",
			goutils.Err(err))

		return nil, err
	}

	goutils.Debug("Serv.Initialize",
		slog.Any("reply", res))

	return res, nil
}

// Play - play game
func (serv *Serv) Play(req *sgc7pb.RequestPlay, stream sgc7pb.GameLogic_PlayServer) error {
	goutils.Debug("Serv.Play",
		slog.Any("req", req))

	res, err := serv.onPlay(req)
	if err != nil {
		goutils.Error("Serv.Play:onPlay",
			goutils.Err(err))

		return err
	}

	serv.LogReplyPlay("Serv.Play", res, zapcore.DebugLevel)

	return stream.Send(res)
}

// Play2 - play game
func (serv *Serv) Play2(ctx context.Context, req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	goutils.Debug("Serv.Play",
		slog.Any("req", req))

	res, err := serv.onPlay(req)
	if err != nil {
		goutils.Error("Serv.Play:onPlay",
			goutils.Err(err))

		return nil, err
	}

	// goutils.Debug("Serv.Play",
	// 	slog.Any("reply", res))
	serv.LogReplyPlay("Serv.Play", res, zapcore.DebugLevel)

	return res, nil
}

// ProcCheat - process cheat
func (serv *Serv) ProcCheat(plugin sgc7plugin.IPlugin, cheat string) error {
	if cheat != "" {
		str := goutils.AppendString("[", cheat, "]")

		rngs := []int{}
		err := sonic.Unmarshal([]byte(str), &rngs)
		if err != nil {
			return err
		}

		plugin.SetCache(rngs)
	}

	return nil
}

// Play - play game
func (serv *Serv) onPlay(req *sgc7pb.RequestPlay) (*sgc7pb.ReplyPlay, error) {
	ips := serv.game.NewPlayerState()
	if req.PlayerState != nil {
		err := serv.service.BuildPlayerStateFromPB(ips, req.PlayerState)
		if err != nil {
			goutils.Error("Serv.onPlay:BuildPlayerStateFromPB",
				goutils.Err(err))

			return nil, err
		}
	}

	plugin := serv.game.NewPlugin()
	defer serv.game.FreePlugin(plugin)

	serv.ProcCheat(plugin, req.Cheat)

	stake := sgc7pbutils.BuildStake(req.Stake)
	err := serv.game.CheckStake(stake)
	if err != nil {
		goutils.Error("Serv.onPlay:CheckStake",
			slog.Any("stake", stake),
			goutils.Err(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}
	gameData := serv.game.NewGameData(stake)
	if gameData == nil {
		goutils.Error("Serv.onPlay:NewGameData",
			goutils.Err(sgc7game.ErrInvalidStake))

		return nil, sgc7game.ErrInvalidStake
	}

	defer serv.game.DeleteGameData(gameData)

	cmd := req.Command

	serv.game.OnBet(plugin, cmd, req.ClientParams, ips, stake, results, gameData)

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := serv.game.Play(plugin, cmd, req.ClientParams, ips, stake, results, gameData)
		if err != nil {
			goutils.Error("Serv.onPlay:Play",
				slog.Int("results", len(results)),
				goutils.Err(err))

			return nil, err
		}

		if pr == nil {
			break
		}

		results = append(results, pr)
		if pr.IsFinish {
			break
		}

		if pr.IsWait {
			break
		}

		if len(pr.NextCmds) > 0 {
			cmd = pr.NextCmds[0]
		} else {
			cmd = ""
		}
	}

	pr := &sgc7pb.ReplyPlay{
		RandomNumbers: sgc7pbutils.BuildPBRngs(plugin.GetUsedRngs()),
	}

	ps, err := serv.service.BuildPBPlayerState(ips)
	if err != nil {
		goutils.Error("Serv.onPlay:BuildPlayerState",
			goutils.Err(err))

		return nil, err
	}

	pr.PlayerState = ps

	if len(results) > 0 {
		AddPlayResult(serv.service, pr, results)

		lastr := results[len(results)-1]

		pr.Finished = lastr.IsFinish
		pr.NextCommands = lastr.NextCmds
		pr.NextCommandParams = lastr.NextCmdParams
	}

	return pr, nil
}

// Play - play game
func (serv *Serv) LogReplyPlay(str string, reply *sgc7pb.ReplyPlay, logLevel zapcore.Level) {
	if logLevel == zapcore.DebugLevel {
		goutils.Debug(str,
			sgc7pbutils.JSON("reply", reply))
	} else if logLevel == zapcore.InfoLevel {
		goutils.Info(str,
			sgc7pbutils.JSON("reply", reply))
	}
}
