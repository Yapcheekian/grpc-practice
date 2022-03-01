package main

import (
	"context"
	"fmt"
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

	res, err := c.Greet(context.Background(), &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yap",
		},
	})
	if err != nil {
		log.Fatalf("failed to greet unary: %v", err)
	}

	fmt.Println(res)
}
