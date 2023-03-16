package main

import (
	"context"
	"io"
	"log"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
)

func doGreet(c pb.GreetServiceClient) {
	log.Println("doGreet was invoked")

	res, err := c.Greet(context.Background(), &pb.GreetRequest{FirstName: "Bibo!"})

	if err != nil {
		log.Fatalf("Could not greet %v\n", err)
	}

	log.Printf("Greeting: %s\n", res.Result)
}

func doGreetStream(c pb.GreetServiceClient) {
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
