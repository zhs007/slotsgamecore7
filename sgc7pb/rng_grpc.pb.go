// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: rng.proto

package sgc7pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Rng_GetRngs_FullMethodName = "/sgc7pb.Rng/getRngs"
)

// RngClient is the client API for Rng service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Rng - RNG Service
type RngClient interface {
	// getRngs - get rngs
	GetRngs(ctx context.Context, in *RequestRngs, opts ...grpc.CallOption) (*ReplyRngs, error)
}

type rngClient struct {
	cc grpc.ClientConnInterface
}

func NewRngClient(cc grpc.ClientConnInterface) RngClient {
	return &rngClient{cc}
}

func (c *rngClient) GetRngs(ctx context.Context, in *RequestRngs, opts ...grpc.CallOption) (*ReplyRngs, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReplyRngs)
	err := c.cc.Invoke(ctx, Rng_GetRngs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RngServer is the server API for Rng service.
// All implementations must embed UnimplementedRngServer
// for forward compatibility.
//
// Rng - RNG Service
type RngServer interface {
	// getRngs - get rngs
	GetRngs(context.Context, *RequestRngs) (*ReplyRngs, error)
	mustEmbedUnimplementedRngServer()
}

// UnimplementedRngServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRngServer struct{}

func (UnimplementedRngServer) GetRngs(context.Context, *RequestRngs) (*ReplyRngs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRngs not implemented")
}
func (UnimplementedRngServer) mustEmbedUnimplementedRngServer() {}
func (UnimplementedRngServer) testEmbeddedByValue()             {}

// UnsafeRngServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RngServer will
// result in compilation errors.
type UnsafeRngServer interface {
	mustEmbedUnimplementedRngServer()
}

func RegisterRngServer(s grpc.ServiceRegistrar, srv RngServer) {
	// If the following call pancis, it indicates UnimplementedRngServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Rng_ServiceDesc, srv)
}

func _Rng_GetRngs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestRngs)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RngServer).GetRngs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Rng_GetRngs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RngServer).GetRngs(ctx, req.(*RequestRngs))
	}
	return interceptor(ctx, in, info, handler)
}

// Rng_ServiceDesc is the grpc.ServiceDesc for Rng service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Rng_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sgc7pb.Rng",
	HandlerType: (*RngServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getRngs",
			Handler:    _Rng_GetRngs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rng.proto",
}
