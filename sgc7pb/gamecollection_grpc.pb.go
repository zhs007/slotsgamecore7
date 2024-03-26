// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: gamecollection.proto

package sgc7pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GameLogicCollection_InitGame_FullMethodName             = "/sgc7pb.GameLogicCollection/initGame"
	GameLogicCollection_GetGameConfig_FullMethodName        = "/sgc7pb.GameLogicCollection/getGameConfig"
	GameLogicCollection_InitializeGamePlayer_FullMethodName = "/sgc7pb.GameLogicCollection/initializeGamePlayer"
	GameLogicCollection_PlayGame_FullMethodName             = "/sgc7pb.GameLogicCollection/playGame"
	GameLogicCollection_PlayGame2_FullMethodName            = "/sgc7pb.GameLogicCollection/playGame2"
)

// GameLogicCollectionClient is the client API for GameLogicCollection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameLogicCollectionClient interface {
	// initGame - initial game
	InitGame(ctx context.Context, in *RequestInitGame, opts ...grpc.CallOption) (*ReplyInitGame, error)
	// getGameConfig - get game config
	GetGameConfig(ctx context.Context, in *RequestGameConfig, opts ...grpc.CallOption) (*ReplyGameConfig, error)
	// initializeGamePlayer - initialize a game player
	InitializeGamePlayer(ctx context.Context, in *RequestInitializeGamePlayer, opts ...grpc.CallOption) (*ReplyInitializeGamePlayer, error)
	// playGame - play game
	PlayGame(ctx context.Context, in *RequestPlayGame, opts ...grpc.CallOption) (GameLogicCollection_PlayGameClient, error)
	// playGame2 - play game v2
	PlayGame2(ctx context.Context, in *RequestPlayGame, opts ...grpc.CallOption) (*ReplyPlayGame, error)
}

type gameLogicCollectionClient struct {
	cc grpc.ClientConnInterface
}

func NewGameLogicCollectionClient(cc grpc.ClientConnInterface) GameLogicCollectionClient {
	return &gameLogicCollectionClient{cc}
}

