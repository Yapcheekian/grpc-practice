package main

import (
	"context"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/Yapcheekian/grpc-practice/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetReponse, error) {
	firstName := in.GetGreeting().FirstName
	if firstName == "" {
		return nil, status.Errorf(codes.InvalidArgument, "first name input is required")
	}
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
			return status.Errorf(codes.Internal, err.Error())
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := "Hello "
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return stream.SendAndClose(&greetpb.LongGreetResponse{
					Result: result,
				})
			}
			log.Printf("failed to receive client stream: %v\n", err)
			return status.Errorf(codes.Internal, err.Error())
		}

		firstName := req.GetGreeting().FirstName
		result += firstName + " "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			log.Printf("failed to receive bidi client stream: %v\n", err)
			return status.Errorf(codes.Internal, err.Error())
		}

		firstName := req.GetGreeting().FirstName
		result := "Hello " + firstName

		res := &greetpb.GreetEveryoneResponse{
			Result: result,
		}
		if err := stream.Send(res); err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
	}
}

func main() {
	l, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(l); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
