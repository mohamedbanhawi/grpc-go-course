package main

import (
	"context"
	"io"
	"log"
	"math"
	"net"

	pb "github.com/mohamedbanhawi/grpc-go-course/calculator/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var addr string = "0.0.0.0:50051"

// Implement CalculatorService Server

type Server struct {
	pb.CalculateServiceServer
}

// Implement CalculateServiceServer Interface Concrete Methods

func (s *Server) Sum(ctx context.Context, request *pb.SumRequest) (*pb.SumResponse, error) {

	log.Printf("Invoked sum on server\n")

	return &pb.SumResponse{
		Result: request.GetFirstNumber() + request.GetSecondNumber(),
	}, nil
}

func (s *Server) Primes(in *pb.PrimesRequest, stream pb.CalculateService_PrimesServer) error {
	log.Printf("Invoked Primes on server\n")

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

func (s *Server) Average(stream pb.CalculateService_AverageServer) error {

	log.Printf("Invoked Average on server\n")
	var numbers []float32

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			var total float32 = 0
			for _, value := range numbers {
				total += value
			}
			if len(numbers) <= 0 {
				stream.SendAndClose(&pb.AverageResponse{})
			}
			stream.SendAndClose(&pb.AverageResponse{Result: total / float32(len(numbers))})
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to recieve number request %v", err)
			return err
		}
		log.Printf("Recieved number %v to Average on server\n", req)
		numbers = append(numbers, req.Number)
	}

}

func (s *Server) Max(stream pb.CalculateService_MaxServer) error {

	log.Printf("Invoked Max on server\n")
	var numbers []float32

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to recieve max number request %v", err)
			return err
		}
		log.Printf("Recieved number %v to max on server\n", req)
		numbers = append(numbers, req.Number)

		// maximum number closure
		max := func() float32 {
			var max float32
			for i, e := range numbers {
				if i == 0 || e > max {
					max = e
				}
			}
			return max
		}()
		stream.Send(&pb.MaxResponse{Result: max})
	}
}

func (s *Server) Sqrt(ctx context.Context, in *pb.SqrtRequest) (*pb.SqrtResponse, error) {
	log.Println("Invoked Sqrt on server")

	number := in.Number

	if number < 0 {
		return nil, status.Errorf(codes.InvalidArgument,
			"Invalid Argurment %d is less than 0", number)
	}

	return &pb.SqrtResponse{Result: math.Sqrt(float64(number))}, nil
}

func main() {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("Failed to listen on: %v\n", err)
	}

	log.Printf("Listening on %s\n", addr)

	// TLS Boilerplate grpc server options
	opts := []grpc.ServerOption{}
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"
	creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)

	if err != nil {
		log.Fatalf("Error when loading Cert/Key %v\n", err)
	}
	opts = append(opts, grpc.Creds(creds))

	s := grpc.NewServer(opts...)

	pb.RegisterCalculateServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}

}
