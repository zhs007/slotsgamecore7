// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: game.proto

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
	DTGameLogic_GetConfig_FullMethodName  = "/sgc7pb.DTGameLogic/getConfig"
	DTGameLogic_Initialize_FullMethodName = "/sgc7pb.DTGameLogic/initialize"
	DTGameLogic_Play_FullMethodName       = "/sgc7pb.DTGameLogic/play"
	DTGameLogic_Play2_FullMethodName      = "/sgc7pb.DTGameLogic/play2"
)

// DTGameLogicClient is the client API for DTGameLogic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DTGameLogicClient interface {
	// getConfig - get config
	GetConfig(ctx context.Context, in *RequestConfig, opts ...grpc.CallOption) (*GameConfig, error)
	// initialize - initialize a player
	Initialize(ctx context.Context, in *RequestInitialize, opts ...grpc.CallOption) (*PlayerState, error)
	// play - play game
	Play(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (DTGameLogic_PlayClient, error)
	// play2 - play game v2
	Play2(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (*ReplyPlay, error)
}

type dTGameLogicClient struct {
	cc grpc.ClientConnInterface
}

func NewDTGameLogicClient(cc grpc.ClientConnInterface) DTGameLogicClient {
	return &dTGameLogicClient{cc}
}

func (c *dTGameLogicClient) GetConfig(ctx context.Context, in *RequestConfig, opts ...grpc.CallOption) (*GameConfig, error) {
	out := new(GameConfig)
	err := c.cc.Invoke(ctx, DTGameLogic_GetConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dTGameLogicClient) Initialize(ctx context.Context, in *RequestInitialize, opts ...grpc.CallOption) (*PlayerState, error) {
	out := new(PlayerState)
	err := c.cc.Invoke(ctx, DTGameLogic_Initialize_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dTGameLogicClient) Play(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (DTGameLogic_PlayClient, error) {
	stream, err := c.cc.NewStream(ctx, &DTGameLogic_ServiceDesc.Streams[0], DTGameLogic_Play_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &dTGameLogicPlayClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DTGameLogic_PlayClient interface {
	Recv() (*ReplyPlay, error)
	grpc.ClientStream
}

type dTGameLogicPlayClient struct {
	grpc.ClientStream
}

func (x *dTGameLogicPlayClient) Recv() (*ReplyPlay, error) {
	m := new(ReplyPlay)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dTGameLogicClient) Play2(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (*ReplyPlay, error) {
	out := new(ReplyPlay)
	err := c.cc.Invoke(ctx, DTGameLogic_Play2_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DTGameLogicServer is the server API for DTGameLogic service.
// All implementations must embed UnimplementedDTGameLogicServer
// for forward compatibility
type DTGameLogicServer interface {
	// getConfig - get config
	GetConfig(context.Context, *RequestConfig) (*GameConfig, error)
	// initialize - initialize a player
	Initialize(context.Context, *RequestInitialize) (*PlayerState, error)
	// play - play game
	Play(*RequestPlay, DTGameLogic_PlayServer) error
	// play2 - play game v2
	Play2(context.Context, *RequestPlay) (*ReplyPlay, error)
	mustEmbedUnimplementedDTGameLogicServer()
}

// UnimplementedDTGameLogicServer must be embedded to have forward compatible implementations.
type UnimplementedDTGameLogicServer struct {
}

func (UnimplementedDTGameLogicServer) GetConfig(context.Context, *RequestConfig) (*GameConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedDTGameLogicServer) Initialize(context.Context, *RequestInitialize) (*PlayerState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Initialize not implemented")
}
func (UnimplementedDTGameLogicServer) Play(*RequestPlay, DTGameLogic_PlayServer) error {
	return status.Errorf(codes.Unimplemented, "method Play not implemented")
}
func (UnimplementedDTGameLogicServer) Play2(context.Context, *RequestPlay) (*ReplyPlay, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Play2 not implemented")
}
func (UnimplementedDTGameLogicServer) mustEmbedUnimplementedDTGameLogicServer() {}

// UnsafeDTGameLogicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DTGameLogicServer will
// result in compilation errors.
type UnsafeDTGameLogicServer interface {
	mustEmbedUnimplementedDTGameLogicServer()
}

func RegisterDTGameLogicServer(s grpc.ServiceRegistrar, srv DTGameLogicServer) {
	s.RegisterService(&DTGameLogic_ServiceDesc, srv)
}

func _DTGameLogic_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DTGameLogicServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DTGameLogic_GetConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DTGameLogicServer).GetConfig(ctx, req.(*RequestConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _DTGameLogic_Initialize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInitialize)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DTGameLogicServer).Initialize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DTGameLogic_Initialize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DTGameLogicServer).Initialize(ctx, req.(*RequestInitialize))
	}
	return interceptor(ctx, in, info, handler)
}

func _DTGameLogic_Play_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RequestPlay)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DTGameLogicServer).Play(m, &dTGameLogicPlayServer{stream})
}

type DTGameLogic_PlayServer interface {
	Send(*ReplyPlay) error
	grpc.ServerStream
}

type dTGameLogicPlayServer struct {
	grpc.ServerStream
}

func (x *dTGameLogicPlayServer) Send(m *ReplyPlay) error {
	return x.ServerStream.SendMsg(m)
}

func _DTGameLogic_Play2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPlay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DTGameLogicServer).Play2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DTGameLogic_Play2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DTGameLogicServer).Play2(ctx, req.(*RequestPlay))
	}
	return interceptor(ctx, in, info, handler)
}

