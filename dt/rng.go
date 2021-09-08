package dt

import (
	"context"

	goutils "github.com/zhs007/goutils"
	dtrngpb "github.com/zhs007/slotsgamecore7/dtrngpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RngClient - DTRngClient
type RngClient struct {
	servAddr string
	gameCode string
	conn     *grpc.ClientConn
	client   dtrngpb.DTRngClient
}

// NewRngClient - new RngClient
func NewRngClient(servAddr string, gameCode string) *RngClient {
	client := &RngClient{
		servAddr: servAddr,
		gameCode: gameCode,
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
		conn, err := grpc.Dial(client.servAddr, grpc.WithInsecure())
		if err != nil {
			goutils.Error("RngClient.GetRngs:grpc.Dial",
				zap.String("server address", client.servAddr),
				zap.Error(err))

			return nil, err
		}

		client.conn = conn
		client.client = dtrngpb.NewDTRngClient(conn)
	}

	res, err := client.client.GetRngs(ctx, &dtrngpb.RequestRngs{
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
