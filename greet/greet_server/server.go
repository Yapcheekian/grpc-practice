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
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, in *greetpb.GreetRequest) (*greetpb.GreetReponse, error) {
	firstName := in.GetGreeting().FirstName
	if firstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first name input is required")
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
			return status.Error(codes.Internal, err.Error())
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
			return status.Error(codes.Internal, err.Error())
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
			return status.Error(codes.Internal, err.Error())
		}

		firstName := req.GetGreeting().FirstName
		result := "Hello " + firstName

		res := &greetpb.GreetEveryoneResponse{
			Result: result,
		}
		if err := stream.Send(res); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (*server) GreetWithDeadline(ctx context.Context, in *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	select {
	case <-ctx.Done():
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	default:
		firstName := in.GetGreeting().FirstName
		if firstName == "" {
			return nil, status.Error(codes.InvalidArgument, "first name input is required")
		}
		result := "Hello " + firstName
		res := greetpb.GreetWithDeadlineResponse{
			Result: result,
		}

		return &res, nil
	}
}

func main() {
	l, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
	}

	creds, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")
	if err != nil {
		log.Printf("failed to load tls from file: %v\n", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	greetpb.RegisterGreetServiceServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(l); err != nil {
		log.Printf("failed to serve: %v\n", err)
	}
}
