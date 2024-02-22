package mathtoolset2

import (
	"context"
	"os"

	"github.com/zhs007/goutils"
	sgc7pb "github.com/zhs007/slotsgamecore7/sgc7pb"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client - MathToolsetClient
type Client struct {
	servAddr string
	conn     *grpc.ClientConn
	client   sgc7pb.MathToolsetClient
}

// NewClient - new MathToolsetClient
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

func (client *Client) onRequest(ctx context.Context) error {
	if client.conn == nil || client.client == nil {
		conn, err := grpc.Dial(client.servAddr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
		if err != nil {
			goutils.Error("Client.onRequest:grpc.Dial",
				zap.String("server address", client.servAddr),
				zap.Error(err))

			return err
		}

		client.conn = conn
		client.client = sgc7pb.NewMathToolsetClient(conn)
	}

	return nil
}

// RunScript - run script
func (client *Client) RunScript(ctx context.Context, script string, mapFiles map[string]string) (*sgc7pb.ReplyRunScript, error) {
	err := client.onRequest(ctx)
	if err != nil {
		goutils.Error("Client.RunScript:onRequest",
			zap.Error(err))

		return nil, err
	}

	mapfd, err := NewFileDataMap("")
	if err != nil {
		goutils.Error("Client.RunScript:NewFileDataMap",
			zap.Error(err))

		return nil, err
	}

	for k, v := range mapFiles {
		file, err := os.Open(v)
		if err != nil {
			goutils.Error("Client.RunScript:Open",
				zap.String("fn", v),
				zap.Error(err))

			return nil, err
		}

		mapfd.AddReader(k, file)
	}

	jsondata, err := mapfd.ToJson()
	if err != nil {
		goutils.Error("Client.RunScript:ToJson",
			zap.Error(err))

		return nil, err
	}

	res, err := client.client.RunScript(ctx, &sgc7pb.RunScript{
		Script:   script,
		MapFiles: jsondata,
	})
	if err != nil {
		goutils.Error("Client.InitGame:InitGame",
			zap.String("server address", client.servAddr),
			zap.String("script", script),
			zap.Error(err))

		client.reset()

		return nil, err
	}

	return res, nil
}
