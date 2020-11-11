package dtserv

import (
	"context"
	"net"

	jsoniter "github.com/json-iterator/go"
	sgc7game "github.com/zhs007/slotsgamecore7/game"
	sgc7plugin "github.com/zhs007/slotsgamecore7/plugin"
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
	game     sgc7game.IGame
}

// NewServ -
func NewServ(service IService, game sgc7game.IGame, bindaddr string, version string) (*Serv, error) {
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
		game:     game,
	}

	sgc7pb.RegisterDTGameLogicServer(grpcServ, serv)

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
	sgc7utils.Debug("Serv.GetConfig",
		sgc7utils.JSON("req", req))

	cfg := serv.game.GetConfig()

	res := BuildPBGameConfig(cfg)

	sgc7utils.Debug("Serv.GetConfig",
		sgc7utils.JSON("reply", res))

	return res, nil
}

// Initialize - initialize a player
func (serv *Serv) Initialize(ctx context.Context, req *sgc7pb.RequestInitialize) (*sgc7pb.PlayerState, error) {
	sgc7utils.Debug("Serv.Initialize",
		sgc7utils.JSON("req", req))

	ps := serv.game.Initialize()
	res, err := serv.service.BuildPBPlayerState(ps)
	if err != nil {
		sgc7utils.Error("Serv.Initialize:BuildPBPlayerState",
			zap.Error(err))

		return nil, err
	}

	sgc7utils.Debug("Serv.Initialize",
		sgc7utils.JSON("reply", res))

	return res, nil
}

// Play - play game
func (serv *Serv) Play(req *sgc7pb.RequestPlay, stream sgc7pb.DTGameLogic_PlayServer) error {
	sgc7utils.Debug("Serv.Play",
		sgc7utils.JSON("req", req))

	res, err := serv.onPlay(req)
	if err != nil {
		sgc7utils.Error("Serv.Play:onPlay",
			zap.Error(err))

		return err
	}

	sgc7utils.Debug("Serv.Play",
		sgc7utils.JSON("reply", res))

	return stream.Send(res)
}

// ProcCheat - process cheat
func (serv *Serv) ProcCheat(plugin sgc7plugin.IPlugin, cheat string) error {
	if cheat != "" {
		str := sgc7utils.AppendString("[", cheat, "]")

		json := jsoniter.ConfigCompatibleWithStandardLibrary

		rngs := []int{}
		err := json.Unmarshal([]byte(str), &rngs)
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
		ips1, err := serv.service.BuildPlayerStateFromPB(req.PlayerState)
		if err != nil {
			sgc7utils.Error("Serv.onPlay:BuildPlayerStateFromPB",
				zap.Error(err))

			return nil, err
		}

		ips = ips1
	}

	plugin := serv.game.NewPlugin()
	defer serv.game.FreePlugin(plugin)

	serv.ProcCheat(plugin, req.Cheat)

	stake := BuildStake(req.Stake)
	err := serv.game.CheckStake(stake)
	if err != nil {
		sgc7utils.Error("Serv.onPlay:CheckStake",
			sgc7utils.JSON("stake", stake),
			zap.Error(err))

		return nil, err
	}

	results := []*sgc7game.PlayResult{}

	cmd := req.Command

	for {
		if cmd == "" {
			cmd = "SPIN"
		}

		pr, err := serv.game.Play(plugin, cmd, req.ClientParams, ips, stake, results)
		if err != nil {
			sgc7utils.Error("Serv.onPlay:Play",
				zap.Int("results", len(results)),
				zap.Error(err))

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
		RandomNumbers: BuildPBRngs(plugin.GetUsedRngs()),
		Stake:         req.Stake,
	}

	ps, err := serv.service.BuildPBPlayerState(ips)
	if err != nil {
		sgc7utils.Error("Serv.onPlay:BuildPlayerState",
			zap.Error(err))

		return nil, err
	}

	pr.PlayerState = ps

	if len(results) > 0 {
		AddPlayResult(serv.service, pr, results)

		lastr := results[len(results)-1]

		pr.Finished = lastr.IsFinish
		pr.NextCommands = lastr.NextCmds
	}

	return pr, nil
}
