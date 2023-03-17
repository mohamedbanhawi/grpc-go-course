package main

import (
	"context"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mohamedbanhawi/grpc-go-course/calculator/proto"
)

var addr string = "localhost:50051"

func doSum(cc pb.CalculateServiceClient) {

	response, err := cc.Sum(context.Background(), &pb.SumRequest{
		FirstNumber: 10, SecondNumber: 2,
	})

	if err != nil {
		log.Fatalf("Failed to send request %v\n", err)
	}

	log.Printf("Invoked Sum Request")

	log.Printf("Result: %2.2f\n", response.Result)

}

func doPrimes(cc pb.CalculateServiceClient) {
	stream, err := cc.Primes(context.Background(), &pb.PrimesRequest{Number: 120})

	if err != nil {
		log.Fatalf("Failed to send Primes request %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break

		} else if err != nil {
			log.Fatalf("Failed to recv primes steam %v\n", err)
		}

		log.Printf("Factor: %d\n", res.GetResult())
	}
}

func main() {

	// connect to GRPC server
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	client := pb.NewCalculateServiceClient(conn)

	doSum(client)
	doPrimes(client)

}
