package main

import (
	"fmt"
	"log"
	"net"
	"os"

	serverpb "github.com/google/cel-go/server"
	cspb "github.com/google/cel-spec/proto/v1/cel_service"
	grpcpb "google.golang.org/grpc"
	reflectionpb "google.golang.org/grpc/reflection"
)

func main() {
	log.Println("Server opening listening port")
	lis, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Server opened port ", lis.Addr())

	// Must print to stdout, so the client can find the port.
	// So, no, this must be 'fmt', not 'log'.
	fmt.Printf("Listening on %v\n", lis.Addr())
	os.Stdout.Sync()
	log.Println("Server wrote address")

	log.Println("Server registering service on port")
	s := grpcpb.NewServer()
	cspb.RegisterCelServiceServer(s, &serverpb.CelServer{})
	log.Println("Server calling Register")
	reflectionpb.Register(s)
	log.Println("Server calling Serve")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
