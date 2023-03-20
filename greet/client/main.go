package main

import (
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
)

var addr string = "localhost:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v\n", err)
	}

	defer conn.Close()

	c := pb.NewGreetServiceClient(conn)

	wg := sync.WaitGroup{}

	wg.Add(3)
	go doGreet(c, &wg)
	go doGreetStream(c, &wg)
	go doGreetLongStream(c, &wg)
	wg.Wait()

	doGreetAll(c)
	doGreetTimed(c, 5*time.Second)
	doGreetTimed(c, 1*time.Second)

}
