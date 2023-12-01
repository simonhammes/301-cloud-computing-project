package main

import (
	"context"
	"fmt"
	"github.com/simonhammes/301-cloud-computing-project/grpc/search"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
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

func (s *server) GeneratePrimes(r *search.PrimeGenerationRequest, stream search.SearchService_GeneratePrimesServer) error {
	primes := sieveOfEratosthenes(int(r.Limit))

	for _, prime := range primes {
		response := search.PrimeResponse{Prime: int32(prime)}
		if err := stream.Send(&response); err != nil {
			return err
		}

		// Sleep for 1 second
		time.Sleep(time.Second)
	}

	return nil
}

// Return list of primes smaller than N
func sieveOfEratosthenes(N int) (primes []int) {
	b := make([]bool, N)
	for i := 2; i < N; i++ {
		if b[i] == true {
			continue
		}
		primes = append(primes, i)
		for k := i * i; k < N; k += i {
			b[k] = true
		}
	}
	return
}

func (s *server) RecordLogMessages(stream search.SearchService_RecordLogMessagesServer) error {
	count := 0

	for {
		logMessage, err := stream.Recv()

		if err == io.EOF {
			summary := search.LogSummary{Count: int32(count)}
			return stream.SendAndClose(&summary)
		}

		if err != nil {
			return err
		}

		count++

		// Process message
		log.Printf("Message: %s", logMessage.Message)
	}
}

func (s *server) LogChat(stream search.SearchService_LogChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// No more incoming messages
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received message: %s", in.Message)

		// Do some work
		time.Sleep(500 * time.Millisecond)

		message := search.LogMessage{
			Message: fmt.Sprintf("Processed log message: %s", in.Message),
		}
		err = stream.Send(&message)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
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