// DTGameLogic_ServiceDesc is the grpc.ServiceDesc for DTGameLogic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DTGameLogic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sgc7pb.DTGameLogic",
	HandlerType: (*DTGameLogicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getConfig",
			Handler:    _DTGameLogic_GetConfig_Handler,
		},
		{
			MethodName: "initialize",
			Handler:    _DTGameLogic_Initialize_Handler,
		},
		{
			MethodName: "play2",
			Handler:    _DTGameLogic_Play2_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "play",
			Handler:       _DTGameLogic_Play_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "game.proto",
}

const (
	GameLogic_GetConfig_FullMethodName  = "/sgc7pb.GameLogic/getConfig"
	GameLogic_Initialize_FullMethodName = "/sgc7pb.GameLogic/initialize"
	GameLogic_Play_FullMethodName       = "/sgc7pb.GameLogic/play"
	GameLogic_Play2_FullMethodName      = "/sgc7pb.GameLogic/play2"
)

// GameLogicClient is the client API for GameLogic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GameLogicClient interface {
	// getConfig - get config
	GetConfig(ctx context.Context, in *RequestConfig, opts ...grpc.CallOption) (*GameConfig, error)
	// initialize - initialize a player
	Initialize(ctx context.Context, in *RequestInitialize, opts ...grpc.CallOption) (*PlayerState, error)
	// play - play game
	Play(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (GameLogic_PlayClient, error)
	// play2 - play game v2
	Play2(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (*ReplyPlay, error)
}

type gameLogicClient struct {
	cc grpc.ClientConnInterface
}

func NewGameLogicClient(cc grpc.ClientConnInterface) GameLogicClient {
	return &gameLogicClient{cc}
}

func (c *gameLogicClient) GetConfig(ctx context.Context, in *RequestConfig, opts ...grpc.CallOption) (*GameConfig, error) {
	out := new(GameConfig)
	err := c.cc.Invoke(ctx, GameLogic_GetConfig_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameLogicClient) Initialize(ctx context.Context, in *RequestInitialize, opts ...grpc.CallOption) (*PlayerState, error) {
	out := new(PlayerState)
	err := c.cc.Invoke(ctx, GameLogic_Initialize_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gameLogicClient) Play(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (GameLogic_PlayClient, error) {
	stream, err := c.cc.NewStream(ctx, &GameLogic_ServiceDesc.Streams[0], GameLogic_Play_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &gameLogicPlayClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GameLogic_PlayClient interface {
	Recv() (*ReplyPlay, error)
	grpc.ClientStream
}

type gameLogicPlayClient struct {
	grpc.ClientStream
}

func (x *gameLogicPlayClient) Recv() (*ReplyPlay, error) {
	m := new(ReplyPlay)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *gameLogicClient) Play2(ctx context.Context, in *RequestPlay, opts ...grpc.CallOption) (*ReplyPlay, error) {
	out := new(ReplyPlay)
	err := c.cc.Invoke(ctx, GameLogic_Play2_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GameLogicServer is the server API for GameLogic service.
// All implementations must embed UnimplementedGameLogicServer
// for forward compatibility
type GameLogicServer interface {
	// getConfig - get config
	GetConfig(context.Context, *RequestConfig) (*GameConfig, error)
	// initialize - initialize a player
	Initialize(context.Context, *RequestInitialize) (*PlayerState, error)
	// play - play game
	Play(*RequestPlay, GameLogic_PlayServer) error
	// play2 - play game v2
	Play2(context.Context, *RequestPlay) (*ReplyPlay, error)
	mustEmbedUnimplementedGameLogicServer()
}

// UnimplementedGameLogicServer must be embedded to have forward compatible implementations.
type UnimplementedGameLogicServer struct {
}

func (UnimplementedGameLogicServer) GetConfig(context.Context, *RequestConfig) (*GameConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConfig not implemented")
}
func (UnimplementedGameLogicServer) Initialize(context.Context, *RequestInitialize) (*PlayerState, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Initialize not implemented")
}
func (UnimplementedGameLogicServer) Play(*RequestPlay, GameLogic_PlayServer) error {
	return status.Errorf(codes.Unimplemented, "method Play not implemented")
}
func (UnimplementedGameLogicServer) Play2(context.Context, *RequestPlay) (*ReplyPlay, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Play2 not implemented")
}
func (UnimplementedGameLogicServer) mustEmbedUnimplementedGameLogicServer() {}

// UnsafeGameLogicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GameLogicServer will
// result in compilation errors.
type UnsafeGameLogicServer interface {
	mustEmbedUnimplementedGameLogicServer()
}

func RegisterGameLogicServer(s grpc.ServiceRegistrar, srv GameLogicServer) {
	s.RegisterService(&GameLogic_ServiceDesc, srv)
}

func _GameLogic_GetConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicServer).GetConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogic_GetConfig_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicServer).GetConfig(ctx, req.(*RequestConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameLogic_Initialize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestInitialize)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicServer).Initialize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogic_Initialize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicServer).Initialize(ctx, req.(*RequestInitialize))
	}
	return interceptor(ctx, in, info, handler)
}

func _GameLogic_Play_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RequestPlay)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GameLogicServer).Play(m, &gameLogicPlayServer{stream})
}

type GameLogic_PlayServer interface {
	Send(*ReplyPlay) error
	grpc.ServerStream
}

type gameLogicPlayServer struct {
	grpc.ServerStream
}

func (x *gameLogicPlayServer) Send(m *ReplyPlay) error {
	return x.ServerStream.SendMsg(m)
}

func _GameLogic_Play2_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPlay)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GameLogicServer).Play2(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GameLogic_Play2_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GameLogicServer).Play2(ctx, req.(*RequestPlay))
	}
	return interceptor(ctx, in, info, handler)
}

// GameLogic_ServiceDesc is the grpc.ServiceDesc for GameLogic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GameLogic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sgc7pb.GameLogic",
	HandlerType: (*GameLogicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getConfig",
			Handler:    _GameLogic_GetConfig_Handler,
		},
		{
			MethodName: "initialize",
			Handler:    _GameLogic_Initialize_Handler,
		},
		{
			MethodName: "play2",
			Handler:    _GameLogic_Play2_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "play",
			Handler:       _GameLogic_Play_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "game.proto",
}
