// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package hello

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

// HelloClient is the client API for Hello service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HelloClient interface {
	Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error)
	Channel(ctx context.Context, opts ...grpc.CallOption) (Hello_ChannelClient, error)
}

type helloClient struct {
	cc grpc.ClientConnInterface
}

func NewHelloClient(cc grpc.ClientConnInterface) HelloClient {
	return &helloClient{cc}
}

func (c *helloClient) Hello(ctx context.Context, in *String, opts ...grpc.CallOption) (*String, error) {
	out := new(String)
	err := c.cc.Invoke(ctx, "/hello.Hello/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *helloClient) Channel(ctx context.Context, opts ...grpc.CallOption) (Hello_ChannelClient, error) {
	stream, err := c.cc.NewStream(ctx, &Hello_ServiceDesc.Streams[0], "/hello.Hello/Channel", opts...)
	if err != nil {
		return nil, err
	}
	x := &helloChannelClient{stream}
	return x, nil
}

type Hello_ChannelClient interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ClientStream
}

type helloChannelClient struct {
	grpc.ClientStream
}

func (x *helloChannelClient) Send(m *String) error {
	return x.ClientStream.SendMsg(m)
}

func (x *helloChannelClient) Recv() (*String, error) {
	m := new(String)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// HelloServer is the server API for Hello service.
// All implementations must embed UnimplementedHelloServer
// for forward compatibility
type HelloServer interface {
	Hello(context.Context, *String) (*String, error)
	Channel(Hello_ChannelServer) error
	mustEmbedUnimplementedHelloServer()
}

// UnimplementedHelloServer must be embedded to have forward compatible implementations.
type UnimplementedHelloServer struct {
}

func (UnimplementedHelloServer) Hello(context.Context, *String) (*String, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedHelloServer) Channel(Hello_ChannelServer) error {
	return status.Errorf(codes.Unimplemented, "method Channel not implemented")
}
func (UnimplementedHelloServer) mustEmbedUnimplementedHelloServer() {}

// UnsafeHelloServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HelloServer will
// result in compilation errors.
type UnsafeHelloServer interface {
	mustEmbedUnimplementedHelloServer()
}

func RegisterHelloServer(s grpc.ServiceRegistrar, srv HelloServer) {
	s.RegisterService(&Hello_ServiceDesc, srv)
}

func _Hello_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(String)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HelloServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/hello.Hello/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HelloServer).Hello(ctx, req.(*String))
	}
	return interceptor(ctx, in, info, handler)
}

func _Hello_Channel_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(HelloServer).Channel(&helloChannelServer{stream})
}

type Hello_ChannelServer interface {
	Send(*String) error
	Recv() (*String, error)
	grpc.ServerStream
}

type helloChannelServer struct {
	grpc.ServerStream
}

func (x *helloChannelServer) Send(m *String) error {
	return x.ServerStream.SendMsg(m)
}

func (x *helloChannelServer) Recv() (*String, error) {
	m := new(String)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Hello_ServiceDesc is the grpc.ServiceDesc for Hello service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Hello_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "hello.Hello",
	HandlerType: (*HelloServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _Hello_Hello_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Channel",
			Handler:       _Hello_Channel_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "hello/hello.proto",
}
