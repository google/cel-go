package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/cel-go/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	log.Println("Server opening listening port")
	lis, err := net.Listen("tcp4", "127.0.0.1:")
	if err != nil {
		lis, err = net.Listen("tcp6", "[::1]:0")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
	}
	log.Println("Server opened port ", lis.Addr())

	// Must print to stdout, so the client can find the port.
	// So, no, this must be 'fmt', not 'log'.
	fmt.Printf("Listening on %v\n", lis.Addr())
	os.Stdout.Sync()
	log.Println("Server wrote address")

	log.Println("Server registering service on port")
	s := grpc.NewServer()
	exprpb.RegisterCelServiceServer(s, &server.CelServer{})
	log.Println("Server calling Register")
	reflection.Register(s)
	log.Println("Server calling Serve")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
