package bgg

import (
	"context"

	goutils "github.com/zhs007/goutils"
	"github.com/zhs007/slotsgamecore7/bggrngpb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RngClient - DTRngClient
type RngClient struct {
	servAddr         string
	gameCode         string
	conn             *grpc.ClientConn
	client           bggrngpb.BGGRngClient
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
					zap.String("server address", client.servAddr),
					zap.Error(err))

				return nil, err
			}
		} else {
			conn, err = grpc.Dial(client.servAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				goutils.Error("RngClient.GetRngs:grpc.Dial",
					zap.String("server address", client.servAddr),
					zap.Error(err))

				return nil, err
			}
		}

		client.conn = conn
		client.client = bggrngpb.NewBGGRngClient(conn)
	}

	res, err := client.client.GetRngs(ctx, &bggrngpb.RequestRngs{
		Nums: int32(nums),
	})
	if err != nil {
		goutils.Error("RngClient.GetRngs:GetRngs",
			zap.String("server address", client.servAddr),
			zap.String("gamecode", client.gameCode),
			zap.Int("nums", nums),
			zap.Error(err))

		client.reset()

		return nil, err
	}

	return res.Rngs, nil
}
