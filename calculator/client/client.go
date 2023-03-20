package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

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

func doMax(cc pb.CalculateServiceClient) {

	stream, err := cc.Max(context.Background())

	if err != nil {
		log.Fatalf("Failed to initaite Average request %v\n", err)
	}

	requests := []*pb.MaxRequest{
		{Number: 1},
		{Number: 5},
		{Number: 3},
		{Number: 6},
		{Number: 2},
		{Number: 20},
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer stream.CloseSend()
		defer wg.Done()
		for _, request := range requests {
			stream.Send(request)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Printf("Failed to recieve response from server%v\n", err)
			}
			fmt.Printf("Max Result %2.2f\n", res.Result)
		}
	}()

	wg.Wait()

}

func doSqrt(cc pb.CalculateServiceClient, number int32) {
	log.Println("New Sqrt calculation invoked")

	res, err := cc.Sqrt(context.Background(), &pb.SqrtRequest{Number: number})

	if err != nil {
		e, ok := status.FromError(err)

		if ok {
			// grpc error
			if e.Code() == codes.InvalidArgument {
				log.Printf("Send an invalid negative number\n")
			} else {
				log.Printf("gRPC error %v:%s", e.Code(), e.Message())
			}
			return
		}
		log.Fatalf("Failed to send Sqrt Request\n")
		return
	}

	log.Printf("Sqrt of (%d) = (%2.2f)\n", number, res.Result)

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
	doMax(client)
	doSqrt(client, 25)
	doSqrt(client, -25)

}