func (c *gameLogicCollectionClient) InitGame(ctx context.Context, in *RequestInitGame, opts ...grpc.CallOption) (*ReplyInitGame, error) {
	out := new(ReplyInitGame)
	err := c.cc.Invoke(ctx, GameLogicCollection_InitGame_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameLogicCollectionClient) GetGameConfig(ctx context.Context, in *RequestGameConfig, opts ...grpc.CallOption) (*ReplyGameConfig, error) {
	out := new(ReplyGameConfig)
	err := c.cc.Invoke(ctx, GameLogicCollection_GetGameConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameLogicCollectionClient) InitializeGamePlayer(ctx context.Context, in *RequestInitializeGamePlayer, opts ...grpc.CallOption) (*ReplyInitializeGamePlayer, error) {
	out := new(ReplyInitializeGamePlayer)
	err := c.cc.Invoke(ctx, GameLogicCollection_InitializeGamePlayer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameLogicCollectionClient) PlayGame(ctx context.Context, in *RequestPlayGame, opts ...grpc.CallOption) (GameLogicCollection_PlayGameClient, error) {
	stream, err := c.cc.NewStream(ctx, &GameLogicCollection_ServiceDesc.Streams[0], GameLogicCollection_PlayGame_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gameLogicCollectionPlayGameClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GameLogicCollection_PlayGameClient interface {
	Recv() (*ReplyPlayGame, error)
	grpc.ClientStream
}

type gameLogicCollectionPlayGameClient struct {
	grpc.ClientStream
}

func (x *gameLogicCollectionPlayGameClient) Recv() (*ReplyPlayGame, error) {
	m := new(ReplyPlayGame)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gameLogicCollectionClient) PlayGame2(ctx context.Context, in *RequestPlayGame, opts ...grpc.CallOption) (*ReplyPlayGame, error) {
	out := new(ReplyPlayGame)
	err := c.cc.Invoke(ctx, GameLogicCollection_PlayGame2_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameLogicCollectionServer is the server API for GameLogicCollection service.
// All implementations must embed UnimplementedGameLogicCollectionServer
// for forward compatibility
type GameLogicCollectionServer interface {
	// initGame - initial game
	InitGame(context.Context, *RequestInitGame) (*ReplyInitGame, error)
	// getGameConfig - get game config
	GetGameConfig(context.Context, *RequestGameConfig) (*ReplyGameConfig, error)
	// initializeGamePlayer - initialize a game player
	InitializeGamePlayer(context.Context, *RequestInitializeGamePlayer) (*ReplyInitializeGamePlayer, error)
	// playGame - play game
	PlayGame(*RequestPlayGame, GameLogicCollection_PlayGameServer) error
	// playGame2 - play game v2
	PlayGame2(context.Context, *RequestPlayGame) (*ReplyPlayGame, error)
	mustEmbedUnimplementedGameLogicCollectionServer()
}

// UnimplementedGameLogicCollectionServer must be embedded to have forward compatible implementations.
type UnimplementedGameLogicCollectionServer struct {
}

func (UnimplementedGameLogicCollectionServer) InitGame(context.Context, *RequestInitGame) (*ReplyInitGame, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitGame not implemented")
}
func (UnimplementedGameLogicCollectionServer) GetGameConfig(context.Context, *RequestGameConfig) (*ReplyGameConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGameConfig not implemented")
}
func (UnimplementedGameLogicCollectionServer) InitializeGamePlayer(context.Context, *RequestInitializeGamePlayer) (*ReplyInitializeGamePlayer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitializeGamePlayer not implemented")
}
func (UnimplementedGameLogicCollectionServer) PlayGame(*RequestPlayGame, GameLogicCollection_PlayGameServer) error {
	return status.Errorf(codes.Unimplemented, "method PlayGame not implemented")
}
func (UnimplementedGameLogicCollectionServer) PlayGame2(context.Context, *RequestPlayGame) (*ReplyPlayGame, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PlayGame2 not implemented")
}
func (UnimplementedGameLogicCollectionServer) mustEmbedUnimplementedGameLogicCollectionServer() {}

// UnsafeGameLogicCollectionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameLogicCollectionServer will
// result in compilation errors.
type UnsafeGameLogicCollectionServer interface {
	mustEmbedUnimplementedGameLogicCollectionServer()
}

func RegisterGameLogicCollectionServer(s grpc.ServiceRegistrar, srv GameLogicCollectionServer) {
	s.RegisterService(&GameLogicCollection_ServiceDesc, srv)
}

func _GameLogicCollection_InitGame_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInitGame)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicCollectionServer).InitGame(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogicCollection_InitGame_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicCollectionServer).InitGame(ctx, req.(*RequestInitGame))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameLogicCollection_GetGameConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestGameConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicCollectionServer).GetGameConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogicCollection_GetGameConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicCollectionServer).GetGameConfig(ctx, req.(*RequestGameConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameLogicCollection_InitializeGamePlayer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInitializeGamePlayer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicCollectionServer).InitializeGamePlayer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogicCollection_InitializeGamePlayer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicCollectionServer).InitializeGamePlayer(ctx, req.(*RequestInitializeGamePlayer))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameLogicCollection_PlayGame_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RequestPlayGame)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GameLogicCollectionServer).PlayGame(m, &gameLogicCollectionPlayGameServer{stream})
}

type GameLogicCollection_PlayGameServer interface {
	Send(*ReplyPlayGame) error
	grpc.ServerStream
}

type gameLogicCollectionPlayGameServer struct {
	grpc.ServerStream
}

func (x *gameLogicCollectionPlayGameServer) Send(m *ReplyPlayGame) error {
	return x.ServerStream.SendMsg(m)
}

func _GameLogicCollection_PlayGame2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPlayGame)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicCollectionServer).PlayGame2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogicCollection_PlayGame2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicCollectionServer).PlayGame2(ctx, req.(*RequestPlayGame))
	}
	return interceptor(ctx, in, info, handler)
}

// GameLogicCollection_ServiceDesc is the grpc.ServiceDesc for GameLogicCollection service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameLogicCollection_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sgc7pb.GameLogicCollection",
	HandlerType: (*GameLogicCollectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "initGame",
			Handler:    _GameLogicCollection_InitGame_Handler,
		},
		{
			MethodName: "getGameConfig",
			Handler:    _GameLogicCollection_GetGameConfig_Handler,
		},
		{
			MethodName: "initializeGamePlayer",
			Handler:    _GameLogicCollection_InitializeGamePlayer_Handler,
		},
		{
			MethodName: "playGame2",
			Handler:    _GameLogicCollection_PlayGame2_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "playGame",
			Handler:       _GameLogicCollection_PlayGame_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "gamecollection.proto",
}
