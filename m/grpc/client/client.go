package main

import (
	"context"
	"fmt"
	"grpc-test/hello"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50052"
	defaultName = "hello"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := hello.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	// r, err := c.Hello(ctx, &hello.String{Value: "Hello"})
	// if err != nil {
	// 	log.Fatalf("could not hello: %v", err)
	// }
	stream, err := c.Channel(ctx)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			if err := stream.Send(&hello.String{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
	}()
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
	// log.Printf("helloing: %s", r.GetValue())
}
