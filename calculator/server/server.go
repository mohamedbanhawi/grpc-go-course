package main

import (
	"context"
	"log"
	"net"

	pb "github.com/mohamedbanhawi/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
)

var addr string = "0.0.0.0:50051"

// Implement CalculatorService Server

type Server struct {
	pb.CalculateServiceServer
}

// Implement Sum Method
func (s *Server) Sum(ctx context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {

	log.Printf("Recieved sum on server")

	return &pb.SumResponse{
		Result: request.GetFirstNumber() + request.GetSecondNumber(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listening on %s\n", addr)

	// grpc Server
	s := grpc.NewServer()

	pb.RegisterCalculateServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}

}
