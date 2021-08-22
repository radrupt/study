package main

import (
	"context"
	"grpc-test/hello"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type HelloServiceImpl struct {
	hello.UnimplementedHelloServer
}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *hello.String,
) (*hello.String, error) {
	reply := &hello.String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloServiceImpl) Channel(stream hello.Hello_ChannelServer) error {
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		reply := &hello.String{Value: "hello:" + args.GetValue()}

		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	hello.RegisterHelloServer(s, &HelloServiceImpl{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
