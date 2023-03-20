package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *Server) GreetLongStream(stream pb.GreetService_GreetLongStreamServer) error {
	log.Printf("GreetStream function invoked\n")

	var result string = ""

	for {
		request, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&pb.GreetResponse{Result: result})
			return nil
		} else if err != nil {
			log.Fatalf("Failed to recieve stream%v\n", err)
		}
		log.Printf("GreetStream function stream recieved%v\n", request)
		result += fmt.Sprintf("Hello, %s!\n", request.FirstName)
	}
}

func (s *Server) GreetAll(stream pb.GreetService_GreetAllServer) error {
	log.Printf("GreetAll function invoked\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			log.Fatalf("Failed to recieve stream%v\n", err)
		}
		log.Printf("GreetAll function stream recieved%v\n", req)
		stream.Send(&pb.GreetResponse{Result: req.FirstName})
	}

}

func (s *Server) GreetTimed(ctx context.Context, in *pb.GreetRequest) (*pb.GreetResponse, error) {
	log.Println("GreetTimed Invoked on Server")

	response := &pb.GreetResponse{Result: fmt.Sprintf("Hello, %s !\n", in.FirstName)}

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Client cancelled the request")
			return nil, status.Error(codes.Canceled, "Client cancelled the request")
		}
		time.Sleep(1 * time.Second)
	}
	return response, nil
}
