package main

import (
	"context"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetReponse, error) {
	firstName := in.GetGreeting().FirstName
	result := "Hello " + firstName
	res := greetpb.GreetReponse{
		Result: result,
	}

	return &res, nil
}

func (*server) GreetManyTimes(in *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := in.GetGreeting().FirstName
	for i := 0; i < 3; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: "Hello " + firstName + " number " + strconv.Itoa(i),
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func main() {
	l, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
