package main

import (
	"log"
	"net"

	pb "github.com/mohamedbanhawi/grpc-go-course/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr string = "0.0.0.0:50051"

type Server struct {
	pb.GreetServiceServer
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

	pb.RegisterGreetServiceServer(s, &Server{})

	if err = s.Serve(lis); err != nil {
		log.Fatalf("Failed to server: %v\n", err)
	}

}
