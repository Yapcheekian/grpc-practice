package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("ssl/ca.crt", "")
	if err != nil {
		log.Fatalf("failed to load tls from file: %v", err)
	}

	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(creds))
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
		respErr, ok := status.FromError(err)
		if ok {
			// user error
			fmt.Println(respErr.Message())
			return
		} else {
			log.Fatalf("failed to greet unary: %v\n", err)
		}
	}

	fmt.Println("unary: ", res)

	// Server Streaming
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

		fmt.Println("server steaming: ", msg)
	}

	// Client streaming
	clientStream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to greet client streaming: %v\n", err)
	}

	for i := 0; i < 3; i++ {
		req := &greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "yap " + strconv.Itoa(i),
			},
		}
		if err := clientStream.Send(req); err != nil {
			log.Fatalf("failed to send client stream: %v\n", err)
		}
	}

	res2, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to close client stream: %v\n", err)
	}

	fmt.Println("client streaming: ", res2)

	// BiDi Streaming
	biDiStream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("fail to greet bidi streaming: %v\n", err)
	}

	waitCh := make(chan struct{})

	go func() {
		for i := 0; i < 3; i++ {
			req := &greetpb.GreetEveryoneRequest{
				Greeting: &greetpb.Greeting{
					FirstName: "yap " + strconv.Itoa(i),
				},
			}
			if err := biDiStream.Send(req); err != nil {
				log.Fatalf("failed to send bidi streaming: %v", err)
			}
		}
		if err := biDiStream.CloseSend(); err != nil {
			log.Fatalf("failed to close client stream: %v", err)
		}
	}()

	go func() {
		for {
			res, err := biDiStream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("failed to receive bidi streaming: %v", err)
				break
			}

			fmt.Println("bidi streaming: ", res)
		}
		close(waitCh)
	}()

	<-waitCh

	// Unary with deadline
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	resDeadline, err := c.GreetWithDeadline(ctx, &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "yap",
		},
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			if respErr.Code() == codes.DeadlineExceeded {
				fmt.Println("timeout was hit, deadline exceeded")
			} else {
				fmt.Println(respErr.Message())
			}
			return
		} else {
			log.Fatalf("failed to greet unary: %v\n", err)
		}
	}

	fmt.Println("unary with deadline: ", resDeadline)
}
