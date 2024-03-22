package grpcserv

import (
	"context"
	"io"
	"log/slog"

	goutils "github.com/zhs007/goutils"
	sgc7pbutils "github.com/zhs007/slotsgamecore7/pbutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"google.golang.org/grpc"
)

// Client - GameLogicClient
type Client struct {
	servAddr string
	conn     *grpc.ClientConn
	client   sgc7pb.GameLogicClient
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

// GetConfig - get config
func (client *Client) GetConfig(ctx context.Context) (*sgc7pb.GameConfig, error) {
	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			goutils.Error("Client.GetConfig:grpc.Dial",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			return nil, err
		}

		client.conn = conn
		client.client = sgc7pb.NewGameLogicClient(conn)
	}

	res, err := client.client.GetConfig(ctx, &sgc7pb.RequestConfig{})
	if err != nil {
		goutils.Error("Client.GetConfig:GetConfig",
			slog.String("server address", client.servAddr),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res, nil
}

// Initialize - initialize a player
func (client *Client) Initialize(ctx context.Context) (*sgc7pb.PlayerState, error) {
	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			goutils.Error("Client.Initialize:grpc.Dial",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			return nil, err
		}

		client.conn = conn
		client.client = sgc7pb.NewGameLogicClient(conn)
	}

	res, err := client.client.Initialize(ctx, &sgc7pb.RequestInitialize{})
	if err != nil {
		goutils.Error("Client.Initialize:Initialize",
			slog.String("server address", client.servAddr),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res, nil
}

// Play - play game
func (client *Client) Play(ctx context.Context, ps *sgc7pb.PlayerState,
	cheat string,
	stake *sgc7pb.Stake,
	clientParams string,
	cmd string) (*sgc7pb.ReplyPlay, error) {

	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			goutils.Error("Client.Play:grpc.Dial",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			return nil, err
		}

		client.conn = conn
		client.client = sgc7pb.NewGameLogicClient(conn)
	}

	stream, err := client.client.Play(ctx, &sgc7pb.RequestPlay{
		PlayerState:  ps,
		Cheat:        cheat,
		Stake:        stake,
		ClientParams: clientParams,
		Command:      cmd,
	})
	if err != nil {
		goutils.Error("Client.Play:Play",
			slog.String("server address", client.servAddr),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	reply := &sgc7pb.ReplyPlay{}

	for {
		rp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return reply, nil
			}

			goutils.Error("Client.Play:Recv",
				slog.String("server address", client.servAddr),
				goutils.Err(err))

			client.reset()

			return nil, err
		}

		if rp != nil {
			sgc7pbutils.MergeReplyPlay(reply, rp)
		}
	}
}
