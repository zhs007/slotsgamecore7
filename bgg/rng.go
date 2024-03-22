package bgg

import (
	"context"
	"log/slog"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RngClient - DTRngClient
type RngClient struct {
	servAddr         string
	gameCode         string
	conn             *grpc.ClientConn
	client           sgc7pb.RngClient
	useOpenTelemetry bool
}

// NewRngClient - new RngClient
func NewRngClient(servAddr string, gameCode string, useOpenTelemetry bool) *RngClient {
	client := &RngClient{
		servAddr:         servAddr,
		gameCode:         gameCode,
		useOpenTelemetry: useOpenTelemetry,
	}

	return client
}

// reset - reset
func (client *RngClient) reset() {
	if client.conn != nil {
		client.conn.Close()
	}

	client.conn = nil
	client.client = nil
}

// GetRngs - get rngs
func (client *RngClient) GetRngs(ctx context.Context, nums int) ([]uint32, error) {
	if client.conn == nil || client.client == nil {
		var conn *grpc.ClientConn
		var err error

		if client.useOpenTelemetry {
			conn, err = grpc.Dial(client.servAddr,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
				grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
				grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
			if err != nil {
				goutils.Error("RngClient.GetRngs:grpc.Dial",
					slog.String("server address", client.servAddr),
					goutils.Err(err))

				return nil, err
			}
		} else {
			conn, err = grpc.Dial(client.servAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				goutils.Error("RngClient.GetRngs:grpc.Dial",
					slog.String("server address", client.servAddr),
					goutils.Err(err))

				return nil, err
			}
		}

		client.conn = conn
		client.client = sgc7pb.NewRngClient(conn)
	}

	res, err := client.client.GetRngs(ctx, &sgc7pb.RequestRngs{
		Nums: int32(nums),
	})
	if err != nil {
		goutils.Error("RngClient.GetRngs:GetRngs",
			slog.String("server address", client.servAddr),
			slog.String("gamecode", client.gameCode),
			slog.Int("nums", nums),
			goutils.Err(err))

		client.reset()

		return nil, err
	}

	return res.Rngs, nil
}
