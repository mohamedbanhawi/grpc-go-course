package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
)

func (s *Server) Greet(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Printf("Greet function invoked with %v\n", in)

	return &pb.GreetResponse{
		Result: "Hello " + in.FirstName,
	}, nil
}

func (s *Server) GreetStream(in *pb.GreetRequest, stream pb.GreetService_GreetStreamServer) error {
	log.Printf("GreetStream function invoked with %v\n", in)

	for i := 0; i < 10; i++ {
		stream.Send(&pb.GreetResponse{
			Result: fmt.Sprintf("Hello %s, #%d", in.FirstName, i)})
	}

	return nil
}
