package main

import (
	"context"
	"io"
	"log"
	"time"

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

func doAverage(cc pb.CalculateServiceClient) {
	stream, err := cc.Average(context.Background())

	if err != nil {
		log.Fatalf("Failed to initaite Average request %v\n", err)
	}

	requests := []*pb.AverageRequest{
		{Number: 12.0},
		{Number: 2.0},
		{Number: 42.0},
	}

	for _, request := range requests {
		err := stream.Send(request)
		time.Sleep(time.Second)
		if err != nil {
			log.Fatalf("Failed to send Average request %v\n", err)
		}
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Failed to recieve Average response %v\n", err)
	}

	log.Printf("Average: %2.2f\n", response.Result)
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
	doAverage(client)

}
