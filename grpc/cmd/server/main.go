package main

import (
	"context"
	"fmt"
	"github.com/simonhammes/301-cloud-computing-project/grpc/search"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	search.SearchServiceServer
}

func (s *server) Search(_ context.Context, _ *search.SearchRequest) (*search.SearchResponse, error) {
	response := &search.SearchResponse{
		NumberOfResults: 42,
	}

	return response, nil
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := server{}

	search.RegisterSearchServiceServer(grpcServer, &server)

	println("Starting server...")
	fmt.Printf("Listening on %s\n", listener.Addr())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
