package main

import (
	"log"
	"net"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &greetpb.UnimplementedGreetServiceServer{})

	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
