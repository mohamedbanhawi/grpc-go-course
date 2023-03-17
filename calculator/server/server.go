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

func (s *Server) Primes(in *pb.PrimesRequest, stream pb.CalculateService_PrimesServer) error {
	log.Printf("Recieved Primes on server")

	var number int32 = in.GetNumber()

	var factor int32 = 2
	for number > 1 {
		if number%factor == 0 {
			// if factor evenly divides into N
			err := stream.SendMsg(&pb.PrimesResponse{Result: factor}) // this is a factor
			if err != nil {
				log.Printf("Failed to send %d as a factor for %d\n%v", factor, number, err)
				return err
			}
			number /= factor // divide N by factor so that we have the rest of the number left.
		} else {
			factor++
		}
	}
	return nil
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
