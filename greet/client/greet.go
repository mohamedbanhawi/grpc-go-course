package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
)

func doGreet(c pb.GreetServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("doGreet was invoked")

	res, err := c.Greet(context.Background(), &pb.GreetRequest{FirstName: "Bibo!"})

	if err != nil {
		log.Fatalf("Could not greet %v\n", err)
	}

	log.Printf("Greeting: %s\n", res.Result)
}

func doGreetStream(c pb.GreetServiceClient, wg *sync.WaitGroup) {

	defer wg.Done()
	log.Println("doGreetStream was invoked")

	stream, err := c.GreetStream(context.Background(), &pb.GreetRequest{FirstName: "Bibo!"})

	if err != nil {
		log.Fatalf("Could not greet %v\n", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break

		} else if err != nil {
			log.Fatalf("Could not greet steam %v\n", err)
		}

		log.Printf("Greeting: %s\n", res.GetResult())
	}
}

func doGreetLongStream(c pb.GreetServiceClient, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Println("doGreetLongStream was invoked")

	stream, err := c.GreetLongStream(context.Background())

	if err != nil {
		log.Fatalf("Could not create greet long steam %v\n", err)

	}

	requests := []*pb.GreetRequest{
		{FirstName: "Bibo"},
		{FirstName: "Miso"},
		{FirstName: "Yoyo"},
		{FirstName: "Toti"},
	}
	for _, request := range requests {
		stream.Send(request)
		log.Printf("Sending %v\n", request)
		time.Sleep(time.Second)
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Could not recieve greet long stream %v\n", err)

	}

	log.Printf("Recieved: %s\n", response.Result)

}
