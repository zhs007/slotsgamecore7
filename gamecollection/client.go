package gamecollection

import (
	"context"
	"io"
	"log/slog"

	"github.com/zhs007/goutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client - GameLogicCollectionClient
type Client struct {
	servAddr string
	conn     *grpc.ClientConn
	client   sgc7pb.GameLogicCollectionClient
}

// NewClient - new GameLogicClient
func NewClient(servAddr string) (*Client, error) {
	client := &Client{
		servAddr: servAddr,
	}

	return client, nil
}

// reset - reset
func (client *Client) reset() {
	if client.conn != nil {
		client.conn.Close()
	}

	client.conn = nil
	client.client = nil
}

func (client *Client) onRequest(_ context.Context) error {
	if client.conn == nil || client.client == nil {
	    conn, err := grpc.Dial(client.servAddr,
		    grpc.WithTransportCredentials(insecure.NewCredentials()),
		    // otelgrpc v0.63.0 uses stats handlers instead of interceptors
		    grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			goutils.Error("Client.onRequest:grpc.Dial",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			return err
		}

		client.conn = conn
		client.client = sgc7pb.NewGameLogicCollectionClient(conn)
	}

	return nil
}

// InitGame - init game
func (client *Client) InitGame(ctx context.Context, gameCode string, data string) (*sgc7pb.ReplyInitGame, error) {
	err := client.onRequest(ctx)
	if err != nil {
		goutils.Error("Client.InitGame:onRequest",
			goutils.Err(err))

		return nil, err
	}

	res, err := client.client.InitGame(ctx, &sgc7pb.RequestInitGame{
		GameCode: gameCode,
		Config:   data,
	})
	if err != nil {
		goutils.Error("Client.InitGame:InitGame",
			slog.String("server address", client.servAddr),
			slog.String("gameCode", gameCode),
			slog.String("data", data),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res, nil
}

// GetGameConfig - get config
func (client *Client) GetGameConfig(ctx context.Context, gameCode string) (*sgc7pb.ReplyGameConfig, error) {
	err := client.onRequest(ctx)
	if err != nil {
		goutils.Error("Client.GetGameConfig:onRequest",
			goutils.Err(err))

		return nil, err
	}

	res, err := client.client.GetGameConfig(ctx, &sgc7pb.RequestGameConfig{
		GameCode: gameCode,
	})
	if err != nil {
		goutils.Error("Client.GetGameConfig:GetGameConfig",
			slog.String("server address", client.servAddr),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res, nil
}

// InitializeGamePlayer - initialize a player
func (client *Client) InitializeGamePlayer(ctx context.Context, gameCode string) (*sgc7pb.ReplyInitializeGamePlayer, error) {
	err := client.onRequest(ctx)
	if err != nil {
		goutils.Error("Client.InitializeGamePlayer:onRequest",
			goutils.Err(err))

		return nil, err
	}

	res, err := client.client.InitializeGamePlayer(ctx, &sgc7pb.RequestInitializeGamePlayer{
		GameCode: gameCode,
	})
	if err != nil {
		goutils.Error("Client.InitializeGamePlayer:InitializeGamePlayer",
			slog.String("server address", client.servAddr),
			slog.String("gameCode", gameCode),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res, nil
}

// PlayGame - play game
func (client *Client) PlayGame(ctx context.Context, gameCode string, ps *sgc7pb.PlayerState,
	cheat string, stake *sgc7pb.Stake, clientParams string, cmd string) (*sgc7pb.ReplyPlayGame, error) {

	err := client.onRequest(ctx)
	if err != nil {
		goutils.Error("Client.PlayGame:onRequest",
			goutils.Err(err))

		return nil, err
	}

	stream, err := client.client.PlayGame(ctx, &sgc7pb.RequestPlayGame{
		GameCode: gameCode,
		Play: &sgc7pb.RequestPlay{
			PlayerState:  ps,
			Cheat:        cheat,
			Stake:        stake,
			ClientParams: clientParams,
			Command:      cmd,
		},
	})
	if err != nil {
		goutils.Error("Client.PlayGame:PlayGame",
			slog.String("server address", client.servAddr),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	// reply := &sgc7pb.ReplyPlayGame{}
	var reply *sgc7pb.ReplyPlayGame
	for {
		rp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return reply, nil
			}

			goutils.Error("Client.PlayGame:Recv",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			client.reset()

			return nil, err
		} else {
			reply = rp
		}

		// return rp, nil

		// if rp != nil {
		// 	sgc7pbutils.MergeReplyPlay(reply, rp)
		// }
	}
}
