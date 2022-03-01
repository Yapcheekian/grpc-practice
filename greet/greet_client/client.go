package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// Unary
	res, err := c.Greet(context.Background(), &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yap",
		},
	})
	if err != nil {
		log.Fatalf("failed to greet unary: %v\n", err)
	}

	fmt.Println("unary: ", res)

	// Streaming server
	stream, err := c.GreetManyTimes(context.Background(), &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yap",
		},
	})
	if err != nil {
		log.Fatalf("failed to greet server streaming: %v\n", err)
	}

	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("faild to receive stream: %v\n", err)
		}

		fmt.Println("steam: ", msg)
	}
}
