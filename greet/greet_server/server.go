package main

import (
	"context"
	"log"
	"net"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetReponse, error) {
	firstName := in.GetGreeting().FirstName
	result := "Hello " + firstName
	res := greetpb.GreetReponse{
		Result: result,
	}

	return &res, nil
}

func main() {
	l, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
