package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mohamedbanhawi/grpc-go-course/calculator/proto"
)

var addr string = "localhost:50051"

func doCalculate(cc pb.CalculateServiceClient) {

	response, err := cc.Calculate(context.Background(), &pb.CalculateRequest{
		FirstNumber: 10, SecondNumber: 2,
	})

	if err != nil {
		log.Fatalf("Failed to send request %v\n", err)
	}

	log.Printf("Invoked Calculate Request")

	log.Printf("Result: %2.2f\n", response.Result)
}

func main() {

	// connect to GRPC server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	client := pb.NewCalculateServiceClient(conn)

	doCalculate(client)

}
